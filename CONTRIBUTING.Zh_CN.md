# Websoft9 贡献指南

欢迎参与 Websoft9 项目的开发！本指南将帮助您了解如何为项目做出贡献。

## 目录

- [项目概述](#项目概述)
- [开发环境搭建](#开发环境搭建)
- [开发规范](#开发规范)
- [贡献流程](#贡献流程)
- [代码审查](#代码审查)
- [测试规范](#测试规范)
- [安全规范](#安全规范)
- [社区参与](#社区参与)

## 项目概述

Websoft9 是一个现代化的云应用管理解决方案平台，采用分层架构设计，提供应用部署、监控、管理等全生命周期服务。

### 核心组件

- **API Service**: 基于 Golang + Gin + GORM 的后端服务
- **Websoft9 Agent**: 部署在服务器节点的客户端代理
- **Web UI**: Vue 3 + Element Plus 前端界面（计划中）

### 技术栈

- **后端**: Go 1.24+, Gin, GORM, SQLite/MySQL, Redis, InfluxDB
- **前端**: Vue 3, TypeScript, Element Plus, Pinia
- **基础设施**: Docker, GitHub Actions

## 开发环境搭建

### 环境要求

- Go 1.24+
  - [golangci-lint 1.64.8](https://github.com/golangci/golangci-lint)
  - [gosec 2.22.7+](https://github.com/securego/gosec)
- Node.js 18+
- Docker 20.10+
- Git 2.30+
- Linux

>
> **仅针对中国站用户，请使用代理**
>
> ```shell
> go env -w GOPROXY=<https://mirrors.aliyun.com/goproxy/,direct>
> ```
>

### 快速开始

1. **Fork 并克隆仓库**

   ```bash
   git clone https://github.com/Websoft9/webox.git
   cd webox
   ```

2. **设置开发环境**

   ```bash
   # 安装 Go 依赖
   cd api-service
   go mod tidy

   # 初始化数据库
   make init-db

   # 启动 API 服务
   make run
   ```

3. **启动 Agent（可选）**

   ```bash
   cd websoft9-agent
   go mod tidy
   make build
   sudo ./websoft9-agent
   ```

### 开发工具推荐

- **IDE**: VS Code, GoLand
- **API 测试**: Apifox, Postman
- **数据库**: DBeaver, TablePlus
- **容器**: Docker Desktop

## 开发规范

### 代码风格

#### Go 代码规范

- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `gofmt` 和 `goimports` 格式化代码
- 使用 `golangci-lint` 进行代码检查
- 使用 `gosec` 进行安全检查

```go
// 正确的函数注释和命名
// CreateUser creates a new user with the given information.
// It returns the created user or an error if the operation fails.
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // 验证输入参数
    if err := s.validateCreateUserRequest(req); err != nil {
        return nil, errors.Wrap(err, "invalid create user request")
    }

    // 创建用户
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, errors.Wrap(err, "failed to create user in database")
    }

    return user, nil
}
```

#### 项目结构

```text
api-service/
├── cmd/server/           # 应用程序入口
├── internal/            # 私有应用程序代码
│   ├── controller/      # API 控制器
│   ├── service/         # 业务逻辑层
│   ├── repository/      # 数据访问层
│   ├── model/           # 数据模型
│   ├── middleware/      # 中间件
│   └── config/          # 配置管理
├── pkg/                 # 公共库代码
├── api/                 # API 文档
├── configs/             # 配置文件
└── docs/                # 项目文档
```

### 命名规范

- **包名**: 小写，简短，有意义的名词
- **变量名**: 驼峰命名法，首字母小写
- **常量名**: 全大写，下划线分隔
- **函数名**: 驼峰命名法，公开函数首字母大写
- **结构体**: 驼峰命名法，首字母大写

### API 设计规范

#### RESTful API

- 使用标准 HTTP 方法：GET、POST、PUT、DELETE、PATCH
- URL 设计遵循 RESTful 原则
- 使用合适的 HTTP 状态码

```go
// 路由定义示例
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

#### 响应格式

```go
type APIResponse struct {
    Success   bool        `json:"success"`
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Error     *APIError   `json:"error,omitempty"`
}

type PaginatedResponse struct {
    Items      interface{} `json:"items"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    TotalPages int         `json:"total_pages"`
}
```

### 错误处理

- 使用标准的 error 接口
- 错误信息应该清晰、具体
- 使用 `api-service/pkg/errors` 添加上下文信息

```go
import "api-service/pkg/errors"

// Demo: GetByID retrieves a permission by ID
func (r *permissionRepository) GetByID(ctx context.Context, id uint) (*model.Permission, error) {
    var permission model.Permission
    err := r.db.WithContext(ctx).Where("status != -1").First(&permission, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.NewAppError(errors.CodeRecordNotFound)
        }
        return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
    }
    return &permission, nil
}
```

### 错误定义

- 使用 `api-service/pkg/errors/codes` 定义错误代码、国际化处理（i18n messages key）、HTTP状态码

```go
// Error code constants definition based on API documentation
// Error codes are organized by category with specific ranges for easy identification
const (
    // System related error codes (6000-6999)
    CodeInternalError ErrorCode = 6001 // System Internal Error
)

// CodeToI18nKey maps error codes to their i18n message keys
// These keys should correspond to entries in the i18n locale files
var CodeToI18nKey = map[ErrorCode]string{
    // System related errors (6000-6999)
    CodeInternalError: "system.internal_error",
}

// codeToHTTPStatus maps business error codes to HTTP status codes
// This mapping ensures consistent HTTP responses for different error types
var CodeToHTTPStatus = map[ErrorCode]HTTPCode{
    // System related errors (6000-6999)
    CodeInternalError: http.StatusInternalServerError,
}
```

- 使用 `api-service/pkg/errors/errors` 创建错误对象

```go
// Predefined common errors with internationalization support
// These errors can be reused throughout the application for consistency
var (
    // System related errors (6000-6999)
    ErrInternalError = NewAppErrorWithI18n(CodeInternalError, CodeToI18nKey[CodeInternalError])
)
```

- 遵循`错误处理规范`进行异常错误的处理：

```go
    //示例1: 创建一个错误信息并包含原始错误
    err := errors.NewAppErrorWrapError(err, errors.ErrInternalError)
    return err

    //示例2: 创建一个错误信息
    err := errors.NewAppError(errors.CodeRecordNotFound)
    return err

    //示例3: 直接使用错误信息
    return errors.ErrInternalError
```

### 国际化处理

- 使用 `api-service/configs/lang/<LANGUAGE>.yaml` 定义多语言的翻译和命名统一
- 使用 `pkg/i18n` 进行多语言支持

```go
// 使用 pkg/i18n 进行多语言支持
import (
    "api-service/pkg/errors"
    "api-service/pkg/logger"
)

func (s *userService) CreateUser(ctx *gin.Context) error {
    var req request.CreateUserRequest

    if req.Email == "" {
        logger.ErrorContext(ctx.Request.Context(), "Email cannot be empty", logger.ErrorField(err))
        // 使用自定义的多语言翻译
        return errors.NewAppErrorWithI18n(errors.CodeInvalidEmailFormat, "user.email_required")

        // 使用错误定义预置的多语言翻译
        // return errors.NewAppError(errors.CodeInvalidEmailFormat)
    }
    // 业务逻辑
}
```

### 日志规范

- 使用结构化日志（推荐 zap，项目中已默认实现）
- 日志级别：DEBUG、INFO、WARN、ERROR、FATAL
- 包含必要的上下文信息

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

## 贡献流程

### Git 工作流

我们采用 **Git Flow** 工作流模型：

```text
main (生产分支)
├── develop (开发分支)
├── release/v1.2.0 (发布分支)
├── feature/user-authentication (功能分支)
└── hotfix/critical-bug-fix (修复分支)
```

### 分支命名规范

| 分支类型 | 命名格式 | 示例 | 用途 |
|----------|----------|------|------|
| 主分支 | `main` | `main` | 生产环境代码 |
| 开发分支 | `develop` | `develop` | 开发环境代码 |
| 功能分支 | `feature/功能描述` | `feature/user-management` | 新功能开发 |
| 发布分支 | `release/版本号` | `release/v1.2.0` | 发布准备 |
| 修复分支 | `hotfix/问题描述` | `hotfix/login-error` | 紧急修复 |
| 修复分支 | `bugfix/问题描述` | `bugfix/api-validation` | 一般修复 |

### 提交信息规范

采用 **Conventional Commits** 规范：

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### 提交类型

| 类型 | 描述 | 示例 |
|------|------|------|
| `feat` | 新功能 | `feat(auth): add JWT token refresh mechanism` |
| `fix` | 修复 bug | `fix(api): handle null pointer in user service` |
| `docs` | 文档更新 | `docs(readme): update installation instructions` |
| `style` | 代码格式调整 | `style(user): format code with gofmt` |
| `refactor` | 代码重构 | `refactor(db): extract connection logic to separate package` |
| `test` | 测试相关 | `test(user): add unit tests for user service` |
| `chore` | 构建过程或辅助工具的变动 | `chore(deps): update golang to 1.24` |
| `perf` | 性能优化 | `perf(api): optimize database queries` |
| `ci` | CI/CD 相关 | `ci(github): add automated testing workflow` |

#### 提交信息示例

```bash
feat(auth): add JWT token refresh mechanism

- Implement automatic token refresh
- Add refresh token storage
- Handle token expiration gracefully

Closes #123
```

### 功能开发流程

1. **创建功能分支**

   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/user-management
   ```

2. **开发功能**

   ```bash
   # 编写代码
   # 添加测试
   # 更新文档

   git add .
   git commit -m "feat(user): add user creation functionality"
   ```

3. **推送分支并创建 PR**

   ```bash
   git push origin feature/user-management
   # 在 GitHub 上创建 Pull Request
   ```

4. **代码审查和合并**

   - 等待代码审查
   - 根据反馈修改代码
   - 审查通过后合并到 develop 分支

5. **清理分支**

   ```bash
   git checkout develop
   git pull origin develop
   git branch -d feature/user-management
   ```

### Pull Request 规范

#### PR 标题格式

```text
<type>[scope]: <description>
```

示例：

- `feat(auth): add OAuth2 integration`
- `fix(api): resolve memory leak in user service`
- `docs(readme): update development setup guide`

#### PR 描述模板

```markdown
## 变更类型
- [ ] 新功能 (feature)
- [ ] Bug 修复 (fix)
- [ ] 文档更新 (docs)
- [ ] 代码重构 (refactor)
- [ ] 性能优化 (perf)
- [ ] 测试相关 (test)
- [ ] 其他 (chore)

## 变更描述
简要描述本次变更的内容和目的。

## 相关 Issue
Closes #123
Fixes #456

## 测试说明
- [ ] 已添加单元测试
- [ ] 已添加集成测试
- [ ] 已进行手动测试
- [ ] 测试覆盖率 ≥ 80%

## 检查清单
- [ ] 代码遵循项目编码规范
- [ ] 已更新相关文档
- [ ] 已通过所有自动化测试
- [ ] 已进行代码自查
- [ ] 无明显性能问题

## 截图/演示
如果涉及 UI 变更，请提供截图或演示视频。

## 其他说明
其他需要说明的内容。
```

## 代码审查

### 审查检查清单

#### 功能性检查

- [ ] 功能是否按需求正确实现
- [ ] 边界条件是否正确处理
- [ ] 错误处理是否完善
- [ ] 性能是否满足要求

#### 代码质量检查

- [ ] 代码是否符合项目规范
- [ ] 是否有重复代码
- [ ] 变量和函数命名是否清晰
- [ ] 注释是否充分和准确
- [ ] 是否遵循 SOLID 原则

#### 安全性检查

- [ ] 是否存在安全漏洞
- [ ] 敏感信息是否正确处理
- [ ] 输入验证是否充分
- [ ] 权限控制是否正确

#### 测试检查

- [ ] 是否有足够的测试覆盖
- [ ] 测试用例是否合理
- [ ] 是否测试了错误场景

### 审查反馈规范

使用以下标签进行反馈：

- `MUST`: 必须修改的问题
- `SHOULD`: 建议修改的问题
- `COULD`: 可选的改进建议
- `QUESTION`: 需要澄清的问题
- `PRAISE`: 值得称赞的代码

示例：

```text
MUST: 这里存在空指针异常的风险，需要添加 nil 检查。

SHOULD: 建议将这个魔法数字提取为常量，提高代码可读性。

COULD: 可以考虑使用更简洁的写法：`return err != nil`

QUESTION: 这个函数的时间复杂度是多少？是否需要优化？

PRAISE: 这个错误处理写得很好，提供了清晰的上下文信息。
```

## 测试规范

### 测试策略

采用测试金字塔模型：

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

### 单元测试

#### Go 单元测试示例

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

### 测试覆盖率要求

- 单元测试覆盖率 ≥ 80%
- 集成测试覆盖率 ≥ 60%
- 关键业务逻辑覆盖率 ≥ 90%

### 运行测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 运行特定包的测试
go test -v ./internal/service/...

# 运行集成测试
make test-integration
```

## 安全规范

### 代码安全

#### 敏感信息处理

**❌ 错误示例 - 硬编码敏感信息**

```go
const (
    DBPassword = "password123"
    APIKey     = "sk-1234567890abcdef"
    JWTSecret  = "my-secret-key"
)
```

**✅ 正确示例 - 使用环境变量**

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

#### 输入验证

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

    // 自定义验证
    if err := ValidatePasswordStrength(r.Password); err != nil {
        return err
    }

    return nil
}

func (r *CreateUserRequest) Sanitize() {
    r.Username = strings.TrimSpace(r.Username)
    r.Email = strings.ToLower(strings.TrimSpace(r.Email))
    // HTML 转义防止 XSS
    r.Username = html.EscapeString(r.Username)
}
```

#### SQL 注入防护

```go
// ❌ 错误示例 - 容易受到 SQL 注入攻击
func GetUserByUsername(username string) (*User, error) {
    query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
    rows, err := db.Query(query)
    // ...
}

// ✅ 正确示例 - 使用参数化查询
func GetUserByUsername(username string) (*User, error) {
    query := "SELECT * FROM users WHERE username = ?"
    rows, err := db.Query(query, username)
    // ...
}

// ✅ 使用 GORM（推荐）
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, errors.Wrap(err, "failed to get user by username")
    }
    return &user, nil
}
```

### 安全检查清单

#### 代码审查安全检查

- [ ] 没有硬编码的密码、密钥或敏感信息
- [ ] 所有用户输入都经过验证和清理
- [ ] 使用参数化查询防止 SQL 注入
- [ ] 实现了适当的认证和授权机制
- [ ] 敏感数据在传输和存储时都进行了加密
- [ ] 实现了适当的错误处理，不泄露敏感信息
- [ ] 使用了安全的随机数生成器
- [ ] 实现了适当的日志记录和审计

## 社区参与

### 报告问题

如果您发现了 bug 或有功能建议，请：

1. 搜索现有的 Issues，避免重复报告
2. 使用合适的 Issue 模板
3. 提供详细的复现步骤和环境信息
4. 如果是安全问题，请私下联系维护者

### 功能请求

1. 在 Issues 中描述您的需求
2. 解释为什么这个功能有用
3. 提供具体的使用场景
4. 考虑向后兼容性

### 文档贡献

- 修复文档中的错误
- 改进现有文档的清晰度
- 添加缺失的文档
- 翻译文档到其他语言

### 社区行为准则

我们致力于为每个人提供友好、安全和欢迎的环境。请：

- 使用友好和包容的语言
- 尊重不同的观点和经验
- 优雅地接受建设性批评
- 关注对社区最有利的事情
- 对其他社区成员表示同理心

## 获取帮助

### 联系方式

- **GitHub Issues**: 报告 bug 和功能请求
- **GitHub Discussions**: 一般讨论和问题

### 资源链接

- [项目文档](./docs/)
- [API 文档](./api-service/docs/)
- [开发规范](./docs/开发规范.md)
- [架构设计](./docs/designs/)

## 致谢

感谢所有为 Websoft9 项目做出贡献的开发者！您的贡献让这个项目变得更好。

---

**Happy Coding! 🚀**
