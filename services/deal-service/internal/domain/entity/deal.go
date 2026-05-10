package entity

import (
"errors"
"time"

"github.com/google/uuid"
)

type DealStatus string
type ViewingStatus string

const (
DealStatusNew      DealStatus = "new"
DealStatusReview   DealStatus = "review"
DealStatusApproved DealStatus = "approved"
DealStatusClosed   DealStatus = "closed"
DealStatusRejected DealStatus = "rejected"

ViewingStatusPending   ViewingStatus = "pending"
ViewingStatusConfirmed ViewingStatus = "confirmed"
ViewingStatusCancelled ViewingStatus = "cancelled"
ViewingStatusCompleted ViewingStatus = "completed"
)

type Deal struct {
ID          uuid.UUID
ListingID   uuid.UUID
BuyerID     uuid.UUID
AgentID     *uuid.UUID
Status      DealStatus
PriceAgreed float64
Currency    string
CreatedAt   time.Time
ClosedAt    *time.Time
}

type ViewingRequest struct {
ID          uuid.UUID
ListingID   uuid.UUID
BuyerID     uuid.UUID
Status      ViewingStatus
Comment     string
ScheduledAt *time.Time
CreatedAt   time.Time
}

type Favorite struct {
ID        uuid.UUID
UserID    uuid.UUID
ListingID uuid.UUID
CreatedAt time.Time
}

func NewDeal(listingID, buyerID uuid.UUID, agentID *uuid.UUID, price float64) (*Deal, error) {
if price <= 0 {
return nil, errors.New("цена сделки должна быть больше нуля")
}
return &Deal{
ID:          uuid.New(),
ListingID:   listingID,
BuyerID:     buyerID,
AgentID:     agentID,
Status:      DealStatusNew,
PriceAgreed: price,
Currency:    "RUB",
CreatedAt:   time.Now(),
}, nil
}

func (d *Deal) Approve() error {
if d.Status != DealStatusReview {
return errors.New("одобрить можно только сделку на рассмотрении")
}
d.Status = DealStatusApproved
return nil
}

func (d *Deal) Close() error {
if d.Status != DealStatusApproved {
return errors.New("закрыть можно только одобренную сделку")
}
now := time.Now()
d.Status = DealStatusClosed
d.ClosedAt = &now
return nil
}

func (d *Deal) Reject() error {
if d.Status == DealStatusClosed {
return errors.New("нельзя отклонить закрытую сделку")
}
d.Status = DealStatusRejected
return nil
}

func NewViewingRequest(listingID, buyerID uuid.UUID, comment string, scheduledAt *time.Time) (*ViewingRequest, error) {
if listingID == uuid.Nil {
return nil, errors.New("не указан id объявления")
}
if buyerID == uuid.Nil {
return nil, errors.New("не указан id покупателя")
}
return &ViewingRequest{
ID:          uuid.New(),
ListingID:   listingID,
BuyerID:     buyerID,
Status:      ViewingStatusPending,
Comment:     comment,
ScheduledAt: scheduledAt,
CreatedAt:   time.Now(),
}, nil
}

func (v *ViewingRequest) Confirm() error {
if v.Status != ViewingStatusPending {
return errors.New("подтвердить можно только ожидающий просмотр")
}
v.Status = ViewingStatusConfirmed
return nil
}

func (v *ViewingRequest) Cancel() error {
if v.Status == ViewingStatusCompleted {
return errors.New("нельзя отменить завершённый просмотр")
}
v.Status = ViewingStatusCancelled
return nil
}

func (v *ViewingRequest) Complete() error {
if v.Status != ViewingStatusConfirmed {
return errors.New("завершить можно только подтверждённый просмотр")
}
v.Status = ViewingStatusCompleted
return nil
}
