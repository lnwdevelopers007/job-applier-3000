package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	w1 := httptest.NewRecorder()

	body, _ := json.Marshal(rawUser("Odindindindindun"))

	req, _ := http.NewRequest("POST", "/users/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w1, req)
	r := regexp.MustCompile(`"InsertedID":"(.+)"`)
	userIDMatches := r.FindStringSubmatch(w1.Body.String())
	userID := userIDMatches[1]

	w2 := httptest.NewRecorder()
	objID, _ := primitive.ObjectIDFromHex(userID)
	body, _ = json.Marshal(rawJobApplication(objID))
	req, _ = http.NewRequest("POST", "/apply/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req)

	assert.Equal(t, http.StatusCreated, w2.Code)
	assert.Contains(t, w2.Body.String(), "InsertedID")
	jobAppIDMatches := r.FindStringSubmatch(w2.Body.String())
	jobAppID := jobAppIDMatches[1]

	w3 := httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/apply/"+jobAppID, nil)
	router.ServeHTTP(w3, req)
	assert.Equal(t, http.StatusOK, w3.Code)
}
