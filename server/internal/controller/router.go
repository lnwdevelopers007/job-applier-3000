package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
)

// NewRouter returns new, default router.
func NewRouter() *gin.Engine {
	router := gin.Default()

	// generic CRUD
	router.GET("/jobs", GetController[schema.JobSchema]("jobs").RetrieveAll())
	router.POST("/jobs", GetController[schema.JobSchema]("jobs").Create())

	// custom query route
	jobCtrl := NewJobController()
	router.GET("/jobs/query", jobCtrl.QueryJobs())

	// health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	return router
}
