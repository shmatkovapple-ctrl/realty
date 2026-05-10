package main

import (
"context"
"fmt"
"log"
"net/http"
"os"
"os/signal"
"syscall"
"time"

"github.com/redis/go-redis/v9"
grpcClients "realty/services/api-gateway/internal/infrastructure/grpc"
httpDelivery "realty/services/api-gateway/internal/delivery/http"
"realty/services/api-gateway/internal/middleware"
"realty/services/api-gateway/pkg/config"
)

func main() {
cfg, err := config.Load()
if err != nil {
log.Fatalf("ошибка загрузки конфига: %v", err)
}

clients, err := grpcClients.NewClients(
cfg.UserServiceAddr,
cfg.ListingServiceAddr,
cfg.DealServiceAddr,
cfg.SearchServiceAddr,
cfg.NotificationServiceAddr,
)
if err != nil {
log.Fatalf("подключение к сервисам: %v", err)
}
defer clients.Close()
log.Println("gRPC клиенты подключены")

opt, err := redis.ParseURL(cfg.RedisURL)
if err != nil {
log.Fatalf("парсинг Redis URL: %v", err)
}
redisClient := redis.NewClient(opt)
defer redisClient.Close()
log.Println("Redis подключён")

authMW, err := middleware.NewAuthMiddleware(cfg.UserServiceAddr)
if err != nil {
log.Fatalf("создание auth middleware: %v", err)
}
rateLimiter := middleware.NewRateLimiter(redisClient, 100, time.Minute)

userHandler         := httpDelivery.NewUserHandler(clients.User)
listingHandler      := httpDelivery.NewListingHandler(clients.Listing)
searchHandler       := httpDelivery.NewSearchHandler(clients.Search)
dealHandler         := httpDelivery.NewDealHandler(clients.Deal)
notificationHandler := httpDelivery.NewNotificationHandler(clients.Notification)

router := httpDelivery.NewRouter(
authMW,
rateLimiter,
userHandler,
listingHandler,
searchHandler,
dealHandler,
notificationHandler,
)

server := &http.Server{
Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
Handler:      router,
ReadTimeout:  15 * time.Second,
WriteTimeout: 15 * time.Second,
IdleTimeout:  60 * time.Second,
}

go func() {
log.Printf("api-gateway HTTP запущен на порту %s", cfg.HTTPPort)
if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
log.Fatalf("ошибка HTTP сервера: %v", err)
}
}()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

log.Println("завершение работы api-gateway...")
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
server.Shutdown(ctx)
}
