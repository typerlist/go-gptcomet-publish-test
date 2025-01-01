# Show available commands
help:
    @just --list

# Generate changelog
changelog:
    git cliff -l --prepend CHANGELOG.md

# Display all available commands
default:
    @just --list

# Run all tests
test:
    go test ./... -v

# Generate coverage report
coverage:
    go test ./... -coverprofile=coverage.out
    go tool cover -func=coverage.out

# Generate HTML coverage report
coverage-html: coverage
    go tool cover -html=coverage.out -o coverage.html

# Run code linting
lint:
    golangci-lint run ./...

# Format code
fmt:
    go fmt ./...
    gofmt -s -w .

# Update dependencies
tidy:
    go mod tidy

# Build project
build:
    go build -v ./...

# Clean build and test artifacts
clean:
    rm -f coverage.out coverage.html
    go clean

# Run all quality checks (test + lint)
check: test lint

# Install development dependencies
dev-deps:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run local development server
run:
    go run main.go
