package kafka

import (
"context"
"fmt"
"log"

"realty/services/search-service/internal/domain/entity"
"realty/services/search-service/internal/usecase"
)

type EventHandler struct {
searchUC usecase.SearchUseCase
}

func NewEventHandler(searchUC usecase.SearchUseCase) *EventHandler {
return &EventHandler{searchUC: searchUC}
}

func (h *EventHandler) HandleListingPublished(ctx context.Context, event map[string]any) error {
listing, err := eventToListingIndex(event)
if err != nil {
return fmt.Errorf("конвертация события: %w", err)
}
listing.Status = "published"

if err := h.searchUC.IndexListing(ctx, listing); err != nil {
return fmt.Errorf("индексация объявления: %w", err)
}

log.Printf("объявление проиндексировано: %s", listing.ListingID)
return nil
}

func (h *EventHandler) HandleListingDeleted(ctx context.Context, event map[string]any) error {
listingID, ok := event["listing_id"].(string)
if !ok || listingID == "" {
return fmt.Errorf("listing_id не найден в событии")
}

if err := h.searchUC.DeleteIndex(ctx, listingID); err != nil {
return fmt.Errorf("удаление из индекса: %w", err)
}

log.Printf("объявление удалено из индекса: %s", listingID)
return nil
}

func eventToListingIndex(event map[string]any) (*entity.ListingIndex, error) {
listing := &entity.ListingIndex{}

if v, ok := event["listing_id"].(string); ok { listing.ListingID = v }
if v, ok := event["type"].(string); ok        { listing.Type = v }
if v, ok := event["title"].(string); ok       { listing.Title = v }
if v, ok := event["description"].(string); ok { listing.Description = v }
if v, ok := event["price"].(float64); ok      { listing.Price = v }
if v, ok := event["area_sqm"].(float64); ok   { listing.AreaSqm = v }
if v, ok := event["rooms"].(float64); ok      { listing.Rooms = int32(v) }
if v, ok := event["city"].(string); ok        { listing.City = v }
if v, ok := event["district"].(string); ok    { listing.District = v }
if v, ok := event["lat"].(float64); ok        { listing.Lat = v }
if v, ok := event["lng"].(float64); ok        { listing.Lng = v }
if v, ok := event["preview_url"].(string); ok { listing.PreviewURL = v }

if listing.ListingID == "" {
return nil, fmt.Errorf("listing_id обязателен")
}

return listing, nil
}
