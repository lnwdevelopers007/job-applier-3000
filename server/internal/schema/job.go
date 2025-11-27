package schema

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title               string             `bson:"title" json:"title" binding:"required,min=3,max=200"`
	CompanyID           primitive.ObjectID `bson:"companyID" json:"companyID" binding:"required"`
	Location            string             `bson:"location" json:"location" binding:"required,min=2,max=200"`
	WorkType            string             `bson:"workType" json:"workType" binding:"required,gte=0"`
	WorkArrangement     string             `bson:"workArrangement" json:"workArrangement" binding:"required,gte=0"`
	Currency            string             `bson:"currency" json:"currency" binding:"required,len=3"`
	MinSalary           float64            `bson:"minSalary" json:"minSalary" binding:"required,gte=0,lte=1000000000"`
	MaxSalary           float64            `bson:"maxSalary" json:"maxSalary" binding:"required,gte=0,lte=1000000000"`
	JobDescription      string             `bson:"jobDescription" json:"jobDescription" binding:"required,min=50,max=10000"`
	JobSummary          string             `bson:"jobSummary" json:"jobSummary" binding:"required,min=20,max=500"`
	RequiredSkills      string             `bson:"requiredSkills" json:"requiredSkills" binding:"required,min=5,max=2000"`
	ExperienceLevel     string             `bson:"experienceLevel" json:"experienceLevel" binding:"required,,gte=0"`
	Education           string             `bson:"education" json:"education" binding:"required,min=2,max=200"`
	NiceToHave          string             `bson:"niceToHave" json:"niceToHave" binding:"omitempty,max=2000"`
	Questions           string             `bson:"questions" json:"questions" binding:"omitempty,max=2000"`
	PostOpenDate        time.Time          `bson:"postOpenDate" json:"postOpenDate" binding:"required"`
	ApplicationDeadline time.Time          `bson:"applicationDeadline" json:"applicationDeadline" binding:"required"`
	NumberOfPositions   int                `bson:"numberOfPositions" json:"numberOfPositions" binding:"required,gte=1,lte=1000"`
	Visibility          string             `bson:"visibility" json:"visibility" binding:"required,gte=0"`
	EmailNotifications  bool               `bson:"emailNotifications" json:"emailNotifications" binding:"boolean"`
	AutoReject          bool               `bson:"autoReject" json:"autoReject" binding:"boolean"`
}

func (j Job) GetCollectionName() string {
	return "jobs"
}

func (j *Job) Validate() error {
	// MinSalary must not exceed MaxSalary
	if j.MinSalary > j.MaxSalary {
		return fmt.Errorf("minSalary (%.2f) cannot be greater than maxSalary (%.2f)", j.MinSalary, j.MaxSalary)
	}

	// PostOpenDate must be before or equal to ApplicationDeadline
	if j.PostOpenDate.After(j.ApplicationDeadline) {
		return fmt.Errorf("postOpenDate (%s) cannot be after applicationDeadline (%s)", j.PostOpenDate, j.ApplicationDeadline)
	}
	return nil
}

func (j Job) ValidatePartial(fields map[string]any) error {
    // Validate min/max salary if present
    if min, ok := fields["minSalary"].(float64); ok {
        if min <= 0 || min > 1000000000 {
            return fmt.Errorf("minSalary must be between 0 and 1,000,000,000")
        }
        if max, ok2 := fields["maxSalary"].(float64); ok2 {
            if max <= 0 || max > 1000000000 {
                return fmt.Errorf("maxSalary must be between 0 and 1,000,000,000")
            }
            if min > max {
                return fmt.Errorf("minSalary cannot be greater than maxSalary")
            }
        }
    }

    // Validate dates if present
    if postStr, ok := fields["postOpenDate"].(string); ok {
        post, err := time.Parse(time.RFC3339, postStr)
        if err != nil {
            return fmt.Errorf("postOpenDate is invalid: %v", err)
        }
        if deadlineStr, ok2 := fields["applicationDeadline"].(string); ok2 {
            deadline, err2 := time.Parse(time.RFC3339, deadlineStr)
            if err2 != nil {
                return fmt.Errorf("applicationDeadline is invalid: %v", err2)
            }
            if post.After(deadline) {
                return fmt.Errorf("postOpenDate cannot be after applicationDeadline")
            }
        }
    }

    return nil
}

