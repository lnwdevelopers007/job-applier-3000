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

// Query godoc
// @Summary      Query job applications
// @Description  Get all job applications that match query parameters
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        applicantID query string false "Applicant ID"
// @Param        jobID query string false "Job ID"
// @Param        companyID query string false "Company ID"
// @Param        status query string false "Status of the application"
// @Success      200  {array}  schema.ApplicationWithApplicant
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /apply/query [get]
func (jc JobApplicationController) Query(c *gin.Context) {

	jobApplicationFilter, shouldReturn := jobApplicationFilter(c)
	if shouldReturn {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	findOpts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	applications, err := repository.FindAll[schema.JobApplication](ctx, jobApplicationFilter, findOpts)
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
	userMap, err := getUsersFromIDs(
		ctx,
		userIDs,
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

func jobApplicationFilter(c *gin.Context) (bson.M, bool) {
	allowedParams := map[string]func(string) (any, error){
		"id": func(v string) (any, error) {
			if v == "" {
				return nil, fmt.Errorf("id parameter is empty")
			}
			return primitive.ObjectIDFromHex(v)
		},
		"applicantID": func(v string) (any, error) {
			if v == "" {
				return nil, fmt.Errorf("applicantID parameter is empty")
			}
			objID, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				return nil, err
			}
			// Try both ObjectID and string format for backwards compatibility
			return bson.M{"$in": []any{objID, v}}, nil
		},
		"jobID": func(v string) (any, error) {
			if v == "" {
				return nil, fmt.Errorf("jobID parameter is empty")
			}
			objID, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				return nil, err
			}
			// Try both ObjectID and string format for backwards compatibility
			return bson.M{"$in": []any{objID, v}}, nil
		},
		"companyID": func(v string) (any, error) {
			if v == "" {
				return nil, fmt.Errorf("companyID parameter is empty")
			}
			objID, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				return nil, err
			}
			// Try both ObjectID and string format for backwards compatibility
			return bson.M{"$in": []any{objID, v}}, nil
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

// Create godoc
// @Summary      Create new job application
// @Description  Submit a job application for a specific job
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        request body schema.JobApplication true "Job Application JSON"
// @Success      200  {object} schema.JobApplication
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /apply/ [post]
func (jc JobApplicationController) Create(c *gin.Context) {
	companyEmail, err := jc.shouldNotifyCompany(c)
	if err {
		return
	}
	jc.baseController.Create(c)

	if companyEmail != "" {
		var raw schema.JobApplication
		if err := c.ShouldBindBodyWithJSON(&raw); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			applicantID := raw.ApplicantID
			jobID := raw.JobID

			var validApplicantID, validJobID primitive.ObjectID
			var convErr error
			// Validate applicant ID
			validApplicantID, convErr = primitive.ObjectIDFromHex(applicantID.Hex())
			if convErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ApplicantID format"})
				return
			}
			// Validate job ID
			validJobID, convErr = primitive.ObjectIDFromHex(jobID.Hex())
			if convErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JobID format"})
				return
			}

			applicant, _ := repository.FindOne[schema.User](ctx, validApplicantID)
			job, _ := repository.FindOne[schema.Job](ctx, validJobID)

			subject := "New applicant applied to your job"
			body := fmt.Sprintf(
				"Hello,\n\n%s has applied to your job \"%s\".\n\nPlease review the application in your applicant board.\n\nBest regards,\nJob Applier 3000",
				applicant.Name,
				job.Title,
			)

			email.Send(companyEmail, subject, body)
		}
	}
}

// shouldNotifyCompany determines whether company should be notified when an applicant applied for a job or not.
func (jc JobApplicationController) shouldNotifyCompany(c *gin.Context) (companyEmail string, err bool) {
	var raw schema.JobApplication
	if err := c.ShouldBindBodyWithJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", true
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	job, jobErr := repository.FindOne[schema.Job](ctx, raw.JobID)
	if jobErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": jobErr.Error()})
		return "", true
	}
	if !job.EmailNotifications {
		return "", false
	}
	company, companyErr := repository.FindOne[schema.User](ctx, job.CompanyID)
	if companyErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": companyErr.Error()})
		return "", true
	}
	return company.Email, false
}

// Update godoc
// @Summary      Update job application
// @Description  Update job application data (including status changes)
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        id path string true "Application ID"
// @Param        request body schema.JobApplication true "Updated Application JSON"
// @Success      200  {object} schema.JobApplication
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /apply/{id} [put]
func (jc JobApplicationController) Update(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	jc.baseController.Update(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updatedApp, err := repository.FindOne[schema.JobApplication](ctx, objID)
	if err != nil {
		fmt.Println("Failed to fetch updated application for notification:", err)
		return
	}

	jc.notifyApplicantOnStatusChange(ctx, updatedApp)
}

// notifyApplicantOnStatusChange notifies the applicant when their application status changes
func (jc JobApplicationController) notifyApplicantOnStatusChange(ctx context.Context, app schema.JobApplication) {
	applicant, err := repository.FindOne[schema.User](ctx, app.ApplicantID)
	if err != nil {
		fmt.Println("Failed to fetch applicant for notification:", err)
		return
	}

	job, err := repository.FindOne[schema.Job](ctx, app.JobID)
	if err != nil {
		fmt.Println("Failed to fetch job for notification:", err)
		return
	}

	var subject, body string

	switch app.Status {
	case "ACCEPTED":
		subject = "Congratulations! Your job application has been accepted"
		body = fmt.Sprintf(
			"Hello %s,\n\nWe are pleased to inform you that your application for the job \"%s\" has been ACCEPTED.\n\nOur team will contact you soon with the next steps.\n\nBest regards,\nJob Applier 3000",
			applicant.Name,
			job.Title,
		)
	case "REJECTED":
		subject = "Update on your job application"
		body = fmt.Sprintf(
			"Hello %s,\n\nWe regret to inform you that your application for the job \"%s\" has been REJECTED.\n\nWe appreciate your interest and encourage you to apply for future opportunities.\n\nBest regards,\nJob Applier 3000",
			applicant.Name,
			job.Title,
		)
	default:
		subject = "Your job application status has been updated"
		body = fmt.Sprintf(
			"Hello %s,\n\nThe status of your application for the job \"%s\" has been updated to: %s\n\nBest regards,\nJob Applier 3000",
			applicant.Name,
			job.Title,
			app.Status,
		)
	}

	email.Send(applicant.Email, subject, body)
}

// Delete godoc
// @Summary      Delete a job application
// @Description  Remove a job application by its ID
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Application ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /apply/{id} [delete]
func (jc JobApplicationController) Delete(c *gin.Context) {
	jc.baseController.Delete(c)
}

// RetrieveAll godoc
// @Summary      Get all job applications
// @Description  Retrieve all job applications in the system
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Success      200  {array}   schema.JobApplication
// @Failure      500  {object}  map[string]string
// @Router       /apply/ [get]
func (jc JobApplicationController) RetrieveAll(c *gin.Context) {
	jc.baseController.RetrieveAll(c)
}

// RetrieveOne godoc
// @Summary      Get a specific job application
// @Description  Retrieve a single job application by its ID
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Application ID"
// @Success      200  {object}  schema.JobApplication
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /apply/{id} [get]
func (jc JobApplicationController) RetrieveOne(c *gin.Context) {
	jc.baseController.RetrieveOne(c)
}
