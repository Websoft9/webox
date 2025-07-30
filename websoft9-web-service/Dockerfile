# This action is used for building webox-server to DockerHub at one time
# It need LABEL version="" at your dockerfile.
# modify time: 202507191210, you can modify here to trigger Docker Build action


FROM golang:1.24-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git bash gcc musl-dev

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod tidy

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

#ENV REDIS_HOST=websoft9-redis
#ENV REDIS_PORT=6379
#ENV JWT_SECRET="570ef696bca45877d1b5da0ede0177703d2f1079f31a1332607a1ad237e5d356"

# 暴露端口
EXPOSE 8080 9090

# 初始化数据库并运行应用
CMD ["sh", "-c", "/root/scripts/init_db.sh && ./main"]