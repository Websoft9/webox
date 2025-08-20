# Websoft9 用户管理 API 文档

## 概述

本文档描述了 Websoft9 平台用户管理模块的 REST API 接口。用户管理模块负责用户账户的全生命周期管理，包括用户创建、信息维护、状态管理、查询筛选等功能。

## 认证授权

所有用户管理接口（除了注册和登录）都需要 JWT Token 认证。

### 请求头格式

```http
Authorization: Bearer <JWT_TOKEN>
```

## 接口列表

### 1. 用户注册

**接口地址：** `POST /api/v1/auth/register`

**功能描述：** 用户注册

**请求参数：**

```json
{
  "username": "testuser",         // 用户名，必填，3-32字符
  "email": "test@example.com",    // 邮箱，必填，唯一
  "password": "TestPass123"       // 密码，必填，6位以上
}
```

**响应示例：**

```json
{
  "code": 200,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2025-08-20T10:30:00Z"
  }
}
```

### 2. 用户登录

**接口地址：** `POST /api/v1/auth/login`

**功能描述：** 用户登录

**请求参数：**

```json
{
  "username": "testuser",    // 用户名，必填
  "password": "TestPass123"  // 密码，必填
}
```

**响应示例：**

```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 3. 获取用户资料

**接口地址：** `GET /api/v1/users/profile`

**功能描述：** 获取当前登录用户的资料

**请求参数：** 无

**响应示例：**

```json
{
  "code": 200,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "group_id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar": "https://example.com/avatars/user.jpg",
    "phone": "13800138000",
    "gender": 1,
    "signature": "这是我的个性签名",
    "status": 1,
    "timezone": "Asia/Shanghai",
    "language": "zh-CN",
    "roles": [
      {
        "id": 1,
        "name": "管理员",
        "code": "admin"
      }
    ],
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z"
  }
}
```

### 4. 获取用户列表

**接口地址：** `GET /api/v1/users`

**功能描述：** 获取用户列表，支持分页、搜索、筛选、排序

**查询参数：**

| 参数      | 类型    | 必填 | 默认值     | 描述                             |
| --------- | ------- | ---- | ---------- | -------------------------------- |
| page      | integer | 否   | 1          | 页码                             |
| page_size | integer | 否   | 20         | 每页数量（1-100）                |
| keyword   | string  | 否   | -          | 搜索关键词（用户名、邮箱、昵称、手机号） |
| status    | integer | 否   | -          | 用户状态（0-禁用，1-启用，-1-已删除） |
| role_id   | integer | 否   | -          | 角色ID筛选                       |
| sort      | string  | 否   | created_at | 排序字段                         |
| order     | string  | 否   | desc       | 排序方向（asc, desc）            |

**请求示例：**

```http
GET /api/v1/users?page=1&page_size=10&keyword=test&status=1&sort=created_at&order=desc
```

**响应示例：**

```json
{
  "code": 200,
  "message": "Users retrieved successfully",
  "data": {
    "items": [
      {
        "id": 1,
        "group_id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "nickname": "系统管理员",
        "avatar": "https://example.com/avatars/admin.jpg",
        "phone": "13800138000",
        "gender": 1,
        "signature": "系统管理员账户",
        "status": 1,
        "timezone": "Asia/Shanghai",
        "language": "zh-CN",
        "last_login_at": "2025-08-20T10:30:00Z",
        "last_login_ip": "192.168.1.100",
        "roles": [
          {
            "id": 1,
            "name": "系统管理员",
            "code": "admin"
          }
        ],
        "created_at": "2025-08-20T10:30:00Z",
        "updated_at": "2025-08-20T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 1,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

### 5. 获取用户详情

**接口地址：** `GET /api/v1/users/{id}`

**功能描述：** 根据用户ID获取用户详细信息

**路径参数：**

| 参数 | 类型    | 必填 | 描述   |
| ---- | ------- | ---- | ------ |
| id   | integer | 是   | 用户ID |

**响应示例：**

```json
{
  "code": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "group_id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "nickname": "系统管理员",
    "avatar": "https://example.com/avatars/admin.jpg",
    "phone": "13800138000",
    "gender": 1,
    "signature": "系统管理员账户",
    "status": 1,
    "timezone": "Asia/Shanghai",
    "language": "zh-CN",
    "roles": [
      {
        "id": 1,
        "name": "系统管理员",
        "code": "admin"
      }
    ],
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z"
  }
}
```

### 6. 创建用户

**接口地址：** `POST /api/v1/users`

**功能描述：** 创建新用户（需要管理员权限）

**请求参数：**

| 参数名    | 数据类型 | 是否可空 | 描述                                 |
| --------- | -------- | -------- | ------------------------------------ |
| group_id  | integer  | 否       | 用户组ID                             |
| username  | string   | 否       | 用户名，3-32字符，字母数字下划线     |
| email     | string   | 否       | 邮箱地址，必须唯一                   |
| password  | string   | 否       | 密码，8-32字符，包含大小写字母和数字 |
| nickname  | string   | 是       | 用户昵称                             |
| phone     | string   | 是       | 手机号码                             |
| gender    | integer  | 是       | 性别（1-男，2-女，0-未知）           |
| signature | string   | 是       | 个性签名                             |
| timezone  | string   | 是       | 时区，默认Asia/Shanghai              |
| language  | string   | 是       | 语言，默认zh-CN                      |
| role_ids  | array    | 是       | 角色ID数组                           |
| status    | integer  | 是       | 状态：0-禁用，1-启用                 |

**请求示例：**

```json
{
  "group_id": 1,
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "NewPass123",
  "nickname": "新用户",
  "phone": "13800138001",
  "gender": 1,
  "signature": "这是我的个性签名",
  "timezone": "Asia/Shanghai",
  "language": "zh-CN",
  "role_ids": [2],
  "status": 1
}
```

**响应示例：**

```json
{
  "code": 200,
  "message": "User created successfully",
  "data": {
    "id": 2,
    "group_id": 1,
    "username": "newuser",
    "email": "newuser@example.com",
    "nickname": "新用户",
    "phone": "13800138001",
    "gender": 1,
    "signature": "这是我的个性签名",
    "status": 1,
    "timezone": "Asia/Shanghai",
    "language": "zh-CN",
    "roles": [
      {
        "id": 2,
        "name": "普通用户",
        "code": "user"
      }
    ],
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z"
  }
}
```

### 7. 更新用户

**接口地址：** `PUT /api/v1/users/{id}`

**功能描述：** 更新用户信息

**路径参数：**

| 参数 | 类型    | 必填 | 描述   |
| ---- | ------- | ---- | ------ |
| id   | integer | 是   | 用户ID |

**请求参数：**（所有参数都是可选的）

| 参数名    | 数据类型 | 是否可空 | 描述                       |
| --------- | -------- | -------- | -------------------------- |
| group_id  | integer  | 是       | 用户组ID                   |
| nickname  | string   | 是       | 用户昵称                   |
| email     | string   | 是       | 邮箱地址，必须唯一         |
| phone     | string   | 是       | 手机号码                   |
| gender    | integer  | 是       | 性别（1-男，2-女，0-未知） |
| signature | string   | 是       | 个性签名                   |
| timezone  | string   | 是       | 时区                       |
| language  | string   | 是       | 语言                       |
| role_ids  | array    | 是       | 角色ID数组                 |
| status    | integer  | 是       | 状态：0-禁用，1-启用       |

**请求示例：**

```json
{
  "nickname": "更新的昵称",
  "signature": "更新的个性签名",
  "status": 1
}
```

**响应示例：**

```json
{
  "code": 200,
  "message": "User updated successfully",
  "data": {
    "id": 2,
    "group_id": 1,
    "username": "newuser",
    "email": "newuser@example.com",
    "nickname": "更新的昵称",
    "phone": "13800138001",
    "gender": 1,
    "signature": "更新的个性签名",
    "status": 1,
    "timezone": "Asia/Shanghai",
    "language": "zh-CN",
    "roles": [
      {
        "id": 2,
        "name": "普通用户",
        "code": "user"
      }
    ],
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T11:00:00Z"
  }
}
```

### 8. 删除用户

**接口地址：** `DELETE /api/v1/users/{id}`

**功能描述：** 删除用户（软删除，设置状态为-1）

**路径参数：**

| 参数 | 类型    | 必填 | 描述   |
| ---- | ------- | ---- | ------ |
| id   | integer | 是   | 用户ID |

**响应示例：**

```json
{
  "code": 200,
  "message": "User deleted successfully",
  "data": null
}
```

### 9. 修改用户密码

**接口地址：** `PUT /api/v1/users/{id}/password`

**功能描述：** 修改用户密码

**路径参数：**

| 参数 | 类型    | 必填 | 描述   |
| ---- | ------- | ---- | ------ |
| id   | integer | 是   | 用户ID |

**请求参数：**

| 参数名           | 数据类型 | 是否可空 | 描述       |
| ---------------- | -------- | -------- | ---------- |
| old_password     | string   | 否       | 原密码     |
| new_password     | string   | 否       | 新密码     |
| confirm_password | string   | 否       | 确认新密码 |

**请求示例：**

```json
{
  "old_password": "OldPass123",
  "new_password": "NewPass123",
  "confirm_password": "NewPass123"
}
```

**响应示例：**

```json
{
  "code": 200,
  "message": "Password changed successfully",
  "data": null
}
```

## 错误响应

### 通用错误格式

```json
{
  "code": 400,
  "message": "错误描述",
  "data": "详细错误信息"
}
```

### 常见错误码

| 错误码 | 描述           | 示例                                   |
| ------ | -------------- | -------------------------------------- |
| 400    | 请求参数错误   | 参数验证失败、格式错误等               |
| 401    | 未授权         | Token无效、Token过期等                 |
| 403    | 权限不足       | 没有权限访问该资源                     |
| 404    | 资源不存在     | 用户不存在                             |
| 409    | 资源冲突       | 用户名或邮箱已存在                     |
| 500    | 服务器内部错误 | 数据库连接失败、系统异常等             |

## 数据字典

### 用户状态（status）

| 值  | 描述   |
| --- | ------ |
| -1  | 已删除 |
| 0   | 禁用   |
| 1   | 启用   |

### 性别（gender）

| 值  | 描述 |
| --- | ---- |
| 0   | 未知 |
| 1   | 男   |
| 2   | 女   |

### 排序字段（sort）

支持的排序字段：
- `id` - 用户ID
- `username` - 用户名
- `email` - 邮箱
- `created_at` - 创建时间
- `updated_at` - 更新时间
- `last_login_at` - 最后登录时间

### 排序方向（order）

| 值   | 描述 |
| ---- | ---- |
| asc  | 升序 |
| desc | 降序 |

## 注意事项

1. **密码复杂度要求：** 密码必须包含大写字母、小写字母和数字，长度8-32位
2. **用户名规则：** 用户名只能包含字母、数字和下划线，长度3-32位
3. **邮箱唯一性：** 邮箱地址必须在系统中唯一
4. **手机号格式：** 支持中国大陆手机号格式（1开头的11位数字）
5. **删除用户：** 删除操作是软删除，用户状态会被设置为-1，但数据仍保留在数据库中
6. **权限控制：** 创建、更新、删除用户操作需要管理员权限
7. **分页限制：** 每页最多返回100条记录
8. **搜索范围：** 关键词搜索会在用户名、邮箱、昵称、手机号字段中进行模糊匹配
