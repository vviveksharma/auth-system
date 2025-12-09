# Separate Server Architecture

## Overview

This system now runs as **two independent Docker containers**:

1. **API Server (Port 8080)** - For application key authenticated clients
2. **UI Server (Port 8081)** - For tenant admin portal

Both servers share the same database and Redis cache but run completely independently.

---

## Architecture Benefits

### 🎯 **Independent Scaling**
- Scale API server separately based on application traffic
- Scale UI server based on admin portal usage
- Different resource allocation per service

### 🔒 **Better Isolation**
- API server crash won't affect UI server
- Independent deployments and rollbacks
- Separate monitoring and logging

### 📊 **Clear Separation of Concerns**
- API Server: Application key authentication, user/role management
- UI Server: Tenant authentication, admin operations

---

## Running the System

### **Option 1: Run Both Servers (Docker Compose)**
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f api-server  # API server logs
docker-compose logs -f ui-server   # UI server logs

# Stop all services
docker-compose down
```

**Access:**
- API Server: http://localhost:8080
- UI Server: http://localhost:8081
- Mailpit UI: http://localhost:8025
- Database Admin: http://localhost:26257

---

### **Option 2: Run Individual Servers**

**API Server Only:**
```bash
docker-compose up -d db redis mailpit api-server
```

**UI Server Only:**
```bash
docker-compose up -d db redis mailpit ui-server
```

**Stop Specific Server:**
```bash
docker-compose stop api-server
docker-compose stop ui-server
```

---

### **Option 3: Local Development (Both Servers)**
```bash
# Run both servers locally (default behavior)
go run main.go

# Or with explicit mode
SERVER_MODE=BOTH go run main.go
```

**Option 4: Local Development (Single Server)**
```bash
# Run only API server
SERVER_MODE=API go run main.go

# Run only UI server
SERVER_MODE=UI go run main.go
```

---

## Docker Build Process

### **Building Individual Images:**

```bash
# Build API server image
docker build -f Dockerfile.api -t auth-api-server .

# Build UI server image
docker build -f Dockerfile.ui -t auth-ui-server .
```

### **Building with Docker Compose:**
```bash
# Build both images
docker-compose build

# Build specific service
docker-compose build api-server
docker-compose build ui-server
```

---

## Configuration

### **Environment Variables**

Both servers share the same configuration but run independently:

```env
# Database
DB_HOST=db
DB_PORT=26257
DB_USER=root
DB_PASSWORD=
DB_NAME=auth_system

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# SMTP
SMTP_HOST=mailpit
SMTP_PORT=1025

# Server Mode (API, UI, or BOTH)
SERVER_MODE=API  # or UI or BOTH
```

### **Docker Compose Override (Optional)**

Create `docker-compose.override.yml` for local customization:

```yaml
version: "3.8"

services:
  api-server:
    environment:
      - DEBUG=true
      - LOG_LEVEL=debug
    ports:
      - "9090:8080"  # Custom port mapping
  
  ui-server:
    environment:
      - DEBUG=true
      - LOG_LEVEL=debug
```

---

## Health Checks

Both servers expose health check endpoints:

```bash
# API Server health
curl http://localhost:8080/health

# UI Server health
curl http://localhost:8081/health
```

**Health Check Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-11-22T...",
  "server": "api",  // or "ui"
  "database": "healthy",
  "cache": "healthy"
}
```

---

## Scaling Strategies

### **Docker Compose Scaling:**
```bash
# Scale API server to 3 instances
docker-compose up -d --scale api-server=3

# Scale UI server to 2 instances
docker-compose up -d --scale ui-server=2
```

**Note:** Remove port mappings from docker-compose.yml and use a reverse proxy (nginx) for load balancing when scaling.

### **Production Scaling with Nginx:**

```nginx
# nginx.conf
upstream api_backend {
    least_conn;
    server api-server-1:8080;
    server api-server-2:8080;
    server api-server-3:8080;
}

upstream ui_backend {
    least_conn;
    server ui-server-1:8081;
    server ui-server-2:8081;
}

server {
    listen 80;
    
    location /api/ {
        proxy_pass http://api_backend/;
    }
    
    location /admin/ {
        proxy_pass http://ui_backend/;
    }
}
```

---

## Monitoring

### **Docker Stats:**
```bash
# Monitor resource usage
docker stats auth-api-server auth-ui-server
```

### **Container Logs:**
```bash
# Follow API server logs
docker logs -f auth-api-server

# Follow UI server logs
docker logs -f auth-ui-server

# View last 100 lines
docker logs --tail 100 auth-api-server
```

### **Docker Compose Logs:**
```bash
# All services
docker-compose logs -f

# Specific services
docker-compose logs -f api-server ui-server
```

---

## Troubleshooting

### **Container Won't Start:**
```bash
# Check container status
docker-compose ps

# View full logs
docker-compose logs api-server
docker-compose logs ui-server

# Restart specific service
docker-compose restart api-server
```

### **Database Connection Issues:**
```bash
# Check if database is healthy
docker-compose ps db

# Test database connection from container
docker exec -it auth-api-server sh
# Inside container:
wget -O- http://db:26257/health?ready=1
```

### **Port Already in Use:**
```bash
# Find process using port
lsof -i :8080
lsof -i :8081

# Change port in docker-compose.yml
ports:
  - "9080:8080"  # Map to different host port
```

---

## Development Workflow

### **1. Local Development (Hot Reload)**
```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### **2. Test Changes in Docker**
```bash
# Rebuild and restart
docker-compose up -d --build api-server

# View logs
docker-compose logs -f api-server
```

### **3. Debug Inside Container**
```bash
# Access container shell
docker exec -it auth-api-server sh

# Check environment variables
env | grep DB_HOST

# Test network connectivity
ping db
ping redis
```

---

## Migration from Single Container

**Old Way (Single Container):**
```bash
docker run -p 8080:8080 -p 8081:8081 auth-system
```

**New Way (Separate Containers):**
```bash
docker-compose up -d api-server ui-server
```

**No Breaking Changes:**
- Same APIs and endpoints
- Same authentication flow
- Same database schema
- Existing clients continue to work

---

## Performance Optimization

### **Resource Limits (Production):**

```yaml
services:
  api-server:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
  
  ui-server:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
        reservations:
          cpus: '0.25'
          memory: 128M
```

### **Connection Pooling:**

Both servers share optimized connection pools:
- Database: Max 100 connections (50 per server recommended)
- Redis: Max 50 connections (25 per server recommended)

---

## Next Steps

✅ **Done:** Separate server containers  
⏭️ **Next:** Implement features from SCALABILITY_ROADMAP.md:
  1. Connection pooling optimization
  2. Redis caching strategy
  3. Rate limiting per tenant
  4. Background job processing
  5. Horizontal scaling with load balancer

See `SCALABILITY_ROADMAP.md` for detailed implementation guide.
