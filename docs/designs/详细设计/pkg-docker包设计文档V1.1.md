# pkg/docker 包设计文档 V1.1

**说明**：pkg/docker 是 Websoft9 平台的 Docker 通用技术能力包，封装 Docker Engine SDK 和 Docker Compose SDK，提供统一的容器操作接口。

## 文档元信息

- **负责人**：
- **审核**：
- **创建日期**：2025-11-12
- **版本**：V1.1

---

## 1. 设计目标

### 1.1 核心定位

- 🎯 **技术能力包**：封装 Docker Engine SDK 和 Compose SDK，提供统一接口
- 🎯 **可复用组件**：可在多个项目和场景中使用（API Service、Agent、Workflow、CLI 工具）
- 🎯 **抽象层**：屏蔽底层 SDK 的复杂性，提供简洁易用的 API
- 🎯 **无状态设计**：不依赖数据库，不包含业务逻辑

### 1.2 技术基础

- **Docker Engine SDK**: `github.com/docker/docker/client` - 管理容器、镜像、网络、卷
- **Docker Compose SDK**: `github.com/docker/compose/v2` - 管理多容器编排
- **Go 版本要求**: Go 1.24+
- **依赖版本**:
  - `github.com/docker/docker` v25.0.0+
  - `github.com/docker/compose/v2` v2.24.0+

### 1.3 设计原则

1. **简单易用**：一行代码创建连接，扁平化 API，合理默认值
2. **单连接设计**：一个 Docker 对象管理一个 Docker 主机，多服务器由上层管理
3. **通用性**：不包含业务逻辑，纯技术操作，可独立使用
4. **无状态**：不依赖数据库，所有数据实时查询

### 1.4 功能列表

#### 1.4.1 容器管理 (Container)

- ✅ **创建容器**：基于镜像创建容器，支持环境变量、端口映射、卷挂载、网络配置
- ✅ **生命周期控制**：启动、停止、重启、删除容器
- ✅ **状态查询**：列表查询、详情查看、运行状态监控
- ✅ **日志管理**：实时日志流、历史日志查询
- ✅ **命令执行**：在运行中的容器内执行命令
- ✅ **资源监控**：CPU、内存、网络、磁盘 I/O 统计

#### 1.4.2 镜像管理 (Image)

- ✅ **镜像获取**：从仓库拉取镜像、从 Dockerfile 构建镜像、导入镜像文件
- ✅ **镜像发布**：推送到仓库、导出为文件、打标签
- ✅ **镜像管理**：列表查询、详情查看、删除镜像、清理未使用镜像

#### 1.4.3 网络管理 (Network)

- ✅ **网络创建**：创建自定义网络（bridge、overlay、host 等）
- ✅ **网络连接**：容器加入/退出网络
- ✅ **网络查询**：列表查询、详情查看
- ✅ **网络清理**：删除网络、清理未使用网络

#### 1.4.4 卷管理 (Volume)

- ✅ **卷创建**：创建数据卷用于持久化存储
- ✅ **卷查询**：列表查询、详情查看
- ✅ **卷清理**：删除卷、清理未使用卷

#### 1.4.5 编排管理 (Compose)

- ✅ **项目部署**：基于 docker-compose.yaml 部署多容器应用
- ✅ **项目控制**：启动、停止、重启、删除整个项目
- ✅ **服务管理**：控制项目中的单个或多个服务
- ✅ **状态查询**：查看项目列表、服务状态、日志输出
- ✅ **镜像管理**：拉取项目所需镜像
- ✅ **配置验证**：验证 docker-compose.yaml 语法正确性

---

## 2. 架构设计

### 2.1 技术架构

**双 SDK 架构**：
```
pkg/docker
├── Container Manager  → Docker Engine SDK
├── Image Manager      → Docker Engine SDK  
├── Network Manager    → Docker Engine SDK
├── Volume Manager     → Docker Engine SDK
└── Compose Manager    → Compose SDK
         ↓
    Docker Daemon
```

### 2.2 包结构

```
pkg/docker/
├── docker.go              # 核心入口
├── options.go             # 配置选项
├── errors.go              # 错误定义
├── container/             # 容器管理
├── image/                 # 镜像管理
├── network/               # 网络管理
├── volume/                # 卷管理
└── compose/               # 编排管理
```

---

## 3. 核心接口

### 3.1 创建 Docker 管理器

**三种创建方式**：

