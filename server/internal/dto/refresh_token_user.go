package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshTokenUser struct {
	Email     string             `bson:"email,omitempty" json:"email,omitempty"`
	Name      string             `bson:"name,omitempty" json:"name,omitempty"`
	AvatarURL string             `bson:"avatarURL,omitempty" json:"avatarURL,omitempty"`
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Role      string             `bson:"role,omitempty" json:"role,omitempty"`
	Verified  bool               `bson:"verified,omitempty" json:"verified,omitempty"`
	Banned    bool               `bson:"banned,omitempty" json:"banned,omitempty"`
}
