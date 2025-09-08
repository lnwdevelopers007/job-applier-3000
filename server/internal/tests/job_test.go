package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/controller"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setUp() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := controller.NewRouter()
	return router
}

func TestRetrieveAllJobs(t *testing.T) {
	router := setUp()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs/", nil)
	router.ServeHTTP(w, req)
	fmt.Println(w)
	assert.Equal(t, 200, w.Code)
}

func newJob() ([]byte, error) {
	now := time.Now().UTC().Truncate(time.Second)
	return json.Marshal(map[string]any{
		"title":        "Brr Brr Engineer",
		"companyID":    primitive.NewObjectID(),
		"location":     "Millenium Science School, Kivotos.",
		"salary":       120000,
		"salaryRate":   "yearly",
		"workType":     "onsite",
		"contractType": "full-time",
		// privacyPolicy is optional
		"publicationInfo": map[string]any{
			"isHiring":  true,
			"createdAt": now.Format(time.RFC3339),
			"startDate": now.Format(time.RFC3339),
			"endDate":   now.AddDate(0, 1, 0).Format(time.RFC3339),
		},
		"criteria": map[string]any{
			"requirements":   []string{"Go", "MongoDB"},
			"qualifications": []string{"Bachelor's Degree"},
			// commonQuestions is optional
		},
		"isApproved": true, // NOTE: with `binding:"required"` on a bool, false will fail validation
	})
}

func TestCreateJob(t *testing.T) {
	router := setUp()
	w := httptest.NewRecorder()

	// Example JSON payload
	body, _ := newJob()

	// Create POST request with JSON body
	req, _ := http.NewRequest("POST", "/jobs/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	router.ServeHTTP(w, req)

	// Debug output (optional)
	fmt.Println("Response code:", w.Code)
	fmt.Println("Response body:", w.Body.String())

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code) // or StatusCreated if you return 201
	assert.Contains(t, w.Body.String(), "InsertedID")
}
