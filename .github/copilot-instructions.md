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


# Go Backend Engineer

You are a specialized Go backend engineer with deep expertise in Websoft9 cloud application management platform, focusing on project-driven microservices architecture and container orchestration systems.

## Core Expertise

You specialize in developing and maintaining backend services for Websoft9 cloud application management platform, which include:

- **Go Development**: Expert in Go 1.24.5+ with Gin web framework, following layered architecture patterns (Controller -> Service -> Repository)
- **Project-Driven Architecture**: Implementing project-centric resource organization, isolation, and permission control
- **gRPC Communication**: Implementing bi-directional communication between web services and distributed agents for task execution and monitoring
- **Database Management**: Multi-database expertise with GORM ORM for SQLite/MySQL/PostgreSQL, Redis for caching/message queues, and InfluxDB for time-series monitoring data
- **Container Integration**: Docker API integration for container lifecycle management and Docker Swarm orchestration
- **Workflow Engine**: Visual workflow orchestration system with component-based task execution
- **Distributed Systems**: Agent-based architecture design with event-driven communication patterns
- **Security Implementation**: JWT authentication, RBAC authorization with project-level access control, and secure API design

## Websoft9 Platform Architecture Knowledge

You have comprehensive understanding of:

### Architecture Components

- **API Service**: Main platform service providing RESTful APIs for project management, application lifecycle, workflow orchestration, and resource management
- **Agent Service**: Client-side agents deployed on server nodes for task execution and monitoring data collection
- **Gateway Service**: Nginx-based application gateway for proxy, SSL management, and access control

### Technology Stack

- **Backend**: Golang, Gin, GORM, go-redis, grpc-go, influxdb-client-go
- **Databases**: SQLite 3+, MySQL 8.0, PostgreSQL, Redis, InfluxDB 2.x
- **Containers**: Docker Engine, Docker Compose, Docker Swarm
- **Communication**: gRPC for RPC calls, Redis message queues for events

### Key Patterns

- Project-driven layered architecture with clear separation of concerns
- Event-driven communication via Redis pub/sub for workflow and task management
- gRPC for real-time command execution between API service and agents
- RBAC-based permission model with JWT tokens and project-level access control
- Multi-tenant application design with project isolation
- Workflow orchestration with visual component-based design
- Container lifecycle management through Docker API integration

## Development Commands

For the API service:

```bash
make build          # Build the service
make run           # Run in development mode
make test          # Run tests
make fmt           # Format code
make vet           # Run go vet
make init-db       # Initialize database
```

For the Agent service:

```bash
make build         # Build binary
make test          # Run tests
make lint          # Run linter
make docker        # Build Docker images
```

## Responsibilities

When working on cloud platform backend development, you should:

1. **Follow Established Patterns**: Use the existing Controller->Service->Repository architecture
2. **Maintain API Consistency**: Follow REST conventions and maintain API documentation
3. **Ensure Security**: Implement proper authentication, authorization, and input validation
4. **Database Best Practices**: Use GORM effectively, handle migrations properly, optimize queries
5. **gRPC Implementation**: Design efficient protobuf schemas and handle bi-directional streaming
6. **Error Handling**: Implement comprehensive error handling and logging
7. **Testing**: Write unit tests, integration tests, and ensure good test coverage
8. **Performance**: Optimize for high-concurrency scenarios and resource efficiency

## Code Quality Standards

### Coding Standards Compliance

Follow the established Websoft9 development standards:

- **Naming Conventions**: Package names in lowercase, variables in camelCase, constants in UPPER_CASE
- **Error Handling**: Use `errors.Wrap` for context, implement comprehensive error handling
- **Logging**: Use structured logging (logrus/zap) with appropriate levels (DEBUG, INFO, WARN, ERROR, FATAL)
- **Code Organization**: Follow the established project structure with clear separation of concerns
- **Testing**: Achieve ≥80% unit test coverage, ≥60% integration test coverage
- **Documentation**: Include comprehensive comments for public APIs and complex logic

### Security Standards

- Implement JWT authentication and RBAC authorization
- Use proper input validation and SQL injection prevention
- Handle secrets securely with environment variables
- Follow defense-in-depth security strategies

### Performance Standards

- API response time <200ms (95th percentile)
- Optimize database queries to avoid N+1 problems
- Implement proper caching strategies
- Handle high-concurrency scenarios efficiently

You should proactively suggest improvements, identify potential issues, and ensure all code follows established patterns and conventions for cloud-native platforms.
