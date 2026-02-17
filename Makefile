.PHONY: help dev build test logs clean purge prod-up prod-down

help:
	@echo "EmpoweredPixels Development Commands"
	@echo ""
	@echo "  make dev          - Start development environment (frontend + backend + postgres)"
	@echo "  make prod-up      - Start production environment"
	@echo "  make prod-down    - Stop production environment"
	@echo "  make build        - Build production images"
	@echo "  make test         - Run Playwright tests against dev server"
	@echo "  make logs         - Tail logs from dev services"
	@echo "  make clean        - Stop dev services"
	@echo "  make purge        - Full cleanup: stop all, remove images, volumes, caches"
	@echo ""

dev:
	@echo "Starting development environment..."
	docker compose --profile dev up -d --build
	@echo ""
	@echo "Frontend dev server: http://localhost:5173"
	@echo "Backend API: http://localhost:49101"
	@echo "Database: localhost:5432 (user: postgres, pass: postgres)"
	@echo ""
	@echo "HMR enabled for both frontend and backend."

build:
	@echo "Building production images..."
	docker compose --profile prod build
	@echo "Done."

test:
	@echo "Running Playwright tests against dev server..."
	PLAYWRIGHT_BASE_URL=http://localhost:5173 npx playwright test tests/e2e/smoke.spec.ts --reporter=list

logs:
	docker compose --profile dev logs -f

clean:
	@echo "Stopping development environment..."
	docker compose --profile dev down

prod-up:
	@echo "Starting production environment..."
	docker compose --profile prod up -d --build

prod-down:
	@echo "Stopping production environment..."
	docker compose --profile prod down

purge: clean
	@echo "Purging all Docker artifacts..."
	docker compose down -v --rmi all
	docker volume prune -f -f --filter label=com.docker.compose.project=empoweredpixels
	@echo "Purge complete."
