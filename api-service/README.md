# Websoft9 API Service

Websoft9云应用管理平台的核心后端API服务，采用分层架构设计，基于Golang和Gin框架构建，提供RESTful API接口，支持与websoft9-agent通过gRPC通信。

## 🚀 技术栈

- **开发语言**: Golang 1.24.5
- **Web框架**: Gin
- **ORM框架**: GORM
- **数据库**: SQLite (开发) / MySQL (生产)
- **缓存**: Redis 6.0+
- **时序数据库**: InfluxDB 2.0+
- **认证**: JWT
- **通信**: gRPC (与Agent通信)
- **国际化**: go-i18n
- **日志**: Zap
- **配置管理**: Viper

## 📁 项目结构

```text
api-service/
├── main.go                 # 应用程序入口点
├── go.mod                  # Go模块依赖管理
├── go.sum                  # Go模块校验和
├── Makefile               # 构建和开发工具
├── Dockerfile             # Docker容器构建文件
├── test_api.sh            # API测试脚本
├── cmd/                   # 命令行工具
│   └── server/           # 服务器启动命令
├── configs/              # 配置文件目录
│   └── config.yaml       # 主配置文件
├── data/                 # 数据存储目录
├── docs/                 # 文档目录
│   └── i18n-standardization.md
├── internal/             # 内部应用代码(不对外暴露)
│   ├── config/          # 配置管理
│   ├── constants/       # 常量定义
│   ├── controller/      # HTTP控制器层
│   │   ├── user.go     # 用户相关API
│   │   └── i18n.go     # 国际化API
│   ├── dto/            # 数据传输对象
│   │   ├── common.go   # 通用DTO
│   │   ├── request/    # 请求DTO
│   │   └── response/   # 响应DTO
│   ├── interface/      # 接口定义(依赖倒置)
│   │   ├── repository/ # 存储层接口
│   │   └── service/    # 服务层接口
│   ├── middleware/     # 中间件
│   │   ├── auth.go    # JWT认证中间件
│   │   ├── cors.go    # 跨域中间件
│   │   ├── error.go   # 错误处理中间件
│   │   ├── i18n.go    # 国际化中间件
│   │   └── logger.go  # 日志中间件
│   ├── model/         # 数据模型定义
│   ├── repository/    # 数据访问层实现
│   ├── router/        # 路由配置
│   └── service/       # 业务逻辑层实现
├── pkg/                  # 可复用的公共包
│   ├── auth/            # JWT认证工具
│   ├── errors/          # 统一错误处理
│   │   ├── codes.go    # 错误码定义
│   │   ├── errors.go   # 错误类型定义
│   │   └── handler.go  # 错误处理器
│   ├── i18n/           # 国际化支持
│   │   └── locales/    # 多语言文件
│   ├── logger/         # 结构化日志
│   ├── response/       # 统一响应格式
│   ├── utils/          # 工具函数
│   └── validator/      # 参数验证
└── scripts/            # 部署和初始化脚本
    ├── init_db.sh     # 数据库初始化
    ├── init_mysql.sql # MySQL初始化脚本
    └── init_sqlite.sql # SQLite初始化脚本
```

## 🚀 快速开始

### 环境要求

- **Go**: 1.24.5 或更高版本
- **SQLite**: 3.0+ (开发环境)
- **Redis**: 6.0+ (可选，用于缓存和会话)
- **InfluxDB**: 2.0+ (可选，用于监控数据)

### 本地开发

#### 1. 克隆项目并安装依赖

```bash
# 进入项目目录
cd api-service

# 安装Go依赖
make deps
# 或者
go mod tidy
```

#### 2. 配置文件

复制并修改配置文件：

```bash
cp configs/config.yaml configs/config.local.yaml
```

主要配置项说明：

- `server.port`: HTTP服务端口 (默认: 8080)
- `server.mode`: 运行模式 (debug/release)
- `database.path`: SQLite数据库文件路径
- `redis.*`: Redis连接配置 (可选)
- `influxdb.*`: InfluxDB连接配置 (可选)
- `jwt.secret`: JWT密钥

#### 3. 初始化数据库

```bash
# 自动创建SQLite数据库和表结构
make run
# 或者使用脚本初始化
./scripts/init_db.sh --init ./scripts/init_data.sql
```

#### 4. 启动服务

```bash
# 开发模式运行
make run

# 或者直接运行
go run main.go

# 使用热重载 (需要先安装air)
make dev
```

#### 5. 验证安装

```bash
# 测试API健康检查
curl http://localhost:8080/health

# 运行API测试
./test_api.sh
```

