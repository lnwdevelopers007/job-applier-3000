package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestDB(t *testing.T) {
	// Setup test database connection
	// You'll need to implement this based on your test DB setup
}

func teardownTestDB(t *testing.T) {
	// Cleanup test database
}

func createTestUser(t *testing.T, role string, banned bool) schema.User {
	user := schema.User{
		ID:       primitive.NewObjectID(),
		UserID:   primitive.NewObjectID().Hex(),
		Email:    "test@example.com",
		Name:     "Test User",
		Role:     role,
		Banned:      banned,
		Verified: true,
	}

	// Insert into test database
	db := database.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.Collection("users").InsertOne(ctx, user)
	assert.NoError(t, err)

	return user
}

func TestAccessControlMiddleware_BannedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("ENABLE_AUTH", "true")
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create banned user
	bannedUser := createTestUser(t, "jobSeeker", true)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Simulate AuthMiddleware setting user context
		c.Set("userID", bannedUser.ID.Hex())
		c.Set("role", bannedUser.Role)
		c.Next()
	})
	router.Use(AccessControlMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "account_banned")
}

func TestAccessControlMiddleware_NotBannedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("ENABLE_AUTH", "true")
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create non-banned user
	user := createTestUser(t, "jobSeeker", false)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", user.ID.Hex())
		c.Set("role", user.Role)
		c.Next()
	})
	router.Use(AccessControlMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAccessControlMiddleware_RolePermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("ENABLE_AUTH", "true")
	setupTestDB(t)
	defer teardownTestDB(t)

	tests := []struct {
		name           string
		userRole       string
		method         string
		path           string
		expectedStatus int
	}{
		// Job routes
		{
			name:           "Company can create job",
			userRole:       "company",
			method:         "POST",
			path:           "/jobs/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "JobSeeker cannot create job",
			userRole:       "jobSeeker",
			method:         "POST",
			path:           "/jobs/",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Admin can create job",
			userRole:       "admin",
			method:         "POST",
			path:           "/jobs/",
			expectedStatus: http.StatusOK,
		},
		// Application routes
		{
			name:           "JobSeeker can apply",
			userRole:       "jobSeeker",
			method:         "POST",
			path:           "/apply/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Company cannot apply",
			userRole:       "company",
			method:         "POST",
			path:           "/apply/",
			expectedStatus: http.StatusForbidden,
		},
		// User routes (admin only)
		{
			name:           "Admin can access users",
			userRole:       "admin",
			method:         "GET",
			path:           "/users/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "JobSeeker cannot access users",
			userRole:       "jobSeeker",
			method:         "GET",
			path:           "/users/",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Company cannot access users",
			userRole:       "company",
			method:         "GET",
			path:           "/users/",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := createTestUser(t, tt.userRole, false)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("userID", user.ID.Hex())
				c.Set("role", user.Role)
				c.Next()
			})
			router.Use(AccessControlMiddleware())
			router.Handle(tt.method, tt.path, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAccessControlMiddleware_OwnershipCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("ENABLE_AUTH", "true")
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create company user and their job
	company := createTestUser(t, "company", false)
	otherCompany := createTestUser(t, "company", false)

	job := schema.Job{
		ID:        primitive.NewObjectID(),
		Title:     "Test Job",
		CompanyID: company.ID,
	}

	db := database.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.Collection("jobs").InsertOne(ctx, job)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		user           schema.User
		jobID          primitive.ObjectID
		expectedStatus int
	}{
		{
			name:           "Owner can update their job",
			user:           company,
			jobID:          job.ID,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-owner company cannot update job",
			user:           otherCompany,
			jobID:          job.ID,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("userID", tt.user.ID.Hex())
				c.Set("role", tt.user.Role)
				c.Next()
			})
			router.Use(AccessControlMiddleware())
			router.PUT("/jobs/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest("PUT", "/jobs/"+tt.jobID.Hex(), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAccessControlMiddleware_DisabledAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("ENABLE_AUTH", "false")

	router := gin.New()
	router.Use(AccessControlMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAccessControlMiddleware_PublicRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("ENABLE_AUTH", "true")

	publicPaths := []string{
		"/auth/google",
		"/auth/google/callback",
		"/health",
		"/swagger/index.html",
	}

	for _, path := range publicPaths {
		t.Run("Public route: "+path, func(t *testing.T) {
			router := gin.New()
			router.Use(AccessControlMiddleware())
			router.GET(path, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest("GET", path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestMatchRoutePattern(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedKey    string
		shouldMatch    bool
	}{
		{
			name:        "Match job by ID",
			method:      "GET",
			path:        "/jobs/507f1f77bcf86cd799439011",
			expectedKey: "GET:/jobs/:id",
			shouldMatch: true,
		},
		{
			name:        "Match file download",
			method:      "GET",
			path:        "/files/download/507f1f77bcf86cd799439011",
			expectedKey: "GET:/files/download/:id",
			shouldMatch: true,
		},
		{
			name:        "Match complex file route",
			method:      "GET",
			path:        "/files/application/507f1f77bcf86cd799439011/download/507f191e810c19729de860ea",
			expectedKey: "GET:/files/application/:applicationId/download/:fileId",
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, _, found := matchRoutePattern(tt.method, tt.path)
			assert.Equal(t, tt.shouldMatch, found)
			if found {
				assert.Equal(t, tt.expectedKey, key)
			}
		})
	}
}

func TestExtractIDFromPath(t *testing.T) {
	tests := []struct {
		path     string
		prefix   string
		expected string
	}{
		{
			path:     "/jobs/507f1f77bcf86cd799439011",
			prefix:   "/jobs/",
			expected: "507f1f77bcf86cd799439011",
		},
		{
			path:     "/files/download/507f1f77bcf86cd799439011",
			prefix:   "/files/download/",
			expected: "507f1f77bcf86cd799439011",
		},
		{
			path:     "/apply/507f1f77bcf86cd799439011",
			prefix:   "/apply/",
			expected: "507f1f77bcf86cd799439011",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := extractIDFromPath(tt.path, tt.prefix)
			assert.Equal(t, tt.expected, result)
		})
	}
}