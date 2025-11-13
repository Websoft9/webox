# GitHub Copilot Instructions for Websoft9

This file provides guidance to GitHub Copilot  when working with code in this repository.

## 项目概述

Websoft9 是一个云应用管理平台，采用分层架构设计，提供完整的应用托管服务。平台遵循项目驱动的架构设计，包含以下核心组件：

### 核心服务架构

- **Websoft9 Gateway Service** (`网关服务层`): 基于 Nginx 的应用网关，提供代理转发、访问控制、SSL证书管理
- **Websoft9 Web Service** (`webox/api-service/`): 平台核心服务，前后端分离架构 (Go + Gin + GORM 后端，Vue 3 + Element Plus 前端)
- **Websoft9 Agent** (`webox/websoft9-agent/`): 部署在各服务器节点的客户端程序，负责任务执行和监控数据采集
- **Websoft9 Storage Service** (`数据存储层`): 提供配置数据存储、缓存数据存储、监控数据存储
- **Docker Runtime** (`基础设施层`): 提供容器化应用的运行支撑

### 平台特点

- **项目驱动**: 以项目为核心的资源组织和权限管理，实现资源隔离和团队协作
- **分层管理**: 平台级管理和项目级管理相结合，满足不同层次的管理需求
- **工作流驱动**: 支持可视化工作流编排，实现复杂业务场景的自动化处理
- **全生命周期**: 覆盖应用从部署到运维的完整生命周期管理
- **安全合规**: 完善的权限控制、审计日志、密钥管理等安全功能
- **多云支持**: 支持公有云、私有云和混合云环境下的统一管理

## 常用开发命令

### API Service (webox/api-service/)

```bash
# 开发流程
cd webox/api-service
make deps          # 下载依赖
make run           # 启动开发服务器 (http://localhost:8080)
make dev           # 热重载开发 (需要 air: go install github.com/air-verse/air@latest)

# 测试和质量控制
make test          # 运行单元测试
make full-test     # 运行 API 集成测试
./test_api.sh      # 运行 API 测试脚本
make fmt           # 代码格式化 (goimports)
make vet           # 代码检查
make lint          # 代码静态分析 (golangci-lint)

# 构建和部署
make init-swag     # 生成 swagger 文档
make build         # 构建二进制文件
make docker-build  # 构建 Docker 镜像
make docker-run    # 运行 Docker 容器
```

### Agent (webox/websoft9-agent/)

```bash
# 开发流程
cd webox/websoft9-agent
make deps          # 下载依赖
make build         # 构建二进制文件
make run           # 运行 agent (需要 root 权限)

# 跨平台构建
make build-linux   # 构建 Linux 版本
make build-all     # 构建所有平台版本

# 测试和质量控制
make test          # 运行测试
make test-coverage # 生成覆盖率报告
make fmt           # 代码格式化
make lint          # 代码分析 (需要 golangci-lint)

# 系统安装
make install       # 安装到 /usr/local/bin/
```

### 项目级命令

```bash
# Git hooks (在 webox/ 目录下运行)
./scripts/pre-commit  # pre-commit 检查: 格式化、lint、测试、安全扫描

# CI/CD 相关
make build-all     # 构建所有组件
make test-all      # 运行全部测试
make security-scan # 安全扫描
```

## 架构和代码结构

### 技术架构设计

Websoft9 采用服务分层、单一职责的设计，每一层对应不同的角色与职责：

```text
用户接入层 (浏览器/移动端)
       ↓
网关服务层 (Nginx 应用网关)
       ↓
平台服务层 (Web Service: Controller → Service → Repository → Model)
       ↓                                      ↗
平台客户端层 (Agent) ←→ 数据存储层 (MySQL/Redis/InfluxDB)
       ↓
基础设施层 (Docker Runtime)
```

### API Service 架构

后端采用 4 层架构设计和依赖注入：

**核心架构原则：**

- **分层架构**: Controller → Service → Repository → Model
- **接口驱动**: Repository 和 Service 层由 `internal/interface/` 中的接口定义
- **依赖注入**: 层间通过构造函数注入实现解耦
- **中间件系统**: 认证、CORS、国际化、错误处理、日志记录
- **统一错误处理**: `pkg/errors/` 中的自定义错误类型
- **结构化日志**: Zap 日志库，支持上下文
- **国际化支持**: 通过 go-i18n 实现多语言支持

