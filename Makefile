.PHONY: help build up down restart logs clean test

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo '$(BLUE)Available commands:$(NC)'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

build: ## Build Docker images
	@echo "$(BLUE)Building Docker images...$(NC)"
	docker-compose build --no-cache

up: ## Start all services
	@echo "$(BLUE)Starting services...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Services started!$(NC)"
	@echo "API available at: http://localhost:3000"

down: ## Stop all services
	@echo "$(BLUE)Stopping services...$(NC)"
	docker-compose down
	@echo "$(GREEN)Services stopped!$(NC)"

restart: down up ## Restart all services

logs: ## Show logs (follow mode)
	docker-compose logs -f --tail=100

logs-app: ## Show app logs only
	docker-compose logs -f app --tail=100

logs-db: ## Show database logs only
	docker-compose logs -f db --tail=100

ps: ## Show running containers
	docker-compose ps

clean: ## Remove containers, volumes, and images
	@echo "$(RED)Warning: This will remove all data!$(NC)"
	docker-compose down -v
	docker system prune -f

rebuild: clean build up ## Clean rebuild and start

shell-app: ## Open shell in app container
	docker-compose exec app sh

shell-db: ## Open MySQL shell in database container
	docker-compose exec db mysql -u wallet_user -p wallet_api

health: ## Check service health status
	@echo "$(BLUE)Checking service health...$(NC)"
	@docker-compose ps
	@echo ""
	@echo "$(BLUE)Testing API endpoint...$(NC)"
	@curl -s http://localhost:3000/api/health || echo "$(RED)API not responding$(NC)"

dev: ## Start services for development
	@echo "$(BLUE)Starting development environment...$(NC)"
	docker-compose up

migrate-up: ## Run database migrations up
	@echo "$(BLUE)Running migrations...$(NC)"
	docker-compose exec app sh -c "migrate -path /app/migrations -database 'mysql://wallet_user:wallet_secure_password@tcp(db:3306)/wallet_api' up"

migrate-down: ## Run database migrations down
	@echo "$(RED)Rolling back migrations...$(NC)"
	docker-compose exec app sh -c "migrate -path /app/migrations -database 'mysql://wallet_user:wallet_secure_password@tcp(db:3306)/wallet_api' down"

stats: ## Show container resource usage
	docker-compose stats
