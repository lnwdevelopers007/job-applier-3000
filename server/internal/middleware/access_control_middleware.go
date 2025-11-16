package middleware

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AccessControlMiddleware checks if user is banned and has permission to access the route
func AccessControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		enableAuth, _ := strconv.ParseBool(os.Getenv("ENABLE_AUTH"))
		
		// Skip access control if auth is disabled (for development)
		if !enableAuth {
			c.Next()
			return
		}

		// Skip for public routes (auth routes)
		if strings.HasPrefix(c.Request.URL.Path, "/auth/") || 
		   strings.HasPrefix(c.Request.URL.Path, "/health") ||
		   strings.HasPrefix(c.Request.URL.Path, "/swagger/") ||
		   strings.HasPrefix(c.Request.URL.Path, "/jobs/public/") ||
		   strings.HasPrefix(c.Request.URL.Path, "/users/public/") {
			c.Next()
			return
		}

		// 1. Check if user is banned
		if err := checkBanStatus(c); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "account_banned",
				"message": "Your account has been banned. Please contact support.",
			})
			c.Abort()
			return
		}

		// 2. Check RBAC permissions
		if err := checkRoutePermission(c); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "access_denied",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkBanStatus verifies if the user is banned
func checkBanStatus(c *gin.Context) error {
	userID, exists := c.Get("userID")
	if !exists {
		// No user ID means not authenticated, let AuthMiddleware handle it
		return nil
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return nil
	}

	// Query database for fresh ban status
	db := database.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil
	}

	var user schema.User
	err = db.Collection("users").FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil // User not found, let other middleware handle
	}

	if user.Banned {
		return ErrUserBanned
	}

	return nil
}

// checkRoutePermission checks if the user has permission to access the route
func checkRoutePermission(c *gin.Context) error {
	role, exists := c.Get("role")
	if !exists {
		return ErrNoRole
	}

	roleStr, ok := role.(string)
	if !ok {
		return ErrInvalidRole
	}

	// Build route key: "METHOD:PATH"
	method := c.Request.Method
	path := c.Request.URL.Path
	routeKey := method + ":" + path

	// Try to match exact route first
	permission, found := RoutePermissions[routeKey]
	
	// If not found, try to match with pattern (e.g., /jobs/:id)
	if !found {
		matchedKey, perm, ok := matchRoutePattern(method, path)
		if ok {
			routeKey, permission, found = matchedKey, perm, ok
			log.Printf("Matched route pattern: %s", routeKey)
		}
	}

	// If route not in permission map, allow (routes without restrictions)
	if !found {
		return nil
	}

	// Check if user's role is allowed
	if !permission.IsRoleAllowed(roleStr) {
		return ErrInsufficientPermissions
	}

	// Check ownership if required
	if permission.RequireOwnership {
		if err := checkOwnership(c, roleStr, path); err != nil {
			return err
		}
	}

	return nil
}

// matchRoutePattern matches dynamic routes like /jobs/:id
func matchRoutePattern(method, path string) (string, Permission, bool) {
	parts := strings.Split(path, "/")
	
	// Try different patterns based on path structure
	patterns := []string{
		method + ":" + path,
	}

	// For paths with IDs (e.g., /jobs/123abc), try pattern with :id
	if len(parts) >= 3 {
		// /resource/:id pattern (e.g., /jobs/:id)
		pattern := method + ":" + strings.Join(parts[:len(parts)-1], "/") + "/:id"
		patterns = append(patterns, pattern)
	}

	if len(parts) >= 4 {
		// /resource/:id/action pattern (e.g., /files/download/:id)
		pattern := method + ":" + parts[1] + "/" + parts[2] + "/:id"
		patterns = append(patterns, pattern)
	}

	if len(parts) >= 5 {
		// /files/application/:applicationId/download/:fileId
		if parts[1] == "files" && parts[2] == "application" {
			pattern := method + ":/files/application/:applicationId/download/:fileId"
			patterns = append(patterns, pattern)
		}
	}

	// Try each pattern
	for _, pattern := range patterns {
		if perm, found := RoutePermissions[pattern]; found {
			return pattern, perm, true
		}
	}

	return "", Permission{}, false
}

