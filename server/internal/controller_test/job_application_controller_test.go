package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createJobApplication(router *gin.Engine, userID string, jobID string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	objID, _ := primitive.ObjectIDFromHex(userID)
	jobIDReal, _ := primitive.ObjectIDFromHex(jobID)
	body, _ := json.Marshal(rawJobApplication(objID, jobIDReal))
	req, _ := http.NewRequest("POST", "/apply/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w, req
}

func rawJobApplication(applicantID primitive.ObjectID, jobID primitive.ObjectID) map[string]any {
	now := time.Now().UTC().Truncate(time.Second)
	return map[string]any{
		"applicantID": applicantID,
		"jobID":       jobID,
		"status":      "waiting for approval",
		"createdAt":   now,
	}
}

func createJob(router *gin.Engine, r *regexp.Regexp) string {
	w := httptest.NewRecorder()

	body, _ := json.Marshal(rawJob("Job for Job Application Creation Test"))

	req, _ := http.NewRequest("POST", "/jobs/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	jobIDMatches := r.FindStringSubmatch(w.Body.String())
	jobID := jobIDMatches[1]
	return jobID
}

func deleteJobApplication(jobAppID string, router *gin.Engine) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/apply/"+jobAppID, nil)
	router.ServeHTTP(w, req)
	return w
}

// Test the GET /apply/query endpoint with valid data (200 status)
func TestQueryJobApplicationsSuccess(t *testing.T) {
	router := getTestRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/apply/", nil)
	router.ServeHTTP(w, req)

	t.Log(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test the GET /apply/query endpoint with applicantID query
func TestQueryJobApplicationsByApplicantID(t *testing.T) {
	router := getTestRouter()
	w := httptest.NewRecorder()

	// Test querying with a specific applicant ID
	req, _ := http.NewRequest("GET", fmt.Sprintf("/apply/query?applicantID=%s", globalApplicantID), nil)
	router.ServeHTTP(w, req)

	fmt.Println("Response body:", w.Body.String())

	// Should succeed
	assert.Equal(t, http.StatusOK, w.Code)
	t.Log(w.Body.String())

	// Response should be JSON array
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

// Test the GET /apply/query endpoint with invalid da
// Should get 400 (Bad request) error
func TestQueryJobApplicationsInvalidParam(t *testing.T) {
	router := getTestRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/apply/query?invalidParam=invalidValue", nil)
	router.ServeHTTP(w, req)

	fmt.Println("Response body:", w.Body.String())

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test querying job applications by applicant's status
func TestQueryJobApplicationsByStatus(t *testing.T) {
	router := getTestRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/apply/query?status=pending", nil)
	router.ServeHTTP(w, req)

	fmt.Println("Response body:", w.Body.String())

	// Should succeed
	assert.Equal(t, http.StatusOK, w.Code)

	// Response should be JSON array
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestCreateAndDeleteJobApplication(t *testing.T) {
	router := getTestRouter()
	r := regexp.MustCompile(`"InsertedID":"(.+)"`)
	userID := createUser(router, r, "Odindindindindun")
	jobID := createJob(router, r)

	w2, _ := createJobApplication(router, userID, jobID)

	assert.Equal(t, http.StatusCreated, w2.Code)
	assert.Contains(t, w2.Body.String(), "InsertedID")
	jobAppIDMatches := r.FindStringSubmatch(w2.Body.String())
	jobAppID := jobAppIDMatches[1]

	w3 := deleteJobApplication(jobAppID, router)
	assert.Equal(t, http.StatusOK, w3.Code)
}
