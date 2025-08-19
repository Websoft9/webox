# GitHub Copilot Instructions for Websoft9 Project

## Project Overview

Websoft9 is a modern cloud application management platform that provides application deployment, monitoring, and management services throughout the entire lifecycle. The platform is organized around projects as core units, supporting multi-cloud environments.

### Core Architecture Components

- **api-service**: Core backend service built with Golang + Gin + GORM, providing RESTful API interfaces
- **websoft9-agent**: Client agent deployed on server nodes, communicating with the server via gRPC for task execution and monitoring data collection

## Code Style and Standards

### Go Development Standards

#### Project Structure

Follow the standard Go project layout:

```text
api-service/
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── controller/        # API routes and handlers
│   ├── service/           # Business logic layer
│   ├── repository/        # Data access layer
│   ├── model/             # Data models
│   ├── middleware/        # Middleware
│   └── config/            # Configuration management
├── pkg/                   # Public library code
├── api/                   # API definition files (OpenAPI/Swagger)
└── configs/               # Configuration file templates
```

#### Naming Conventions

- Package names: lowercase, short, meaningful nouns
- Variable names: camelCase, starting with lowercase
- Constants: ALL_CAPS with underscores
- Functions: camelCase, starting with uppercase (public) or lowercase (private)
- Structs: camelCase, starting with uppercase

#### Error Handling

- Use standard error interface
- Error messages should be clear and specific
- Use errors.Wrap to add context information

```go
import "github.com/pkg/errors"

func (s *UserService) GetUser(id int64) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, errors.Wrapf(err, "failed to get user with id %d", id)
    }
    return user, nil
}
```

#### Logging Standards

- Use structured logging (logrus or zap)
- Log levels: DEBUG, INFO, WARN, ERROR, FATAL
- Include necessary context information

```go
import "github.com/sirupsen/logrus"

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    logger := logrus.WithFields(logrus.Fields{
        "operation": "CreateUser",
        "username":  req.Username,
    })
    
    logger.Info("Creating new user")
    // Implementation logic
}
```

### API Design Standards

#### RESTful API Design

- Use standard HTTP methods: GET, POST, PUT, DELETE, PATCH
- URL design follows RESTful principles
- Use appropriate HTTP status codes

```go
// Route definition example
func SetupRoutes(r *gin.Engine) {
    api := r.Group("/api/v1")
    {
        users := api.Group("/users")
        {
            users.GET("", userHandler.ListUsers)           // GET /api/v1/users
            users.POST("", userHandler.CreateUser)         // POST /api/v1/users
            users.GET("/:id", userHandler.GetUser)         // GET /api/v1/users/:id
            users.PUT("/:id", userHandler.UpdateUser)      // PUT /api/v1/users/:id
            users.DELETE("/:id", userHandler.DeleteUser)   // DELETE /api/v1/users/:id
        }
    }
}
```

#### Standardized Response Format

```go
type APIResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
    APIResponse
    Pagination *PaginationInfo `json:"pagination,omitempty"`
}
```

### Database Standards

#### GORM Usage Standards

```go
type User struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    
    Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
    Email    string `gorm:"uniqueIndex;size:100;not null" json:"email"`
    Password string `gorm:"size:255;not null" json:"-"`
    Status   int    `gorm:"default:1" json:"status"`
}

// Table name
func (User) TableName() string {
    return "users"
}
```

## Security Guidelines

### Sensitive Information Handling

- Never hardcode sensitive information
- Use environment variables for configuration
- Encrypt sensitive data in storage

```go
// Correct example - using environment variables
type Config struct {
    DBPassword string `env:"DB_PASSWORD,required"`
    APIKey     string `env:"API_KEY,required"`
    JWTSecret  string `env:"JWT_SECRET,required"`
}
```

### Input Validation and Sanitization

- Validate all external inputs
- Use parameterized queries to prevent SQL injection
- Implement proper error handling

```go
// Correct example - using parameterized queries
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, errors.Wrap(err, "failed to get user by username")
    }
    return &user, nil
}
```

### Authentication and Authorization

- Implement JWT-based authentication
- Use RBAC (Role-Based Access Control) for authorization
- Implement proper session management

## Testing Standards

### Unit Testing

- Use table-driven tests
- Mock external dependencies
- Aim for 80%+ test coverage

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        request *CreateUserRequest
        setup   func(*MockUserRepository)
        want    *User
        wantErr bool
    }{
        {
            name: "successful user creation",
            request: &CreateUserRequest{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "password123",
            },
            setup: func(repo *MockUserRepository) {
                repo.On("Create", mock.Anything, mock.Anything).
                    Return(&User{ID: 1, Username: "testuser"}, nil)
            },
            want: &User{ID: 1, Username: "testuser"},
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Git Workflow

### Commit Message Format

Use Conventional Commits specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation updates
- `style`: Code formatting
- `refactor`: Code refactoring
- `test`: Test-related changes
- `chore`: Build process or auxiliary tool changes

### Branch Naming

- `main`: Production branch
- `develop`: Development branch
- `feature/feature-description`: Feature branches
- `release/version-number`: Release branches
- `hotfix/issue-description`: Hotfix branches

## Performance Requirements

- API response time < 200ms (95th percentile)
- Database query optimization, avoid N+1 problems
- Proper use of caching mechanisms
- Implement connection pooling for databases

## Documentation Standards

### Code Comments

- Use clear and concise comments
- Document public APIs and complex logic
- Include examples where appropriate

```go
// UserService handles user-related business logic.
type UserService struct {
    repo UserRepository
}

// CreateUser creates a new user with the given information.
// It returns the created user or an error if the operation fails.
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Implementation
}
```

### API Documentation

- Use OpenAPI/Swagger specifications
- Include request/response examples
- Document error codes and messages

## Project-Specific Context

### Core Business Concepts

- **Project**: Core organizational unit for resource isolation and permission control
- **Application**: Deployable software packages managed within projects
- **Workflow**: Visual orchestration of tasks and components
- **Resource Group**: Logical grouping of infrastructure resources
- **Agent**: Client-side component for task execution and monitoring

### Key Features to Implement

1. **Project Management**: Project creation, member management, resource allocation
2. **Application Lifecycle**: Deploy, update, configure, migrate, uninstall applications
3. **Workflow Engine**: Visual workflow designer with predefined components
4. **Resource Management**: Servers, databases, certificates, cloud resources
5. **Monitoring & Alerting**: System metrics, application health, alert rules
6. **Security & Audit**: RBAC, audit logs, security scanning

### Technology Stack Preferences

- **Backend**: Go 1.24+, Gin, GORM, SQLite/MySQL, Redis, InfluxDB
- **Communication**: gRPC for internal services, REST for external APIs
- **Authentication**: JWT tokens with refresh mechanism
- **Containerization**: Docker with multi-stage builds
- **CI/CD**: GitHub Actions with automated testing and security scanning

## Code Generation Guidelines

When generating code, please:

1. Follow the established project structure and patterns
2. Include proper error handling and logging
3. Add appropriate validation and security measures
4. Write accompanying unit tests
5. Include relevant documentation and comments
6. Consider performance implications
7. Ensure code is production-ready and follows best practices
8. Use dependency injection patterns where appropriate
9. Implement proper resource cleanup (defer statements)
10. Follow the principle of least privilege for security
