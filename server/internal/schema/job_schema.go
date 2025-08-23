// Package database handles all connection between web server and mongodb.
package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

// Schema for various datatypes.
// more info here: https://gin-gonic.com/en/docs/examples/binding-and-validation/

type JobSchema struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title    string             `bson:"title" json:"title" binding:"required"`
	Company  string             `bson:"company" json:"company"`
	Location string             `bson:"location" json:"location"`
	Salary   int                `bson:"salary" json:"salary" binding:"required,gte=0"`
}
