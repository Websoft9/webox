# Websoft9 平台

现代化的云应用管理解决方案平台，提供应用部署、监控、管理等全生命周期服务。

## 项目概述

Websoft9 是一个以项目为核心组织单元的云应用管理平台，采用分层架构设计，支持多云环境下的应用全生命周期管理。平台通过项目维度实现资源隔离、权限控制、成本管控和团队协作。

### 核心架构组件

- **api-service**: 核心后端服务，基于 Golang + Gin + GORM 构建，提供 RESTful API 接口
- **websoft9-agent**: 部署在服务器节点的客户端代理，通过 gRPC 与服务端通信，负责任务执行和监控数据采集

## 平台功能模块

### 1. 平台主页

- **总览仪表盘**: 项目总览看板、监控总览看板
- **应用快捷导航**: 已部署应用的快速访问入口

### 2. 项目管理 (核心模块)

- **仪表盘**: 监控看板、任务看板、资源看板
- **文件夹**: 项目文件夹、个人文件夹管理
- **应用管理**: 查询、卸载、更新、发布、克隆、配置、迁移
- **工作流**: 画布编排、组件管理、任务调度
- **项目资源**: 资源组、服务器、密钥、数据库、应用网关、证书、云资源
- **项目团队**: 成员管理、角色分配、权限控制
- **项目设置**: 环境变量、项目配置

### 3. 应用市场

- **应用市场列表**: 应用分类、搜索、部署
- **应用心愿单**: 需求提交、投票、悬赏机制

### 4. 平台管理

- **项目管理**: 项目创建、修改、删除
- **平台设置**: 基本设置、系统更新、系统服务、SMTP、短信、Webhook等
- **安全管理**: 角色管理、权限管理、认证管理
- **用户管理**: 用户增删改查
- **告警通知**: 告警查询、规则配置、推送管理
- **个人中心**: 个人资料、个人设置
- **审计日志**: 操作记录、日志查询

## 技术架构

### API Service (后端服务)

基于 Golang + Gin + GORM 构建的核心后端服务。

**技术栈：**

- Golang 1.24+
- Gin Web 框架
- GORM ORM 框架
- SQLite 数据库
- Redis 缓存
- InfluxDB 时序数据库
- JWT 认证
- gRPC 通信

### Websoft9 Agent (客户端代理)

部署在服务器节点的客户端代理，负责任务执行和监控数据采集。

**技术栈：**

- Golang
- gRPC (grpc-go)
- Redis (go-redis)
- InfluxDB (influxdb-client-go)
- Docker Engine API

## 快速开始

### 环境要求

- Go 1.24+
- Docker 20.10+
- SQLite 3.0+ / MySQL 8.0+
- Redis 6.0+
- InfluxDB 2.0+

### 开发环境搭建

1. **克隆仓库**

   ```bash
   git clone <repository-url>
   cd webox
   ```

2. **启动 API 服务**

   ```bash
   cd api-service
   go mod tidy
   make init-db    # 初始化数据库
   make run        # 启动开发服务
   ```

3. **启动 Agent 代理**

   ```bash
   cd websoft9-agent
   go mod tidy
   make build
   sudo ./websoft9-agent
   ```

### Docker 部署

```bash
# 构建并运行 API 服务
cd api-service
docker build -t websoft9/api-service .
docker run -p 8080:8080 -p 9090:9090 websoft9/api-service

# 构建并运行 Agent
cd websoft9-agent
docker build -t websoft9/agent .
docker run --privileged websoft9/agent
```

## 开发工具链

- **API 测试**: [Apifox](https://apifox.com/)
- **CI/CD**: GitHub Actions
- **AI 编程**: [GitHub Copilot](https://github.com/features/copilot), [Claude Code](https://docs.anthropic.com/zh-CN/docs/claude-code/overview)

## 项目结构

```text
webox/
├── api-service/              # 后端 API 服务
│   ├── main.go              # 应用入口文件
│   ├── cmd/                 # 命令行工具
│   │   └── server/         # 服务器启动命令
│   ├── configs/            # 配置文件
│   │   └── config.yaml     # 应用配置
│   ├── internal/           # 内部包
│   │   ├── config/         # 配置管理
│   │   ├── constants/      # 常量定义
│   │   ├── controller/     # 控制器层 (HTTP 处理)
│   │   ├── dto/           # 数据传输对象
│   │   │   ├── request/   # 请求 DTO
│   │   │   └── response/  # 响应 DTO
│   │   ├── interface/     # 接口定义
│   │   │   ├── repository/ # 仓储接口
│   │   │   └── service/   # 服务接口
│   │   ├── middleware/    # 中间件
│   │   ├── model/         # 数据模型
│   │   ├── repository/    # 数据访问层
│   │   ├── router/        # 路由配置
│   │   └── service/       # 业务逻辑层
│   ├── pkg/               # 公共包
│   │   ├── auth/          # 认证相关
│   │   ├── errors/        # 错误处理
│   │   ├── logger/        # 日志处理
│   │   ├── response/      # 响应封装
│   │   ├── utils/         # 工具函数
│   │   └── validator/     # 数据验证
│   ├── scripts/           # 脚本文件
│   └── docs/              # API 文档
└── websoft9-agent/         # 客户端代理
    ├── cmd/                # 命令行入口
    │   └── agent/         # 代理启动命令
    ├── configs/           # 配置文件
    │   └── agent.yaml     # 代理配置
    ├── internal/          # 内部包
    │   ├── agent/         # 代理核心逻辑
    │   ├── communication/ # 通信管理 (gRPC)
    │   ├── config/        # 配置管理
    │   ├── constants/     # 常量定义
    │   ├── monitor/       # 监控数据采集
    │   └── task/          # 任务执行器
    ├── pkg/               # 公共包
    │   └── security/      # 安全验证
    └── scripts/           # 脚本文件
```

## 开发规范

1. 遵循 Go 语言编码规范
2. 使用依赖注入模式
3. 保持接口与实现分离
4. 实现完善的错误处理
5. 添加适当的日志记录
