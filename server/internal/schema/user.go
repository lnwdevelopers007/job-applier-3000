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
	UserID    string             `bson:"userID"`
	Provider  string             `bson:"provider"`
	Email     string             `bson:"email"`
	Name      string             `bson:"name"`
	AvatarURL string             `bson:"avatarURL"`
	Role      string             `bson:"role"`
	Verified  bool               `bson:"verified"`
	UpdatedAt time.Time          `bson:"updatedAt"`
	CreatedAt time.Time          `bson:"createdAt"`
	UserInfo  any                `bson:"userInfo"`
}
