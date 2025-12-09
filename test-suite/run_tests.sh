#!/bin/bash

# Test runner script with isolated database support
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TEST_MODE="${1:-isolated}"
COMPOSE_FILE="$PROJECT_ROOT/test-suite/docker-compose.test.yml"

echo -e "${BLUE}===================================${NC}"
echo -e "${BLUE}   Auth System Test Runner${NC}"
echo -e "${BLUE}===================================${NC}"
echo ""

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo -e "${RED}❌ Docker is not running. Please start Docker Desktop.${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Docker is running${NC}"
}

# Function to check if docker-compose is available
check_docker_compose() {
    if command -v docker-compose > /dev/null 2>&1; then
        DOCKER_COMPOSE="docker-compose"
    elif docker compose version > /dev/null 2>&1; then
        DOCKER_COMPOSE="docker compose"
    else
        echo -e "${RED}❌ docker-compose or 'docker compose' not found${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Docker Compose is available${NC}"
}

# Function to clean up test environment
cleanup_test_env() {
    echo -e "${YELLOW}🧹 Cleaning up test environment...${NC}"
    cd "$PROJECT_ROOT"
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" down -v --remove-orphans 2>/dev/null || true
    echo -e "${GREEN}✅ Cleanup complete${NC}"
}

# Function to start isolated test environment
start_isolated_env() {
    echo -e "${BLUE}🚀 Starting isolated test environment...${NC}"
    cd "$PROJECT_ROOT"
    
    # Stop any existing test containers and remove volumes
    echo -e "${YELLOW}🗑️  Removing old containers and volumes...${NC}"
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" down -v --remove-orphans 2>/dev/null || true
    
    # Remove any orphaned volumes
    docker volume rm test_db_data 2>/dev/null || true
    
    # Start test environment
    echo -e "${YELLOW}📦 Building and starting containers...${NC}"
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" up -d --build
    
    # Wait for database to initialize
    echo -e "${YELLOW}⏳ Waiting for database initialization...${NC}"
    sleep 5
    
    echo -e "${YELLOW}⏳ Waiting for services to be ready (this may take 30-60 seconds)...${NC}"
    
    # Wait for API server health check
    echo -n "   Waiting for API server (localhost:8080)... "
    for i in {1..60}; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            echo -e "${GREEN}Ready!${NC}"
            break
        fi
        sleep 1
        if [ $i -eq 60 ]; then
            echo -e "${RED}Timeout!${NC}"
            echo -e "${RED}❌ API server did not start in time${NC}"
            echo "Logs:"
            $DOCKER_COMPOSE -f "$COMPOSE_FILE" logs test-api-server
            exit 1
        fi
    done
    
    # Wait for UI server health check
    echo -n "   Waiting for UI server (localhost:8081)... "
    for i in {1..60}; do
        if curl -s http://localhost:8081/health > /dev/null 2>&1; then
            echo -e "${GREEN}Ready!${NC}"
            break
        fi
        sleep 1
        if [ $i -eq 60 ]; then
            echo -e "${RED}Timeout!${NC}"
            echo -e "${RED}❌ UI server did not start in time${NC}"
            echo "Logs:"
            $DOCKER_COMPOSE -f "$COMPOSE_FILE" logs test-ui-server
            exit 1
        fi
    done
    
    # Seed test data (dummy application key and test tenant)
    echo -e "${YELLOW}📦 Seeding test data...${NC}"
    if [ -f "$SCRIPT_DIR/seed_test_data.sql" ]; then
        # Execute SQL via docker exec
        docker exec auth-test-db cockroach sql --insecure --database=defaultdb < "$SCRIPT_DIR/seed_test_data.sql" 2>/dev/null || echo -e "${YELLOW}   Note: Some seed data may already exist${NC}"
        echo -e "${GREEN}   ✅ Test data seeded${NC}"
    else
        echo -e "${YELLOW}   ⚠️  seed_test_data.sql not found, skipping${NC}"
    fi
    
    echo -e "${GREEN}✅ Test environment is ready!${NC}"
    echo ""
    echo -e "${BLUE}Test Services:${NC}"
    echo "   • API Server:     http://localhost:8080"
    echo "   • UI Server:      http://localhost:8081"
    echo "   • Test Database:  localhost:5433"
    echo "   • Test Redis:     localhost:6380"
    echo "   • Test Mailpit:   http://localhost:8026"
    echo ""
}

