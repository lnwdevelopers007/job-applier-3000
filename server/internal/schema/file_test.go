package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- helpers ---

func fileValidPayload() map[string]any {
	return map[string]any{
		"parentID":      primitive.NewObjectID(),
		"parentColl":    "jobs",
		"content":       primitive.Binary{Subtype: 0x00, Data: []byte("hello world")},
		"fileExtension": "pdf",
	}
}

func bindMockFile(t *testing.T, payload map[string]any) (File, error) {
	return bindMockRequest[File](t, payload)
}

// --- valid case ---

func TestValidFile(t *testing.T) {
	payload := fileValidPayload()
	_, err := bindMockFile(t, payload)
	assert.NoError(t, err)
}

// --- invalid cases ---

func TestFileMissingParentID(t *testing.T) {
	payload := fileValidPayload()
	delete(payload, "parentID")

	_, err := bindMockFile(t, payload)
	assert.Error(t, err)
}

func TestFileMissingParentColl(t *testing.T) {
	payload := fileValidPayload()
	delete(payload, "parentColl")

	_, err := bindMockFile(t, payload)
	assert.Error(t, err)
}

func TestFileInvalidFileMissingContent(t *testing.T) {
	payload := fileValidPayload()
	delete(payload, "content")

	_, err := bindMockFile(t, payload)
	assert.Error(t, err)
}

func TestFileMissingFileExtension(t *testing.T) {
	payload := fileValidPayload()
	delete(payload, "fileExtension")

	_, err := bindMockFile(t, payload)
	assert.Error(t, err)
}
