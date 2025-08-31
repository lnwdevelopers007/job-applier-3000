package schema

import (
	"time"
)

type User struct {
	ID        string    `bson:"_id"`
	Provider  string    `bson:"provider"`
	Email     string    `bson:"email"`
	Name      string    `bson:"name"`
	AvatarURL string    `bson:"avatarURL"`
	UpdatedAt time.Time `bson:"updatedAt"`
	CreatedAt time.Time `bson:"createdAt"`
}
