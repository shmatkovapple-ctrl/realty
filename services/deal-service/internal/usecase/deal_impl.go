package usecase

import (
"context"
"fmt"
"time"

"github.com/google/uuid"
"realty/services/deal-service/internal/domain/entity"
"realty/services/deal-service/internal/domain/event"
"realty/services/deal-service/internal/domain/repository"
kafkaInfra "realty/services/deal-service/internal/infrastructure/kafka"
)

type dealUseCase struct {
dealRepo   repository.DealRepository
publisher  *kafkaInfra.Publisher
}

type viewingUseCase struct {
viewingRepo repository.ViewingRepository
publisher   *kafkaInfra.Publisher
}

type favoriteUseCase struct {
favoriteRepo repository.FavoriteRepository
}

func NewDealUseCase(dealRepo repository.DealRepository, publisher *kafkaInfra.Publisher) DealUseCase {
return &dealUseCase{dealRepo: dealRepo, publisher: publisher}
}

func NewViewingUseCase(viewingRepo repository.ViewingRepository, publisher *kafkaInfra.Publisher) ViewingUseCase {
return &viewingUseCase{viewingRepo: viewingRepo, publisher: publisher}
}

func NewFavoriteUseCase(favoriteRepo repository.FavoriteRepository) FavoriteUseCase {
return &favoriteUseCase{favoriteRepo: favoriteRepo}
}

func (uc *dealUseCase) CreateDeal(ctx context.Context, input CreateDealInput) (*entity.Deal, error) {
deal, err := entity.NewDeal(input.ListingID, input.BuyerID, input.AgentID, input.PriceAgreed)
if err != nil {
return nil, fmt.Errorf("создание сделки: %w", err)
}

if input.Currency != "" {
deal.Currency = input.Currency
}

if err := uc.dealRepo.Save(ctx, deal); err != nil {
return nil, fmt.Errorf("сохранение сделки: %w", err)
}

_ = uc.publisher.Publish(ctx, event.EventDealCreated, event.DealCreated{
DealID:    deal.ID,
ListingID: deal.ListingID,
BuyerID:   deal.BuyerID,
Price:     deal.PriceAgreed,
OccuredAt: time.Now(),
})

return deal, nil
}

func (uc *dealUseCase) UpdateDeal(ctx context.Context, input UpdateDealInput) (*entity.Deal, error) {
deal, err := uc.dealRepo.FindByID(ctx, input.ID)
if err != nil {
return nil, fmt.Errorf("сделка не найдена: %w", err)
}

switch input.Status {
case entity.DealStatusApproved:
if err := deal.Approve(); err != nil {
return nil, err
}
case entity.DealStatusClosed:
if err := deal.Close(); err != nil {
return nil, err
}
case entity.DealStatusRejected:
if err := deal.Reject(); err != nil {
return nil, err
}
case entity.DealStatusReview:
deal.Status = entity.DealStatusReview
}

if input.PriceAgreed > 0 {
deal.PriceAgreed = input.PriceAgreed
}

if err := uc.dealRepo.Update(ctx, deal); err != nil {
return nil, fmt.Errorf("обновление сделки: %w", err)
}

_ = uc.publisher.Publish(ctx, event.EventDealStatusChanged, event.DealStatusChanged{
DealID:    deal.ID,
BuyerID:   deal.BuyerID,
Status:    string(deal.Status),
OccuredAt: time.Now(),
})

return deal, nil
}

func (uc *dealUseCase) GetDeal(ctx context.Context, id uuid.UUID) (*entity.Deal, error) {
deal, err := uc.dealRepo.FindByID(ctx, id)
if err != nil {
return nil, fmt.Errorf("сделка не найдена: %w", err)
}
return deal, nil
}

func (uc *dealUseCase) ListDeals(ctx context.Context, input ListDealsInput) ([]*entity.Deal, int, error) {
if input.Page < 1 {
input.Page = 1
}
if input.Limit < 1 || input.Limit > 100 {
input.Limit = 20
}

if input.BuyerID != nil {
return uc.dealRepo.FindByBuyer(ctx, *input.BuyerID, input.Page, input.Limit)
}
if input.ListingID != nil {
return uc.dealRepo.FindByListing(ctx, *input.ListingID, input.Page, input.Limit)
}

return nil, 0, fmt.Errorf("необходимо указать buyer_id или listing_id")
}

