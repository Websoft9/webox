FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git bash

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM alpine:latest

LABEL maintainer="Websoft9<help@websoft9.com>"
LABEL version="0.1.1"

# 安装运行时依赖，包括 sqlite3
RUN apk --no-cache add ca-certificates bash sqlite

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/scripts ./scripts

# 创建数据目录
RUN mkdir -p /root/data
VOLUME /root/data

# 暴露端口
EXPOSE 8080 9090

# 初始化数据库并运行应用
CMD ["sh", "-c", "/root/scripts/init_db.sh && ./main"]