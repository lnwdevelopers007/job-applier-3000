package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type Company struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name" binding:"required"`
	UserID      primitive.ObjectID `bson:"userID" json:"userID" binding:"required"`
	AboutUs     string             `bson:"aboutUs" json:"aboutUs" binding:"required"`
	CompanyType string             `bson:"companyType" json:"companyType" binding:"required"`
}
