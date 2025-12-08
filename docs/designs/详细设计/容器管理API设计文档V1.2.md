# 容器管理 API 设计文档

**说明**：容器管理模块提供 Docker 和 Docker Compose 相关的 RESTful API 服务，供前端实现 Docker 管理界面使用。

## 文档元信息

- **负责人**：
- **审核**：
- **创建日期**：2025-11-13
- **更新日期**：2025-11-13
- **版本**：V1.2

---

## 1. 设计概述

### 1.1 核心定位

容器管理模块是 **Websoft9 平台的 Docker API 服务层**：

**核心职责**：
- 🎯 **Docker API 封装**：提供 Docker 和 Docker Compose 的 RESTful API
- 🎯 **前端服务**：为前端 Docker 管理界面提供数据接口
- 🎯 **多服务器管理**：支持管理多台服务器的 Docker 资源
- 🎯 **Agent 协作**：**所有 Docker 操作通过 Agent 执行**（不直接连接 Docker）

### 1.2 架构设计

```
┌─────────────────────────────────────────────┐
│         前端 Docker 管理界面                   │
│     (调用统一的 API Service 入口)              │
└──────────────────┬──────────────────────────┘
                   │ HTTP REST
                   │ GET /api/v1/servers/{server_id}/containers
                   ▼
┌─────────────────────────────────────────────┐
│      API Service (容器管理 API 模块)          │
│  • 接收前端请求                               │
│  • 统一认证/鉴权 (JWT + RBAC)                 │
│  • 参数验证                                   │
│  • 根据 server_id 路由到对应的 Agent         │
│  • 数据聚合和处理                             │
│  • 审计日志                                   │
└──────┬────────┬────────┬─────────────────────┘
       │        │        │ gRPC 通信
       │        │        │ (发送 Docker 操作指令)
       ▼        ▼        ▼
   ┌────────┬────────┬────────┐
   │Agent-1 │Agent-2 │Agent-3 │  (部署在各目标服务器)
   │        │        │        │
   │pkg/docker│pkg/docker│pkg/docker│  (本地调用)
   └───┬────┴───┬────┴───┬────┘
       │        │        │ Unix Socket
       │        │        │ /var/run/docker.sock
       ▼        ▼        ▼
   [Docker] [Docker] [Docker]
   服务器1   服务器2   服务器3
```

**执行流程**：

1. **前端发起请求**：
   ```
   GET /api/v1/servers/1/containers
   ```

2. **API Service 处理**：
   - 验证用户 JWT Token
   - 检查用户对服务器1的访问权限
   - 确定目标 Agent（服务器1的 Agent）

3. **调用 Agent**：
   ```
   API Service → gRPC → Agent-1
   请求：ListContainers()
   ```

4. **Agent 执行**：
   ```go
   // Agent 在本地调用 pkg/docker
   docker := docker.NewLocal()  // 连接本地 Docker
   containers, err := docker.Container.List(ctx, options)
   ```

5. **返回结果**：
   ```
   Agent-1 → gRPC → API Service → HTTP → 前端
   ```

**为什么需要 API Service 层？**

| 功能 | 说明 |
|------|------|
| **统一入口** | 前端只需对接一个 API 地址，无需管理多个 Agent |
| **多服务器路由** | 根据 `server_id` 自动路由到对应 Agent |
| **统一认证** | JWT + RBAC 权限控制，用户只能访问授权的服务器 |
| **数据聚合** | 可聚合多台服务器数据（如：查询所有容器）|
| **审计日志** | 记录所有操作日志（谁在什么时间对哪台服务器做了什么）|
| **资源配额** | 限制用户创建容器的数量、资源配置等 |
| **统一错误处理** | 标准化错误响应和国际化 |

**使用场景**：
- ✅ **前端 Docker 管理界面**：调用本模块 API 实现容器管理功能
- ✅ **跨服务器容器管理**：通过 API 管理多台服务器的 Docker 资源
- ❌ **工作流等模块**：可直接使用底层 `pkg/docker` 包，无需通过本 API

