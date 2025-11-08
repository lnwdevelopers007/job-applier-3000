package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/email"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/repository"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// JobController is a custom controller for JobSchema
type JobController struct {
	baseController BaseController[schema.Job]
}

// NewJobController initializes a JobController
func NewJobController() JobController {
	return JobController{
		baseController: BaseController[schema.Job]{
			collectionName: "jobs",
			displayName:    "Job",
		},
	}
}

// Query godoc
// @Summary Query jobs with filters
// @Description Get jobs by filtering fields like title, location, salary, company, etc.
// @Tags jobs
// @Accept  json
// @Produce  json
// @Param id query string false "Job ID (ObjectID)"
// @Param title query string false "Title (regex match)"
// @Param companyID query string false "Company ID (ObjectID)"
// @Param location query string false "Location (regex match)"
// @Param minSalary query integer false "Minimum salary"
// @Param maxSalary query integer false "Maximum salary"
// @Param workType query string false "Work type (e.g., Full-time, Part-time)"
// @Param workArrangement query string false "Work arrangement (e.g., Remote, On-site)"
// @Param postOpenDate query string false "Post open date (1d or 6w)"
// @Param latest query bool false "If true, returns the latest 3 jobs"
// @Param sort query string false "Sorting: dateAsc | dateDesc | title"
// @Success 200 {array} schema.Job
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/query [get]
func (jc JobController) Query(c *gin.Context) {
	// Allowed query params
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
		"workType": func(v string) (interface{}, error) {
			if v == "" {
				return nil, nil
			}
			return bson.M{"$regex": v, "$options": "i"}, nil
		},
		"workArrangement": func(v string) (interface{}, error) {
			if v == "" {
				return nil, nil
			}
			return bson.M{"$regex": v, "$options": "i"}, nil
		},
		"postOpenDate": func(v string) (interface{}, error) {
			if v == "" {
				return nil, nil
			}

			now := time.Now()
			midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

			switch v {
			case "1d":
				return bson.M{"$gte": midnight.AddDate(0, 0, -1)}, nil
			case "6w":
				return bson.M{"$gte": midnight.AddDate(0, 0, -7*6)}, nil
			default:
				return nil, nil
			}
		},
		"latest": func(v string) (interface{}, error) {
			if v == "true" {
				return true, nil
			}
			return nil, nil
		},
		"sort": func(v string) (interface{}, error) {
			switch v {
			case "dateAsc", "dateDesc", "title":
				return v, nil
			case "", "null":
				return nil, nil
			default:
				return nil, fmt.Errorf("invalid sort value")
			}
		},
	}

	filter := bson.M{}
	salaryFilter := bson.M{}
	findOptions := options.Find()

	// Apply query params
	for key, values := range c.Request.URL.Query() {
		if key == "sort" || key == "latest" {
			continue
		}
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
				case "postOpenDate":
					filter["postOpenDate"] = val
				case "latest":
					now := time.Now()
					filter["postOpenDate"] = bson.M{"$lte": now}
					findOptions.SetSort(bson.D{{Key: "postOpenDate", Value: -1}})
					findOptions.SetLimit(3)
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

	// Handle latest parameter
	latestParam := c.Query("latest")
	if latestParam == "true" {
		now := time.Now()
		filter["postOpenDate"] = bson.M{"$lte": now}
		findOptions.SetSort(bson.D{{Key: "postOpenDate", Value: -1}})
		findOptions.SetLimit(3)
	}

	// Handle sort parameter
	sortParam := c.Query("sort")
	switch sortParam {
	case "dateAsc":
		findOptions.SetSort(bson.D{{Key: "postOpenDate", Value: 1}})
	case "dateDesc":
		findOptions.SetSort(bson.D{{Key: "postOpenDate", Value: -1}})
	case "title":
		findOptions.SetSort(bson.D{{Key: "title", Value: 1}})
	case "", "null":
		// no sorting
	default:
		// ignore invalid sort
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobs, err := repository.FindAll[schema.Job](ctx, jc.baseController.collectionName, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(jobs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No jobs found"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// small helper
func toInt(s string) int {
	var n int
	_, _ = fmt.Sscan(s, &n)
	return n
}

// Create godoc
// @Summary Create a new job
// @Description Add a new job posting to the database
// @Tags jobs
// @Accept  json
// @Produce  json
// @Param job body schema.Job true "Job data"
// @Success 201 {object} schema.Job
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/ [post]
func (jc JobController) Create(c *gin.Context) {
	jc.baseController.Create(c)
}

// RetrieveAll godoc
// @Summary Get all jobs
// @Description Retrieve all job postings
// @Tags jobs
// @Produce  json
// @Success 200 {array} schema.Job
// @Failure 500 {object} map[string]string
// @Router /jobs/ [get]
func (jc JobController) RetrieveAll(c *gin.Context) {
	jc.baseController.RetrieveAll(c)
}

// Update godoc
// @Summary Update a job
// @Description Update a job posting by ID
// @Tags jobs
// @Accept  json
// @Produce  json
// @Param id path string true "Job ID"
// @Param job body schema.Job true "Updated job data"
// @Success 200 {object} schema.Job
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/{id} [put]
func (jc JobController) Update(c *gin.Context) {
	jc.baseController.Update(c)
}

// Delete godoc
// @Summary Delete a job
// @Description Delete a job posting by ID and notify applicants
// @Tags jobs
// @Accept  json
// @Produce  json
// @Param id path string true "Job ID"
// @Param reason body object false "Reason for deletion"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/{id} [delete]
func (jc JobController) Delete(c *gin.Context) {
	shouldReturn := notifyJobDeletion(c)
	if shouldReturn {
		return
	}
	jc.baseController.Delete(c)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	repository.DeleteMany[schema.JobApplication](ctx, "job_applications", bson.M{"jobID": bson.M{"$eq": id}})
}

// notifyJobDeletion send emails to all applicants when a job they applied to got deleted.
func notifyJobDeletion(c *gin.Context) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Job ID"})
		return true
	}

	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return true
	}
	reason := body.Reason
	if reason == "" {
		reason = "No reason provided."
	}

	job, err := repository.FindOne[schema.Job](ctx, "jobs", jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No job Found"})
		return true
	}

	filter := bson.M{"jobID": bson.M{"$eq": job.ID}}
	jobApplications, err := repository.FindAll[schema.JobApplication](ctx, "job_applications", filter)
	if err == mongo.ErrNoDocuments {
		return false
	}
	if err != mongo.ErrNoDocuments && err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Problems when finding job applications"})
		return true
	}

	visited := make(map[primitive.ObjectID]bool)
	var applicantIDs []primitive.ObjectID
	for _, jobApp := range jobApplications {
		if !visited[jobApp.ApplicantID] {
			visited[jobApp.ApplicantID] = true
			applicantIDs = append(applicantIDs, jobApp.ApplicantID)
		}
	}

	if len(applicantIDs) == 0 {
		return false
	}

	applicants, err := getUsersFromID(ctx, applicantIDs)
	if err == mongo.ErrNoDocuments {
		return false
	}
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return true
	}

	for _, applicant := range applicants {
		emailBody := fmt.Sprintf(
			"The job '%s' you applied for has been deleted.\n\nReason: %s",
			job.Title,
			reason,
		)
		email.Send(applicant.Email, "Job Deletion Notice", emailBody)
	}

	return false
}

// RetrieveOne godoc
// @Summary Get a job by ID
// @Description Retrieve details of a specific job posting
// @Tags jobs
// @Produce  json
// @Param id path string true "Job ID"
// @Success 200 {object} schema.Job
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /jobs/{id} [get]
func (jc JobController) RetrieveOne(c *gin.Context) {
	jc.baseController.RetrieveOne(c)
}
