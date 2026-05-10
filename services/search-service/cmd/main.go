package main

import (
"context"
"fmt"
"log"
"net"
"os"
"os/signal"
"syscall"

"github.com/elastic/go-elasticsearch/v8"
"github.com/redis/go-redis/v9"
"google.golang.org/grpc"
"google.golang.org/grpc/reflection"
searchv1 "realty/api/gen/search/v1"
grpcDelivery "realty/services/search-service/internal/delivery/grpc"
kafkaDelivery "realty/services/search-service/internal/delivery/kafka"
elasticInfra "realty/services/search-service/internal/infrastructure/elastic"
kafkaInfra "realty/services/search-service/internal/infrastructure/kafka"
redisInfra "realty/services/search-service/internal/infrastructure/redis"
"realty/services/search-service/internal/usecase"
"realty/services/search-service/pkg/config"
)

func main() {
cfg, err := config.Load()
if err != nil {
log.Fatalf("ошибка загрузки конфига: %v", err)
}

ctx := context.Background()

esClient, err := elasticsearch.NewClient(elasticsearch.Config{
Addresses: []string{cfg.ElasticURL},
})
if err != nil {
log.Fatalf("подключение к Elasticsearch: %v", err)
}

if err := elasticInfra.EnsureIndex(ctx, esClient); err != nil {
log.Fatalf("создание индекса Elasticsearch: %v", err)
}
log.Println("Elasticsearch подключён")

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

searchRepo  := elasticInfra.NewSearchRepository(esClient)
searchCache := redisInfra.NewSearchCache(redisClient)
uc          := usecase.NewSearchUseCase(searchRepo, searchCache)

eventHandler := kafkaDelivery.NewEventHandler(uc)

publishedConsumer := kafkaInfra.NewConsumer(
cfg.KafkaURL,
"listing.published",
"search-service",
eventHandler.HandleListingPublished,
)
publishedConsumer.Start(ctx)
defer publishedConsumer.Close()

deletedConsumer := kafkaInfra.NewConsumer(
cfg.KafkaURL,
"listing.deleted",
"search-service",
eventHandler.HandleListingDeleted,
)
deletedConsumer.Start(ctx)
defer deletedConsumer.Close()

log.Println("Kafka consumers запущены")

handler    := grpcDelivery.NewSearchHandler(uc)
grpcServer := grpc.NewServer()
searchv1.RegisterSearchServiceServer(grpcServer, handler)
reflection.Register(grpcServer)

lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
if err != nil {
log.Fatalf("запуск listener: %v", err)
}

go func() {
log.Printf("search-service gRPC запущен на порту %s", cfg.GRPCPort)
if err := grpcServer.Serve(lis); err != nil {
log.Fatalf("ошибка gRPC сервера: %v", err)
}
}()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

log.Println("завершение работы search-service...")
grpcServer.GracefulStop()
}
