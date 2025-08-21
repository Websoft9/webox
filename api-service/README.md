# Websoft9 API Service

Websoft9 平台的核心后端服务，基于 Golang 和 Gin 框架构建，提供 RESTful API 接口，与 websoft9-agent 通过 gRPC 协议通信。

## 技术栈

- **开发语言**: Go 1.24.5+
- **Web框架**: Gin
- **ORM框架**: GORM
- **数据库**: SQLite (开发环境) / MySQL (生产环境)
- **缓存**: Redis
- **时序数据库**: InfluxDB
- **认证**: JWT
- **通信协议**: gRPC (与 Agent 通信)
- **日志**: Zap
- **配置管理**: Viper

## 项目架构

采用分层架构 + 依赖注入模式：

```
Client → Controller → Service → Repository → Model
                ↓
        Middleware (认证/日志/错误处理)
                ↓
        gRPC Client → Agent (系统操作)
```

### 核心组件

- **Controller层**: HTTP 请求处理和参数验证
- **Service层**: 业务逻辑实现
- **Repository层**: 数据访问抽象
- **Model层**: 数据模型定义
- **Middleware**: 中间件处理 (认证、日志、CORS、错误处理)
- **DTO**: 数据传输对象 (Request/Response)

## 项目结构

```text
api-service/
├── main.go                 # 应用入口
├── go.mod                  # Go模块文件
├── go.sum                  # Go模块依赖锁定文件
├── Dockerfile              # Docker构建文件
├── Makefile                # 构建脚本
├── test_api.sh            # API测试脚本
├── cmd/                   # 命令行应用
│   └── server/            # 服务器入口
├── configs/               # 配置文件
│   └── config.yaml
├── data/                  # 数据存储目录
├── docs/                  # 文档目录
├── scripts/               # 脚本文件
│   ├── init_db.sh        # 数据库初始化脚本
│   ├── init_mysql.sql    # MySQL初始化SQL
│   └── init_sqlite.sql   # SQLite初始化SQL
├── internal/              # 内部代码
│   ├── config/           # 配置管理
│   ├── constants/        # 常量定义
│   ├── controller/       # 控制器层(HTTP处理)
│   ├── dto/              # 数据传输对象
│   │   ├── common.go     # 通用DTO定义
│   │   ├── request/      # 请求DTO
│   │   └── response/     # 响应DTO
│   ├── interface/        # 接口定义
│   │   ├── repository/   # 仓储接口
│   │   └── service/      # 服务接口
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   ├── repository/       # 数据访问层
│   ├── service/          # 业务逻辑层
│   └── router/           # 路由配置
└── pkg/                   # 公共包
    ├── auth/             # JWT认证
    ├── errors/           # 错误处理
    ├── logger/           # 日志组件
    ├── response/         # 统一响应格式
    ├── utils/            # 工具函数
    └── validator/        # 输入验证
```

## 快速开始

### 环境要求

- **Go**: 1.24.5+
- **SQLite**: 3.0+ (开发环境)
- **Redis**: 6.0+ (可选，缓存功能)
- **InfluxDB**: 2.0+ (可选，监控数据)

### 本地开发

1. **安装依赖**
   ```bash
   make deps
   # 或者
   go mod tidy
   ```

2. **初始化数据库**
   ```bash
   make init-db
   ```

3. **运行开发服务器**
   ```bash
   make run
   # 或者
   go run main.go
   ```

4. **开发模式 (热重载)**
   ```bash
   # 需要先安装 air: go install github.com/air-verse/air@latest
   make dev
   ```

### 构建和测试

```bash
# 构建应用
make build

# 运行测试
make test

# 完整测试 (包含API测试)
make full-test

# 代码格式化和检查
make fmt
make vet
```

### Docker 部署

```bash
# 构建Docker镜像
make docker-build

# 运行Docker容器
make docker-run
```

### 配置文件

主要配置文件位于 `configs/config.yaml`：

```yaml
server:
  port: "8080"           # HTTP服务端口
  mode: "debug"          # 运行模式: debug/release

database:
  path: "./data/websoft9.db"  # SQLite数据库路径

redis:
  host: "localhost"      # Redis主机
  port: "6379"          # Redis端口
  password: ""          # Redis密码
  db: 0                 # Redis数据库

influxdb:
  url: "http://localhost:8086"    # InfluxDB地址
  token: "your-influxdb-token"    # InfluxDB访问令牌
  org: "websoft9"                 # 组织名
  bucket: "metrics"               # 存储桶

jwt:
  secret: "websoft9-jwt-secret-key"  # JWT密钥
  expire_time: 3600                  # 过期时间(秒)

grpc:
  port: "9090"          # gRPC服务端口
```

