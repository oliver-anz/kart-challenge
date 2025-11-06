.PHONY: help build run init test test-coverage clean generate

.DEFAULT_GOAL := help

# Help target - shows all available commands
help: ## Show this help message
	@echo "Backend Challenge - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# Initialize database
init: ## Initialize SQLite database with schema and seed data
	@echo "Initializing database..."
	@mkdir -p data
	@sqlite3 data/store.db < data/init.sql
	@echo "Database initialized: data/store.db"

# Build the application
build: ## Build the backend server
	@echo "Building backend-challenge..."
	@go build -o backend-challenge .
	@echo "Build complete!"

# Run the server
run: build ## Build and run the server
	@./backend-challenge

# Run all tests
test: ## Run all tests
	@echo "Running all tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total

# Clean build artifacts and test files
clean: ## Clean build artifacts and test databases
	@echo "Cleaning..."
	@rm -f backend-challenge
	@rm -f coverage.out coverage.html
	@rm -f test*.db
	@rm -f api/test*.db
	@rm -f db/test*.db
	@echo "Clean complete!"

# Generate mocks
generate: ## Generate mocks for testing
	@echo "Generating mocks..."
	@go install go.uber.org/mock/mockgen@latest
	@go generate ./...
	@echo "Mocks generated!"
