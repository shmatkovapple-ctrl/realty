package main

import (
	"fmt"
	"log"

	"realty/services/api-gateway/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ошибка загрузки конфига: %v", err)
	}
	fmt.Printf("api-gateway запущен на порту %s\n", cfg.HTTPPort)
}
