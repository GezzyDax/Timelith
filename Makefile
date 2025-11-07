.PHONY: help dev-check quick-start test-all clean build up down logs

# Default target
help:
	@echo "Timelith Development Commands"
	@echo "=============================="
	@echo ""
	@echo "Development:"
	@echo "  make quick-start    - Start infrastructure (PostgreSQL, Redis)"
	@echo "  make dev-check      - Run all checks before committing"
	@echo "  make test-all       - Run all tests (Go, Docker builds, integration)"
	@echo ""
	@echo "Docker Compose:"
	@echo "  make build          - Build all Docker images"
	@echo "  make up             - Start all services"
	@echo "  make down           - Stop all services"
	@echo "  make logs           - Show logs from all services"
	@echo "  make restart        - Restart all services"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean          - Clean all build artifacts and caches"
	@echo "  make clean-docker   - Remove all Docker containers and volumes"
	@echo ""
	@echo "Versioning:"
	@echo "  make bump-version   - Manually bump version"
	@echo ""
	@echo "Backend (Go):"
	@echo "  make backend-build  - Build Go backend"
	@echo "  make backend-test   - Run Go tests"
	@echo "  make backend-lint   - Run Go linter"
	@echo "  make backend-run    - Run Go backend locally"
	@echo ""
	@echo "Frontend (Next.js):"
	@echo "  make web-install    - Install npm dependencies"
	@echo "  make web-dev        - Run Next.js in dev mode"
	@echo "  make web-build      - Build Next.js for production"
	@echo "  make web-lint       - Run ESLint"
	@echo ""

# Development scripts
dev-check:
	@./scripts/dev-check.sh

quick-start:
	@./scripts/quick-start.sh

test-all:
	@./scripts/test-all.sh

clean:
	@./scripts/clean-all.sh

bump-version:
	@./scripts/bump-version.sh

# Docker Compose
build:
	docker compose build --parallel

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

restart: down up

clean-docker:
	docker compose down -v
	docker system prune -f

# Backend commands
backend-build:
	cd go-backend && make build

backend-test:
	cd go-backend && go test -v -race ./...

backend-lint:
	cd go-backend && golangci-lint run

backend-run:
	cd go-backend && go run ./cmd/server

backend-fmt:
	cd go-backend && gofmt -w .

# Frontend commands
web-install:
	cd web-ui && npm ci

web-dev:
	cd web-ui && npm run dev

web-build:
	cd web-ui && npm run build

web-lint:
	cd web-ui && npm run lint

web-type-check:
	cd web-ui && npx tsc --noEmit

# CI commands (what runs in GitHub Actions)
ci-check: backend-lint backend-test web-lint web-type-check web-build build
