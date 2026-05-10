package event

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventUserRegistered = "user.registered"
	EventUserLoggedIn   = "user.logged_in"
	EventUserBlocked    = "user.blocked"
)

type UserRegistered struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	OccuredAt time.Time `json:"occured_at"`
}

type UserLoggedIn struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	OccuredAt time.Time `json:"occured_at"`
}
