package usecase

import (
"context"

"github.com/google/uuid"
"realty/services/notification-service/internal/domain/entity"
)

type NotificationUseCase interface {
Send(ctx context.Context, input SendInput) error
List(ctx context.Context, userID uuid.UUID, unreadOnly bool, page, limit int) ([]*entity.Notification, int, int, error)
MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error
MarkAllAsRead(ctx context.Context, userID uuid.UUID) (int, error)
}

type SendInput struct {
UserID  uuid.UUID
Type    entity.NotificationType
Channel entity.NotificationChannel
Payload map[string]string
}
