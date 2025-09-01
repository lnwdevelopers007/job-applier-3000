package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/auth"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
)

// NewRouter returns new, default router.
func NewRouter() *gin.Engine {
	router := gin.Default()
	cfg := setUpCors()

	router.Use(cors.New(cfg))

	// --- Auth routes ---
	// These are handled by browser redirects, not JavaScript.
	// They DO NOT need CORS middleware.

	authGroup := router.Group("/auth")
	{
		authGroup.GET("/:provider", auth.Login)
		authGroup.GET("/:provider/callback", auth.OAuthCallback)
		authGroup.GET("/:provider/logout", auth.Logout)
	}

	// --- API routes ---
	// Create a group for all other API endpoints that WILL be called from your frontend JS.
	api := router.Group("/api")

	{
		// Move your API endpoints into this group
		api.GET("/jobs", GetController[schema.Job]("jobs").RetrieveAll)
		api.POST("/jobs", GetController[schema.Job]("jobs").Create)

		// custom query route
		jobCtrl := NewJobController()
		router.GET("/jobs/query", jobCtrl.QueryJobs())

		// health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})
	}

	return router
}

func setUpCors() cors.Config {
	cfg := cors.Config{
		AllowOrigins: []string{
			os.Getenv("FRONTEND"),
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cfg
}
