package controller

import (
	"context"
	"net/http"
	"time"
    "fmt"

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

func (jc JobApplicationController) QueryApplications() gin.HandlerFunc {
    return func(c *gin.Context) {
        db := database.GetDatabase()
        collection := db.Collection(jc.collectionName)

        allowedParams := map[string]func(string) (interface{}, error){
            "id": func(v string) (interface{}, error) {
                if v == "" {
                    return nil, fmt.Errorf("id parameter is empty")
                }
                return primitive.ObjectIDFromHex(v)
            },
            "applicantID": func(v string) (interface{}, error) {
                return primitive.ObjectIDFromHex(v)
            },
            "jobID": func(v string) (interface{}, error) {
                return primitive.ObjectIDFromHex(v)
            },
            "companyID": func(v string) (interface{}, error) {
                return primitive.ObjectIDFromHex(v)
            },
            "status": func(v string) (interface{}, error) {
                if v == "" {
                    return nil, nil
                }
                return v, nil
            },
        }

        filter := bson.M{}

        // Loop through query params
        for key, value := range c.Request.URL.Query() {
            if fn, ok := allowedParams[key]; ok {
                val, err := fn(value[0])
                if err != nil {
                    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                    return
                }
                if val != nil {
                    if key == "id" {
                        filter["_id"] = val
                    } else {
                        filter[key] = val
                    }
                }
            } else {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported query parameter: " + key})
                return
            }
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

        if len(applications) == 0 {
            c.JSON(http.StatusNotFound, gin.H{"message": "No applications found"})
            return
        }

        c.JSON(http.StatusOK, applications)
    }
}




// CreateApplication adds a new job application
func (jc JobApplicationController) CreateApplication() gin.HandlerFunc {
    return func(c *gin.Context) {
        db := database.GetDatabase()
        collection := db.Collection(jc.collectionName)

        var application schema.JobApplication
        if err := c.ShouldBindJSON(&application); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        application.ID = primitive.NewObjectID()
        application.CreatedAt = time.Now()

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        res, err := collection.InsertOne(ctx, application)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
            return
        }

        c.JSON(http.StatusCreated, res)
    }
}

// UpdateApplication updates a job application by ID
func (jc JobApplicationController) UpdateApplication() gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        objID, err := primitive.ObjectIDFromHex(id)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
            return
        }

        var application schema.JobApplication
        if err := c.ShouldBindJSON(&application); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        db := database.GetDatabase()
        collection := db.Collection(jc.collectionName)

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        update := bson.M{"$set": application}
        res, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application"})
            return
        }

        if res.MatchedCount == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Application updated successfully"})
    }
}

// DeleteApplication deletes a job application by ID
func (jc JobApplicationController) DeleteApplication() gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        objID, err := primitive.ObjectIDFromHex(id)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
            return
        }

        db := database.GetDatabase()
        collection := db.Collection(jc.collectionName)

        res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        if res.DeletedCount == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Application deleted successfully"})
    }
}
