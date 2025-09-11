# Websoft9 功能详细设计说明书 V2.0 —Git 仓库集成功能


## 2 功能概述

### 2.1 功能定位

Websoft9 使用 Git 仓库存储工作负载所需的部署模板、配置文件和脚本文件，实现集中式的 GitOps 管理能力。Websoft9 不直接实现底层 Git 操作，而是与外部成熟 Git 仓库产品（例如：Gitea）进行深度集成。包括封装一套统一的内外部 API 接口，以及对 Gitea UI 进行细粒度控制以供用户访问。

### 2.2 核心特性

- **Gitea 深度集成**：将 Gitea 作为底层 Git 服务提供者，复用其成熟的仓库管理、版本控制、文件操作等功能
- **用户与权限同步**：平台用户系统与 Gitea 用户/组织/团队的双向同步，确保权限一致性
- **UI 集成与控制**：通过代理中间件嵌入 Gitea Web UI，并基于平台权限控制访问范围
- **仓库分类**：项目成员公用一个项目仓库，用户个人默认有一个私有仓库
- **单点登录**：用户在平台登录后自动获得 Gitea 访问权限

### 2.3 技术架构

#### 2.3.1 整体架构分层

```
┌─────────────────────────────────────────────────────────────┐
│                    用户层 (浏览器/客户端)                     │
└─────────────────────┬───────────────────────────────────────┘
                      │ HTTP/HTTPS
┌─────────────────────▼───────────────────────────────────────┐
│                  网关层 (Nginx)                            │
│              代理转发 | SSL | 限流                          │
└─────────────────────┬───────────────────────────────────────┘
                      │ REST API
┌─────────────────────▼───────────────────────────────────────┐
│                  API层 (Gin框架)                           │
│         路由 | 鉴权 | 校验 | 响应格式 | 审计拦截             │
└─────────────────────┬───────────────────────────────────────┘
                      │ 服务调用
┌─────────────────────▼───────────────────────────────────────┐
│                   服务层 (业务逻辑)                         │
│  空间管理 | 文件操作 | 版本控制 | 权限管理 | 搜索 | 存储      │
└─────────┬───────────────────────────────────────────┬───────┘
          │                                           │
          ▼                                           ▼
┌─────────────────────┐                    ┌──────────────────┐
│    Git 仓库存储      │                    │   MySQL 数据库   │
│  main分支(工作区)    │                    │   元数据/索引    │
│  recycle分支(回收)   │                    │   权限关系       │
│  版本历史记录        │                    │   回收站索引     │
└─────────────────────┘                    └──────────────────┘
```

#### 2.3.2 核心组件与职责

**网关层 (Nginx)**
- 统一入口与反向代理
- SSL/TLS 终结
- 基础安全头与限流

**API层 (Gin + 中间件)**
- RESTful 路由管理
- JWT 认证与权限校验
- 请求参数验证
- 统一响应格式化
- 审计中间件自动记录
- Gitea UI 访问控制中间件（ProxyMiddleware）: 通过平台中间件对嵌入或代理到 Gitea 的 URL 列表进行控制（允许/阻止），并负责会话/身份在反向代理请求中的注入与隔离。

**Gitea 服务（外部/可部署的子系统）**
- 作为正式的 Git 托管与管理系统（仓库、分支、提交、PR、Issues、Web UI）
- 提供 REST API、Web UI，以及钩子（webhooks）用于事件通知
- Gitea 的组织（Organization）用于实现项目空间映射，个人用户用于个人空间

**服务层 (核心业务组件)**
- **Gitea 集成服务**: GiteaIntegrationService（仓库/组织/成员创建、查询、权限同步）
- **权限同步服务**: PermissionSyncService（平台与 Gitea 权限双向同步）
- **用户映射服务**: UserMappingService（平台用户与 Gitea 用户映射管理）
- **SSO 认证服务**: SSOService（单点登录令牌管理）
- **空间管理服务**: SpaceService（空间与 Gitea 仓库/组织映射）
- **UI 代理中间件**: ProxyMiddleware（Gitea UI 访问控制、URL 白名单/黑名单）
- **配额监控服务**: QuotaService（监控 Gitea 仓库配额使用情况）

**存储层 (混合存储策略)**
- **Gitea 仓库**: 实际存储文件内容、历史记录与分支（main/recycle 分支由 Gitea 仓库管理）
- **MySQL**: 元数据索引 + 权限关系 + 回收站索引 + Gitea 仓库映射（例如记录 Gitea repo id、org id、clone 地址、hook 状态）
- **文件系统**: Gitea 的裸仓库物理存储（由 Gitea 进程管理）

#### 2.3.3 关键数据流

**文件上传流程:**
```
用户请求 → 鉴权校验 → 配额检查 → 文件保存 → Git提交(聚合) → 元数据更新 → 响应返回
```

**文件删除流程:**
```
删除请求 → 权限验证 → Git移动到recycle分支 → 更新recycle_index → 返回确认
```

**文件恢复流程:**
```
恢复请求 → 权限验证 → 从recycle分支cherry-pick → 更新索引状态 → 返回结果
```

**版本回滚流程:**
```
回滚请求 → 权限验证 → Git checkout指定版本 → 创建新提交 → 审计记录 → 返回确认
```

#### 2.3.4 核心设计原则

**存储策略**
- Git作为唯一文件版本真实来源，避免重复实现版本控制逻辑
- MySQL仅存储必要的元数据与索引，加速查询性能
- 回收站通过独立recycle分支实现，配合数据库索引优化列表查询

**版本控制策略**
- 自动提交聚合机制(默认20秒)，减少碎片化提交
- 仅对白名单文本文件类型提供差异比较功能
- 保持完整的Git提交历史，确保可追溯性

**权限控制策略**
- 基于RBAC模型的多层权限体系
- 权限验证集中在服务层入口
- 支持权限继承与显式覆盖机制

**简化原则**
- 当前版本不引入缓存、消息队列、高可用等复杂组件
- 单实例部署，降低实现与运维复杂度
- 保留清晰的扩展接入点，支持后续演进

#### 2.3.5 技术栈选择

- **数据库**: MySQL 8.0+ (生产) / SQLite (开发)
- **版本控制**: Git 2.0+ 命令行工具
- **认证**: JWT Token
- **日志**: Zap 结构化日志
- **文档**: Swagger/OpenAPI
- **部署**: Docker + Docker Compose

