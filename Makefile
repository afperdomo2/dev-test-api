.PHONY: install dev build swagger clean run fmt fmt-check vet check setup test test-cover db-seed \
        staged-check \
        migrate-up migrate-down migrate-status migrate-create \
        fe-install fe-dev fe-build fe-lint fe-format-check fe-format fe-check fe-test fe-test-cover \
        e2e e2e-ui e2e-report \
        mobile-start mobile-android mobile-ios mobile-web mobile-typecheck mobile-lint

# ── Backend ──────────────────────────────────────────────

install:
	@echo "📦 Installing Go dependencies..."
	go mod tidy
	@echo "✅ Go dependencies installed"

dev:
	@echo "🔥 Starting Go live reload server..."
	go tool air

build:
	@echo "📦 Building Go binary..."
	go build -o ./tmp/main.exe .
	@echo "✅ Binary built at ./tmp/main.exe"

run:
	@echo "🚀 Running Go server..."
	go run main.go

swagger:
	@echo "📖 Generating Swagger docs..."
	go tool swag init -g main.go
	@echo "✅ Docs generated in docs/"

clean:
	@echo "🧹 Cleaning..."
	rm -rf tmp/ docs/
	@echo "✅ Cleaned"

fmt:
	@echo "🔧 Formatting Go..."
	go fmt ./...

fmt-check:
	@echo "🔍 Checking Go formatting..."
	bash scripts/check-fmt.sh

vet:
	@echo "🔍 Vetting Go..."
	go vet ./...

check: fmt vet
	@echo "✅ All Go checks passed"

test:
	@echo "🧪 Running unit tests..."
	go test -short -count=1 -timeout 60s ./...
	@echo "✅ Unit tests passed"

test-cover:
	@echo "🧪 Running tests with coverage..."
	go test -short -count=1 -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out | tail -1
	@echo "✅ Coverage report at coverage.out"

setup:
	@echo "🔗 Installing git hooks..."
	go tool lefthook install
	@echo "✅ Hooks installed"

staged-check:
	@echo "🔎 Verificando archivos en staging (lefthook pre-commit)..."
	@go tool lefthook run pre-commit

db-seed:
	@echo "🌱 Seeding 50 preguntas..."
	@bash scripts/seed-db.sh 2>nul || powershell -ExecutionPolicy Bypass -File scripts/seed-db.ps1
	@echo "✅ Seed completado"

# ── Migrations (goose) ────────────────────────────────

MIGRATE_DIR = internal/database/migrations

# La CLI goose se ejecuta via script por SO (mismo patron que db-seed):
#   scripts/migrate.sh  → Unix / Git Bash / WSL
#   scripts/migrate.ps1 → Windows PowerShell
ifeq ($(OS),Windows_NT)
MIGRATE_RUN = powershell -ExecutionPolicy Bypass -File scripts/migrate.ps1
else
MIGRATE_RUN = bash scripts/migrate.sh
endif

# Carga .env y ejecuta la CLI goose (scripts/migrate.*). Comandos:
#   make migrate-up          — aplica migraciones pendientes
#   make migrate-down        — revierte la última migración (down)
#   make migrate-status      — estado de las migraciones
#   make migrate-create name=<desc> — crea NNNNN_<desc>.sql (up+down)
migrate-up:
	@echo "🛢️ Aplicando migraciones pendientes..."
	@$(MIGRATE_RUN) up
	@echo "✅ Migraciones aplicadas"

migrate-down:
	@echo "🛢️ Revirtiendo última migración..."
	@$(MIGRATE_RUN) down
	@echo "✅ Migración revertida"

migrate-status:
	@echo "🛢️ Estado de las migraciones:"
	@$(MIGRATE_RUN) status

migrate-create:
	@$(MIGRATE_RUN) create "$(name)"
	@echo "✅ Migración creada en $(MIGRATE_DIR)/"

# ── Frontend ────────────────────────────────────────────

fe-install:
	@echo "📦 Installing frontend dependencies..."
	cd apps/web && pnpm install
	@echo "✅ Frontend dependencies installed"

fe-dev:
	@echo "🔥 Starting frontend dev server..."
	cd apps/web && pnpm dev

fe-build:
	@echo "📦 Building frontend..."
	cd apps/web && pnpm build
	@echo "✅ Frontend built at apps/web/dist/"

fe-lint:
	@echo "🔍 Linting frontend..."
	cd apps/web && pnpm lint
	@echo "✅ Frontend lint OK"

fe-format-check:
	@echo "🎨 Checking frontend formatting..."
	cd apps/web && pnpm prettier --check src/
	@echo "✅ Frontend format OK"

fe-format:
	@echo "🎨 Formatting frontend..."
	cd apps/web && pnpm prettier --write src/
	@echo "✅ Frontend formatted"

fe-check:
	@echo "🔍 Type-checking frontend..."
	cd apps/web && pnpm type-check
	@echo "✅ Frontend type-check OK"

fe-test:
	@echo "🧪 Running frontend unit tests..."
	cd apps/web && pnpm test:run
	@echo "✅ Frontend unit tests passed"

fe-test-cover:
	@echo "🧪 Running frontend tests with coverage..."
	cd apps/web && pnpm test:cover
	@echo "✅ Frontend coverage report"

# ── E2E (Playwright) ───────────────────────────────────
# Requiere DB dedicada y limpia + backend corriendo:
#   DB_NAME=devtest_e2e make migrate-up
#   DB_NAME=devtest_e2e make run
# El web lo levanta Playwright (o reutiliza make fe-dev).

e2e:
	@echo "🧪 Running E2E (Playwright)..."
	@pnpm --filter dev-test-frontend test:e2e

e2e-ui:
	@echo "🎭 Opening Playwright UI..."
	@pnpm --filter dev-test-frontend test:e2e:ui

e2e-report:
	@echo "📊 Opening Playwright report..."
	@pnpm --filter dev-test-frontend test:e2e:report

# ── Mobile (apps/mobile/) ──────────────────────────────

mobile-start:
	@echo "🔥 Starting Expo dev server..."
	pnpm --filter mobile start

mobile-android:
	@echo "🤖 Starting Expo on Android..."
	pnpm --filter mobile android

mobile-ios:
	@echo "🍎 Starting Expo on iOS..."
	pnpm --filter mobile ios

mobile-web:
	@echo "🌐 Starting Expo on web..."
	pnpm --filter mobile web

mobile-typecheck:
	@echo "🔍 Type-checking mobile..."
	pnpm --filter mobile typecheck
	@echo "✅ Mobile type-check OK"

mobile-lint:
	@echo "🔍 Linting mobile..."
	pnpm --filter mobile lint
	@echo "✅ Mobile lint OK"
