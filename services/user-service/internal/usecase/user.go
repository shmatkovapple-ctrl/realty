package usecase

import (
	"context"

	"github.com/google/uuid"
	"realty/services/user-service/internal/domain/entity"
)

type UserUseCase interface {
	Register(ctx context.Context, input RegisterInput) (*AuthOutput, error)
	Login(ctx context.Context, input LoginInput) (*AuthOutput, error)
	Logout(ctx context.Context, userID, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*AuthOutput, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	UpdateProfile(ctx context.Context, input UpdateProfileInput) (*entity.User, error)
	ValidateToken(ctx context.Context, accessToken string) (*TokenClaims, error)
}

type RegisterInput struct {
	Email    string
	Password string
	Phone    string
	Role     string
}

type LoginInput struct {
	Email      string
	Password   string
	DeviceInfo string
}

type UpdateProfileInput struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
	Bio       string
}

type AuthOutput struct {
	AccessToken  string
	RefreshToken string
	User         *entity.User
}

type TokenClaims struct {
	UserID string
	Role   string
}