## 3 功能依赖关系

空间功能作为平台的核心模块，与其他功能特性和基础服务存在依赖关系：

### 3.1 Feature级别依赖（功能特性依赖）
| 依赖Feature | 依赖关系 | 依赖说明 | 影响范围 |
|-------------|----------|----------|----------|
| **Gitea 服务** | 强依赖 | Gitea 提供仓库管理、版本控制、Web 界面和 API。平台需要与其对接并保证高可用 | 所有与仓库/版本控制/UI 集成的功能 |
| **用户管理Feature** | 强依赖 | 需要用户身份验证、用户信息查询，以及与 Gitea 的用户/组织/团队同步 | 全部功能的用户身份与权限验证 |
| **项目管理Feature** | 强依赖 | 项目空间归属于具体项目，需要项目信息 | 项目空间的创建和归属管理 |
| **权限管理Feature** | 强依赖 | 文件和空间的权限控制 | 文件访问控制、操作权限验证 |
| **工作流Feature** | 弱依赖 | 工作流文件和模板的存储管理 | 工作流文件管理功能 |
| **应用管理Feature** | 弱依赖 | 应用配置文件和部署模板存储 | 应用相关文件管理 |

### 3.2 基础服务依赖（Infrastructure依赖）

| 依赖服务 | 依赖关系 | 依赖说明 | 技术要求 |
|----------|----------|----------|----------|
| **认证服务** | 强依赖 | 用户身份验证和会话管理 | JWT Token认证 |
| **权限服务** | 强依赖 | RBAC权限验证和角色管理 | 权限验证中间件 |
| **审计中间件** | 强依赖 | 通过统一中间件记录文件操作 | 异步审计记录 |
| **通知服务** | 弱依赖 | 文件操作通知、空间变更通知 | 邮件/短信通知 |

### 3.3 基础设施依赖（Infrastructure依赖）

| 基础设施 | 依赖关系 | 依赖说明 | 技术要求 |
|----------|----------|----------|----------|
| **Git服务** | 强依赖 | 文件版本控制、存储管理、内容搜索 | Git 2.0+ |
| **文件系统** | 强依赖 | Git仓库的物理存储载体 | 支持大文件、权限控制 |
| **数据库服务** | 强依赖 | 元数据存储 | MySQL 8.0+ |

### 3.4 依赖接口清单

- **认证服务接口：** POST /api/v1/auth/validate # 用户验证
- **权限服务接口：** POST /api/v1/permissions/check # 权限检查
- **通知服务接口：** POST /api/v1/notifications/send  # 发送通知
- **用户管理接口：** GET /api/v1/users/{id} # 获取用户信息
- **项目管理接口：** GET /api/v1/projects/{id} # 获取项目信息

**审计记录说明**：所有API操作自动通过项目统一的审计中间件记录，无需业务代码主动调用审计接口。

### 3.5 开发优先级建议

**Phase 0: 基础服务准备（开发前置条件）**

- 认证服务（用户身份验证）
- 权限服务（基础权限验证）
- 审计服务（操作日志记录）
- 数据库服务（元数据存储）

**Phase 1: 核心功能开发**

- 项目空间基础功能开发
- 与认证、权限服务的集成
- 基础文件操作功能

**Phase 2: Feature间集成**

- 与用户管理的数据集成
- 与项目管理的关联
- 工作流和应用管理的文件存储集成


### 3.6 风险评估与应对

**高风险依赖：**
- **认证服务不可用**：项目空间完全无法使用
  - 应对：实现认证服务的高可用部署，设计降级策略
- **权限服务故障**：可能导致权限控制失效
  - 应对：实现权限缓存机制，设置默认安全策略

**中风险依赖：**
- **用户管理变更**：可能影响用户信息获取
  - 应对：定义稳定的数据接口，使用数据适配层
- **项目管理调整**：可能影响项目关联逻辑
  - 应对：通过事件驱动架构解耦，避免直接依赖

**低风险依赖：**
- **通知服务异常**：仅影响通知功能，不影响核心业务
  - 应对：异步处理，允许通知失败

## 4 Business（业务目标与用户故事）

### 4.1 业务目标概述

- 基于 Gitea 的成熟 Git 平台，实现空间管理与权限控制的统一
- 通过用户权限同步，实现平台与 Gitea 的无缝集成
- 通过 UI 代理实现 Gitea 界面的安全嵌入与访问控制

### 4.2 核心用户故事拆分

1. **US001 - 用户权限同步**
   - 作为平台管理员，我希望平台用户能自动映射到 Gitea 用户，无需重复创建账户
   - 作为项目管理员，我希望项目成员权限能自动同步到 Gitea 组织权限，保持一致性
   - 作为用户，我希望在平台登录后能直接访问 Gitea 界面，无需重复登录

2. **US002 - 空间映射管理**
   - 作为项目管理员，我希望创建项目空间时自动在 Gitea 创建对应组织和仓库
   - 作为用户，我希望个人空间能映射到 Gitea 个人仓库，享受完整的 Git 功能
   - 作为管理员，我希望删除空间时能安全清理对应的 Gitea 资源

3. **US003 - UI 集成与访问控制**
   - 作为用户，我希望在平台内直接访问 Gitea 界面，体验无缝集成
   - 作为管理员，我希望控制用户在 Gitea UI 中的访问范围，确保安全性
   - 作为项目成员，我希望只能访问有权限的仓库界面，其他内容被隐藏

## 5 Module（模块拆分与职责）

### 5.1 模块架构设计

```
空间 Feature 集成架构
├── Gitea 集成模块（Gitea Integration）
├── 用户权限同步模块（User & Permission Sync）
├── UI 代理与访问控制模块（UI Proxy & Access Control）
├── 空间映射管理模块（Space Mapping）
└── 单点登录模块（SSO Authentication）
```

### 5.2 模块功能映射

