# Default target
default: build

# Build the application
.PHONY: build

build:
	go build ./...

# Test target to run all tests in the project
.PHONY: test

test: build
	go test ./... -v

# Lint target to run the linter
.PHONY: lint

lint:
	golangci-lint run

# Run linter with fixes
.PHONY: lint-fix

lint-fix:
	golangci-lint run --fix

# Check Go code formatting without modifying files
.PHONY: fmt-check

fmt-check:
	@echo "Checking Go formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "The following files need formatting:"; \
		gofmt -l .; \
		exit 1; \
	fi
	@if [ -n "$$(goimports -l .)" ]; then \
		echo "The following files need import formatting:"; \
		goimports -l .; \
		exit 1; \
	fi
	@echo "All files are properly formatted!"

# Format Go code
.PHONY: fmt

fmt:
	gofmt -w .
	goimports -w .

# Install the application
.PHONY: install

install: build
	go install

# Check for Go modules updates
.PHONY: update

update:
	go get -u ./...
	go mod tidy


# Ensure pre-commit hook is executable
.PHONY: config-pre-commit

config-pre-commit:
	@echo "Setting up pre-commit hook..."
	git config --local core.hooksPath .githooks/

# Pre-commit checks (format, lint, test)
.PHONY: pre-commit

pre-commit: fmt lint test
	@echo "Pre-commit checks passed!"