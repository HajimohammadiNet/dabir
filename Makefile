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