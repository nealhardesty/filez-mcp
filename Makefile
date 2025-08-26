# MCP Directory Walker Server Makefile

# Build target - compile the binary
build:
	@echo "Building directory-walker..."
	go build -o directory-walker main.go
	@echo "Build complete: ./directory-walker"

# Run target - build and run with current directory as root
run: build
	@echo "Starting directory-walker with current directory as root..."
	./directory-walker .

# Test target - run go test
test:
	@echo "Running tests..."
	go test -v ./...

# Clean target - remove built binaries
clean:
	@echo "Cleaning up..."
	rm -f directory-walker
	@echo "Clean complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download
	@echo "Dependencies installed"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

.PHONY: build run test clean deps fmt lint