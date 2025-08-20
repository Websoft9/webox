## 用户模块开发流程 - 模块协作详解

### 1. 请求流程 (以用户注册为例)

```
客户端请求 → 路由层 → 中间件 → 控制器 → 服务层 → 仓储层 → 数据库
     ↓         ↓        ↓        ↓       ↓       ↓        ↓
   JSON    路由匹配   验证/日志  参数绑定  业务逻辑  数据操作   持久化
     ↑         ↑        ↑        ↑       ↑       ↑        ↑
客户端响应 ← 响应格式 ← 错误处理 ← DTO转换 ← 模型转换 ← 查询结果 ← 数据返回
```

### 2. 详细协作流程

#### 步骤1: 客户端发起请求
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com", 
  "password": "123456"
}
```

#### 步骤2: 路由层处理 (router/router.go)
- 匹配路由规则: `/api/v1/auth/register`
- 调用对应的控制器方法: `controllers.UserController.Register`

#### 步骤3: 中间件处理
- **日志中间件**: 记录请求信息
- **CORS中间件**: 处理跨域
- **验证中间件**: 验证请求格式
- **错误处理中间件**: 统一错误处理

#### 步骤4: 控制器处理 (controller/user.go)
```go
func (c *UserController) Register(ctx *gin.Context) {
    // 1. 参数绑定和验证
    var req request.UserRegisterRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        // 返回验证错误
    }
    
    // 2. 调用服务层
    user, err := c.userService.Register(ctx, &req)
    if err != nil {
        // 错误处理
    }
    
    // 3. 返回响应
    pkg_response.Success(ctx, "用户注册成功", user)
}
```

#### 步骤5: 服务层处理 (service/user.go)
```go
func (s *userService) Register(ctx context.Context, req *request.UserRegisterRequest) (*response.UserResponse, error) {
    // 1. 业务验证
    if err := validator.ValidateUsername(req.Username); err != nil {
        return nil, err
    }
    
    // 2. 检查重复
    exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
    if exists {
        return nil, errors.NewAppError(errors.CodeUserExists, "用户名已存在")
    }
    
    // 3. 密码加密
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    
    // 4. 创建用户模型
    user := &model.User{
        Username: req.Username,
        Email:    req.Email,
        Password: string(hashedPassword),
        Status:   "active",
        Role:     "user",
    }
    
    // 5. 调用仓储层保存
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // 6. 转换为响应DTO
    return &response.UserResponse{
        ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
        Status:   user.Status,
        Role:     user.Role,
    }, nil
}
```

#### 步骤6: 仓储层处理 (repository/user.go)
```go
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
    // 执行数据库操作
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&model.User{}).
        Where("username = ?", username).Count(&count).Error
    return count > 0, err
}
```

### 3. 数据流转过程

#### 3.1 请求数据流
```
HTTP JSON → request.UserRegisterRequest → model.User → 数据库记录
```

#### 3.2 响应数据流  
```
数据库记录 → model.User → response.UserResponse → HTTP JSON
```

### 4. 错误处理流程

#### 4.1 验证错误
```
控制器 → 参数绑定失败 → 错误中间件 → 统一错误响应
```

#### 4.2 业务错误
```
服务层 → 业务验证失败 → 自定义错误 → 错误中间件 → 统一错误响应
```

#### 4.3 数据库错误
```
仓储层 → GORM错误 → 包装为业务错误 → 错误中间件 → 统一错误响应
```

### 5. 各层职责总结

| 层次 | 职责 | 主要工作 |
|------|------|----------|
| 路由层 | 请求分发 | URL匹配、路由组管理 |
| 中间件 | 横切关注点 | 认证、日志、CORS、错误处理 |
| 控制器 | HTTP处理 | 参数绑定、调用服务、响应格式 |
| 服务层 | 业务逻辑 | 业务验证、数据转换、事务协调 |
| 仓储层 | 数据访问 | CRUD操作、查询优化 |
| 模型层 | 数据结构 | 实体定义、业务方法 |
| DTO层 | 数据传输 | 请求验证、响应格式、数据隔离 |

### 6. 开发最佳实践

#### 6.1 分层原则
- 每层只关注自己的职责
- 上层可以调用下层，不能反向调用
- 通过接口实现解耦

#### 6.2 错误处理
- 在合适的层次处理错误
- 使用统一的错误类型
- 提供详细的错误信息

#### 6.3 日志记录
- 在关键节点记录日志
- 使用结构化日志
- 包含请求上下文信息

#### 6.4 数据验证
- 多层验证（参数验证、业务验证）
- 使用专门的验证器
- 提供友好的错误信息

### 7. 扩展性考虑

#### 7.1 新增功能
- 遵循现有分层结构
- 实现对应的接口
- 更新路由配置

#### 7.2 性能优化
- 仓储层添加缓存
- 服务层优化业务逻辑
- 数据库查询优化

#### 7.3 测试策略
- 单元测试：每层独立测试
- 集成测试：测试层间协作
- API测试：端到端测试
