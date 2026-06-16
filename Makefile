.PHONY: install dev build swagger clean run fmt vet check setup \
        fe-install fe-dev fe-build fe-lint fe-check

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
	go build -o ./tmp/main .
	@echo "✅ Binary built at ./tmp/main"

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

setup:
	@echo "🔗 Installing git hooks..."
	go tool lefthook install
	@echo "✅ Hooks installed"

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
