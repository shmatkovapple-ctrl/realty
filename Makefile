# ============================================================
#  Makefile — realty project (Windows / PowerShell)
# ============================================================

PROTO_SRC  = api/proto
PROTO_GEN  = api/gen
MIGRATIONS = migrations

DB_HOST     = localhost
DB_PORT     = 6432
DB_USER     = usr
DB_PASSWORD = pass
DB_NAME     = lets_goto_it
DB_SSL      = disable
DB_URL      = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL)

SERVICES = user-service listing-service deal-service search-service notification-service api-gateway

.DEFAULT_GOAL := help

.PHONY: help check-tools proto-gen proto-clean \
        migrate-up migrate-down migrate-status \
        migrate-create migrate-version migrate-force \
        build run-user run-listing run-deal \
        run-search run-notification run-gateway \
        infra-up infra-down infra-ps

# --- Help ------------------------------------------------------

help:
	@echo.
	@echo   realty -- dostupnye komandy
	@echo.
	@echo   Proto:
	@echo     make proto-gen              -- sgenerirovaty Go-kod iz vsex .proto
	@echo     make proto-clean            -- udalit sgenerirovannyi kod
	@echo.
	@echo   Migracii:
	@echo     make migrate-up             -- primenit vse novye migracii
	@echo     make migrate-down           -- otkatit poslednyuyu
	@echo     make migrate-status         -- tekushaya versiya sxemy
	@echo     make migrate-create name=X  -- sozdat novuyu migraciyu X
	@echo     make migrate-force v=X      -- vstavit versiyu X
	@echo.
	@echo   Sborka:
	@echo     make build                  -- sobrat vse servisy
	@echo.
	@echo   Zapusk servisov:
	@echo     make run-user               -- zapustit user-service
	@echo     make run-listing            -- zapustit listing-service
	@echo     make run-deal               -- zapustit deal-service
	@echo     make run-search             -- zapustit search-service
	@echo     make run-notification       -- zapustit notification-service
	@echo     make run-gateway            -- zapustit api-gateway
	@echo.
	@echo   Infrastruktura:
	@echo     make infra-up               -- zapustit vse docker konteinery
	@echo     make infra-down             -- ostanovit konteinery
	@echo     make infra-ps               -- status konteinerov
	@echo.
	@echo     make check-tools            -- proverit nalichie zavisimostei
	@echo.

# --- Проверка инструментов ------------------------------------

check-tools:
	@echo Proveryaem protoc...
	@protoc --version
	@echo Proveryaem protoc-gen-go...
	@protoc-gen-go --version
	@echo Proveryaem protoc-gen-go-grpc...
	@protoc-gen-go-grpc --version
	@echo Proveryaem migrate...
	@migrate -version
	@echo Vse instrumenty naideny.

# --- Proto генерация ------------------------------------------

proto-gen:
	@echo Generiruem Go-kod iz proto...
	@powershell -ExecutionPolicy Bypass -File scripts/proto-gen.ps1

proto-clean:
	@echo Udalyaem sgenerirovannyi kod...
	@powershell -Command "Get-ChildItem -Path $(PROTO_GEN) -Recurse -Filter '*.pb.go' | Remove-Item -Force"
	@echo Ochisheno.

# --- Миграции -------------------------------------------------

migrate-up:
	@echo Primenyaem migracii...
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" up
	@echo Migracii primeneny.

migrate-down:
	@echo Otkativaem poslednyuyu migraciyu...
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" down 1

migrate-status:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" version

migrate-version:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" version

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS) -seq $(name)
	@echo Faily sozdany v ./$(MIGRATIONS)

migrate-force:
	migrate -path $(MIGRATIONS) -database "$(DB_URL)" force $(v)

# --- Сборка ---------------------------------------------------

build:
	@echo Sobiraem vse servisy...
	@powershell -ExecutionPolicy Bypass -File scripts/build.ps1
	@echo Sborka zavershena.

# --- Запуск сервисов ------------------------------------------

run-user:
	go run realty/services/user-service/cmd

run-listing:
	go run realty/services/listing-service/cmd

run-deal:
	go run realty/services/deal-service/cmd

run-search:
	go run realty/services/search-service/cmd

run-notification:
	go run realty/services/notification-service/cmd

run-gateway:
	go run realty/services/api-gateway/cmd

# --- Инфраструктура -------------------------------------------

infra-up:
	docker compose up -d
	@echo Infrastruktura zapushena.

infra-down:
	docker compose down

infra-ps:
	docker compose ps

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-rebuild:
	docker compose down
	docker compose build --no-cache
	docker compose up -d