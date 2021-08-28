package publishers

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/xeviknal/background-commons/database"
	"github.com/xeviknal/background-commons/models"
)

func TestJobsCanceller_Publish(t *testing.T) {
	asserter := assert.New(t)
	db := database.GetDb()

	// Prepare the environment
	// Create a new Job first
	prepareCancellerEnvironment(t, false)

	count, err := db.SelectInt(fmt.Sprintf("SELECT count(*) FROM jobs WHERE started_at IS NOT NULL "+
		"AND TIMEDIFF(NOW(), started_at) > %d", defaultTimeout))
	asserter.NoError(err)
	asserter.Equal(count, int64(10))

	// Call the task
	JobCancellerPublisher{}.Publish()

	count, err = db.SelectInt("SELECT count(*) FROM jobs WHERE cancelled_at IS NOT NULL")
	asserter.NoError(err)
	asserter.Equal(count, int64(10))

	var jobs []models.Job
	_, err = db.Select(&jobs, "SELECT * FROM jobs WHERE cancelled_at IS NOT NULL")
	asserter.NoError(err)

	for _, job := range jobs {
		asserter.NotNil(job.CancelledAt)
		asserter.NotNil(job.StartedAt)
		asserter.NotNil(job.QueuedAt)
		asserter.Nil(job.FinishedAt)
	}
}

func TestJobsCanceller_PublishAlreadyCancelled(t *testing.T) {
	asserter := assert.New(t)
	db := database.GetDb()

	// Prepare the environment
	// Create a new Job first
	prepareCancellerEnvironment(t, true)

	count, err := db.SelectInt(fmt.Sprintf("SELECT count(*) FROM jobs WHERE started_at IS NOT NULL "+
		"AND TIMEDIFF(NOW(), started_at) > %d", defaultTimeout))
	asserter.NoError(err)
	asserter.Equal(count, int64(10))

	// Call the task
	JobCancellerPublisher{}.Publish()

	var jobs []models.Job
	_, err = db.Select(&jobs, "SELECT * FROM jobs WHERE cancelled_at IS NOT NULL")
	asserter.NoError(err)

	for _, job := range jobs {
		// CancelledAt is not modified twice (implying that the publisher is run twice)
		asserter.Equal(job.QueuedAt, job.CancelledAt)
	}
}

func prepareCancellerEnvironment(t *testing.T, cancelled bool) {
	database.SetConnectionConfig("test", "test", "test")
	db := database.GetDb()

	defaultTimeout = 1 // 1 second

	for i := 0; i < 10; i++ {
		job := models.NewJob(10)
		now := time.Now()
		job.QueuedAt = &now
		job.StartedAt = &now
		if cancelled {
			job.CancelledAt = &now
		}
		err := db.Insert(&job)
		if err != nil {
			t.Fatalf("error while creating test data: %v", err)
		}
	}

	time.Sleep(1 * time.Second)
}
