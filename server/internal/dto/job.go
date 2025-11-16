package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID                  *primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title               *string             `bson:"title,omitempty" json:"title,omitempty"`
	Location            *string             `bson:"location,omitempty" json:"location,omitempty"`
	WorkType            *string             `bson:"workType,omitempty" json:"workType,omitempty"`
	WorkArrangement     *string             `bson:"workArrangement,omitempty" json:"workArrangement,omitempty"`
	Currency            *string             `bson:"currency,omitempty" json:"currency,omitempty"`
	MinSalary           *float64            `bson:"minSalary,omitempty" json:"minSalary,omitempty"`
	MaxSalary           *float64            `bson:"maxSalary,omitempty" json:"maxSalary,omitempty"`
	JobDescription      *string             `bson:"jobDescription,omitempty" json:"jobDescription,omitempty"`
	JobSummary          *string             `bson:"jobSummary,omitempty" json:"jobSummary,omitempty"`
	RequiredSkills      *string             `bson:"requiredSkills,omitempty" json:"requiredSkills,omitempty"`
	ExperienceLevel     *string             `bson:"experienceLevel,omitempty" json:"experienceLevel,omitempty"`
	Education           *string             `bson:"education,omitempty" json:"education,omitempty"`
	NiceToHave          *string             `bson:"niceToHave,omitempty" json:"niceToHave,omitempty"`
	Questions           *string             `bson:"questions,omitempty" json:"questions,omitempty"`
	PostOpenDate        *time.Time          `bson:"postOpenDate,omitempty" json:"postOpenDate,omitempty"`
	ApplicationDeadline *time.Time          `bson:"applicationDeadline,omitempty" json:"applicationDeadline,omitempty"`
	NumberOfPositions   *int                `bson:"numberOfPositions,omitempty" json:"numberOfPositions,omitempty"`
	Visibility          *string             `bson:"visibility,omitempty" json:"visibility,omitempty"`
	EmailNotifications  *bool               `bson:"emailNotifications,omitempty" json:"emailNotifications,omitempty"`
	AutoReject          *bool               `bson:"autoReject,omitempty" json:"autoReject,omitempty"`
}
