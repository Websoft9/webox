# Websoft9 API Service Dockerfile
# 多阶段构建，优化镜像大小和安全性

# 定义InfluxDB版本变量
ARG INFLUXDB_VERSION=2.7.12

# 构建阶段
FROM golang:1.24-alpine AS builder

# 声明变量
ARG INFLUXDB_VERSION

# 安装必要的工具
RUN apk add --no-cache git bash gcc musl-dev ca-certificates curl

# 创建非特权用户
RUN adduser -D -s /bin/sh appuser

# 设置工作目录
WORKDIR /app

# 下载并解压 InfluxDB 二进制
RUN curl -L https://dl.influxdata.com/influxdb/releases/v${INFLUXDB_VERSION}/influxdb${INFLUXDB_VERSION%%.*}-${INFLUXDB_VERSION}_linux_amd64.tar.gz | tar -xz

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
    .

# 运行阶段
FROM alpine:latest

LABEL version="0.1.2"

# 声明变量
ARG INFLUXDB_VERSION

# 安装 CA 证书和时区数据
RUN apk --no-cache add ca-certificates bash sqlite curl supervisor redis

# 创建非特权用户
RUN adduser -D -s /bin/sh appuser

# 设置工作目录
WORKDIR /home/appuser/

# 从构建阶段复制二进制文件
COPY --from=builder /app/api-service .

# 复制配置文件（如果存在）
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/scripts ./scripts
RUN chmod +x ./scripts/init_db.sh ./scripts/init_influxdb.sh ./scripts/healthcheck.sh

# 从构建阶段复制 influxd
COPY --from=builder /app/influxdb${INFLUXDB_VERSION%%.*}-${INFLUXDB_VERSION}/usr/bin/influxd /usr/local/bin/

# 创建数据目录
RUN mkdir -p /home/appuser/data /home/appuser/logs /home/appuser/data/influxdb && chown -R appuser:appuser /home/appuser/

# 复制supervisord 配置文件和健康检查脚本
COPY supervisord.conf /etc/supervisord.conf.d/supervisord.conf

# 切换到非特权用户
USER appuser

# 暴露端口
EXPOSE 8080 9090 6379 8086

# 启动应用
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf.d/supervisord.conf"]