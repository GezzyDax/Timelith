.PHONY: help build up down restart logs clean setup

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Initial setup (copy .env.example to .env)
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file. Please edit it with your configuration."; \
	else \
		echo ".env file already exists."; \
	fi

build: ## Build all containers
	docker-compose build

up: ## Start all services
	docker-compose up -d

down: ## Stop all services
	docker-compose down

restart: ## Restart all services
	docker-compose restart

logs: ## Show logs from all services
	docker-compose logs -f

logs-rails: ## Show Rails logs
	docker-compose logs -f rails-app

logs-go: ## Show Go backend logs
	docker-compose logs -f go-backend

shell-rails: ## Open Rails console
	docker-compose exec rails-app bundle exec rails console

shell-go: ## Open Go backend shell
	docker-compose exec go-backend sh

db-migrate: ## Run Rails database migrations
	docker-compose exec rails-app bundle exec rails db:migrate

db-console: ## Open PostgreSQL console
	docker-compose exec postgres psql -U timelith -d timelith_production

redis-cli: ## Open Redis CLI
	docker-compose exec redis redis-cli

clean: ## Remove all containers, volumes, and images
	docker-compose down -v --rmi all

ps: ## Show running containers
	docker-compose ps

create-admin: ## Create admin user (interactive)
	docker-compose exec rails-app bundle exec rails runner "User.create!(email: 'admin@example.com', password: 'admin123', role: 'admin'); puts 'Admin user created: admin@example.com / admin123'"
