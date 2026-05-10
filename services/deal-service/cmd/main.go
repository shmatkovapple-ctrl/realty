package main

import (
	"fmt"
	"log"

	"realty/services/deal-service/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ошибка загрузки конфига: %v", err)
	}
	fmt.Printf("deal-service запущен на порту %s\n", cfg.GRPCPort)
}