# Function to run tests
run_tests() {
    echo -e "${BLUE}🧪 Running integration tests...${NC}"
    cd "$SCRIPT_DIR"
    
    # Set test environment variables
    export API_BASE_URL="http://localhost:8080"
    export UI_BASE_URL="http://localhost:8081"
    
    # Run tests with verbose output
    if go test -v -timeout 300s; then
        echo ""
        echo -e "${GREEN}✅ All tests passed!${NC}"
        return 0
    else
        echo ""
        echo -e "${RED}❌ Some tests failed${NC}"
        return 1
    fi
}

# Function to show logs
show_logs() {
    echo -e "${BLUE}📋 Test Environment Logs:${NC}"
    echo ""
    cd "$PROJECT_ROOT"
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" logs --tail=100
}

# Main execution flow
main() {
    case "$TEST_MODE" in
        isolated)
            echo -e "${YELLOW}Mode: Isolated Test Environment${NC}"
            echo -e "${YELLOW}This will start a fresh database and servers${NC}"
            echo ""
            
            check_docker
            check_docker_compose
            
            # Start isolated environment
            start_isolated_env
            
            # Run tests
            if run_tests; then
                TEST_EXIT_CODE=0
            else
                TEST_EXIT_CODE=1
            fi
            
            # Ask if user wants to keep environment running
            echo ""
            echo -e "${YELLOW}Keep test environment running? (y/n)${NC}"
            read -r -t 10 KEEP_RUNNING || KEEP_RUNNING="n"
            
            if [[ "$KEEP_RUNNING" =~ ^[Yy]$ ]]; then
                echo -e "${GREEN}Test environment is still running. To stop:${NC}"
                echo "  cd $PROJECT_ROOT"
                echo "  $DOCKER_COMPOSE -f test-suite/docker-compose.test.yml down -v"
            fi
            
            exit $TEST_EXIT_CODE
            ;;
            
        quick)
            echo -e "${YELLOW}Mode: Quick Test (uses existing servers)${NC}"
            echo -e "${YELLOW}Assumes servers are already running on 8080 and 8081${NC}"
            echo ""
            
            # Check if servers are running
            if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
                echo -e "${RED}❌ API server not running on port 8080${NC}"
                exit 1
            fi
            
            if ! curl -s http://localhost:8081/health > /dev/null 2>&1; then
                echo -e "${RED}❌ UI server not running on port 8081${NC}"
                exit 1
            fi
            
            echo -e "${GREEN}✅ Servers are running${NC}"
            echo ""
            
            run_tests
            ;;
            
        logs)
            show_logs
            ;;
            
        cleanup)
            check_docker_compose
            cleanup_test_env
            ;;
            
        *)
            echo -e "${RED}Unknown mode: $TEST_MODE${NC}"
            echo ""
            echo "Usage: $0 [MODE]"
            echo ""
            echo "Modes:"
            echo "  isolated  - Start isolated test environment with fresh DB (default)"
            echo "  quick     - Run tests against existing servers (faster)"
            echo "  logs      - Show test environment logs"
            echo "  cleanup   - Clean up test environment"
            echo ""
            echo "Examples:"
            echo "  $0 isolated     # Full isolated test run"
            echo "  $0 quick        # Quick test with existing servers"
            echo "  $0 cleanup      # Clean up test containers"
            echo "  $0 logs         # View logs"
            exit 1
            ;;
    esac
}

# Trap Ctrl+C and cleanup
trap cleanup_test_env INT

main
