package config

import "os"

type Config struct {
	GRPCPort    string
	HTTPPort    string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
}

func Load() (*Config, error) {
	return &Config{
		GRPCPort:    getEnv("GRPC_PORT", "50051"),
		HTTPPort:    getEnv("HTTP_PORT", "8081"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://usr:pass@localhost:6432/lets_goto_it?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "supersecretkey_change_in_prod"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