### 🐳 Docker 开发

#### 构建并运行容器

```bash
# 构建Docker镜像
make docker-build

# 运行容器
make docker-run

# 或者使用docker-compose (包含Redis等依赖服务)
cd ../docker
docker-compose up -d
```

#### Docker环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `SERVER_PORT` | HTTP服务端口 | 8080 |
| `DATABASE_PATH` | 数据库文件路径 | ./data/websoft9.db |
| `REDIS_HOST` | Redis主机地址 | localhost |
| `REDIS_PORT` | Redis端口 | 6379 |
| `JWT_SECRET` | JWT密钥 | websoft9-secret |

## 🏗️ 架构设计

### 分层架构模式

项目采用经典的四层架构设计，通过依赖注入实现层间解耦：

```text
┌─────────────────┐
│   Controller    │ ← HTTP请求处理，参数验证，响应格式化
├─────────────────┤
│    Service      │ ← 业务逻辑处理，数据转换，规则验证
├─────────────────┤
│   Repository    │ ← 数据访问抽象，查询封装，事务管理
├─────────────────┤
│     Model       │ ← 数据模型定义，数据库映射
└─────────────────┘
```

### 核心特性

1. **依赖注入**: 使用接口实现层间解耦，便于测试和扩展
2. **统一错误处理**: 自定义错误类型和中间件统一处理
3. **国际化支持**: 支持中英文多语言切换
4. **结构化日志**: 使用Zap提供高性能日志记录
5. **参数验证**: 基于struct tags的请求参数验证
6. **JWT认证**: 无状态的用户认证和授权
7. **中间件系统**: 可插拔的请求处理中间件

### 通信架构

```text
Client ──HTTP──> API Service ──gRPC──> Websoft9 Agent
   │                  │                       │
   │                  │                       │
   └──WebSocket───────┘                       │
                      │                       │
                   Database              Docker Engine
                   (SQLite)              System Commands
```

## 🔧 开发指南

### 编码规范

1. **文件命名**: 使用单数名词 (`user.go` 而不是 `users.go`)
2. **包导入**: 使用绝对路径导入内部包
3. **错误处理**: 使用 `pkg/errors` 包进行错误包装和处理
4. **日志记录**: 使用结构化日志，包含上下文信息
5. **接口设计**: 遵循接口隔离原则，定义小而专注的接口

### 代码示例

#### 添加新的API端点

1. **定义请求/响应DTO**:

```go
// internal/dto/request/example.go
type CreateExampleRequest struct {
    Name        string `json:"name" binding:"required,min=1,max=100"`
    Description string `json:"description" binding:"max=500"`
}

// internal/dto/response/example.go
type ExampleResponse struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    CreatedAt   string `json:"created_at"`
}
```

2. **定义服务接口**:

```go
// internal/interface/service/example.go
type ExampleService interface {
    Create(ctx context.Context, req *request.CreateExampleRequest) (*response.ExampleResponse, error)
    GetByID(ctx context.Context, id uint) (*response.ExampleResponse, error)
}
```

3. **实现服务逻辑**:

```go
// internal/service/example.go
func (s *exampleService) Create(ctx context.Context, req *request.CreateExampleRequest) (*response.ExampleResponse, error) {
    s.logger.InfoContext(ctx, "创建示例开始", logger.String("name", req.Name))

    // 业务逻辑处理
    example := &model.Example{
        Name:        req.Name,
        Description: req.Description,
    }

    if err := s.exampleRepo.Create(ctx, example); err != nil {
        return nil, errors.WrapError(err, errors.CodeInternalError, "创建示例失败")
    }

    return &response.ExampleResponse{
        ID:          example.ID,
        Name:        example.Name,
        Description: example.Description,
        CreatedAt:   example.CreatedAt.Format(time.RFC3339),
    }, nil
}
```

4. **添加控制器处理**:

```go
// internal/controller/example.go
func (c *ExampleController) Create(ctx *gin.Context) {
    var req request.CreateExampleRequest
    if err := c.bindAndValidateRequest(ctx, &req); err != nil {
        return
    }

    result, err := c.exampleService.Create(ctx.Request.Context(), &req)
    if err != nil {
        errors.HandleError(ctx, err)
        return
    }

    pkg_response.Success(ctx, result)
}
```

### 测试策略

- **单元测试**: 每个service和repository都应有对应的测试
- **集成测试**: 测试API端点的完整流程
- **基准测试**: 性能关键路径的基准测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./internal/service/...

