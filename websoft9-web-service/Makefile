.PHONY: build run clean test deps docker-build docker-run

# 构建应用
build:
	go build -v -o websoft9-web-service .

# 运行应用
run:
	go run main.go

# 清理构建文件
clean:
	rm -f websoft9-web-service
	rm -f data/websoft9.db

# 下载依赖
deps:
	go mod tidy
	go mod download

# 运行测试
test:
	go test -v ./...

# 格式化代码
fmt:
	go fmt ./...

# 检查代码
vet:
	go vet ./...

# 构建Docker镜像
docker-build:
	docker build -t websoft9-web-service:latest .

# 运行Docker容器
docker-run:
	docker run -p 8080:8080 -p 9090:9090 websoft9-web-service:latest

# 开发模式（热重载需要安装air: go install github.com/air-verse/air@latest）
dev:
	air

# 初始化项目
init: deps build
	@echo "Project initialized successfully!"

# 完整测试
full-test: build
	./test_api.sh