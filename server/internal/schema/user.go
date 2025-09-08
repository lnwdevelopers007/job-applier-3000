package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string             `bson:"userID"`
	Provider  string             `bson:"provider"`
	Email     string             `bson:"email"`
	Name      string             `bson:"name"`
	AvatarURL string             `bson:"avatarURL"`
	UpdatedAt time.Time          `bson:"updatedAt"`
	CreatedAt time.Time          `bson:"createdAt"`
}
