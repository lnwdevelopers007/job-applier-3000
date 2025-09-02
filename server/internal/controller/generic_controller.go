// Package controller is like the "controller" in MVC design pattern.
package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenericController implements IController interface.
// It is a basic controller for basic CRUD operations
// involving only one specific collection in the database.
type GenericController[Schema any] struct {
	collectionName string
}

// Create() inserts one document (row) to collectionName collection.
func (controller GenericController[Schema]) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDatabase()
		collection := db.Collection(controller.collectionName)

		var raw Schema
		if err := c.ShouldBindJSON(&raw); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := collection.InsertOne(ctx, raw)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert job"})
			return
		}

		c.JSON(http.StatusCreated, res)
	}
}

// RetrieveAll retrieves all documents (row) and all of its attirbutes
// from collectionName collection
func (controller GenericController[Schema]) RetrieveAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDatabase()
		collection := db.Collection(controller.collectionName)
		cursor, err := collection.Find(context.Background(), bson.D{{}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var jobs []bson.M
		if err = cursor.All(context.Background(), &jobs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, jobs)
	}
}

// Update() updates a resource by ID.
func (controller GenericController[Schema]) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		db := database.GetDatabase()
		collection := db.Collection(controller.collectionName)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// use $set to update only provided fields
		update := bson.M{"$set": raw}

		res, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job"})
			return
		}

		if res.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Job updated successfully"})
	}
}

// Delete() deletes a resource by ID.
func (controller GenericController[Schema]) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
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
			c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
	}
}