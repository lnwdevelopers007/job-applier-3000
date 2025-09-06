package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	JobApplicationID primitive.ObjectID `bson:"jobApplicationID" json:"jobApplicationID" binding:"required"`
	Content          string             `bson:"name" json:"name" binding:"required"`
	Timestamp        time.Time          `bson:"timestamp" json:"timestamp" binding:"required"`
}