| 模块名称 | 核心功能 | 业务场景 | 技术实现要点 |
|----------|----------|----------|-------------|
| Gitea 集成模块 | Gitea API 调用、Webhooks 处理 | 与 Gitea 通信的所有场景 | HTTP 客户端、API 封装、错误处理 |
| 用户权限同步模块 | 平台与 Gitea 用户权限同步 | 用户管理、权限变更同步 | 双向同步、权限映射、冲突解决 |
| UI 代理与访问控制模块 | Gitea UI 代理、URL 访问控制 | Gitea 界面嵌入、安全控制 | 反向代理、URL 过滤、会话注入 |
| 空间映射管理模块 | 空间与 Gitea 资源映射 | 空间创建、删除、迁移 | 资源映射、生命周期管理 |
| 单点登录模块 | SSO 令牌管理、身份传递 | 用户登录、会话管理 | JWT 处理、令牌生成、身份验证 |

### 5.3 详细模块拆分

#### 5.3.1 Gitea 集成模块（Gitea Integration）
**核心职责**：负责与 Gitea 的所有 API 交互、资源管理、Webhooks 处理

**详细设计**：
- **Gitea API 客户端**：
  - 仓库管理：创建、删除、更新仓库信息
  - 组织管理：创建、删除组织，管理组织成员
  - 用户管理：创建、删除 Gitea 用户，管理用户信息
  - 权限管理：设置仓库权限、团队权限、成员权限
  
- **Webhooks 处理**：
  - 注册 Webhooks：为仓库注册推送、PR 等事件的 Webhooks
  - 事件处理：接收并处理 Gitea 发送的 Webhooks 事件
  - 事件转换：将 Gitea 事件转换为平台内部事件

- **资源同步**：
  - 仓库信息同步：定期同步仓库元数据到平台数据库
  - 状态监控：监控 Gitea 服务状态和仓库健康状况

**交付物**：
- Service：GiteaAPIClient, GiteaWebhookService, GiteaResourceSyncService
- Models：GiteaRepo, GiteaOrg, GiteaUser
- Handlers：WebhookHandler

#### 5.3.2 用户权限同步模块（User & Permission Sync）
**核心职责**：管理平台用户与 Gitea 用户的映射关系，确保权限一致性

**详细设计**：
- **用户映射管理**：
  - 用户创建同步：平台创建用户时自动在 Gitea 创建对应用户
  - 用户信息同步：双向同步用户基本信息（用户名、邮箱、头像等）
  - 用户状态同步：同步用户启用/禁用状态

- **权限映射策略**：
  - 平台角色到 Gitea 权限映射：
    - 平台管理员 → Gitea 管理员
    - 项目 Owner → Gitea 组织 Owner
    - 项目成员 → Gitea 组织成员/团队成员
  - 细粒度权限控制：将平台的 RBAC 权限转换为 Gitea 的仓库权限

- **权限同步机制**：
  - 增量同步：监听权限变更事件，实时同步到 Gitea
  - 全量同步：定期执行全量权限校验和同步
  - 冲突解决：处理平台与 Gitea 权限不一致的情况

**交付物**：
- Service：UserMappingService, PermissionSyncService, PermissionConflictResolver
- Models：UserMapping, PermissionMapping
- Jobs：PermissionSyncJob, UserSyncJob

#### 5.3.3 UI 代理与访问控制模块（UI Proxy & Access Control）
**核心职责**：提供 Gitea Web UI 的安全代理访问，实现细粒度的 URL 访问控制

**详细设计**：
- **反向代理中间件**：
  - 请求拦截：拦截所有到 Gitea UI 的请求
  - 身份注入：将平台用户身份转换为 Gitea 身份并注入请求
  - 响应处理：处理 Gitea 返回的响应，进行必要的内容过滤

- **URL 访问控制**：
  - 白名单机制：配置允许访问的 Gitea URL 模式
  - 黑名单机制：配置禁止访问的敏感 URL
  - 动态权限检查：基于用户权限动态允许/拒绝访问特定 URL

- **会话管理**：
  - 会话同步：将平台会话与 Gitea 会话关联
  - 自动登录：用户访问 Gitea UI 时自动完成登录
  - 会话失效处理：处理会话过期和刷新

**交付物**：
- Middleware：GiteaProxyMiddleware, URLAccessControlMiddleware, SessionSyncMiddleware
- Service：URLFilterService, SessionMappingService
- Configuration：URLAllowList, URLBlockList

#### 5.3.4 空间映射管理模块（Space Mapping）
**核心职责**：管理平台空间与 Gitea 资源的映射关系

**详细设计**：
- **空间类型映射**：
  - 项目空间 → Gitea 组织 + 仓库
  - 个人空间 → Gitea 个人仓库
  - 映射关系持久化：在数据库中维护映射关系表

- **资源生命周期管理**：
  - 创建映射：空间创建时自动创建对应的 Gitea 资源
  - 删除映射：空间删除时安全清理 Gitea 资源
  - 更新映射：空间信息变更时同步更新 Gitea 资源

- **映射验证**：
  - 一致性检查：定期验证映射关系的有效性
  - 孤儿资源清理：清理无效的映射关系和孤儿资源
  - 映射修复：自动修复损坏的映射关系

**交付物**：
- Service：SpaceMappingService, ResourceLifecycleService, MappingValidationService
- Models：SpaceMapping, MappingValidation
- Jobs：MappingValidationJob, OrphanResourceCleanupJob

#### 5.3.5 单点登录模块（SSO Authentication）
**核心职责**：实现平台与 Gitea 的单点登录集成

**详细设计**：
- **SSO 策略选择**：
  - OAuth2 集成：配置 Gitea 作为 OAuth2 提供方或消费方
  - JWT 令牌传递：通过 JWT 令牌在平台和 Gitea 间传递身份信息
  - Header 注入：在代理请求中注入认证 Header

- **令牌管理**：
  - 令牌生成：为平台用户生成访问 Gitea 的临时令牌
  - 令牌刷新：自动刷新即将过期的令牌
  - 令牌撤销：用户登出时撤销相关令牌

- **认证流程**：
  - 登录重定向：用户访问 Gitea 时重定向到平台登录
  - 自动认证：已登录用户自动获得 Gitea 访问权限
  - 权限验证：验证用户是否有权访问特定 Gitea 资源

**交付物**：
- Service：SSOService, TokenService, AuthFlowService
- Models：SSOToken, AuthSession
- Middleware：SSOAuthMiddleware

### 5.4 所需中间件清单

#### 5.4.1 核心中间件

