package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the struct representing user (both company and jobSeeker).
// UserInfo field can be either CompanyInfo or JobSeekerInfo, but due to its implementation
// requireing the use of generic, which will break everything in the code,
// I'll leave it as any type.
// field validation to be added (in case we support non-google oauth and things go wrong).
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string             `bson:"userID,omitempty" json:"userID,omitempty"`
	Provider  string             `bson:"provider" json:"provider"`
	Email     string             `bson:"email,omitempty" json:"email,omitempty"`
	Name      string             `bson:"name,omitempty" json:"name,omitempty"`
	AvatarURL string             `bson:"avatarURL,omitempty" json:"avatarURL,omitempty"`
	Role      string             `bson:"role,omitempty" json:"role,omitempty"`
	Verified  bool               `bson:"verified,omitempty" json:"verified,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UserInfo  any                `bson:"userInfo,omitempty" json:"userInfo,omitempty"`
}
