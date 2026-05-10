package repository

import (
"context"

"github.com/google/uuid"
"realty/services/listing-service/internal/domain/entity"
)

type ListingRepository interface {
Save(ctx context.Context, listing *entity.Listing) error
FindByID(ctx context.Context, id uuid.UUID) (*entity.Listing, error)
Update(ctx context.Context, listing *entity.Listing) error
Delete(ctx context.Context, id uuid.UUID) error
FindBySeller(ctx context.Context, sellerID uuid.UUID, page, limit int) ([]*entity.Listing, int, error)
}

type ListingMediaRepository interface {
Save(ctx context.Context, listingID uuid.UUID, url, mediaType string, sortOrder int) error
FindByListingID(ctx context.Context, listingID uuid.UUID) ([]string, error)
DeleteByListingID(ctx context.Context, listingID uuid.UUID) error
}
