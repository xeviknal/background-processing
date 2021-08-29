package consumers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xeviknal/background-commons/database"
	"github.com/xeviknal/background-commons/models"
)

func TestJobCancellerConsumer_Perform(t *testing.T) {
	asserter := assert.New(t)
	db := database.GetDb()

	// Prepare the environment
	// Create a new Job first
	job := prepareCancellerEnvironment(t)

	// Call the task
	JobCancellerConsumer{}.Perform(map[string]interface{}{"id": job.Id})

	// Asserts
	obj, err := db.Get(models.Job{}, job.Id)
	updatedJob := obj.(*models.Job)
	asserter.NoError(err)

	// Job results
	asserter.Equal("CANCELLED!", *updatedJob.Status)

	// Timestamp modified
	asserter.NotNil(updatedJob.StartedAt)
	asserter.NotNil(updatedJob.FinishedAt)
	asserter.NotNil(updatedJob.CancelledAt)
}

func prepareCancellerEnvironment(t *testing.T) *models.Job {
	database.SetConnectionConfig("test", "test", "test")
	db := database.GetDb()

	job := models.NewJob(10)
	now := time.Now()
	job.CancelledAt = &now
	err := db.Insert(&job)
	if err != nil {
		t.Fatalf("error while creating test data: %v", err)
	}

	return &job
}
