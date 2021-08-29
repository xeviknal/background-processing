package publishers

import (
	"fmt"
	"log"
	"time"

	"github.com/xeviknal/background-commons/database"
	"github.com/xeviknal/background-commons/models"
	"github.com/xeviknal/background-processing/buffer"
)

type JobCancellerPublisher struct{}
type JobCancellerConsumer struct{}

// Checking that JobCancellerPublisher implements
//  Publisher interface at compile time
var _ Publisher = (*JobCancellerPublisher)(nil)

var defaultTimeout = 10

func (jc JobCancellerPublisher) Publish() {
	db := database.GetDb()

	// Starting a Db transaction
	trans, err := db.Begin()
	if err != nil {
		log.Fatalf("error when starting a transaction for %T: %v", jc, err)
	}

	// Fetching the jobs that exceeded maximum execution time allowed
	var jobs []models.Job
	_, err = trans.Select(&jobs, fmt.Sprintf("SELECT * FROM jobs WHERE started_at IS NOT NULL "+
		"AND TIMEDIFF(NOW(), started_at) > SEC_TO_TIME(%d) AND cancelled_at IS NULL AND finished_at IS NULL FOR UPDATE", defaultTimeout))
	if err != nil {
		log.Fatalf("error querying for %T: %v", jc, err)
	}

	queue := buffer.GetBuffer()
	for _, job := range jobs {
		descriptor := buffer.Descriptor(buffer.TaskDescriptor{
			Consumer:   "JobCancellerConsumer",
			Descriptor: map[string]interface{}{"id": job.Id},
		})
		queue.Add(&descriptor)

		// Bookkeeping: job cancellation time
		now := time.Now()
		job.CancelledAt = &now
		trans.Update(&job)
	}
	trans.Commit()
}
