# Simple Makefile for a Go project

# Load environment variables from .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Default target
all: build

# Build the application binary
build:
	@echo "Building..."
	go build -o main cmd/api/main.go

# Run the application
run:
	@echo "Running..."
	@go run cmd/api/main.go

# Clean the built binary
clean:
	@echo "Cleaning..."
	@rm -f main
	@rm -rf tmp

# Watch for changes and live reload using Air
watch:
	@echo "Cleaning up old builds and cache..."
	@rm -f main
	@rm -rf tmp/
	@mkdir -p tmp
	@echo "Checking Air installation..."
	@if ! command -v air > /dev/null; then \
		echo "Installing Air..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@echo "Testing build command..."
	@go build -o ./tmp/main ./cmd/api/main.go
	@echo "Build successful, starting Air..."
	@air

# Alternative watch command using basic go run with file monitoring
watch-simple:
	@echo "Starting simple file watcher..."
	@while true; do \
		echo "Building and running..."; \
		go run cmd/api/main.go & \
		PID=$$!; \
		inotifywait -e modify -r --exclude='tmp|\.git' . 2>/dev/null || \
		(echo "inotifywait not available, using basic sleep method..."; sleep 5); \
		echo "Files changed, restarting..."; \
		kill $$PID 2>/dev/null || true; \
		sleep 1; \
	done

# Database migration commands
db-migrate:
	@echo "Running database migrations..."
	@migrate -path migrations -database "$(DATABASE_URL)" up

db-rollback:
	@echo "Rolling back last migration..."
	@migrate -path migrations -database "$(DATABASE_URL)" down 1

db-reset:
	@echo "Resetting database..."
	@migrate -path migrations -database "$(DATABASE_URL)" down -all
	@migrate -path migrations -database "$(DATABASE_URL)" up

db-status:
	@echo "Checking migration status..."
	@migrate -path migrations -database "$(DATABASE_URL)" version

db-create-migration:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Development workflow helpers
dev-setup: db-migrate
	@echo "Database ready for development"

# Production deployment
deploy-db:
	@echo "Running production migrations..."
	@migrate -path migrations -database "$(DATABASE_URL)" up

# Debug command to check environment variables
debug-env:
	@echo "DATABASE_URL: $(DATABASE_URL)"
	@echo "First 20 chars: $(shell echo '$(DATABASE_URL)' | cut -c1-20)"

.PHONY: all build run clean watch db-migrate db-rollback db-reset db-status db-create-migration dev-setup deploy-db debug-env
