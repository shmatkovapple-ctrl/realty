package entity

import (
"errors"
"time"

"github.com/google/uuid"
)

type ListingType string
type ListingStatus string

const (
TypeApartment  ListingType = "apartment"
TypeHouse      ListingType = "house"
TypeCommercial ListingType = "commercial"
TypeLand       ListingType = "land"

StatusDraft     ListingStatus = "draft"
StatusPublished ListingStatus = "published"
StatusArchived  ListingStatus = "archived"
StatusSold      ListingStatus = "sold"
)

type Listing struct {
ID          uuid.UUID
SellerID    uuid.UUID
AgentID     *uuid.UUID
Type        ListingType
Status      ListingStatus
Title       string
Description string
Price       float64
Currency    string
AreaSqm     float64
Rooms       int32
Floor       int32
FloorsTotal int32
Address     Address
MediaURLs   []string
PublishedAt *time.Time
CreatedAt   time.Time
UpdatedAt   time.Time
}

type Address struct {
Country  string
City     string
District string
Street   string
Building string
Lat      float64
Lng      float64
}

func NewListing(sellerID uuid.UUID, t ListingType, title string, price float64, area float64) (*Listing, error) {
if title == "" {
return nil, errors.New("заголовок не может быть пустым")
}
if price <= 0 {
return nil, errors.New("цена должна быть больше нуля")
}
if area <= 0 {
return nil, errors.New("площадь должна быть больше нуля")
}
if !t.IsValid() {
return nil, errors.New("недопустимый тип объекта")
}

now := time.Now()
return &Listing{
ID:        uuid.New(),
SellerID:  sellerID,
Type:      t,
Status:    StatusDraft,
Title:     title,
Price:     price,
Currency:  "RUB",
AreaSqm:   area,
CreatedAt: now,
UpdatedAt: now,
}, nil
}

func (l *Listing) Publish() error {
if l.Status == StatusPublished {
return errors.New("объявление уже опубликовано")
}
if l.Title == "" {
return errors.New("нельзя опубликовать объявление без заголовка")
}
if l.Address.City == "" {
return errors.New("нельзя опубликовать объявление без города")
}
if len(l.MediaURLs) == 0 {
return errors.New("нельзя опубликовать объявление без фотографий")
}

now := time.Now()
l.Status = StatusPublished
l.PublishedAt = &now
l.UpdatedAt = now
return nil
}

func (l *Listing) Archive() error {
if l.Status == StatusArchived {
return errors.New("объявление уже архивировано")
}
l.Status = StatusArchived
l.UpdatedAt = time.Now()
return nil
}

func (l *Listing) MarkAsSold() error {
if l.Status != StatusPublished {
return errors.New("пометить как проданное можно только опубликованное объявление")
}
l.Status = StatusSold
l.UpdatedAt = time.Now()
return nil
}

func (l *Listing) IsPublished() bool {
return l.Status == StatusPublished
}

func (t ListingType) IsValid() bool {
switch t {
case TypeApartment, TypeHouse, TypeCommercial, TypeLand:
return true
}
return false
}
