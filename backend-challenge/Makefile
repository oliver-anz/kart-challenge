.PHONY: build run test test-coverage clean generate

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
