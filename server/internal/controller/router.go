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
		authGroup.POST("/:provider/logout", auth.Logout)
		authGroup.POST("/refresh", auth.RefreshRefreshToken)
		authGroup.GET("/me", middleware.AuthMiddleware(), auth.Me)
	}

	// Job controller
	jobCtrl := NewJobController()

	// Public job routes (no auth required)
	publicJobs := router.Group("/jobs/public")
	{
		publicJobs.GET("/latest", jobCtrl.GetLatestPublic)
	}

	publicuserController := NewUserController()
	publicUser := router.Group("/users/public")
	{
		publicUser.GET("/:id", publicuserController.GetPublicInfo)
	}

	// Apply AuthMiddleware and AccessControlMiddleware to all protected routes
	// Order matters: AuthMiddleware first, then AccessControlMiddleware
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.AccessControlMiddleware())

	// Protected job routes
	jobs := protected.Group("/jobs")
	{
		jobs.GET("/query", jobCtrl.Query)
		jobs.GET("/", jobCtrl.RetrieveAll)
		jobs.POST("/", jobCtrl.Create)
		jobs.PUT("/:id", jobCtrl.Update)
		jobs.DELETE("/:id", jobCtrl.Delete)
		jobs.GET("/:id", jobCtrl.RetrieveOne)
	}

	// Job application routes
	applicationController := NewJobApplicationController()
	applyRoutes := protected.Group("/apply")
	{
		applyRoutes.GET("/query", applicationController.Query)
		applyRoutes.POST("/", applicationController.Create)
		applyRoutes.PUT("/:id", applicationController.Update)
		applyRoutes.DELETE("/:id", applicationController.Delete)
		applyRoutes.GET("/:id", applicationController.RetrieveOne)
	}

	// User routes (admin only)
	userController := NewUserController()
	userRoutes := protected.Group("/users")
	{
		userRoutes.GET("/query", userController.Query)
		userRoutes.GET("/", userController.RetrieveAll)
		userRoutes.POST("/", userController.Create)
		userRoutes.PUT("/:id", userController.Update)
		userRoutes.DELETE("/:id", userController.Delete)
		userRoutes.GET("/:id", userController.RetrieveOne)
		userRoutes.PATCH("/:id/verify", userController.VerifyUser)
		userRoutes.PATCH("/:id/role", userController.EditPermission)
	}

	// File routes
	file := NewFileController()
	fileRoutes := protected.Group("/files")
	{
		fileRoutes.POST("/upload", file.Upload)
		fileRoutes.GET("/download/:id", file.Download)
		fileRoutes.GET("/user/:userId", file.ListByUser)
		fileRoutes.DELETE("/:id", file.Delete)
		fileRoutes.GET("/application/:applicationId", file.GetApplicantFiles)
		fileRoutes.GET("/application/:applicationId/download/:fileId", file.DownloadApplicantFile)
	}

	// Public routes (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	note := NewNoteController()
	noteRoutes := router.Group("/notes")
	noteRoutes.Use(middleware.AuthMiddleware())
	{
		noteRoutes.GET("/", note.Query)
		noteRoutes.GET("/query", note.Query)
		noteRoutes.GET("/:id", note.RetrieveOne)
		noteRoutes.POST("/", note.Create)
		noteRoutes.PUT("/:id", note.Update)
		noteRoutes.DELETE("/:id", note.Delete)
	}

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
