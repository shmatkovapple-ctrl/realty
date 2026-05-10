package entity

import (
"time"

"github.com/google/uuid"
)

type NotificationType string
type NotificationChannel string

const (
TypeViewingCreated   NotificationType = "viewing.created"
TypeViewingConfirmed NotificationType = "viewing.confirmed"
TypeViewingCancelled NotificationType = "viewing.cancelled"
TypeDealCreated      NotificationType = "deal.created"
TypeDealStatus       NotificationType = "deal.status_changed"
TypeListingPublished NotificationType = "listing.published"
TypeSystem           NotificationType = "system"

ChannelEmail NotificationChannel = "email"
ChannelInApp NotificationChannel = "in_app"
ChannelBoth  NotificationChannel = "both"
)

type Notification struct {
ID        uuid.UUID
UserID    uuid.UUID
Type      NotificationType
Title     string
Body      string
Link      string
IsRead    bool
CreatedAt time.Time
}

func NewNotification(userID uuid.UUID, t NotificationType, title, body, link string) *Notification {
return &Notification{
ID:        uuid.New(),
UserID:    userID,
Type:      t,
Title:     title,
Body:      body,
Link:      link,
IsRead:    false,
CreatedAt: time.Now(),
}
}

func (n *Notification) MarkAsRead() {
n.IsRead = true
}

var NotificationTitles = map[NotificationType]string{
TypeViewingCreated:   "Новая заявка на просмотр",
TypeViewingConfirmed: "Просмотр подтверждён",
TypeViewingCancelled: "Просмотр отменён",
TypeDealCreated:      "Новая сделка",
TypeDealStatus:       "Статус сделки изменён",
TypeListingPublished: "Объявление опубликовано",
TypeSystem:           "Системное уведомление",
}
