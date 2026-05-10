package main

import (
	"fmt"
	"log"

	"realty/services/notification-service/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ошибка загрузки конфига: %v", err)
	}
	fmt.Printf("notification-service запущен на порту %s\n", cfg.GRPCPort)
}
