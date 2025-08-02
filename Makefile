.PHONY: help docker-up docker-down docker-build run migrate-up seed test lint clean

help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

docker-up: ## Start Docker containers (production mode)
	docker-compose up -d

run-dev: ## Start Docker containers (development mode with hot reload)
	docker-compose -f docker-compose.dev.yml up

stop: ## Stop Docker containers
	docker-compose down

docker-build: ## Build Docker images
	docker-compose build

docker-logs: ## View Docker container logs
	docker-compose logs -f

db-shell: ## Access PostgreSQL shell
	docker exec -it go-practice-1-postgres-1 psql -U devuser -d school_db

db-reset: ## Reset database (removes all data)
	docker-compose down -v
	docker-compose up -d

run: ## Run the application locally (requires local PostgreSQL)
	go run cmd/api/main.go

migrate-up: ## Run database migrations
	@echo "Migrations are run automatically on startup"

seed: ## Seed the database
	@echo "Database is seeded automatically on startup"

test: ## Run tests
	go test -v ./...

lint: ## Run golangci-lint
	golangci-lint run

clean: ## Clean build artifacts
	go clean
	rm -f main go-practice-1

deps: ## Download dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o main cmd/api/main.go