.PHONY: build run clean test deps docker-build docker-run

# 构建应用
build:
	go build -v -o api-service .

# 运行应用
run:
	go run main.go

# 清理构建文件
clean:
	rm -f api-service
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
	docker build -t api-service:latest .

# 运行Docker容器
docker-run:
	docker run -p 8080:8080 -p 9090:9090 api-service:latest

# 开发模式（热重载需要安装air: go install github.com/air-verse/air@latest）
dev:
	air

# 初始化项目
init: deps build
	@echo "Project initialized successfully!"

# 完整测试
full-test: build
	./test_api.sh

# 数据库初始化
init-db:
	./scripts/init_db.sh

# MySQL数据库初始化
init-mysql:
	./scripts/init_db.sh -t mysql -u root -p password

# 重置数据库
reset-db:
	rm -f data/websoft9.db
	./scripts/init_db.sh