**功能边界**：
- ✅ **本模块负责**：
  - 提供容器、镜像、网络、卷、Compose 的 RESTful API
  - 管理多服务器的 Docker 资源
  - **通过 Agent gRPC 调用执行所有 Docker 操作**
  
- ❌ **不负责**：
  - 业务编排逻辑（由工作流模块负责）
  - Docker 安装和配置（由其他模块负责）
  - 批量操作（由前端自行实现）
  - **直接连接 Docker**（一律通过 Agent）

---

### 1.3 核心功能

提供基础 Docker 操作的 RESTful API：

#### 1.2.1 容器管理 (Container)

- ✅ **列出容器**：查询容器列表（支持过滤）
- ✅ **创建容器**：基于镜像创建容器
- ✅ **容器详情**：获取容器详细信息
- ✅ **启动/停止/重启**：容器生命周期控制
- ✅ **删除容器**：删除已停止的容器
- ✅ **容器日志**：获取容器日志
- ✅ **容器统计**：获取容器资源使用情况

**说明**：所有容器操作最终由 Agent 调用 `pkg/docker` 在目标服务器上执行。

#### 1.3.2 镜像管理 (Image)

- ✅ **列出镜像**：查询镜像列表
- ✅ **拉取镜像**：从仓库拉取镜像
- ✅ **删除镜像**：删除本地镜像
- ✅ **镜像详情**：获取镜像详细信息

**说明**：所有镜像操作最终由 Agent 调用 `pkg/docker` 在目标服务器上执行。

#### 1.3.3 网络管理 (Network)

- ✅ **列出网络**：查询网络列表
- ✅ **创建网络**：创建自定义网络
- ✅ **删除网络**：删除网络
- ✅ **网络详情**：获取网络详细信息

**说明**：所有网络操作最终由 Agent 调用 `pkg/docker` 在目标服务器上执行。

#### 1.3.4 卷管理 (Volume)

- ✅ **列出卷**：查询卷列表
- ✅ **创建卷**：创建数据卷
- ✅ **删除卷**：删除卷
- ✅ **卷详情**：获取卷详细信息

**说明**：所有卷操作最终由 Agent 调用 `pkg/docker` 在目标服务器上执行。

#### 1.3.5 Compose 管理

- ✅ **部署项目**：基于 docker-compose.yaml 部署项目
- ✅ **列出项目**：查询 Compose 项目列表
- ✅ **停止项目**：停止 Compose 项目
- ✅ **启动项目**：启动已停止的项目
- ✅ **重新部署项目**：重新拉取镜像并重新部署（Re-pull image and redeploy）
- ✅ **删除项目**：删除 Compose 项目

**说明**：所有 Compose 操作最终由 Agent 调用 `pkg/docker` 在目标服务器上执行。

---

## 2. 非功能性需求

### 2.1 性能要求

| 指标 | 要求 | 说明 |
|------|------|------|
| **API 响应时间** | P95 < 200ms | 除长时操作（拉取镜像、日志流等） |
| **并发 API 请求** | 100 QPS | 单个 API Service 实例 |

### 2.2 可靠性要求

| 指标 | 要求 | 说明 |
|------|------|------|
| **服务可用性** | 99.9% | API Service 可用性 |
| **Agent 通信** | 自动重连 | Agent 断线后自动重连 |

### 2.3 安全性要求

| 要求 | 说明 |
|------|------|
| **API 认证** | 所有 API 需要 JWT 认证 |
| **权限控制** | 基于角色的操作权限 |
| **审计日志** | 记录所有操作日志 |

---

## 3. API 设计

### 3.1 API 概览

容器管理模块提供 **5 大类基础 Docker API**：

| 类别 | 章节 | API 数量 | 说明 |
|------|------|---------|------|
| 容器管理 | 3.2 | 7 | 容器生命周期、日志、统计等 |
| 镜像管理 | 3.3 | 4 | 镜像拉取、查询、删除等 |
| 网络管理 | 3.4 | 4 | 网络创建、查询、删除等 |
| 卷管理 | 3.5 | 4 | 卷创建、查询、删除等 |
| Compose 管理 | 3.6 | 6 | 项目部署、启动、停止、重建、删除等 |

