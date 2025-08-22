// Package database handles all connection between web server and mongodb.
package database

import "go.mongodb.org/mongo-driver/bson/primitive"

// Schema for various datatypes.
// We don't need to put "All" fields inside here.
// Just the ones that requires validation.

type JobSchema struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title    string             `bson:"title" json:"title" binding:"required"`
	Company  string             `bson:"company" json:"company"`
	Location string             `bson:"location" json:"location"`
	Salary   int                `bson:"salary" json:"salary" binding:"required,gte=0"`
}
