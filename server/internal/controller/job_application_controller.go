package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/email"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	jobApplicationFilter, shouldReturn := JobApplicationFilter(c)
	if shouldReturn {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	findOpts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	applications, err := findAll[schema.JobApplication](ctx, jc.baseController.collectionName, jobApplicationFilter, findOpts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(applications) == 0 {
		c.JSON(http.StatusOK, []any{})
		return
	}

	// Extract job seeker IDs from applications
	userIDs := extractUnique(
		applications,
		func(app schema.JobApplication) primitive.ObjectID { return app.ApplicantID },
	)

	// query users
	userMap, err := getUsers(
		ctx,
		userIDs,
		"users",
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// // Combine results
	combinedResults := make([]schema.ApplicationWithApplicant, 0, len(applications))
	for _, app := range applications {
		result := schema.ApplicationWithApplicant{
			JobApplication: app,
		}
		var applicant = schema.User{}

		if user, exists := userMap[app.ApplicantID]; exists {
			applicant = user
		}

		result.Applicant = applicant

		combinedResults = append(combinedResults, result)
	}

	c.JSON(http.StatusOK, combinedResults)
}

func getUsers(
	ctx context.Context,
	userIDs []primitive.ObjectID,
	collectionName string,
) (
	map[primitive.ObjectID]schema.User,
	error,
) {
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	result, err := findAll[schema.User](ctx, collectionName, filter)

	// Create a map of job seekers for easy lookup
	resultMap := make(map[primitive.ObjectID]schema.User)
	for _, js := range result {
		resultMap[js.ID] = js
	}

	return resultMap, err
}

func JobApplicationFilter(c *gin.Context) (bson.M, bool) {
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
				return nil, true
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
			return nil, true
		}
	}
	return filter, false
}

// Create adds a new job application
func (jc JobApplicationController) Create(c *gin.Context) {
	shouldReturn := jc.sendEmail(c)
	if shouldReturn {
		return
	}
	jc.baseController.Create(c)
}

func (jc JobApplicationController) sendEmail(c *gin.Context) bool {
	var raw schema.JobApplication
	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return true
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	job, jobErr := findOne[schema.Job](ctx, "jobs", raw.JobID)
	if jobErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": jobErr.Error()})
		return true
	}
	if !job.EmailNotifications {
		return false
	}
	company, companyErr := findOne[schema.User](ctx, "users", job.CompanyID)
	if companyErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": companyErr.Error()})
		return true
	}
	email.Send(company.Email, "New applicant applied", "yay")
	return false
}

// Update updates a job application by ID
func (jc JobApplicationController) Update(c *gin.Context) {
	jc.baseController.Update(c)
}

// Delete deletes a job application by ID
func (jc JobApplicationController) Delete(c *gin.Context) {
	jc.baseController.Delete(c)
}

// RetrieveAll fetches all job applications from the database.
func (jc JobApplicationController) RetrieveAll(c *gin.Context) {
	jc.baseController.RetrieveAll(c)
}

// RetrieveOne fetches a single job application by ID.
func (jc JobApplicationController) RetrieveOne(c *gin.Context) {
	jc.baseController.RetrieveOne(c)
}
