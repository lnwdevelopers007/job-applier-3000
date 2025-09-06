package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

type File struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ParentID      primitive.ObjectID `bson:"parentID" json:"parentID" binding:"required"`
	ParentColl    string             `bson:"parentColl" json:"parentColl" binding:"required"`
	Content       primitive.ObjectID `bson:"content" json:"content" binding:"required"`
	FileExtension string             `bson:"fileExtension" json:"fileExtension" binding:"required"`
}
