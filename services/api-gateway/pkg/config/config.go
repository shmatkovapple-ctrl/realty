package config

import "os"

type Config struct {
HTTPPort                string
RedisURL                string
UserServiceAddr         string
ListingServiceAddr      string
DealServiceAddr         string
SearchServiceAddr       string
NotificationServiceAddr string
}

func Load() (*Config, error) {
return &Config{
HTTPPort:                getEnv("HTTP_PORT", "8080"),
RedisURL:                getEnv("REDIS_URL", "redis://localhost:6379"),
UserServiceAddr:         getEnv("USER_SERVICE_ADDR", "localhost:50051"),
ListingServiceAddr:      getEnv("LISTING_SERVICE_ADDR", "localhost:50052"),
DealServiceAddr:         getEnv("DEAL_SERVICE_ADDR", "localhost:50053"),
SearchServiceAddr:       getEnv("SEARCH_SERVICE_ADDR", "localhost:50054"),
NotificationServiceAddr: getEnv("NOTIFICATION_SERVICE_ADDR", "localhost:50055"),
}, nil
}

func getEnv(key, defaultVal string) string {
if val := os.Getenv(key); val != "" {
return val
}
return defaultVal
}
