package controller

import (
	"io"
	"net/http"
	"os"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileController struct {
	baseController BaseController[schema.File]
}

func NewFileController() FileController {
	return FileController{
		baseController: BaseController[schema.File]{
			collectionName: "files",
			displayName:    "File",
		},
	}
}
// getUserFromContext extracts user info from context (set by auth middleware)
func getUserFromContext(c *gin.Context) (userID primitive.ObjectID, role string, err error) {
	enableAuth, _ := strconv.ParseBool(os.Getenv("ENABLE_AUTH"))

	if enableAuth {
		// Get from context (set by auth middleware)
		userIDStr, exists := c.Get("userID")
		if !exists {
			err = http.ErrNotSupported
			return
		}

		userID, err = primitive.ObjectIDFromHex(userIDStr.(string))
		if err != nil {
			return
		}

		roleVal, _ := c.Get("role")
		role, _ = roleVal.(string)
	} else {
		// When auth disabled, use context from middleware (simulated headers)
		userIDStr, exists := c.Get("userID")
		if !exists {
			userID = primitive.NewObjectID()
			role = "jobSeeker"
			err = nil
			return
		}

		userID, err = primitive.ObjectIDFromHex(userIDStr.(string))
		if err != nil {
			userID = primitive.NewObjectID()
		}

		roleVal, _ := c.Get("role")
		role, _ = roleVal.(string)
		if role == "" {
			role = "jobSeeker"
		}
		err = nil
	}

	return
}

// Upload handles file upload
func (fc FileController) Upload(c *gin.Context) {
	// Get authenticated user
	userID, userRole, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	// Validate uploaded file (size, type)
	if err := schema.ValidateUploadedFile(header); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get metadata from form
	category := schema.FileCategory(c.PostForm("category"))

	// Validate category for user role
	if err := schema.ValidateFileCategory(category, userRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	// Create file document
	fileDoc := schema.File{
		UserID:        userID,
		Content:       fileBytes,
		FileExtension: "pdf",
		Filename:      header.Filename,
		ContentType:   header.Header.Get("Content-Type"),
		Size:          header.Size,
		Category:      category,
		UploadDate:    time.Now(),
	}

	// Save to database
	db := database.GetDatabase()
	collection := db.Collection(fc.baseController.collectionName)
	result, err := collection.InsertOne(c.Request.Context(), fileDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	fileDoc.ID = result.InsertedID.(primitive.ObjectID)

	// Return metadata only
	c.JSON(http.StatusCreated, gin.H{
		"id":            fileDoc.ID,
		"userID":        fileDoc.UserID,
		"filename":      fileDoc.Filename,
		"fileExtension": fileDoc.FileExtension,
		"contentType":   fileDoc.ContentType,
		"size":          fileDoc.Size,
		"category":      fileDoc.Category,
		"uploadDate":    fileDoc.UploadDate,
	})
}

// Download retrieves a single file (download)
func (fc FileController) Download(c *gin.Context) {
	fileID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	// Get authenticated user
	requestingUserID, _, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	db := database.GetDatabase()
	collection := db.Collection(fc.baseController.collectionName)

	var fileDoc schema.File
	err = collection.FindOne(c.Request.Context(), bson.M{"_id": objectID}).Decode(&fileDoc)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// Authorization check: User can only download their own files
	if fileDoc.UserID != requestingUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you do not have permission to access this file",
		})
		return
	}

	// Set headers for file download
	c.Header("Content-Type", fileDoc.ContentType)
	c.Header("Content-Disposition", "attachment; filename="+fileDoc.Filename)
	c.Data(http.StatusOK, fileDoc.ContentType, fileDoc.Content)
}


