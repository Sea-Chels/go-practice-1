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

db-tables: ## List all tables in the database
	docker exec go-practice-1-postgres-1 psql -U devuser -d school_db -c "\dt"

db-reset: ## Reset database (removes all data)
	docker-compose down -v
	docker-compose up -d

run: ## Run the application locally (requires local PostgreSQL)
	go run cmd/api/main.go

migrate-up: ## Run database migrations
	@echo "Migrations are run automatically on startup"

migrate-create: ## Create a new migration file (usage: make migrate-create name=create_table_name)
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a migration name. Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi; \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	seq=$$(ls migrations/*.sql 2>/dev/null | wc -l | xargs); \
	seq=$$(printf "%03d" $$(($$seq + 1))); \
	filename="migrations/$${seq}_$(name).sql"; \
	echo "-- Migration: $(name)" > $$filename; \
	echo "-- Created at: $$(date)" >> $$filename; \
	echo "" >> $$filename; \
	echo "-- Write your SQL migration here" >> $$filename; \
	echo "-- Example:" >> $$filename; \
	echo "-- CREATE TABLE example (" >> $$filename; \
	echo "--     id SERIAL PRIMARY KEY," >> $$filename; \
	echo "--     name VARCHAR(255) NOT NULL," >> $$filename; \
	echo "--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP" >> $$filename; \
	echo "-- );" >> $$filename; \
	echo "" >> $$filename; \
	echo "Created migration file: $$filename"; \
	echo "Edit this file with your SQL commands"

migrate-run: ## Run pending migrations while the server is running
	@echo "Running migrations..."; \
	go run cmd/migrate/main.go

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