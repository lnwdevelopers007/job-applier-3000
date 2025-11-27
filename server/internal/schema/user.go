package schema

import (
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the struct representing user (both company and jobSeeker).
// UserInfo field can be either CompanyInfo or JobSeekerInfo, but due to its implementation
// requireing the use of generic, which will break everything in the code,
// I'll leave it as any type.
// field validation to be added (in case we support non-google oauth and things go wrong).
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string             `bson:"userID,omitempty" json:"userID,omitempty" binding:"omitempty,max=100"`
	Provider  string             `bson:"provider" json:"provider"`
	Email     string             `bson:"email,omitempty" json:"email,omitempty" binding:"omitempty,email,max=255"`
	Name      string             `bson:"name,omitempty" json:"name,omitempty" binding:"omitempty,min=1,max=200"`
	AvatarURL string             `bson:"avatarURL,omitempty" json:"avatarURL,omitempty" binding:"omitempty,url,max=500"`
	Role      string             `bson:"role,omitempty" json:"role,omitempty" binding:"omitempty,oneof=jobSeeker company faculty admin"`
	Verified  bool               `bson:"verified" json:"verified"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UserInfo  bson.M             `bson:"userInfo,omitempty" json:"userInfo,omitempty"`
	Banned    bool               `bson:"banned,omitempty" json:"banned,omitempty"`
}

func (u User) GetCollectionName() string {
	return "users"
}