**API 总数**：约 **25 个基础 API**

**说明**：
- 所有 API 路径格式：`/api/v1/servers/{server_id}/{resource}`
- `server_id`: 目标服务器 ID
- 前端可基于基础 API 实现批量操作等高级功能

---

### 3.2 容器管理 API

#### 3.2.1 列出容器

```
GET /api/v1/servers/{server_id}/containers
```

**Query 参数**：
```
all=true         // 是否包含停止的容器
limit=50         // 返回数量限制
```

**响应**：
```json
{
  "code": 0,
  "data": [
    {
      "id": "abc123",
      "name": "my-app",
      "image": "nginx:latest",
      "state": "running",
      "status": "Up 2 hours",
      "created": "2025-11-13T10:00:00Z",
      "ports": [
        {"container_port": 80, "host_port": 8080, "protocol": "tcp"}
      ]
    }
  ]
}
```

#### 3.2.2 创建容器

```
POST /api/v1/servers/{server_id}/containers
```

**请求体**：
```json
{
  "name": "my-app",
  "image": "nginx:latest",
  "env": {"ENV_VAR": "value"},
  "ports": [
    {"container_port": 80, "host_port": 8080, "protocol": "tcp"}
  ],
  "volumes": [
    {"type": "volume", "source": "my-data", "target": "/data"}
  ],
  "networks": ["my-network"],
  "restart_policy": "unless-stopped"
}
```

#### 3.2.3 容器详情

```
GET /api/v1/servers/{server_id}/containers/{container_id}
```

#### 3.2.4 启动容器

```
POST /api/v1/servers/{server_id}/containers/{container_id}/start
```

#### 3.2.5 停止容器

```
POST /api/v1/servers/{server_id}/containers/{container_id}/stop
```

#### 3.2.6 删除容器

```
DELETE /api/v1/servers/{server_id}/containers/{container_id}
```

**Query 参数**：
```
force=true  // 是否强制删除运行中的容器
```

#### 3.2.7 获取容器日志

```
GET /api/v1/servers/{server_id}/containers/{container_id}/logs
```

**Query 参数**：
```
tail=100         // 最后 N 行日志
follow=false     // 是否跟随日志流
timestamps=true  // 是否显示时间戳
```

---

### 3.3 镜像管理 API

#### 3.3.1 列出镜像

```
GET /api/v1/servers/{server_id}/images
```

#### 3.3.2 拉取镜像

```
POST /api/v1/servers/{server_id}/images/pull
```

**请求体**：
```json
{
  "image": "nginx:latest",
  "auth": {
    "username": "user",
    "password": "pass"
  }
}
```

#### 3.3.3 镜像详情

```
GET /api/v1/servers/{server_id}/images/{image_id}
```

#### 3.3.4 删除镜像

```
DELETE /api/v1/servers/{server_id}/images/{image_id}
```

---

### 3.4 网络管理 API

#### 3.4.1 列出网络

```
GET /api/v1/servers/{server_id}/networks
```

#### 3.4.2 创建网络

```
POST /api/v1/servers/{server_id}/networks
```

**请求体**：
```json
{
  "name": "my-network",
  "driver": "bridge",
  "ipam": {
    "subnet": "172.20.0.0/16"
  }
}
```

#### 3.4.3 网络详情

```
GET /api/v1/servers/{server_id}/networks/{network_id}
```

#### 3.4.4 删除网络

```
DELETE /api/v1/servers/{server_id}/networks/{network_id}
```

---

### 3.5 卷管理 API

#### 3.5.1 列出卷

```
GET /api/v1/servers/{server_id}/volumes
```

#### 3.5.2 创建卷

```
POST /api/v1/servers/{server_id}/volumes
```

**请求体**：
```json
{
  "name": "my-data",
  "driver": "local",
  "labels": {
    "project": "my-project"
  }
}
```

#### 3.5.3 卷详情

```
GET /api/v1/servers/{server_id}/volumes/{volume_name}
```

#### 3.5.4 删除卷

```
DELETE /api/v1/servers/{server_id}/volumes/{volume_name}
```

