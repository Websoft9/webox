.PHONY: build run clean test deps docker-build docker-run init-postgres test-db health-check lint security-scan init-swag

# 构建应用
build:
	go build -v -o api-service ./cmd/server

# 构建用户创建工具
build-create-user:
	mkdir -p tools/create-user
	go build -v -o tools/create-user/create-user ./tools/create-user

# 运行应用
run:
	go run ./cmd/server

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

# 测试覆盖率
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 格式化代码
fmt:
	goimports -w .
	go fmt ./...

# 检查代码
vet:
	go vet ./...

# 代码静态分析（需要安装 golangci-lint）
lint:
	golangci-lint run --config ../.golangci.yml ./...

# 安全扫描（需要安装 gosec）
security-scan:
	gosec ./...

# 初始化swagger docs
init-swag:
	swag init -g cmd/server/main.go

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

# 数据库初始化（默认SQLite）
init-db:
	./scripts/init_db.sh --init ./scripts/init_data.sql

# MySQL数据库初始化
init-mysql:
	./scripts/init_db.sh -t mysql -u root -p $(or $(DB_PASSWORD),password) --init ./scripts/init_data.sql

# PostgreSQL数据库初始化
init-postgres:
	./scripts/init_db.sh -t postgres -u postgres -p $(or $(DB_PASSWORD),password) --init ./scripts/init_data.sql

# 测试数据库连接
test-db:
	@echo "Testing database connection..."
	@go run ./tools/test-db/main.go

# 重置SQLite数据库
reset-db:
	rm -f data/websoft9.db
	./scripts/init_db.sh --init ./scripts/init_data.sql

# 健康检查
health-check:
	@echo "Checking application health..."
	@curl -f http://localhost:8080/health || echo "Application not responding"
	@curl -f http://localhost:8080/health/database || echo "Database health check failed"

# 启动开发环境（包含数据库初始化）
setup-dev: deps init-db
	@echo "Development environment setup complete!"
	@echo "Run 'make dev' to start the development server"

# 启动生产环境
setup-prod: deps build init-db
	@echo "Production environment setup complete!"
	@echo "Run './api-service' to start the server"

# Docker Compose 相关命令
docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

# 数据库迁移
migrate:
	@echo "Running database migrations..."
	@go run ./tools/migrate-db/main.go

# 显示帮助信息
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run in development mode with hot reload"
	@echo "  clean         - Clean build files"
	@echo "  deps          - Download dependencies"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt           - Format code"
	@echo "  vet           - Check code"
	@echo "  lint          - Run static code analysis"
	@echo "  security-scan - Run security scan"
	@echo "  init-swag     - Initialize swagger docs"
	@echo "  init-db       - Initialize SQLite database"
	@echo "  init-mysql    - Initialize MySQL database"
	@echo "  init-postgres - Initialize PostgreSQL database"
	@echo "  test-db       - Test database connection"
	@echo "  reset-db      - Reset SQLite database"
	@echo "  migrate       - Run database migrations"
	@echo "  health-check  - Check application health"
	@echo "  setup-dev     - Setup development environment"
	@echo "  setup-prod    - Setup production environment"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
