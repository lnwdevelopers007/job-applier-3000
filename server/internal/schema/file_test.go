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
		"parentID":      primitive.NewObjectID(),
		"parentColl":    "job_seekers",
		"filename":      "resume.pdf",
		"content":       base64.StdEncoding.EncodeToString(data),
		"contentType":   "application/pdf",
		"fileExtension": "pdf",
		"size":          int64(len(data)),
		"category":      "resume",
		"uploadDate":    time.Now(),
	}
}

func bindMockFile(t *testing.T, payload map[string]any) (File, error) {
	return bindMockRequest[File](t, payload)
}

// --- valid case ---

func TestValidFile(t *testing.T) {
	payload := fileValidPayload()
	file, err := bindMockFile(t, payload)
	assert.NoError(t, err)
	
	// Test validation
	err = file.Validate()
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

// --- validation tests ---

func TestFileValidationFileTooLarge(t *testing.T) {
	payload := fileValidPayload()
	payload["size"] = int64(11 * 1024 * 1024) // 11MB
	
	file, err := bindMockFile(t, payload)
	assert.NoError(t, err)
	
	err = file.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum")
}

func TestFileValidationInvalidContentType(t *testing.T) {
	payload := fileValidPayload()
	payload["contentType"] = "image/png"
	
	file, err := bindMockFile(t, payload)
	assert.NoError(t, err)
	
	err = file.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PDF")
}

func TestFileValidationInvalidCategory(t *testing.T) {
	payload := fileValidPayload()
	payload["category"] = "invalid_category"
	
	file, err := bindMockFile(t, payload)
	assert.NoError(t, err)
	
	err = file.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid file category")
}

func TestFileValidationInvalidParentColl(t *testing.T) {
	payload := fileValidPayload()
	payload["parentColl"] = "invalid_collection"
	
	file, err := bindMockFile(t, payload)
	assert.NoError(t, err)
	
	err = file.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid parent collection")
}
