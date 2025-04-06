package models

import (
	"go_auth/serializers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string
}

// Serializers for User model
// These methods convert the User model to its public and private views.
func (u *User) ToPublicView() serializers.UserPublicView {
	return serializers.UserPublicView{
		ID:       u.ID.String(),
		UserName: u.UserName,
	}
}

func (u *User) ToPrivateView() serializers.UserPrivateView {
	return serializers.UserPrivateView{
		ID:        u.ID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		UserName:  u.UserName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
