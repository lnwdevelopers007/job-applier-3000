package schema

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- helpers ---

func noteValidPayload() map[string]any {
	now := time.Now().UTC().Truncate(time.Second)
	return map[string]any{
		"id":               primitive.NewObjectID(),
		"jobApplicationID": primitive.NewObjectID(),
		"content":          "This is a test note",
		"timestamp":        now,
	}
}

func bindMockNote(t *testing.T, payload map[string]any) (Note, error) {
	return bindMockRequest[Note](t, payload)
}

// --- valid case ---

func TestValidNote(t *testing.T) {
	payload := noteValidPayload()
	_, err := bindMockNote(t, payload)
	assert.NoError(t, err)
}

// --- JobApplicationID ---

func TestMissingJobApplicationID(t *testing.T) {
	payload := noteValidPayload()
	delete(payload, "jobApplicationID")
	_, err := bindMockNote(t, payload)
	assert.Error(t, err)
}

// --- Content ---

func TestMissingContent(t *testing.T) {
	payload := noteValidPayload()
	delete(payload, "content")
	_, err := bindMockNote(t, payload)
	assert.Error(t, err)
}

// --- Timestamp ---

func TestMissingTimestamp(t *testing.T) {
	payload := noteValidPayload()
	delete(payload, "timestamp")
	_, err := bindMockNote(t, payload)
	assert.Error(t, err)
}