```go
// 方式1：本地自动检测（Agent 场景）
func NewLocal(opts ...Option) (*Docker, error)

// 方式2：远程连接（API Service 场景）
func NewRemote(host string, port int, opts ...Option) (*Docker, error)

// 方式3：通用接口（底层实现）
func New(host string, opts ...Option) (*Docker, error)
```

**配置选项**：

```go
// WithTLS 配置 TLS 连接
func WithTLS(certPath, keyPath, caPath string) Option

// WithTimeout 配置请求超时时间
func WithTimeout(timeout time.Duration) Option

// WithAPIVersion 配置 Docker API 版本
func WithAPIVersion(version string) Option

// WithHTTPClient 配置自定义 HTTP 客户端
func WithHTTPClient(client *http.Client) Option

// WithDialContext 配置自定义拨号上下文
func WithDialContext(dialContext func(ctx context.Context, network, addr string) (net.Conn, error)) Option
```

**Docker 对象**：

```go
type Docker struct {
    Container *ContainerManager  // 容器操作
    Image     *ImageManager      // 镜像操作
    Network   *NetworkManager    // 网络操作
    Volume    *VolumeManager     // 卷操作
    Compose   *ComposeManager    // 编排操作
}

// 连接管理
func (d *Docker) Close() error                                // 关闭连接
func (d *Docker) Ping(ctx context.Context) error              // 测试连接

// 信息查询
func (d *Docker) Info(ctx context.Context) (*SystemInfo, error)      // 获取系统信息
func (d *Docker) Version(ctx context.Context) (*VersionInfo, error)  // 获取 Docker 版本
```

---

### 3.2 容器管理接口

```go
type ContainerManager struct {
    // 生命周期管理
    Create(ctx context.Context, config *ContainerConfig) (containerID string, error)    // 创建容器
    Start(ctx context.Context, containerID string) error                                 // 启动容器
    Stop(ctx context.Context, containerID string, timeout *int) error                   // 停止容器
    Restart(ctx context.Context, containerID string) error                               // 重启容器
    Remove(ctx context.Context, containerID string, force bool) error                   // 删除容器
    Pause(ctx context.Context, containerID string) error                                 // 暂停容器
    Unpause(ctx context.Context, containerID string) error                               // 恢复容器
    Kill(ctx context.Context, containerID string, signal string) error                  // 强制终止容器
    Wait(ctx context.Context, containerID string) (<-chan WaitResponse, <-chan error)   // 等待容器停止
    Update(ctx context.Context, containerID string, config UpdateConfig) error          // 更新容器配置
    
    // 查询和监控
    List(ctx context.Context, options ListOptions) ([]Container, error)                 // 列出所有容器
    Inspect(ctx context.Context, containerID string) (*ContainerInfo, error)            // 查看容器详情
    Logs(ctx context.Context, containerID string, options LogOptions) (io.ReadCloser, error)  // 获取容器日志
    Stats(ctx context.Context, containerID string, stream bool) (<-chan Stats, error)   // 获取容器资源统计
    
    // 执行命令
    Exec(ctx context.Context, containerID string, cmd []string) (output string, error)  // 在容器内执行命令
}
```

---

### 3.3 镜像管理接口

```go
type ImageManager struct {
    // 镜像获取
    Pull(ctx context.Context, ref string, options PullOptions) (io.ReadCloser, error)           // 从仓库拉取镜像
    Build(ctx context.Context, buildContext io.Reader, options BuildOptions) (io.ReadCloser, error)  // 从 Dockerfile 构建镜像
    Load(ctx context.Context, input io.Reader) error                                             // 从文件导入镜像
    
    // 镜像发布
    Push(ctx context.Context, ref string, options PushOptions) (io.ReadCloser, error)           // 推送镜像到仓库
    Save(ctx context.Context, images []string) (io.ReadCloser, error)                           // 导出镜像为文件
    Tag(ctx context.Context, source, target string) error                                        // 为镜像打标签
    
    // 镜像管理
    List(ctx context.Context, options ListOptions) ([]Image, error)                             // 列出所有镜像
    Inspect(ctx context.Context, imageID string) (*ImageInfo, error)                            // 查看镜像详情
    Remove(ctx context.Context, imageID string, force bool) error                               // 删除镜像
    Prune(ctx context.Context, options PruneOptions) (*PruneReport, error)                      // 清理未使用的镜像
}
```

---

### 3.4 网络管理接口

