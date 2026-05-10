package config

import "os"

type Config struct {
	GRPCPort   string
	HTTPPort   string
	ElasticURL string
	RedisURL   string
	KafkaURL   string
}

func Load() (*Config, error) {
	return &Config{
		GRPCPort:   getEnv("GRPC_PORT", "50054"),
		HTTPPort:   getEnv("HTTP_PORT", "8084"),
		ElasticURL: getEnv("ELASTIC_URL", "http://localhost:9200"),
		RedisURL:   getEnv("REDIS_URL", "redis://localhost:6379"),
		KafkaURL:   getEnv("KAFKA_URL", "localhost:9092"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
