package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrieveAllJobs(t *testing.T) {
	router := getTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// t.Log(w.Body.String())
	assert.GreaterOrEqual(t, w.Body.Len(), 1)
}

func TestCreateValidJob(t *testing.T) {
	router := getTestRouter()
	w := httptest.NewRecorder()

	// Example JSON payload
	body, _ := json.Marshal(rawJob("TestCreateValidJob"))

	// Create POST request with JSON body
	req, _ := http.NewRequest("POST", "/jobs/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	router.ServeHTTP(w, req)

	// Debug output (optional)
	// t.Log("Response body:", w.Body.String())

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code) // or StatusCreated if you return 201
	assert.Contains(t, w.Body.String(), "InsertedID")
}

func TestCreateJobWithWrongCompanyID(t *testing.T) {
	router := getTestRouter()
	w := httptest.NewRecorder()

	// Example JSON payload
	raw := rawJob("TestCreateJobWithWrongCompanyID")
	raw["companyID"] = "64f3a2b7e1d3a8c1b0f9d2xx"
	body, _ := json.Marshal(raw)

	// Create POST request with JSON body
	req, _ := http.NewRequest("POST", "/jobs/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	router.ServeHTTP(w, req)

	// t.Log("Response body:", w.Body.String())

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code) // or StatusCreated if you return 201
}

func TestUpdateJob(t *testing.T) {
	router := getTestRouter()
	w1 := httptest.NewRecorder()

	// Example JSON payload
	raw := rawJob("Job Title Before Update")
	body, _ := json.Marshal(raw)

	// Create POST request with JSON body
	req, _ := http.NewRequest("POST", "/jobs/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	router.ServeHTTP(w1, req)

	// Debug output (optional)
	r := regexp.MustCompile(`"InsertedID":"(.+)"`)
	matches := r.FindStringSubmatch(w1.Body.String())
	id := matches[1]

	// t.Log(id)

	w2 := httptest.NewRecorder()

	raw["title"] = "Job Title After Update"
	body, _ = json.Marshal(raw)

	// Create POST request with JSON body
	req, _ = http.NewRequest("PUT", "/jobs/"+id, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	router.ServeHTTP(w2, req)
	assert.Contains(t, w2.Body.String(), "updated")
	assert.Equal(t, w2.Code, http.StatusOK)
	// t.Log(w2.Body.String())

}

func TestDeleteJob(t *testing.T) {
	router := getTestRouter()
	w1 := httptest.NewRecorder()

	// Example JSON payload
	raw := rawJob("You should not see this in mongo collection!")
	body, _ := json.Marshal(raw)

	// Create POST request with JSON body
	req, _ := http.NewRequest("POST", "/jobs/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	router.ServeHTTP(w1, req)

	// Debug output (optional)
	r := regexp.MustCompile(`"InsertedID":"(.+)"`)
	matches := r.FindStringSubmatch(w1.Body.String())
	id := matches[1]

	// Create POST request with JSON body
	w2 := httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/jobs/"+id, bytes.NewReader([]byte(`{"reason": "testing delete"}`)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w2, req)
	t.Log(w2.Body)
	assert.Equal(t, w2.Code, http.StatusOK)
}
