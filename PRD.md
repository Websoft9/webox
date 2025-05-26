# PRD

## 概要

- 产品名称：应用聚合与托管平台
- 价值体现：聚合/接入/运行/连接一切应用或服务，打开 Websoft9 控制台，开启一天轻松基于 Web 随时随地办公模式
- 业务功能：寻找、部署、资源接入、运行、发布、管理、安全、订阅、咨询

## 功能需求

功能需求描述要设计的系统的功能。它描述了系统将是什么以及它将如何发挥作用以满足用户需求。简单说，它解决了客户的业务问题。

### 寻找

基于 AI + 知识库，根据用户需求推荐具体的软件集。

- 支持主流 AI 平台 API
- 接入不同的知识库
  - Websoft9 在线知识库
  - 软件测评、软件独立博客网站知识库
- 交互式推荐软件或软件组合

### 部署

将软件部署到目标服务器：

- 多服务器支持
- 目标服务器环境准备
- Docker/Podman 容器引擎支持
- 部署编排模式
  - docker-compose
  - docker run
  - 基于源码构建镜像
- 部署流程可视化
- 异步
- 部署物准备
  - 镜像
  - 软件仓库
  - Websoft9 制品库
- 私有云

### 运行

### 管理

### 发布

- 支持将应用发布到多种类型的网关（包含 Websoft9 默认的网关）
- 绑定域名
- 无域名下的访问管理

### 网关

内置轻量级网关，可以完成：

- 反向代理到容器的端口
- 自动申请 HTTPS 证书
- 绑定域名（包含通配符域名）
- 客户端到后端的流量可控、可配置
- 支持个性化的配置策略

### 安全

### 智库

提供应用和科技相关知识，结合 RAG/LLM 供用户使用

### 资源

使用 Websoft9 过程中，需要接入、维护的对象。将外部应用全局接入到 Websoft9，更方便快捷登录使用

- 服务器
- 数据库
- 域名
- 证书
- SFTP
- S3 以及兼容
- File (shell 或 Python 脚本,zip, html，.sql, ini, md 等)
- LLM
- 账号或密钥
- 制品库（pip, node, docker 等仓库）
- SaaS 应用
- HTTP request（包含 webhook）
- git

### 监控

- 提供 Websoft9 控制台的资源占用情况
- 提供应用所有容器的资源占用情况
- 监控支持多种时间周期
- 资源包括：CPU/内存/网络/磁盘空间

### License 管理

- 支持导入/更新 License，激活商业版本
- License 离线可用
- 系统根据 License 类型，呈现不同的功能特征

### 插件

用户通过安装插件，提升 Websoft9 的功能：

- 应该是前端程序与 Websoft9 plugin 相关的 API
- 支持 declarative（声明式编程） 和 programmatic（函数式编程）
- 插件入口设计

### 用户贡献应用

- 用户可以为 AppStore 贡献应用
- 用户可以将自定义应用上架的 AppStore，但保持私有

## 非功能需求

非功能需求定义了软件系统的质量属性，解释了要设计的系统的限制和约束。

### 性能

- 页面在 1s 内响应
- 尽量支持高并发，在高负载下仍然可以正常响应
- 应用启动、重建和部署的时间应尽可能短，以提升用户体验
- 拉取镜像的速度尽可能加快

### 可用性

- 安装或构建时，在线制品库保持 99% 的可用性
- 

### 可维护

- 易于升级，向下兼容
- 提供多种无用资源的清理策略
- 完整的日志记录、查询、分析、告警与归档

### 安全性

- 支持授权用户/角色访问相关功能
- 支持非授权用户有限的功能使用（例如：部署软件）
- 密码不允许以明文形式显示

### 审计

- 需记录到审计日志的重要操作：应用管理的所有动作、容器操作、资源操作

### 国际化

- 系统默认支持多语言
- 中文和英文的翻译是 100% 可用的
- 

### 组网

所有的节点之间如何具备公网 IP 或可达的局域网 IP，那么使用原生的容器编排技术进行组网

- docker swarm
- k8s
- nomad

如果网络不通，多服务器组建需组件虚拟局域网，实现集中管理。  

可能的技术方案：

- 中央服务器转发：VPN,SSH 隧道，流量完全走中央服务器
- P2P机制：NAT穿透+中央协调（headscale）
- CNI: flannel

### 事件驱动与微服务编排（胶水代码）

将事务拆分为在多个微服务执行，暂定采用 [temporal](https://temporal.io/)

### API 网关

通过一个统一的 API 入口，便于各种客户端调用。暂定：[krakend](https://www.krakend.io/) + Redis

### 配置管理

### 身份认证与访问授权

- 存储微服务的凭证
- 更新凭证（事件/轮询/Webhook 三种机制之一）
- 多用户、多角色

### API/CLI

- API 优先，先有 API，才有界面和 CLI
- CLI是 HTTP API 的精简包装器，CLI 命令在内部直接映射到 HTTP API

## 架构

架构设计中关注的要点：

- 技术功能：多用户多权限、Vault, Git, SSO, Variables, Orchestration, AI workfolw, Connection pipeline, Embeddeding applications/micro frontends
- 架构组件：前端(UI/CLI)、server, BFF 组件、API 网关，Agent, Controller, Event Stream, Storage, Runtime, LLMs provider, Web Components, iPaaS or [CI/CD](https://www.lambdatest.com/blog/best-ci-cd-tools/)
- 横切关注点：日志记录、安全性、事务管理、缓存、消息队列、有限状态机、特征与配置管理、[微服务架构框架](https://microservices.io/)
- 开发框架：使用全站框架(SvelteKit, T3, [remix](https://remix.run/)等)
- 可配置性：决策表、树、图、DSL 的灵活运用
- 插件：集成 CI/CD, 集成 Github, GitLab 等
- 开发语言：[Golang](https://golang.halfiisland.com/)
