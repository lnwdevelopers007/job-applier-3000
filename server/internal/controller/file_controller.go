package controller

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileController struct {
	baseController BaseController[schema.File, schema.File]
}

func NewFileController() FileController {
	return FileController{
		baseController: BaseController[schema.File, schema.File]{
			collectionName: "files",
			displayName:    "File",
		},
	}
}

// Upload godoc
// @Summary      Upload a file
// @Description  Uploads a PDF or image file for a specific user. Supports categories such as resume, profile_picture, etc.
// @Tags         Files
// @Accept       multipart/form-data
// @Produce      json
// @Param        userID    formData  string  true   "User ID (hex ObjectID)"
// @Param        userRole  formData  string  true   "User role (jobSeeker or company)"
// @Param        category  formData  string  true   "File category (e.g., resume, certificate, logo)"
// @Param        file      formData  file    true   "File to upload"
// @Success      201  {object}  schema.File  "File uploaded successfully"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /files/upload [post]
func (fc FileController) Upload(c *gin.Context) {
	// Get authenticated user
	userID, userRole, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Check user role in database matches the role from context (only when auth is enabled)
	db := database.GetDatabase()
	enableAuth, _ := strconv.ParseBool(os.Getenv("ENABLE_AUTH"))
	if enableAuth {
		userCollection := db.Collection("users")
		var userDoc struct {
			Role string `bson:"role"`
		}
		err = userCollection.FindOne(c.Request.Context(), bson.M{"_id": userID}).Decode(&userDoc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify user role"})
			return
		}
		if !strings.EqualFold(userDoc.Role, userRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "role mismatch: user role in system does not match current context"})
			return
		}
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

// Download godoc
// @Summary      Download a file
// @Description  Downloads a file by its ID. Only the owner can access their files.
// @Tags         Files
// @Accept       json
// @Produce      octet-stream
// @Param        id               path      string  true  "File ID"
// @Param        requestingUserID query     string  true  "ID of the requesting user"
// @Success      200  {file}      binary
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /files/{id}/download [get]
func (fc FileController) Download(c *gin.Context) {
	fileID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	// Get authenticated user
	requestingUserID, userRole, err := getUserFromContext(c)
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

	// Authorization check: User can only download their own files or admin can download any file
	if userRole != "admin" && fileDoc.UserID != requestingUserID {
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

// ListByUser godoc
// @Summary      List all files for a user
// @Description  Retrieves metadata for all files belonging to a specific user. Only the owner can view their files.
// @Tags         Files
// @Accept       json
// @Produce      json
// @Param        userId           path      string  true  "User ID"
// @Param        requestingUserID query     string  true  "ID of the requesting user"
// @Success      200  {object}  map[string][]schema.File  "List of files"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /files/user/{userId} [get]
func (fc FileController) ListByUser(c *gin.Context) {
	userID := c.Param("userId")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get authenticated user
	requestingUserID, userRole, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Check if requesting user is trying to access their own files or is an admin
	enableAuth, _ := strconv.ParseBool(os.Getenv("ENABLE_AUTH"))
	if enableAuth && userRole != "admin" && objectID != requestingUserID {
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

// Delete godoc
// @Summary      Delete a file
// @Description  Deletes a file by its ID. Only the owner of the file can delete it.
// @Tags         Files
// @Accept       json
// @Produce      json
// @Param        id               path      string  true  "File ID"
// @Param        requestingUserID query     string  true  "ID of the requesting user"
// @Success      200  {object}  map[string]string  "File deleted successfully"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /files/{id} [delete]
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

// GetApplicantFiles godoc
// @Summary      Get applicant files for a job application
// @Description  Allows a company to view files (resume, transcript, certification) of an applicant for a specific job application. Only the company who owns the job can access.
// @Tags         Files
// @Accept       json
// @Produce      json
// @Param        applicationId   path      string  true  "Job Application ID"
// @Param        requestingUserID query     string  true  "ID of the requesting user (company)"
// @Success      200 {object} map[string]interface{} "Metadata of applicant's files"
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /files/applicant/{applicationId} [get]
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

	// 4. Get applicant's files (only relevant categories: resume, transcript, certification)
	cursor, err := db.Collection("files").Find(
		c.Request.Context(),
		bson.M{
			"userID": application.ApplicantID,
			"category": bson.M{
				"$in": []string{"resume", "transcript", "certification"},
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

// DownloadApplicantFile godoc
// @Summary      Download an applicant's file
// @Description  Allows a company to download a specific file from an applicant for a job application. Only the company who owns the job can access.
// @Tags         Files
// @Accept       json
// @Produce      octet-stream
// @Param        applicationId   path      string  true  "Job Application ID"
// @Param        fileId          path      string  true  "File ID"
// @Success      200  {file}      binary
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /files/application/{applicationId}/download/{fileId} [get]
func (fc FileController) DownloadApplicantFile(c *gin.Context) {
	applicationID := c.Param("applicationId")
	appObjectID, err := primitive.ObjectIDFromHex(applicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application ID"})
		return
	}

	fileID := c.Param("fileId")
	fileObjectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
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

	// 4. Get the file and verify it belongs to the applicant
	var fileDoc schema.File
	err = db.Collection("files").FindOne(
		c.Request.Context(),
		bson.M{"_id": fileObjectID},
	).Decode(&fileDoc)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// 5. Verify the file belongs to the applicant
	if fileDoc.UserID != application.ApplicantID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "file does not belong to this applicant",
		})
		return
	}

	// 6. Serve the file
	c.Header("Content-Type", fileDoc.ContentType)
	c.Header("Content-Disposition", "attachment; filename="+fileDoc.Filename)
	c.Data(http.StatusOK, fileDoc.ContentType, fileDoc.Content)
}