首次运行时会自动创建数据库表结构。

## API 接口

### 认证接口

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/refresh` - 刷新令牌

### 用户管理

- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `PUT /api/v1/users/:id` - 更新用户信息
- `DELETE /api/v1/users/:id` - 删除用户

### 健康检查

- `GET /health` - 服务健康状态

详细的 API 文档请参考项目文档或使用 Swagger UI。

## 开发规范

### 代码组织

项目遵循以下规范：

1. **分层架构**: Controller → Service → Repository → Model
2. **依赖注入**: 使用构造函数注入，接口与实现分离
3. **错误处理**: 统一错误码和错误消息格式
4. **日志记录**: 结构化日志，包含请求上下文
5. **输入验证**: 请求参数验证和业务规则验证

### 文件命名

- 单数名词: `user.go` (不是 `users.go`)
- 接口文件: `internal/interface/{layer}/{domain}.go`
- DTO文件: `internal/dto/{request|response}/{domain}.go`

### 响应格式

统一的 JSON 响应格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {...}
}
```

错误响应格式：

```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": "详细错误信息"
}
```

## 部署指南

### 开发环境

使用本地 SQLite 数据库，快速启动：

```bash
# 克隆项目
git clone <repository-url>
cd api-service

# 安装依赖并初始化
make init

# 启动开发服务器
make dev
```

### 生产环境

#### Docker 部署 (推荐)

```bash
# 1. 构建镜像
docker build -t websoft9/api-service:latest .

# 2. 运行容器
docker run -d \
  --name api-service \
  -p 8080:8080 \
  -p 9090:9090 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/configs:/app/configs \
  websoft9/api-service:latest
```

#### Docker Compose 部署

```yaml
version: '3.8'
services:
  api-service:
    build: .
    ports:
      - "8080:8080"
      - "9090:9090"
    volumes:
      - ./data:/app/data
      - ./configs:/app/configs
    environment:
      - GIN_MODE=release
    depends_on:
      - redis
      - influxdb

  redis:
    image: redis:6-alpine
    ports:
      - "6379:6379"

  influxdb:
    image: influxdb:2.0
    ports:
      - "8086:8086"
    environment:
      - DOCKER_INFLUXDB_INIT_MODE=setup
      - DOCKER_INFLUXDB_INIT_USERNAME=admin
      - DOCKER_INFLUXDB_INIT_PASSWORD=password
```

#### 二进制部署

```bash
# 1. 构建应用
make build

# 2. 配置环境
export GIN_MODE=release

# 3. 启动服务
./api-service
```

### 环境变量

支持通过环境变量覆盖配置：

- `SERVER_PORT`: HTTP服务端口
- `DATABASE_PATH`: 数据库文件路径
- `REDIS_HOST`: Redis主机地址
- `REDIS_PORT`: Redis端口
- `JWT_SECRET`: JWT密钥
- `GIN_MODE`: Gin运行模式 (debug/release)

## 监控和日志

### 日志

日志文件位置：
- 开发环境: 控制台输出
- 生产环境: `/var/log/api-service/`

日志级别：DEBUG, INFO, WARN, ERROR

### 健康检查

```bash
# 服务健康状态
curl http://localhost:8080/health

# 响应示例
{
  "status": "ok",
  "version": "1.0.0",
  "timestamp": "2025-08-21T10:00:00Z",
  "dependencies": {
    "database": "ok",
    "redis": "ok",
    "influxdb": "ok"
  }
}
```

## 故障排查

### 常见问题

1. **数据库连接失败**
   - 检查 SQLite 文件权限
   - 确认 data 目录存在

2. **Redis 连接失败**
   - 检查 Redis 服务状态
   - 验证连接配置

3. **JWT 认证失败**
   - 检查 JWT 密钥配置
   - 确认令牌未过期

4. **gRPC 连接失败**
   - 检查 Agent 服务状态
   - 验证网络连接

### 调试模式

```bash
# 启用详细日志
export GIN_MODE=debug

# 查看数据库查询日志
export DB_LOG_LEVEL=info
```

## 贡献指南

1. Fork 项目
2. 创建特性分支: `git checkout -b feature/new-feature`
3. 提交更改: `git commit -am 'Add new feature'`
4. 推送分支: `git push origin feature/new-feature`
5. 创建 Pull Request

### 代码提交规范

使用约定式提交格式：

```
type(scope): description

feat(user): add user registration API
fix(auth): resolve JWT token validation issue
docs(readme): update installation guide
```

## 许可证

本项目采用 [MIT License](LICENSE) 许可证。
