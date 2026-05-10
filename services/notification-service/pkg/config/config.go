package config

import "os"

type Config struct {
	GRPCPort    string
	HTTPPort    string
	DatabaseURL string
	KafkaURL    string
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
}

func Load() (*Config, error) {
	return &Config{
		GRPCPort:    getEnv("GRPC_PORT", "50055"),
		HTTPPort:    getEnv("HTTP_PORT", "8085"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://usr:pass@localhost:6432/lets_goto_it?sslmode=disable"),
		KafkaURL:    getEnv("KAFKA_URL", "localhost:9092"),
		SMTPHost:    getEnv("SMTP_HOST", "localhost"),
		SMTPPort:    getEnv("SMTP_PORT", "587"),
		SMTPUser:    getEnv("SMTP_USER", ""),
		SMTPPass:    getEnv("SMTP_PASS", ""),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
