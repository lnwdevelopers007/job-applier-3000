package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type JobSeeker struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID  primitive.ObjectID `bson:"userID" json:"userID" binding:"required"`
	Contact ContactInfo        `bson:"contact" json:"contact" binding:"required"`
}

type ContactInfo struct {
	Location string `bson:"location" json:"location" binding:"required"`
	Phone    string `bson:"phone" json:"phone" binding:"required"`
	LinkedIn string `bson:"linkedIn" json:"linkedIn" binding:"required"`
}

type JobSeekerDenormalised struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	User    User               `bson:"user" json:"user"`
	Contact ContactInfo        `bson:"contact" json:"contact" binding:"required"`
}
