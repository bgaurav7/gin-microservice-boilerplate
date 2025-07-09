.PHONY: run build test lint clean

# Default target
all: build

# Run the application with live reload
run:
	@echo "Running application with live reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air is not installed. Installing..."; \
		go install github.com/air-verse/air@latest; \
		$(GOPATH)/bin/air || ~/go/bin/air; \
	fi

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/server ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint is not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin tmp
