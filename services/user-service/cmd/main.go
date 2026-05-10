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
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	userv1 "realty/api/gen/user/v1"
	grpcDelivery "realty/services/user-service/internal/delivery/grpc"
	postgresInfra "realty/services/user-service/internal/infrastructure/postgres"
	redisInfra "realty/services/user-service/internal/infrastructure/redis"
	"realty/services/user-service/internal/usecase"
	"realty/services/user-service/pkg/config"
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

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("парсинг Redis URL: %v", err)
	}
	redisClient := redis.NewClient(opt)
	defer redisClient.Close()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("ping Redis: %v", err)
	}
	log.Println("Redis подключён")

	userRepo := postgresInfra.NewUserRepository(db)
	tokenCache := redisInfra.NewTokenCache(redisClient)
	uc := usecase.NewUserUseCase(userRepo, tokenCache, cfg.JWTSecret)
	handler := grpcDelivery.NewUserHandler(uc)

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("запуск listener: %v", err)
	}

	go func() {
		log.Printf("user-service gRPC запущен на порту %s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("ошибка gRPC сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("завершение работы user-service...")
	grpcServer.GracefulStop()
}
