package kafka

import (
"context"
"fmt"
"log"

"github.com/google/uuid"
"realty/services/notification-service/internal/domain/entity"
"realty/services/notification-service/internal/usecase"
)

type EventHandler struct {
uc usecase.NotificationUseCase
}

func NewEventHandler(uc usecase.NotificationUseCase) *EventHandler {
return &EventHandler{uc: uc}
}

func (h *EventHandler) HandleDealCreated(ctx context.Context, event map[string]any) error {
buyerID, err := parseUUID(event, "buyer_id")
if err != nil {
return err
}

payload := map[string]string{
"deal_id": getString(event, "deal_id"),
"price":   fmt.Sprintf("%.0f", getFloat(event, "price")),
}

return h.uc.Send(ctx, usecase.SendInput{
UserID:  buyerID,
Type:    entity.TypeDealCreated,
Channel: entity.ChannelInApp,
Payload: payload,
})
}

func (h *EventHandler) HandleDealStatusChanged(ctx context.Context, event map[string]any) error {
buyerID, err := parseUUID(event, "buyer_id")
if err != nil {
return err
}

payload := map[string]string{
"deal_id": getString(event, "deal_id"),
"status":  getString(event, "status"),
}

return h.uc.Send(ctx, usecase.SendInput{
UserID:  buyerID,
Type:    entity.TypeDealStatus,
Channel: entity.ChannelInApp,
Payload: payload,
})
}

func (h *EventHandler) HandleViewingCreated(ctx context.Context, event map[string]any) error {
buyerID, err := parseUUID(event, "buyer_id")
if err != nil {
return err
}

payload := map[string]string{
"listing_id": getString(event, "listing_id"),
"viewing_id": getString(event, "viewing_id"),
}

return h.uc.Send(ctx, usecase.SendInput{
UserID:  buyerID,
Type:    entity.TypeViewingCreated,
Channel: entity.ChannelInApp,
Payload: payload,
})
}

func (h *EventHandler) HandleViewingStatusChanged(ctx context.Context, event map[string]any) error {
buyerID, err := parseUUID(event, "buyer_id")
if err != nil {
return err
}

statusStr := getString(event, "status")
notifType := entity.TypeViewingConfirmed
if statusStr == "cancelled" {
notifType = entity.TypeViewingCancelled
}

payload := map[string]string{
"viewing_id": getString(event, "viewing_id"),
"status":     statusStr,
}

return h.uc.Send(ctx, usecase.SendInput{
UserID:  buyerID,
Type:    notifType,
Channel: entity.ChannelInApp,
Payload: payload,
})
}

func (h *EventHandler) HandleListingPublished(ctx context.Context, event map[string]any) error {
sellerID, err := parseUUID(event, "seller_id")
if err != nil {
log.Printf("seller_id не найден в событии listing.published: %v", err)
return nil
}

payload := map[string]string{
"listing_id": getString(event, "listing_id"),
}

return h.uc.Send(ctx, usecase.SendInput{
UserID:  sellerID,
Type:    entity.TypeListingPublished,
Channel: entity.ChannelInApp,
Payload: payload,
})
}

func parseUUID(event map[string]any, key string) (uuid.UUID, error) {
val, ok := event[key].(string)
if !ok || val == "" {
return uuid.Nil, fmt.Errorf("%s не найден в событии", key)
}
id, err := uuid.Parse(val)
if err != nil {
return uuid.Nil, fmt.Errorf("неверный формат %s: %w", key, err)
}
return id, nil
}

func getString(event map[string]any, key string) string {
val, _ := event[key].(string)
return val
}

func getFloat(event map[string]any, key string) float64 {
val, _ := event[key].(float64)
return val
}