```go
type NetworkManager struct {
    // 网络生命周期
    Create(ctx context.Context, name string, options CreateOptions) (networkID string, error)  // 创建自定义网络
    Remove(ctx context.Context, networkID string) error                                        // 删除网络
    
    // 网络连接
    Connect(ctx context.Context, networkID, containerID string) error                          // 容器加入网络
    Disconnect(ctx context.Context, networkID, containerID string, force bool) error           // 容器退出网络
    
    // 查询
    List(ctx context.Context, options ListOptions) ([]Network, error)                          // 列出所有网络
    Inspect(ctx context.Context, networkID string) (*NetworkInfo, error)                       // 查看网络详情
    Prune(ctx context.Context, options PruneOptions) (*PruneReport, error)                     // 清理未使用的网络
}
```

---

### 3.5 卷管理接口

```go
type VolumeManager struct {
    // 卷生命周期
    Create(ctx context.Context, name string, options CreateOptions) (*Volume, error)   // 创建数据卷
    Remove(ctx context.Context, volumeName string, force bool) error                   // 删除数据卷
    
    // 查询
    List(ctx context.Context, options ListOptions) ([]*Volume, error)                  // 列出所有数据卷
    Inspect(ctx context.Context, volumeName string) (*Volume, error)                   // 查看数据卷详情
    
    // 维护
    Prune(ctx context.Context, options PruneOptions) (*PruneReport, error)             // 清理未使用的数据卷
}
```

---

### 3.6 Compose 管理接口

```go
type ComposeManager struct {
    // 项目生命周期
    Deploy(ctx context.Context, projectName string, composeYAML string, options DeployOptions) error  // 部署 Compose 项目
    Down(ctx context.Context, projectName string, options DownOptions) error                          // 删除 Compose 项目
    
    // 服务控制
    Start(ctx context.Context, projectName string, services []string) error                           // 启动项目服务
    Stop(ctx context.Context, projectName string, services []string) error                            // 停止项目服务
    Restart(ctx context.Context, projectName string, services []string) error                         // 重启项目服务
    
    // 查询
    List(ctx context.Context) ([]Project, error)                                                      // 列出所有项目
    Ps(ctx context.Context, projectName string) ([]Service, error)                                    // 查看项目服务状态
    Logs(ctx context.Context, projectName string, services []string, options LogOptions) (io.ReadCloser, error)  // 获取项目日志
    
    // 镜像管理
    Pull(ctx context.Context, projectName string) error                                               // 拉取项目所需镜像
    
    // 验证
    Validate(composeYAML string) error                                                                 // 验证 Compose 配置
}
```

---

## 4. 错误处理

pkg/docker 的错误处理遵循 API Service 的规范，详见 `api-service/pkg/errors`。

**核心原则**：
- 使用自定义错误类型，包含错误码、消息、上下文信息
- 区分连接错误、资源错误、操作错误
- 透传底层 Docker SDK 错误，添加操作上下文
- 所有错误以 `docker:` 前缀标识

---

## 5. 常量定义

### 5.1 默认配置

```go
const (
    // 连接相关
    DefaultLocalHost   = "unix:///var/run/docker.sock"  // 本地 Socket 路径
    DefaultRemotePort  = 2375                            // HTTP 端口
    DefaultTLSPort     = 2376                            // HTTPS 端口
    
    // 超时配置
    DefaultTimeout     = 30 * time.Second                // 默认请求超时
    DefaultDialTimeout = 5 * time.Second                 // 连接超时
    
    // API 版本
    DefaultAPIVersion  = "1.43"                          // 默认 API 版本
    MinAPIVersion      = "1.40"                          // 最低支持版本
)
```

### 5.2 限制配置

```go
const (
    // 日志限制
    MaxLogLines        = 1000                            // 最大日志行数
    MaxLogSize         = 10 * 1024 * 1024               // 最大日志大小 10MB
    
    // 构建限制
    MaxBuildContextSize = 100 * 1024 * 1024             // 最大构建上下文 100MB
    
    // 并发限制
    MaxConcurrentPulls  = 5                              // 最大并发拉取数
)
```

---

## 6. 日志处理

pkg/docker 使用项目统一的日志包 `api-service/pkg/logger`，保持与 API Service 一致的日志风格。

**日志规范**：
- 使用结构化日志（Zap）
- 支持上下文日志（自动提取 request_id、user_id）
- 日志级别：Debug、Info、Warn、Error、Fatal
- 使用 logger.String()、logger.Int()、logger.Duration() 等字段方法

