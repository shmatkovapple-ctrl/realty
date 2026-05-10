package repository

import (
	"context"

	"github.com/google/uuid"
	"realty/services/user-service/internal/domain/entity"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserProfileRepository interface {
	Save(ctx context.Context, profile *entity.UserProfile) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error)
	Update(ctx context.Context, profile *entity.UserProfile) error
}
