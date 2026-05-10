package postgres

import (
"context"
"fmt"

"github.com/google/uuid"
"github.com/jackc/pgx/v5/pgxpool"
"realty/services/notification-service/internal/domain/entity"
"realty/services/notification-service/internal/domain/repository"
)

type notificationRepository struct {
db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) repository.NotificationRepository {
return &notificationRepository{db: db}
}

func (r *notificationRepository) Save(ctx context.Context, n *entity.Notification) error {
_, err := r.db.Exec(ctx, `
INSERT INTO notifications (id, user_id, type, title, body, link, is_read, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`, n.ID, n.UserID, n.Type, n.Title, n.Body, n.Link, n.IsRead, n.CreatedAt)
if err != nil {
return fmt.Errorf("сохранение уведомления: %w", err)
}
return nil
}

func (r *notificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
n := &entity.Notification{}
err := r.db.QueryRow(ctx, `
SELECT id, user_id, type, title, body, link, is_read, created_at
FROM notifications WHERE id = $1
`, id).Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.Link, &n.IsRead, &n.CreatedAt)
if err != nil {
return nil, fmt.Errorf("поиск уведомления: %w", err)
}
return n, nil
}

func (r *notificationRepository) FindByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool, page, limit int) ([]*entity.Notification, int, int, error) {
offset := (page - 1) * limit

query := `
SELECT id, user_id, type, title, body, link, is_read, created_at
FROM notifications WHERE user_id = $1
`
args := []any{userID}

if unreadOnly {
query += " AND is_read = false"
}

query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
args = append(args, limit, offset)

rows, err := r.db.Query(ctx, query, args...)
if err != nil {
return nil, 0, 0, fmt.Errorf("поиск уведомлений: %w", err)
}
defer rows.Close()

var notifications []*entity.Notification
for rows.Next() {
n := &entity.Notification{}
if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.Link, &n.IsRead, &n.CreatedAt); err != nil {
return nil, 0, 0, fmt.Errorf("сканирование уведомления: %w", err)
}
notifications = append(notifications, n)
}

var total, unread int
r.db.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE user_id=$1`, userID).Scan(&total)
r.db.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE user_id=$1 AND is_read=false`, userID).Scan(&unread)

return notifications, total, unread, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id, userID uuid.UUID) error {
_, err := r.db.Exec(ctx, `
UPDATE notifications SET is_read=true WHERE id=$1 AND user_id=$2
`, id, userID)
if err != nil {
return fmt.Errorf("пометка уведомления как прочитанного: %w", err)
}
return nil
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) (int, error) {
result, err := r.db.Exec(ctx, `
UPDATE notifications SET is_read=true WHERE user_id=$1 AND is_read=false
`, userID)
if err != nil {
return 0, fmt.Errorf("пометка всех уведомлений: %w", err)
}
return int(result.RowsAffected()), nil
}
