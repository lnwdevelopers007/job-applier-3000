package schema

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- helpers ---

func fileValidPayload() map[string]any {
	data := []byte("hello world")
	return map[string]any{
		"parentID":      primitive.NewObjectID(),
		"parentColl":    "jobs",
		"content":       base64.StdEncoding.EncodeToString(data), // JSON string → []byte
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
