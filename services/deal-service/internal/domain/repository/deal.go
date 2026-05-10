package repository

import (
"context"

"github.com/google/uuid"
"realty/services/deal-service/internal/domain/entity"
)

type DealRepository interface {
Save(ctx context.Context, deal *entity.Deal) error
FindByID(ctx context.Context, id uuid.UUID) (*entity.Deal, error)
Update(ctx context.Context, deal *entity.Deal) error
FindByBuyer(ctx context.Context, buyerID uuid.UUID, page, limit int) ([]*entity.Deal, int, error)
FindByListing(ctx context.Context, listingID uuid.UUID, page, limit int) ([]*entity.Deal, int, error)
}

type ViewingRepository interface {
Save(ctx context.Context, viewing *entity.ViewingRequest) error
FindByID(ctx context.Context, id uuid.UUID) (*entity.ViewingRequest, error)
Update(ctx context.Context, viewing *entity.ViewingRequest) error
FindByListing(ctx context.Context, listingID uuid.UUID, page, limit int) ([]*entity.ViewingRequest, int, error)
FindByBuyer(ctx context.Context, buyerID uuid.UUID, page, limit int) ([]*entity.ViewingRequest, int, error)
}

type FavoriteRepository interface {
Save(ctx context.Context, favorite *entity.Favorite) error
Delete(ctx context.Context, userID, listingID uuid.UUID) error
FindByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.Favorite, int, error)
Exists(ctx context.Context, userID, listingID uuid.UUID) (bool, error)
}
