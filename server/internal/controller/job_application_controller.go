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

// JobApplicationController handles JobApplication CRUD operations
type JobApplicationController struct {
	baseController BaseController[schema.JobApplication]
}

// NewJobApplicationController initializes a JobApplicationController
func NewJobApplicationController() JobApplicationController {
	return JobApplicationController{
		baseController: BaseController[schema.JobApplication]{
			collectionName: "job_applications",
			displayName:    "Application",
		},
	}
}

func (jc JobApplicationController) Query(c *gin.Context) {
	db := database.GetDatabase()
	collection := db.Collection(jc.baseController.collectionName)

	allowedParams := map[string]func(string) (any, error){
		"id": func(v string) (any, error) {
			if v == "" {
				return nil, fmt.Errorf("id parameter is empty")
			}
			return primitive.ObjectIDFromHex(v)
		},
		"applicantID": func(v string) (any, error) {
			return primitive.ObjectIDFromHex(v)
		},
		"jobID": func(v string) (any, error) {
			return primitive.ObjectIDFromHex(v)
		},
		"companyID": func(v string) (any, error) {
			return primitive.ObjectIDFromHex(v)
		},
		"status": func(v string) (any, error) {
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
		c.JSON(http.StatusOK, []any{})
		return
	}

	// Extract applicant IDs from applications
	applicantIDs := make([]primitive.ObjectID, 0, len(applications))
	applicantIDMap := make(map[primitive.ObjectID]bool)
	for _, app := range applications {
		if _, exists := applicantIDMap[app.ApplicantID]; !exists {
			applicantIDMap[app.ApplicantID] = true
			applicantIDs = append(applicantIDs, app.ApplicantID)
		}
	}

	fmt.Println(applicantIDs)

	// Query job seekers
	jobSeekerCollection := db.Collection("jobSeeker")

	// Create filter for job seekers
	jobSeekerFilter := bson.M{"_id": bson.M{"$in": applicantIDs}}
	jobSeekerCursor, err := jobSeekerCollection.Find(ctx, jobSeekerFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var jobSeekers []schema.JobSeeker
	if err := jobSeekerCursor.All(ctx, &jobSeekers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(jobSeekers)

	// Create a map of job seekers for easy lookup
	jobSeekerMap := make(map[primitive.ObjectID]schema.JobSeeker)
	for _, js := range jobSeekers {
		jobSeekerMap[js.ID] = js
	}

	// Combine applications with job seeker data
	type ApplicationWithSeekerInfo struct {
		ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		ApplicantID   primitive.ObjectID `bson:"applicantID" json:"applicantID" binding:"required"`
		JobID         primitive.ObjectID `bson:"jobID" json:"jobID" binding:"required"`
		CompanyID     primitive.ObjectID `bson:"companyID" json:"companyID" binding:"required"`
		Status        string             `bson:"status" json:"status" binding:"required"`
		CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
		ApplicantInfo schema.JobSeeker   `json:"applicantInfo"`
	}

	combinedResults := make([]ApplicationWithSeekerInfo, 0, len(applications))
	for _, app := range applications {
		result := ApplicationWithSeekerInfo{
			ID:          app.ID,
			ApplicantID: app.ApplicantID,
			JobID:       app.JobID,
			CompanyID:   app.CompanyID,
			Status:      app.Status,
			CreatedAt:   app.CreatedAt,
		}

		if jobSeeker, exists := jobSeekerMap[app.ApplicantID]; exists {
			result.ApplicantInfo = jobSeeker
		}

		combinedResults = append(combinedResults, result)
	}

	c.JSON(http.StatusOK, combinedResults)
}

// Create adds a new job application
func (jc JobApplicationController) Create(c *gin.Context) {
	jc.baseController.Create(c)
}

// Update updates a job application by ID
func (jc JobApplicationController) Update(c *gin.Context) {
	jc.baseController.Update(c)
}

// Delete deletes a job application by ID
func (jc JobApplicationController) Delete(c *gin.Context) {
	jc.baseController.Delete(c)
}
