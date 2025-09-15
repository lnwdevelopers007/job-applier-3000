package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JobApplicationController handles JobApplication CRUD operations
type JobApplicationController struct {
	collectionName string
}

// NewJobApplicationController initializes a JobApplicationController
func NewJobApplicationController() JobApplicationController {
	return JobApplicationController{collectionName: "job_applications"}
}

// QueryApplications supports querying applications with optional filters
func (jc JobApplicationController) QueryApplications() gin.HandlerFunc {
    return func(c *gin.Context) {
        db := database.GetDatabase()
        collection := db.Collection(jc.collectionName)

        // Filters from query params
        applicantID := c.Query("applicantID")
        jobID := c.Query("jobID")
        companyID := c.Query("companyID")
        status := c.Query("status")

        filter := bson.M{}

        if applicantID != "" {
            if objID, err := primitive.ObjectIDFromHex(applicantID); err == nil {
                filter["applicantID"] = objID
            }
        }
        if jobID != "" {
            if objID, err := primitive.ObjectIDFromHex(jobID); err == nil {
                filter["jobID"] = objID
            }
        }
        if companyID != "" {
            if objID, err := primitive.ObjectIDFromHex(companyID); err == nil {
                filter["companyID"] = objID
            }
        }
        if status != "" {
            filter["status"] = status
        }

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        cursor, err := collection.Find(ctx, filter)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        var applications []schema.JobApplication
        if err := cursor.All(ctx, &applications); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, applications)
    }
}
