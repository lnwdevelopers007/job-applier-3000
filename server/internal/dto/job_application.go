package dto

import (
	"time"
)

type JobApplication struct {
	Status    *string    `bson:"status,omitempty" json:"status,omitempty"`
	CreatedAt *time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}
