package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title               string             `bson:"title" json:"title" binding:"required,gte=0"`
	CompanyID           primitive.ObjectID `bson:"companyID" json:"companyID" binding:"required"`
	Location            string             `bson:"location" json:"location" binding:"required,gte=0"`
	WorkType            string             `bson:"workType" json:"workType" binding:"required,gte=0"`
	WorkArrangement     string             `bson:"workArrangement" json:"workArrangement" binding:"required,gte=0"`
	Currency            string             `bson:"currency" json:"currency" binding:"required,gte=0"`
	MinSalary           float64            `bson:"minSalary" json:"minSalary" binding:"required,gte=0"`
	MaxSalary           float64            `bson:"maxSalary" json:"maxSalary" binding:"required,gte=0"`
	JobDescription      string             `bson:"jobDescription" json:"jobDescription" binding:"required,gte=0"`
	JobSummary          string             `bson:"jobSummary" json:"jobSummary" binding:"required,gte=0"`
	RequiredSkills      string             `bson:"requiredSkills" json:"requiredSkills" binding:"required,gte=0"`
	ExperienceLevel     string             `bson:"experienceLevel" json:"experienceLevel" binding:"required,gte=0"`
	Education           string             `bson:"education" json:"education" binding:"required,gte=0"`
	NiceToHave          string             `bson:"niceToHave" json:"niceToHave"`
	Questions           string             `bson:"questions" json:"questions"`
	ApplicationDeadline time.Time          `bson:"applicationDeadline" json:"applicationDeadline" binding:"required"`
	NumberOfPositions   int                `bson:"numberOfPositions" json:"numberOfPositions" binding:"required"`
	Visibility          string             `bson:"visibility" json:"visibility" binding:"required,gte=0"`
	EmailNotifications  bool               `bson:"emailNotifications" json:"emailNotifications" binding:"boolean"`
	AutoReject          bool               `bson:"autoReject" json:"autoReject" binding:"boolean"`
}
