package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type JobSeeker struct {
	ID     primitive.ObjectID
	UserID primitive.ObjectID
}

type ContactInfo struct {
	Location string
	Phone    string
	LinkedIn string
}
