package schema

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- helpers ---

func fileValidPayload() map[string]any {
	data := []byte("hello world")
	return map[string]any{
		"userID":        primitive.NewObjectID(),
		"content":       base64.StdEncoding.EncodeToString(data),
		"fileExtension": "pdf",
		"filename":      "resume.pdf",
		"contentType":   "application/pdf",
		"size":          int64(len(data)),
		"category":      "resume",
		"uploadDate":    time.Now(),
	}
}

func bindMockFile(t *testing.T, payload map[string]any) (File, error) {
	return bindMockRequest[File](t, payload)
}

// --- valid binding cases ---

func TestValidFile(t *testing.T) {
	payload := fileValidPayload()
	_, err := bindMockFile(t, payload)
	assert.NoError(t, err)
}

// --- missing required fields ---

func TestFileMissingUserID(t *testing.T) {
	payload := fileValidPayload()
	delete(payload, "userID")

	_, err := bindMockFile(t, payload)
	assert.Error(t, err)
}

func TestFileMissingFileExtension(t *testing.T) {
	payload := fileValidPayload()
	delete(payload, "fileExtension")

	_, err := bindMockFile(t, payload)
	assert.Error(t, err)
}

func TestFileMissingFilename(t *testing.T) {
	payload := fileValidPayload()
	delete(payload, "filename")

	_, err := bindMockFile(t, payload)
	assert.Error(t, err)
}

func TestFileMissingContentType(t *testing.T) {
	payload := fileValidPayload()
	delete(payload, "contentType")

	_, err := bindMockFile(t, payload)
	assert.Error(t, err)
}

func TestFileMissingSize(t *testing.T) {
	payload := fileValidPayload()
	delete(payload, "size")

	_, err := bindMockFile(t, payload)
	assert.Error(t, err)
}

func TestFileMissingCategory(t *testing.T) {
	payload := fileValidPayload()
	delete(payload, "category")

	_, err := bindMockFile(t, payload)
	assert.Error(t, err)
}