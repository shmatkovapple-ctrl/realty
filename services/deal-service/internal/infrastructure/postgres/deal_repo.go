package postgres

import (
"context"
"fmt"

"github.com/google/uuid"
"github.com/jackc/pgx/v5/pgxpool"
"realty/services/deal-service/internal/domain/entity"
"realty/services/deal-service/internal/domain/repository"
)

type dealRepository struct {
db *pgxpool.Pool
}

func NewDealRepository(db *pgxpool.Pool) repository.DealRepository {
return &dealRepository{db: db}
}

func (r *dealRepository) Save(ctx context.Context, d *entity.Deal) error {
_, err := r.db.Exec(ctx, `
INSERT INTO deals (id, listing_id, buyer_id, agent_id, status, price_agreed, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`, d.ID, d.ListingID, d.BuyerID, d.AgentID, d.Status, d.PriceAgreed, d.CreatedAt)
if err != nil {
return fmt.Errorf("сохранение сделки: %w", err)
}
return nil
}

func (r *dealRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Deal, error) {
d := &entity.Deal{}
err := r.db.QueryRow(ctx, `
SELECT id, listing_id, buyer_id, agent_id, status, price_agreed, created_at, closed_at
FROM deals WHERE id = $1
`, id).Scan(&d.ID, &d.ListingID, &d.BuyerID, &d.AgentID, &d.Status, &d.PriceAgreed, &d.CreatedAt, &d.ClosedAt)
if err != nil {
return nil, fmt.Errorf("поиск сделки: %w", err)
}
return d, nil
}

func (r *dealRepository) Update(ctx context.Context, d *entity.Deal) error {
_, err := r.db.Exec(ctx, `
UPDATE deals SET status=$1, price_agreed=$2, closed_at=$3 WHERE id=$4
`, d.Status, d.PriceAgreed, d.ClosedAt, d.ID)
if err != nil {
return fmt.Errorf("обновление сделки: %w", err)
}
return nil
}

func (r *dealRepository) FindByBuyer(ctx context.Context, buyerID uuid.UUID, page, limit int) ([]*entity.Deal, int, error) {
offset := (page - 1) * limit
rows, err := r.db.Query(ctx, `
SELECT id, listing_id, buyer_id, agent_id, status, price_agreed, created_at, closed_at
FROM deals WHERE buyer_id = $1
ORDER BY created_at DESC LIMIT $2 OFFSET $3
`, buyerID, limit, offset)
if err != nil {
return nil, 0, fmt.Errorf("поиск сделок покупателя: %w", err)
}
defer rows.Close()

var deals []*entity.Deal
for rows.Next() {
d := &entity.Deal{}
if err := rows.Scan(&d.ID, &d.ListingID, &d.BuyerID, &d.AgentID, &d.Status, &d.PriceAgreed, &d.CreatedAt, &d.ClosedAt); err != nil {
return nil, 0, fmt.Errorf("сканирование сделки: %w", err)
}
deals = append(deals, d)
}

var total int
r.db.QueryRow(ctx, `SELECT COUNT(*) FROM deals WHERE buyer_id=$1`, buyerID).Scan(&total)
return deals, total, nil
}

func (r *dealRepository) FindByListing(ctx context.Context, listingID uuid.UUID, page, limit int) ([]*entity.Deal, int, error) {
offset := (page - 1) * limit
rows, err := r.db.Query(ctx, `
SELECT id, listing_id, buyer_id, agent_id, status, price_agreed, created_at, closed_at
FROM deals WHERE listing_id = $1
ORDER BY created_at DESC LIMIT $2 OFFSET $3
`, listingID, limit, offset)
if err != nil {
return nil, 0, fmt.Errorf("поиск сделок по объявлению: %w", err)
}
defer rows.Close()

var deals []*entity.Deal
for rows.Next() {
d := &entity.Deal{}
if err := rows.Scan(&d.ID, &d.ListingID, &d.BuyerID, &d.AgentID, &d.Status, &d.PriceAgreed, &d.CreatedAt, &d.ClosedAt); err != nil {
return nil, 0, fmt.Errorf("сканирование сделки: %w", err)
}
deals = append(deals, d)
}

var total int
r.db.QueryRow(ctx, `SELECT COUNT(*) FROM deals WHERE listing_id=$1`, listingID).Scan(&total)
return deals, total, nil
}
