# Test Suite Makefile
# Convenience commands for running tests

.PHONY: help test-isolated test-quick test-api test-ui test-clean test-logs setup

# Default target
help:
	@echo "===================================="
	@echo "   Auth System Test Commands"
	@echo "===================================="
	@echo ""
	@echo "Available commands:"
	@echo "  make setup          - Install test dependencies"
	@echo "  make test-isolated  - Run tests with isolated test database (recommended)"
	@echo "  make test-quick     - Run tests against existing servers (fast)"
	@echo "  make test-api       - Run API server tests only"
	@echo "  make test-ui        - Run UI server tests only"
	@echo "  make test-logs      - Show test environment logs"
	@echo "  make test-clean     - Clean up test environment"
	@echo ""
	@echo "Examples:"
	@echo "  make setup          # First time setup"
	@echo "  make test-isolated  # Full test with isolated DB"
	@echo "  make test-quick     # Quick test (servers must be running)"
	@echo ""

# Install dependencies
setup:
	@echo "🔧 Setting up test dependencies..."
	@./setup.sh

# Run tests with isolated test database
test-isolated:
	@echo "🧪 Running tests with isolated test database..."
	@./run_tests.sh isolated -v

# Run tests against existing servers (faster)
test-quick:
	@echo "⚡ Running quick tests (using existing servers)..."
	@./run_tests.sh quick

# Run API server tests only
test-api:
	@echo "🧪 Running API server tests only..."
	@cd .. && go test ./test-suite -v -run "Test.*_API_.*" -timeout 300s

# Run UI server tests only
test-ui:
	@echo "🧪 Running UI server tests only..."
	@cd .. && go test ./test-suite -v -run "Test.*_UI_.*" -timeout 300s

# Show test environment logs
test-logs:
	@echo "📋 Showing test environment logs..."
	@./run_tests.sh logs

# Clean up test environment
test-clean:
	@echo "🧹 Cleaning up test environment..."
	@./run_tests.sh cleanup
	@echo "✅ Test environment cleaned"

# Run a specific test scenario by number
test-scenario-%:
	@echo "🧪 Running test scenario $*..."
	@cd .. && go test ./test-suite -v -run "Test$(shell printf '%02d' $*)_.*" -timeout 60s

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	@cd .. && go test ./test-suite -v -coverprofile=test-suite/coverage.out -timeout 300s
	@cd .. && go tool cover -html=test-suite/coverage.out -o test-suite/coverage.html
	@echo "✅ Coverage report generated: test-suite/coverage.html"

# Watch mode - rerun tests on changes (requires entr)
test-watch:
	@echo "👀 Watching for changes..."
	@find .. -name "*.go" | entr -c make test-quick