# 生成测试覆盖率报告
make coverage
```

## 📚 API文档

### 认证

API使用JWT Bearer Token进行认证：

```bash
# 获取访问令牌
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# 使用令牌访问受保护的API
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 主要端点

| 方法 | 端点 | 描述 | 认证 |
|------|------|------|------|
| GET | `/health` | 健康检查 | 否 |
| GET | `/api/v1/i18n/messages` | 获取国际化消息 | 否 |
| POST | `/api/v1/auth/login` | 用户登录 | 否 |
| POST | `/api/v1/auth/register` | 用户注册 | 否 |
| GET | `/api/v1/users` | 获取用户列表 | 是 |
| GET | `/api/v1/users/:id` | 获取用户详情 | 是 |
| PUT | `/api/v1/users/:id` | 更新用户信息 | 是 |
| DELETE | `/api/v1/users/:id` | 删除用户 | 是 |

### 响应格式

所有API响应都遵循统一格式：

```json
{
  "code": 200,
  "message": "success",
  "data": { ... },
  "timestamp": "2025-08-21T10:30:00Z"
}
```

错误响应格式：

```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": "字段'name'是必需的",
  "timestamp": "2025-08-21T10:30:00Z"
}
```

## 🚀 部署指南

### 生产环境部署

#### 1. 环境准备

```bash
# 安装Go环境
# 配置MySQL/PostgreSQL数据库
# 配置Redis集群
# 配置InfluxDB集群
```

#### 2. 构建应用

```bash
# 构建生产版本
make build

# 或者构建Docker镜像
make docker-build
```

#### 3. 配置文件

生产环境配置示例：

```yaml
server:
  port: "8080"
  mode: "release"

database:
  driver: "mysql"
  dsn: "user:password@tcp(localhost:3306)/websoft9?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  host: "redis-cluster.example.com"
  port: "6379"
  password: "your-redis-password"

jwt:
  secret: "your-very-secure-jwt-secret-key"
  expire_hours: 24

log:
  level: "info"
  format: "json"
```

#### 4. 使用Docker Compose部署

```yaml
version: '3.8'
services:
  api-service:
    image: websoft9/api-service:latest
    ports:
      - "8080:8080"
    environment:
      - SERVER_MODE=release
      - DATABASE_DSN=mysql://user:pass@db:3306/websoft9
      - REDIS_HOST=redis
      - JWT_SECRET=your-secret-key
    depends_on:
      - db
      - redis
    restart: unless-stopped

  db:
    image: mysql:8.0
    environment:
      MYSQL_DATABASE: websoft9
      MYSQL_USER: websoft9
      MYSQL_PASSWORD: password
      MYSQL_ROOT_PASSWORD: rootpassword
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass your-redis-password
    restart: unless-stopped

volumes:
  mysql_data:
```

### 监控和维护

#### 健康检查

```bash
# 应用健康检查
curl http://localhost:8080/health

# 数据库连接检查
curl http://localhost:8080/health/db

# Redis连接检查
curl http://localhost:8080/health/redis
```

#### 日志管理

```bash
# 查看应用日志
docker logs api-service

# 实时查看日志
docker logs -f api-service

# 查看错误日志
docker logs api-service 2>&1 | grep ERROR
```

## 🛠️ 开发工具

### Makefile 命令

```bash
make build        # 构建应用程序
make run          # 运行开发服务器
make dev          # 热重载开发模式
make test         # 运行测试
make coverage     # 生成测试覆盖率报告
make clean        # 清理构建文件
make deps         # 下载依赖
make docker-build # 构建Docker镜像
make docker-run   # 运行Docker容器
make lint         # 代码静态检查
make fmt          # 格式化代码
```

### 开发环境设置

#### VS Code 推荐扩展

- Go (Google官方)
- REST Client (API测试)
- Docker (容器管理)
- GitLens (Git增强)

#### 调试配置

`.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch API Service",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/main.go",
            "env": {
                "GO_ENV": "development"
            },
            "args": []
        }
    ]
}
```

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

### 提交规范

```bash
# 功能开发
git commit -m "feat: 添加用户管理API"

# 问题修复
git commit -m "fix: 修复用户注册验证问题"

# 文档更新
git commit -m "docs: 更新API文档"

# 代码重构
git commit -m "refactor: 重构用户服务层代码"
```

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](../LICENSE) 文件了解详情。

## 🔗 相关链接

- [Websoft9 官网](https://www.websoft9.com)
- [项目文档](../docs/)
- [问题反馈](https://github.com/Websoft9/webox/issues)
- [更新日志](../CHANGELOG.md)
