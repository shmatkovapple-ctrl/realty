package config

import "os"

type Config struct {
GRPCPort      string
HTTPPort      string
DatabaseURL   string
RedisURL      string
KafkaURL      string
MinioURL      string
MinioUser     string
MinioPassword string
}

func Load() (*Config, error) {
return &Config{
GRPCPort:      getEnv("GRPC_PORT", "50052"),
HTTPPort:      getEnv("HTTP_PORT", "8082"),
DatabaseURL:   getEnv("DATABASE_URL", "postgres://usr:pass@localhost:6432/lets_goto_it?sslmode=disable"),
RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379"),
KafkaURL:      getEnv("KAFKA_URL", "localhost:9092"),
MinioURL:      getEnv("MINIO_URL", "localhost:9000"),
MinioUser:     getEnv("MINIO_USER", "minioadmin"),
MinioPassword: getEnv("MINIO_PASSWORD", "minioadmin123"),
}, nil
}

func getEnv(key, defaultVal string) string {
if val := os.Getenv(key); val != "" {
return val
}
return defaultVal
}
