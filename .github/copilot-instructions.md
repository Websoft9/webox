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
├── dto/                    # Data Transfer Objects
│   ├── request/user.go     # Input validation with binding tags
│   └── response/user.go    # Output formatting, exclude sensitive fields
└── interface/              # Contracts for dependency inversion
    ├── service/user.go     # Service interfaces
    └── repository/user.go  # Repository interfaces
```

### Dependency Injection Pattern
Main initialization in `main.go` follows strict order:

```go
// 1. Logger → 2. Config → 3. Database → 4. External Services → 5. Layers → 6. Router
jwtAuth := auth.NewJWTAuth(cfg.JWT.Secret, cfg.JWT.ExpireTime)
userRepo := repository.NewUserRepository(db)
userService := service.NewUserService(userRepo, jwtAuth, zapLogger)
userController := controller.NewUserController(userService, zapLogger)
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

// Use predefined errors for common cases
if user == nil {
    return nil, errors.ErrUserNotFound
}

// Create custom business errors with specific codes
return nil, errors.NewAppError(errors.CodeUserAlreadyExists, "用户名已存在")
```

### Validation & Request Binding Pattern
Follow strict validation in controllers:

```go
// Controller method pattern
func (c *UserController) Register(ctx *gin.Context) {
    var req request.UserRegisterRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        c.logger.WarnContext(ctx, "请求参数绑定失败", logger.ErrorField(err))
        errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, "请求参数无效"))
        return
    }
    // Business logic...
}
```

### Logging Pattern
Use structured logging with context throughout:

```go
s.logger.InfoContext(ctx, "用户注册开始", 
    logger.String("username", req.Username))
```

### Router & Middleware Pattern
Routes organized by protection level in `internal/router/router.go`:

```go
// Global middleware order: Logger → CORS → ErrorHandler → RequestValidator
r.Use(middleware.LoggerMiddleware(log))
r.Use(middleware.CORS())
r.Use(middleware.ErrorHandler(log))
r.Use(middleware.RequestValidator(log))

// Public routes (no auth)
auth := v1.Group("/auth")
auth.POST("/register", controllers.UserController.Register)

// Protected routes (JWT required)
protected := v1.Group("/")
protected.Use(middleware.JWTAuth(cfg))
users := protected.Group("/users")
```

### DTO Pattern
Always use request/response DTOs for API boundaries:
- `internal/dto/request/` - Input validation with binding tags
- `internal/dto/response/` - Output formatting, exclude sensitive fields
- `internal/dto/common.go` - Base structures (pagination, sorting, search)

### Pagination & Search Pattern
Use standardized list request pattern:
```go
// Embed BaseListRequest for pagination, sorting, search
type UserListRequest struct {
    dto.BaseListRequest
    Status string `form:"status" binding:"omitempty,oneof=active inactive banned"`
    Role   string `form:"role" binding:"omitempty,oneof=admin user guest"`
}

// Use NewPaginationResponse for consistent list responses
response := dto.NewPaginationResponse(req.Page, req.PageSize, total, users)
```

### Response Format Pattern
Use consistent JSON response structure:

```go
// Success response
response.Success(ctx, "操作成功", data)

// Error response (use response.Error directly or through error handler)
response.Error(ctx, http.StatusBadRequest, "参数错误", err.Error())
// OR via centralized error handling
errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, "参数错误"))
```

## Development Workflows

### Building & Testing
```bash
# API Service
cd api-service
make build          # Build binary
make run            # Development server  
make test           # Run tests
make dev            # Hot reload (requires air)
make init-db        # Initialize SQLite database
make reset-db       # Reset database and reinitialize

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

## Gin Framework Best Practices

### Middleware Chain Order
Critical: middleware order affects behavior
```go
// Correct order in router.go:
r.Use(middleware.LoggerMiddleware(log))    // 1. Logging first
r.Use(middleware.CORS())                   // 2. CORS handling
r.Use(middleware.ErrorHandler(log))        // 3. Error recovery
r.Use(middleware.RequestValidator(log))    // 4. Request validation
// Route-specific middleware like JWT comes after
```

### Context & Request Handling
Always pass gin.Context through layers:
```go
// Controller → Service → Repository (all use ctx context.Context)
func (c *UserController) Register(ctx *gin.Context) {
    user, err := c.userService.Register(ctx, &req)  // Pass gin.Context
}
```

### Input Validation Best Practices
Multi-layer validation approach:
```go
// 1. JSON binding validation (struct tags)
type UserRegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// 2. Custom business validation in service layer
func (s *userService) Register(ctx context.Context, req *request.UserRegisterRequest) {
    if err := validator.ValidateUsername(req.Username); err != nil {
        return nil, err
    }
    if err := validator.ValidateEmail(req.Email); err != nil {
        return nil, err
    }
}
```

### Error Response Format
Standardized error structure:
```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": "validation error details"
}
```

## Agent-Specific Patterns

### Task Execution
Agent follows this pattern for all operations:
```go
// internal/task/executor.go - handles tasks from gRPC
type TaskResult struct {
    TaskID   string                 `json:"task_id"`
    Status   string                 `json:"status"`  // success/failed/timeout
    Message  string                 `json:"message"`
    Data     map[string]interface{} `json:"data"`
    Duration int64                  `json:"duration"` // 执行时间(毫秒)
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
6. **Middleware Order**: Logger → CORS → ErrorHandler → RequestValidator → JWT (route-specific)
7. **Gin Context**: Never store gin.Context in structs, always pass as parameter
8. **Response Consistency**: Always use `response.Success()` and `errors.HandleError()`
9. **Validation**: Combine struct tags with business validation in service layer
10. **Database Migration**: Add new models to `main.go` AutoMigrate call

## Security Best Practices

### JWT Implementation
```go
// JWT middleware extracts and validates tokens
func JWTAuth(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        claims, err := jwtAuth.ValidateToken(tokenString)
        if err != nil {
            response.Error(c, http.StatusUnauthorized, "Invalid token", err.Error())
            c.Abort()
            return
        }
        
        // Set user context for downstream handlers
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        c.Next()
    }
}
```

### Input Sanitization
Always validate and sanitize inputs:
```go
// Clean inputs before processing
req.Username = strings.TrimSpace(req.Username)
req.Email = strings.ToLower(strings.TrimSpace(req.Email))

// Use pkg/validator for complex validation rules
if err := validator.ValidateEmail(req.Email); err != nil {
    return errors.NewAppError(errors.CodeValidationError, "邮箱格式无效")
}
```
