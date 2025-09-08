package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobApplication struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ApplicantID primitive.ObjectID `bson:"applicantID" json:"applicantID" binding:"required"`
	JobID       primitive.ObjectID `bson:"jobID" json:"jobID" binding:"required"`
	CompanyID   primitive.ObjectID `bson:"companyID" json:"companyID" binding:"required"`
	Status      string             `bson:"status" json:"status" binding:"required"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt" binding:"required"`
}
