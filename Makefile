APP_NAME=dabir-api
COMPOSE_FILE=deployments/docker-compose.yml

DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= dabir
DB_PASSWORD ?= change-this-db-password
DB_NAME ?= dabir
DB_SSLMODE ?= disable

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: run
run:
	go run ./cmd/api

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o bin/$(APP_NAME) ./cmd/api

.PHONY: postgres-up
postgres-up:
	docker compose -f $(COMPOSE_FILE) up -d postgres

.PHONY: postgres-down
postgres-down:
	docker compose -f $(COMPOSE_FILE) down

.PHONY: postgres-logs
postgres-logs:
	docker compose -f $(COMPOSE_FILE) logs -f postgres

.PHONY: compose-up
compose-up:
	docker compose --env-file .env -f $(COMPOSE_FILE) up -d --build

.PHONY: compose-down
compose-down:
	docker compose --env-file .env -f $(COMPOSE_FILE) down

.PHONY: compose-logs
compose-logs:
	docker compose --env-file .env -f $(COMPOSE_FILE) logs -f

.PHONY: api-logs
api-logs:
	docker compose --env-file .env -f $(COMPOSE_FILE) logs -f api

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

.PHONY: migrate-version
migrate-version:
	migrate -path migrations -database "$(DB_URL)" version

.PHONY: migrate-force
migrate-force:
	migrate -path migrations -database "$(DB_URL)" force $(VERSION)

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir migrations -seq $(NAME)

.PHONY: dev-reset
dev-reset:
	docker compose --env-file .env -f $(COMPOSE_FILE) down -v
	docker compose --env-file .env -f $(COMPOSE_FILE) up -d postgres
	sleep 3
	migrate -path migrations -database "$(DB_URL)" up

COMPOSE_FILE=deployments/docker-compose.yml

web-logs:
	docker compose -f $(COMPOSE_FILE) --env-file .env logs -f web

compose-ps:
	docker compose -f $(COMPOSE_FILE) --env-file .env ps
