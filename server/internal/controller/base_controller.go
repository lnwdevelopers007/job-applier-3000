// Package controller is like the "controller" in MVC design pattern.
package controller

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
		msg := "Create " + controller.displayName + " failed: incorrect request body"
		slog.Warn(userInfo + msg)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	// Validate input
	if v, ok := any(&raw).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := repository.InsertOne(ctx, raw)
	if err != nil {
		msg := "Create " + controller.displayName + " failed"
		slog.Error(userInfo + msg + ": " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	slog.Info(userInfo + "Created " + controller.displayName)

	c.JSON(http.StatusCreated, res)
}

// RetrieveAll retrieves all documents (row) and all of its attirbutes
// from collectionName collection
func (controller BaseController[Schema, DTO]) RetrieveAll(c *gin.Context) {
	userInfo := getUserForLogging(c)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := repository.FindAll[Schema](ctx, bson.M{})
	if err != nil {
		msg := "Retrieve All " + controller.displayName + " failed"
		slog.Error(userInfo + msg + ": " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	slog.Info(userInfo + "Retrieved All " + controller.displayName)
	c.JSON(http.StatusOK, res)
}

// Update() updates a resource by ID.
func (controller BaseController[Schema, DTO]) Update(c *gin.Context) {
	userInfo := getUserForLogging(c)
	id := c.Param("id") // get :id from URL
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		msg := "Update " + controller.displayName + " failed: invalid ID"
		slog.Warn(userInfo + msg + ": " + id)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	var newData DTO
	if err := c.ShouldBindBodyWithJSON(&newData); err != nil {
		msg := "Update " + controller.displayName + " failed: incorrect request body"
		slog.Warn(userInfo + msg)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := repository.Update[Schema](ctx, objID, newData)

	if err != nil {
		msg := "Update " + controller.displayName + " failed"
		slog.Error(userInfo + msg + ": " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	if res.MatchedCount == 0 {
		msg := "Update " + controller.displayName + " failed: resource not found"
		slog.Warn(userInfo + msg)
		c.JSON(http.StatusNotFound, gin.H{"error": msg})
		return
	}

	msg := "Updated " + controller.displayName + ": " + id
	slog.Info(userInfo + msg)
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// Delete() deletes a resource by ID.
func (controller BaseController[Schema, DTO]) Delete(c *gin.Context) {
	userInfo := getUserForLogging(c)
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		msg := "Delete " + controller.displayName + "failed: invalid ID"
		slog.Warn(userInfo + msg + ": " + id)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := repository.DeleteOne[Schema](ctx, objID)

	if err != nil {
		msg := "Delete " + controller.displayName + "failed"
		slog.Error(userInfo + msg + ": " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	if result.DeletedCount == 0 {
		msg := "Delete " + controller.displayName + "failed: resource not found"
		slog.Warn(userInfo + msg)
		c.JSON(http.StatusNotFound, gin.H{"error": msg})
		return
	}
	msg := "Deleted " + controller.displayName + ": " + id
	slog.Info(userInfo + msg)
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// RetrieveOne retrieves a single document by ID from collectionName collection.
func (controller BaseController[Schema, DTO]) RetrieveOne(c *gin.Context) {
	userInfo := getUserForLogging(c)
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		msg := "Retrieve " + controller.displayName + "failed: invalid ID"
		slog.Warn(userInfo + msg + ": " + id)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := repository.FindOne[Schema](ctx, objID)
	if err != nil {
		msg := "Retrieve " + controller.displayName + "failed: resource not found"
		slog.Warn(userInfo + msg)
		c.JSON(http.StatusNotFound, gin.H{"error": msg})
		return
	}
	msg := "Retrieved " + controller.displayName + ": " + id
	slog.Info(msg)
	c.JSON(http.StatusOK, res)
}
