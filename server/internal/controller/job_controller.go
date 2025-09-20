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

// QueryJobs adds advanced filtering (title, company, location, salary range, etc.)
func (jc JobController) QueryJobs() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDatabase()
		collection := db.Collection(jc.collectionName)

		// Only accept specific query params
		allowedParams := map[string]func(string) (interface{}, error){
			"id": func(v string) (interface{}, error) {
				if v == "" {
					return nil, fmt.Errorf("id parameter is empty")
				}
				return primitive.ObjectIDFromHex(v)
			},
			"title": func(v string) (interface{}, error) {
				if v == "" {
					return nil, nil
				}
				return bson.M{"$regex": v, "$options": "i"}, nil
			},
			"companyID": func(v string) (interface{}, error) {
				if v == "" {
					return nil, nil
				}
				return primitive.ObjectIDFromHex(v)
			},
			"location": func(v string) (interface{}, error) {
				if v == "" {
					return nil, nil
				}
				return bson.M{"$regex": v, "$options": "i"}, nil
			},
			"minSalary": func(v string) (interface{}, error) {
				if v == "" {
					return nil, nil
				}
				return toInt(v), nil
			},
			"maxSalary": func(v string) (interface{}, error) {
				if v == "" {
					return nil, nil
				}
				return toInt(v), nil
			},
		}

		filter := bson.M{}
		salaryFilter := bson.M{}

		// Loop through query params
		for key, values := range c.Request.URL.Query() {
			if fn, ok := allowedParams[key]; ok {
				val, err := fn(values[0])
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				if val != nil {
					switch key {
					case "id":
						filter["_id"] = val
					case "minSalary":
						salaryFilter["$gte"] = val
					case "maxSalary":
						salaryFilter["$lte"] = val
					default:
						filter[key] = val
					}
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported query parameter: " + key})
				return
			}
		}

		if len(salaryFilter) > 0 {
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

		if len(jobs) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No jobs found"})
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
