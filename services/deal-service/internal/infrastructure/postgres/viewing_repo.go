package postgres

import (
"context"
"fmt"

"github.com/google/uuid"
"github.com/jackc/pgx/v5/pgxpool"
"realty/services/deal-service/internal/domain/entity"
"realty/services/deal-service/internal/domain/repository"
)

type viewingRepository struct {
db *pgxpool.Pool
}

func NewViewingRepository(db *pgxpool.Pool) repository.ViewingRepository {
return &viewingRepository{db: db}
}

func (r *viewingRepository) Save(ctx context.Context, v *entity.ViewingRequest) error {
_, err := r.db.Exec(ctx, `
INSERT INTO viewing_requests (id, listing_id, buyer_id, status, comment, scheduled_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`, v.ID, v.ListingID, v.BuyerID, v.Status, v.Comment, v.ScheduledAt, v.CreatedAt)
if err != nil {
return fmt.Errorf("сохранение заявки на просмотр: %w", err)
}
return nil
}

func (r *viewingRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.ViewingRequest, error) {
v := &entity.ViewingRequest{}
err := r.db.QueryRow(ctx, `
SELECT id, listing_id, buyer_id, status, comment, scheduled_at, created_at
FROM viewing_requests WHERE id = $1
`, id).Scan(&v.ID, &v.ListingID, &v.BuyerID, &v.Status, &v.Comment, &v.ScheduledAt, &v.CreatedAt)
if err != nil {
return nil, fmt.Errorf("поиск заявки на просмотр: %w", err)
}
return v, nil
}

func (r *viewingRepository) Update(ctx context.Context, v *entity.ViewingRequest) error {
_, err := r.db.Exec(ctx, `
UPDATE viewing_requests SET status=$1, scheduled_at=$2 WHERE id=$3
`, v.Status, v.ScheduledAt, v.ID)
if err != nil {
return fmt.Errorf("обновление заявки на просмотр: %w", err)
}
return nil
}

func (r *viewingRepository) FindByListing(ctx context.Context, listingID uuid.UUID, page, limit int) ([]*entity.ViewingRequest, int, error) {
offset := (page - 1) * limit
rows, err := r.db.Query(ctx, `
SELECT id, listing_id, buyer_id, status, comment, scheduled_at, created_at
FROM viewing_requests WHERE listing_id = $1
ORDER BY created_at DESC LIMIT $2 OFFSET $3
`, listingID, limit, offset)
if err != nil {
return nil, 0, fmt.Errorf("поиск заявок по объявлению: %w", err)
}
defer rows.Close()

var viewings []*entity.ViewingRequest
for rows.Next() {
v := &entity.ViewingRequest{}
if err := rows.Scan(&v.ID, &v.ListingID, &v.BuyerID, &v.Status, &v.Comment, &v.ScheduledAt, &v.CreatedAt); err != nil {
return nil, 0, fmt.Errorf("сканирование заявки: %w", err)
}
viewings = append(viewings, v)
}

var total int
r.db.QueryRow(ctx, `SELECT COUNT(*) FROM viewing_requests WHERE listing_id=$1`, listingID).Scan(&total)
return viewings, total, nil
}

func (r *viewingRepository) FindByBuyer(ctx context.Context, buyerID uuid.UUID, page, limit int) ([]*entity.ViewingRequest, int, error) {
offset := (page - 1) * limit
rows, err := r.db.Query(ctx, `
SELECT id, listing_id, buyer_id, status, comment, scheduled_at, created_at
FROM viewing_requests WHERE buyer_id = $1
ORDER BY created_at DESC LIMIT $2 OFFSET $3
`, buyerID, limit, offset)
if err != nil {
return nil, 0, fmt.Errorf("поиск заявок покупателя: %w", err)
}
defer rows.Close()

var viewings []*entity.ViewingRequest
for rows.Next() {
v := &entity.ViewingRequest{}
if err := rows.Scan(&v.ID, &v.ListingID, &v.BuyerID, &v.Status, &v.Comment, &v.ScheduledAt, &v.CreatedAt); err != nil {
return nil, 0, fmt.Errorf("сканирование заявки: %w", err)
}
viewings = append(viewings, v)
}

var total int
r.db.QueryRow(ctx, `SELECT COUNT(*) FROM viewing_requests WHERE buyer_id=$1`, buyerID).Scan(&total)
return viewings, total, nil
}
