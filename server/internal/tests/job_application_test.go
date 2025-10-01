package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/controller"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setUpJobApplicationTests() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := controller.NewRouter()
	return router
}

// Test the GET /apply/ endpoint which calls the Query function
func TestGetJobApplications(t *testing.T) {
	router := setUpJobApplicationTests()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/apply/", nil)
	router.ServeHTTP(w, req)

	fmt.Println("Response body:", w.Body.String())

	assert.Equal(t, 200, w.Code)
}

// Test the GET /apply/query endpoint with valid data (200 status)
func TestQueryJobApplicationsSuccess(t *testing.T) {
	router := setUpJobApplicationTests()
	w := httptest.NewRecorder()

	// Test querying with a specific status that should return data
	req, _ := http.NewRequest("GET", "/apply/", nil)
	router.ServeHTTP(w, req)

	fmt.Println("Response body:", w.Body.String())

	// Should succeed
	assert.Equal(t, 200, w.Code)

	// Response should be JSON array
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

// Test the GET /apply/query endpoint with applicantID query
func TestQueryJobApplicationsByApplicantID(t *testing.T) {
	router := setUpJobApplicationTests()
	w := httptest.NewRecorder()

	// Test querying with a specific applicant ID
	applicantID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/apply/query?applicantID=%s", applicantID), nil)
	router.ServeHTTP(w, req)

	fmt.Println("Response body:", w.Body.String())

	// Should succeed
	assert.Equal(t, http.StatusOK, w.Code)

	// Response should be JSON array
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

// Test the GET /apply/query endpoint with invalid da
// Should get 400 (Bad request) error
func TestQueryJobApplicationsServerError(t *testing.T) {
	router := setUpJobApplicationTests()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/apply/query?invalidParam=invalidValue", nil)
	router.ServeHTTP(w, req)

	fmt.Println("Response body:", w.Body.String())

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test querying job applications with applicant info
func TestQueryJobApplicationsWithApplicantInfo(t *testing.T) {
	router := setUpJobApplicationTests()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/apply/query?status=pending", nil)
	router.ServeHTTP(w, req)

	fmt.Println("Response body:", w.Body.String())

	// Should succeed
	assert.Equal(t, http.StatusOK, w.Code)

	// Response should be JSON array
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}