func (uc *viewingUseCase) CreateViewingRequest(ctx context.Context, input CreateViewingInput) (*entity.ViewingRequest, error) {
viewing, err := entity.NewViewingRequest(input.ListingID, input.BuyerID, input.Comment, input.ScheduledAt)
if err != nil {
return nil, fmt.Errorf("создание заявки на просмотр: %w", err)
}

if err := uc.viewingRepo.Save(ctx, viewing); err != nil {
return nil, fmt.Errorf("сохранение заявки: %w", err)
}

_ = uc.publisher.Publish(ctx, event.EventViewingCreated, event.ViewingCreated{
ViewingID:   viewing.ID,
ListingID:   viewing.ListingID,
BuyerID:     viewing.BuyerID,
ScheduledAt: viewing.ScheduledAt,
OccuredAt:   time.Now(),
})

return viewing, nil
}

func (uc *viewingUseCase) UpdateViewingRequest(ctx context.Context, input UpdateViewingInput) (*entity.ViewingRequest, error) {
viewing, err := uc.viewingRepo.FindByID(ctx, input.ID)
if err != nil {
return nil, fmt.Errorf("заявка не найдена: %w", err)
}

switch input.Status {
case entity.ViewingStatusConfirmed:
if err := viewing.Confirm(); err != nil {
return nil, err
}
case entity.ViewingStatusCancelled:
if err := viewing.Cancel(); err != nil {
return nil, err
}
case entity.ViewingStatusCompleted:
if err := viewing.Complete(); err != nil {
return nil, err
}
}

if input.ScheduledAt != nil {
viewing.ScheduledAt = input.ScheduledAt
}

if err := uc.viewingRepo.Update(ctx, viewing); err != nil {
return nil, fmt.Errorf("обновление заявки: %w", err)
}

_ = uc.publisher.Publish(ctx, event.EventViewingStatusChanged, event.ViewingStatusChanged{
ViewingID: viewing.ID,
BuyerID:   viewing.BuyerID,
Status:    string(viewing.Status),
OccuredAt: time.Now(),
})

return viewing, nil
}

func (uc *viewingUseCase) ListViewingRequests(ctx context.Context, input ListViewingsInput) ([]*entity.ViewingRequest, int, error) {
if input.Page < 1 {
input.Page = 1
}
if input.Limit < 1 || input.Limit > 100 {
input.Limit = 20
}

if input.ListingID != nil {
return uc.viewingRepo.FindByListing(ctx, *input.ListingID, input.Page, input.Limit)
}
if input.BuyerID != nil {
return uc.viewingRepo.FindByBuyer(ctx, *input.BuyerID, input.Page, input.Limit)
}

return nil, 0, fmt.Errorf("необходимо указать listing_id или buyer_id")
}

func (uc *favoriteUseCase) AddToFavorites(ctx context.Context, userID, listingID uuid.UUID) (*entity.Favorite, error) {
exists, err := uc.favoriteRepo.Exists(ctx, userID, listingID)
if err != nil {
return nil, fmt.Errorf("проверка избранного: %w", err)
}
if exists {
return nil, fmt.Errorf("объявление уже в избранном")
}

favorite := &entity.Favorite{
ID:        uuid.New(),
UserID:    userID,
ListingID: listingID,
CreatedAt: time.Now(),
}

if err := uc.favoriteRepo.Save(ctx, favorite); err != nil {
return nil, fmt.Errorf("добавление в избранное: %w", err)
}

return favorite, nil
}

func (uc *favoriteUseCase) RemoveFromFavorites(ctx context.Context, userID, listingID uuid.UUID) error {
if err := uc.favoriteRepo.Delete(ctx, userID, listingID); err != nil {
return fmt.Errorf("удаление из избранного: %w", err)
}
return nil
}

func (uc *favoriteUseCase) ListFavorites(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.Favorite, int, error) {
if page < 1 {
page = 1
}
if limit < 1 || limit > 100 {
limit = 20
}
return uc.favoriteRepo.FindByUser(ctx, userID, page, limit)
}
