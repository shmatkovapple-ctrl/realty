package postgres

import (
"context"
"fmt"
"time"

"github.com/google/uuid"
"github.com/jackc/pgx/v5/pgxpool"
"realty/services/deal-service/internal/domain/entity"
"realty/services/deal-service/internal/domain/repository"
)

type favoriteRepository struct {
db *pgxpool.Pool
}

func NewFavoriteRepository(db *pgxpool.Pool) repository.FavoriteRepository {
return &favoriteRepository{db: db}
}

func (r *favoriteRepository) Save(ctx context.Context, f *entity.Favorite) error {
_, err := r.db.Exec(ctx, `
INSERT INTO favorites (id, user_id, listing_id, created_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, listing_id) DO NOTHING
`, f.ID, f.UserID, f.ListingID, f.CreatedAt)
if err != nil {
return fmt.Errorf("сохранение избранного: %w", err)
}
return nil
}

func (r *favoriteRepository) Delete(ctx context.Context, userID, listingID uuid.UUID) error {
_, err := r.db.Exec(ctx, `
DELETE FROM favorites WHERE user_id=$1 AND listing_id=$2
`, userID, listingID)
if err != nil {
return fmt.Errorf("удаление из избранного: %w", err)
}
return nil
}

func (r *favoriteRepository) FindByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.Favorite, int, error) {
offset := (page - 1) * limit
rows, err := r.db.Query(ctx, `
SELECT id, user_id, listing_id, created_at
FROM favorites WHERE user_id = $1
ORDER BY created_at DESC LIMIT $2 OFFSET $3
`, userID, limit, offset)
if err != nil {
return nil, 0, fmt.Errorf("поиск избранного: %w", err)
}
defer rows.Close()

var favorites []*entity.Favorite
for rows.Next() {
f := &entity.Favorite{}
if err := rows.Scan(&f.ID, &f.UserID, &f.ListingID, &f.CreatedAt); err != nil {
return nil, 0, fmt.Errorf("сканирование избранного: %w", err)
}
favorites = append(favorites, f)
}

var total int
r.db.QueryRow(ctx, `SELECT COUNT(*) FROM favorites WHERE user_id=$1`, userID).Scan(&total)
return favorites, total, nil
}

func (r *favoriteRepository) Exists(ctx context.Context, userID, listingID uuid.UUID) (bool, error) {
var count int
err := r.db.QueryRow(ctx, `
SELECT COUNT(*) FROM favorites WHERE user_id=$1 AND listing_id=$2
`, userID, listingID).Scan(&count)
if err != nil {
return false, fmt.Errorf("проверка избранного: %w", err)
}
return count > 0, nil
}

func NewFavoriteEntity(userID, listingID uuid.UUID) *entity.Favorite {
return &entity.Favorite{
ID:        uuid.New(),
UserID:    userID,
ListingID: listingID,
CreatedAt: time.Now(),
}
}
