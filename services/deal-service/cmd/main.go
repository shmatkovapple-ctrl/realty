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
dealv1 "realty/api/gen/deal/v1"
grpcDelivery "realty/services/deal-service/internal/delivery/grpc"
kafkaInfra "realty/services/deal-service/internal/infrastructure/kafka"
postgresInfra "realty/services/deal-service/internal/infrastructure/postgres"
"realty/services/deal-service/internal/usecase"
"realty/services/deal-service/pkg/config"
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

dealRepo     := postgresInfra.NewDealRepository(db)
viewingRepo  := postgresInfra.NewViewingRepository(db)
favoriteRepo := postgresInfra.NewFavoriteRepository(db)

dealUC     := usecase.NewDealUseCase(dealRepo, publisher)
viewingUC  := usecase.NewViewingUseCase(viewingRepo, publisher)
favoriteUC := usecase.NewFavoriteUseCase(favoriteRepo)

handler := grpcDelivery.NewDealHandler(dealUC, viewingUC, favoriteUC)

grpcServer := grpc.NewServer()
dealv1.RegisterDealServiceServer(grpcServer, handler)
reflection.Register(grpcServer)

lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
if err != nil {
log.Fatalf("запуск listener: %v", err)
}

go func() {
log.Printf("deal-service gRPC запущен на порту %s", cfg.GRPCPort)
if err := grpcServer.Serve(lis); err != nil {
log.Fatalf("ошибка gRPC сервера: %v", err)
}
}()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

log.Println("завершение работы deal-service...")
grpcServer.GracefulStop()
}