**前端架构特点：(未来规划)**

- **Vue 3**: 基于 Composition API 的现代化单页应用
- **Element Plus**: UI 组件库，确保界面一致性和美观性
- **Pinia**: 状态管理，支持模块化的状态组织
- **Vue Router**: 单页应用的路由管理
- **国际化**: 支持多语言界面

### 目录结构

```text
api-service/
├── cmd/server/main.go       # 应用程序入口
├── internal/                # 私有应用程序代码
│   ├── config/              # 配置管理
│   ├── constants/           # 常量定义
│   ├── controller/          # API 路由和处理器
│   ├── dto/                 # 数据传输对象
│   │   ├── request/         # 请求 DTO
│   │   └── response/        # 响应 DTO
│   ├── interface/           # 接口定义
│   │   ├── repository/      # 存储层接口
│   │   └── service/         # 服务层接口
│   ├── middleware/          # 中间件
│   ├── model/               # 数据模型
│   ├── repository/          # 数据访问层
│   ├── router/              # 路由配置
│   └── service/             # 业务逻辑层
├── pkg/                     # 可被外部应用使用的库代码
│   ├── auth/                # 全局认证模块 (JWT)
│   ├── database/            # 数据库模块
│   ├── email/               # 通用邮件模块
│   ├── errors/              # 全局错误处理模块
│   ├── i18n/                # 国际化模块
│   ├── logger/              # 结构化日志模块 (Zap)
│   ├── redis/               # 缓存数据库模块（Redis）
│   ├── utils/               # 通用工具模块
├── scripts/                 # 部署和初始化脚本
├── configs/                 # 配置文件模板
└── docs/                    # API文档目录

websoft9-agent/
├── cmd/agent/main.go        # Agent 入口程序
├── internal/
│   ├── agent/               # Agent 核心逻辑
│   ├── monitor/             # 系统监控模块
│   ├── task/                # 任务执行模块
│   ├── workflow/            # 工作流模块
│   └── communication/       # gRPC 通信模块
├── pkg/security/            # 安全验证
├── proto/                   # gRPC 协议定义
├── configs/                 # 配置文件
└── scripts/                 # 部署脚本
```

### 数据库和存储

**数据存储层架构：**

- **Config DB** (配置数据库):
  - 开发环境: SQLite (自动生成于 `data/websoft9.db`)
  - 生产环境: MySQL 8.0/PostgreSQL (支持 GORM 多数据库)
  - 支持主从复制、读写分离等高可用配置

- **Cache DB** (缓存数据库):
  - Redis: 存储用户会话信息、权限缓存、消息队列等
  - 支持数据持久化配置，确保重要缓存数据可靠性

- **Monitor DB** (监控数据库):
  - InfluxDB 2.x: 存储服务器性能指标、应用监控数据、告警事件等时序数据
  - 支持高效的时间范围查询和数据聚合分析

**数据迁移：**

- 开发环境: GORM AutoMigrate 自动迁移
- 生产环境: 使用专门的迁移脚本
- 所有迁移操作必须可回滚

### 通信架构

**系统间通信：**

- **用户-网关**: HTTP/HTTPS (应用访问)
- **网关-API**: HTTP/HTTPS (RESTful API 调用)
- **前端-API**: HTTP REST + WebSocket (实时通信)
- **API-Agent**: gRPC 通信 (任务指令下发和状态上报)
- **服务端-客户端**: Redis 消息队列 (事件消息交互)

**消息队列类型：**

- 应用运行状态 (Container status)
- 应用健康状态 (App health status)
- 客户端状态 (agent heartbeat)
- 运行时异常 (runtime error)

## 开发指导原则

### 代码规范标准

**Go 后端开发规范：**

1. **命名规范**:
   - 包名：小写，简短，有意义的名词
   - 文件命名：使用单数名词 (`user.go`，不是 `users.go`)
   - 变量名：驼峰命名法，首字母小写
   - 常量名：全大写，下划线分隔
   - 函数名：驼峰命名法，首字母大写（公开）或小写（私有）
   - 结构体：驼峰命名法，首字母大写

2. **包导入**: 内部包使用绝对路径