| 中间件名称 | 功能描述 | 应用场景 | 技术要求 |
|-----------|----------|----------|----------|
| GiteaProxyMiddleware | Gitea UI 反向代理 | 所有 Gitea UI 访问 | Gin 中间件、HTTP 代理 |
| URLAccessControlMiddleware | URL 访问控制 | Gitea UI 安全控制 | 正则匹配、权限验证 |
| SessionSyncMiddleware | 会话同步 | 平台与 Gitea 会话管理 | Redis 会话存储 |
| SSOAuthMiddleware | 单点登录认证 | 用户身份验证 | JWT 处理、OAuth2 |
| PermissionSyncMiddleware | 权限同步 | 权限变更时触发同步 | 异步任务队列 |

#### 5.4.2 辅助中间件

| 中间件名称 | 功能描述 | 应用场景 | 技术要求 |
|-----------|----------|----------|----------|
| GiteaAPIMiddleware | Gitea API 请求增强 | 所有 Gitea API 调用 | 重试机制、错误处理 |
| WebhookValidationMiddleware | Webhook 签名验证 | Gitea Webhook 接收 | HMAC 签名验证 |
| ResourceMappingMiddleware | 资源映射检查 | 空间操作前验证 | 数据库查询、缓存 |
| ConfigurationMiddleware | 动态配置管理 | URL 过滤规则更新 | 配置热重载 |

### 5.5 接口设计

| 接口名称 | 请求方式 | 请求路径 | 描述 |
|----------|----------|----------|------|
| 创建空间映射 | POST | /api/v1/spaces/{id}/mapping | 创建空间与 Gitea 资源的映射 |
| 获取空间映射 | GET | /api/v1/spaces/{id}/mapping | 获取空间的 Gitea 映射信息 |
| 删除空间映射 | DELETE | /api/v1/spaces/{id}/mapping | 删除空间映射并清理 Gitea 资源 |
| 用户权限同步 | POST | /api/v1/sync/permissions | 手动触发用户权限同步 |
| 获取同步状态 | GET | /api/v1/sync/status | 获取权限同步状态 |
| 配置 URL 访问控制 | PUT | /api/v1/gitea/access-control | 配置 Gitea UI 访问控制规则 |
| 获取 Gitea 代理状态 | GET | /api/v1/gitea/proxy/status | 获取 Gitea 代理服务状态 |
| SSO 令牌生成 | POST | /api/v1/sso/token | 为用户生成 Gitea 访问令牌 |
| 获取 Webhook 配置 | GET | /api/v1/gitea/webhooks | 获取 Gitea Webhook 配置 |
| 更新 Webhook 配置 | PUT | /api/v1/gitea/webhooks | 更新 Gitea Webhook 配置 |

## 6 API（关键接口与契约）

### 6.1 接口分层设计

```
External APIs（对外接口）
├── Space Mapping APIs
├── User Permission Sync APIs
├── UI Proxy Control APIs
├── SSO Token APIs
└── Gitea Integration APIs

Internal APIs（内部接口）
├── Gitea API Client
├── Permission Mapping
├── Session Management
└── Resource Lifecycle
```

### 6.2 核心 API 设计

#### 1. 空间映射管理 APIs
```
POST   /api/v1/spaces/{id}/mapping          # 创建空间与 Gitea 资源映射
GET    /api/v1/spaces/{id}/mapping          # 获取空间映射信息
DELETE /api/v1/spaces/{id}/mapping          # 删除空间映射
PUT    /api/v1/spaces/{id}/mapping/sync     # 手动同步映射状态
```

#### 2. 用户权限同步 APIs
```
POST   /api/v1/sync/users                   # 同步用户到 Gitea
POST   /api/v1/sync/permissions             # 同步权限到 Gitea
GET    /api/v1/sync/status                  # 获取同步状态
POST   /api/v1/sync/validate                # 验证同步一致性
```

#### 3. UI 代理控制 APIs
```
GET    /api/v1/gitea/proxy/config           # 获取代理配置
PUT    /api/v1/gitea/proxy/config           # 更新代理配置
POST   /api/v1/gitea/proxy/allowlist        # 配置允许访问的 URL
POST   /api/v1/gitea/proxy/blocklist         # 配置禁止访问的 URL
GET    /api/v1/gitea/proxy/status           # 获取代理状态
```

#### 4. SSO 认证 APIs
```
POST   /api/v1/sso/token                    # 生成 Gitea 访问令牌
GET    /api/v1/sso/session                  # 获取会话状态
DELETE /api/v1/sso/session                  # 销毁会话
POST   /api/v1/sso/refresh                  # 刷新令牌
```

#### 5. Gitea 集成 APIs
```
POST   /api/v1/gitea/repos                  # 在 Gitea 创建仓库
GET    /api/v1/gitea/repos                  # 获取 Gitea 仓库列表
DELETE /api/v1/gitea/repos/{id}             # 删除 Gitea 仓库
POST   /api/v1/gitea/orgs                   # 创建 Gitea 组织
GET    /api/v1/gitea/orgs/{id}/members      # 获取组织成员
POST   /api/v1/gitea/webhooks               # 配置 Webhook
```

### 6.3 中间件集成要点

#### 1. GiteaProxyMiddleware 实现要点

```go
type GiteaProxyMiddleware struct {
    giteaURL     string
    allowList    []string
    blockList    []string
    defaultPolicy string
    ssoService   SSOService
}

func (m *GiteaProxyMiddleware) Handle(c *gin.Context) {
    // 1. 检查 URL 访问控制
    if !m.isURLAllowed(c.Request.URL.Path) {
        c.JSON(403, gin.H{"error": "Access denied"})
        return
    }
    
    // 2. 注入 SSO 认证信息
    token, err := m.ssoService.GetGiteaToken(c)
    if err != nil {
        c.JSON(401, gin.H{"error": "Authentication required"})
        return
    }
    
    // 3. 转发请求到 Gitea
    m.proxyToGitea(c, token)
}
```

#### 2. 权限同步中间件要点

```go
type PermissionSyncMiddleware struct {
    syncService PermissionSyncService
    queue       AsyncQueue
}

func (m *PermissionSyncMiddleware) Handle(c *gin.Context) {
    // 监听权限变更操作
    if isPermissionOperation(c.Request.URL.Path) {
        // 异步触发权限同步
        m.queue.Push(SyncTask{
            Type: "permission_sync",
            SpaceID: getSpaceID(c),
            UserID: getUserID(c),
        })
    }
    c.Next()
}
```

### 6.4 API 设计规范（Gitea 集成版）

