package consumers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xeviknal/background-commons/database"
	"github.com/xeviknal/background-commons/models"
)

func TestJobsConsumer_Perform(t *testing.T) {
	asserter := assert.New(t)
	db := database.GetDb()

	// Prepare the environment
	// Create a new Job first
	job := prepareEnvironment(t)

	// Call the task
	JobsConsumer{}.Perform(map[string]interface{}{"id": job.Id})

	// Asserts
	obj, err := db.Get(models.Job{}, job.Id)
	updatedJob := obj.(*models.Job)
	asserter.NoError(err)

	// Job results
	asserter.GreaterOrEqual(3*time.Second, *updatedJob.Sleep)
	asserter.Equal("DONE!", *updatedJob.Status)

	// Timestamp modified
	asserter.NotNil(updatedJob.QueuedAt)
	asserter.NotNil(updatedJob.StartedAt)
	asserter.NotNil(updatedJob.FinishedAt)
}

func prepareEnvironment(t *testing.T) *models.Job {
	database.SetConnectionConfig("test", "test", "test")
	db := database.GetDb()
	maxSleepDuration = 3
	minSleepDuration = 1

	job := models.NewJob(10)
	now := time.Now()
	job.QueuedAt = &now
	err := db.Insert(&job)
	if err != nil {
		t.Fatalf("error while creating test data: %v", err)
	}

	return &job
}
