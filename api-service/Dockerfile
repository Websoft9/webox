# Websoft9 API Service Dockerfile
# 多阶段构建，优化镜像大小和安全性

# 构建阶段
FROM golang:1.24-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git bash gcc musl-dev ca-certificates

# 创建非特权用户
RUN adduser -D -s /bin/sh appuser

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

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME} -w -s" \
    -a -installsuffix cgo \
    -o api-service \
    main

# 运行阶段
FROM alpine:latest

# 安装 CA 证书和时区数据
RUN apk --no-cache add ca-certificates bash sqlite

# 创建非特权用户
RUN adduser -D -s /bin/sh appuser

# 设置工作目录
WORKDIR /home/appuser/

# 从构建阶段复制二进制文件
COPY --from=builder /app/api-service .

# 复制配置文件（如果存在）
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/scripts ./scripts

# 创建数据目录
RUN mkdir -p data logs && chown -R appuser:appuser /home/appuser/

# 切换到非特权用户
USER appuser

# 暴露端口
EXPOSE 8080 9090

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 启动应用
CMD ["sh", "-c", "/app/scripts/init_db.sh && ./api-service"]