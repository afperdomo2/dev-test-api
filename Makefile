.PHONY: install dev build swagger clean run

install:
	@echo "📦 Installing dependencies..."
	go mod tidy
	@echo "✅ Dependencies installed"

dev:
	@echo "🔥 Starting live reload server..."
	go tool air

build:
	@echo "📦 Building..."
	go build -o ./tmp/main .
	@echo "✅ Binary built at ./tmp/main"

run:
	@echo "🚀 Running server..."
	go run main.go

swagger:
	@echo "📖 Generating Swagger docs..."
	swag init -g main.go
	@echo "✅ Docs generated in docs/"

clean:
	@echo "🧹 Cleaning..."
	rm -rf tmp/ docs/
	@echo "✅ Cleaned"
