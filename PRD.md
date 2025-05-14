# PRD

## 概要

- 产品名称：应用聚合与托管平台
- 价值体现：聚合/接入/运行/连接一切应用，打开 Websoft9 控制台，开启一天轻松基于 Web 随时随地办公模式
- 业务功能：寻找、部署、接入、运行、发布、管理、安全、订阅、咨询
- 技术功能：多用户多权限、Vault, Git, SSO, Variables, Orchestration, AI workfolw, Connection pipeline, Embeddeding applications/micro frontends
- 架构组件：前端(UI/CLI)、server, BFF 组件、API 网关，Agent, Controller, Event Stream, Storage, Runtime, LLMs provider, Web Components, iPaaS or [CI/CD](https://www.lambdatest.com/blog/best-ci-cd-tools/)
- 横切关注点：日志记录、安全性、事务管理、缓存、消息队列、有限状态机、特征与配置管理、[微服务架构框架](https://microservices.io/)
- 开发框架：使用全站框架(SvelteKit, T3, [remix](https://remix.run/)等)

## 业务需求

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

### 接入

将外部应用接入到 Websoft9，更方便快捷登录使用

### 运行

### 管理

### 发布

### 安全

### 智库

提供应用和科技相关知识，结合 RAG/LLM 供用户使用

## 技术需求

### 组网

多服务器组建虚拟局域网，实现集中管理。可能的技术方案：
- 中央服务器转发：VPN,SSH 隧道，流量完全走中央服务器
- P2P机制：NAT穿透+中央协调（headscale）
- CNI: flannel

### 微服务编排

将事务拆分为在多个微服务执行，暂定采用 n8n

### 凭证

- 存储微服务的凭证
- 更新凭证（事件/轮询/Webhook 三种机制之一）

### 事件消息

事件机制在软件架构中必不可少。






