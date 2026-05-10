package event

import (
"time"

"github.com/google/uuid"
)

const (
EventListingCreated   = "listing.created"
EventListingPublished = "listing.published"
EventListingArchived  = "listing.archived"
EventListingSold      = "listing.sold"
EventListingDeleted   = "listing.deleted"
)

type ListingCreated struct {
ListingID uuid.UUID `json:"listing_id"`
SellerID  uuid.UUID `json:"seller_id"`
Type      string    `json:"type"`
Title     string    `json:"title"`
Price     float64   `json:"price"`
City      string    `json:"city"`
OccuredAt time.Time `json:"occured_at"`
}

type ListingPublished struct {
ListingID   uuid.UUID `json:"listing_id"`
SellerID    uuid.UUID `json:"seller_id"`
Type        string    `json:"type"`
Title       string    `json:"title"`
Price       float64   `json:"price"`
City        string    `json:"city"`
District    string    `json:"district"`
Lat         float64   `json:"lat"`
Lng         float64   `json:"lng"`
PreviewURL  string    `json:"preview_url"`
OccuredAt   time.Time `json:"occured_at"`
}

type ListingDeleted struct {
ListingID uuid.UUID `json:"listing_id"`
OccuredAt time.Time `json:"occured_at"`
}
