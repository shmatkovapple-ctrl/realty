package config

import "os"

type Config struct {
	GRPCPort    string
	HTTPPort    string
	DatabaseURL string
	KafkaURL    string
}

func Load() (*Config, error) {
	return &Config{
		GRPCPort:    getEnv("GRPC_PORT", "50053"),
		HTTPPort:    getEnv("HTTP_PORT", "8083"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://usr:pass@localhost:6432/lets_goto_it?sslmode=disable"),
		KafkaURL:    getEnv("KAFKA_URL", "localhost:9092"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
