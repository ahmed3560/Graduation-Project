.PHONY: build test run up down migrate-up migrate-down migrate-version

build:
	go build -o bin/api ./cmd/api && go build -o bin/worker ./cmd/worker && go build -o bin/migrate ./cmd/migrate

test:
	go vet ./... && go test ./...

run:
	go run ./cmd/api

up:
	docker compose up -d --build

down:
	docker compose down

# Needs DATABASE_URL in the environment.
migrate-up:
	go run ./cmd/migrate up

# Usage: make migrate-down N=1
migrate-down:
	go run ./cmd/migrate down $(N)

migrate-version:
	go run ./cmd/migrate version
