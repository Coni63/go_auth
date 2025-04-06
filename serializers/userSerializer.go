package serializers

import (
	"time"
)

type UserPublicView struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
}

type UserPrivateView struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// No password field
}
