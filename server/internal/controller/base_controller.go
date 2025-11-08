// Package controller is like the "controller" in MVC design pattern.
package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseController implements IController interface.
// It is a basic controller for basic CRUD operations
// involving only one specific collection in the database.
type BaseController[Schema any] struct {
	collectionName string
	displayName    string
}

// Create() inserts one document (row) to collectionName collection.
func (controller BaseController[Schema]) Create(c *gin.Context) {
	db := database.GetDatabase()
	collection := db.Collection(controller.collectionName)

	var raw Schema
	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, raw)
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
	res, err := findAll[Schema](ctx, controller.collectionName, bson.M{})
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

	// Use map to only update fields that are actually provided in the request
	var updateFields map[string]any
	if err := c.ShouldBindJSON(&updateFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add updatedAt timestamp
	updateFields["updatedAt"] = time.Now()

	db := database.GetDatabase()
	collection := db.Collection(controller.collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// use $set to update only provided fields
	update := bson.M{"$set": updateFields}

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

	db := database.GetDatabase()
	collection := db.Collection(controller.collectionName)

	// Delete the document from MongoDB
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No " + controller.displayName + " found",
		})
		return
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
	res, err := findOne[Schema](ctx, controller.collectionName, objID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": controller.displayName + " not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (controller BaseController[Schema]) PatchOne(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
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
