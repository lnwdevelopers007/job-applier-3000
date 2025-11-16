package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type CompanyInfo struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID       primitive.ObjectID `bson:"userID" json:"userID" binding:"required"`
	Name         string             `bson:"name" json:"name" binding:"required"`
	AboutUs      string             `bson:"aboutUs" json:"aboutUs" binding:"required"`
	Industry     string             `bson:"industry" json:"industry" binding:"required"`
	Size         string             `bson:"size" json:"size" binding:"required"`
	Website      string             `bson:"website" json:"website" binding:"required"`
	Logo         string             `bson:"logo,omitempty" json:"logo,omitempty"`
	FoundedYear  string             `bson:"foundedYear,omitempty" json:"foundedYear,omitempty"`
	Headquarters string             `bson:"headquarters,omitempty" json:"headquarters,omitempty"`
	LinkedIn     string             `bson:"linkedIn,omitempty" json:"linkedIn,omitempty"`
}
