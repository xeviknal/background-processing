package publishers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xeviknal/background-commons/database"
	"github.com/xeviknal/background-commons/models"
)

func TestJobsTask_Publish(t *testing.T) {
	asserter := assert.New(t)
	db := database.GetDb()

	// Prepare the environment
	// Create a new Job first
	prepareEnvironment(t)

	count, err := db.SelectInt("SELECT count(*) FROM jobs WHERE queued_at IS NULL")
	asserter.NoError(err)
	asserter.Equal(count, int64(10))

	// Call the task
	JobsPublisher{}.Publish()

	count, err = db.SelectInt("SELECT count(*) FROM jobs WHERE queued_at IS NOT NULL")
	asserter.NoError(err)
	asserter.Equal(count, int64(10))

	var jobs []models.Job
	_, err = db.Select(&jobs, "SELECT * FROM jobs WHERE queued_at IS NOT NULL")
	asserter.NoError(err)

	for _, job := range jobs {
		asserter.NotNil(job.QueuedAt)
		asserter.Nil(job.StartedAt)
		asserter.Nil(job.FinishedAt)
	}
}

func prepareEnvironment(t *testing.T) {
	database.SetConnectionConfig("test", "test", "test")
	db := database.GetDb()

	for i := 0; i < 10; i++ {
		job := models.NewJob(10)
		err := db.Insert(&job)
		if err != nil {
			t.Fatalf("error while creating test data: %v", err)
		}
	}
}
