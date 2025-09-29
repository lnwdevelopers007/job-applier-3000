package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- helpers ---

func jobSeekerValidPayload() map[string]any {
	return map[string]any{
		"id":     primitive.NewObjectID(),
		"userID": primitive.NewObjectID(),
		"contact": map[string]any{
			"location": "Bangkok, Thailand",
			"phone":    "+66123456789",
			"linkedIn": "https://linkedin.com/in/johndoe",
		},
	}
}

func bindMockJobSeeker(t *testing.T, payload map[string]any) (JobSeeker, error) {
	return bindMockRequest[JobSeeker](t, payload)
}

// --- valid case ---

func TestValidJobSeeker(t *testing.T) {
	payload := jobSeekerValidPayload()
	_, err := bindMockJobSeeker(t, payload)
	assert.NoError(t, err)
}

// --- UserID ---

func TestJobSeekerMissingUserIDInJobSeeker(t *testing.T) {
	payload := jobSeekerValidPayload()
	delete(payload, "userID")
	_, err := bindMockJobSeeker(t, payload)
	assert.Error(t, err)
}

// --- Contact ---

func TestJobSeekerMissingContactInJobSeeker(t *testing.T) {
	payload := jobSeekerValidPayload()
	delete(payload, "contact")
	_, err := bindMockJobSeeker(t, payload)
	assert.Error(t, err)
}

func TestJobSeekerMissingLocationInContact(t *testing.T) {
	payload := jobSeekerValidPayload()
	contact := payload["contact"].(map[string]any)
	delete(contact, "location")
	_, err := bindMockJobSeeker(t, payload)
	assert.Error(t, err)
}

func TestJobSeekerMissingPhoneInContact(t *testing.T) {
	payload := jobSeekerValidPayload()
	contact := payload["contact"].(map[string]any)
	delete(contact, "phone")
	_, err := bindMockJobSeeker(t, payload)
	assert.Error(t, err)
}

func TestJobSeekerMissingLinkedInInContact(t *testing.T) {
	payload := jobSeekerValidPayload()
	contact := payload["contact"].(map[string]any)
	delete(contact, "linkedIn")
	_, err := bindMockJobSeeker(t, payload)
	assert.NoError(t, err)
}
