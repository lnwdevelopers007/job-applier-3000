// Package controller is like the "controller" in MVC design pattern.
package controller

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/repository"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseController implements IController interface.
// It is a basic controller for basic CRUD operations
// involving only one specific collection in the database.
type BaseController[Schema schema.CollectionEntity, DTO any] struct {
	collectionName string
	displayName    string
}

// Create() inserts one document (row) to collectionName collection.
func (controller BaseController[Schema, DTO]) Create(c *gin.Context) {
	userInfo := getUserForLogging(c)
	var raw Schema
	if err := c.ShouldBindBodyWithJSON(&raw); err != nil {
		msg := "Cannot create " + controller.displayName + ", incorrect schema"
		slog.Warn(userInfo + msg)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := repository.InsertOne(ctx, raw)
	if err != nil {
		msg := "Failed to create " + controller.displayName
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		slog.Error(userInfo + msg)
		return
	}

	slog.Info(userInfo + "Created " + controller.displayName + " Successfully")

	c.JSON(http.StatusCreated, res)
}

// RetrieveAll retrieves all documents (row) and all of its attirbutes
// from collectionName collection
func (controller BaseController[Schema, DTO]) RetrieveAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := repository.FindAll[Schema](ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cannot retrieve " + controller.displayName})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Update() updates a resource by ID.
func (controller BaseController[Schema, DTO]) Update(c *gin.Context) {
	id := c.Param("id") // get :id from URL
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var newData DTO
	if err := c.ShouldBindBodyWithJSON(&newData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := repository.Update[Schema](ctx, objID, newData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update " + controller.displayName,
		})
		return
	}

	if res.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": controller.displayName + " not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": controller.displayName + " updated successfully",
	})
}

// Delete() deletes a resource by ID.
func (controller BaseController[Schema, DTO]) Delete(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := repository.DeleteOne[Schema](ctx, objID)

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No " + controller.displayName + " found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Delete Request",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": controller.displayName + " Deleted successfully",
	})
}

// RetrieveOne retrieves a single document by ID from collectionName collection.
func (controller BaseController[Schema, DTO]) RetrieveOne(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := repository.FindOne[Schema](ctx, objID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": controller.displayName + " not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (controller BaseController[Schema, DTO]) PatchOne(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty update payload"})
		return
	}

	update := bson.M{"$set": body}

	db := database.GetDatabase()
	collection := db.Collection(controller.collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update " + controller.displayName,
		})
		return
	}

	if res.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": controller.displayName + " not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": controller.displayName + " patched successfully",
	})
}
