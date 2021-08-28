package consumers

import (
	"log"
	"math/rand"
	"time"

	"gopkg.in/gorp.v1"

	"github.com/xeviknal/background-commons/database"
	"github.com/xeviknal/background-commons/models"
)

type JobsConsumer struct {
	trans *gorp.Transaction
	job   *models.Job
}

// Checking that JobsConsumer implements
//  Consumer interface at compile time
var _ Consumer = (*JobsConsumer)(nil)

var maxSleepDuration = 40
var minSleepDuration = 15

func (jc JobsConsumer) Perform(params map[string]interface{}) {
	jobId, found := params["id"]
	if !found {
		log.Fatalf("wrong parameters for %T: id not found", jc)
		return
	}

	// Before starting the task itself, notify that the task has been started
	// Using transaction to block the job row in database and avoid
	// multiple executions of the jobs
	if err := jc.StartTransaction(); err != nil {
		// TODO: retry? Set to QUEUE and increment retries
		log.Fatalf("error starting a transaction for %T: %v", jc, err)
		return
	}

	if err := jc.RetrieveJob(jobId.(int64)); err != nil {
		// TODO: retry? Set to QUEUE and increment retries
		log.Fatalf("error retrieveing info for %T job with id: %d. error: %v", jc, jobId, err)
		return
	}

	// TODO: Validations
	// Only one task for object_id every 5 minutes
	// Fetch last Job with object_id = job.objectId and check time difference
	// If too close, abort mission

	// Notifying task start
	if err := jc.NotifyStartedAt(); err != nil {
		// TODO: retry? Set to QUEUE and increment retries
		log.Fatalf("error notifying start for job %T with id: %d. error: %v", jc, jobId, err)
		return
	}
	// Committing the reception of the job
	if err := jc.FinishTransaction(); err != nil {
		// TODO: retry? Set to QUEUE and increment retries
		log.Fatalf("error finishing transactio for job %T with id: %d. error: %v", jc, jobId, err)
		return
	}
	// Finishing transaction. Releasing the lock. Allowing publishers to see
	// this job is already being processed.

	// Starting the task itself. It doesn't need any lock.
	jc.PerformTask()

	// Once the task is done. Notify the end.
	// Done in a transaction not to race condition with other publishers.
	if err := jc.StartTransaction(); err != nil {
		// TODO: retry? Set to QUEUE and increment retries
		log.Fatalf("error starting transaction for job %T with id: %d. error: %v", jc, jobId, err)
		return
	}

	// Notify the end of the task. To allow publishers know this job is done
	if err := jc.NotifyEndTask(); err != nil {
		// TODO: retry? Set to QUEUE and increment retries
		log.Fatalf("error updating job %T with id: %d. error: %v", jc, jobId, err)
		return
	}

	if err := jc.FinishTransaction(); err != nil {
		// TODO: retry? Set to QUEUE and increment retries
		log.Fatalf("error finishing transaction for job %T with id: %d. error: %v", jc, jobId, err)
		return
	}
}

func sleepDuration() time.Duration {
	return time.Duration(rand.Intn(maxSleepDuration-minSleepDuration)+minSleepDuration) * time.Second
}

func (jc *JobsConsumer) StartTransaction() error {
	db := database.GetDb()

	// Starting a Db transaction
	trans, err := db.Begin()
	if err != nil {
		log.Fatalf("error when starting a transaction for %T: %v", jc, err)
		return err
	}
	jc.trans = trans
	return nil
}

func (jc *JobsConsumer) FinishTransaction() error {
	return jc.trans.Commit()
}

func (jc *JobsConsumer) RetrieveJob(jobId int64) error {
	jc.job = &models.Job{}
	err := jc.trans.SelectOne(jc.job, "SELECT * FROM jobs WHERE id = ? FOR UPDATE", jobId)
	if err != nil {
		log.Fatalf("error retrieving job #%d info for %T: %v", jobId, jc, err)
		return err
	}
	return nil
}

func (jc *JobsConsumer) NotifyStartedAt() error {
	now := time.Now()
	jc.job.StartedAt = &now
	_, err := jc.trans.Update(jc.job)
	return err
}

func (jc *JobsConsumer) PerformTask() {
	dur := sleepDuration()
	time.Sleep(dur)

	jc.job.Sleep = &dur
	status := "DONE!"
	jc.job.Status = &status
}

func (jc *JobsConsumer) NotifyEndTask() error {
	now := time.Now()
	jc.job.FinishedAt = &now
	_, err := jc.trans.Update(jc.job)
	return err
}
