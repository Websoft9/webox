# PR Check 工作流测试

这是一个测试 PR，用于验证 GitHub Actions 的 PR Check 工作流是否正常工作。

## 测试内容

### 1. 提交消息格式检查
- 使用符合 Conventional Commits 规范的提交消息
- 格式：`test(ci): add test files for PR check workflow validation`

### 2. 代码格式检查
- 添加了格式正确的 Go 代码
- 包含 `api-service` 和 `websoft9-agent` 两个组件

### 3. 单元测试
- 为两个组件都添加了基本的单元测试
- 测试覆盖率应该达到合理水平

### 4. 构建检查
- 添加了 `go.mod` 文件确保项目可以构建
- 代码结构符合 Go 项目规范

### 5. 安全扫描
- 代码中没有明显的安全问题
- Gosec 扫描应该通过

## 预期结果

所有检查项都应该通过：
- ✅ 提交消息检查
- ✅ 代码格式检查  
- ✅ 单元测试
- ✅ 构建检查
- ✅ 安全扫描

## 测试文件

### API Service
- `api-service/main.go` - 主要代码
- `api-service/main_test.go` - 单元测试
- `api-service/go.mod` - Go 模块定义

### Websoft9 Agent
- `websoft9-agent/agent.go` - 主要代码
- `websoft9-agent/agent_test.go` - 单元测试
- `websoft9-agent/go.mod` - Go 模块定义

这些文件提供了足够的代码来测试工作流的所有功能，同时保持简单和可维护。