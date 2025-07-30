# Websoft9 Agent

Websoft9 Agent 是部署在各服务器节点的客户端程序，负责执行服务端下发的任务指令和采集监控数据。

## 功能特性

- **Monitor** - 服务器/应用监控及数据采集
  - 系统资源监控（CPU、内存、磁盘、网络）
  - 容器应用监控（容器状态、资源使用）
  - 应用健康检查（HTTP检查、TCP检查、自定义脚本）
  - 日志采集和转发

- **Task** - 服务端任务指令执行
  - 应用部署和管理操作
  - 系统命令执行
  - 文件传输和管理
  - 系统服务管理

- **Workflows** - 工作流任务调度
  - 工作流任务的本地执行
  - 任务状态反馈和日志上报
  - 任务超时和异常处理

- **Communication** - 通信管理
  - 与服务端的gRPC通信
  - 心跳保持和状态上报
  - 消息队列事件处理

## 技术栈

- 开发语言：Golang
- gRPC：grpc-go
- Redis：go-redis
- InfluxDB：influxdb-client-go
- Docker API：Docker Engine API

## 部署要求

- 必须使用 root 用户启动
- 以系统服务方式运行
- 支持容器化部署

## 构建和运行

```bash
# 构建
go build -o websoft9-agent ./cmd/agent

# 运行
./websoft9-agent
```
