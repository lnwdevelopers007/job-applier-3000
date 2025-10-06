package controller

import (
	"io"
	"net/http"
	"time"

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

// Upload handles file upload
func (fc FileController) Upload(c *gin.Context) {
	// TODO: Get authenticated user from context (after auth middleware is implemented)
	// For now, accept userID from form
	userIDStr := c.PostForm("userID")
	userRole := c.PostForm("userRole") // "jobSeeker" or "company"

	// Validate userID
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
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

	// TODO: Get authenticated user from context (after auth middleware)
	// For now, accept userID from query param for testing
	requestingUserID := c.Query("requestingUserID")
	if requestingUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	requestingUserObjectID, err := primitive.ObjectIDFromHex(requestingUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requesting user ID"})
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
	if fileDoc.UserID != requestingUserObjectID {
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

	// TODO: Authorization check - user can only list their own files
	// For now, accept requestingUserID from query param
	requestingUserID := c.Query("requestingUserID")
	if requestingUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	requestingUserObjectID, err := primitive.ObjectIDFromHex(requestingUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requesting user ID"})
		return
	}

	// Check if requesting user is trying to access their own files
	if objectID != requestingUserObjectID {
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

	// TODO: Get authenticated user from context (after auth middleware)
	// For now, accept userID from query param
	requestingUserID := c.Query("requestingUserID")
	if requestingUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	requestingUserObjectID, err := primitive.ObjectIDFromHex(requestingUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requesting user ID"})
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
	if fileDoc.UserID != requestingUserObjectID {
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