- **响应格式**：统一使用 `{code, message, data}` 结构
- **权限验证**：所有 API 都需要通过权限中间件验证
- **Gitea 集成**：底层操作通过 Gitea API 实现，平台 API 提供业务层封装
- **同步机制**：权限和用户变更自动触发与 Gitea 的同步
- **错误处理**：包含 Gitea API 调用失败的错误处理和重试机制

#### 1. 获取空间文件列表

**GET /api/v1/spaces/{id}/files**

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述         |
| --------- | ------- | ---- | ------ | ------------ |
| path      | string  | 否   | /      | 目录路径     |
| page      | integer | 否   | 1      | 页码         |
| page_size | integer | 否   | 50     | 每页数量     |
| file_type | string  | 否   | -      | 文件类型筛选 |
| keyword   | string  | 否   | -      | 搜索关键词   |
| sort      | string  | 否   | name   | 排序字段     |
| order     | string  | 否   | asc    | 排序方向     |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "current_path": "/documents",
    "parent_path": "/",
    "items": [
      {
        "id": 1,
        "name": "项目文档.pdf",
        "path": "/documents/项目文档.pdf",
        "type": "FILE",
        "size": 2048576,
        "mime_type": "application/pdf",
        "download_count": 5,
        "parent_id": null,
        "storage_path": "/storage/space_1/documents/项目文档.pdf",
        "checksum": "sha256:abc123...",
        "version_tracked": true,
        "current_version": "v1.2",
        "permissions": {
          "read": true,
          "write": true,
          "delete": false
        },
        "created_by": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员"
        },
        "updated_by": {
          "id": 2,
          "username": "editor",
          "nickname": "编辑者"
        },
        "created_at": "2025-01-15T10:30:00Z",
        "updated_at": "2025-01-22T10:30:00Z"
      },
      {
        "id": 2,
        "name": "images",
        "path": "/documents/images",
        "type": "DIRECTORY",
        "size": 0,
        "file_count": 15,
        "subfolder_count": 2,
        "parent_id": null,
        "storage_path": "/storage/space_1/documents/images",
        "permissions": {
          "read": true,
          "write": true,
          "delete": true
        },
        "created_by": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员"
        },
        "created_at": "2025-01-10T10:30:00Z",
        "updated_at": "2025-01-20T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 50,
      "total": 2,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    },
    "storage_info": {
      "used_space": 104857600,
      "total_space": 1073741824,
      "usage_percentage": 9.77,
      "file_count": 125,
      "folder_count": 8
    },
    "breadcrumb": [
      {"name": "根目录", "path": "/"},
      {"name": "documents", "path": "/documents"}
    ]
  }
}
```

#### 2. 上传文件到空间

**POST /api/v1/spaces/{id}/files**

请求体（multipart/form-data）：

| 参数名      | 类型    | 必填 | 默认值 | 描述               |
| ----------- | ------- | ---- | ------ | ------------------ |
| file        | file    | 是   | -      | 文件内容           |
| path        | string  | 否   | -      | 目标路径，默认根目录 |
| filename    | string  | 否   | -      | 自定义文件名       |
| description | string  | 否   | -      | 文件描述           |
| overwrite   | boolean | 否   | -      | 是否覆盖同名文件   |

响应示例：

```json
{
  "code": 200,
  "message": "文件上传成功",
  "data": {
    "id": 123,
    "name": "新文档.pdf",
    "path": "/documents/新文档.pdf",
    "size": 2048576,
    "mime_type": "application/pdf",
    "checksum": "sha256:def456...",
    "upload_progress": 100,
    "version_tracked": true,
    "current_version": "v1.0",
    "created_at": "2025-01-22T10:35:00Z"
  }
}
```

#### 3. 文件搜索

**GET /api/v1/spaces/{id}/search**

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述                                   |
| --------- | ------- | ---- | ------ | -------------------------------------- |
| q         | string  | 是   | -      | 搜索关键词                             |
| type      | string  | 否   | all    | 搜索类型（filename, content, all）     |
| file_type | string  | 否   | -      | 文件类型筛选                           |
| path      | string  | 否   | -      | 搜索路径范围                           |
| page      | integer | 否   | 1      | 页码                                   |
| page_size | integer | 否   | 20     | 每页数量                               |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "keyword": "项目",
    "total_results": 15,
    "search_time": 150,
    "items": [
      {
        "id": 1,
        "name": "项目文档.pdf",
        "path": "/documents/项目文档.pdf",
        "type": "FILE",
        "size": 2048576,
        "match_type": "filename",
        "highlight": "<mark>项目</mark>文档.pdf",
        "content_preview": "这是一个关于<mark>项目</mark>的重要文档...",
        "relevance_score": 0.95,
        "updated_at": "2025-01-22T10:30:00Z"
      }
    ],
    "suggestions": ["项目文档", "项目计划", "项目总结"],
    "filters": {
      "file_types": [
        {"type": "pdf", "count": 8},
        {"type": "docx", "count": 5},
        {"type": "txt", "count": 2}
      ],
      "paths": [
        {"path": "/documents", "count": 10},
        {"path": "/projects", "count": 5}
      ]
    }
  }
}
```

#### 4. 批量文件操作

**POST /api/v1/spaces/{id}/files/batch**

请求体参数：

| 参数名    | 数据类型  | 是否可空 | 描述                                          |
| --------- | --------- | -------- | --------------------------------------------- |
| action    | string    | 否       | 操作类型（delete, move, copy, download）      |
| file_paths | string[]  | 否       | 文件路径数组                                  |
| target_path | string    | 是       | 目标路径（move/copy操作必填）                 |
| force     | boolean   | 是       | 是否强制操作                                  |

响应示例：

```json
{
  "code": 200,
  "message": "批量操作已启动",
  "data": {
    "batch_id": "batch_123456789",
    "action": "move",
    "total_files": 5,
    "status": "PROCESSING",
    "progress": 20,
    "results": [
      {
        "file_path": "/documents/file1.pdf",
        "status": "SUCCESS",
        "message": "移动成功"
      },
      {
        "file_path": "/documents/file2.pdf",
        "status": "PROCESSING",
        "message": "正在移动..."
      },
      {
        "file_path": "/documents/file3.pdf",
        "status": "FAILED",
        "message": "权限不足",
        "error_code": "PERMISSION_DENIED"
      }
    ],
    "estimated_time": 30,
    "started_at": "2025-01-22T10:35:00Z"
  }
}
```

