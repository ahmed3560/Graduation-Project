.PHONY: build test run up down

build:
	go build -o bin/api ./cmd/api && go build -o bin/worker ./cmd/worker

test:
	go vet ./... && go test ./...

run:
	go run ./cmd/api

up:
	docker compose up -d --build

down:
	docker compose down
