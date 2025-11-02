package schema

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileCategory string

const (
	CategoryResume        FileCategory = "resume"
	CategoryTranscript   FileCategory = "transcript"
	CategoryCertification FileCategory = "certification"
	CategoryVerification  FileCategory = "verification"
)

type File struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID        primitive.ObjectID `bson:"userID" json:"userID" binding:"required"`
	Content       []byte             `bson:"content" json:"-"` // Don't expose in JSON
	FileExtension string             `bson:"fileExtension" json:"fileExtension" binding:"required"`
	Filename      string             `bson:"filename" json:"filename" binding:"required"`
	ContentType   string             `bson:"contentType" json:"contentType" binding:"required"`
	Size          int64              `bson:"size" json:"size" binding:"required"`
	Category      FileCategory       `bson:"category" json:"category" binding:"required"`
	UploadDate    time.Time          `bson:"uploadDate" json:"uploadDate"`
}