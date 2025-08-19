# Websoft9 Platform

A modern cloud application management solution platform providing comprehensive application deployment, monitoring, and management services throughout the entire lifecycle.

## Project Overview

Websoft9 is a cloud application management platform with project-centric organization, featuring layered architecture design and supporting full lifecycle application management in multi-cloud environments. The platform achieves resource isolation, permission control, cost management, and team collaboration through project-based organization.

### Core Architecture Components

- **api-service**: Core backend service built with Golang + Gin + GORM, providing RESTful API interfaces
- **websoft9-agent**: Client agent deployed on server nodes, communicating with the server via gRPC, responsible for task execution and monitoring data collection

## Platform Feature Modules

### 1. Platform Homepage

- **Overview Dashboard**: Project overview board, monitoring overview board
- **Application Quick Navigation**: Quick access portal for deployed applications

### 2. Project Management (Core Module)

- **Dashboard**: Monitoring board, task board, resource board
- **Folders**: Project folders, personal folder management
- **Application Management**: Query, uninstall, update, publish, clone, configure, migrate
- **Workflows**: Canvas orchestration, component management, task scheduling
- **Project Resources**: Resource groups, servers, secrets, databases, application gateways, certificates, cloud resources
- **Project Teams**: Member management, role assignment, permission control
- **Project Settings**: Environment variables, project configuration

### 3. Application Marketplace

- **Application Marketplace List**: Application categorization, search, deployment
- **Application Wishlist**: Requirement submission, voting, bounty mechanism

### 4. Platform Management

- **Project Management**: Project creation, modification, deletion
- **Platform Settings**: Basic settings, system updates, system services, SMTP, SMS, Webhook, etc.
- **Security Management**: Role management, permission management, authentication management
- **User Management**: User CRUD operations
- **Alert Notifications**: Alert queries, rule configuration, push management
- **Personal Center**: Personal profile, personal settings
- **Audit Logs**: Operation records, log queries

## Technical Architecture

### API Service (Backend Service)

Core backend service built with Golang + Gin + GORM.

**Tech Stack:**

- Golang 1.24+
- Gin Web Framework
- GORM ORM Framework
- SQLite Database
- Redis Cache
- InfluxDB Time Series Database
- JWT Authentication
- gRPC Communication

### Websoft9 Agent (Client Agent)

Client agent deployed on server nodes, responsible for task execution and monitoring data collection.

**Tech Stack:**

- Golang
- gRPC (grpc-go)
- Redis (go-redis)
- InfluxDB (influxdb-client-go)
- Docker Engine API

## Quick Start

### Environment Requirements

- Go 1.24+
- Docker 20.10+
- SQLite 3.0+ / MySQL 8.0+
- Redis 6.0+
- InfluxDB 2.0+

### Development Environment Setup

1. **Clone Repository**

   ```bash
   git clone <repository-url>
   cd webox
   ```

2. **Start API Service**

   ```bash
   cd api-service
   go mod tidy
   make init-db    # Initialize database
   make run        # Start development service
   ```

3. **Start Agent**

   ```bash
   cd websoft9-agent
   go mod tidy
   make build
   sudo ./websoft9-agent
   ```

### Docker Deployment

```bash
# Build and run API service
cd api-service
docker build -t websoft9/api-service .
docker run -p 8080:8080 -p 9090:9090 websoft9/api-service

# Build and run Agent
cd websoft9-agent
docker build -t websoft9/agent .
docker run --privileged websoft9/agent
```

## Development Toolchain

- **API Testing**: [Apifox](https://apifox.com/)
- **CI/CD**: GitHub Actions
- **AI Programming**: [GitHub Copilot](https://github.com/features/copilot), [Claude Code](https://docs.anthropic.com/zh-CN/docs/claude-code/overview)

## Project Structure

```text
webox/
├── api-service/              # Backend API service
│   ├── main.go
│   ├── internal/
│   │   ├── config/          # Configuration management
│   │   ├── controller/      # Controller layer
│   │   ├── service/         # Business logic layer
│   │   ├── repository/      # Data access layer
│   │   ├── middleware/      # Middleware
│   │   └── model/          # Data models
│   ├── pkg/                # Common packages
│   └── docs/               # API documentation
└── websoft9-agent/         # Client agent
    ├── cmd/                # Command line entry
    ├── internal/           # Internal packages
    │   ├── agent/         # Agent core logic
    │   ├── monitor/       # Monitoring data collection
    │   └── executor/      # Task executor
    └── pkg/               # Common packages
```

## Development Standards

1. Follow Go language coding standards
2. Use dependency injection patterns
3. Maintain interface-implementation separation
4. Implement comprehensive error handling
5. Add appropriate logging
