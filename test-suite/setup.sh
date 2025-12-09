#!/bin/bash

# Setup script for test suite dependencies
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}===================================${NC}"
echo -e "${BLUE}   Test Suite Setup${NC}"
echo -e "${BLUE}===================================${NC}"
echo ""

# Check Go installation
echo -e "${YELLOW}Checking Go installation...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    echo "Please install Go 1.21 or later from https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✅ Go is installed: $GO_VERSION${NC}"

# Check Docker installation
echo -e "${YELLOW}Checking Docker installation...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed${NC}"
    echo "Please install Docker Desktop from https://www.docker.com/products/docker-desktop"
    exit 1
fi

DOCKER_VERSION=$(docker --version)
echo -e "${GREEN}✅ Docker is installed: $DOCKER_VERSION${NC}"

# Check docker-compose
echo -e "${YELLOW}Checking Docker Compose...${NC}"
if command -v docker-compose &> /dev/null; then
    COMPOSE_VERSION=$(docker-compose --version)
    echo -e "${GREEN}✅ Docker Compose is installed: $COMPOSE_VERSION${NC}"
elif docker compose version &> /dev/null; then
    COMPOSE_VERSION=$(docker compose version)
    echo -e "${GREEN}✅ Docker Compose is available: $COMPOSE_VERSION${NC}"
else
    echo -e "${RED}❌ Docker Compose is not available${NC}"
    exit 1
fi

# Install test dependencies
echo ""
echo -e "${YELLOW}Installing Go test dependencies...${NC}"
cd "$PROJECT_ROOT"

# Install testify
echo "   Installing testify/suite..."
go get github.com/stretchr/testify/suite
go get github.com/stretchr/testify/assert

# Download all dependencies
echo "   Downloading all Go modules..."
go mod download

# Verify dependencies
echo "   Verifying Go modules..."
go mod verify

echo -e "${GREEN}✅ All dependencies installed${NC}"

# Make test scripts executable
echo ""
echo -e "${YELLOW}Making test scripts executable...${NC}"
chmod +x "$SCRIPT_DIR/run_tests.sh"
chmod +x "$SCRIPT_DIR/setup.sh"
echo -e "${GREEN}✅ Scripts are now executable${NC}"

# Summary
echo ""
echo -e "${GREEN}===================================${NC}"
echo -e "${GREEN}   Setup Complete!${NC}"
echo -e "${GREEN}===================================${NC}"
echo ""
echo -e "${BLUE}You can now run tests with:${NC}"
echo ""
echo "  # Run with isolated test database (recommended)"
echo "  cd $SCRIPT_DIR"
echo "  ./run_tests.sh isolated"
echo ""
echo "  # Or use make commands from project root:"
echo "  cd $PROJECT_ROOT"
echo "  make test-isolated"
echo ""
echo -e "${YELLOW}Note:${NC} The isolated mode will:"
echo "  • Start a fresh PostgreSQL test database on port 5433"
echo "  • Start test API server on port 8080"
echo "  • Start test UI server on port 8081"
echo "  • Run all 30 test scenarios"
echo "  • Clean up after tests (optional)"
echo ""
echo -e "${GREEN}Happy testing! 🧪${NC}"
