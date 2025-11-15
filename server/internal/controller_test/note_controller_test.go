package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func rawNote(jobApplicationID primitive.ObjectID, content string) map[string]any {
	now := time.Now().UTC().Truncate(time.Second)
	return map[string]any{
		"jobApplicationID": jobApplicationID,
		"content":          content,
		"timestamp":        now,
	}
}

func createNote(router *gin.Engine, r *regexp.Regexp, jobApplicationID primitive.ObjectID, content string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()

	body, _ := json.Marshal(rawNote(jobApplicationID, content))

	req, _ := http.NewRequest("POST", "/notes/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	return w, req
}

func TestSavingValidNote(t *testing.T) {
	router := getTestRouter()

	r := regexp.MustCompile(`"InsertedID":"(.+)"`)

	jobApplier := createUser(router, r, "ibuki")
	recruiter := createUser(router, r, "suspicous van")

	job := createJob(router, r, recruiter)

	w, _ := createJobApplication(router, jobApplier, job)

	jobApplication := r.FindStringSubmatch(w.Body.String())[1]

	jobApplicationObjID, _ := primitive.ObjectIDFromHex(jobApplication)

	w2, _ := createNote(router, r, jobApplicationObjID, "declicious")
	assert.Equal(t, http.StatusCreated, w2.Code)
}
