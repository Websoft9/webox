# Webox

Websoft9 平台的下一代版本 - 云应用管理解决方案。

## 概述

Webox 是一个现代化的云平台，由两个主要组件构成：

- **websoft9-web-service**: 核心后端服务，提供 RESTful API
- **websoft9-agent**: 部署在服务器节点的客户端代理，负责任务执行和监控

## 组件介绍

### Websoft9 Web Service

基于 Golang 和 Gin 框架构建的核心后端服务。

**技术栈：**
- Golang 1.24+
- Gin Web 框架
- GORM ORM 框架
- SQLite 数据库
- Redis 缓存
- InfluxDB 时序数据库
- JWT 认证
- gRPC 通信

### Websoft9 Agent

部署在服务器节点的客户端代理，负责任务执行和监控。

**技术栈：**
- Golang
- gRPC (grpc-go)
- Redis (go-redis)
- InfluxDB (influxdb-client-go)
- Docker Engine API

## 快速开始

### 环境要求

- Go 1.24+
- Docker
- SQLite 3.0+
- Redis 6.0+
- InfluxDB 2.0+

### 开发环境搭建

1. **克隆仓库**
   ```bash
   git clone <repository-url>
   cd webox
   ```

2. **启动 Web 服务**
   ```bash
   cd websoft9-web-service
   go mod tidy
   go run main.go
   ```

3. **启动代理**
   ```bash
   cd websoft9-agent
   go mod tidy
   go build -o websoft9-agent ./cmd/agent
   sudo ./websoft9-agent
   ```

### Docker 部署

```bash
# 构建并运行 web 服务
cd websoft9-web-service
docker build -t websoft9-web-service .
docker run -p 8080:8080 -p 9090:9090 websoft9-web-service

# 构建并运行代理
cd websoft9-agent
docker build -t websoft9-agent .
docker run --privileged websoft9-agent
```

## 开发工具链

- **API 测试**: [Apifox](https://apifox.com/)
- **CI/CD**: GitHub Actions
- **AI 编程**: [GitHub Copilot](https://github.com/features/copilot), [Claude Code](https://docs.anthropic.com/zh-CN/docs/claude-code/overview)

## 项目结构

```
webox/
├── websoft9-web-service/    # 后端 API 服务
│   ├── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── controller/
│   │   ├── service/
│   │   └── repository/
│   └── pkg/
└── websoft9-agent/          # 客户端代理
    ├── cmd/
    ├── internal/
    └── pkg/
```

## 开发规范

1. 遵循 Go 语言编码规范
2. 使用依赖注入模式
3. 保持接口与实现分离
4. 实现完善的错误处理
5. 添加适当的日志记录
