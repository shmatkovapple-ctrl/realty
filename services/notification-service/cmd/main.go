package main

import (
"context"
"fmt"
"log"
"net"
"os"
"os/signal"
"syscall"

"github.com/jackc/pgx/v5/pgxpool"
"google.golang.org/grpc"
"google.golang.org/grpc/reflection"
notificationv1 "realty/api/gen/notification/v1"
grpcDelivery "realty/services/notification-service/internal/delivery/grpc"
kafkaDelivery "realty/services/notification-service/internal/delivery/kafka"
kafkaInfra "realty/services/notification-service/internal/infrastructure/kafka"
postgresInfra "realty/services/notification-service/internal/infrastructure/postgres"
smtpInfra "realty/services/notification-service/internal/infrastructure/smtp"
"realty/services/notification-service/internal/usecase"
"realty/services/notification-service/pkg/config"
)

func main() {
cfg, err := config.Load()
if err != nil {
log.Fatalf("ошибка загрузки конфига: %v", err)
}

ctx := context.Background()

db, err := pgxpool.New(ctx, cfg.DatabaseURL)
if err != nil {
log.Fatalf("подключение к PostgreSQL: %v", err)
}
defer db.Close()

if err := db.Ping(ctx); err != nil {
log.Fatalf("ping PostgreSQL: %v", err)
}
log.Println("PostgreSQL подключён")

notifRepo   := postgresInfra.NewNotificationRepository(db)
emailSender := smtpInfra.NewEmailSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)
uc          := usecase.NewNotificationUseCase(notifRepo, emailSender)

eventHandler := kafkaDelivery.NewEventHandler(uc)

consumers := []*kafkaInfra.Consumer{
kafkaInfra.NewConsumer(cfg.KafkaURL, "deal.created",          "notification-service", eventHandler.HandleDealCreated),
kafkaInfra.NewConsumer(cfg.KafkaURL, "deal.status_changed",   "notification-service", eventHandler.HandleDealStatusChanged),
kafkaInfra.NewConsumer(cfg.KafkaURL, "viewing.created",       "notification-service", eventHandler.HandleViewingCreated),
kafkaInfra.NewConsumer(cfg.KafkaURL, "viewing.status_changed","notification-service", eventHandler.HandleViewingStatusChanged),
kafkaInfra.NewConsumer(cfg.KafkaURL, "listing.published",     "notification-service", eventHandler.HandleListingPublished),
}

for _, c := range consumers {
c.Start(ctx)
}
defer func() {
for _, c := range consumers {
c.Close()
}
}()

log.Println("Kafka consumers запущены — слушаем 5 топиков")

handler    := grpcDelivery.NewNotificationHandler(uc)
grpcServer := grpc.NewServer()
notificationv1.RegisterNotificationServiceServer(grpcServer, handler)
reflection.Register(grpcServer)

lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
if err != nil {
log.Fatalf("запуск listener: %v", err)
}

go func() {
log.Printf("notification-service gRPC запущен на порту %s", cfg.GRPCPort)
if err := grpcServer.Serve(lis); err != nil {
log.Fatalf("ошибка gRPC сервера: %v", err)
}
}()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

log.Println("завершение работы notification-service...")
grpcServer.GracefulStop()
}
