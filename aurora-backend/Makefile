.PHONY: dev build run test lint clean migrate-up migrate-down docker-build docker-up docker-down

APP_NAME = aurora
BUILD_DIR = ./build
CMD_DIR = ./cmd/aurora

dev:
	go run $(CMD_DIR)/main.go -config configs/config.example.yaml

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

run: build
	$(BUILD_DIR)/$(APP_NAME) -config configs/config.yaml

test:
	go test -v -race -cover ./...

lint:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR)

migrate-up:
	migrate -path internal/adapter/postgres/migrations -database "postgres://aurora:aurora@localhost:5432/aurora?sslmode=disable" up

migrate-down:
	migrate -path internal/adapter/postgres/migrations -database "postgres://aurora:aurora@localhost:5432/aurora?sslmode=disable" down

docker-build:
	docker build -t aurora-backend -f deployments/Dockerfile .

docker-up:
	docker compose -f deployments/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker-compose.yml down

docker-logs:
	docker compose -f deployments/docker-compose.yml logs -f
