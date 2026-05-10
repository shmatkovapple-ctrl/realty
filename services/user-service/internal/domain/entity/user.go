package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Role string
type Status string

const (
	RoleBuyer  Role = "buyer"
	RoleSeller Role = "seller"
	RoleAgent  Role = "agent"
	RoleAdmin  Role = "admin"

	StatusActive  Status = "active"
	StatusBlocked Status = "blocked"
	StatusPending Status = "pending"
)

type User struct {
	ID           uuid.UUID
	Email        string
	Phone        string
	PasswordHash string
	Role         Role
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Profile      *UserProfile
}

type UserProfile struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	FirstName  string
	LastName   string
	AvatarURL  string
	Bio        string
	VerifiedAt *time.Time
}

func NewUser(email, passwordHash, phone string, role Role) (*User, error) {
	if email == "" {
		return nil, errors.New("email не может быть пустым")
	}
	if passwordHash == "" {
		return nil, errors.New("пароль не может быть пустым")
	}
	if !role.IsValid() {
		return nil, errors.New("недопустимая роль пользователя")
	}

	now := time.Now()
	return &User{
		ID:           uuid.New(),
		Email:        email,
		Phone:        phone,
		PasswordHash: passwordHash,
		Role:         role,
		Status:       StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
		Profile: &UserProfile{
			ID:     uuid.New(),
			UserID: uuid.New(),
		},
	}, nil
}

func (u *User) Block() error {
	if u.Status == StatusBlocked {
		return errors.New("пользователь уже заблокирован")
	}
	u.Status = StatusBlocked
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Activate() error {
	if u.Status == StatusActive {
		return errors.New("пользователь уже активен")
	}
	u.Status = StatusActive
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

func (r Role) IsValid() bool {
	switch r {
	case RoleBuyer, RoleSeller, RoleAgent, RoleAdmin:
		return true
	}
	return false
}
