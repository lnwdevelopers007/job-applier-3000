package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/controller"
	"github.com/stretchr/testify/assert"
)

func setUp() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := controller.NewRouter()
	return router
}

func TestGetJobs(t *testing.T) {
	router := setUp()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/jobs", nil)
	router.ServeHTTP(w, req)
	fmt.Println(w)
	assert.Equal(t, 200, w.Code)
}

func TestPostJob(t *testing.T) {
	router := setUp()
	w := httptest.NewRecorder()

	// Example JSON payload
	body := `{"title": "Software Engineer", "salary": 1000}`

	// Create POST request with JSON body
	req, _ := http.NewRequest("POST", "/api/jobs", strings.NewReader(body))
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
