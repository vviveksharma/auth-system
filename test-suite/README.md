# Integration Test Suite

Comprehensive integration tests for the Auth System with **completely isolated test infrastructure**.

## 🎯 Overview

This test suite provides **30 comprehensive integration tests** with dedicated test infrastructure:

### Test Environment
- **Dedicated CockroachDB** (localhost:5433) - Fresh test database, isolated from production
- **Dedicated Redis** (localhost:6380) - Isolated cache instance
- **Test API Server** (port 8080) - For application key authentication
- **Test UI Server** (port 8081) - Tenant admin portal
- **Pre-seeded dummy data** - Application key: `test_app_key_12345678901234567890`

### Test Coverage
- **API Server**: User authentication, profile management, role management, permissions
- **UI Server**: Tenant management, admin dashboard, message approval, role assignments

## 🚀 Quick Start

### First Time Setup
```bash
# Install dependencies and build images
make test-setup

# Or manually
cd test-suite
./setup.sh
```

### Run All Tests
```bash
# From project root
make test-isolated

# Or from test-suite folder
./run_tests.sh
```

## 📋 Test Scenarios

The suite includes **30 tests** organized in 15 scenarios:

### UI Server Tests (Tenant Management)
1. ✅ Tenant Registration & Login (Tests 1-2)
2. ✅ Application Token Management (Tests 3-5)
3. ✅ Message Management (Tests 18-20)
4. ✅ Dashboard & Analytics (Tests 21-22)
5. ✅ Tenant Role Management (Tests 23-24)

### API Server Tests (User Management)
6. ✅ User Registration & Authentication (Tests 6-7)
7. ✅ Profile Management (Tests 8-9)
8. ✅ Role Management (Tests 10-15)
9. ✅ Role Change Requests (Tests 16-17)
10. ✅ Token Refresh (Test 25)
11. ✅ Password Reset (Test 26)
12. ✅ Cleanup (Test 27)
13. ✅ Logout & Validation (Tests 28-29)
14. ✅ Health Checks (Test 30)

See [TEST_SCENARIOS.md](./TEST_SCENARIOS.md) for detailed documentation.

## 🛠️ Available Commands

### Main Commands
```bash
make test-isolated    # Run full test suite with isolated environment
make test-quick       # Run tests with existing containers (faster)
make test-setup       # First-time setup
```

### Container Management
```bash
make test-start       # Start test containers
make test-stop        # Stop test containers
make test-clean       # Remove containers and volumes
make test-status      # Show container status
make test-logs        # Follow container logs
```

### From test-suite folder
```bash
./setup.sh           # Setup environment
./run_tests.sh       # Run all tests
make test            # Full test suite
make test-quick      # Quick run
make coverage        # Generate coverage report
```

## 🗂️ Test Environment

### Isolated Infrastructure
The test suite uses completely isolated infrastructure:
- **Test Database**: `auth_test_db` on port **5433** (not 5432)
- **Test Redis**: Port **6380** (not 6379)
- **Test Servers**: Separate Docker containers with test configuration

### No Production Impact
- ✅ Tests run on isolated database
- ✅ Unique test data with timestamps
- ✅ Automatic cleanup option
- ✅ Production data remains untouched

## 📊 Test Output

```bash
🧪 Starting Integration Test Suite
====================================
⏳ Waiting for API Server at http://localhost:8080...
✅ API Server is ready!
⏳ Waiting for UI Server at http://localhost:8081...
✅ UI Server is ready!

✅ Both servers are ready!

📝 Test 1: Tenant Registration
   ✅ Tenant registered: tenant_1732262400@test.com

🔐 Test 2: Tenant Login
   ✅ Tenant logged in successfully

🔑 Test 3: Create Application Token
   ✅ Application token created: abc123...

... (30 tests total)

💚 Test 30: Final Health Checks
   ✅ API Server healthy
   ✅ UI Server healthy

🏁 Test Suite Completed!
=======================
```

## 🔍 Debugging Failed Tests

### View Logs
```bash
# All logs
make test-logs

# Specific service
docker logs auth-test-api-server
docker logs auth-test-ui-server
docker logs auth-test-db
```

### Run Specific Test
```bash
cd test-suite
go test -v -run TestIntegrationSuite/Test07_API_UserLogin
```

### Check Container Status
```bash
make test-status
# or
cd test-suite && docker-compose -f docker-compose.test.yml ps
```

### Interactive Debugging
```bash
# Open shell in API server
docker exec -it auth-test-api-server sh

# Open shell in database
docker exec -it auth-test-db psql -U authuser -d auth_test_db
```

## 📦 File Structure

```
test-suite/
├── integration_test.go          # Main test file (30 tests)
├── docker-compose.test.yml      # Isolated test environment
├── setup.sh                     # Environment setup script
├── run_tests.sh                 # Test runner script
├── Makefile                     # Test commands
├── TEST_SCENARIOS.md            # Detailed test documentation
└── README.md                    # This file
```

## 🧪 Test Data Management

### Unique Data Generation
Each test run creates unique data:
```go
timestamp := time.Now().Unix()
tenantEmail := fmt.Sprintf("tenant_%d@test.com", timestamp)
userEmail := fmt.Sprintf("user_%d@test.com", timestamp)
```

### Database Cleanup
```bash
# Remove all test data
make test-clean

# Or manually
cd test-suite
docker-compose -f docker-compose.test.yml down -v
```

## ✅ Success Criteria

All tests should pass with:
- ✅ Correct HTTP status codes (200, 401, 404, etc.)
- ✅ Valid response data structures
- ✅ Token generation and validation
- ✅ Database operations successful
- ✅ Proper error handling
- ✅ Both servers healthy throughout

## 🔧 Requirements

- **Go**: 1.21 or higher
- **Docker**: Latest version
- **Docker Compose**: v2.0+
- **Dependencies**:
  - `github.com/stretchr/testify/suite`
  - `github.com/stretchr/testify/assert`

## 📈 Coverage

Generate coverage report:
```bash
cd test-suite
make coverage
# Opens coverage.html in browser
```

## 🤝 Contributing

When adding new tests:
1. Add test method with descriptive name
2. Document in TEST_SCENARIOS.md
3. Update test count in README
4. Ensure test independence (no cross-test dependencies)
5. Use unique test data (timestamps)

## 📝 Notes

- Tests run **sequentially** to maintain dependencies
- Each test is **independent** with unique data
- Test environment is **ephemeral** and can be recreated
- Logs are available for debugging
- Test data is **NOT** cleaned up automatically (for debugging)

## 🆘 Troubleshooting

### Tests Timeout
- Increase timeout: `go test -v -timeout 10m`
- Check container health: `make test-status`

### Connection Refused
- Ensure containers are running: `make test-start`
- Wait for services: containers need ~30 seconds to be ready

### Port Conflicts
- Check if ports are in use: `lsof -i :8080,8081,5433,6380`
- Stop conflicting services

### Database Issues
- Reset database: `make test-clean && make test-start`
- Check logs: `docker logs auth-test-db`

## 📚 Related Documentation

- [TEST_SCENARIOS.md](./TEST_SCENARIOS.md) - Detailed test scenarios
- [../DEPLOYMENT.md](../DEPLOYMENT.md) - Production deployment guide
- [../SCALABILITY_ROADMAP.md](../SCALABILITY_ROADMAP.md) - Scaling recommendations

---

**Last Updated**: November 2024  
**Test Count**: 30 scenarios  
**Coverage**: Both API and UI servers
