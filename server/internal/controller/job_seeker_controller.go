package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
)

type JobSeekerController struct {
	baseController BaseController[schema.JobSeeker]
}

func NewJobSeekerController() JobSeekerController {
	return JobSeekerController{
		baseController: BaseController[schema.JobSeeker]{
			collectionName: "jobSeeker",
			displayName:    "Job Seeker",
		},
	}
}

func (jc JobSeekerController) Create(c *gin.Context) {
	jc.baseController.Create(c)
}

// Update() updates a job by ID.
func (jc JobSeekerController) Update(c *gin.Context) {
	jc.baseController.Update(c)
}

// Delete() deletes a job by ID.
func (jc JobSeekerController) Delete(c *gin.Context) {
	jc.baseController.Delete(c)
}

// RetrieveOne fetches a single job seeker by ID.
func (jc JobSeekerController) RetrieveOne(c *gin.Context) {
	jc.baseController.RetrieveOne(c)
}