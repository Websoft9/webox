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