3. **错误处理**: 始终使用 `pkg/errors` 添加上下文信息，使用标准 error 接口

4. **日志记录**: 使用`pkg/logger`结构化日志 (Zap)，包含必要的上下文信息

5. **魔法值处理**: 多次重复出现的魔法值，始终使用`internal/constants` 添加`const`常量定义，局部出现的魔法值，在go模块头部添加`const`常量定义

6. **国际化（i18n）**: 使用`pkg/i18n`对需要输出给用户的业务日志或信息进行国际化翻译，包括多语言文件`api-service/configs/lang`的翻译和命名统一

7. **测试**: 遵循表驱动测试模式

**Vue 3 前端开发规范：**

1. **组件命名**:
   - 组件文件名使用 PascalCase
   - 组件在模板中使用 kebab-case

2. **Composition API**: 优先使用 `<script setup>` 语法

3. **样式规范**:
   - 使用 BEM 命名方法论
   - 使用 kebab-case
   - CSS 类命名

### 安全要求

**安全设计原则：**

- **最小权限原则**: 所有服务、用户、进程仅授予完成任务所需的最小权限
- **数据加密**: 敏感数据在传输和存储过程中均采用加密措施
- **认证与授权**: 采用基于 RBAC 的权限模型，所有 API 均需身份认证和权限校验
- **安全审计**: 对关键操作和敏感事件进行日志记录
- **输入校验**: 所有外部输入均进行严格校验

**安全功能：**

- JWT 认证和可配置过期时间
- bcrypt 密码加密
- 双因子认证支持 (TOTP)
- 结构体标签和自定义验证器进行输入验证
- RBAC 角色权限控制系统
- gosec pre-commit 安全扫描

### 性能要求

**后端性能指标：**

- API 响应时间 < 200ms (95%)
- 数据库查询优化，避免 N+1 问题
- 合理使用缓存
- 单个函数不超过 100 行
- 圈复杂度不超过 10
- 嵌套层级不超过 4 层

**前端性能指标：**

- 首屏加载时间 < 2s
- 路由懒加载
- 图片懒加载和压缩

### 新增功能流程

**添加新 API 端点时：**

1. **定义 DTO 对象** 在 `internal/dto/request/` 和 `internal/dto/response/`
2. **创建服务接口** 在 `internal/interface/service/`
3. **实现服务层** 在 `internal/service/`
4. **添加仓库接口** (如需) 在 `internal/interface/repository/`
5. **实现仓库层** (如需) 在 `internal/repository/`
6. **创建控制器** 在 `internal/controller/`
7. **注册路由** 在 `internal/router/router.go`
8. **编写测试** 为 service 和 repository 层
9. **API 文档** 使用 `make init-swag` 更新 API 文档

**添加新前端组件时：**

1. **查看现有组件** 了解编写风格和模式
2. **遵循框架约定** 命名规范、类型定义等
3. **创建组件** 在适当的目录下
4. **编写测试** 为组件功能
5. **更新文档** 同步更新相关文档

### 配置管理

**服务配置：**

- **API 服务**: `configs/config.yaml` (基于 Viper 的配置)
- **Agent**: `configs/agent.yaml`
- **网关服务**: Nginx 配置文件
- **前端**: 环境变量和 Vite 配置

**i18n国际化配置：**

- **多语言**: `api-service/configs/lang/<LANGUAGE>.yaml`

**环境支持：**

- **开发环境**: 复制 `config.yaml` 为 `config.local.yaml` 进行本地覆盖
- **Docker 部署**: 支持环境变量覆盖
- **生产环境**: 使用环境变量和配置文件管理

### 测试策略

采用测试金字塔模型：

```text
       /\
      /  \  E2E 测试 (10%)
     /____\
    /      \
   /        \ 集成测试 (20%)
  \________/
  \        /
   \______/ 单元测试 (70%)
```

**测试覆盖率要求：**

- 单元测试覆盖率 ≥ 80%
- 集成测试覆盖率 ≥ 60%
- 关键业务逻辑覆盖率 ≥ 90%

**测试类型：**

- **单元测试**: Service 和 Repository 层，使用 Go 的 table-driven 模式
- **集成测试**: 完整 API 端点测试
- **E2E 测试**: 使用 Playwright 进行端到端测试
- **安全测试**: 使用 `scripts/test_security_*.sh` 专门的安全测试脚本

