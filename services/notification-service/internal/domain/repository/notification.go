package repository

import (
"context"

"github.com/google/uuid"
"realty/services/notification-service/internal/domain/entity"
)

type NotificationRepository interface {
Save(ctx context.Context, n *entity.Notification) error
FindByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
FindByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool, page, limit int) ([]*entity.Notification, int, int, error)
MarkAsRead(ctx context.Context, id, userID uuid.UUID) error
MarkAllAsRead(ctx context.Context, userID uuid.UUID) (int, error)
}
