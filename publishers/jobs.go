package publishers

import (
	"log"
	"time"

	"github.com/xeviknal/background-commons/database"
	"github.com/xeviknal/background-commons/models"
	"github.com/xeviknal/background-processing/buffer"
)

type JobsPublisher struct{}
type JobsConsumer struct{}

// Checking that JobsPublisher implements
//  Publisher interface at compile time
var _ Publisher = (*JobsPublisher)(nil)

func (t JobsPublisher) Publish() {
	db := database.GetDb()

	// Starting a Db transaction
	trans, err := db.Begin()
	if err != nil {
		log.Fatalf("error when starting a transaction for %T: %v", t, err)
	}

	// Fetching the jobs that haven't been queued
	var jobs []models.Job
	_, err = trans.Select(&jobs, "SELECT * FROM jobs WHERE queued_at IS NULL FOR UPDATE")
	if err != nil {
		log.Fatalf("error querying for %T: %v", t, err)
	}

	queue := buffer.GetBuffer()
	for _, job := range jobs {
		descriptor := buffer.Descriptor(buffer.TaskDescriptor{
			Consumer:   JobsConsumer{},
			Descriptor: job.Id,
		})
		queue.Add(&descriptor)

		// Bookkeeping: job queueing time
		now := time.Now()
		job.QueuedAt = &now
		trans.Update(&job)
	}
	trans.Commit()
}
