package event

import (
"time"

"github.com/google/uuid"
)

const (
EventDealCreated         = "deal.created"
EventDealStatusChanged   = "deal.status_changed"
EventViewingCreated      = "viewing.created"
EventViewingStatusChanged = "viewing.status_changed"
)

type DealCreated struct {
DealID    uuid.UUID `json:"deal_id"`
ListingID uuid.UUID `json:"listing_id"`
BuyerID   uuid.UUID `json:"buyer_id"`
Price     float64   `json:"price"`
OccuredAt time.Time `json:"occured_at"`
}

type DealStatusChanged struct {
DealID    uuid.UUID `json:"deal_id"`
BuyerID   uuid.UUID `json:"buyer_id"`
Status    string    `json:"status"`
OccuredAt time.Time `json:"occured_at"`
}

type ViewingCreated struct {
ViewingID   uuid.UUID `json:"viewing_id"`
ListingID   uuid.UUID `json:"listing_id"`
BuyerID     uuid.UUID `json:"buyer_id"`
ScheduledAt *time.Time `json:"scheduled_at"`
OccuredAt   time.Time  `json:"occured_at"`
}

type ViewingStatusChanged struct {
ViewingID uuid.UUID `json:"viewing_id"`
BuyerID   uuid.UUID `json:"buyer_id"`
Status    string    `json:"status"`
OccuredAt time.Time `json:"occured_at"`
}
