APP_NAME=paybridge-transaction-service

include .env
export $(shell sed 's/=.*//' .env)

.PHONY: run migrate-up migrate-down migrate-create build

# Run main server
run:
	go run ./cmd/server

# Apply all migrations
migrate-up:
	go run ./cmd/migrate up

# Rollback last migration
migrate-down:
	go run ./cmd/migrate down

# Create new migration file
migrate-create:
	@if [ -z "$(n)" ]; then \
		echo "Usage: make migrate-create n=create_wallets"; \
		exit 1; \
	fi
	migrate create -ext sql -dir ./migrations $(n)

# Build binary
build:
	go build -o bin/$(APP_NAME) main.go
