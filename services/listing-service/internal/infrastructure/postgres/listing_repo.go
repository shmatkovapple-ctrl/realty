package postgres

import (
"context"
"fmt"

"github.com/google/uuid"
"github.com/jackc/pgx/v5/pgxpool"
"realty/services/listing-service/internal/domain/entity"
"realty/services/listing-service/internal/domain/repository"
)

type listingRepository struct {
db *pgxpool.Pool
}

func NewListingRepository(db *pgxpool.Pool) repository.ListingRepository {
return &listingRepository{db: db}
}

func (r *listingRepository) Save(ctx context.Context, l *entity.Listing) error {
tx, err := r.db.Begin(ctx)
if err != nil {
return fmt.Errorf("начало транзакции: %w", err)
}
defer tx.Rollback(ctx)

_, err = tx.Exec(ctx, `
INSERT INTO listings
(id, seller_id, agent_id, type, status, title, description, price, currency, area_sqm, rooms, floor, floors_total, published_at, created_at, updated_at)
VALUES
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
`, l.ID, l.SellerID, l.AgentID, l.Type, l.Status, l.Title, l.Description,
l.Price, l.Currency, l.AreaSqm, l.Rooms, l.Floor, l.FloorsTotal,
l.PublishedAt, l.CreatedAt, l.UpdatedAt)
if err != nil {
return fmt.Errorf("сохранение объявления: %w", err)
}

_, err = tx.Exec(ctx, `
INSERT INTO addresses (id, listing_id, country, city, district, street, building, lat, lng)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`, uuid.New(), l.ID, l.Address.Country, l.Address.City, l.Address.District,
l.Address.Street, l.Address.Building, l.Address.Lat, l.Address.Lng)
if err != nil {
return fmt.Errorf("сохранение адреса: %w", err)
}

return tx.Commit(ctx)
}

func (r *listingRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Listing, error) {
l := &entity.Listing{}

err := r.db.QueryRow(ctx, `
SELECT l.id, l.seller_id, l.agent_id, l.type, l.status, l.title, l.description,
       l.price, l.currency, l.area_sqm, l.rooms, l.floor, l.floors_total,
       l.published_at, l.created_at, l.updated_at,
       a.country, a.city, a.district, a.street, a.building, a.lat, a.lng
FROM listings l
LEFT JOIN addresses a ON a.listing_id = l.id
WHERE l.id = $1
`, id).Scan(
&l.ID, &l.SellerID, &l.AgentID, &l.Type, &l.Status, &l.Title, &l.Description,
&l.Price, &l.Currency, &l.AreaSqm, &l.Rooms, &l.Floor, &l.FloorsTotal,
&l.PublishedAt, &l.CreatedAt, &l.UpdatedAt,
&l.Address.Country, &l.Address.City, &l.Address.District,
&l.Address.Street, &l.Address.Building, &l.Address.Lat, &l.Address.Lng,
)
if err != nil {
return nil, fmt.Errorf("поиск объявления по id: %w", err)
}

return l, nil
}

func (r *listingRepository) Update(ctx context.Context, l *entity.Listing) error {
tx, err := r.db.Begin(ctx)
if err != nil {
return fmt.Errorf("начало транзакции: %w", err)
}
defer tx.Rollback(ctx)

_, err = tx.Exec(ctx, `
UPDATE listings SET
type=$1, status=$2, title=$3, description=$4, price=$5,
currency=$6, area_sqm=$7, rooms=$8, floor=$9, floors_total=$10,
published_at=$11, updated_at=$12
WHERE id=$13
`, l.Type, l.Status, l.Title, l.Description, l.Price,
l.Currency, l.AreaSqm, l.Rooms, l.Floor, l.FloorsTotal,
l.PublishedAt, l.UpdatedAt, l.ID)
if err != nil {
return fmt.Errorf("обновление объявления: %w", err)
}

_, err = tx.Exec(ctx, `
UPDATE addresses SET
country=$1, city=$2, district=$3, street=$4, building=$5, lat=$6, lng=$7
WHERE listing_id=$8
`, l.Address.Country, l.Address.City, l.Address.District,
l.Address.Street, l.Address.Building, l.Address.Lat, l.Address.Lng, l.ID)
if err != nil {
return fmt.Errorf("обновление адреса: %w", err)
}

return tx.Commit(ctx)
}

func (r *listingRepository) Delete(ctx context.Context, id uuid.UUID) error {
_, err := r.db.Exec(ctx, `DELETE FROM listings WHERE id=$1`, id)
if err != nil {
return fmt.Errorf("удаление объявления: %w", err)
}
return nil
}

func (r *listingRepository) FindBySeller(ctx context.Context, sellerID uuid.UUID, page, limit int) ([]*entity.Listing, int, error) {
offset := (page - 1) * limit

rows, err := r.db.Query(ctx, `
SELECT l.id, l.seller_id, l.agent_id, l.type, l.status, l.title, l.description,
       l.price, l.currency, l.area_sqm, l.rooms, l.floor, l.floors_total,
       l.published_at, l.created_at, l.updated_at,
       a.country, a.city, a.district, a.street, a.building, a.lat, a.lng
FROM listings l
LEFT JOIN addresses a ON a.listing_id = l.id
WHERE l.seller_id = $1
ORDER BY l.created_at DESC
LIMIT $2 OFFSET $3
`, sellerID, limit, offset)
if err != nil {
return nil, 0, fmt.Errorf("поиск объявлений продавца: %w", err)
}
defer rows.Close()

var listings []*entity.Listing
for rows.Next() {
l := &entity.Listing{}
err := rows.Scan(
&l.ID, &l.SellerID, &l.AgentID, &l.Type, &l.Status, &l.Title, &l.Description,
&l.Price, &l.Currency, &l.AreaSqm, &l.Rooms, &l.Floor, &l.FloorsTotal,
&l.PublishedAt, &l.CreatedAt, &l.UpdatedAt,
&l.Address.Country, &l.Address.City, &l.Address.District,
&l.Address.Street, &l.Address.Building, &l.Address.Lat, &l.Address.Lng,
)
if err != nil {
return nil, 0, fmt.Errorf("сканирование объявления: %w", err)
}
listings = append(listings, l)
}

var total int
err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM listings WHERE seller_id=$1`, sellerID).Scan(&total)
if err != nil {
return nil, 0, fmt.Errorf("подсчёт объявлений: %w", err)
}

return listings, total, nil
}
