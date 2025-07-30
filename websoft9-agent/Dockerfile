# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o websoft9-agent ./cmd/agent

# 运行阶段
FROM alpine:latest

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/websoft9-agent .

# 复制配置文件
COPY --from=builder /app/configs/agent.yaml /etc/websoft9/agent.yaml

# 创建工作目录
RUN mkdir -p /var/lib/websoft9/agent
RUN mkdir -p /var/log/websoft9

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口（如果需要）
# EXPOSE 8080

# 运行应用
CMD ["./websoft9-agent", "-config", "/etc/websoft9/agent.yaml"]