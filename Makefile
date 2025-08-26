# Makefile for filez-mcp directory walker

# Variables
BINARY_NAME=directory-walker
TEST_HTTP_BINARY=test/test-http
TEST_STDIO_BINARY=test/test-stdio
GO_FLAGS=-ldflags="-s -w"

# Default target
.PHONY: all
all: build

# Build the main binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(GO_FLAGS) -o $(BINARY_NAME) main.go
	@echo "✓ Build completed: $(BINARY_NAME)"

# Build test tools
.PHONY: build-tests
build-tests:
	@echo "Building test tools..."
	go build -o $(TEST_HTTP_BINARY) test/test-http.go
	go build -o $(TEST_STDIO_BINARY) test/test-stdio.go
	@echo "✓ Test tools built: $(TEST_HTTP_BINARY), $(TEST_STDIO_BINARY)"

# Run the server with current directory as root
.PHONY: run
run: build
	@echo "Starting $(BINARY_NAME) with current directory as root..."
	@echo "Press Ctrl+C to stop the server"
	@echo "Server will be available at http://localhost:5001/mcp"
	./$(BINARY_NAME) .

# Run the server in stdio mode
.PHONY: run-stdio
run-stdio: build
	@echo "Starting $(BINARY_NAME) in stdio mode with current directory as root..."
	@echo "Press Ctrl+C to stop the server"
	./$(BINARY_NAME) -s .

# Run unit tests
.PHONY: test
test:
	@echo "Running unit tests..."
	go test -v -race -cover
	@echo "✓ All tests passed"

# Run unit tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# Test HTTP transport
.PHONY: test-http
test-http: build build-tests
	@echo "Testing HTTP transport..."
	./$(TEST_HTTP_BINARY) ./$(BINARY_NAME) . 5002

# Test stdio transport
.PHONY: test-stdio
test-stdio: build build-tests
	@echo "Testing stdio transport..."
	./$(TEST_STDIO_BINARY) ./$(BINARY_NAME) .

# Run all integration tests
.PHONY: test-integration
test-integration: test-http test-stdio
	@echo "✓ All integration tests completed"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f $(TEST_HTTP_BINARY)
	rm -f $(TEST_STDIO_BINARY)
	rm -f coverage.out coverage.html
	@echo "✓ Clean completed"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Code formatted"

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	go vet .
	@echo "✓ Linting completed"

# Check code quality
.PHONY: check
check: fmt lint test
	@echo "✓ All checks passed"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download
	@echo "✓ Dependencies installed"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build           - Build the directory-walker binary"
	@echo "  build-tests     - Build test tools"
	@echo "  run             - Run server in HTTP mode with current directory"
	@echo "  run-stdio       - Run server in stdio mode with current directory"
	@echo "  test            - Run unit tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-http       - Test HTTP transport"
	@echo "  test-stdio      - Test stdio transport"
	@echo "  test-integration- Run all integration tests"
	@echo "  clean           - Remove build artifacts"
	@echo "  fmt             - Format code"
	@echo "  lint            - Run linter"
	@echo "  check           - Run fmt, lint, and test"
	@echo "  deps            - Install/update dependencies"
	@echo "  help            - Show this help message"
