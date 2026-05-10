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
listingv1 "realty/api/gen/listing/v1"
grpcDelivery "realty/services/listing-service/internal/delivery/grpc"
kafkaInfra "realty/services/listing-service/internal/infrastructure/kafka"
minioInfra "realty/services/listing-service/internal/infrastructure/minio"
postgresInfra "realty/services/listing-service/internal/infrastructure/postgres"
"realty/services/listing-service/internal/usecase"
"realty/services/listing-service/pkg/config"
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

publisher := kafkaInfra.NewPublisher(cfg.KafkaURL)
defer publisher.Close()
log.Println("Kafka publisher создан")

storage, err := minioInfra.NewStorage(cfg.MinioURL, cfg.MinioUser, cfg.MinioPassword)
if err != nil {
log.Fatalf("подключение к MinIO: %v", err)
}

if err := storage.EnsureBucket(ctx); err != nil {
log.Fatalf("создание бакета MinIO: %v", err)
}
log.Println("MinIO подключён")

listingRepo := postgresInfra.NewListingRepository(db)
uc          := usecase.NewListingUseCase(listingRepo, publisher, storage)
handler     := grpcDelivery.NewListingHandler(uc)

grpcServer := grpc.NewServer()
listingv1.RegisterListingServiceServer(grpcServer, handler)
reflection.Register(grpcServer)

lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
if err != nil {
log.Fatalf("запуск listener: %v", err)
}

go func() {
log.Printf("listing-service gRPC запущен на порту %s", cfg.GRPCPort)
if err := grpcServer.Serve(lis); err != nil {
log.Fatalf("ошибка gRPC сервера: %v", err)
}
}()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

log.Println("завершение работы listing-service...")
grpcServer.GracefulStop()
}
