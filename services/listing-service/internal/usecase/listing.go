package usecase

import (
"context"

"github.com/google/uuid"
"realty/services/listing-service/internal/domain/entity"
)

type ListingUseCase interface {
Create(ctx context.Context, input CreateInput) (*entity.Listing, error)
GetByID(ctx context.Context, id uuid.UUID) (*entity.Listing, error)
Update(ctx context.Context, input UpdateInput) (*entity.Listing, error)
Delete(ctx context.Context, id, sellerID uuid.UUID) error
Publish(ctx context.Context, id, sellerID uuid.UUID) (*entity.Listing, error)
GetUploadURL(ctx context.Context, listingID, filename, contentType string) (uploadURL, fileURL string, err error)
ListBySeller(ctx context.Context, sellerID uuid.UUID, page, limit int) ([]*entity.Listing, int, error)
}

type CreateInput struct {
SellerID    uuid.UUID
Type        string
Title       string
Description string
Price       float64
AreaSqm     float64
Rooms       int32
Floor       int32
FloorsTotal int32
Address     entity.Address
}

type UpdateInput struct {
ID          uuid.UUID
SellerID    uuid.UUID
Title       string
Description string
Price       float64
AreaSqm     float64
Rooms       int32
Floor       int32
FloorsTotal int32
Address     entity.Address
}