#### 5. 获取文件版本历史

**GET /api/v1/spaces/{id}/files/{fileId}/versions**

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述     |
| --------- | ------- | ---- | ------ | -------- |
| page      | integer | 否   | 1      | 页码     |
| page_size | integer | 否   | 20     | 每页数量 |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "file_info": {
      "id": 1,
      "name": "项目文档.pdf",
      "path": "/documents/项目文档.pdf",
      "current_version": "v1.3"
    },
    "versions": [
      {
        "id": 3,
        "version": "v1.3",
        "version_hash": "abc123def456",
        "commit_message": "更新项目进度",
        "size": 2048576,
        "created_by": {
          "id": 2,
          "username": "editor",
          "nickname": "编辑者"
        },
        "created_at": "2025-01-22T10:30:00Z",
        "is_current": true,
        "download_url": "/api/v1/spaces/1/files/1/versions/3/download"
      },
      {
        "id": 2,
        "version": "v1.2",
        "version_hash": "def456ghi789",
        "commit_message": "修复文档错误",
        "size": 2035467,
        "created_by": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员"
        },
        "created_at": "2025-01-20T15:20:00Z",
        "is_current": false,
        "download_url": "/api/v1/spaces/1/files/1/versions/2/download"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 3,
      "total_pages": 1
    },
    "statistics": {
      "total_versions": 3,
      "size_trend": "increasing",
      "update_frequency": "weekly"
    }
  }
}
```

#### 6. 通用文件上传

**POST /api/v1/upload**

请求体（multipart/form-data）：

| 参数名        | 类型    | 必填 | 默认值 | 描述                      |
| ------------- | ------- | ---- | ------ | ------------------------- |
| file          | file    | 是   | -      | 文件内容                  |
| type          | string  | 否   | other  | 上传类型                  |
| max_size      | integer | 否   | -      | 最大文件大小（字节）      |
| allowed_types | string  | 否   | -      | 允许的文件类型（逗号分隔） |

响应示例：

```json
{
  "code": 200,
  "message": "文件上传成功",
  "data": {
    "file_id": "file_123456789",
    "filename": "document.pdf",
    "original_name": "项目文档.pdf",
    "size": 2048576,
    "mime_type": "application/pdf",
    "url": "https://cdn.websoft9.com/files/file_123456789.pdf",
    "thumbnail_url": "https://cdn.websoft9.com/thumbnails/file_123456789.jpg",
    "checksum": "sha256:abc123def456",
    "uploaded_at": "2025-01-22T10:35:00Z",
    "expires_at": "2025-01-22T11:35:00Z"
  }
}
```

#### 7. 下载文件

**GET /api/v1/spaces/{id}/files/download**

查询参数：

| 参数 | 类型   | 必填 | 默认值 | 描述     |
| ---- | ------ | ---- | ------ | -------- |
| path | string | 是   | -      | 文件路径 |

响应：
- Content-Type: application/octet-stream
- Content-Disposition: attachment; filename="filename.pdf"
- 文件二进制流

#### 8. 删除文件

**DELETE /api/v1/spaces/{id}/files**

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述                      |
| ------ | -------- | -------- | ------------------------- |
| paths  | string[] | 否       | 要删除的文件/目录路径数组 |

响应示例：

```json
{
  "code": 200,
  "message": "文件删除成功",
  "data": {
    "deleted_count": 3,
    "moved_to_recycle_bin": true,
    "results": [
      {
        "path": "/documents/file1.pdf",
        "status": "SUCCESS",
        "message": "已移入回收站"
      },
      {
        "path": "/documents/file2.pdf",
        "status": "FAILED",
        "message": "权限不足",
        "error_code": "PERMISSION_DENIED"
      }
    ]
  }
}
```

#### 9. 创建目录

**POST /api/v1/spaces/{id}/directories**

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述     |
| ----------- | -------- | -------- | -------- |
| path        | string   | 否       | 目录路径 |
| description | string   | 是       | 目录描述 |

响应示例：

```json
{
  "code": 200,
  "message": "目录创建成功",
  "data": {
    "id": 456,
    "name": "新目录",
    "path": "/documents/新目录",
    "type": "DIRECTORY",
    "permissions": {
      "read": true,
      "write": true,
      "delete": true
    },
    "created_at": "2025-01-22T10:35:00Z"
  }
}
```

#### 10. 移动文件

**PUT /api/v1/spaces/{id}/files/move**

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述         |
| ----------- | -------- | -------- | ------------ |
| source_path | string   | 否       | 源文件路径   |
| target_path | string   | 否       | 目标文件路径 |

响应示例：

```json
{
  "code": 200,
  "message": "文件移动成功",
  "data": {
    "source_path": "/documents/old_file.pdf",
    "target_path": "/archive/old_file.pdf",
    "file_id": 123,
    "moved_at": "2025-01-22T10:35:00Z"
  }
}
```

#### 11. 复制文件

**POST /api/v1/spaces/{id}/files/copy**

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述         |
| ----------- | -------- | -------- | ------------ |
| source_path | string   | 否       | 源文件路径   |
| target_path | string   | 否       | 目标文件路径 |

响应示例：

```json
{
  "code": 200,
  "message": "文件复制成功",
  "data": {
    "source_file_id": 123,
    "new_file_id": 456,
    "source_path": "/documents/file.pdf",
    "target_path": "/backup/file_copy.pdf",
    "copied_at": "2025-01-22T10:35:00Z"
  }
}
```

#### 12. 回收站操作

**获取回收站文件列表**

**GET /api/v1/spaces/{id}/recycle-bin**

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述     |
| --------- | ------- | ---- | ------ | -------- |
| page      | integer | 否   | 1      | 页码     |
| page_size | integer | 否   | 20     | 每页数量 |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "original_path": "/documents/deleted_file.pdf",
        "file_name": "deleted_file.pdf",
        "file_size": 2048576,
        "recycle_commit_hash": "abc123def456",
        "deleted_by": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员"
        },
        "deleted_at": "2025-01-20T10:30:00Z",
        "expires_at": "2025-02-19T10:30:00Z",
        "days_remaining": 22
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 5,
      "total_pages": 1
    },
    "statistics": {
      "total_size": 104857600,
      "total_count": 5,
      "expires_soon_count": 2
    }
  }
}
```

