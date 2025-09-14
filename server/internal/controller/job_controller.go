package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JobController is a custom controller for JobSchema
type JobController struct {
	collectionName string
}

// NewJobController initializes a JobController
func NewJobController() JobController {
	return JobController{collectionName: "jobs"}
}
func (jc JobController) Create(c *gin.Context) {
	db := database.GetDatabase()
	collection := db.Collection(jc.collectionName)

	// Step 1: Bind all fields except companyID
	var job schema.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 2: Extract companyID separately
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyIDStr, ok := raw["companyID"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companyID is required"})
		return
	}

	companyID, err := primitive.ObjectIDFromHex(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid companyID"})
		return
	}

	// Step 3: Overwrite CompanyID in job struct
	job.CompanyID = companyID

	// Step 4: Insert into MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert job"})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// QueryJobs adds advanced filtering (title, company, location, salary range, etc.)
func (jc JobController) QueryJobs() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDatabase()
		collection := db.Collection(jc.collectionName)

		// collect filters from query params
		title := c.Query("title")
		company := c.Query("company")
		location := c.Query("location")
		minSalary := c.Query("minSalary")
		maxSalary := c.Query("maxSalary")

		filter := bson.M{}

		if title != "" {
			filter["title"] = bson.M{"$regex": title, "$options": "i"}
		}
		if company != "" {
			filter["company"] = bson.M{"$regex": company, "$options": "i"}
		}
		if location != "" {
			filter["location"] = bson.M{"$regex": location, "$options": "i"}
		}
		if minSalary != "" || maxSalary != "" {
			salaryFilter := bson.M{}
			if minSalary != "" {
				// convert string to int
				// ignore error handling for now
				// TODO: handle parsing errors safely
				salaryFilter["$gte"] = toInt(minSalary)
			}
			if maxSalary != "" {
				salaryFilter["$lte"] = toInt(maxSalary)
			}
			filter["salary"] = salaryFilter
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var jobs []schema.Job
		if err := cursor.All(ctx, &jobs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, jobs)
	}
}

// small helper
func toInt(s string) int {
	var n int
	_, _ = fmt.Sscan(s, &n)
	return n
}
