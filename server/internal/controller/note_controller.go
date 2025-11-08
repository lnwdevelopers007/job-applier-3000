package controller

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/repository"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
)

type NoteController struct {
	baseController BaseController[schema.Note]
}

func NewNoteController() NoteController {
	return NoteController{
		baseController: BaseController[schema.Note]{
			collectionName: "notes",
			displayName:    "Note",
		},
	}
}

func (nc NoteController) validateNoteOwner(c *gin.Context) (shouldReturn bool) {
	enableAuth, _ := strconv.ParseBool(os.Getenv("ENABLE_AUTH"))
	if !enableAuth {
		return false
	}

	var raw schema.Note
	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect Note Schema"})
		return true
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobApplication, err := repository.FindOne[schema.JobApplication](
		ctx,
		"job_applications",
		raw.JobApplicationID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot find Job Application"})
		return true
	}

	job, err := repository.FindOne[schema.Job](
		ctx,
		"jobs",
		jobApplication.JobID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot find Job"})
		return true
	}

	middlewareUserID, _, err := getUserFromMiddleware(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot get user from middleware"})
		return true
	}

	if middlewareUserID != job.CompanyID {
		userMismatcherdErr := errors.New("the requested user and the job poster are not the same person")
		c.JSON(http.StatusBadRequest, gin.H{"error": userMismatcherdErr.Error()})
		return true
	}

	return false
}

func (nc NoteController) Create(c *gin.Context) {
	shouldReturn := nc.validateNoteOwner(c)
	if shouldReturn {
		return
	}
	nc.baseController.Create(c)
}

func (nc NoteController) Delete(c *gin.Context) {
	shouldReturn := nc.validateNoteOwner(c)
	if shouldReturn {
		return
	}
	nc.baseController.Delete(c)
}

func (nc NoteController) Update(c *gin.Context) {
	shouldReturn := nc.validateNoteOwner(c)
	if shouldReturn {
		return
	}
	nc.baseController.Update(c)
}

func (nc NoteController) RetrieveAll(c *gin.Context) {
	nc.baseController.RetrieveAll(c)
}

func (nc NoteController) RetrieveOne(c *gin.Context) {
	nc.baseController.RetrieveOne(c)
}
