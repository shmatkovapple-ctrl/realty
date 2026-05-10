package grpc

import (
"context"

"github.com/google/uuid"
"google.golang.org/grpc/codes"
"google.golang.org/grpc/status"
notificationv1 "realty/api/gen/notification/v1"
"realty/services/notification-service/internal/domain/entity"
"realty/services/notification-service/internal/usecase"
"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationHandler struct {
notificationv1.UnimplementedNotificationServiceServer
uc usecase.NotificationUseCase
}

func NewNotificationHandler(uc usecase.NotificationUseCase) *NotificationHandler {
return &NotificationHandler{uc: uc}
}

func (h *NotificationHandler) SendNotification(ctx context.Context, req *notificationv1.SendNotificationReq) (*notificationv1.SendNotificationRes, error) {
userID, err := uuid.Parse(req.UserId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный user_id")
}

if err := h.uc.Send(ctx, usecase.SendInput{
UserID:  userID,
Type:    entity.NotificationType(req.Type.String()),
Channel: entity.NotificationChannel(req.Channel.String()),
Payload: req.Payload,
}); err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &notificationv1.SendNotificationRes{Success: true}, nil
}

func (h *NotificationHandler) ListNotifications(ctx context.Context, req *notificationv1.ListNotificationsReq) (*notificationv1.ListNotificationsRes, error) {
userID, err := uuid.Parse(req.UserId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный user_id")
}

notifications, total, unread, err := h.uc.List(ctx, userID, req.UnreadOnly, int(req.Page), int(req.Limit))
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

var protoNotifications []*notificationv1.Notification
for _, n := range notifications {
protoNotifications = append(protoNotifications, &notificationv1.Notification{
Id:        n.ID.String(),
UserId:    n.UserID.String(),
Title:     n.Title,
Body:      n.Body,
Link:      n.Link,
IsRead:    n.IsRead,
CreatedAt: timestamppb.New(n.CreatedAt),
})
}

return &notificationv1.ListNotificationsRes{
Notifications: protoNotifications,
Total:         int32(total),
Unread:        int32(unread),
}, nil
}

func (h *NotificationHandler) MarkAsRead(ctx context.Context, req *notificationv1.MarkAsReadReq) (*notificationv1.MarkAsReadRes, error) {
notifID, err := uuid.Parse(req.NotificationId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный notification_id")
}

userID, err := uuid.Parse(req.UserId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный user_id")
}

if err := h.uc.MarkAsRead(ctx, notifID, userID); err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &notificationv1.MarkAsReadRes{Success: true}, nil
}

func (h *NotificationHandler) MarkAllAsRead(ctx context.Context, req *notificationv1.MarkAllAsReadReq) (*notificationv1.MarkAllAsReadRes, error) {
userID, err := uuid.Parse(req.UserId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный user_id")
}

count, err := h.uc.MarkAllAsRead(ctx, userID)
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &notificationv1.MarkAllAsReadRes{MarkedCount: int32(count)}, nil
}
