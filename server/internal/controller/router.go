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

	authGroup := router.Group("/auth")
	{
		authGroup.GET("/:provider", auth.Login)
		authGroup.GET("/:provider/callback", auth.OAuthCallback)
		authGroup.GET("/:provider/logout", auth.Logout)
	}

	jobs := router.Group("/jobs")
	jobCtrl := NewJobController()
	{
		jobs.GET("/", GetController[schema.Job]("jobs").RetrieveAll)
		jobs.POST("/", GetController[schema.Job]("jobs").Create)
		jobs.PUT("/:id", GetController[schema.Job]("jobs").Update)
		jobs.DELETE("/:id", GetController[schema.Job]("jobs").Delete)

		jobs.GET("/query", jobCtrl.QueryJobs())

	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

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
