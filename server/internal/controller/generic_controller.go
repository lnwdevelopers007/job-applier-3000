// Package controller is like the "controller" in MVC design pattern.
package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

// GenericController implements IController interface.
// It is a basic controller for basic CRUD operations
// involving only one specific collection in the database.
type GenericController[Schema any] struct {
	collectionName string
}

// Create() inserts one document (row) to collectionName collection.
func (controller GenericController[Schema]) Create(c *gin.Context) {
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

// RetrieveAll retrieves all documents (row) and all of its attirbutes
// from collectionName collection
func (controller GenericController[Schema]) RetrieveAll(c *gin.Context) {
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

// Update() to be implemented.
func (controller GenericController[Schema]) Update(c *gin.Context) {
	panic("Not Implemented")
}

// Delete() to be implemented.
func (controller GenericController[Schema]) Delete(c *gin.Context) {
	panic("Not Implemented")
}
