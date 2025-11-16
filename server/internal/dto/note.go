package dto

import (
	"time"
)

type Note struct {
	Content   *string    `bson:"content,omitempty" json:"content,omitempty"`
	Timestamp *time.Time `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}
