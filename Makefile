.PHONY: help docker-up docker-down migrate-up migrate-down migrate-create sqlc-generate run build test clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

docker-up: ## Start Docker containers
	docker-compose up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## View Docker container logs
	docker-compose logs -f

migrate-up: ## Run database migrations
	migrate -path db/schema -database "postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable" up

migrate-down: ## Rollback last database migration
	migrate -path db/schema -database "postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable" down 1

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	migrate -path db/schema -database "postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable" force $(version)

migrate-up-prod: ## Run migrations on production (requires PRODUCTION_DATABASE_URL env var)
	@if [ -z "$$PRODUCTION_DATABASE_URL" ]; then \
		echo "❌ Error: PRODUCTION_DATABASE_URL environment variable is not set"; \
		echo "Example: export PRODUCTION_DATABASE_URL='postgres://user:pass@host:5432/dbname?sslmode=require'"; \
		exit 1; \
	fi
	@echo "🔄 Creating automatic backup before production migration..."
	@./scripts/backup-prod-db.sh
	@echo ""
	@echo "⚠️  Running migrations on PRODUCTION database..."
	@echo "Database: $$PRODUCTION_DATABASE_URL" | sed 's/:\/\/.*@/:\/\/***@/'
	@read -p "Continue? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	migrate -path db/schema -database "$$PRODUCTION_DATABASE_URL" up
	@echo "✅ Production migrations completed"

migrate-down-prod: ## Rollback production migration (requires PRODUCTION_DATABASE_URL env var)
	@if [ -z "$$PRODUCTION_DATABASE_URL" ]; then \
		echo "❌ Error: PRODUCTION_DATABASE_URL environment variable is not set"; \
		exit 1; \
	fi
	@echo "⚠️  WARNING: Rolling back production migration!"
	@read -p "Are you sure? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	migrate -path db/schema -database "$$PRODUCTION_DATABASE_URL" down 1

migrate-status-prod: ## Check production migration status (requires PRODUCTION_DATABASE_URL env var)
	@if [ -z "$$PRODUCTION_DATABASE_URL" ]; then \
		echo "❌ Error: PRODUCTION_DATABASE_URL environment variable is not set"; \
		exit 1; \
	fi
	migrate -path db/schema -database "$$PRODUCTION_DATABASE_URL" version

migrate-create: ## Create new migration (usage: make migrate-create name=create_table)
	migrate create -ext sql -dir db/schema -seq $(name)

sqlc-generate: ## Generate sqlc code
	sqlc generate

run: ## Run the application
	LOG_LEVEL=debug go run cmd/api/main.go

build: ## Build the application
	@mkdir -p bin
	go build -o bin/api cmd/api/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

deps: ## Download dependencies
	go mod download
	go mod tidy

dev: docker-up migrate-up sqlc-generate run ## Start development environment

setup: deps docker-up migrate-up sqlc-generate ## Initial project setup
	@echo "Setup complete! Run 'make run' to start the server."

seed: ## Seed database with fake data (usage: make seed users=50 orders=200)
	go run cmd/seeder/main.go -users=$(or $(users),50) -orders=$(or $(orders),200)

seed-clear: ## Clear and reseed database with fake data
	go run cmd/seeder/main.go -clear -users=$(or $(users),50) -orders=$(or $(orders),200)

# Database backup/restore
db-backup: ## Backup database to file (usage: make db-backup file=backups/backup.sql)
	@mkdir -p backups
	@./scripts/backup-db.sh $(or $(file),backups/db-backup-$(shell date +%Y%m%d-%H%M%S).sql)

db-restore: ## Restore database from file (usage: make db-restore file=backups/backup.sql)
	@./scripts/restore-db.sh $(file)
