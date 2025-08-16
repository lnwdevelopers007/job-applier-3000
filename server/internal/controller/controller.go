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

// GetJobs returns all jobs in the database.
func GetJobs(c *gin.Context) {
	db := database.GetInstance().Database("ja-3000")
	cursor, err := db.Collection("jobs").Find(context.Background(), bson.D{{}})
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

// PostJobs save jobs to the database.
func PostJobs(c *gin.Context) {
	db := database.GetInstance().Database("ja-3000")
	collection := db.Collection("jobs")

	raw, err := database.ValidateRequestContent[database.JobSchema](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := collection.InsertOne(ctx, raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert job"})
		return
	}

	c.JSON(http.StatusCreated, raw)
}
