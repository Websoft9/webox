# GitHub Copilot Instructions for Websoft9

## Project Architecture Overview

Websoft9 is a cloud application management platform with two main Go services:

- **api-service**: Core backend API (Gin + GORM + SQLite/Redis/InfluxDB)
- **websoft9-agent**: Node agent for task execution and monitoring (gRPC client)

Communication flow: `Client → API Service → gRPC → Agent → Docker/System operations`

## Code Patterns & Conventions

### Layered Architecture (api-service)
Follow the established 4-layer pattern with dependency injection:

```go
// Controller → Service → Repository → Model
internal/
├── controller/user.go       # HTTP handlers, single noun naming
├── service/user.go         # Business logic implementation  
├── repository/user.go      # Data access with context.Context
└── interface/              # Contracts for dependency inversion
    ├── service/user.go     # Service interfaces
    └── repository/user.go  # Repository interfaces
```

### Error Handling Pattern
Use the centralized error system in `pkg/errors/`:

```go
// Service layer - wrap and return custom errors
if err != nil {
    return nil, errors.WrapError(err, errors.CodeInternalError, "创建用户失败")
}

// Controller layer - let middleware handle errors
if err != nil {
    errors.HandleError(ctx, err)
    return
}
```

### Logging Pattern
Use structured logging with context throughout:

```go
s.logger.InfoContext(ctx, "用户注册开始", 
    logger.String("username", req.Username))
```

### DTO Pattern
Always use request/response DTOs for API boundaries:
- `internal/dto/request/` - Input validation with binding tags
- `internal/dto/response/` - Output formatting, exclude sensitive fields

## Development Workflows

### Building & Testing
```bash
# API Service
cd api-service
make build          # Build binary
make run            # Development server
make test           # Run tests
make dev            # Hot reload (requires air)

# Agent
cd websoft9-agent  
make build          # Build agent
make test           # Run tests
```

### Database Operations
- SQLite for development: `data/websoft9.db` 
- Auto-migration in `main.go`: `db.AutoMigrate(&model.User{})`
- Repository pattern with context: `func (r *repo) Create(ctx context.Context, user *model.User) error`

### Docker Development
```bash
# Full stack with Redis
docker-compose -f docker/docker-compose.yml up

# Single service builds
make docker-build && make docker-run
```

## Agent-Specific Patterns

### Task Execution
Agent follows this pattern for all operations:
```go
// internal/task/executor.go - handles tasks from gRPC
type TaskResult struct {
    ID     string `json:"id"`
    Status string `json:"status"`  // success/failed/running
    Output string `json:"output"`
}
```

### Monitoring Collection
- System metrics → InfluxDB (CPU, memory, disk, network)
- Container monitoring via Docker API
- Health checks (HTTP/TCP probes)

### Security Validation
All agent operations use `pkg/security/validator.go`:
- Command validation before execution
- Path sanitization for file operations
- Input validation for all external data

## Dependencies & Integration

### Key External Services
- **Redis**: Session storage and task queues
- **InfluxDB**: Time-series monitoring data
- **gRPC**: Agent ↔ API Service communication
- **Docker Engine API**: Container operations

### Testing Strategy
- Unit tests for business logic: `*_test.go`
- Integration tests in CI pipeline
- Security validation tests for agent operations

## Common Gotchas

1. **Context Usage**: Always pass context.Context through all layers for cancellation
2. **Error Wrapping**: Use `errors.WrapError()` not standard library wrapping
3. **Agent Privileges**: Agent must run as root for system operations
4. **File Naming**: Use single nouns (`user.go` not `user_controller.go`)
5. **gRPC Health**: Agent maintains heartbeat to API service
