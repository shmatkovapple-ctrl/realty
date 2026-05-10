package usecase

import (
"context"
"time"

"github.com/google/uuid"
"realty/services/deal-service/internal/domain/entity"
)

type DealUseCase interface {
CreateDeal(ctx context.Context, input CreateDealInput) (*entity.Deal, error)
UpdateDeal(ctx context.Context, input UpdateDealInput) (*entity.Deal, error)
GetDeal(ctx context.Context, id uuid.UUID) (*entity.Deal, error)
ListDeals(ctx context.Context, input ListDealsInput) ([]*entity.Deal, int, error)
}

type ViewingUseCase interface {
CreateViewingRequest(ctx context.Context, input CreateViewingInput) (*entity.ViewingRequest, error)
UpdateViewingRequest(ctx context.Context, input UpdateViewingInput) (*entity.ViewingRequest, error)
ListViewingRequests(ctx context.Context, input ListViewingsInput) ([]*entity.ViewingRequest, int, error)
}

type FavoriteUseCase interface {
AddToFavorites(ctx context.Context, userID, listingID uuid.UUID) (*entity.Favorite, error)
RemoveFromFavorites(ctx context.Context, userID, listingID uuid.UUID) error
ListFavorites(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.Favorite, int, error)
}

type CreateDealInput struct {
ListingID   uuid.UUID
BuyerID     uuid.UUID
AgentID     *uuid.UUID
PriceAgreed float64
Currency    string
}

type UpdateDealInput struct {
ID          uuid.UUID
Status      entity.DealStatus
PriceAgreed float64
}

type ListDealsInput struct {
BuyerID   *uuid.UUID
SellerID  *uuid.UUID
ListingID *uuid.UUID
Page      int
Limit     int
}

type CreateViewingInput struct {
ListingID   uuid.UUID
BuyerID     uuid.UUID
Comment     string
ScheduledAt *time.Time
}

type UpdateViewingInput struct {
ID          uuid.UUID
Status      entity.ViewingStatus
ScheduledAt *time.Time
}

type ListViewingsInput struct {
ListingID *uuid.UUID
BuyerID   *uuid.UUID
Page      int
Limit     int
}
