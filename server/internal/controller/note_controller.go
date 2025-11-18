package controller

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/dto"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/repository"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteController struct {
	baseController BaseController[schema.Note, dto.Note]
}

func NewNoteController() NoteController {
	return NoteController{
		baseController: BaseController[schema.Note, dto.Note]{
			collectionName: "notes",
			displayName:    "Note",
		},
	}
}

// getMiddlewareUserID get userID from middleware
func getMiddlewareUserID(c *gin.Context) (primitive.ObjectID, error) {
	middlewareUserID, _, err := getUserFromMiddleware(c)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return middlewareUserID, nil
}

// getJobAndApplication to get JobApplication and its parent Job
func getJobAndApplication(ctx context.Context, jobAppID primitive.ObjectID) (schema.JobApplication, schema.Job, error) {
	jobApplication, err := repository.FindOne[schema.JobApplication](ctx, jobAppID)
	if err != nil {
		return schema.JobApplication{}, schema.Job{}, err
	}
	job, err := repository.FindOne[schema.Job](ctx, jobApplication.JobID)
	if err != nil {
		return schema.JobApplication{}, schema.Job{}, err
	}
	return jobApplication, job, nil
}

// helper to get all jobs owned by a company user
func getJobsByCompanyID(ctx context.Context, companyID primitive.ObjectID) ([]schema.Job, error) {
	return repository.FindAll[schema.Job](ctx, bson.M{"companyID": companyID})
}

// validateNoteOwner ensures the requesting user is the job owner
func (nc NoteController) validateNoteOwner(c *gin.Context) (shouldReturn bool) {
	enableAuth, _ := strconv.ParseBool(os.Getenv("ENABLE_AUTH"))
	if !enableAuth {
		return false
	}

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return true
	}

	// Fetch existing note to get jobApplicationID
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingNote, err := repository.FindOne[schema.Note](ctx, objID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Note not found"})
		return true
	}

	jobApplication, job, err := getJobAndApplication(ctx, existingNote.JobApplicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot find Job or Job Application"})
		return true
	}

	middlewareUserID, err := getMiddlewareUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot get user from middleware"})
		return true
	}

	if middlewareUserID != job.CompanyID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "the requested user and the job poster are not the same person"})
		return true
	}

	_ = jobApplication
	return false
}

// Create godoc
// @Summary      Create a note
// @Description  Create a new note to a specific job application.
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        Note  body      schema.Note  true  "Note object"
// @Success      201   {object}  schema.Note
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notes/ [post]
func (nc NoteController) Create(c *gin.Context) {
	shouldReturn := nc.validateNoteOwner(c)
	if shouldReturn {
		return
	}
	nc.baseController.Create(c)
}

// Delete godoc
// @Summary      Delete a note
// @Description  Delete a note.
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Note ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id} [delete]
func (nc NoteController) Delete(c *gin.Context) {
	shouldReturn := nc.validateNoteOwner(c)
	if shouldReturn {
		return
	}
	nc.baseController.Delete(c)
}

// Update godoc
// @Summary      Update an existing Note
// @Description  Modify note data by providing the note ID in the path.
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "Note ID"
// @Param        note  body      schema.Note  true  "Updated note data"
// @Success      200   {object}  schema.Note
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notes/{id} [put]
func (nc NoteController) Update(c *gin.Context) {
	shouldReturn := nc.validateNoteOwner(c)
	if shouldReturn {
		return
	}
	nc.baseController.Update(c)
}

// RetrieveOne godoc
// @Summary      Get a ntoe by ID
// @Description  Retrieve a single note document by its ID.
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Note ID"
// @Success      200  {object}  schema.Note
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id} [get]
func (nc NoteController) RetrieveOne(c *gin.Context) {
	nc.baseController.RetrieveOne(c)
}

// Query godoc
// @Summary      RETRIEVES all notes linked to all jobs owned by user + query by jobApplicationID.
// @Description  RETRIEVES all notes linked to all jobs owned by the authenticated company user, and QUERY by jobApplicationID.
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param        jobApplicationID  query     string  false  "Filter by a specific job application ID"
// @Success      200   {array}   schema.Note
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notes/ [get]
func (nc NoteController) Query(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	middlewareUserID, err := getMiddlewareUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot get user from middleware"})
		return
	}

	jobAppParam := c.Query("jobApplicationID")
	var jobAppFilter bson.M

	// 1. Find all jobs owned by this company
	jobs, err := getJobsByCompanyID(ctx, middlewareUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot fetch jobs"})
		return
	}
	if len(jobs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No jobs found"})
		return
	}

	jobIDs := make([]primitive.ObjectID, 0, len(jobs))
	for _, job := range jobs {
		jobIDs = append(jobIDs, job.ID)
	}

	// 2. Find job applications linked to these jobs
	if jobAppParam != "" {
		objID, err := primitive.ObjectIDFromHex(jobAppParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jobApplicationID"})
			return
		}
		jobAppFilter = bson.M{"_id": objID, "jobID": bson.M{"$in": jobIDs}}
	} else {
		jobAppFilter = bson.M{"jobID": bson.M{"$in": jobIDs}}
	}

	jobApplications, err := repository.FindAll[schema.JobApplication](ctx, jobAppFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot fetch job applications"})
		return
	}
	if len(jobApplications) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No job applications found"})
		return
	}

	jobAppIDs := make([]primitive.ObjectID, 0, len(jobApplications))
	for _, app := range jobApplications {
		jobAppIDs = append(jobAppIDs, app.ID)
	}

	// 3. Find all notes linked to those job applications
	notes, err := repository.FindAll[schema.Note](ctx, bson.M{
		"jobApplicationID": bson.M{"$in": jobAppIDs},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot fetch notes"})
		return
	}

	c.JSON(http.StatusOK, notes)
}