---

### 3.6 Compose 管理 API

#### 3.6.1 部署 Compose 项目

```
POST /api/v1/servers/{server_id}/compose
```

**请求体**：
```json
{
  "project_name": "my-app",
  "compose_content": "version: '3'\nservices:\n  web:\n    image: nginx\n    ports:\n      - 80:80"
}
```

#### 3.6.2 列出 Compose 项目

```
GET /api/v1/servers/{server_id}/compose
```

#### 3.6.3 停止 Compose 项目

```
POST /api/v1/servers/{server_id}/compose/{project}/stop
```

#### 3.6.4 启动 Compose 项目

```
POST /api/v1/servers/{server_id}/compose/{project}/start
```

#### 3.6.5 重新部署 Compose 项目

```
POST /api/v1/servers/{server_id}/compose/{project}/redeploy
```

**说明**：重新拉取镜像并强制重新部署所有服务（Re-pull image and redeploy）

**Query 参数**：
```
no_cache=false   // 是否不使用缓存构建镜像（如果有 build）
pull=true        // 是否重新拉取镜像（默认 true）
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "project_name": "my-app",
    "redeployed_services": ["web", "db"],
    "pulled_images": ["nginx:latest", "mysql:8.0"],
    "status": "running"
  },
  "message": "Project redeployed successfully"
}
```

**操作流程**：
1. 拉取最新镜像：`docker compose pull`
2. 强制重建并启动：`docker compose up -d --force-recreate`
3. 保留数据卷，仅重建容器

**使用场景**：
- 更新应用到最新版本
- 强制刷新所有容器
- 修复容器配置问题

#### 3.6.6 删除 Compose 项目

```
DELETE /api/v1/servers/{server_id}/compose/{project}
```

**Query 参数**：
```
volumes=true  // 是否删除关联的卷
```

---

## 4. 数据模型

### 4.1 请求/响应 DTO

参考 `api-service/internal/dto/` 设计规范：

```go
// CreateContainerRequest 创建容器请求
type CreateContainerRequest struct {
    Name          string            `json:"name" binding:"required"`
    Image         string            `json:"image" binding:"required"`
    Env           map[string]string `json:"env"`
    Ports         []PortMapping     `json:"ports"`
    Volumes       []VolumeMount     `json:"volumes"`
    Networks      []string          `json:"networks"`
    RestartPolicy string            `json:"restart_policy"`
}

// ContainerResponse 容器响应
type ContainerResponse struct {
    ID      string    `json:"id"`
    Name    string    `json:"name"`
    Image   string    `json:"image"`
    State   string    `json:"state"`
    Status  string    `json:"status"`
    Created time.Time `json:"created"`
}
```

---

## 5. 错误处理

遵循 `api-service/pkg/errors` 错误处理规范：

```go
var (
    ErrContainerNotFound = errors.New("container not found")
    ErrImageNotFound     = errors.New("image not found")
    ErrPermissionDenied  = errors.New("permission denied")
)
```

**HTTP 状态码映射**：
- `400 Bad Request`: 参数验证失败
- `401 Unauthorized`: 未认证
- `403 Forbidden`: 权限不足
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误

---

## 6. 测试策略

### 6.1 单元测试

- Controller 层测试（Mock Service）
- Service 层测试（Mock Agent Client）

### 6.2 集成测试

- API → Agent → Docker 端到端测试
- 使用真实 Docker 环境

---

## 7. 版本记录

| 版本 | 日期 | 说明 |
|------|------|------|
| V1.2 | 2025-11-13 | 简化设计，聚焦基础 Docker API |
| V1.1 | 2025-11-13 | 添加高级功能（已删除） |
| V1.0 | 2025-11-12 | 初始版本 |

---

## 8. 参考文档

- **pkg/docker 设计文档**: `docs/designs/详细设计/pkg-docker包设计文档V1.1.md`
- **API Service 架构**: `api-service/` 目录
- **Docker SDK**: https://pkg.go.dev/github.com/docker/docker/client
- **Compose SDK**: https://github.com/docker/compose/blob/main/docs/sdk.md