详细说明参考：`api-service/pkg/logger`

---

## 7. 使用场景

### 7.1 主要场景：Agent 本地执行（推荐）

**架构流程**：
```
用户操作 → API Service → gRPC/消息队列 → Agent → pkg/docker.NewLocal() → Docker Daemon
         (管理服务器)                    (应用服务器)  (本地 socket)
```

**特点**：
- ✅ **安全**：无需暴露 Docker API 端口
- ✅ **简单**：本地 Unix Socket 通信
- ✅ **高性能**：无网络开销
- ✅ **防火墙友好**：只需 Agent 心跳连接

**使用方式**：
```go
// Agent 启动时创建 Docker 连接
docker, _ := docker.NewLocal()
defer docker.Close()

// 接收 API Service 下发的任务
task := <-taskQueue
docker.Container.Create(ctx, task.Config)
```

---

### 7.2 工作流场景

**适用场景**：
- 自动化部署流程
- 容器编排任务
- 定时运维操作
- 批量容器管理

**架构流程**：
```
工作流引擎 → 工作流节点 → pkg/docker → Docker Daemon
(Workflow)   (Task Node)   (本地/远程)
```

**特点**：
- 🔄 **串行/并行执行**：支持复杂的容器操作流程
- 🎯 **条件分支**：根据容器状态执行不同操作
- 📊 **状态跟踪**：记录每个步骤的执行结果
- 🔁 **错误重试**：失败自动重试机制

**使用方式**：
```go
// 工作流节点中调用 pkg/docker
func DeployWorkflowNode(ctx context.Context) error {
    docker, _ := docker.NewLocal()
    defer docker.Close()
    
    // 1. 拉取镜像
    docker.Image.Pull(ctx, "nginx:latest", nil)
    
    // 2. 创建网络
    networkID, _ := docker.Network.Create(ctx, "app-network", nil)
    
    // 3. 部署应用
    docker.Compose.Deploy(ctx, "my-app", composeYAML, nil)
    
    return nil
}
```

---

### 7.3 备用场景：远程直连（特殊情况）

**适用场景**：
- 开发调试环境
- 单机部署模式（所有服务在一台机器）
- 临时运维操作

**架构流程**：
```
API Service → pkg/docker.NewRemote() → Docker Daemon (TCP 2375/2376)
(管理服务器)                          (应用服务器)
```

**注意事项**：
- ⚠️ 需要开放 Docker API 端口（2375 或 2376）
- ⚠️ 必须配置 TLS 加密（生产环境）
- ⚠️ 需要配置防火墙规则

**使用方式**：
```go
// 远程连接需要配置 TLS
docker, _ := docker.NewRemote("192.168.1.10", 2376, 
    docker.WithTLS("/path/to/cert.pem", "/path/to/key.pem", "/path/to/ca.pem"))
defer docker.Close()
```

---

## 8. 性能优化

### 8.1 连接管理

- **连接复用**：上层服务维护连接池，避免重复创建
- **超时控制**：合理设置超时时间，避免长时间阻塞
- **并发控制**：限制并发请求数，避免资源耗尽

### 8.2 流式处理

- **日志流**：使用 io.ReadCloser 流式读取，避免内存占用
- **统计流**：使用 channel 实时推送统计数据
- **构建输出**：流式输出构建进度，实时反馈

---

## 9. 测试策略

### 9.1 单元测试

- 使用 Mock Docker Client 测试各个 Manager
- 覆盖正常流程和异常流程
- 验证错误处理和边界条件

### 9.2 集成测试

- 使用真实 Docker 环境测试
- 验证端到端功能
- 测试并发场景

---

## 10. 版本记录

| 版本 | 日期 | 说明 |
|------|------|------|
| V1.1 | 2025-11-12 | 精简版：删除详细代码实现，保留核心接口设计；补充错误处理、常量定义、依赖说明、日志处理 |
| V1.0 | 2025-11-12 | 初始版本 |

---

## 11. 参考文档

- **Docker Engine SDK**: https://pkg.go.dev/github.com/docker/docker/client
- **Docker Compose SDK**: https://github.com/docker/compose/blob/main/docs/sdk.md
- **Docker API**: https://docs.docker.com/engine/api/
- **Websoft9 Logger**: `api-service/pkg/logger` - 项目统一日志包
