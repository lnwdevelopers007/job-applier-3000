package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Provider     string             `bson:"provider" json:"provider"`
	ProviderID   string             `bson:"provider_id" json:"provider_id"`
	Email        string             `bson:"email" json:"email"`
	Name         string             `bson:"name" json:"name"`
	AvatarURL    string             `bson:"avatar_url" json:"avatar_url"`
	AccessToken  string             `bson:"access_token,omitempty" json:"-"`
	RefreshToken string             `bson:"refresh_token,omitempty" json:"-"`
	ExpiresAt    *time.Time         `bson:"expires_at,omitempty" json:"-"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
