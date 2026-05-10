package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"realty/services/user-service/internal/domain/entity"
	"realty/services/user-service/internal/domain/repository"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, u *entity.User) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("начало транзакции: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
INSERT INTO users (id, email, phone, password_hash, role, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`, u.ID, u.Email, u.Phone, u.PasswordHash, u.Role, u.Status, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("сохранение пользователя: %w", err)
	}

	if u.Profile != nil {
		u.Profile.UserID = u.ID
		_, err = tx.Exec(ctx, `
INSERT INTO user_profiles (id, user_id, first_name, last_name, avatar_url, bio)
VALUES ($1, $2, $3, $4, $5, $6)
`, u.Profile.ID, u.Profile.UserID, u.Profile.FirstName, u.Profile.LastName, u.Profile.AvatarURL, u.Profile.Bio)
		if err != nil {
			return fmt.Errorf("сохранение профиля: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	u := &entity.User{Profile: &entity.UserProfile{}}

	err := r.db.QueryRow(ctx, `
SELECT u.id, u.email, u.phone, u.password_hash, u.role, u.status, u.created_at, u.updated_at,
       p.id, p.first_name, p.last_name, p.avatar_url, p.bio, p.verified_at
FROM users u
LEFT JOIN user_profiles p ON p.user_id = u.id
WHERE u.id = $1
`, id).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt,
		&u.Profile.ID, &u.Profile.FirstName, &u.Profile.LastName, &u.Profile.AvatarURL, &u.Profile.Bio, &u.Profile.VerifiedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("поиск пользователя по id: %w", err)
	}

	return u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	u := &entity.User{Profile: &entity.UserProfile{}}

	err := r.db.QueryRow(ctx, `
SELECT u.id, u.email, u.phone, u.password_hash, u.role, u.status, u.created_at, u.updated_at,
       p.id, p.first_name, p.last_name, p.avatar_url, p.bio, p.verified_at
FROM users u
LEFT JOIN user_profiles p ON p.user_id = u.id
WHERE u.email = $1
`, email).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.Role, &u.Status, &u.CreatedAt, &u.UpdatedAt,
		&u.Profile.ID, &u.Profile.FirstName, &u.Profile.LastName, &u.Profile.AvatarURL, &u.Profile.Bio, &u.Profile.VerifiedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("поиск пользователя по email: %w", err)
	}

	return u, nil
}

func (r *userRepository) Update(ctx context.Context, u *entity.User) error {
	_, err := r.db.Exec(ctx, `
UPDATE users SET email=$1, phone=$2, role=$3, status=$4, updated_at=$5
WHERE id=$6
`, u.Email, u.Phone, u.Role, u.Status, u.UpdatedAt, u.ID)
	if err != nil {
		return fmt.Errorf("обновление пользователя: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("удаление пользователя: %w", err)
	}
	return nil
}