// checkOwnership verifies if the user owns the resource they're trying to access
func checkOwnership(c *gin.Context, role, path string) error {
	userID, _ := c.Get("userID")
	userIDStr, _ := userID.(string)
	userObjID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return ErrInvalidUserID
	}

	// Admin has access to everything
	if role == "admin" {
		return nil
	}

	db := database.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// ===== USER ROUTES - SPECIAL HANDLING =====
	if strings.HasPrefix(path, "/users/") && !strings.Contains(path, "verify") && !strings.Contains(path, "role") {
		var targetUserID string
		
		// Handle /users/query?id=xxx
		if strings.Contains(path, "/users/query") {
			targetUserID = c.Query("id")
			if targetUserID == "" {
				// If no ID in query, this might be listing users - let it through for now
				// Admin role check already passed if we got here
				return nil
			}
		} else {
			// Handle /users/:id
			targetUserID = extractIDFromPath(path, "/users/")
		}
		
		// Check if viewing/editing own profile
		if targetUserID == userIDStr {
			return nil // Always allow accessing own profile
		}

		// For GET requests (viewing other users), check role-based permissions
		if c.Request.Method == "GET" {
			targetObjID, err := primitive.ObjectIDFromHex(targetUserID)
			if err != nil {
				return ErrInvalidResourceID
			}

			// Get target user's role from database
			var targetUser schema.User
			err = db.Collection("users").FindOne(ctx, bson.M{"_id": targetObjID}).Decode(&targetUser)
			if err != nil {
				return ErrResourceNotFound
			}

			// Check if viewer can see target user based on roles
			if !canViewUserRole(role, targetUser.Role) {
				return ErrNotResourceOwner
			}
			return nil // Allowed to view
		}

		// For PUT/PATCH/DELETE, must be own profile (non-admin)
		return ErrNotResourceOwner
	}

	// ===== JOB ROUTES =====
	if strings.HasPrefix(path, "/jobs/") && !strings.Contains(path, "query") {
		// Job ownership check
		jobID := extractIDFromPath(path, "/jobs/")
		jobObjID, err := primitive.ObjectIDFromHex(jobID)
		if err != nil {
			return ErrInvalidResourceID
		}

		var job schema.Job
		err = db.Collection("jobs").FindOne(ctx, bson.M{"_id": jobObjID}).Decode(&job)
		if err != nil {
			return ErrResourceNotFound
		}

		if job.CompanyID != userObjID {
			return ErrNotResourceOwner
		}
	} else if strings.HasPrefix(path, "/apply/") && !strings.Contains(path, "query") {
		// ===== APPLICATION ROUTES =====
		appID := extractIDFromPath(path, "/apply/")
		appObjID, err := primitive.ObjectIDFromHex(appID)
		if err != nil {
			return ErrInvalidResourceID
		}

		var app schema.JobApplication
		err = db.Collection("job_applications").FindOne(ctx, bson.M{"_id": appObjID}).Decode(&app)
		if err != nil {
			return ErrResourceNotFound
		}

		// Job seeker can access their own applications
		// Company can access applications for their jobs
		if role == "jobSeeker" && app.ApplicantID != userObjID {
			return ErrNotResourceOwner
		} else if role == "company" {
			// Check if job belongs to company
			var job schema.Job
			err = db.Collection("jobs").FindOne(ctx, bson.M{"_id": app.JobID}).Decode(&job)
			if err != nil {
				return ErrResourceNotFound
			}
			if job.CompanyID != userObjID {
				return ErrNotResourceOwner
			}
		}
	} else if strings.HasPrefix(path, "/files/") {
		// ===== FILE ROUTES =====
		if strings.Contains(path, "/files/user/") {
			// /files/user/:userId - check if accessing own files
			targetUserID := extractIDFromPath(path, "/files/user/")
			if targetUserID != userIDStr {
				return ErrNotResourceOwner
			}
		} else if strings.Contains(path, "/files/download/") || strings.Contains(path, "/files/") && c.Request.Method == "DELETE" {
			// /files/download/:id or DELETE /files/:id
			fileID := extractIDFromPath(path, "/files/download/")
			if fileID == "" {
				fileID = extractIDFromPath(path, "/files/")
			}
			fileObjID, err := primitive.ObjectIDFromHex(fileID)
			if err != nil {
				return ErrInvalidResourceID
			}

			var file schema.File
			err = db.Collection("files").FindOne(ctx, bson.M{"_id": fileObjID}).Decode(&file)
			if err != nil {
				return ErrResourceNotFound
			}

			if file.UserID != userObjID {
				return ErrNotResourceOwner
			}
		}
	}

	return nil
}

// canViewUserRole checks if the viewer can see a target user's profile
// This implements the role-based viewing matrix for user profiles
func canViewUserRole(viewerRole string, targetRole string) bool {
	// Admin can view everyone (already handled before this function)
	if viewerRole == "admin" {
		return true
	}

	// Job Seeker can view: company profiles only (NOT other job seekers, faculty, or admin)
	if viewerRole == "jobSeeker" {
		return targetRole == "company"
	}

	// Company can view: job seeker profiles AND other company profiles (for job postings)
	// BUT NOT faculty or admin
	if viewerRole == "company" {
		return targetRole == "jobSeeker" || targetRole == "company"
	}

	// Faculty can view: company profiles (for job browsing)
	if viewerRole == "faculty" {
		return targetRole == "company"
	}

	// Default: deny
	return false
}

// extractIDFromPath extracts the ID from a path like /jobs/123abc
func extractIDFromPath(path, prefix string) string {
	path = strings.TrimPrefix(path, prefix)
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// Custom errors
var (
	ErrUserBanned              = &AccessError{Code: "user_banned", Message: "User account is banned"}
	ErrNoRole                  = &AccessError{Code: "no_role", Message: "User role not found"}
	ErrInvalidRole             = &AccessError{Code: "invalid_role", Message: "Invalid user role"}
	ErrInsufficientPermissions = &AccessError{Code: "insufficient_permissions", Message: "You do not have permission to access this resource"}
	ErrInvalidUserID           = &AccessError{Code: "invalid_user_id", Message: "Invalid user ID"}
	ErrInvalidResourceID       = &AccessError{Code: "invalid_resource_id", Message: "Invalid resource ID"}
	ErrResourceNotFound        = &AccessError{Code: "resource_not_found", Message: "Resource not found"}
	ErrNotResourceOwner        = &AccessError{Code: "not_owner", Message: "You do not own this resource"}
)

// AccessError represents an access control error
type AccessError struct {
	Code    string
	Message string
}

func (e *AccessError) Error() string {
	return e.Message
}
