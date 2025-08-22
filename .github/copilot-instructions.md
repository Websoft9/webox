# GitHub Copilot Instructions for Websoft9

## Project Overview
Websoft9: Cloud application management platform with two Go services:
- **api-service**: REST API (Gin + GORM + SQLite/Redis + i18n)  
- **websoft9-agent**: System agent (monitoring + task execution)

## Core Patterns (MANDATORY)

### 1. Layered Architecture
```go
internal/
├── controller/     # HTTP handlers
├── service/        # Business logic  
├── repository/     # Data access
├── interface/      # Contracts
├── model/          # Domain models
└── dto/           # Request/Response DTOs
```

### 2. Error Handling (REQUIRED)
```go
// Use centralized errors
if existingUser != nil {
    return nil, errors.NewAppError(errors.CodeConflict, "用户已存在")
}

// Controllers use middleware
if err != nil {
    errors.HandleError(ctx, err)
    return
}
```

### 3. Context Usage (CRITICAL)
```go
// ALWAYS pass context through all layers
func (s *service) Method(ctx context.Context, req *dto.Request) error {
    return s.repo.Method(ctx, data)
}

// ALWAYS use WithContext for database
return r.db.WithContext(ctx).Create(user).Error
```

## Go Rules

### Naming
- Files: `user.go` (not `user_controller.go`)
- Interfaces: `UserService` (not `IUserService`)
- Methods: PascalCase exported, camelCase private

### Memory & Performance
```go
// Use pointers for large structs
func (s *service) UpdateUser(ctx context.Context, user *model.User) error

// Always close resources
defer file.Close()

// Prefer interfaces for dependencies
type UserService interface { /* methods */ }
```

### Error Patterns
```go
// Early return
if err := validate(req); err != nil {
    return nil, err
}

// Never ignore errors
result, err := operation()
if err != nil {
    return errors.NewAppError(errors.CodeInternalError, "操作失败")
}
```

## API Design
```go
// Current routes
POST   /api/v1/auth/register
POST   /api/v1/auth/login
GET    /api/v1/users/profile
PUT    /api/v1/users/profile

// DTOs with validation
type UserRegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

## Anti-Patterns (NEVER)
```go
// ❌ Don't ignore context
func (r *repo) Get(id uint) error

// ❌ Don't use panic in business logic  
if user == nil { panic("user is nil") }

// ❌ Don't ignore errors
result, _ := operation()

// ✅ DO this instead
func (r *repo) Get(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).First(&user, id).Error
}
```

## Key Dependencies
- **Gin**: HTTP framework
- **GORM**: ORM with SQLite
- **Redis**: Sessions & caching
- **i18n**: en-US (default), zh-CN

## Development
```bash
# Quick start
cd api-service && make dev
./scripts/run-tests.sh
cd docker && docker-compose up -d
```

---
*Keep examples project-specific and concrete.*