// ListByUser lists all files for a specific user
func (fc FileController) ListByUser(c *gin.Context) {
	userID := c.Param("userId")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get authenticated user
	requestingUserID, _, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Check if requesting user is trying to access their own files
	enableAuth, _ := strconv.ParseBool(os.Getenv("ENABLE_AUTH"))
	if enableAuth && objectID != requestingUserID {
    c.JSON(http.StatusForbidden, gin.H{
        "error": "you can only access your own files",
    })
    return
}

	db := database.GetDatabase()
	collection := db.Collection(fc.baseController.collectionName)

	cursor, err := collection.Find(
		c.Request.Context(),
		bson.M{"userID": objectID},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve files"})
		return
	}
	defer cursor.Close(c.Request.Context())

	var files []gin.H
	for cursor.Next(c.Request.Context()) {
		var file schema.File
		if err := cursor.Decode(&file); err != nil {
			continue
		}

		// Return metadata only (no binary content)
		files = append(files, gin.H{
			"id":            file.ID,
			"userID":        file.UserID,
			"filename":      file.Filename,
			"fileExtension": file.FileExtension,
			"contentType":   file.ContentType,
			"size":          file.Size,
			"category":      file.Category,
			"uploadDate":    file.UploadDate,
		})
	}

	if files == nil {
		files = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

// Delete removes a file by ID
func (fc FileController) Delete(c *gin.Context) {
	fileID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	// Get authenticated user
	requestingUserID, _, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	db := database.GetDatabase()
	collection := db.Collection(fc.baseController.collectionName)

	// First, check if file exists and user owns it
	var fileDoc schema.File
	err = collection.FindOne(c.Request.Context(), bson.M{"_id": objectID}).Decode(&fileDoc)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// Authorization check
	if fileDoc.UserID != requestingUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you do not have permission to delete this file",
		})
		return
	}

	// Delete the file
	result, err := collection.DeleteOne(c.Request.Context(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted successfully"})
}

// GetApplicantFiles allows companies to view files of applicants who applied to their jobs
func (fc FileController) GetApplicantFiles(c *gin.Context) {
	applicationID := c.Param("applicationId")
	appObjectID, err := primitive.ObjectIDFromHex(applicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application ID"})
		return
	}

	// Get authenticated user
	requestingUserID, requestingUserRole, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Only companies can access this endpoint
	if requestingUserRole != "company" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "only companies can access applicant files",
		})
		return
	}

	db := database.GetDatabase()

	// 1. Find the job application
	var application schema.JobApplication
	err = db.Collection("job_applications").FindOne(
		c.Request.Context(),
		bson.M{"_id": appObjectID},
	).Decode(&application)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	// 2. Find the job to verify company ownership
	var job schema.Job
	err = db.Collection("jobs").FindOne(
		c.Request.Context(),
		bson.M{"_id": application.JobID},
	).Decode(&job)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	// 3. Verify the requesting user (company) owns the job
	if job.CompanyID != requestingUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you can only access files for applications to your own jobs",
		})
		return
	}

	// 4. Get applicant's files (only relevant categories: resume, cover_letter, certification)
	cursor, err := db.Collection("files").Find(
		c.Request.Context(),
		bson.M{
			"userID": application.ApplicantID,
			"category": bson.M{
				"$in": []string{"resume", "cover_letter", "certification"},
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve files"})
		return
	}
	defer cursor.Close(c.Request.Context())

	var files []gin.H
	for cursor.Next(c.Request.Context()) {
		var file schema.File
		if err := cursor.Decode(&file); err != nil {
			continue
		}

		// Return metadata only
		files = append(files, gin.H{
			"id":            file.ID,
			"userID":        file.UserID,
			"filename":      file.Filename,
			"fileExtension": file.FileExtension,
			"contentType":   file.ContentType,
			"size":          file.Size,
			"category":      file.Category,
			"uploadDate":    file.UploadDate,
		})
	}

	if files == nil {
		files = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"applicationID": application.ID,
		"applicantID":   application.ApplicantID,
		"jobID":         application.JobID,
		"files":         files,
	})
}