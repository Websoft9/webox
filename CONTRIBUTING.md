# Websoft9 Contributing Guide

Welcome to contribute to the Websoft9 project! This guide will help you understand how to contribute to the project.

## Table of Contents

- [Project Overview](#project-overview)
- [Development Environment Setup](#development-environment-setup)
- [Development Standards](#development-standards)
- [Contribution Process](#contribution-process)
- [Code Review](#code-review)
- [Testing Standards](#testing-standards)
- [Security Standards](#security-standards)
- [Community Participation](#community-participation)

## Project Overview

Websoft9 is a modern cloud application management solution platform with layered architecture design, providing full lifecycle services for application deployment, monitoring, and management.

### Core Components

- **API Service**: Backend service based on Golang + Gin + GORM
- **Websoft9 Agent**: Client agent deployed on server nodes
- **Web UI**: Vue 3 + Element Plus frontend interface (planned)

### Technology Stack

- **Backend**: Go 1.24+, Gin, GORM, SQLite/MySQL, Redis, InfluxDB
- **Frontend**: Vue 3, TypeScript, Element Plus, Pinia
- **Infrastructure**: Docker, GitHub Actions

## Development Environment Setup

### Environment Requirements

- Go 1.24+
  - [golangci-lint 1.64.8](https://github.com/golangci/golangci-lint)
  - [gosec 2.22.7+](https://github.com/securego/gosec)
- Node.js 18+
- Docker 20.10+
- Git 2.30+
- Linux

### Quick Start

1. **Fork and Clone Repository**

   ```bash
   git clone https://github.com/Websoft9/webox.git
   cd webox
   ```

2. **Setup Development Environment**

   ```bash
   # Install Go dependencies
   cd api-service
   go mod tidy

   # Initialize database
   make init-db

   # Start API service
   make run
   ```

3. **Start Agent (Optional)**

   ```bash
   cd websoft9-agent
   go mod tidy
   make build
   sudo ./websoft9-agent
   ```

### Recommended Development Tools

- **IDE**: VS Code, GoLand
- **API Testing**: Apifox, Postman
- **Database**: DBeaver, TablePlus
- **Container**: Docker Desktop

## Development Standards

### Code Style

#### Go Code Standards

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` and `goimports` to format code
- Use `golangci-lint` for code checking
- Use `gosec` for security checks

```go
// Correct function comments and naming
// CreateUser creates a new user with the given information.
// It returns the created user or an error if the operation fails.
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Validate input parameters
    if err := s.validateCreateUserRequest(req); err != nil {
        return nil, errors.Wrap(err, "invalid create user request")
    }

    // Create user
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, errors.Wrap(err, "failed to create user in database")
    }

    return user, nil
}
```

#### Project Structure

```text
api-service/
├── cmd/server/           # Application entry point
├── internal/            # Private application code
│   ├── controller/      # API controllers
│   ├── service/         # Business logic layer
│   ├── repository/      # Data access layer
│   ├── model/           # Data models
│   ├── middleware/      # Middleware
│   └── config/          # Configuration management
├── pkg/                 # Public library code
├── api/                 # API documentation
├── configs/             # Configuration files
└── docs/                # Project documentation
```

### Naming Conventions

- **Package names**: Lowercase, short, meaningful nouns
- **Variable names**: CamelCase, starting with lowercase
- **Constant names**: All uppercase, underscore separated
- **Function names**: CamelCase, public functions start with uppercase
- **Structs**: CamelCase, starting with uppercase

### API Design Standards

#### RESTful API

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

#### Response Format

```go
type APIResponse struct {
    Success   bool        `json:"success"`
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Error     *APIError   `json:"error,omitempty"`
    Timestamp int64       `json:"timestamp"`
}

type PaginatedResponse struct {
    Items      interface{} `json:"items"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    TotalPages int         `json:"total_pages"`
}
```

### Error Handling

- Use standard error interface
- Error messages should be clear and specific
- Use `api-service/pkg/i18n` to add internationalization
- Use `api-service/pkg/errors` to add context information

```go
import "api-service/pkg/errors"

// GetByID retrieves a permission by ID
func (r *permissionRepository) GetByID(ctx context.Context, id uint) (*model.Permission, error) {
    var permission model.Permission
    err := r.db.WithContext(ctx).Where("status != -1").First(&permission, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.NewAppErrorWithMessage(errors.CodeRecordNotFound, "permission not found")
        }
        return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get permission by ID")
    }
    return &permission, nil
}
```

### Logging Standards

- Use structured logging (recommend zap, it's implemented by default.)
- Log levels: DEBUG, INFO, WARN, ERROR, FATAL
- Include necessary context information

```go
import "api-service/pkg/logger"

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    s.logger.InfoContext(ctx, "Creating user", logger.String("username", req.Username))

    logger.Info("Creating new user")

    user, err := s.repo.Create(req)
    if err != nil {
        s.logger.ErrorContext(ctx, "Failed to create user", logger.ErrorField(err))

        return nil, errors.WrapError(err, errors.CodeInternalError, "Failed to create user")
    }

    s.logger.InfoContext(ctx, "User created successfully", logger.Uint("user_id", user.ID))
    return s.buildUserResponse(user), nil
}
```

## Contribution Process

### Git Workflow

We adopt the **Git Flow** workflow model:

```text
main (production branch)
├── develop (development branch)
├── release/v1.2.0 (release branch)
├── feature/user-authentication (feature branch)
└── hotfix/critical-bug-fix (hotfix branch)
```

### Branch Naming Conventions

| Branch Type | Naming Format | Example | Purpose |
|-------------|---------------|---------|---------|
| Main branch | `main` | `main` | Production environment code |
| Development branch | `develop` | `develop` | Development environment code |
| Feature branch | `feature/feature-description` | `feature/user-management` | New feature development |
| Release branch | `release/version-number` | `release/v1.2.0` | Release preparation |
| Hotfix branch | `hotfix/issue-description` | `hotfix/login-error` | Emergency fixes |
| Bugfix branch | `bugfix/issue-description` | `bugfix/api-validation` | General fixes |

### Commit Message Standards

Adopt **Conventional Commits** specification:

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### Commit Types

| Type | Description | Example |
|------|-------------|---------|
| `feat` | New feature | `feat(auth): add JWT token refresh mechanism` |
| `fix` | Bug fix | `fix(api): handle null pointer in user service` |
| `docs` | Documentation update | `docs(readme): update installation instructions` |
| `style` | Code formatting | `style(user): format code with gofmt` |
| `refactor` | Code refactoring | `refactor(db): extract connection logic to separate package` |
| `test` | Test related | `test(user): add unit tests for user service` |
| `chore` | Build process or auxiliary tool changes | `chore(deps): update golang to 1.24` |
| `perf` | Performance optimization | `perf(api): optimize database queries` |
| `ci` | CI/CD related | `ci(github): add automated testing workflow` |

#### Commit Message Examples

```bash
feat(auth): add JWT token refresh mechanism

- Implement automatic token refresh
- Add refresh token storage
- Handle token expiration gracefully

Closes #123
```

### Feature Development Process

1. **Create Feature Branch**

   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/user-management
   ```

2. **Develop Feature**

   ```bash
   # Write code
   # Add tests
   # Update documentation

   git add .
   git commit -m "feat(user): add user creation functionality"
   ```

3. **Push Branch and Create PR**

   ```bash
   git push origin feature/user-management
   # Create Pull Request on GitHub
   ```

4. **Code Review and Merge**

   - Wait for code review
   - Modify code based on feedback
   - Merge to develop branch after review approval

5. **Clean Up Branch**

   ```bash
   git checkout develop
   git pull origin develop
   git branch -d feature/user-management
   ```

### Pull Request Standards

#### PR Title Format

```text
<type>[scope]: <description>
```

Examples:

- `feat(auth): add OAuth2 integration`
- `fix(api): resolve memory leak in user service`
- `docs(readme): update development setup guide`

#### PR Description Template

```markdown
## Change Type
- [ ] New feature (feature)
- [ ] Bug fix (fix)
- [ ] Documentation update (docs)
- [ ] Code refactoring (refactor)
- [ ] Performance optimization (perf)
- [ ] Test related (test)
- [ ] Other (chore)

## Change Description
Brief description of the changes and purpose.

## Related Issues
Closes #123
Fixes #456

## Testing Instructions
- [ ] Added unit tests
- [ ] Added integration tests
- [ ] Performed manual testing
- [ ] Test coverage ≥ 80%

## Checklist
- [ ] Code follows project coding standards
- [ ] Updated relevant documentation
- [ ] Passed all automated tests
- [ ] Performed code self-review
- [ ] No obvious performance issues

## Screenshots/Demo
If UI changes are involved, please provide screenshots or demo videos.

## Additional Notes
Any other information that needs to be explained.
```

## Code Review

### Review Checklist

#### Functionality Check

- [ ] Is the functionality correctly implemented according to requirements
- [ ] Are boundary conditions properly handled
- [ ] Is error handling comprehensive
- [ ] Does performance meet requirements

#### Code Quality Check

- [ ] Does code follow project standards
- [ ] Is there duplicate code
- [ ] Are variable and function names clear
- [ ] Are comments sufficient and accurate
- [ ] Does it follow SOLID principles

#### Security Check

- [ ] Are there security vulnerabilities
- [ ] Is sensitive information properly handled
- [ ] Is input validation sufficient
- [ ] Is permission control correct

#### Test Check

- [ ] Is there sufficient test coverage
- [ ] Are test cases reasonable
- [ ] Are error scenarios tested

### Review Feedback Standards

Use the following tags for feedback:

- `MUST`: Issues that must be fixed
- `SHOULD`: Issues that should be fixed
- `COULD`: Optional improvement suggestions
- `QUESTION`: Questions that need clarification
- `PRAISE`: Code worth praising

Example:

```text
MUST: There's a risk of null pointer exception here, need to add nil check.

SHOULD: Suggest extracting this magic number as a constant to improve code readability.

COULD: Consider using a more concise approach: `return err != nil`

QUESTION: What's the time complexity of this function? Does it need optimization?

PRAISE: This error handling is well written, providing clear context information.
```

## Testing Standards

### Testing Strategy

Adopt the test pyramid model:

```text
    /\
   /  \  E2E Tests (10%)
  /____\
 /      \
/        \ Integration Tests (20%)
\________/
\        /
 \______/ Unit Tests (70%)
```

### Unit Testing

#### Go Unit Test Example

```go
// user_service_test.go
package user

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

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
        {
            name: "invalid email format",
            request: &CreateUserRequest{
                Username: "testuser",
                Email:    "invalid-email",
                Password: "password123",
            },
            setup: func(repo *MockUserRepository) {},
            want: nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := &MockUserRepository{}
            tt.setup(repo)

            service := NewUserService(repo)
            got, err := service.CreateUser(context.Background(), tt.request)

            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }

            repo.AssertExpectations(t)
        })
    }
}
```

### Test Coverage Requirements

- Unit test coverage ≥ 80%
- Integration test coverage ≥ 60%
- Critical business logic coverage ≥ 90%

### Running Tests

```bash
# Run all tests
make test

# Run tests and generate coverage report
make test-coverage

# Run tests for specific package
go test -v ./internal/service/...

# Run integration tests
make test-integration
```

## Security Standards

### Code Security

#### Sensitive Information Handling

**❌ Wrong Example - Hardcoded Sensitive Information**

```go
const (
    DBPassword = "password123"
    APIKey     = "sk-1234567890abcdef"
    JWTSecret  = "my-secret-key"
)
```

**✅ Correct Example - Using Environment Variables**

```go
type Config struct {
    DBPassword string `env:"DB_PASSWORD,required"`
    APIKey     string `env:"API_KEY,required"`
    JWTSecret  string `env:"JWT_SECRET,required"`
}

func LoadConfig() (*Config, error) {
    var config Config
    if err := env.Parse(&config); err != nil {
        return nil, errors.Wrap(err, "failed to parse config")
    }
    return &config, nil
}
```

#### Input Validation

```go
import (
    "github.com/go-playground/validator/v10"
    "html"
    "strings"
)

var validate = validator.New()

type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func (r *CreateUserRequest) Validate() error {
    if err := validate.Struct(r); err != nil {
        return errors.Wrap(err, "validation failed")
    }

    // Custom validation
    if err := ValidatePasswordStrength(r.Password); err != nil {
        return err
    }

    return nil
}

func (r *CreateUserRequest) Sanitize() {
    r.Username = strings.TrimSpace(r.Username)
    r.Email = strings.ToLower(strings.TrimSpace(r.Email))
    // HTML escape to prevent XSS
    r.Username = html.EscapeString(r.Username)
}
```

#### SQL Injection Protection

```go
// ❌ Wrong Example - Vulnerable to SQL injection
func GetUserByUsername(username string) (*User, error) {
    query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
    rows, err := db.Query(query)
    // ...
}

// ✅ Correct Example - Using parameterized queries
func GetUserByUsername(username string) (*User, error) {
    query := "SELECT * FROM users WHERE username = ?"
    rows, err := db.Query(query, username)
    // ...
}

// ✅ Using GORM (Recommended)
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, errors.Wrap(err, "failed to get user by username")
    }
    return &user, nil
}
```

### Security Checklist

#### Code Review Security Check

- [ ] No hardcoded passwords, keys, or sensitive information
- [ ] All user inputs are validated and sanitized
- [ ] Use parameterized queries to prevent SQL injection
- [ ] Implement appropriate authentication and authorization mechanisms
- [ ] Sensitive data is encrypted during transmission and storage
- [ ] Implement appropriate error handling without leaking sensitive information
- [ ] Use secure random number generators
- [ ] Implement appropriate logging and auditing

## Community Participation

### Reporting Issues

If you find bugs or have feature suggestions, please:

1. Search existing Issues to avoid duplicate reports
2. Use appropriate Issue templates
3. Provide detailed reproduction steps and environment information
4. For security issues, please contact maintainers privately

### Feature Requests

1. Describe your needs in Issues
2. Explain why this feature is useful
3. Provide specific use cases
4. Consider backward compatibility

### Documentation Contributions

- Fix errors in documentation
- Improve clarity of existing documentation
- Add missing documentation
- Translate documentation to other languages

### Community Code of Conduct

We are committed to providing a friendly, safe, and welcoming environment for everyone. Please:

- Use friendly and inclusive language
- Respect different viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Help

### Contact Information

- **GitHub Issues**: Report bugs and feature requests
- **GitHub Discussions**: General discussions and questions

### Resource Links

- [Project Documentation](./docs/)
- [API Documentation](./api-service/docs/)
- [Development Standards](./docs/开发规范.md)
- [Architecture Design](./docs/designs/)

## Acknowledgments

Thanks to all developers who have contributed to the Websoft9 project! Your contributions make this project better.

---

**Happy Coding! 🚀**
