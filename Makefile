APP_NAME=dabir-api
COMPOSE_FILE=deployments/docker-compose.yml

.PHONY: run
run:
	go run ./cmd/api

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: postgres-up
postgres-up:
	docker compose -f $(COMPOSE_FILE) up -d postgres

.PHONY: postgres-down
postgres-down:
	docker compose -f $(COMPOSE_FILE) down

.PHONY: postgres-logs
postgres-logs:
	docker compose -f $(COMPOSE_FILE) logs -f postgres

.PHONY: build
build:
	go build -o bin/$(APP_NAME) ./cmd/api

DB_URL=postgres://dabir:dabir_secret@localhost:5432/dabir?sslmode=disable

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

.PHONY: migrate-force
migrate-force:
	migrate -path migrations -database "$(DB_URL)" force $(VERSION)

.PHONY: migrate-version
migrate-version:
	migrate -path migrations -database "$(DB_URL)" version

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir migrations -seq $(NAME)