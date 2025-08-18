package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
	req, _ := http.NewRequest("GET", "/jobs", nil)
	router.ServeHTTP(w, req)
	fmt.Println(w)
	assert.Equal(t, 200, w.Code)
}
