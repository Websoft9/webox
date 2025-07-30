# Websoft9 Agent Makefile

# 变量定义
BINARY_NAME=websoft9-agent
MAIN_PATH=./cmd/agent
BUILD_DIR=./build
VERSION?=v1.0.0
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S)

# Go 相关变量
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 构建标志
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all build clean test deps docker help

# 默认目标
all: clean deps test build

# 构建二进制文件
build:
	@echo "构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# 构建 Linux 版本
build-linux:
	@echo "构建 Linux 版本..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)

# 构建所有平台版本
build-all: build-linux
	@echo "构建 Windows 版本..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe $(MAIN_PATH)
	@echo "构建 macOS 版本..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin $(MAIN_PATH)

# 清理构建文件
clean:
	@echo "清理构建文件..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# 运行测试
test:
	@echo "运行测试..."
	$(GOTEST) -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# 下载依赖
deps:
	@echo "下载依赖..."
	$(GOMOD) download
	$(GOMOD) tidy

# 更新依赖
update-deps:
	@echo "更新依赖..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# 构建 Docker 镜像
docker:
	@echo "构建 Docker 镜像..."
	docker build -t websoft9/agent:$(VERSION) .
	docker build -t websoft9/agent:latest .

# 运行应用
run:
	@echo "运行应用..."
	$(BUILD_DIR)/$(BINARY_NAME) -config ./configs/agent.yaml

# 安装到系统
install: build
	@echo "安装到系统..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	sudo mkdir -p /etc/websoft9
	sudo cp configs/agent.yaml /etc/websoft9/

# 格式化代码
fmt:
	@echo "格式化代码..."
	$(GOCMD) fmt ./...

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run

# 生成 mock 文件
mock:
	@echo "生成 mock 文件..."
	mockery --all --output ./mocks

# 显示帮助信息
help:
	@echo "可用的 make 目标:"
	@echo "  build        - 构建二进制文件"
	@echo "  build-linux  - 构建 Linux 版本"
	@echo "  build-all    - 构建所有平台版本"
	@echo "  clean        - 清理构建文件"
	@echo "  test         - 运行测试"
	@echo "  test-coverage- 运行测试并生成覆盖率报告"
	@echo "  deps         - 下载依赖"
	@echo "  update-deps  - 更新依赖"
	@echo "  docker       - 构建 Docker 镜像"
	@echo "  run          - 运行应用"
	@echo "  install      - 安装到系统"
	@echo "  fmt          - 格式化代码"
	@echo "  lint         - 代码检查"
	@echo "  mock         - 生成 mock 文件"
	@echo "  help         - 显示此帮助信息"