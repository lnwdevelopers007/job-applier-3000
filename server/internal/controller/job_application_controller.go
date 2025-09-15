package controller


// JobApplicationController handles JobApplication CRUD operations
type JobApplicationController struct {
	collectionName string
}

// NewJobApplicationController initializes a JobApplicationController
func NewJobApplicationController() JobApplicationController {
	return JobApplicationController{collectionName: "job_applications"}
}

