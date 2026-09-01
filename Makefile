.PHONY: install dev build swagger clean run fmt vet check setup test test-cover db-seed \
        fe-install fe-dev fe-build fe-lint fe-check fe-test fe-test-cover

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
	swag init -g main.go
	@echo "✅ Docs generated in docs/"

clean:
	@echo "🧹 Cleaning..."
	rm -rf tmp/ docs/
	@echo "✅ Cleaned"

fmt:
	@echo "🔧 Formatting Go..."
	go fmt ./...

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

db-seed:
	@echo "🌱 Seeding 50 preguntas..."
	@bash scripts/seed-db.sh 2>nul || powershell -ExecutionPolicy Bypass -File scripts/seed-db.ps1
	@echo "✅ Seed completado"

# ── Frontend ────────────────────────────────────────────

fe-install:
	@echo "📦 Installing frontend dependencies..."
	cd frontend && pnpm install
	@echo "✅ Frontend dependencies installed"

fe-dev:
	@echo "🔥 Starting frontend dev server..."
	cd frontend && pnpm dev

fe-build:
	@echo "📦 Building frontend..."
	cd frontend && pnpm build
	@echo "✅ Frontend built at frontend/dist/"

fe-lint:
	@echo "🔍 Linting frontend..."
	cd frontend && pnpm lint
	@echo "✅ Frontend lint OK"

fe-check:
	@echo "🔍 Type-checking frontend..."
	cd frontend && pnpm type-check
	@echo "✅ Frontend type-check OK"

fe-test:
	@echo "🧪 Running frontend unit tests..."
	cd frontend && pnpm test:run
	@echo "✅ Frontend unit tests passed"

fe-test-cover:
	@echo "🧪 Running frontend tests with coverage..."
	cd frontend && pnpm test:cover
	@echo "✅ Frontend coverage report"
