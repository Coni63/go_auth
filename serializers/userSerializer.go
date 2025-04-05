package serializers

import (
	"time"
)

type UserPublicView struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	// No password field
}

type UserPrivateView struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Lastname  string    `json:"lastname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// No password field
}
