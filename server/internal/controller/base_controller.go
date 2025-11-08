// Package controller is like the "controller" in MVC design pattern.
package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/repository"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseController implements IController interface.
// It is a basic controller for basic CRUD operations
// involving only one specific collection in the database.
type BaseController[Schema schema.CollectionEntity] struct {
	collectionName string
	displayName    string
}

// Create() inserts one document (row) to collectionName collection.
func (controller BaseController[Schema]) Create(c *gin.Context) {

	var raw Schema
	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := repository.InsertOne(ctx, raw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create " + controller.displayName,
		})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// RetrieveAll retrieves all documents (row) and all of its attirbutes
// from collectionName collection
func (controller BaseController[Schema]) RetrieveAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := repository.FindAll[Schema](ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Update() updates a resource by ID.
func (controller BaseController[Schema]) Update(c *gin.Context) {
	id := c.Param("id") // get :id from URL
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var raw Schema
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// use $set to update only provided fields
	res, err := repository.Update(ctx, objID, raw)
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
func (controller BaseController[Schema]) Delete(c *gin.Context) {
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
func (controller BaseController[Schema]) RetrieveOne(c *gin.Context) {
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
