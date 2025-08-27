# Makefile for MCP Directory Walker Server

# Variables
BINARY_NAME=directory-walker
TEST_STDIO_BINARY=test/test-stdio
TEST_HTTP_BINARY=test/test-http
GO_FILES=$(shell find . -name "*.go" -not -path "./test/*" -not -name "*_test.go")
TEST_FILES=$(shell find . -name "*_test.go" -not -path "./test/*")

# Default target
.PHONY: all
all: build

# Build the main binary
.PHONY: build
build: $(BINARY_NAME)

$(BINARY_NAME): $(GO_FILES) go.mod go.sum
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

# Build test tools
.PHONY: build-test-tools
build-test-tools: $(TEST_STDIO_BINARY) $(TEST_HTTP_BINARY)

$(TEST_STDIO_BINARY): test/test-stdio.go go.mod go.sum
	@echo "Building test-stdio..."
	cd test && go build -o test-stdio test-stdio.go

$(TEST_HTTP_BINARY): test/test-http.go go.mod go.sum
	@echo "Building test-http..."
	cd test && go build -o test-http test-http.go

# Run the server with current directory as root (HTTP mode)
.PHONY: run
run: build
	@echo "Starting MCP Directory Walker Server (HTTP) for current directory..."
	./$(BINARY_NAME) .

# Run the server in stdio mode
.PHONY: run-stdio
run-stdio: build
	@echo "Starting MCP Directory Walker Server (stdio) for current directory..."
	./$(BINARY_NAME) -s .

# Run unit tests
.PHONY: test
test:
	@echo "Running unit tests..."
	go test -v ./...

# Run integration tests (requires server to be built)
.PHONY: test-integration
test-integration: build build-test-tools
	@echo "Running integration tests..."
	@echo "Testing stdio transport..."
	$(TEST_STDIO_BINARY) ./$(BINARY_NAME) .
	@echo ""
	@echo "For HTTP tests, start the server first with 'make run' in another terminal,"
	@echo "then run: $(TEST_HTTP_BINARY)"

# Test stdio transport specifically
.PHONY: test-stdio
test-stdio: build build-test-tools
	@echo "Testing stdio transport..."
	$(TEST_STDIO_BINARY) ./$(BINARY_NAME) .

# Test HTTP transport (assumes server is already running)
.PHONY: test-http
test-http: build-test-tools
	@echo "Testing HTTP transport (server should be running on localhost:5001)..."
	$(TEST_HTTP_BINARY)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f $(TEST_STDIO_BINARY)
	rm -f $(TEST_HTTP_BINARY)
	go clean

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build              - Build the directory-walker binary"
	@echo "  build-test-tools   - Build test tools (test-stdio and test-http)"
	@echo "  run                - Build and run server in HTTP mode for current directory"
	@echo "  run-stdio          - Build and run server in stdio mode for current directory"
	@echo "  test               - Run unit tests"
	@echo "  test-integration   - Run integration tests (builds everything first)"
	@echo "  test-stdio         - Test stdio transport specifically"
	@echo "  test-http          - Test HTTP transport (assumes server is running)"
	@echo "  clean              - Remove build artifacts"
	@echo "  deps               - Install and tidy dependencies"
	@echo "  help               - Show this help message"

# Development targets
.PHONY: dev-http
dev-http: build
	@echo "Starting development HTTP server with verbose output..."
	./$(BINARY_NAME) . 2>&1

.PHONY: dev-stdio  
dev-stdio: build
	@echo "Starting development stdio server..."
	./$(BINARY_NAME) -s .

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint to be installed)
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run ./...
