package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/auth"
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
		authGroup.POST("/refresh", auth.RefreshToken)
	}

	jobs := router.Group("/jobs")
	jobCtrl := NewJobController()
	{
		jobs.GET("/query", jobCtrl.Query)
		jobs.GET("/", jobCtrl.RetrieveAll)
		jobs.POST("/", jobCtrl.Create)
		jobs.PUT("/:id", jobCtrl.Update)
		jobs.DELETE("/:id", jobCtrl.Delete)
		jobs.GET("/:id", jobCtrl.RetrieveOne)
	}

	applicationController := NewJobApplicationController()
	applyRoutes := router.Group("/apply")
	{
		applyRoutes.GET("/query", applicationController.Query)
		applyRoutes.GET("/", applicationController.RetrieveAll)
		applyRoutes.POST("/", applicationController.Create)
		applyRoutes.PUT("/:id", applicationController.Update)
		applyRoutes.DELETE("/:id", applicationController.Delete)
		applyRoutes.GET("/:id", applicationController.RetrieveOne)
	}

	jobSeeker := NewJobSeekerController()
	jobSeekerRoutes := router.Group("/seeker")
	{
		jobSeekerRoutes.GET("/", jobSeeker.RetrieveAll)
		jobSeekerRoutes.POST("/", jobSeeker.Create)
    jobSeekerRoutes.PUT("/:id", jobSeeker.Update)
		jobSeekerRoutes.DELETE("/:id", jobSeeker.Delete)
		jobSeekerRoutes.GET("/:id", jobSeeker.RetrieveOne)
	}

	company := NewCompanyController()
	companyRoutes := router.Group("/company")
	{
		companyRoutes.GET("/", company.RetrieveAll)
		companyRoutes.POST("/", company.Create)
		companyRoutes.PUT("/:id", company.Update)
		companyRoutes.DELETE("/:id", company.Delete)
		companyRoutes.GET("/:id", company.RetrieveOne)
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
