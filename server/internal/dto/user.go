package dto

import "time"

type User struct {
	Email     *string    `bson:"email,omitempty" json:"email,omitempty"`
	Name      *string    `bson:"name,omitempty" json:"name,omitempty"`
	AvatarURL *string    `bson:"avatarURL,omitempty" json:"avatarURL,omitempty"`
	Role      *string    `bson:"role,omitempty" json:"role,omitempty"`
	Verified  *bool      `bson:"verified,omitempty" json:"verified,omitempty"`
	UserInfo  *any       `bson:"userInfo,omitempty" json:"userInfo,omitempty"`
	Banned    *bool      `bson:"banned,omitempty" json:"banned,omitempty"`
	UpdatedAt *time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
