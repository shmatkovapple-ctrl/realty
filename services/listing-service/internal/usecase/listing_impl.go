package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"realty/services/listing-service/internal/domain/entity"
	"realty/services/listing-service/internal/domain/event"
	"realty/services/listing-service/internal/domain/repository"
	kafkaInfra "realty/services/listing-service/internal/infrastructure/kafka"
	minioInfra "realty/services/listing-service/internal/infrastructure/minio"

	"github.com/google/uuid"
)

type listingUseCase struct {
	listingRepo repository.ListingRepository
	publisher   *kafkaInfra.Publisher
	storage     *minioInfra.Storage
}

func NewListingUseCase(
	listingRepo repository.ListingRepository,
	publisher *kafkaInfra.Publisher,
	storage *minioInfra.Storage,
) ListingUseCase {
	return &listingUseCase{
		listingRepo: listingRepo,
		publisher:   publisher,
		storage:     storage,
	}
}

func (uc *listingUseCase) Create(ctx context.Context, input CreateInput) (*entity.Listing, error) {
	listing, err := entity.NewListing(
		input.SellerID,
		entity.ListingType(input.Type),
		input.Title,
		input.Price,
		input.AreaSqm,
	)
	if err != nil {
		return nil, fmt.Errorf("создание объявления: %w", err)
	}

	listing.Description = input.Description
	listing.Rooms = input.Rooms
	listing.Floor = input.Floor
	listing.FloorsTotal = input.FloorsTotal
	listing.Address = input.Address

	if err := uc.listingRepo.Save(ctx, listing); err != nil {
		return nil, fmt.Errorf("сохранение объявления: %w", err)
	}

	_ = uc.publisher.Publish(ctx, event.EventListingCreated, event.ListingCreated{
		ListingID: listing.ID,
		SellerID:  listing.SellerID,
		Type:      string(listing.Type),
		Title:     listing.Title,
		Price:     listing.Price,
		City:      listing.Address.City,
		OccuredAt: time.Now(),
	})

	return listing, nil
}

func (uc *listingUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Listing, error) {
	listing, err := uc.listingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("объявление не найдено: %w", err)
	}
	return listing, nil
}

func (uc *listingUseCase) Update(ctx context.Context, input UpdateInput) (*entity.Listing, error) {
	listing, err := uc.listingRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("объявление не найдено: %w", err)
	}

	if listing.SellerID != input.SellerID {
		return nil, errors.New("нет прав для редактирования этого объявления")
	}

	if input.Title != "" {
		listing.Title = input.Title
	}
	if input.Description != "" {
		listing.Description = input.Description
	}
	if input.Price > 0 {
		listing.Price = input.Price
	}
	if input.AreaSqm > 0 {
		listing.AreaSqm = input.AreaSqm
	}
	if input.Rooms > 0 {
		listing.Rooms = input.Rooms
	}
	if input.Floor > 0 {
		listing.Floor = input.Floor
	}
	if input.FloorsTotal > 0 {
		listing.FloorsTotal = input.FloorsTotal
	}
	if input.Address.City != "" {
		listing.Address = input.Address
	}
	if len(input.MediaURLs) > 0 {
		listing.MediaURLs = input.MediaURLs
	}

	listing.UpdatedAt = time.Now()

	if err := uc.listingRepo.Update(ctx, listing); err != nil {
		return nil, fmt.Errorf("обновление объявления: %w", err)
	}

	return listing, nil
}

func (uc *listingUseCase) Delete(ctx context.Context, id, sellerID uuid.UUID) error {
	listing, err := uc.listingRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("объявление не найдено: %w", err)
	}

	if listing.SellerID != sellerID {
		return errors.New("нет прав для удаления этого объявления")
	}

	if err := uc.listingRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("удаление объявления: %w", err)
	}

	_ = uc.publisher.Publish(ctx, event.EventListingDeleted, event.ListingDeleted{
		ListingID: id,
		OccuredAt: time.Now(),
	})

	return nil
}

func (uc *listingUseCase) Publish(ctx context.Context, id, sellerID uuid.UUID) (*entity.Listing, error) {
	listing, err := uc.listingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("объявление не найдено: %w", err)
	}

	if listing.SellerID != sellerID {
		return nil, errors.New("нет прав для публикации этого объявления")
	}

	if err := listing.Publish(); err != nil {
		return nil, fmt.Errorf("публикация объявления: %w", err)
	}

	if err := uc.listingRepo.Update(ctx, listing); err != nil {
		return nil, fmt.Errorf("сохранение после публикации: %w", err)
	}

	previewURL := ""
	if len(listing.MediaURLs) > 0 {
		previewURL = listing.MediaURLs[0]
	}

	_ = uc.publisher.Publish(ctx, event.EventListingPublished, event.ListingPublished{
		ListingID:  listing.ID,
		SellerID:   listing.SellerID,
		Type:       string(listing.Type),
		Title:      listing.Title,
		Price:      listing.Price,
		City:       listing.Address.City,
		District:   listing.Address.District,
		Lat:        listing.Address.Lat,
		Lng:        listing.Address.Lng,
		PreviewURL: previewURL,
		OccuredAt:  time.Now(),
	})

	return listing, nil
}

func (uc *listingUseCase) GetUploadURL(ctx context.Context, listingID, filename, contentType string) (string, string, error) {
	uploadURL, fileURL, err := uc.storage.GetUploadURL(ctx, listingID, filename, contentType)
	if err != nil {
		return "", "", fmt.Errorf("получение URL для загрузки: %w", err)
	}
	return uploadURL, fileURL, nil
}

func (uc *listingUseCase) SaveMedia(ctx context.Context, listingID uuid.UUID, urls []string) error {
	if err := uc.listingRepo.SaveMedia(ctx, listingID, urls); err != nil {
		return fmt.Errorf(`сохранение медиа: %w`, err)
	}
	return nil
}

func (uc *listingUseCase) ListBySeller(ctx context.Context, sellerID uuid.UUID, page, limit int) ([]*entity.Listing, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	listings, total, err := uc.listingRepo.FindBySeller(ctx, sellerID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("получение объявлений продавца: %w", err)
	}

	return listings, total, nil
}
