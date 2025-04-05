package models

import (
	"go_auth/serializers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name     string    `json:"name"`
	Lastname string    `json:"lastname"`
	Password string
}

// Serializers for User model
// These methods convert the User model to its public and private views.
func (u *User) ToPublicView() serializers.UserPublicView {
	return serializers.UserPublicView{
		ID:       u.ID.String(),
		Name:     u.Name,
		Lastname: u.Lastname,
	}
}

func (u *User) ToPrivateView() serializers.UserPrivateView {
	return serializers.UserPrivateView{
		ID:        u.ID.String(),
		Name:      u.Name,
		Lastname:  u.Lastname,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
