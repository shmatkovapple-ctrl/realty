package usecase

import (
"context"
"fmt"
"log"

"github.com/google/uuid"
"realty/services/notification-service/internal/domain/entity"
"realty/services/notification-service/internal/domain/repository"
smtpInfra "realty/services/notification-service/internal/infrastructure/smtp"
)

type notificationUseCase struct {
repo        repository.NotificationRepository
emailSender *smtpInfra.EmailSender
}

func NewNotificationUseCase(
repo repository.NotificationRepository,
emailSender *smtpInfra.EmailSender,
) NotificationUseCase {
return &notificationUseCase{
repo:        repo,
emailSender: emailSender,
}
}

func (uc *notificationUseCase) Send(ctx context.Context, input SendInput) error {
title := entity.NotificationTitles[input.Type]
if title == "" {
title = "Уведомление"
}

body := buildBody(input.Type, input.Payload)
link := buildLink(input.Type, input.Payload)

n := entity.NewNotification(input.UserID, input.Type, title, body, link)

if err := uc.repo.Save(ctx, n); err != nil {
return fmt.Errorf("сохранение уведомления: %w", err)
}

if input.Channel == entity.ChannelEmail || input.Channel == entity.ChannelBoth {
go func() {
if email, ok := input.Payload["email"]; ok && email != "" {
if err := uc.emailSender.Send(email, title, body); err != nil {
log.Printf("ошибка отправки email: %v", err)
}
}
}()
}

return nil
}

func (uc *notificationUseCase) List(ctx context.Context, userID uuid.UUID, unreadOnly bool, page, limit int) ([]*entity.Notification, int, int, error) {
if page < 1 {
page = 1
}
if limit < 1 || limit > 100 {
limit = 20
}
return uc.repo.FindByUser(ctx, userID, unreadOnly, page, limit)
}

func (uc *notificationUseCase) MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error {
return uc.repo.MarkAsRead(ctx, notificationID, userID)
}

func (uc *notificationUseCase) MarkAllAsRead(ctx context.Context, userID uuid.UUID) (int, error) {
return uc.repo.MarkAllAsRead(ctx, userID)
}

func buildBody(t entity.NotificationType, payload map[string]string) string {
switch t {
case entity.TypeViewingCreated:
return fmt.Sprintf("Новая заявка на просмотр объявления %s", payload["listing_id"])
case entity.TypeViewingConfirmed:
return "Ваша заявка на просмотр подтверждена"
case entity.TypeViewingCancelled:
return "Ваша заявка на просмотр отменена"
case entity.TypeDealCreated:
return fmt.Sprintf("Создана новая сделка на сумму %s руб.", payload["price"])
case entity.TypeDealStatus:
return fmt.Sprintf("Статус вашей сделки изменён на: %s", payload["status"])
case entity.TypeListingPublished:
return "Ваше объявление успешно опубликовано"
default:
if msg, ok := payload["message"]; ok {
return msg
}
return "У вас новое уведомление"
}
}

func buildLink(t entity.NotificationType, payload map[string]string) string {
switch t {
case entity.TypeViewingCreated, entity.TypeViewingConfirmed, entity.TypeViewingCancelled:
if id, ok := payload["listing_id"]; ok {
return fmt.Sprintf("/listings/%s", id)
}
case entity.TypeDealCreated, entity.TypeDealStatus:
if id, ok := payload["deal_id"]; ok {
return fmt.Sprintf("/deals/%s", id)
}
case entity.TypeListingPublished:
if id, ok := payload["listing_id"]; ok {
return fmt.Sprintf("/listings/%s", id)
}
}
return "/"
}
