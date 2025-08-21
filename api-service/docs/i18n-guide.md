# API Service 国际化 (i18n) 功能说明

## 概述

本项目已成功集成了完整的国际化(i18n)功能，支持中文和英文两种语言，默认使用英文。i18n功能已全面集成到错误处理、用户认证和响应消息中。

## 特性

- ✅ 支持英文(en)和中文(zh)两种语言
- ✅ 默认语言为英文
- ✅ 自动从HTTP Header检测语言偏好
- ✅ 支持查询参数和自定义Header指定语言
- ✅ 统一的错误处理和响应消息多语言支持
- ✅ 嵌入式语言包，无需外部文件
- ✅ 提供i18n管理和测试API端点

## 项目结构

```
pkg/i18n/
├── i18n.go              # 核心i18n实现
└── locales/             # 语言文件
    ├── en.yaml          # 英文语言包
    └── zh.yaml          # 中文语言包

internal/middleware/
└── i18n.go             # i18n中间件

internal/controller/
└── i18n.go             # i18n管理API控制器
```

## 语言检测优先级

1. 查询参数 `lang`：`/api/v1/users?lang=zh`
2. 自定义Header `X-Language`：`X-Language: zh`
3. 标准Header `Accept-Language`：`Accept-Language: zh-CN,zh;q=0.9,en;q=0.8`
4. 默认语言：`en`

## API使用示例

### 1. 语言管理端点

#### 获取支持的语言列表
```bash
GET /api/v1/i18n/languages

# 响应
{
  "code": 200,
  "message": "Operation successful",
  "data": {
    "default": "en",
    "languages": [
      {
        "code": "en",
        "name": "English",
        "native_name": "English",
        "supported": true,
        "default": true
      },
      {
        "code": "zh", 
        "name": "Chinese",
        "native_name": "中文",
        "supported": true,
        "default": false
      }
    ]
  }
}
```

#### 获取指定语言的翻译
```bash
GET /api/v1/i18n/translations/zh

# 响应包含该语言的翻译示例
```

#### 测试i18n功能
```bash
GET /api/v1/i18n/test
Accept-Language: zh-CN

# 响应包含各种场景的多语言消息测试
```

### 2. 用户认证API多语言示例

#### 用户注册 - 中文
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Accept-Language: zh-CN" \
  -d '{"username":"user","password":"pass123","email":"user@example.com"}' \
  http://localhost:8080/api/v1/auth/register

# 成功响应
{
  "code": 200,
  "message": "用户创建成功",
  "data": {...}
}

# 错误响应
{
  "code": 400,
  "message": "数据验证失败"
}
```

#### 用户注册 - 英文
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Accept-Language: en" \
  -d '{"username":"user","password":"pass123","email":"user@example.com"}' \
  http://localhost:8080/api/v1/auth/register

# 成功响应
{
  "code": 200,
  "message": "User created successfully", 
  "data": {...}
}

# 错误响应
{
  "code": 400,
  "message": "Validation failed"
}
```

## 语言包配置

### 英文语言包 (en.yaml)
```yaml
user:
  not_found: "User not found"
  created_success: "User created successfully"
  login_success: "Login successful"
  invalid_credentials: "Invalid username or password"

auth:
  unauthorized: "Unauthorized access"
  token_expired: "Token has expired"

common:
  success: "Operation successful"
  validation_failed: "Validation failed"
  not_found: "Resource not found"
```

### 中文语言包 (zh.yaml)
```yaml
user:
  not_found: "用户不存在"
  created_success: "用户创建成功"
  login_success: "登录成功"
  invalid_credentials: "用户名或密码错误"

auth:
  unauthorized: "未授权访问"
  token_expired: "令牌已过期"

common:
  success: "操作成功"
  validation_failed: "数据验证失败"
  not_found: "资源不存在"
```

## 开发者使用指南

### 在Controller中使用i18n

```go
import "api-service/internal/middleware"

func (c *UserController) SomeMethod(ctx *gin.Context) {
    // 获取翻译后的消息
    successMsg := middleware.T(ctx, "user.created_success")
    errorMsg := middleware.T(ctx, "user.not_found")
    
    // 在响应中使用
    pkg_response.Success(ctx, successMsg, data)
}
```

### 在Service中使用i18n

```go
import "api-service/pkg/i18n"

func (s *userService) SomeMethod(ctx context.Context, lang string) error {
    // 直接使用翻译函数
    message := i18n.T("user.not_found", lang)
    return errors.NewAppErrorWithI18n(errors.CodeUserNotFound, message, "user.not_found")
}
```

### 添加新的翻译

1. 在 `pkg/i18n/locales/en.yaml` 中添加英文翻译
2. 在 `pkg/i18n/locales/zh.yaml` 中添加中文翻译
3. 重新编译应用程序

### 错误处理中的i18n

所有错误现在都支持i18n：

```go
// 使用预定义的i18n错误
return errors.ErrUserNotFound

// 创建带i18n的自定义错误
return errors.NewAppErrorWithI18n(
    errors.CodeUserNotFound, 
    "User not found", 
    "user.not_found"
)
```

## 配置

在 `configs/config.yaml` 中可以配置i18n设置：

```yaml
i18n:
  default_language: "en"
  supported_languages: ["en", "zh"]
```

## 技术细节

- 使用 `github.com/nicksnyder/go-i18n/v2` 库
- 语言文件使用YAML格式，编译时嵌入二进制文件
- 支持复数形式和模板变量
- 通过中间件自动检测和设置语言
- 统一的错误处理支持多语言

## 测试验证

已通过以下测试验证：

1. ✅ 语言检测和切换
2. ✅ 用户注册/登录的多语言响应
3. ✅ 错误消息的多语言显示
4. ✅ API端点的i18n功能
5. ✅ 默认语言回退机制

项目现在完全支持中英文双语言，为后续国际化扩展奠定了坚实基础。