## 版本控制和协作规范

### Git 工作流

采用 **Git Flow** 工作流模型：

```text
main (生产分支)
├── develop (开发分支)
├── release/v1.2.0 (发布分支)
└── hotfix/critical-bug-fix (修复分支)
```

### 分支命名规范

| 分支类型 | 命名格式 | 示例 |
|----------|----------|------|
| 开发分支 | `develop/版本号` | `develop/v1.2.1` |
| 修复分支 | `hotfix/问题描述` | `hotfix/login-error-handling` |
| 发布分支 | `release/版本号` | `release/v1.2.0` |

### 提交信息规范

使用 **Conventional Commits** 规范：

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**提交类型：**

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

### 代码审查规范

**Pull Request 要求：**

- PR 标题清晰描述变更内容
- 包含详细的变更说明
- 关联相关的 Issue
- 通过所有自动化测试
- 至少一个团队成员审查通过

## 部署架构和环境管理

### 部署架构类型

#### 1. 标准部署架构

按服务角色和资源需求划分为管理服务器和应用服务器：

- **管理服务器**: 部署 Websoft9 平台服务
- **应用服务器**: 部署客户容器应用
- **网关服务**: 提供应用访问控制

#### 2. 最小化部署架构

单机部署方式，仅使用一台服务器部署所有组件，适用于：

- 单一应用的快速部署
- 测试和验证场景
- 开发环境

#### 3. 混合云部署架构

支持跨物理网络机房的部署：

- **公有云和私有云混合部署**
- **多公有云混合部署**
- **网络可访问性**: 保障服务间的网络通信

### 版本号规范

采用 **语义化版本控制（SemVer）**：

```text
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]

例如：
1.0.0        # 正式版本
1.1.0-alpha  # 预发布版本
1.1.0-beta.1 # Beta 版本
1.1.0+20240731 # 带构建信息的版本
```

**版本号递增规则：**

- MAJOR：不兼容的 API 修改
- MINOR：向下兼容的功能性新增
- PATCH：向下兼容的问题修正

### CI/CD 流水线

**GitHub Actions 配置示例：**

- **测试阶段**: 代码格式化、静态分析、单元测试、安全扫描
- **构建阶段**: Docker 镜像构建、镜像仓库推送
- **部署阶段**: 自动化部署到对应环境

## 最佳实践和注意事项

### 开发环境设置

1. **安装依赖工具**:

   ```bash
   # Go 工具链
   go install github.com/air-verse/air@latest        # 热重载
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

   # 前端工具
   npm install -g @vue/cli
   npm install -g vite
   ```

2. **环境变量配置**:

   ```bash
   # 开发环境配置
   export WEBSOFT9_ENV=development
   export DB_TYPE=sqlite
   export REDIS_ENABLED=false
   ```

### 性能优化建议

1. **数据库优化**:
   - 使用数据库连接池
   - 合理设置索引
   - 避免 N+1 查询问题
   - 使用查询缓存

2. **前端优化**:
   - 使用 CDN 加速静态资源
   - 实现代码分割和懒加载
   - 优化图片和静态资源
   - 使用 Gzip 压缩

### 安全最佳实践

1. **敏感信息处理**:
   - 禁止在代码中硬编码敏感信息
   - 使用环境变量或密钥管理服务
   - 定期轮换密钥和证书

2. **API 安全**:
   - 实现请求限流
   - 使用 HTTPS 加密通信
   - 验证所有输入参数
   - 实现安全头配置

### 监控和日志

1. **日志等级使用**:
   - DEBUG：详细的调试信息
   - INFO：一般信息记录
   - WARN：警告信息
   - ERROR：错误信息
   - FATAL：致命错误

2. **监控指标**:
   - 系统资源使用率 (CPU、内存、磁盘、网络)
   - 应用性能指标 (响应时间、QPS、错误率)

### 项目文档

- **产品需求**: `webox/docs/designs/总体方案/产品需求说明书V1.1.md`
- **架构设计**: `webox/docs/designs/架构设计/`
- **详细设计**: `webox/docs/designs/详细设计/`
- **原型设计**: `webox/docs/designs/原型设计/`
- **开发规范**: `webox/docs/开发规范.md`