**恢复回收站文件**

**POST /api/v1/spaces/{id}/recycle-bin/restore**

请求体参数：

| 参数名              | 数据类型  | 是否可空 | 描述                         |
| ------------------- | --------- | -------- | ------------------------------ |
| recycle_commit_hashs | string[]  | 否       | 回收站Git提交哈希数组        |
| restore_to_original | boolean   | 是       | 是否恢复到原路径，默认true   |
| target_path         | string    | 是       | 指定恢复路径                 |

响应示例：

```json
{
  "code": 200,
  "message": "文件恢复完成",
  "data": {
    "restored_count": 2,
    "failed_count": 0,
    "results": [
      {
        "commit_hash": "abc123def456",
        "original_path": "/documents/file1.pdf",
        "restored_path": "/documents/file1.pdf",
        "status": "SUCCESS"
      },
      {
        "commit_hash": "def456ghi789",
        "original_path": "/documents/file2.pdf",
        "restored_path": "/documents/file2.pdf",
        "status": "SUCCESS"
      }
    ]
  }
}
```

**清空回收站**

**DELETE /api/v1/spaces/{id}/recycle-bin**

请求体参数：

| 参数名              | 数据类型  | 是否可空 | 描述                             |
| ------------------- | --------- | -------- | -------------------------------- |
| recycle_commit_hashs | string[]  | 是       | 指定要清空的提交哈希，为空表示全部 |
| confirm             | boolean   | 否       | 确认清空                         |

响应示例：

```json
{
  "code": 200,
  "message": "回收站清空成功",
  "data": {
    "deleted_count": 5,
    "freed_space": 104857600,
    "deleted_at": "2025-01-22T10:35:00Z"
  }
}
```

### 6.4 Git配置管理API详细示例

#### 获取Git配置

