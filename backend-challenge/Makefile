.PHONY: test test-unit test-integration test-coverage test-verbose build clean run help

# Default target
.DEFAULT_GOAL := help

# Build the application
build: ## Build the backend server
	@echo "Building backend-challenge..."
	@go build -o backend-challenge .
	@echo "Build complete!"

# Run all tests
test: ## Run all tests
	@echo "Running all tests..."
	@go test -v ./...

# Run unit tests only (exclude integration tests)
test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test -v -short ./api/... ./db/... ./models/...

# Run integration tests
test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@go test -v -run Integration .

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

# Run the server
run: build ## Build and run the server
	@./backend-challenge

# Initialize the database
db-init: ## Initialize the database from init.sql
	@echo "Initializing database..."
	@sqlite3 data/store.db < init.sql
	@echo "Database initialized!"

# Check if database exists
db-check: ## Check if database is populated
	@echo "Checking database..."
	@sqlite3 data/store.db "SELECT COUNT(*) as products FROM products; SELECT COUNT(*) as coupons FROM valid_coupons;"

# Run linter
lint: ## Run golangci-lint
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

# Vet code
vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Run all checks (fmt, vet, test)
check: fmt vet test ## Run fmt, vet, and tests

# Install dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Help target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
