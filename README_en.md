# Webox

Next generation of Websoft9 platform - A comprehensive cloud application management solution.

## Overview

Webox is a modern cloud platform consisting of two main components:

- **api-service**: Core backend service providing RESTful APIs
- **websoft9-agent**: Client agent deployed on server nodes for task execution and monitoring

## Components

### Websoft9 Web Service

Core backend service built with Golang and Gin framework.

**Tech Stack:**

- Golang 1.24+
- Gin web framework
- GORM ORM
- SQLite database
- Redis cache
- InfluxDB
- JWT authentication
- gRPC

### Websoft9 Agent

Client agent deployed on server nodes for task execution and monitoring.

**Tech Stack:**

- Golang
- gRPC (grpc-go)
- Redis (go-redis)
- InfluxDB (influxdb-client-go)
- Docker Engine API

## Quick Start

### Prerequisites

- Go 1.24+
- Docker
- SQLite 3.0+
- Redis 6.0+
- InfluxDB 2.0+

### Development Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd webox
   ```

2. **Start Web Service**

   ```bash
   cd api-service
   go mod tidy
   go run main.go
   ```

3. **Start Agent**

   ```bash
   cd websoft9-agent
   go mod tidy
   go build -o websoft9-agent ./cmd/agent
   sudo ./websoft9-agent
   ```

### Docker Deployment

```bash
# Build and run web service
cd api-service
docker build -t api-service .
docker run -p 8080:8080 -p 9090:9090 api-service

# Build and run agent
cd websoft9-agent
docker build -t websoft9-agent .
docker run --privileged websoft9-agent
```

## Development Tools

- **API Testing**: [Apifox](https://apifox.com/)
- **CI/CD**: GitHub Actions
- **AI Coding**: [GitHub Copilot](https://github.com/features/copilot), [Claude Code](https://docs.anthropic.com/zh-CN/docs/claude-code/overview)

## Project Structure

```
webox/
├── api-service/    # Backend API service
│   ├── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── controller/
│   │   ├── service/
│   │   └── repository/
│   └── pkg/
└── websoft9-agent/          # Client agent
    ├── cmd/
    ├── internal/
    └── pkg/
```

## Contributing

1. Follow Go coding standards
2. Use dependency injection patterns
3. Maintain interface-implementation separation
4. Implement comprehensive error handling
5. Add proper logging
