package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRetrieveAllJobSeeker(t *testing.T) {
	router := setUp()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/seeker/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func rawJobSeeker() map[string]any {
	return map[string]any{
		"userID": primitive.NewObjectID(),
		"contact": map[string]any{
			"location": "Bangkok, Thailand",
			"phone":    "+66123456789",
			"linkedIn": "https://linkedin.com/in/johndoe",
		},
	}
}

func TestCreateInvalidJobSeeker(t *testing.T) {
	router := setUp()
	w := httptest.NewRecorder()

	// Example JSON payload
	body, _ := json.Marshal(rawJob())

	// Create POST request with JSON body
	req, _ := http.NewRequest("POST", "/seeker/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	router.ServeHTTP(w, req)

	// Debug output (optional)
	fmt.Println("Response body:", w.Body.String())

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestCreateValidJobSeeker(t *testing.T) {
	router := setUp()
	w := httptest.NewRecorder()

	// Example JSON payload
	body, _ := json.Marshal(rawJobSeeker())

	// Create POST request with JSON body
	req, _ := http.NewRequest("POST", "/seeker/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	router.ServeHTTP(w, req)

	// Debug output (optional)
	fmt.Println("Response body:", w.Body.String())

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}
