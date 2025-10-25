package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lnwdevelopers007/job-applier-3000/server/docs"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/auth"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		authGroup.GET("/me", middleware.AuthMiddleware(), auth.Me)
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
		applyRoutes.GET("/", applicationController.Query)
		applyRoutes.POST("/", applicationController.Create)
		applyRoutes.PUT("/:id", applicationController.Update)
		applyRoutes.DELETE("/:id", applicationController.Delete)
		applyRoutes.GET("/:id", applicationController.RetrieveOne)
	}

	userController := NewUserController()
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/query", userController.Query)
		userRoutes.GET("/", userController.RetrieveAll)
		userRoutes.POST("/", userController.Create)
		userRoutes.PUT("/:id", userController.Update)
		userRoutes.DELETE("/:id", userController.Delete)
		userRoutes.GET("/:id", userController.RetrieveOne)

	}

	file := NewFileController()
	fileRoutes := router.Group("/files")
	fileRoutes.Use(middleware.AuthMiddleware())
	{
		fileRoutes.POST("/upload", file.Upload)
		fileRoutes.GET("/download/:id", file.Download)
		fileRoutes.GET("/user/:userId", file.ListByUser)
		fileRoutes.DELETE("/:id", file.Delete)
		fileRoutes.GET("/application/:applicationId", file.GetApplicantFiles)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func setUpCors() cors.Config {
	cfg := cors.Config{
		AllowOrigins: []string{
			os.Getenv("FRONTEND"),
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-Id", "X-User-Role"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cfg
}
