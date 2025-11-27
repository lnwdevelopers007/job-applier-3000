package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type CompanyInfo struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID       primitive.ObjectID `bson:"userID" json:"userID" binding:"required"`
	Name         string             `bson:"name" json:"name" binding:"required,min=2,max=200"`
	AboutUs      string             `bson:"aboutUs" json:"aboutUs" binding:"required,min=10,max=5000"`
	Industry     string             `bson:"industry" json:"industry" binding:"required,min=2,max=100"`
	Size         string             `bson:"size" json:"size" binding:"required"`
	Website      string             `bson:"website" json:"website" binding:"required,url,max=200"`
	Logo         string             `bson:"logo,omitempty" json:"logo,omitempty" binding:"omitempty,url,max=500"`
	FoundedYear  string             `bson:"foundedYear,omitempty" json:"foundedYear,omitempty" binding:"omitempty,len=4,numeric"`
	Headquarters string             `bson:"headquarters,omitempty" json:"headquarters,omitempty" binding:"omitempty,min=2,max=200"`
	LinkedIn     string             `bson:"linkedIn,omitempty" json:"linkedIn,omitempty" binding:"omitempty,url,max=200"`
}