**GET /api/v1/spaces/{id}/git-config**

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_name": "张三",
    "user_email": "zhangsan@websoft9.com",
    "auto_commit_interval": 20,
    "recycle_retention_days": 30,
    "config_source": "space_custom",
    "available_sources": {
      "space_custom": {
        "user_name": "张三",
        "user_email": "zhangsan@websoft9.com"
      },
      "user_profile": {
        "user_name": "admin",
        "user_email": "admin@websoft9.com"
      },
      "git_global": {
        "user_name": "Websoft9 System",
        "user_email": "system@websoft9.com"
      }
    }
  }
}
```

#### 设置Git配置

**PUT /api/v1/spaces/{id}/git-config**

请求体参数：

| 参数名               | 数据类型 | 是否可空 | 描述                           |
| -------------------- | -------- | -------- | ------------------------------ |
| user_name            | string   | 是       | Git用户名，为空则从用户信息获取 |
| user_email           | string   | 是       | Git邮箱，为空则从用户信息获取   |
| auto_commit_interval | integer  | 是       | 自动提交间隔（秒）             |
| recycle_retention_days | integer | 是       | 回收站保留天数                 |
| use_user_profile     | boolean  | 是       | 是否使用用户资料作为Git信息     |

请求体示例：

```json
{
  "user_name": "张三",
  "user_email": "zhangsan@websoft9.com",
  "auto_commit_interval": 30,
  "recycle_retention_days": 30,
  "use_user_profile": false
}
```

响应示例：

```json
{
  "code": 200,
  "message": "Git配置更新成功",
  "data": {
    "user_name": "张三",
    "user_email": "zhangsan@websoft9.com",
    "auto_commit_interval": 30,
    "recycle_retention_days": 30,
    "updated_at": "2025-01-22T10:35:00Z"
  }
}
```

### 6.5 Gitea集成API设计

#### 1. 创建 Gitea 仓库

**POST /api/v1/gitea/repos**

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述         |
| ----------- | -------- | -------- | ------------ |
| name        | string   | 否       | 仓库名称     |
| org_id      | string   | 是       | 组织 ID（项目空间对应的 Gitea 组织） |
| private      | boolean  | 是       | 是否私有仓库 |
| description | string   | 是       | 仓库描述     |

响应示例：

```json
{
  "code": 200,
  "message": "仓库创建成功",
  "data": {
    "id": 123,
    "name": "新建仓库",
    "full_name": "org_name/新建仓库",
    "private": true,
    "description": "这是一个新建的仓库",
    "clone_url": "https://gitea.example.com/org_name/新建仓库.git",
    "html_url": "https://gitea.example.com/org_name/新建仓库",
    "ssh_url": "git@gitea.example.com:org_name/新建仓库.git",
    "created_at": "2025-01-22T10:35:00Z",
    "updated_at": "2025-01-22T10:35:00Z"
  }
}
```

#### 2. 获取 Gitea 仓库列表

**GET /api/v1/gitea/repos**

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述         |
| --------- | ------- | ---- | ------ | ------------ |
| page      | integer | 否   | 1      | 页码         |
| page_size | integer | 否   | 20     | 每页数量     |
| org_id    | string  | 是   | -      | 组织 ID（项目空间对应的 Gitea 组织） |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_count": 2,
    "repos": [
      {
        "id": 123,
        "name": "新建仓库",
        "full_name": "org_name/新建仓库",
        "private": true,
        "description": "这是一个新建的仓库",
        "clone_url": "https://gitea.example.com/org_name/新建仓库.git",
        "html_url": "https://gitea.example.com/org_name/新建仓库",
        "ssh_url": "git@gitea.example.com:org_name/新建仓库.git",
        "created_at": "2025-01-22T10:35:00Z",
        "updated_at": "2025-01-22T10:35:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 2,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

#### 3. 删除 Gitea 仓库

**DELETE /api/v1/gitea/repos/{id}**

响应示例：

```json
{
  "code": 200,
  "message": "仓库删除成功",
  "data": null
}
```

#### 4. 添加 Gitea Webhook

**POST /api/v1/gitea/hooks**

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述         |
| ----------- | -------- | -------- | ------------ |
| repo_id     | string   | 否       | 仓库 ID      |
| url         | string   | 否       | Webhook 地址 |
| events      | string[] | 否       | 触发事件列表 |
| active      | boolean  | 是       | 是否激活    |

响应示例：

```json
{
  "code": 200,
  "message": "Webhook 添加成功",
  "data": {
    "id": 456,
    "type": "gitea",
    "url": "http://example.com/webhook",
    "events": ["push", "pull_request"],
    "active": true,
    "created_at": "2025-01-22T10:35:00Z",
    "updated_at": "2025-01-22T10:35:00Z"
  }
}
```

#### 5. 删除 Gitea Webhook

**DELETE /api/v1/gitea/hooks/{id}**

响应示例：

```json
{
  "code": 200,
  "message": "Webhook 删除成功",
  "data": null
}
```

## 7 Data（数据模型与存储设计）

### 7.1 数据架构设计

```
数据层架构（Gitea 集成版）
├── MySQL（映射关系和索引）
│   ├── spaces（空间元数据）
│   ├── space_gitea_mapping（空间与 Gitea 资源映射）
│   ├── user_gitea_mapping（用户映射关系）
│   └── permission_sync_log（权限同步日志）
├── Gitea Repository（主要存储）
│   ├── 仓库内容：文件和版本历史
│   ├── 组织管理：项目空间映射
│   └── 用户权限：团队和成员管理
└── 审计中间件（已有基础设施）
```

### 7.2 数据库表设计

1. **spaces（空间元数据表）**
```sql
CREATE TABLE spaces (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    space_id VARCHAR(255) NOT NULL UNIQUE,  -- UUID
    name VARCHAR(255) NOT NULL,
    description TEXT,
    space_type ENUM('project', 'personal') NOT NULL,
    owner_id BIGINT NOT NULL,
    project_id BIGINT,                      -- 项目空间关联的项目ID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_owner_id (owner_id),
    INDEX idx_project_id (project_id),
    INDEX idx_space_type (space_type)
);
```

2. **space_gitea_mapping（空间与 Gitea 资源映射表）**
```sql
CREATE TABLE space_gitea_mapping (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    space_id BIGINT NOT NULL,
    gitea_repo_id BIGINT NOT NULL,          -- Gitea 仓库 ID
    gitea_org_id BIGINT,                    -- Gitea 组织 ID（项目空间）
    clone_url VARCHAR(512) NOT NULL,        -- 克隆地址
    html_url VARCHAR(512) NOT NULL,         -- Web 访问地址
    ssh_url VARCHAR(512),                   -- SSH 克隆地址
    webhook_id BIGINT,                      -- Webhook ID
    sync_status ENUM('active', 'syncing', 'error') DEFAULT 'active',
    last_sync_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_space_gitea (space_id),
    INDEX idx_gitea_repo_id (gitea_repo_id),
    INDEX idx_gitea_org_id (gitea_org_id),
    INDEX idx_sync_status (sync_status)
);
```

3. **user_gitea_mapping（用户映射关系表）**
```sql
CREATE TABLE user_gitea_mapping (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    platform_user_id BIGINT NOT NULL,
    gitea_user_id BIGINT NOT NULL,
    gitea_username VARCHAR(255) NOT NULL,
    gitea_email VARCHAR(255) NOT NULL,
    sync_status ENUM('active', 'syncing', 'error') DEFAULT 'active',
    last_sync_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_platform_user (platform_user_id),
    UNIQUE KEY uk_gitea_user (gitea_user_id),
    INDEX idx_gitea_username (gitea_username),
    INDEX idx_sync_status (sync_status)
);
```

4. **permission_sync_log（权限同步日志表）**
```sql
CREATE TABLE permission_sync_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    sync_id VARCHAR(255) NOT NULL,          -- 同步任务 ID
    space_id BIGINT,
    user_id BIGINT,
    sync_type ENUM('user', 'permission', 'full') NOT NULL,
    status ENUM('pending', 'running', 'success', 'failed') NOT NULL,
    total_items INT DEFAULT 0,
    processed_items INT DEFAULT 0,
    failed_items INT DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_sync_id (sync_id),
    INDEX idx_space_id (space_id),
    INDEX idx_status (status),
    INDEX idx_sync_type (sync_type)
);
```

## 8 实施路线图建议

### Phase 0.5: Gitea 准备与集成验证（1 周）
- 部署或接入可用的 Gitea 实例（自托管或托管服务）并完成基本连接验证
- 设计并实现 GiteaIntegrationService 的最小可用 API（创建仓库/组织，查询 repo 元数据）
- 实现并验证 ProxyMiddleware 原型（简单的 URL 白名单/黑名单 + 反向代理）

### Phase 1: 基础集成模块（2-3周）
- **Gitea 集成模块**：完成 Gitea API 客户端、基础 CRUD 操作
- **用户权限同步模块**：实现用户映射、基础权限同步
- **空间映射管理模块**：实现空间与 Gitea 资源的映射关系

### Phase 2: UI 集成与访问控制（2-3周）
- **UI 代理模块**：完成 Gitea UI 的反向代理和访问控制
- **单点登录模块**：实现 SSO 认证和令牌管理
- **前端集成**：Gitea UI 嵌入到平台界面

### Phase 3: 高级特性与优化（1-2周）
- **Webhooks 处理**：接收和处理 Gitea 事件
- **权限冲突解决**：处理平台与 Gitea 权限不一致
- **性能优化与监控**：同步性能优化、状态监控

### 8.1 各阶段交付要求

每个 Phase 都应包含：
- 后端服务与中间件开发
- 数据库 Migration
- 单元测试与集成测试
- Gitea 集成测试
- 用户验收测试

### 8.2 关键中间件开发清单

**必须开发的中间件**：
1. **GiteaProxyMiddleware** - Gitea UI 反向代理
2. **URLAccessControlMiddleware** - URL 访问控制
3. **SessionSyncMiddleware** - 会话同步
4. **SSOAuthMiddleware** - 单点登录认证
5. **PermissionSyncMiddleware** - 权限同步触发器

**辅助中间件**：
1. **GiteaAPIMiddleware** - Gitea API 请求增强
2. **WebhookValidationMiddleware** - Webhook 签名验证
3. **ResourceMappingMiddleware** - 资源映射检查

### 8.3 安全注意事项

- **代理安全**：ProxyMiddleware 在注入或转发身份信息时必须严格控制敏感头与 token 的生命周期
- **权限一致性**：平台与 Gitea 的权限模型不应产生冲突，需明确冲突解决策略
- **会话管理**：确保平台会话与 Gitea 会话的安全关联和失效处理
- **API 安全**：所有 Gitea API 调用需要通过平台权限验证
