# Websoft9 Web Service

Websoft9平台的核心后端服务，基于Golang和Gin框架构建，提供RESTful API接口。

## 技术栈

- **开发语言**: Golang 1.24+
- **Web框架**: Gin
- **ORM框架**: GORM
- **数据库**: SQLite
- **缓存**: Redis
- **时序数据库**: InfluxDB
- **认证**: JWT
- **通信**: gRPC

## 项目结构

```text
api-service/
├── main.go                 # 应用入口
├── go.mod                  # Go模块文件
├── Dockerfile              # Docker构建文件
├── configs/                # 配置文件
│   └── config.yaml
├── internal/               # 内部代码
│   ├── config/            # 配置管理
│   ├── database/          # 数据库连接
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层
│   ├── service/           # 业务逻辑层
│   ├── controller/        # 控制器层
│   ├── middleware/        # 中间件
│   └── router/            # 路由配置
└── pkg/                   # 公共包
    ├── auth/              # 认证相关
    ├── response/          # 响应格式
    └── utils/             # 工具函数
```

## 快速开始

### 环境要求

- Go 1.24+
- SQLite 3.0+
- Redis 6.0+
- InfluxDB 2.0+

### 安装依赖

```bash
go mod tidy
```

### 配置

配置文件位于 `configs/config.yaml`，SQLite数据库将自动创建在 `./data/websoft9.db`。

首次运行时会自动创建数据库表结构。

### 运行

```bash
go run main.go
```

### Docker运行

```bash
docker build -t api-service .
docker run -p 8080:8080 -p 9090:9090 api-service
```

## 架构设计

项目采用分层架构设计：

1. **Controller层**: 处理HTTP请求，参数验证
2. **Service层**: 业务逻辑处理
3. **Repository层**: 数据访问抽象
4. **Model层**: 数据模型定义

## 开发规范

- 遵循Go语言编码规范
- 使用依赖注入模式
- 接口与实现分离
- 统一的错误处理和响应格式
- 完善的日志记录

## 部署

### 生产环境部署

1. 构建Docker镜像
2. 配置环境变量
3. 部署到容器平台

### 配置说明

主要配置项：

- `server.port`: HTTP服务端口
- `database.path`: SQLite数据库文件路径
- `redis`: Redis连接配置
- `influxdb`: InfluxDB连接配置
- `jwt.secret`: JWT密钥
- `grpc.port`: gRPC服务端口
