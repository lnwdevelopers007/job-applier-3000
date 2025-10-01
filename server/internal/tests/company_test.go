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

// Helper function to create a valid company payload
func rawCompany() map[string]any {
	return map[string]any{
		"name":        "Tech Innovations Inc.",
		"userID":      primitive.NewObjectID(),
		"aboutUs":     "We are a leading technology company specializing in innovative solutions.",
		"companyType": "Technology",
	}
}

// Test retrieving all companies
func TestRetrieveAllCompanies(t *testing.T) {
	router := setUp()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/company/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

// Test creating a valid company
func TestCreateValidCompany(t *testing.T) {
	router := setUp()
	w := httptest.NewRecorder()

	body, _ := json.Marshal(rawCompany())
	req, _ := http.NewRequest("POST", "/company/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	fmt.Println("Response body:", w.Body.String())

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "InsertedID")
}

// Test creating an invalid company (missing required fields)
func TestCreateInvalidCompany(t *testing.T) {
	router := setUp()
	w := httptest.NewRecorder()

	// Create invalid company payload with missing fields
	payload := map[string]any{
		"name": "Incomplete Company",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/company/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	fmt.Println("Response body:", w.Body.String())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

// Test retrieving a specific company by ID
func TestGetCompanyByID(t *testing.T) {
	router := setUp()

	// First create a company
	w := httptest.NewRecorder()
	body, _ := json.Marshal(rawCompany())
	req, _ := http.NewRequest("POST", "/company/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Extract the inserted ID
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	companyID := response["InsertedID"]

	// Retrieve the company by ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/company/%s", companyID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

// Test retrieving a company with invalid ID
func TestGetCompanyByInvalidID(t *testing.T) {
	router := setUp()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/company/invalid-id", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}