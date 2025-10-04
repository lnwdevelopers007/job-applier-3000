package schema

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func jobApplicationValidPayload() map[string]any {
	now := time.Now().UTC().Truncate(time.Second)
	return map[string]any{
		"id":          primitive.NewObjectID(),
		"applicantID": primitive.NewObjectID(),
		"jobID":       primitive.NewObjectID(),
		"status":      "waiting for approval",
		"createdAt":   now,
	}
}

func bindMockJobApplication(t *testing.T, payload map[string]any) (JobApplication, error) {
	return bindMockRequest[JobApplication](t, payload)
}

// should always pass
func TestValidJobApplication(t *testing.T) {
	payload := jobApplicationValidPayload()
	_, err := bindMockJobApplication(t, payload)
	assert.NoError(t, err)
}

// error when there's no applicant id
func TestJobApplicationMissingApplicationID(t *testing.T) {
	payload := jobApplicationValidPayload()
	delete(payload, "applicantID")
	_, err := bindMockJobApplication(t, payload)
	assert.Error(t, err)
}

// error when there's no job id
func TestJobApplicationMissingJobID(t *testing.T) {
	payload := jobApplicationValidPayload()
	delete(payload, "jobID")
	_, err := bindMockJobApplication(t, payload)
	assert.Error(t, err)
}

func TestJobApplicationMissingStatus(t *testing.T) {
	payload := jobApplicationValidPayload()
	delete(payload, "status")
	_, err := bindMockJobApplication(t, payload)
	assert.Error(t, err)
}

func TestJobApplicationWrongStatusDataType(t *testing.T) {
	payload := jobApplicationValidPayload()
	payload["status"] = 123
	_, err := bindMockJobApplication(t, payload)
	assert.Error(t, err)
}

func TestJobApplicationMissingCreatedAt(t *testing.T) {
	payload := jobApplicationValidPayload()
	delete(payload, "createdAt")
	_, err := bindMockJobApplication(t, payload)
	assert.NoError(t, err)
}
