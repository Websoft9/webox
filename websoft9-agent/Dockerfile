# Websoft9 Agent Dockerfile

# 构建阶段
FROM golang:1.24-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
ARG VERSION=dev
ARG COMMIT=1.0.0
ARG BUILD_TIME=2025-08-01T00:00:00Z

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -w -s" \
    -a -installsuffix cgo \
    -o websoft9-agent \
    ./cmd/agent

# 运行阶段
FROM alpine:latest

# 安装必要的工具
RUN apk --no-cache add ca-certificates tzdata curl docker-cli

# 注意：Agent 需要 root 权限来管理 Docker 和系统资源
# 在生产环境中，应该通过适当的权限管理来限制访问

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/websoft9-agent .

# 复制配置文件（如果存在）
COPY --from=builder /app/configs ./configs

# 创建必要的目录
RUN mkdir -p data logs

# 暴露端口（如果 Agent 需要监听端口）
EXPOSE 9091

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD pgrep websoft9-agent || exit 1

# 启动应用
CMD ["./websoft9-agent"]