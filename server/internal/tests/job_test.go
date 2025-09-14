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
		// basic info
		"title":           "Brr Brr Engineer",
		"companyID":       "64f3a2b7e1d3a8c1b0f9d2a1",
		"location":        "Millenium Science School, Kivotos.",
		"workType":        "onsite",
		"workArrangement": "full-time",
		"currency":        "THB",
		"minSalary":       2000.34,
		"maxSalary":       300000.213213,

		// description
		"jobDescription": "long",
		"jobSummary":     "longer",

		// requirements
		"requiredSkills":  "none",
		"experienceLevel": "a lot",
		"education":       "maybe",
		"niceToHave":      "noting",

		// post settings
		"applicationDeadline": now,
		"numberOfPositions":   1,
		"visibility":          "public",
		"emailNotifications":  true,
		"autoReject":          false,
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
