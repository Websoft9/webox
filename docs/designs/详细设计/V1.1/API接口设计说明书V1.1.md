# Websoft9 API接口设计说明书 V1.1

**目录**

- [Websoft9 API接口设计说明书 V1.1](#websoft9-api接口设计说明书-v11)
  - [1. 引言](#1-引言)
  - [2. API设计原则](#2-api设计原则)
    - [2.1 RESTful设计规范](#21-restful设计规范)
    - [2.2 URL设计规范](#22-url设计规范)
    - [2.3 统一响应格式](#23-统一响应格式)
      - [成功响应格式](#成功响应格式)
      - [分页响应格式](#分页响应格式)
      - [错误响应格式](#错误响应格式)
  - [3. 认证与授权](#3-认证与授权)
    - [3.1 JWT Token认证](#31-jwt-token认证)
      - [用户登录](#用户登录)
      - [Token使用](#token使用)
      - [Token刷新](#token刷新)
      - [用户登出](#用户登出)
      - [获取当前用户信息](#获取当前用户信息)
    - [3.2 权限控制](#32-权限控制)
  - [4. 接口详细设计](#4-接口详细设计)
    - [4.1 平台主页](#41-平台主页)
      - [4.1.1 项目总览看板](#411-项目总览看板)
      - [4.1.2 监控总览看板](#412-监控总览看板)
      - [4.1.3 应用快捷导航](#413-应用快捷导航)
    - [4.2 项目管理](#42-项目管理)
      - [4.2.1 监控看板](#421-监控看板)
        - [应用监控看板](#应用监控看板)
        - [服务器监控看板](#服务器监控看板)
      - [4.2.2 任务看板](#422-任务看板)
      - [4.2.3 资源看板](#423-资源看板)
      - [4.2.4 文件夹](#424-文件夹)
        - [通用文件上传下载接口](#通用文件上传下载接口)
      - [4.2.5 应用](#425-应用)
      - [4.2.6 工作流](#426-工作流)
        - [工作流执行](#工作流执行)
        - [工作流组件](#工作流组件)
        - [任务](#任务)
      - [4.2.7 资源组](#427-资源组)
      - [4.2.8 服务器](#428-服务器)
      - [4.2.9 密钥管理](#429-密钥管理)
      - [4.2.10 数据库](#4210-数据库)
      - [4.2.11 应用网关](#4211-应用网关)
      - [4.2.12 证书管理](#4212-证书管理)
      - [4.2.13 云资源](#4213-云资源)
      - [4.2.14 项目团队](#4214-项目团队)
      - [4.2.15 项目设置](#4215-项目设置)
    - [4.3 应用市场](#43-应用市场)
      - [4.3.1 应用市场列表](#431-应用市场列表)
        - [应用分类目录](#应用分类目录)
        - [应用列表](#应用列表)
        - [应用部署](#应用部署)
      - [4.3.2 应用心愿单](#432-应用心愿单)
        - [用户收藏和评价](#用户收藏和评价)
    - [4.4 平台管理](#44-平台管理)
      - [4.4.1 项目管理](#441-项目管理)
      - [4.4.2 平台设置](#442-平台设置)
      - [4.4.3 安全管理](#443-安全管理)
        - [角色管理](#角色管理)
        - [权限管理](#权限管理)
        - [认证管理](#认证管理)
      - [4.4.4 用户管理](#444-用户管理)
      - [4.4.5 告警通知](#445-告警通知)
      - [4.4.6 个人中心](#446-个人中心)
      - [4.4.7 审计日志](#447-审计日志)
  - [5. 接口错误处理](#5-接口错误处理)
    - [5.1 标准错误码](#51-标准错误码)
    - [5.2 详细错误码定义](#52-详细错误码定义)
      - [5.2.1 认证相关错误（1000-1999）](#521-认证相关错误1000-1999)
      - [5.2.2 权限相关错误（2000-2999）](#522-权限相关错误2000-2999)
      - [5.2.3 参数验证错误（3000-3999）](#523-参数验证错误3000-3999)
      - [5.2.4 资源相关错误（4000-4999）](#524-资源相关错误4000-4999)
      - [5.2.5 业务逻辑错误（5000-5999）](#525-业务逻辑错误5000-5999)
      - [5.2.6 系统相关错误（6000-6999）](#526-系统相关错误6000-6999)
    - [5.3 错误响应示例](#53-错误响应示例)
      - [5.3.1 参数验证错误](#531-参数验证错误)
      - [5.3.2 业务逻辑错误](#532-业务逻辑错误)
      - [5.3.3 系统错误](#533-系统错误)
  - [6. 接口安全设计](#6-接口安全设计)
    - [6.1 输入验证](#61-输入验证)
      - [6.1.1 参数验证规则](#611-参数验证规则)
      - [6.1.2 SQL注入防护](#612-sql注入防护)
      - [6.1.3 XSS防护](#613-xss防护)
    - [6.2 认证与授权](#62-认证与授权)
      - [6.2.1 JWT Token安全](#621-jwt-token安全)
      - [6.2.2 权限控制](#622-权限控制)
      - [6.2.3 会话管理](#623-会话管理)
    - [6.3 数据加密](#63-数据加密)
      - [6.3.1 传输加密](#631-传输加密)
      - [6.3.2 存储加密](#632-存储加密)
      - [6.3.3 密钥管理](#633-密钥管理)

## 1. 引言

本文档为 **Websoft9架构升级** 的API接口设计文档，本文档的编写目的在于明确 **Websoft9架构升级** 需求的开发途径以及应用方法，通过此文档，为此 **Websoft9架构升级** 需求的维护提供清晰、详细的设计，为下阶段开发工作的开展起到指导作用。

本文档的预期读者是 **Websoft9架构升级** 需求相关的后端开发人员、前端开发人员以及第三方集成开发者。

## 2. API设计原则

### 2.1 RESTful设计规范

- **资源导向**：URL表示资源，使用名词而非动词
- **HTTP动词**：使用标准HTTP方法表示操作
  - GET：查询资源
  - POST：创建资源
  - PUT：更新资源（完整更新）
  - PATCH：更新资源（部分更新）
  - DELETE：删除资源
- **状态码**：使用标准HTTP状态码表示结果
- **无状态**：每个请求包含所有必要信息

### 2.2 URL设计规范

```text
https://api.websoft9.com/v1/{resource}[/{id}][/{sub-resource}]
```

示例：

- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/123` - 获取ID为123的用户
- `GET /api/v1/users/123/roles` - 获取用户123的角色列表
- `POST /api/v1/servers/456/actions` - 对服务器456执行操作

### 2.3 统一响应格式

#### 成功响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    // 响应数据
  }
}
```

#### 分页响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      // 数据项列表
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

#### 错误响应格式

```json
{
  "code": 400,
  "message": "请求参数错误",
  "error": {
    "type": "VALIDATION_ERROR",
    "code": "INVALID_PARAMETER",
    "details": [
      {
        "field": "email",
        "message": "邮箱格式不正确",
        "code": "INVALID_EMAIL_FORMAT"
      }
    ]
  }
}
```

## 3. 认证与授权

### 3.1 JWT Token认证

#### 用户登录

```text
POST /api/v1/auth/login
```

请求体参数：

| 参数名   | 数据类型 | 是否可空 | 描述             |
| -------- | -------- | -------- | ---------------- |
| username | string   | 否       | 用户名           |
| password | string   | 否       | 密码（加密结果） |

响应：

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "nickname": "管理员",
      "roles": ["admin"]
    }
  }
}
```

#### Token使用

所有需要认证的API请求必须在Header中携带Token：

```text
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Token刷新

```text
POST /api/v1/auth/refresh
```

请求体参数：

| 参数名        | 数据类型 | 是否可空 | 描述      |
| ------------- | -------- | -------- | --------- |
| refresh_token | string   | 否       | 刷新Token |

#### 用户登出

```text
POST /api/v1/auth/logout
```

#### 获取当前用户信息

```text
GET /api/v1/auth/me
```

### 3.2 权限控制

基于RBAC模型，每个API接口定义所需权限：

```json
{
  "endpoint": "GET /api/v1/users",
  "required_permissions": ["user:read"],
  "resource_level": false
}
```

## 4. 接口详细设计

### 4.1 平台主页

#### 4.1.1 项目总览看板

**获取项目总览看板数据**

```text
GET /api/v1/dashboard/project-overview
```

查询参数：

| 参数       | 类型   | 必填 | 默认值 | 描述                         |
| ---------- | ------ | ---- | ------ | ---------------------------- |
| time_range | string | 否   | 30d    | 时间范围（1d, 7d, 30d, 90d） |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "project_statistics": {
      "total_projects": 15,
      "active_projects": 12,
      "paused_projects": 2,
      "error_projects": 1,
      "total_resources": 125,
      "total_applications": 45,
      "storage_used": 512000000000
    },
    "resource_overview": {
      "servers": {
        "total": 25,
        "online": 23,
        "offline": 2
      },
      "applications": {
        "total": 45,
        "running": 38,
        "stopped": 5,
        "error": 2
      },
      "storage_usage": {
        "used": 512,
        "total": 1024,
        "unit": "GB",
        "usage_percentage": 50.0
      }
    },
    "recent_activities": [
      {
        "id": 1,
        "type": "APP_DEPLOYMENT",
        "title": "部署应用：WordPress博客",
        "description": "在服务器Web-01上成功部署WordPress应用",
        "project_id": 1,
        "project_name": "生产环境",
        "user": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员"
        },
        "status": "SUCCESS",
        "created_at": "2025-01-22T10:30:00Z"
      },
      {
        "id": 2,
        "type": "SERVER_ADDED",
        "title": "添加服务器：Web-02",
        "description": "成功添加新服务器Web-02到生产环境",
        "project_id": 1,
        "project_name": "生产环境",
        "user": {
          "id": 2,
          "username": "operator",
          "nickname": "运维员"
        },
        "status": "SUCCESS",
        "created_at": "2025-01-22T09:15:00Z"
      }
    ],
    "project_quick_access": [
      {
        "id": 1,
        "name": "生产环境",
        "description": "生产环境项目",
        "status": "ACTIVE",
        "resource_count": 25,
        "app_count": 15,
        "last_accessed": "2025-01-22T10:30:00Z"
      },
      {
        "id": 2,
        "name": "测试环境",
        "description": "测试环境项目",
        "status": "ACTIVE",
        "resource_count": 10,
        "app_count": 8,
        "last_accessed": "2025-01-22T09:45:00Z"
      }
    ],
    "pending_tasks": [
      {
        "id": 1,
        "type": "CERTIFICATE_RENEWAL",
        "title": "SSL证书即将过期",
        "description": "example.com的SSL证书将在7天后过期，请及时续期",
        "priority": "HIGH",
        "due_date": "2025-01-29T00:00:00Z",
        "project_id": 1,
        "project_name": "生产环境"
      },
      {
        "id": 2,
        "type": "SYSTEM_UPDATE",
        "title": "系统更新可用",
        "description": "发现新的系统更新版本v1.2.0，建议升级",
        "priority": "MEDIUM",
        "due_date": null,
        "project_id": null,
        "project_name": null
      }
    ]
  }
}
```

**获取项目统计趋势**

```text
GET /api/v1/dashboard/project-trends
```

查询参数：

| 参数       | 类型   | 必填 | 默认值 | 描述                               |
| ---------- | ------ | ---- | ------ | ---------------------------------- |
| start_time | string | 否   | -      | 开始时间                           |
| end_time   | string | 否   | -      | 结束时间                           |
| group_by   | string | 否   | day    | 分组方式（hour, day, week, month） |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "project_creation_trend": [
      {
        "date": "2025-01-22",
        "count": 2
      },
      {
        "date": "2025-01-21",
        "count": 1
      }
    ],
    "resource_usage_trend": [
      {
        "date": "2025-01-22",
        "cpu_avg": 45.2,
        "memory_avg": 68.5,
        "storage_avg": 50.0
      }
    ],
    "application_deployment_trend": [
      {
        "date": "2025-01-22",
        "deployments": 5,
        "success_rate": 100.0
      }
    ]
  }
}
```

#### 4.1.2 监控总览看板

**获取全局监控指标**

```text
GET /api/v1/monitoring/global-metrics
```

查询参数：

| 参数              | 类型    | 必填 | 默认值 | 描述                         |
| ----------------- | ------- | ---- | ------ | ---------------------------- |
| time_range        | string  | 否   | 24h    | 时间范围（1h, 24h, 7d, 30d） |
| resource_group_id | integer | 否   | -      | 资源组ID                     |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "summary": {
      "app_instances": {
        "total": 25,
        "running": 20,
        "stopped": 3,
        "error": 2
      },
      "servers": {
        "total": 10,
        "online": 9,
        "offline": 1
      },
      "alerts": {
        "total": 3,
        "critical": 1,
        "warning": 2,
        "info": 0
      },
      "workflows": {
        "total": 15,
        "active": 12,
        "paused": 3
      }
    },
    "trends": {
      "app_deployment": [
        {
          "date": "2025-01-22",
          "count": 5
        }
      ],
      "resource_usage": {
        "cpu_avg": 45.2,
        "memory_avg": 68.5,
        "disk_avg": 32.1
      }
    }
  }
}
```

**获取服务器监控数据**

```text
GET /api/v1/monitoring/servers/{id}/metrics
```

查询参数：

| 参数       | 类型   | 必填 | 默认值 | 描述                                   |
| ---------- | ------ | ---- | ------ | -------------------------------------- |
| metric     | string | 否   | -      | 指标类型（cpu, memory, disk, network） |
| start_time | string | 否   | -      | 开始时间（ISO 8601格式）               |
| end_time   | string | 否   | -      | 结束时间（ISO 8601格式）               |
| interval   | string | 否   | 5m     | 数据间隔（1m, 5m, 1h, 1d）             |

**获取应用监控数据**

```text
GET /api/v1/monitoring/app-instances/{id}/metrics
```

查询参数：

| 参数       | 类型   | 必填 | 默认值 | 描述                                   |
| ---------- | ------ | ---- | ------ | -------------------------------------- |
| metric     | string | 否   | -      | 指标类型（cpu, memory, disk, network） |
| start_time | string | 否   | -      | 开始时间（ISO 8601格式）               |
| end_time   | string | 否   | -      | 结束时间（ISO 8601格式）               |
| interval   | string | 否   | 5m     | 数据间隔（1m, 5m, 1h, 1d）             |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "metrics": [
      {
        "timestamp": "2025-07-15T10:30:00Z",
        "cpu_usage": 45.2,
        "memory_usage": 68.5,
        "disk_usage": 32.1,
        "network_in": 1024000,
        "network_out": 512000
      }
    ],
    "container_info": {
      "container_id": "abc123",
      "image": "wordpress:6.4.2",
      "status": "RUNNING",
      "uptime": 86400
    }
  }
}
```

#### 4.1.3 应用快捷导航

**获取应用快捷导航列表**

```text
GET /api/v1/app-shortcuts
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述       |
| --------- | ------- | ---- | ------ | ---------- |
| page      | integer | 否   | 1      | 页码       |
| page_size | integer | 否   | 20     | 每页数量   |
| keyword   | string  | 否   | -      | 搜索关键词 |
| category  | string  | 否   | -      | 应用分类   |
| status    | string  | 否   | -      | 发布状态   |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "app_instance_id": 123,
        "name": "我的博客",
        "description": "个人技术博客网站",
        "icon": "https://cdn.websoft9.com/icons/wordpress.png",
        "url": "https://blog.example.com",
        "internal_url": "http://192.168.1.100:8080",
        "status": "PUBLISHED",
        "template": {
          "id": 1,
          "name": "WordPress",
          "category": "内容管理"
        },
        "server": {
          "id": 1,
          "name": "Web服务器01",
          "ip_address": "192.168.1.100"
        },
        "health_status": "HEALTHY",
        "last_accessed": "2025-07-15T10:30:00Z",
        "access_count": 150,
        "created_at": "2025-07-15T10:30:00Z"
      }
    ],
    "categories": [
      {
        "name": "内容管理",
        "count": 5
      },
      {
        "name": "开发工具",
        "count": 3
      }
    ]
  }
}
```

**添加应用到快捷导航**

```text
POST /api/v1/app-shortcuts
```

请求体参数：

| 参数名          | 数据类型 | 是否可空 | 描述       |
| --------------- | -------- | -------- | ---------- |
| app_instance_id | integer  | 否       | 应用实例ID |
| name            | string   | 是       | 自定义名称 |
| description     | string   | 是       | 自定义描述 |
| icon            | string   | 是       | 自定义图标 |
| sort_order      | integer  | 是       | 排序顺序   |

**从快捷导航移除应用**

```text
DELETE /api/v1/app-shortcuts/{id}
```

### 4.2 项目管理

#### 4.2.1 监控看板

##### 应用监控看板

**获取应用监控看板数据**

```text
GET /api/v1/monitoring/app-dashboard
```

查询参数：

| 参数              | 类型    | 必填 | 默认值 | 描述     |
| ----------------- | ------- | ---- | ------ | -------- |
| resource_group_id | integer | 否   | -      | 资源组ID |
| status            | string  | 否   | -      | 应用状态 |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "summary": {
      "total_apps": 25,
      "running_apps": 20,
      "stopped_apps": 3,
      "error_apps": 2,
      "avg_cpu_usage": 45.2,
      "avg_memory_usage": 68.5,
      "total_storage_used": 102400
    },
    "app_list": [
      {
        "id": 1,
        "name": "我的博客",
        "template_name": "WordPress",
        "status": "RUNNING",
        "cpu_usage": 25.5,
        "memory_usage": 512,
        "disk_usage": 2048,
        "uptime": 86400,
        "health_status": "HEALTHY",
        "server": {
          "id": 1,
          "name": "Web服务器01",
          "ip_address": "192.168.1.100"
        }
      }
    ],
    "resource_trends": {
      "cpu_trend": [
        {
          "timestamp": "2025-01-22T10:00:00Z",
          "value": 42.1
        }
      ],
      "memory_trend": [
        {
          "timestamp": "2025-01-22T10:00:00Z",
          "value": 65.3
        }
      ]
    }
  }
}
```

**获取应用健康检查状态**

```text
GET /api/v1/monitoring/app-instances/{id}/health
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "HEALTHY",
    "checks": [
      {
        "name": "HTTP响应检查",
        "status": "PASS",
        "response_time": 150,
        "last_check": "2025-07-15T10:30:00Z"
      },
      {
        "name": "数据库连接检查",
        "status": "PASS",
        "response_time": 50,
        "last_check": "2025-07-15T10:30:00Z"
      }
    ],
    "uptime_percentage": 99.9,
    "last_downtime": "2025-01-20T15:30:00Z"
  }
}
```

##### 服务器监控看板

**获取服务器监控看板数据**

```text
GET /api/v1/monitoring/server-dashboard
```

查询参数：

| 参数              | 类型    | 必填 | 默认值 | 描述       |
| ----------------- | ------- | ---- | ------ | ---------- |
| resource_group_id | integer | 否   | -      | 资源组ID   |
| status            | string  | 否   | -      | 服务器状态 |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "summary": {
      "total_servers": 10,
      "online_servers": 9,
      "offline_servers": 1,
      "avg_cpu_usage": 35.8,
      "avg_memory_usage": 62.3,
      "avg_disk_usage": 45.1,
      "total_apps": 25
    },
    "server_list": [
      {
        "id": 1,
        "name": "Web服务器01",
        "ip_address": "192.168.1.100",
        "status": "ONLINE",
        "cpu_usage": 45.2,
        "memory_usage": 68.5,
        "disk_usage": 32.1,
        "load_average": 1.25,
        "uptime": 2592000,
        "app_count": 3,
        "last_heartbeat": "2025-07-15T10:30:00Z"
      }
    ],
    "resource_trends": {
      "cpu_trend": [
        {
          "timestamp": "2025-01-22T10:00:00Z",
          "value": 35.8
        }
      ],
      "memory_trend": [
        {
          "timestamp": "2025-01-22T10:00:00Z",
          "value": 62.3
        }
      ],
      "disk_trend": [
        {
          "timestamp": "2025-01-22T10:00:00Z",
          "value": 45.1
        }
      ]
    }
  }
}
```

**获取服务器系统信息**

```text
GET /api/v1/monitoring/servers/{id}/system-info
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "basic_info": {
      "hostname": "web01.example.com",
      "os_type": "Ubuntu",
      "os_version": "20.04.3 LTS",
      "kernel_version": "5.4.0-91-generic",
      "architecture": "x86_64",
      "uptime": 2592000
    },
    "hardware_info": {
      "cpu_model": "Intel(R) Xeon(R) CPU E5-2680 v4",
      "cpu_cores": 4,
      "cpu_threads": 8,
      "memory_total": 8192,
      "disk_total": 102400,
      "network_interfaces": [
        {
          "name": "eth0",
          "ip": "192.168.1.100",
          "mac": "00:50:56:12:34:56"
        }
      ]
    },
    "process_info": {
      "total_processes": 156,
      "running_processes": 2,
      "sleeping_processes": 154,
      "zombie_processes": 0
    }
  }
}
```

#### 4.2.2 任务看板

**获取任务看板数据**

```text
GET /api/v1/task-dashboard
```

查询参数：

| 参数              | 类型    | 必填 | 默认值 | 描述     |
| ----------------- | ------- | ---- | ------ | -------- |
| resource_group_id | integer | 否   | -      | 资源组ID |
| time_range        | string  | 否   | 24h    | 时间范围 |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_statistics": {
      "today_executions": 25,
      "success_executions": 22,
      "failed_executions": 2,
      "running_executions": 1,
      "success_rate": 88.0,
      "avg_execution_time": 145
    },
    "task_trends": [
      {
        "date": "2025-01-22",
        "total_executions": 25,
        "success_executions": 22,
        "failed_executions": 2,
        "running_executions": 1,
        "success_rate": 88.0,
        "avg_execution_time": 145
      },
      {
        "date": "2025-01-21",
        "total_executions": 18,
        "success_executions": 16,
        "failed_executions": 2,
        "running_executions": 0,
        "success_rate": 88.9,
        "avg_execution_time": 132
      }
    ],
    "active_tasks": [
      {
        "id": 1,
        "name": "应用自动部署",
        "workflow_id": 1,
        "workflow_name": "CI/CD流水线",
        "status": "RUNNING",
        "progress": 65,
        "start_time": "2025-01-22T10:30:00Z",
        "estimated_completion": "2025-01-22T10:35:00Z",
        "executor": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员"
        }
      }
    ],
    "failed_tasks": [
      {

       "id": 2,
        "name": "数据库备份",
        "workflow_id": 2,
        "workflow_name": "定时备份",
        "status": "FAILED",
        "error_message": "数据库连接超时",
        "start_time": "2025-01-22T02:00:00Z",
        "end_time": "2025-01-22T02:05:00Z",
        "duration": 300,
        "retry_count": 3,
        "max_retries": 3,
        "can_retry": false,
        "executor": {
          "id": 1,
          "username": "system",
          "nickname": "系统"
        }
      }
    ],
    "scheduled_tasks": [
      {
        "id": 3,
        "name": "系统健康检查",
        "workflow_id": 3,
        "workflow_name": "健康检查",
        "schedule_type": "CRON",
        "cron_expression": "0 */5 * * * *",
        "next_execution": "2025-01-22T10:35:00Z",
        "last_execution": "2025-01-22T10:30:00Z",
        "last_status": "SUCCESS",
        "is_enabled": true
      },
      {
        "id": 4,
        "name": "日志清理",
        "workflow_id": 4,
        "workflow_name": "日志维护",
        "schedule_type": "CRON",
        "cron_expression": "0 0 2 * * *",
        "next_execution": "2025-01-23T02:00:00Z",
        "last_execution": "2025-01-22T02:00:00Z",
        "last_status": "FAILED",
        "is_enabled": true
      }
    ],
    "task_performance": {
      "avg_execution_time_by_type": [
        {
          "task_type": "DEPLOYMENT",
          "avg_time": 180,
          "count": 10
        },
        {
          "task_type": "BACKUP",
          "avg_time": 300,
          "count": 5
        },
        {
          "task_type": "MONITORING",
          "avg_time": 30,
          "count": 20
        }
      ],
      "execution_time_trend": [
        {
          "date": "2025-01-22",
          "avg_time": 145
        },
        {
          "date": "2025-01-21",
          "avg_time": 132
        }
      ]
    }
  }
}
```

**重试失败任务**

```text
POST /api/v1/workflow-executions/{id}/retry
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述     |
| ------ | -------- | -------- | -------- |
| reason | string   | 是       | 重试原因 |

响应示例：

```json
{
  "code": 200,
  "message": "任务重试已启动",
  "data": {
    "execution_id": "exec_123456790",
    "original_execution_id": "exec_123456789",
    "status": "RUNNING",
    "retry_count": 1,
    "start_time": "2025-01-22T10:40:00Z"
  }
}
```

**批量操作任务**

```text
POST /api/v1/tasks/batch-actions
```

请求体参数：

| 参数名   | 数据类型  | 是否可空 | 描述                                     |
| -------- | --------- | -------- | ---------------------------------------- |
| task_ids | integer[] | 否       | 任务ID数组                               |
| action   | string    | 否       | 操作类型（start, stop, enable, disable） |

响应示例：

```json
{
  "code": 200,
  "message": "批量操作已完成",
  "data": {
    "success_count": 3,
    "failed_count": 1,
    "results": [
      {
        "task_id": 1,
        "status": "SUCCESS",
        "message": "任务启动成功"
      },
      {
        "task_id": 2,
        "status": "FAILED",
        "message": "任务正在运行中，无法重复启动"
      }
    ]
  }
}
```

#### 4.2.3 资源看板

**获取资源看板数据**

```text
GET /api/v1/resource-dashboard
```

查询参数：

| 参数              | 类型    | 必填 | 默认值 | 描述     |
| ----------------- | ------- | ---- | ------ | -------- |
| resource_group_id | integer | 否   | -      | 资源组ID |
| time_range        | string  | 否   | 24h    | 时间范围 |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "cloud_resources": {
      "storage_usage": {
        "used": 151,
        "total": 200,
        "unit": "GB",
        "usage_percentage": 75.5
      },
      "bandwidth_usage": {
        "used": 850,
        "total": 1000,
        "unit": "GB",
        "usage_percentage": 85.0
      },
      "monthly_cost": {
        "current_month": 2621.00,
        "last_month": 2450.00,
        "currency": "CNY",
        "growth_rate": 6.98
      }
    },
    "application_overview": {
      "total_apps": 25,
      "running_apps": 20,
      "stopped_apps": 3,
      "error_apps": 2,
      "published_apps": 15,
      "recent_deployments": [
        {
          "id": 123,
          "name": "我的博客",
          "template_name": "WordPress",
          "status": "RUNNING",
          "deployed_at": "2025-07-15T10:30:00Z"
        }
      ]
    },
    "server_overview": {
      "total_servers": 10,
      "online_servers": 9,
      "offline_servers": 1,
      "avg_cpu_usage": 45.2,
      "avg_memory_usage": 68.5,
      "avg_disk_usage": 32.1,
      "recent_servers": [
        {
          "id": 1,
          "name": "Web服务器01",
          "ip_address": "192.168.1.100",
          "status": "ONLINE",
          "last_heartbeat": "2025-07-15T10:30:00Z"
        }
      ]
    },
    "resource_trends": {
      "cpu_trend": [
        {
          "timestamp": "2025-01-22T10:00:00Z",
          "value": 42.1
        }
      ],
      "memory_trend": [
        {
          "timestamp": "2025-01-22T10:00:00Z",
          "value": 65.3
        }
      ],
      "storage_trend": [
        {
          "timestamp": "2025-01-22T10:00:00Z",
          "value": 75.5
        }
      ]
    },
    "alerts_summary": {
      "total_alerts": 5,
      "critical_alerts": 1,
      "warning_alerts": 3,
      "info_alerts": 1,
      "recent_alerts": [
        {
          "id": 1,
          "title": "服务器CPU使用率过高",
          "severity": "WARNING",
          "fired_at": "2025-07-15T10:30:00Z"
        }
      ]
    }
  }
}
```

**获取资源使用统计**

```text
GET /api/v1/resource-dashboard/statistics
```

查询参数：

| 参数       | 类型   | 必填 | 默认值 | 描述                                |
| ---------- | ------ | ---- | ------ | ----------------------------------- |
| start_time | string | 否   | -      | 开始时间                            |
| end_time   | string | 否   | -      | 结束时间                            |
| group_by   | string | 否   | day    | 分组方式（hour, day, week, month）  |
| metric     | string | 否   | -      | 指标类型（cpu, memory, disk, cost） |

**获取成本分析**

```text
GET /api/v1/resource-dashboard/cost-analysis
```

查询参数：

| 参数       | 类型   | 必填 | 默认值 | 描述     |
| ---------- | ------ | ---- | ------ | -------- |
| start_time | string | 否   | -      | 开始时间 |
| end_time   | string | 否   | -      | 结束时间 |
| group_by   | string | 否   | day    | 分组方式 |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_cost": 2621.00,
    "cost_breakdown": [
      {
        "category": "计算资源",
        "cost": 1500.00,
        "percentage": 57.3
      },
      {
        "category": "存储资源",
        "cost": 800.00,
        "percentage": 30.5
      },
      {
        "category": "网络资源",
        "cost": 321.00,
        "percentage": 12.2
      }
    ],
    "cost_timeline": [
      {
        "date": "2025-01-22",
        "cost": 85.50
      }
    ],
    "cost_forecast": {
      "next_month_estimate": 2800.00,
      "confidence": 85.5
    }
  }
}
```

**获取资源配额信息**

```text
GET /api/v1/resource-dashboard/quotas
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "quotas": [
      {
        "resource_type": "CPU",
        "used": 45.2,
        "total": 100.0,
        "unit": "cores",
        "usage_percentage": 45.2
      },
      {
        "resource_type": "MEMORY",
        "used": 68.5,
        "total": 100.0,
        "unit": "GB",
        "usage_percentage": 68.5
      },
      {
        "resource_type": "STORAGE",
        "used": 1500,
        "total": 2000,
        "unit": "GB",
        "usage_percentage": 75.0
      }
    ],
    "warnings": [
      {
        "resource_type": "MEMORY",
        "message": "内存使用率接近70%，建议关注"
      }
    ]
  }
}
```

#### 4.2.4 文件夹

**获取文件列表**

```text
GET /api/v1/my-space/files
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述         |
| --------- | ------- | ---- | ------ | ------------ |
| path      | string  | 否   | /      | 目录路径     |
| page      | integer | 否   | 1      | 页码         |
| page_size | integer | 否   | 50     | 每页数量     |
| file_type | string  | 否   | -      | 文件类型筛选 |
| keyword   | string  | 否   | -      | 搜索关键词   |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "current_path": "/documents",
    "items": [
      {
        "id": 1,
        "name": "项目文档.pdf",
        "path": "/documents/项目文档.pdf",
        "type": "FILE",
        "size": 2048576,
        "mime_type": "application/pdf",
        "download_count": 5,
        "parent_id": null,
        "storage_path": "/storage/user_1/documents/项目文档.pdf",
        "checksum": "sha256:abc123...",
        "created_at": "2025-07-15T10:30:00Z",
        "updated_at": "2025-07-15T10:30:00Z"
      },
      {
        "id": 2,
        "name": "images",
        "path": "/documents/images",
        "type": "DIRECTORY",
        "size": 0,
        "parent_id": null,
        "storage_path": "/storage/user_1/documents/images",
        "created_at": "2025-07-15T10:30:00Z",
        "updated_at": "2025-07-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 50,
      "total": 2,
      "total_pages": 1
    },
    "storage_info": {
      "used_space": 104857600,
      "total_space": 1073741824,
      "usage_percentage": 9.77
    }
  }
}
```

**上传文件**

```text
POST /api/v1/my-space/files
```

请求体（multipart/form-data）：

```text
file: [文件内容]
path: /documents/
is_public: false
description: 项目相关文档
```

**下载文件**

```text
GET /api/v1/my-space/files/download
```

查询参数：

| 参数 | 类型   | 必填 | 默认值 | 描述     |
| ---- | ------ | ---- | ------ | -------- |
| path | string | 是   | -      | 文件路径 |

**删除文件**

```text
DELETE /api/v1/my-space/files
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述                      |
| ------ | -------- | -------- | ------------------------- |
| paths  | string[] | 否       | 要删除的文件/目录路径数组 |

**创建目录**

```text
POST /api/v1/my-space/directories
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述     |
| ----------- | -------- | -------- | -------- |
| path        | string   | 否       | 目录路径 |
| description | string   | 是       | 目录描述 |

**移动文件**

```text
PUT /api/v1/my-space/files/move
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述         |
| ----------- | -------- | -------- | ------------ |
| source_path | string   | 否       | 源文件路径   |
| target_path | string   | 否       | 目标文件路径 |

**复制文件**

```text
POST /api/v1/my-space/files/copy
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述         |
| ----------- | -------- | -------- | ------------ |
| source_path | string   | 否       | 源文件路径   |
| target_path | string   | 否       | 目标文件路径 |

##### 通用文件上传下载接口

**通用文件上传**

```text
POST /api/v1/upload
```

请求体（multipart/form-data）：

```text
file: [文件内容]
type: avatar|document|image|other
max_size: 10485760
allowed_types: jpg,png,gif,pdf,doc,docx
```

响应：

```json
{
  "code": 200,
  "message": "文件上传成功",
  "data": {
    "file_id": "file_123456789",
    "filename": "document.pdf",
    "original_name": "项目文档.pdf",
    "size": 2048576,
    "mime_type": "application/pdf",
    "url": "https://cdn.websoft9.com/files/file_123456789.pdf",
    "thumbnail_url": "https://cdn.websoft9.com/thumbnails/file_123456789.jpg"
  }
}
```

**获取文件信息**

```text
GET /api/v1/files/{file_id}
```

**删除文件**

```text
DELETE /api/v1/files/{file_id}
```

#### 4.2.5 应用

**获取应用实例列表**

```text
GET /api/v1/app-instances
```

查询参数：

| 参数              | 类型    | 必填 | 默认值 | 描述       |
| ----------------- | ------- | ---- | ------ | ---------- |
| page              | integer | 否   | 1      | 页码       |
| page_size         | integer | 否   | 20     | 每页数量   |
| status            | string  | 否   | -      | 应用状态   |
| resource_group_id | integer | 否   | -      | 资源组ID   |
| owner_id          | integer | 否   | -      | 所有者ID   |
| keyword           | string  | 否   | -      | 搜索关键词 |

**应用操作**

```text
POST /api/v1/app-instances/{id}/actions
```

请求体参数：

| 参数名     | 数据类型 | 是否可空 | 描述                                                      |
| ---------- | -------- | -------- | --------------------------------------------------------- |
| action     | string   | 否       | 操作类型（start, stop, restart, update, backup, restore） |
| parameters | object   | 是       | 操作参数，包含force（是否强制）、timeout（超时时间）等    |

响应示例：

```json
{
  "code": 200,
  "message": "操作执行成功",
  "data": {
    "action": "start",
    "status": "SUCCESS",
    "message": "应用启动成功"
  }
}
```

**获取应用日志**

```text
GET /api/v1/app-instances/{id}/logs
```

查询参数：

| 参数  | 类型    | 必填 | 默认值 | 描述     |
| ----- | ------- | ---- | ------ | -------- |
| lines | integer | 否   | 100    | 日志行数 |
| since | string  | 否   | -      | 开始时间 |

**应用发布**

```text
POST /api/v1/app-instances/{id}/publish
```

请求体参数：

| 参数名               | 数据类型 | 是否可空 | 描述                     |
| -------------------- | -------- | -------- | ------------------------ |
| app_gateway_id       | integer  | 否       | 应用网关ID               |
| service_domain       | string   | 否       | 服务域名                 |
| service_port         | integer  | 是       | 服务端口，默认8080       |
| ssl_certificate_id   | integer  | 是       | SSL证书ID                |
| alert_rule_id        | integer  | 是       | 监控告警规则ID           |
| limit_rules          | object   | 是       | 访问控制策略（JSON格式） |
| health_check_enabled | boolean  | 是       | 是否启用健康检查         |
| audit_log_enabled    | boolean  | 是       | 是否启用审计日志         |

响应示例：

```json
{
  "code": 200,
  "message": "应用发布成功",
  "data": {
    "id": 1,
    "app_instance_id": 123,
    "app_gateway_id": 1,
    "service_domain": "blog.example.com",
    "service_port": 8080,
    "ssl_certificate_id": 1,
    "status": "PUBLISHED",
    "published_url": "https://blog.example.com",
    "created_at": "2025-07-15T10:30:00Z"
  }
}
```

**取消应用发布**

```text
DELETE /api/v1/app-instances/{id}/publish
```

响应示例：

```json
{
  "code": 200,
  "message": "应用已下线",
  "data": {
    "app_instance_id": 123,
    "status": "UNPUBLISHED",
    "unpublished_at": "2025-07-15T10:30:00Z"
  }
}
```

**应用克隆**

```text
POST /api/v1/app-instances/{id}/clone
```

请求体参数：

| 参数名            | 数据类型 | 是否可空 | 描述         |
| ----------------- | -------- | -------- | ------------ |
| name              | string   | 否       | 新实例名称   |
| server_id         | integer  | 是       | 目标服务器ID |
| resource_group_id | integer  | 是       | 资源组ID     |
| owner_id          | integer  | 否       | 所有者ID     |

**应用配置更新**

```text
PUT /api/v1/app-instances/{id}
```

请求体参数：

| 参数名            | 数据类型 | 是否可空 | 描述     |
| ----------------- | -------- | -------- | -------- |
| name              | string   | 是       | 实例名称 |
| resource_group_id | integer  | 是       | 资源组ID |

**获取应用详情**

```text
GET /api/v1/app-instances/{id}
```

**应用下线**

```text
DELETE /api/v1/app-instances/{id}/publish
```

**删除应用实例**

```text
DELETE /api/v1/app-instances/{id}
```

请求体参数：

| 参数名        | 数据类型 | 是否可空 | 描述           |
| ------------- | -------- | -------- | -------------- |
| force         | boolean  | 是       | 是否强制删除   |
| backup_before | boolean  | 是       | 删除前是否备份 |

**应用终端访问**

```text
POST /api/v1/app-instances/{id}/terminal
```

响应：

```json
{
  "code": 200,
  "message": "终端连接已建立",
  "data": {
    "terminal_id": "term_123456789",
    "websocket_url": "wss://api.websoft9.com/v1/terminals/term_123456789",
    "expires_at": "2025-07-15T11:30:00Z"
  }
}
```

#### 4.2.6 工作流

**获取工作流列表**

```text
GET /api/v1/workflows
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述       |
| --------- | ------- | ---- | ------ | ---------- |
| page      | integer | 否   | 1      | 页码       |
| page_size | integer | 否   | 20     | 每页数量   |
| keyword   | string  | 否   | -      | 搜索关键词 |
| status    | string  | 否   | -      | 工作流状态 |
| owner_id  | integer | 否   | -      | 所有者ID   |

**获取工作流详情**

```text
GET /api/v1/workflows/{id}
```

**创建工作流**

```text
POST /api/v1/workflows
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述                   |
| ----------- | -------- | -------- | ---------------------- |
| name        | string   | 否       | 工作流名称             |
| code        | string   | 否       | 工作流编码，唯一标识   |
| description | string   | 是       | 工作流描述             |
| definition  | object   | 否       | 工作流定义（JSON格式） |
| owner_id    | integer  | 否       | 所有者ID               |

**更新工作流**

```text
PUT /api/v1/workflows/{id}
```

**删除工作流**

```text
DELETE /api/v1/workflows/{id}
```

**发布工作流**

```text
POST /api/v1/workflows/{id}/publish
```

请求体参数：

| 参数名    | 数据类型 | 是否可空 | 描述         |
| --------- | -------- | -------- | ------------ |
| status    | string   | 否       | 发布状态     |
| changelog | string   | 是       | 版本更新日志 |

##### 工作流执行

**手动执行工作流**

```text
POST /api/v1/workflows/{id}/execute
```

请求体参数：

| 参数名     | 数据类型 | 是否可空 | 描述                   |
| ---------- | -------- | -------- | ---------------------- |
| task_id    | integer  | 是       | 关联任务ID             |
| input      | object   | 是       | 工作流输入参数         |
| async      | boolean  | 是       | 是否异步执行，默认true |
| trigger_by | integer  | 是       | 触发者ID               |

响应：

```json
{
  "code": 200,
  "message": "工作流执行已启动",
  "data": {
    "id": 1,
    "execution_id": "exec_123456789",
    "task_id": 1,
    "status": "RUNNING",
    "trigger_type": "MANUAL",
    "trigger_by": 1,
    "start_time": "2025-07-15T10:30:00Z"
  }
}
```

**获取执行历史**

```text
GET /api/v1/workflow-executions
```

查询参数：

| 参数        | 类型    | 必填 | 默认值 | 描述     |
| ----------- | ------- | ---- | ------ | -------- |
| workflow_id | integer | 否   | -      | 工作流ID |
| status      | string  | 否   | -      | 执行状态 |
| start_time  | string  | 否   | -      | 开始时间 |
| end_time    | string  | 否   | -      | 结束时间 |

**获取执行详情**

```text
GET /api/v1/workflow-executions/{id}
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "task_id": 1,
    "execution_id": "exec_123456789",
    "status": "SUCCESS",
    "trigger_type": "MANUAL",
    "trigger_by": 1,
    "start_time": "2025-07-15T10:30:00Z",
    "end_time": "2025-07-15T10:30:00Z",
    "duration": 151,
    "error_message": null,
    "execution_log": "执行成功完成",
    "workflow": {
      "id": 1,
      "name": "应用自动部署",
      "code": "app_deploy"
    },
    "task": {
      "id": 1,
      "name": "每日部署任务"
    },
    "trigger_user": {
      "id": 1,
      "username": "admin"
    },
    "created_at": "2025-07-15T10:30:00Z",
    "updated_at": "2025-07-15T10:30:00Z"
  }
}
```

**停止执行**

```text
POST /api/v1/workflow-executions/{id}/stop
```

**获取执行日志**

```text
GET /api/v1/workflow-executions/{id}/logs
```

查询参数：

| 参数         | 类型   | 必填 | 默认值 | 描述     |
| ------------ | ------ | ---- | ------ | -------- |
| component_id | string | 否   | -      | 组件ID   |
| level        | string | 否   | -      | 日志级别 |

##### 工作流组件

**获取组件类型列表**

```text
GET /api/v1/workflow-components
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "type": "START",
      "name": "开始",
      "category": "control",
      "icon": "start",
      "description": "工作流开始节点",
      "input_ports": [],
      "output_ports": ["default"],
      "config_schema": null
    },
    {
      "type": "DEPLOY_APP",
      "name": "部署应用",
      "category": "application",
      "icon": "deploy",
      "description": "部署应用实例",
      "input_ports": ["default"],
      "output_ports": ["success", "failure"],
      "config_schema": {
        "type": "object",
        "properties": {
          "template_id": {
            "type": "integer",
            "title": "应用模板ID",
            "required": true
          },
          "server_id": {
            "type": "integer",
            "title": "目标服务器ID",
            "required": true
          },
          "app_name": {
            "type": "string",
            "title": "应用名称",
            "required": true
          }
        }
      }
    },
    {
      "type": "SEND_EMAIL",
      "name": "发送邮件",
      "category": "notification",
      "icon": "email",
      "description": "发送邮件通知",
      "input_ports": ["default"],
      "output_ports": ["success", "failure"],
      "config_schema": {
        "type": "object",
        "properties": {
          "to": {
            "type": "string",
            "title": "收件人",
            "required": true
          },
          "subject": {
            "type": "string",
            "title": "邮件主题",
            "required": true
          },
          "body": {
            "type": "string",
            "title": "邮件内容",
            "required": true,
            "format": "textarea"
          }
        }
      }
    }
  ]
}
```

##### 任务

**获取任务列表**

```text
GET /api/v1/workflow-tasks
```

查询参数：

| 参数          | 类型    | 必填 | 默认值 | 描述                                  |
| ------------- | ------- | ---- | ------ | ------------------------------------- |
| page          | integer | 否   | 1      | 页码                                  |
| page_size     | integer | 否   | 20     | 每页数量                              |
| workflow_id   | integer | 否   | -      | 工作流ID                              |
| status        | string  | 否   | -      | 任务状态（ACTIVE, PAUSED, STOPPED）   |
| schedule_type | string  | 否   | -      | 调度类型（MANUAL, SCHEDULE, TRIGGER） |
| keyword       | string  | 否   | -      | 搜索关键词                            |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "每日数据备份",
        "workflow": {
          "id": 1,
          "name": "数据备份工作流",
          "version": "1.0.0"
        },
        "schedule_type": "SCHEDULE",
        "cron_expression": "0 2 * * *",
        "status": "ACTIVE",
        "next_run_at": "2025-01-23T02:00:00Z",
        "last_run_at": "2025-01-22T02:00:00Z",
        "last_execution_status": "SUCCESS",
        "run_count": 30,
        "success_count": 29,
        "failure_count": 1,
        "success_rate": 96.67,
        "avg_duration": 180,
        "owner": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员"
        },
        "notification_config": {
          "on_success": false,
          "on_failure": true,
          "channels": ["email", "webhook"]
        },
        "created_at": "2025-01-01T10:30:00Z",
        "updated_at": "2025-07-15T10:30:00Z"
      }
    ]
  }
}
```

**创建任务**

```text
POST /api/v1/workflow-tasks
```

请求体参数：

| 参数名          | 数据类型 | 是否可空 | 描述                                   |
| --------------- | -------- | -------- | -------------------------------------- |
| name            | string   | 否       | 任务名称                               |
| workflow_id     | integer  | 否       | 工作流ID                               |
| schedule_type   | string   | 否       | 调度类型（MANUAL, SCHEDULE, TRIGGER）  |
| cron_expression | string   | 是       | Cron表达式（调度类型为SCHEDULE时必填） |
| description     | string   | 是       | 任务描述                               |
| owner_id        | integer  | 否       | 所有者ID                               |

请求示例：

```json
{
  "name": "每日数据备份",
  "workflow_id": 1,
  "schedule_type": "SCHEDULE",
  "cron_expression": "0 2 * * *",
  "description": "每日凌晨2点执行数据备份任务",
  "owner_id": 1
}
```

**获取任务详情**

```text
GET /api/v1/workflow-tasks/{id}
```

**更新任务**

```text
PUT /api/v1/workflow-tasks/{id}
```

**删除任务**

```text
DELETE /api/v1/workflow-tasks/{id}
```

**启动/暂停/停止任务**

```text
POST /api/v1/workflow-tasks/{id}/actions
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述                           |
| ------ | -------- | -------- | ------------------------------ |
| action | string   | 否       | 操作类型（start, pause, stop） |

**手动执行任务**

```text
POST /api/v1/workflow-tasks/{id}/execute
```

请求体参数：

| 参数名     | 数据类型 | 是否可空 | 描述         |
| ---------- | -------- | -------- | ------------ |
| trigger_by | integer  | 是       | 触发者ID     |
| async      | boolean  | 是       | 是否异步执行 |

**获取任务执行历史**

```text
GET /api/v1/workflow-tasks/{id}/executions
```

查询参数：

| 参数       | 类型    | 必填 | 默认值 | 描述                                  |
| ---------- | ------- | ---- | ------ | ------------------------------------- |
| page       | integer | 否   | 1      | 页码                                  |
| page_size  | integer | 否   | 20     | 每页数量                              |
| status     | string  | 否   | -      | 执行状态（SUCCESS, FAILURE, RUNNING） |
| start_time | string  | 否   | -      | 开始时间                              |
| end_time   | string  | 否   | -      | 结束时间                              |

**获取任务统计信息**

```text
GET /api/v1/workflow-tasks/{id}/statistics
```

查询参数：

| 参数       | 类型   | 必填 | 默认值 | 描述                               |
| ---------- | ------ | ---- | ------ | ---------------------------------- |
| start_time | string | 否   | -      | 开始时间                           |
| end_time   | string | 否   | -      | 结束时间                           |
| group_by   | string | 否   | day    | 分组方式（hour, day, week, month） |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "execution_summary": {
      "total_executions": 30,
      "success_executions": 29,
      "failure_executions": 1,
      "success_rate": 96.67,
      "avg_duration": 180,
      "min_duration": 120,
      "max_duration": 300
    },
    "execution_timeline": [
      {
        "date": "2025-01-22",
        "success_count": 1,
        "failure_count": 0,
        "avg_duration": 165
      }
    ],
    "failure_analysis": [
      {
        "error_type": "TIMEOUT",
        "count": 1,
        "percentage": 100.0
      }
    ]
  }
}
```

**批量操作任务**

```text
POST /api/v1/workflow-tasks/batch-actions
```

请求体参数：

| 参数名   | 数据类型  | 是否可空 | 描述                           |
| -------- | --------- | -------- | ------------------------------ |
| task_ids | integer[] | 否       | 任务ID数组                     |
| action   | string    | 否       | 操作类型（start, pause, stop） |

#### 4.2.7 资源组

**获取资源组列表**

```text
GET /api/v1/resource-groups
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述       |
| --------- | ------- | ---- | ------ | ---------- |
| page      | integer | 否   | 1      | 页码       |
| page_size | integer | 否   | 20     | 每页数量   |
| keyword   | string  | 否   | -      | 搜索关键词 |
| owner_id  | integer | 否   | -      | 所有者ID   |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "生产环境",
        "code": "production",
        "description": "生产环境资源组",
        "owner": {
          "id": 1,
          "username": "admin",
          "nickname": "系统管理员"
        },
        "resource_counts": {
          "servers": 5,
          "app_instances": 15,
          "databases": 3,
          "workflows": 8
        },
        "members": [
          {
            "user_id": 1,
            "username": "admin",
            "role": "OWNER"
          },
          {
            "user_id": 2,
            "username": "dev1",
            "role": "MEMBER"
          }
        ],
        "permissions": [
          "resource:read",
          "resource:write",
          "resource:delete"
        ],
        "created_at": "2025-07-15T10:30:00Z",
        "updated_at": "2025-07-15T10:30:00Z"
      }
    ]
  }
}
```

**创建资源组**

```text
POST /api/v1/resource-groups
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述       |
| ----------- | -------- | -------- | ---------- |
| name        | string   | 否       | 资源组名称 |
| code        | string   | 否       | 资源组编码 |
| description | string   | 是       | 资源组描述 |
| sort_order  | integer  | 是       | 排序顺序   |

**获取资源组详情**

```text
GET /api/v1/resource-groups/{id}
```

**更新资源组**

```text
PUT /api/v1/resource-groups/{id}
```

**删除资源组**

```text
DELETE /api/v1/resource-groups/{id}
```

**获取资源组资源列表**

```text
GET /api/v1/resource-groups/{id}/resources
```

查询参数：

| 参数          | 类型   | 必填 | 默认值 | 描述                                       |
| ------------- | ------ | ---- | ------ | ------------------------------------------ |
| resource_type | string | 否   | -      | 资源类型（SERVER, APP_INSTANCE, DATABASE） |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "servers": [
      {
        "id": 1,
        "name": "Web服务器01",
        "ip_address": "192.168.1.100",
        "status": "ONLINE"
      }
    ],
    "app_instances": [
      {
        "id": 123,
        "name": "我的博客",
        "template_name": "WordPress",
        "status": "RUNNING"
      }
    ],
    "databases": [
      {
        "id": 1,
        "name": "主数据库",
        "db_type": "mysql",
        "status": "CONNECTED"
      }
    ]
  }
}
```

**移动资源到资源组**

```text
POST /api/v1/resource-groups/{id}/resources
```

请求体参数：

| 参数名        | 数据类型  | 是否可空 | 描述                                       |
| ------------- | --------- | -------- | ------------------------------------------ |
| resource_type | string    | 否       | 资源类型（SERVER, APP_INSTANCE, DATABASE） |
| resource_ids  | integer[] | 否       | 资源ID数组                                 |

**从资源组移除资源**

```text
DELETE /api/v1/resource-groups/{id}/resources
```

请求体参数：

| 参数名        | 数据类型  | 是否可空 | 描述                                       |
| ------------- | --------- | -------- | ------------------------------------------ |
| resource_type | string    | 否       | 资源类型（SERVER, APP_INSTANCE, DATABASE） |
| resource_ids  | integer[] | 否       | 资源ID数组                                 |

#### 4.2.8 服务器

**获取服务器列表**

```text
GET /api/v1/servers
```

查询参数：

| 参数              | 类型    | 必填 | 默认值 | 描述                         |
| ----------------- | ------- | ---- | ------ | ---------------------------- |
| page              | integer | 否   | 1      | 页码                         |
| page_size         | integer | 否   | 20     | 每页数量                     |
| keyword           | string  | 否   | -      | 搜索关键词（服务器名称、IP） |
| status            | string  | 否   | -      | 服务器状态                   |
| resource_group_id | integer | 否   | -      | 资源组ID                     |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "Web服务器01",
        "hostname": "web01.example.com",
        "ip_address": "192.168.1.100",
        "internal_ip": "10.0.1.100",
        "ssh_port": 22,
        "os_type": "Ubuntu",
        "os_version": "20.04.3 LTS",
        "kernel_version": "5.4.0-91-generic",
        "cpu_cores": 4,
        "memory_total": 8192,
        "disk_total": 102400,
        "architecture": "x86_64",
        "status": "RUNNING",
        "last_heartbeat_at": "2025-07-15T10:30:00Z",
        "resource_group": {
          "id": 1,
          "name": "生产环境",
          "code": "production"
        },
        "owner": {
          "id": 1,
          "username": "admin",
          "nickname": "系统管理员"
        },
        "agent": {
          "id": 1,
          "status": "ONLINE",
          "version": "1.0.0"
        },
        "description": "生产环境Web服务器",
        "created_at": "2025-07-15T10:30:00Z",
        "updated_at": "2025-07-15T10:30:00Z"
      }
    ]
  }
}
```

**注册服务器**

```text
POST /api/v1/servers
```

请求体参数：

| 参数名            | 数据类型 | 是否可空 | 描述            |
| ----------------- | -------- | -------- | --------------- |
| name              | string   | 否       | 服务器名称      |
| hostname          | string   | 否       | 主机名          |
| ip_address        | string   | 否       | 外网IP地址      |
| internal_ip       | string   | 是       | 内网IP地址      |
| ssh_port          | integer  | 是       | SSH端口，默认22 |
| os_type           | string   | 否       | 操作系统类型    |
| resource_group_id | integer  | 是       | 资源组ID        |
| owner_id          | integer  | 否       | 所有者ID        |
| description       | string   | 是       | 服务器描述      |

**服务器操作**

```text
POST /api/v1/servers/{id}/actions
```

请求体参数：

| 参数名     | 数据类型 | 是否可空 | 描述                                                 |
| ---------- | -------- | -------- | ---------------------------------------------------- |
| action     | string   | 否       | 操作类型（reboot, shutdown, start, upgrade_agent）   |
| parameters | object   | 是       | 操作参数，包含force（是否强制）、delay（延迟秒数）等 |

支持的操作：

- `reboot`: 重启服务器
- `shutdown`: 关闭服务器
- `start`: 启动服务器
- `upgrade_agent`: 升级客户端

**获取服务器终端访问**

```text
GET /api/v1/servers/{id}/terminal
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "terminal_url": "wss://api.websoft9.com/terminals/server_1",
    "session_id": "term_123456789",
    "expires_at": "2025-07-15T11:30:00Z"
  }
}
```

**服务器文件管理**

```text
GET /api/v1/servers/{id}/files
```

查询参数：

| 参数 | 类型   | 必填 | 默认值 | 描述     |
| ---- | ------ | ---- | ------ | -------- |
| path | string | 否   | /      | 目录路径 |

**上传文件到服务器**

```text
POST /api/v1/servers/{id}/files
```

请求体（multipart/form-data）：

```text
file: [文件内容]
path: /home/user/
```

**从服务器下载文件**

```text
GET /api/v1/servers/{id}/files/download
```

查询参数：

| 参数 | 类型   | 必填 | 默认值 | 描述     |
| ---- | ------ | ---- | ------ | -------- |
| path | string | 是   | -      | 文件路径 |

**删除服务器文件**

```text
DELETE /api/v1/servers/{id}/files
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述                      |
| ------ | -------- | -------- | ------------------------- |
| paths  | string[] | 否       | 要删除的文件/目录路径数组 |

**获取服务器系统服务列表**

```text
GET /api/v1/servers/{id}/services
```

**管理服务器系统服务**

```text
POST /api/v1/servers/{id}/services/{service_name}/actions
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述                             |
| ------ | -------- | -------- | -------------------------------- |
| action | string   | 否       | 操作类型（start, stop, restart） |

**获取服务器详情**

```text
GET /api/v1/servers/{id}
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Web服务器01",
    "hostname": "web01.example.com",
    "ip_address": "192.168.1.100",
    "internal_ip": "10.0.1.100",
    "ssh_port": 22,
    "os_type": "Ubuntu",
    "os_version": "20.04.3 LTS",
    "kernel_version": "5.4.0-91-generic",
    "cpu_cores": 4,
    "memory_total": 8192,
    "disk_total": 102400,
    "architecture": "x86_64",
    "status": "RUNNING",
    "last_heartbeat_at": "2025-07-15T10:30:00Z",
    "resource_group": {
      "id": 1,
      "name": "生产环境",
      "code": "production"
    },
    "owner": {
      "id": 1,
      "username": "admin",
      "nickname": "系统管理员"
    },
    "agent": {
      "id": 1,
      "status": "ONLINE",
      "version": "1.0.0"
    },
    "stats": {
      "cpu_usage": 45.2,
      "memory_usage": 68.5,
      "disk_usage": 32.1,
      "load_average": 1.25,
      "app_count": 3,
      "running_app_count": 2
    },
    "apps": [
      {
        "id": 123,
        "name": "我的博客",
        "template_name": "WordPress",
        "status": "RUNNING"
      }
    ],
    "created_at": "2025-07-15T10:30:00Z",
    "updated_at": "2025-07-15T10:30:00Z"
  }
}
```

**更新服务器信息**

```text
PUT /api/v1/servers/{id}
```

请求体参数：

| 参数名            | 数据类型 | 是否可空 | 描述       |
| ----------------- | -------- | -------- | ---------- |
| name              | string   | 是       | 服务器名称 |
| resource_group_id | integer  | 是       | 资源组ID   |
| description       | string   | 是       | 服务器描述 |

**删除服务器**

```text
DELETE /api/v1/servers/{id}
```

请求体参数：

| 参数名          | 数据类型 | 是否可空 | 描述           |
| --------------- | -------- | -------- | -------------- |
| force           | boolean  | 是       | 是否强制删除   |
| cleanup_apps    | boolean  | 是       | 是否清理应用   |
| uninstall_agent | boolean  | 是       | 是否卸载客户端 |

**服务器终端访问**

```text
POST /api/v1/servers/{id}/terminal
```

#### 4.2.9 密钥管理

**获取密钥列表**

```text
GET /api/v1/secret-keys
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述     |
| --------- | ------- | ---- | ------ | -------- |
| page      | integer | 否   | 1      | 页码     |
| page_size | integer | 否   | 20     | 每页数量 |
| key_type  | string  | 否   | -      | 密钥类型 |
| usage     | string  | 否   | -      | 使用场景 |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "数据库连接密钥",
        "key_type": "DATABASE",
        "usage": "MYSQL_CONNECTION",
        "description": "MySQL数据库连接密钥",
        "is_encrypted": true,
        "custom_fields": {"rotation_interval": 90},
        "authorized_users": [1, 2],
        "authorized_groups": [1],
        "expires_at": "2025-04-01T00:00:00Z",
        "owner_id": 1,
        "created_at": "2024-10-01T00:00:00Z",
        "updated_at": "2025-07-15T10:30:00Z"
      }
    ]
  }
}
```

**创建密钥**

```text
POST /api/v1/secret-keys
```

请求体参数：

| 参数名            | 数据类型 | 是否可空 | 描述                                 |
| ----------------- | -------- | -------- | ------------------------------------ |
| name              | string   | 否       | 密钥名称                             |
| key_type          | string   | 否       | 密钥类型（API_KEY, DATABASE, SSH等） |
| encrypted_value   | string   | 否       | 加密后的密钥值                       |
| description       | string   | 是       | 密钥描述                             |
| custom_fields     | object   | 是       | 自定义字段                           |
| authorized_users  | array    | 是       | 授权用户ID列表                       |
| authorized_groups | array    | 是       | 授权用户组ID列表                     |
| expires_at        | string   | 是       | 过期时间                             |
| owner_id          | integer  | 否       | 所有者ID                             |

**更新密钥**

```text
PUT /api/v1/secret-keys/{id}
```

**删除密钥**

```text
DELETE /api/v1/secret-keys/{id}
```

**获取密钥值**

```text
GET /api/v1/secret-keys/{id}/value
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "value": "sk-1234567890abcdef",
    "expires_at": "2025-02-21T10:30:00Z"
  }
}
```

#### 4.2.10 数据库

**获取数据库连接列表**

```text
GET /api/v1/database-connections
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述                                   |
| --------- | ------- | ---- | ------ | -------------------------------------- |
| page      | integer | 否   | 1      | 页码                                   |
| page_size | integer | 否   | 20     | 每页数量                               |
| db_type   | string  | 否   | -      | 数据库类型（mysql, postgresql, redis） |
| status    | string  | 否   | -      | 连接状态（CONNECTED, DISCONNECTED）    |
| keyword   | string  | 否   | -      | 搜索关键词                             |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "主数据库",
        "db_type": "mysql",
        "host": "192.168.1.100",
        "port": 3306,
        "database": "websoft9",
        "username": "root",
        "status": "CONNECTED",
        "version": "8.0.32",
        "charset": "utf8mb4",
        "connection_pool": {
          "max_connections": 100,
          "active_connections": 15,
          "idle_connections": 5
        },
        "performance": {
          "queries_per_second": 125.5,
          "slow_queries": 2,
          "uptime": 2592000
        },
        "storage": {
          "total_size": 10240,
          "used_size": 2048,
          "usage_percentage": 20.0
        },
        "last_connected_at": "2025-01-22T02:00:00Z",
        "owner_id": 1,
        "resource_group_id": 1,
        "created_at": "2025-07-15T10:30:00Z",
        "updated_at": "2025-07-15T10:30:00Z"
      }
    ]
  }
}
```

**添加数据库连接**

```text
POST /api/v1/database-connections
```

请求体参数：

| 参数名             | 数据类型 | 是否可空 | 描述                                   |
| ------------------ | -------- | -------- | -------------------------------------- |
| name               | string   | 否       | 连接名称                               |
| db_type            | string   | 否       | 数据库类型（mysql, postgresql, redis） |
| host               | string   | 否       | 主机地址                               |
| port               | integer  | 否       | 端口号                                 |
| database           | string   | 是       | 数据库名称                             |
| username           | string   | 否       | 用户名                                 |
| password           | string   | 否       | 密码                                   |
| ssl_enabled        | boolean  | 是       | 是否启用SSL，默认false                 |
| connection_timeout | integer  | 是       | 连接超时时间（秒），默认30             |
| max_connections    | integer  | 是       | 最大连接数，默认10                     |
| description        | string   | 是       | 描述信息                               |

**测试数据库连接**

```text
POST /api/v1/database-connections/test-connection
```

请求体参数：

| 参数名   | 数据类型 | 是否可空 | 描述     |
| -------- | -------- | -------- | -------- |
| host     | string   | 否       | 主机地址 |
| port     | integer  | 否       | 端口号   |
| database | string   | 是       | 数据库名 |
| username | string   | 否       | 用户名   |
| password | string   | 否       | 密码     |

响应示例：

```json
{
  "code": 200,
  "message": "数据库连接测试成功",
  "data": {
    "is_connected": true,
    "response_time": 50,
    "server_version": "8.0.32",
    "server_info": {
      "charset": "utf8mb4",
      "timezone": "Asia/Shanghai"
    }
  }
}
```

**获取数据库详情**

```text
GET /api/v1/database-connections/{id}
```

**更新数据库连接**

```text
PUT /api/v1/database-connections/{id}
```

**删除数据库连接**

```text
DELETE /api/v1/database-connections/{id}
```

#### 4.2.11 应用网关

**获取网关列表**

```text
GET /api/v1/app-gateways
```

**获取网关详情**

```text
GET /api/v1/app-gateways/{id}
```

**创建网关**

```text
POST /api/v1/app-gateways
```

请求体参数：

| 参数名            | 数据类型 | 是否可空 | 描述     |
| ----------------- | -------- | -------- | -------- |
| name              | string   | 否       | 网关名称 |
| server_id         | integer  | 否       | 服务器ID |
| description       | string   | 是       | 网关描述 |
| resource_group_id | integer  | 是       | 资源组ID |
| owner_id          | integer  | 否       | 所有者ID |

**获取网关发布列表**

```text
GET /api/v1/app-gateways/{id}/publishes
```

**创建网关发布**

```text
POST /api/v1/app-gateways/{id}/publishes
```

请求体参数：

| 参数名               | 数据类型 | 是否可空 | 描述                     |
| -------------------- | -------- | -------- | ------------------------ |
| app_instance_id      | integer  | 否       | 应用实例ID               |
| service_domain       | string   | 否       | 服务域名                 |
| service_port         | integer  | 是       | 服务端口，默认8080       |
| ssl_certificate_id   | integer  | 是       | SSL证书ID                |
| alert_rule_id        | integer  | 是       | 监控告警规则ID           |
| limit_rules          | object   | 是       | 访问控制策略（JSON格式） |
| health_check_enabled | boolean  | 是       | 健康检查开启，默认true   |
| audit_log_enabled    | boolean  | 是       | 审计日志开启，默认true   |
| owner_id             | integer  | 否       | 所有者ID                 |

响应示例：

```json
{
  "code": 200,
  "message": "网关发布创建成功",
  "data": {
    "id": 1,
    "app_instance_id": 123,
    "app_gateway_id": 1,
    "service_domain": "blog.example.com",
    "service_port": 8080,
    "ssl_certificate_id": 1,
    "status": "PUBLISHED",
    "published_url": "https://blog.example.com",
    "health_check_enabled": true,
    "audit_log_enabled": true,
    "owner_id": 1,
    "created_at": "2025-07-15T10:30:00Z"
  }
}
```

#### 4.2.12 证书管理

**获取SSL证书列表**

```text
GET /api/v1/ssl-certificates
```

查询参数：

| 参数            | 类型    | 必填 | 默认值 | 描述         |
| --------------- | ------- | ---- | ------ | ------------ |
| page            | integer | 否   | 1      | 页码         |
| page_size       | integer | 否   | 20     | 每页数量     |
| domain          | string  | 否   | -      | 域名筛选     |
| status          | string  | 否   | -      | 证书状态     |
| expires_in_days | integer | 否   | -      | 即将过期天数 |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "example.com证书",
        "domain": "example.com",
        "certificate_type": "LETS_ENCRYPT",
        "issuer": "Let's Encrypt",
        "subject": "CN=example.com",
        "serial_number": "03A1B2C3D4E5F6",
        "status": "VALID",
        "auto_renew": true,
        "not_before": "2025-01-01T00:00:00Z",
        "not_after": "2025-04-01T00:00:00Z",
        "days_until_expiry": 69,
        "owner_id": 1,
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-07-15T10:30:00Z"
      }
    ]
  }
}
```

**获取SSL证书详情**

```text
GET /api/v1/ssl-certificates/{id}
```

**申请SSL证书**

```text
POST /api/v1/ssl-certificates
```

请求体参数：

| 参数名           | 数据类型 | 是否可空 | 描述                                                                |
| ---------------- | -------- | -------- | ------------------------------------------------------------------- |
| name             | string   | 否       | 证书名称                                                            |
| domain           | string   | 否       | 主域名                                                              |
| certificate_type | string   | 是       | 证书类型（LETS_ENCRYPT, COMMERCIAL, SELF_SIGNED），默认LETS_ENCRYPT |
| auto_renew       | boolean  | 是       | 是否自动续期，默认true                                              |
| owner_id         | integer  | 否       | 所有者ID                                                            |

**更新SSL证书**

```text
PUT /api/v1/ssl-certificates/{id}
```

**删除SSL证书**

```text
DELETE /api/v1/ssl-certificates/{id}
```

**手动续期证书**

```text
POST /api/v1/ssl-certificates/{id}/renew
```

**验证证书**

```text
POST /api/v1/ssl-certificates/{id}/verify
```

响应：

```json
{
  "code": 200,
  "message": "证书验证成功",
  "data": {
    "is_valid": true,
    "expires_at": "2025-04-01T00:00:00Z",
    "days_until_expiry": 69,
    "chain_valid": true,
    "domain_match": true
  }
}
```

#### 4.2.13 云资源

**获取云账号列表**

```text
GET /api/v1/cloud-accounts
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述                                            |
| --------- | ------- | ---- | ------ | ----------------------------------------------- |
| page      | integer | 否   | 1      | 页码                                            |
| page_size | integer | 否   | 20     | 每页数量                                        |
| provider  | string  | 否   | -      | 云服务商（aliyun, tencent, aws, huawei, azure） |
| status    | string  | 否   | -      | 账号状态（ACTIVE, INACTIVE, ERROR）             |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "阿里云主账号",
        "provider": "aliyun",
        "account_id": "123456789012",
        "status": "ACTIVE",
        "regions": ["cn-hangzhou", "cn-beijing", "cn-shanghai"],
        "default_region": "cn-hangzhou",
        "resource_count": 25,
        "monthly_cost": 1500.00,
        "last_sync_at": "2025-01-22T10:30:00Z",
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-22T10:30:00Z"
      }
    ]
  }
}
```

**添加云账号**

```text
POST /api/v1/cloud-accounts
```

请求体参数：

| 参数名         | 数据类型 | 是否可空 | 描述                                            |
| -------------- | -------- | -------- | ----------------------------------------------- |
| name           | string   | 否       | 账号名称                                        |
| provider       | string   | 否       | 云服务商（aliyun, tencent, aws, huawei, azure） |
| credentials    | object   | 否       | 认证信息（加密存储）                            |
| regions        | array    | 是       | 支持的地域列表                                  |
| default_region | string   | 是       | 默认地域                                        |
| description    | string   | 是       | 账号描述                                        |

请求示例：

```json
{
  "name": "阿里云主账号",
  "provider": "aliyun",
  "credentials": {
    "access_key_id": "LTAI4G...",
    "access_key_secret": "xxx..."
  },
  "regions": ["cn-hangzhou", "cn-beijing"],
  "default_region": "cn-hangzhou",
  "description": "公司阿里云主账号"
}
```

**测试云账号连接**

```text
POST /api/v1/cloud-accounts/test-connection
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述     |
| ----------- | -------- | -------- | -------- |
| provider    | string   | 否       | 云服务商 |
| credentials | object   | 否       | 认证信息 |
| region      | string   | 是       | 测试地域 |

响应示例：

```json
{
  "code": 200,
  "message": "云账号连接测试成功",
  "data": {
    "is_connected": true,
    "account_info": {
      "account_id": "123456789012",
      "account_name": "示例公司",
      "available_regions": ["cn-hangzhou", "cn-beijing", "cn-shanghai"]
    },
    "permissions": {
      "ecs": true,
      "oss": true,
      "rds": false
    }
  }
}
```

**获取云资源列表**

```text
GET /api/v1/cloud-resources
```

查询参数：

| 参数             | 类型    | 必填 | 默认值 | 描述                                   |
| ---------------- | ------- | ---- | ------ | -------------------------------------- |
| cloud_account_id | integer | 否   | -      | 云账号ID                               |
| resource_type    | string  | 否   | -      | 资源类型（ecs, oss, rds, slb, vpc）    |
| region           | string  | 否   | -      | 地域                                   |
| status           | string  | 否   | -      | 资源状态（Running, Stopped, Creating） |
| page             | integer | 否   | 1      | 页码                                   |
| page_size        | integer | 否   | 20     | 每页数量                               |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "cloud_account_id": 1,
        "resource_id": "i-bp1234567890abcdef",
        "resource_name": "web-server-01",
        "resource_type": "ecs",
        "region": "cn-hangzhou",
        "zone": "cn-hangzhou-a",
        "status": "Running",
        "specifications": {
          "instance_type": "ecs.t5-lc1m1.small",
          "cpu": 1,
          "memory": 1024,
          "disk_size": 40,
          "bandwidth": 1
        },
        "network_info": {
          "vpc_id": "vpc-bp1234567890abcdef",
          "vswitch_id": "vsw-bp1234567890abcdef",
          "private_ip": "172.16.0.10",
          "public_ip": "47.96.123.456"
        },
        "cost_info": {
          "billing_method": "PostPaid",
          "monthly_cost": 85.50,
          "currency": "CNY"
        },
        "tags": {
          "Environment": "Production",
          "Project": "WebApp"
        },
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-22T10:30:00Z"
      }
    ]
  }
}
```

**云资源操作**

```text
POST /api/v1/cloud-resources/{id}/actions
```

请求体参数：

| 参数名     | 数据类型 | 是否可空 | 描述                                             |
| ---------- | -------- | -------- | ------------------------------------------------ |
| action     | string   | 否       | 操作类型（start, stop, restart, delete, resize） |
| parameters | object   | 是       | 操作参数，根据不同操作类型包含不同参数           |

请求示例：

```json
{
  "action": "resize",
  "parameters": {
    "instance_type": "ecs.t5-lc1m2.small",
    "force": false
  }
}
```

响应示例：

```json
{
  "code": 200,
  "message": "云资源操作已提交",
  "data": {
    "operation_id": "op_123456789",
    "status": "PROCESSING",
    "estimated_time": 300,
    "message": "正在调整实例规格..."
  }
}
```

**获取云资源操作历史**

```text
GET /api/v1/cloud-resources/{id}/operations
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述     |
| --------- | ------- | ---- | ------ | -------- |
| page      | integer | 否   | 1      | 页码     |
| page_size | integer | 否   | 20     | 每页数量 |
| action    | string  | 否   | -      | 操作类型 |

**同步云资源**

```text
POST /api/v1/cloud-accounts/{id}/sync
```

请求体参数：

| 参数名         | 数据类型 | 是否可空 | 描述                                 |
| -------------- | -------- | -------- | ------------------------------------ |
| resource_types | string[] | 是       | 要同步的资源类型，为空则同步所有类型 |
| regions        | string[] | 是       | 要同步的地域，为空则同步所有地域     |
| force          | boolean  | 是       | 是否强制同步，默认false              |

响应示例：

```json
{
  "code": 200,
  "message": "云资源同步已启动",
  "data": {
    "sync_id": "sync_123456789",
    "status": "RUNNING",
    "progress": 0,
    "estimated_time": 600,
    "sync_tasks": [
      {
        "resource_type": "ecs",
        "region": "cn-hangzhou",
        "status": "PENDING"
      },
      {
        "resource_type": "oss",
        "region": "cn-hangzhou",
        "status": "PENDING"
      }
    ]
  }
}
```

**获取云资源成本分析**

```text
GET /api/v1/cloud-resources/cost-analysis
```

查询参数：

| 参数             | 类型    | 必填 | 默认值 | 描述                         |
| ---------------- | ------- | ---- | ------ | ---------------------------- |
| cloud_account_id | integer | 否   | -      | 云账号ID                     |
| start_time       | string  | 否   | -      | 开始时间                     |
| end_time         | string  | 否   | -      | 结束时间                     |
| group_by         | string  | 否   | day    | 分组方式（day, week, month） |
| resource_type    | string  | 否   | -      | 资源类型                     |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_cost": 2500.00,
    "currency": "CNY",
    "cost_by_resource_type": [
      {
        "resource_type": "ecs",
        "cost": 1500.00,
        "percentage": 60.0
      },
      {
        "resource_type": "oss",
        "cost": 500.00,
        "percentage": 20.0
      },
      {
        "resource_type": "rds",
        "cost": 500.00,
        "percentage": 20.0
      }
    ],
    "cost_timeline": [
      {
        "date": "2025-01-22",
        "cost": 85.50
      },
      {
        "date": "2025-01-21",
        "cost": 82.30
      }
    ],
    "cost_forecast": {
      "next_month_estimate": 2650.00,
      "confidence": 85.5,
      "trend": "INCREASING"
    }
  }
}
```

#### 4.2.14 项目团队

**获取项目成员列表**

```text
GET /api/v1/projects/{project_id}/members
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述                                           |
| --------- | ------- | ---- | ------ | ---------------------------------------------- |
| page      | integer | 否   | 1      | 页码                                           |
| page_size | integer | 否   | 20     | 每页数量                                       |
| role      | string  | 否   | -      | 角色筛选（admin, developer, operator, viewer） |
| status    | string  | 否   | -      | 状态筛选（ACTIVE, INACTIVE, PENDING）          |
| keyword   | string  | 否   | -      | 搜索关键词（用户名、昵称）                     |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "user_id": 1,
        "project_id": 1,
        "role": "admin",
        "status": "ACTIVE",
        "permissions": [
          "project:read",
          "project:write",
          "project:delete",
          "member:manage",
          "resource:manage"
        ],
        "user": {
          "id": 1,
          "username": "admin",
          "nickname": "项目管理员",
          "email": "admin@example.com",
          "avatar": "https://example.com/avatars/admin.jpg",
          "last_login_at": "2025-01-22T10:30:00Z"
        },
        "joined_at": "2025-01-01T00:00:00Z",
        "last_active_at": "2025-01-22T10:30:00Z",
        "activity_stats": {
          "login_count": 150,
          "operation_count": 89,
          "last_operation": "部署应用：WordPress"
        }
      },
      {
        "id": 2,
        "user_id": 2,
        "project_id": 1,
        "role": "developer",
        "status": "ACTIVE",
        "permissions": [
          "project:read",
          "app:deploy",
          "app:manage",
          "workflow:create",
          "workflow:execute"
        ],
        "user": {
          "id": 2,
          "username": "developer",
          "nickname": "开发工程师",
          "email": "dev@example.com",
          "avatar": "https://example.com/avatars/dev.jpg",
          "last_login_at": "2025-01-22T09:15:00Z"
        },
        "joined_at": "2025-01-05T00:00:00Z",
        "last_active_at": "2025-01-22T09:15:00Z",
        "activity_stats": {
          "login_count": 85,
          "operation_count": 156,
          "last_operation": "执行工作流：CI/CD部署"
        }
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 2,
      "total_pages": 1
    },
    "statistics": {
      "total_members": 2,
      "active_members": 2,
      "pending_members": 0,
      "role_distribution": {
        "admin": 1,
        "developer": 1,
        "operator": 0,
        "viewer": 0
      }
    }
  }
}
```

**邀请项目成员**

```text
POST /api/v1/projects/{project_id}/members/invite
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述                                           |
| ----------- | -------- | -------- | ---------------------------------------------- |
| email       | string   | 否       | 被邀请人邮箱                                   |
| role        | string   | 否       | 项目角色（admin, developer, operator, viewer） |
| permissions | string[] | 是       | 自定义权限列表                                 |
| message     | string   | 是       | 邀请消息                                       |
| expires_in  | integer  | 是       | 邀请有效期（小时），默认72小时                 |

响应示例：

```json
{
  "code": 200,
  "message": "邀请已发送",
  "data": {
    "invitation_id": "inv_123456789",
    "email": "newuser@example.com",
    "role": "developer",
    "status": "PENDING",
    "invite_url": "https://platform.websoft9.com/invitations/inv_123456789",
    "expires_at": "2025-01-25T10:30:00Z",
    "created_at": "2025-01-22T10:30:00Z"
  }
}
```

**获取邀请列表**

```text
GET /api/v1/projects/{project_id}/invitations
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述                                              |
| --------- | ------- | ---- | ------ | ------------------------------------------------- |
| page      | integer | 否   | 1      | 页码                                              |
| page_size | integer | 否   | 20     | 每页数量                                          |
| status    | string  | 否   | -      | 邀请状态（PENDING, ACCEPTED, EXPIRED, CANCELLED） |

**取消邀请**

```text
DELETE /api/v1/projects/{project_id}/invitations/{invitation_id}
```

**接受邀请**

```text
POST /api/v1/invitations/{invitation_id}/accept
```

请求体参数：

| 参数名   | 数据类型 | 是否可空 | 描述     |
| -------- | -------- | -------- | -------- |
| password | string   | 是       | 用户密码 |

**拒绝邀请**

```text
POST /api/v1/invitations/{invitation_id}/decline
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述     |
| ------ | -------- | -------- | -------- |
| reason | string   | 是       | 拒绝原因 |

**更新成员角色**

```text
PUT /api/v1/projects/{project_id}/members/{member_id}
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述           |
| ----------- | -------- | -------- | -------------- |
| role        | string   | 是       | 项目角色       |
| permissions | string[] | 是       | 自定义权限列表 |
| status      | string   | 是       | 成员状态       |

**移除项目成员**

```text
DELETE /api/v1/projects/{project_id}/members/{member_id}
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述     |
| ------ | -------- | -------- | -------- |
| reason | string   | 是       | 移除原因 |

**获取成员活动记录**

```text
GET /api/v1/projects/{project_id}/members/{member_id}/activities
```

查询参数：

| 参数       | 类型    | 必填 | 默认值 | 描述     |
| ---------- | ------- | ---- | ------ | -------- |
| page       | integer | 否   | 1      | 页码     |
| page_size  | integer | 否   | 20     | 每页数量 |
| start_time | string  | 否   | -      | 开始时间 |
| end_time   | string  | 否   | -      | 结束时间 |
| action     | string  | 否   | -      | 操作类型 |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "action": "APP_DEPLOY",
        "resource_type": "APPLICATION",
        "resource_id": 123,
        "resource_name": "WordPress博客",
        "description": "部署应用：WordPress博客",
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
        "result": "SUCCESS",
        "created_at": "2025-01-22T10:30:00Z"
      }
    ]
  }
}
```

**获取项目角色权限配置**

```text
GET /api/v1/projects/{project_id}/role-permissions
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "roles": [
      {
        "role": "admin",
        "name": "项目管理员",
        "description": "拥有项目的完整管理权限",
        "permissions": [
          "project:read",
          "project:write",
          "project:delete",
          "member:manage",
          "resource:manage",
          "app:deploy",
          "app:manage",
          "workflow:create",
          "workflow:execute",
          "server:manage"
        ]
      },
      {
        "role": "developer",
        "name": "开发者",
        "description": "拥有应用开发和部署权限",
        "permissions": [
          "project:read",
          "app:deploy",
          "app:manage",
          "workflow:create",
          "workflow:execute"
        ]
      },
      {
        "role": "operator",
        "name": "运维者",
        "description": "拥有服务器和监控管理权限",
        "permissions": [
          "project:read",
          "server:manage",
          "monitor:read",
          "alert:manage"
        ]
      },
      {
        "role": "viewer",
        "name": "查看者",
        "description": "只能查看项目信息",
        "permissions": [
          "project:read"
        ]
      }
    ]
  }
}
```

#### 4.2.15 项目设置

**获取项目环境变量**

```text
GET /api/v1/projects/{project_id}/environment-variables
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述                                |
| --------- | ------- | ---- | ------ | ----------------------------------- |
| page      | integer | 否   | 1      | 页码                                |
| page_size | integer | 否   | 50     | 每页数量                            |
| scope     | string  | 否   | -      | 作用域筛选（global, app, workflow） |
| keyword   | string  | 否   | -      | 搜索关键词                          |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "project_id": 1,
        "name": "DATABASE_URL",
        "value": "mysql://user:pass@localhost:3306/db",
        "type": "SENSITIVE",
        "scope": "global",
        "description": "数据库连接地址",
        "is_encrypted": true,
        "usage_count": 5,
        "last_used_at": "2025-01-22T10:30:00Z",
        "created_by": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员"
        },
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-22T10:30:00Z"
      },
      {
        "id": 2,
        "project_id": 1,
        "name": "APP_VERSION",
        "value": "1.2.0",
        "type": "NORMAL",
        "scope": "app",
        "description": "应用版本号",
        "is_encrypted": false,
        "usage_count": 12,
        "last_used_at": "2025-01-22T09:15:00Z",
        "created_by": {
          "id": 2,
          "username": "developer",
          "nickname": "开发者"
        },
        "created_at": "2025-01-15T00:00:00Z",
        "updated_at": "2025-01-20T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 50,
      "total": 2,
      "total_pages": 1
    },
    "statistics": {
      "total_variables": 2,
      "sensitive_variables": 1,
      "global_variables": 1,
      "app_variables": 1,
      "workflow_variables": 0
    }
  }
}
```

**创建环境变量**

```text
POST /api/v1/projects/{project_id}/environment-variables
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述                                        |
| ----------- | -------- | -------- | ------------------------------------------- |
| name        | string   | 否       | 变量名称，大写字母和下划线                  |
| value       | string   | 否       | 变量值                                      |
| type        | string   | 是       | 变量类型（NORMAL, SENSITIVE），默认NORMAL   |
| scope       | string   | 是       | 作用域（global, app, workflow），默认global |
| description | string   | 是       | 变量描述                                    |

验证规则：

- `name`: 必填，1-100字符，大写字母、数字、下划线，不能以数字开头
- `value`: 必填，最大10000字符
- `type`: NORMAL（普通）或 SENSITIVE（敏感）
- `scope`: global（全局）、app（应用）或 workflow（工作流）

**更新环境变量**

```text
PUT /api/v1/projects/{project_id}/environment-variables/{id}
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述     |
| ----------- | -------- | -------- | -------- |
| value       | string   | 是       | 变量值   |
| type        | string   | 是       | 变量类型 |
| scope       | string   | 是       | 作用域   |
| description | string   | 是       | 变量描述 |

**删除环境变量**

```text
DELETE /api/v1/projects/{project_id}/environment-variables/{id}
```

**批量操作环境变量**

```text
POST /api/v1/projects/{project_id}/environment-variables/batch
```

请求体参数：

| 参数名  | 数据类型  | 是否可空 | 描述                       |
| ------- | --------- | -------- | -------------------------- |
| action  | string    | 否       | 操作类型（delete, export） |
| var_ids | integer[] | 否       | 环境变量ID数组             |

**获取项目配置**

```text
GET /api/v1/projects/{project_id}/settings
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "basic_info": {
      "name": "生产环境",
      "identifier": "production",
      "description": "生产环境项目",
      "tags": ["production", "web"],
      "icon": "https://example.com/icons/project.png"
    },
    "advanced_settings": {
      "default_resource_group": "production",
      "default_timezone": "Asia/Shanghai",
      "log_retention_days": 30,
      "backup_strategy": "daily",
      "auto_scaling_enabled": true,
      "monitoring_enabled": true
    },
    "security_settings": {
      "access_control": "members",
      "api_access_enabled": true,
      "audit_enabled": true,
      "two_factor_required": false,
      "ip_whitelist": [],
      "session_timeout": 3600
    },
    "notification_settings": {
      "email_notifications": true,
      "sms_notifications": false,
      "webhook_notifications": true,
      "notification_channels": [
        {
          "type": "email",
          "enabled": true,
          "config": {
            "recipients": ["admin@example.com"]
          }
        }
      ]
    },
    "integration_settings": {
      "git_integration": {
        "enabled": true,
        "repository_url": "https://github.com/company/project.git",
        "branch": "main"
      },
      "ci_cd_integration": {
        "enabled": true,
        "provider": "jenkins",
        "webhook_url": "https://jenkins.example.com/webhook"
      }
    }
  }
}
```

**更新项目配置**

```text
PUT /api/v1/projects/{project_id}/settings
```

请求体参数：

| 参数名                | 数据类型 | 是否可空 | 描述     |
| --------------------- | -------- | -------- | -------- |
| basic_info            | object   | 是       | 基本信息 |
| advanced_settings     | object   | 是       | 高级设置 |
| security_settings     | object   | 是       | 安全设置 |
| notification_settings | object   | 是       | 通知设置 |
| integration_settings  | object   | 是       | 集成设置 |

**获取项目使用统计**

```text
GET /api/v1/projects/{project_id}/usage-statistics
```

查询参数：

| 参数       | 类型   | 必填 | 默认值 | 描述                               |
| ---------- | ------ | ---- | ------ | ---------------------------------- |
| start_time | string | 否   | -      | 开始时间                           |
| end_time   | string | 否   | -      | 结束时间                           |
| group_by   | string | 否   | day    | 分组方式（hour, day, week, month） |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "resource_usage": {
      "servers": 5,
      "applications": 15,
      "workflows": 8,
      "storage_used": 102400000000,
      "bandwidth_used": 51200000000
    },
    "activity_statistics": {
      "total_operations": 1250,
      "successful_operations": 1200,
      "failed_operations": 50,
      "active_users": 8,
      "total_logins": 156
    },
    "cost_analysis": {
      "total_cost": 1500.00,
      "currency": "CNY",
      "cost_breakdown": [
        {
          "category": "计算资源",
          "cost": 800.00,
          "percentage": 53.3
        },
        {
          "category": "存储资源",
          "cost": 400.00,
          "percentage": 26.7
        },
        {
          "category": "网络资源",
          "cost": 300.00,
          "percentage": 20.0
        }
      ]
    },
    "usage_trends": [
      {
        "date": "2025-01-22",
        "operations": 85,
        "active_users": 8,
        "resource_usage": 75.5
      }
    ]
  }
}
```

### 4.3 应用市场

#### 4.3.1 应用市场列表

##### 应用分类目录

**获取应用分类列表**

```text
GET /api/v1/app-store-categories
```

**获取应用分类树**

```text
GET /api/v1/app-store-categories/tree
```

**创建应用分类**

```text
POST /api/v1/app-store-categories
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述                   |
| ----------- | -------- | -------- | ---------------------- |
| name        | string   | 否       | 分类名称               |
| code        | string   | 否       | 分类代码，唯一标识     |
| parent_id   | integer  | 是       | 父分类ID，null表示顶级 |
| icon        | string   | 是       | 分类图标URL            |
| description | string   | 是       | 分类描述               |
| sort_order  | integer  | 是       | 排序顺序，默认0        |
| status      | integer  | 是       | 状态：0-禁用，1-启用   |

##### 应用列表

**获取应用列表**

```text
GET /api/v1/app-store-templates
```

查询参数：

| 参数        | 类型    | 必填 | 默认值         | 描述         |
| ----------- | ------- | ---- | -------------- | ------------ |
| category_id | integer | 否   | -              | 分类ID       |
| keyword     | string  | 否   | -              | 搜索关键词   |
| is_official | boolean | 否   | -              | 是否官方应用 |
| is_featured | boolean | 否   | -              | 是否推荐应用 |
| status      | integer | 否   | -              | 状态筛选     |
| sort        | string  | 否   | download_count | 排序字段     |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "WordPress",
        "code": "wordpress",
        "category_id": 1,
        "version": "6.4.2",
        "icon": "https://cdn.websoft9.com/icons/wordpress.png",
        "description": "WordPress是一个功能强大的内容管理系统",
        "official_url": "https://wordpress.org",
        "source_url": "https://github.com/WordPress/WordPress",
        "compose_template": "version: '3.8'\nservices:\n  wordpress:\n    image: wordpress:6.4.2\n    ...",
        "download_count": 1250,
        "star_count": 89,
        "rating": 4.5,
        "is_official": true,
        "is_featured": true,
        "status": 1,
        "category": {
          "id": 1,
          "name": "内容管理",
          "code": "cms"
        },
        "created_at": "2025-07-15T10:30:00Z",
        "updated_at": "2025-07-15T10:30:00Z"
      }
    ]
  }
}
```

**获取应用详情**

```text
GET /api/v1/app-store-templates/{id}
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "WordPress",
    "code": "wordpress",
    "category_id": 1,
    "version": "6.4.2",
    "icon": "https://cdn.websoft9.com/icons/wordpress.png",
    "description": "WordPress是一个功能强大的内容管理系统",
    "official_url": "https://wordpress.org",
    "source_url": "https://github.com/WordPress/WordPress",
    "compose_template": "version: '3.8'\nservices:\n  wordpress:\n    image: wordpress:6.4.2\n    ...",
    "download_count": 1250,
    "star_count": 89,
    "rating": 4.5,
    "is_official": true,
    "is_featured": true,
    "status": 1,
    "category": {
      "id": 1,
      "name": "内容管理",
      "code": "cms"
    },
    "created_at": "2025-07-15T10:30:00Z",
    "updated_at": "2025-07-15T10:30:00Z"
  }
}
```

**应用评价**

```text
POST /api/v1/app-store-templates/{id}/reviews
```

请求体参数：

| 参数名  | 数据类型 | 是否可空 | 描述          |
| ------- | -------- | -------- | ------------- |
| rating  | integer  | 否       | 评分（1-5分） |
| content | string   | 是       | 评价内容      |
| tags    | array    | 是       | 评价标签      |

**获取应用评价**

```text
GET /api/v1/app-store-templates/{id}/reviews
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述     |
| --------- | ------- | ---- | ------ | -------- |
| page      | integer | 否   | 1      | 页码     |
| page_size | integer | 否   | 20     | 每页数量 |
| rating    | integer | 否   | -      | 评分筛选 |

**应用点赞**

```text
POST /api/v1/app-store-templates/{id}/star
```

**取消应用点赞**

```text
DELETE /api/v1/app-store-templates/{id}/star
```

**收藏应用**

```text
POST /api/v1/app-store-templates/{id}/favorite
```

**取消收藏应用**

```text
DELETE /api/v1/app-store-templates/{id}/favorite
```

**举报应用**

```text
POST /api/v1/app-store-templates/{id}/report
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述                                       |
| ------ | -------- | -------- | ------------------------------------------ |
| reason | string   | 否       | 举报原因（SPAM, INAPPROPRIATE, COPYRIGHT） |
| detail | string   | 是       | 详细说明                                   |

**下载应用模板**

```text
GET /api/v1/app-store-templates/{id}/download
```

响应：

```json
{
  "code": 200,
  "message": "模板下载成功",
  "data": {
    "download_url": "https://cdn.websoft9.com/templates/wordpress-6.4.2.w9f",
    "filename": "wordpress-6.4.2.w9f",
    "size": 20485,
    "expires_at": "2025-07-22T11:30:00Z"
  }
}
```

##### 应用部署

**一键部署应用**

```text
POST /api/v1/app-store-templates/{id}/deploy
```

请求体参数：

| 参数名            | 数据类型 | 是否可空 | 描述                    |
| ----------------- | -------- | -------- | ----------------------- |
| name              | string   | 否       | 应用实例名称            |
| server_id         | integer  | 否       | 目标服务器ID            |
| resource_group_id | integer  | 是       | 资源组ID                |
| config_data       | object   | 是       | 部署配置数据            |
| auto_start        | boolean  | 是       | 是否自动启动，默认true  |
| auto_publish      | boolean  | 是       | 是否自动发布，默认false |

请求示例：

```json
{
  "name": "我的博客",
  "server_id": 1,
  "resource_group_id": 1,
  "config_data": {
    "environment": {
      "WORDPRESS_DB_HOST": "mysql:3306",
      "WORDPRESS_DB_USER": "wordpress",
      "WORDPRESS_DB_PASSWORD": "secure_password",
      "WORDPRESS_DB_NAME": "wordpress"
    },
    "resources": {
      "cpu_limit": "1000m",
      "memory_limit": "1Gi"
    },
    "ports": [
      {
        "container_port": 80,
        "host_port": 8080,
        "protocol": "TCP"
      }
    ],
    "volumes": [
      {
        "name": "wordpress-data",
        "mount_path": "/var/www/html",
        "size": "10Gi"
      }
    ]
  },
  "auto_start": true,
  "auto_publish": false
}
```

响应示例：

```json
{
  "code": 200,
  "message": "应用部署任务已创建",
  "data": {
    "deployment_id": "deploy_123456789",
    "app_instance_id": 123,
    "status": "DEPLOYING",
    "estimated_time": 300,
    "progress": 0,
    "steps": [
      {
        "name": "下载镜像",
        "status": "PENDING"
      },
      {
        "name": "创建容器",
        "status": "PENDING"
      },
      {
        "name": "启动服务",
        "status": "PENDING"
      }
    ]
  }
}
```

**获取部署状态**

```text
GET /api/v1/app-deployments/{deployment_id}
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "deployment_id": "deploy_123456789",
    "template_id": 1,
    "app_instance_id": 123,
    "server_id": 1,
    "status": "DEPLOYING",
    "progress": 65,
    "estimated_time": 300,
    "start_time": "2025-07-15T10:30:00Z",
    "end_time": null,
    "error_message": null,
    "deployment_log": "正在下载镜像...\n创建容器...\n启动服务...",
    "config_data": {
      "environment": {
        "WORDPRESS_DB_HOST": "mysql:3306"
      }
    },
    "owner_id": 1,
    "created_at": "2025-07-15T10:30:00Z",
    "updated_at": "2025-07-15T10:32:00Z"
  }
}
```

**取消部署**

```text
DELETE /api/v1/app-deployments/{deployment_id}
```

响应示例：

```json
{
  "code": 200,
  "message": "部署已取消",
  "data": {
    "deployment_id": "deploy_123456789",
    "status": "CANCELLED",
    "cancelled_at": "2025-07-15T10:35:00Z"
  }
}
```

**获取应用部署模板配置**

```text
GET /api/v1/app-store-templates/{id}/config
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "compose_template": "version: '3.8'\nservices:\n  wordpress:\n    image: wordpress:6.4.2\n    environment:\n      WORDPRESS_DB_HOST: ${WORDPRESS_DB_HOST}\n    ports:\n      - \"80:80\"\n    volumes:\n      - wordpress_data:/var/www/html",
    "environment_variables": [
      {
        "name": "WORDPRESS_DB_HOST",
        "label": "数据库主机",
        "type": "string",
        "required": true,
        "default_value": "mysql:3306",
        "description": "MySQL数据库连接地址"
      }
    ],
    "port_mappings": [
      {
        "container_port": 80,
        "protocol": "HTTP",
        "description": "Web服务端口",
        "required": true
      }
    ],
    "volume_mounts": [
      {
        "path": "/var/www/html",
        "description": "WordPress文件目录",
        "required": true,
        "default_size": "10Gi"
      }
    ],
    "resource_requirements": {
      "min_cpu": "500m",
      "min_memory": "512Mi",
      "recommended_cpu": "1000m",
      "recommended_memory": "1Gi"
    }
  }
}
```

#### 4.3.2 应用心愿单

**获取应用心愿单列表**

```text
GET /api/v1/app-store-wishlists
```

查询参数：

| 参数      | 类型    | 必填 | 默认值     | 描述                                              |
| --------- | ------- | ---- | ---------- | ------------------------------------------------- |
| page      | integer | 否   | 1          | 页码                                              |
| page_size | integer | 否   | 20         | 每页数量                                          |
| keyword   | string  | 否   | -          | 搜索关键词                                        |
| status    | string  | 否   | -          | 状态（PENDING, IN_PROGRESS, COMPLETED, EXPIRED）  |
| priority  | integer | 否   | -          | 优先级（1-高，2-中，3-低）                        |
| sort      | string  | 否   | vote_count | 排序字段（vote_count, created_at, reward_amount） |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "NextCloud",
        "version": "28.0.0",
        "description": "希望能够支持NextCloud私有云盘的一键部署",
        "source_url": "https://nextcloud.com",
        "reward_amount": 500.00,
        "priority": 1,
        "status": "PENDING",
        "view_count": 150,
        "like_count": 25,
        "vote_count": 18,
        "comment_count": 8,
        "submitter": {
          "id": 2,
          "username": "user123",
          "nickname": "云端用户"
        },
        "completed_at": null,
        "expires_at": "2025-02-22T10:30:00Z",
        "created_at": "2025-07-15T10:30:00Z",
        "updated_at": "2025-07-15T10:30:00Z"
      }
    ],
    "statistics": {
      "total_wishlists": 45,
      "pending_count": 30,
      "in_progress_count": 10,
      "completed_count": 5,
      "total_reward": 15000.00
    }
  }
}
```

**提交应用心愿单**

```text
POST /api/v1/app-store-wishlists
```

请求体参数：

| 参数名        | 数据类型 | 是否可空 | 描述                |
| ------------- | -------- | -------- | ------------------- |
| name          | string   | 否       | 应用名称            |
| version       | string   | 是       | 版本号              |
| source_url    | string   | 是       | 来源地址            |
| description   | string   | 否       | 需求描述            |
| reward_amount | number   | 是       | 悬赏金额，默认0     |
| priority      | integer  | 是       | 优先级，默认3（低） |

**获取心愿单详情**

```text
GET /api/v1/app-store-wishlists/{id}
```

**更新心愿单**

```text
PUT /api/v1/app-store-wishlists/{id}
```

**删除心愿单**

```text
DELETE /api/v1/app-store-wishlists/{id}
```

**心愿单投票**

```text
POST /api/v1/app-store-wishlists/{id}/vote
```

**心愿单点赞**

```text
POST /api/v1/app-store-wishlists/{id}/like
```

**获取心愿单评论**

```text
GET /api/v1/app-store-wishlists/{id}/comments
```

**添加心愿单评论**

```text
POST /api/v1/app-store-wishlists/{id}/comments
```

请求体参数：

| 参数名    | 数据类型 | 是否可空 | 描述     |
| --------- | -------- | -------- | -------- |
| content   | string   | 否       | 评论内容 |
| parent_id | integer  | 是       | 父评论ID |

**举报心愿单**

```text
POST /api/v1/app-store-wishlists/{id}/report
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述                                       |
| ------ | -------- | -------- | ------------------------------------------ |
| reason | string   | 否       | 举报原因（SPAM, INAPPROPRIATE, DUPLICATE） |
| detail | string   | 是       | 详细说明                                   |

##### 用户收藏和评价

**获取用户收藏的应用模板**

```text
GET /api/v1/user/favorites/app-templates
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述     |
| --------- | ------- | ---- | ------ | -------- |
| page      | integer | 否   | 1      | 页码     |
| page_size | integer | 否   | 20     | 每页数量 |

**获取用户评价的应用模板**

```text
GET /api/v1/user/reviews/app-templates
```

**获取用户的心愿单**

```text
GET /api/v1/user/wishlists
```

**获取用户点赞的应用模板**

```text
GET /api/v1/user/starred/app-templates
```

### 4.4 平台管理

#### 4.4.1 项目管理

**获取项目列表**

```text
GET /api/v1/admin/projects
```

查询参数：

| 参数       | 类型    | 必填 | 默认值 | 描述                                    |
| ---------- | ------- | ---- | ------ | --------------------------------------- |
| page       | integer | 否   | 1      | 页码                                    |
| page_size  | integer | 否   | 20     | 每页数量                                |
| keyword    | string  | 否   | -      | 搜索关键词（项目名称、标识符）          |
| status     | string  | 否   | -      | 项目状态（ACTIVE, INACTIVE, SUSPENDED） |
| owner_id   | integer | 否   | -      | 项目所有者ID                            |
| start_time | string  | 否   | -      | 创建开始时间                            |
| end_time   | string  | 否   | -      | 创建结束时间                            |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "生产环境",
        "identifier": "production",
        "description": "生产环境项目",
        "status": "ACTIVE",
        "tags": ["production", "web"],
        "icon": "https://example.com/icons/project.png",
        "owner": {
          "id": 1,
          "username": "admin",
          "nickname": "系统管理员"
        },
        "statistics": {
          "member_count": 8,
          "server_count": 5,
          "app_count": 15,
          "workflow_count": 8,
          "storage_used": 102400000000
        },
        "resource_usage": {
          "cpu_usage": 45.2,
          "memory_usage": 68.5,
          "storage_usage": 50.0
        },
        "last_activity": {
          "action": "APP_DEPLOY",
          "user": "developer",
          "timestamp": "2025-01-22T10:30:00Z"
        },
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-22T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 1,
      "total_pages": 1
    },
    "statistics": {
      "total_projects": 15,
      "active_projects": 12,
      "inactive_projects": 2,
      "suspended_projects": 1,
      "total_members": 45,
      "total_resources": 125
    }
  }
}
```

**创建项目**

```text
POST /api/v1/admin/projects
```

请求体参数：

| 参数名                 | 数据类型 | 是否可空 | 描述                                 |
| ---------------------- | -------- | -------- | ------------------------------------ |
| name                   | string   | 否       | 项目名称，3-50字符                   |
| identifier             | string   | 否       | 项目标识符，3-50字符，字母数字下划线 |
| description            | string   | 是       | 项目描述                             |
| tags                   | array    | 是       | 项目标签                             |
| icon                   | string   | 是       | 项目图标URL                          |
| owner_id               | integer  | 否       | 项目所有者ID                         |
| default_resource_group | string   | 是       | 默认资源组，默认"default"            |
| default_timezone       | string   | 是       | 默认时区，默认"Asia/Shanghai"        |
| log_retention_days     | integer  | 是       | 日志保留天数，默认30                 |
| backup_strategy        | string   | 是       | 备份策略，默认"daily"                |
| access_control         | string   | 是       | 访问控制，默认"members"              |
| api_access             | boolean  | 是       | API访问权限，默认true                |
| audit_enabled          | boolean  | 是       | 审计日志开启，默认true               |

验证规则：

- `name`: 必填，3-50字符，不能重复
- `identifier`: 必填，3-50字符，字母数字下划线，不能重复，不能以数字开头
- `owner_id`: 必填，必须是有效的用户ID

**获取项目详情**

```text
GET /api/v1/admin/projects/{id}
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "生产环境",
    "identifier": "production",
    "description": "生产环境项目",
    "status": "ACTIVE",
    "tags": ["production", "web"],
    "icon": "https://example.com/icons/project.png",
    "owner": {
      "id": 1,
      "username": "admin",
      "nickname": "系统管理员",
      "email": "admin@example.com"
    },
    "settings": {
      "default_resource_group": "production",
      "default_timezone": "Asia/Shanghai",
      "log_retention_days": 30,
      "backup_strategy": "daily",
      "access_control": "members",
      "api_access": true,
      "audit_enabled": true
    },
    "statistics": {
      "member_count": 8,
      "server_count": 5,
      "app_count": 15,
      "workflow_count": 8,
      "storage_used": 102400000000,
      "total_operations": 1250,
      "successful_operations": 1200,
      "failed_operations": 50
    },
    "resource_usage": {
      "cpu_usage": 45.2,
      "memory_usage": 68.5,
      "storage_usage": 50.0,

     "bandwidth_usage": 25.8
    },
    "members": [
      {
        "id": 1,
        "user_id": 1,
        "role": "admin",
        "status": "ACTIVE",
        "user": {
          "username": "admin",
          "nickname": "系统管理员"
        },
        "joined_at": "2025-01-01T00:00:00Z"
      }
    ],
    "recent_activities": [
      {
        "action": "APP_DEPLOY",
        "description": "部署应用：WordPress",
        "user": "developer",
        "timestamp": "2025-01-22T10:30:00Z"
      }
    ],
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-22T10:30:00Z"
  }
}
```

**更新项目**

```text
PUT /api/v1/admin/projects/{id}
```

请求体参数：

| 参数名                 | 数据类型 | 是否可空 | 描述         |
| ---------------------- | -------- | -------- | ------------ |
| name                   | string   | 是       | 项目名称     |
| description            | string   | 是       | 项目描述     |
| tags                   | array    | 是       | 项目标签     |
| icon                   | string   | 是       | 项目图标URL  |
| owner_id               | integer  | 是       | 项目所有者ID |
| status                 | string   | 是       | 项目状态     |
| default_resource_group | string   | 是       | 默认资源组   |
| default_timezone       | string   | 是       | 默认时区     |
| log_retention_days     | integer  | 是       | 日志保留天数 |
| backup_strategy        | string   | 是       | 备份策略     |
| access_control         | string   | 是       | 访问控制     |
| api_access             | boolean  | 是       | API访问权限  |
| audit_enabled          | boolean  | 是       | 审计日志开启 |

**删除项目**

```text
DELETE /api/v1/admin/projects/{id}
```

请求体参数：

| 参数名           | 数据类型 | 是否可空 | 描述                     |
| -orce            | boolean  | 是       | 是否强制删除，默认false  |
| backup_data      | boolean  | 是       | 删除前是否备份，默认true |
| transfer_owner   | integer  | 是       | 资源转移目标用户ID       |

响应示例：

```json
{
  "code": 200,
  "message": "项目删除成功",
  "data": {
    "project_id": 1,
    "status": "DELETED",
    "backup_location": "/backups/project_1_20250122.tar.gz",
    "deleted_at": "2025-01-22T10:30:00Z"
  }
}
```

**项目状态管理**

```text
PUT /api/v1/admin/projects/{id}/status
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述                                    |
| ------ | -------- | -------- | --------------------------------------- |
| status | string   | 否       | 项目状态（ACTIVE, INACTIVE, SUSPENDED） |
| reason | string   | 是       | 状态变更原因                            |

**批量操作项目**

```text
POST /api/v1/admin/projects/batch-actions
```

请求体参数：

| 参数名      | 数据类型  | 是否可空 | 描述                                              |
| ----------- | --------- | -------- | ------------------------------------------------- |
| project_ids | integer[] | 否       | 项目ID数组                                        |
| action      | string    | 否       | 操作类型（activate, deactivate, suspend, delete） |
| reason      | string    | 是       | 操作原因                                          |

#### 4.4.2 平台设置

**获取系统配置**

```text
GET /api/v1/system-configs
```

查询参数：

| 参数     | 类型   | 必填 | 默认值 | 描述                                   |
| -------- | ------ | ---- | ------ | -------------------------------------- |
| category | string | 否   | -      | 配置分类（basic, smtp, sms, security） |
| keyword  | string | 否   | -      | 搜索关键词                             |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "config_key": "system_name",
        "config_value": "Websoft9平台",
        "config_type": "STRING",
        "category": "basic",
        "description": "系统名称",
        "is_readonly": false,
        "is_encrypted": false,
        "default_value": "Websoft9",
        "sort_order": 1
      },
      {
        "id": 2,
        "config_key": "smtp_enabled",
        "config_value": "true",
        "config_type": "BOOLEAN",
        "category": "smtp",
        "description": "是否启用SMTP",
        "is_readonly": false,
        "is_encrypted": false,
        "default_value": "false",
        "sort_order": 1
      }
    ]
  }
}
```

**更新系统配置**

```text
PUT /api/v1/system-configs
```

请求体参数：

| 参数名       | 数据类型 | 是否可空 | 描述     |
| ------------ | -------- | -------- | -------- |
| config_key   | string   | 否       | 配置键   |
| config_value | string   | 否       | 配置值   |
| category     | string   | 否       | 配置分类 |

**测试SMTP配置**

```text
POST /api/v1/system-configs/smtp/test
```

请求体参数：

| 参数名     | 数据类型 | 是否可空 | 描述     |
| ---------- | -------- | -------- | -------- |
| test_email | string   | 否       | 测试邮箱 |

**测试短信配置**

```text
POST /api/v1/system-configs/sms/test
```

请求体参数：

| 参数名     | 数据类型 | 是否可空 | 描述       |
| ---------- | -------- | -------- | ---------- |
| test_phone | string   | 否       | 测试手机号 |

**获取License信息**

```text
GET /api/v1/system-configs/license
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "license_type": "ENTERPRISE",
    "company_name": "示例公司",
    "max_users": 100,
    "max_servers": 50,
    "max_apps": 500,
    "features": [
      "MONITORING",
      "WORKFLOW",
      "GATEWAY",
      "BACKUP"
    ],
    "issued_at": "2025-01-01T00:00:00Z",
    "expires_at": "2025-12-31T23:59:59Z",
    "days_remaining": 343,
    "is_valid": true,
    "usage": {
      "users": 25,
      "servers": 10,
      "apps": 125
    }
  }
}
```

**更新License**

```text
POST /api/v1/system-configs/license
```

请求体（multipart/form-data）：

```text
license_file: [License文件]
```

**获取Webhook配置**

```text
GET /api/v1/webhook-configs
```

**更新Webhook配置**

```text
PUT /api/v1/webhook-configs/{id}
```

**创建Webhook配置**

```text
POST /api/v1/webhook-configs
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述         |
| ----------- | -------- | -------- | ------------ |
| name        | string   | 否       | Webhook名称  |
| url         | string   | 否       | 回调URL      |
| secret      | string   | 是       | 签名密钥     |
| events      | array    | 否       | 监听事件     |
| headers     | object   | 是       | 自定义请求头 |
| timeout     | integer  | 是       | 超时时间(秒) |
| retry_count | integer  | 是       | 重试次数     |
| is_active   | boolean  | 是       | 是否启用     |
| owner_id    | integer  | 否       | 所有者ID     |

**获取软件源配置**

```text
GET /api/v1/system-configs/repositories
```

**更新软件源配置**

```text
PUT /api/v1/system-configs/repositories
```

请求体参数：

| 参数名          | 数据类型 | 是否可空 | 描述           |
| --------------- | -------- | -------- | -------------- |
| docker_registry | string   | 是       | Docker镜像仓库 |
| apt_mirror      | string   | 是       | APT软件源      |
| yum_mirror      | string   | 是       | YUM软件源      |

**获取系统状态**

```text
GET /api/v1/system-maintenance/status
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "system_info": {
      "version": "1.0.0",
      "build": "20250122",
      "uptime": 2592000,
      "environment": "production"
    },
    "service_status": {
      "api_server": "RUNNING",
      "database": "RUNNING",
      "redis": "RUNNING",
      "nginx": "RUNNING",
      "monitoring": "RUNNING"
    },
    "resource_usage": {
      "cpu_usage": 25.5,
      "memory_usage": 68.5,
      "disk_usage": 45.2,
      "network_io": {
        "in": 1024000,
        "out": 512000
      }
    },
    "health_checks": [
      {
        "name": "数据库连接",
        "status": "PASS",
        "response_time": 15
      },
      {
        "name": "Redis连接",
        "status": "PASS",
        "response_time": 5
      }
    ]
  }
}
```

**获取容器集群状态**

```text
GET /api/v1/system-maintenance/docker-swarm
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "cluster_info": {
      "id": "swarm_123456",
      "status": "ACTIVE",
      "node_count": 3,
      "manager_count": 1,
      "worker_count": 2
    },
    "nodes": [
      {
        "id": "node_123",
        "hostname": "manager-01",
        "role": "MANAGER",
        "status": "READY",
        "availability": "ACTIVE",
        "engine_version": "20.10.21"
      }
    ],
    "services": [
      {
        "id": "service_123",
        "name": "websoft9-api",
        "mode": "REPLICATED",
        "replicas": "2/2",
        "image": "websoft9/api:1.0.0",
        "status": "RUNNING"
      }
    ]
  }
}
```

**管理容器集群**

```text
POST /api/v1/system-maintenance/docker-swarm/actions
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述                                  |
| ------ | -------- | -------- | ------------------------------------- |
| action | string   | 否       | 操作类型（init, join, leave, update） |
| params | object   | 是       | 操作参数                              |

**获取容器镜像列表**

```text
GET /api/v1/system-maintenance/docker-images
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述                           |
| --------- | ------- | ---- | ------ | ------------------------------ |
| page      | integer | 否   | 1      | 页码                           |
| page_size | integer | 否   | 20     | 每页数量                       |
| keyword   | string  | 否   | -      | 搜索关键词                     |
| status    | string  | 否   | -      | 状态（USED, UNUSED, DANGLING） |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "sha256:abc123",
        "repository": "wordpress",
        "tag": "6.4.2",
        "size": 615000000,
        "created": "2025-07-15T10:30:00Z",
        "status": "USED",
        "containers": [
          {
            "id": "container_123",
            "name": "my-blog",
            "status": "RUNNING"
          }
        ]
      }
    ],
    "summary": {
      "total_images": 25,
      "used_images": 15,
      "unused_images": 8,
      "dangling_images": 2,
      "total_size": 15360000000
    }
  }
}
```

**清理容器镜像**

```text
POST /api/v1/system-maintenance/docker-images/cleanup
```

请求体参数：

| 参数名       | 数据类型 | 是否可空 | 描述                              |
| ------------ | -------- | -------- | --------------------------------- |
| cleanup_type | string   | 否       | 清理类型（UNUSED, DANGLING, ALL） |
| image_ids    | string[] | 是       | 指定镜像ID列表                    |
| force        | boolean  | 是       | 是否强制删除，默认false           |

**获取数据存储信息**

```text
GET /api/v1/system-maintenance/storage
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "workspace_storage": {
      "total_space": 107374182400,
      "used_space": 10737418240,
      "available_space": 96636764160,
      "usage_percentage": 10.0,
      "user_quotas": [
        {
          "user_id": 1,
          "username": "admin",
          "used_space": 1073741824,
          "quota": 10737418240,
          "usage_percentage": 10.0
        }
      ]
    },
    "database_storage": {
      "mysql_size": 2147483648,
      "redis_size": 104857600,
      "influxdb_size": 1073741824
    },
    "backup_storage": {
      "total_backups": 50,
      "total_size": 5368709120,
      "oldest_backup": "2025-01-01T00:00:00Z",
      "newest_backup": "2025-01-22T02:00:00Z"
    }
  }
}
```

**清理存储空间**

```text
POST /api/v1/system-maintenance/storage/cleanup
```

请求体参数：

| 参数名       | 数据类型 | 是否可空 | 描述                                   |
| ------------ | -------- | -------- | -------------------------------------- |
| cleanup_type | string   | 否       | 清理类型（TEMP, LOGS, BACKUPS, CACHE） |
| days_before  | integer  | 是       | 清理多少天前的数据，默认30             |

**获取系统任务列表**

```text
GET /api/v1/system-maintenance/tasks
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "数据库备份",
        "description": "定期备份系统数据库",
        "schedule": "0 2 * * *",
        "status": "ACTIVE",
        "last_run": "2025-01-22T02:00:00Z",
        "next_run": "2025-01-23T02:00:00Z",
        "last_result": "SUCCESS",
        "run_count": 30,
        "success_count": 30,
        "failure_count": 0
      }
    ]
  }
}
```

**手动执行系统任务**

```text
POST /api/v1/system-maintenance/tasks/{id}/execute
```

**获取系统更新信息**

```text
GET /api/v1/system-maintenance/updates
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "current_version": "1.0.0",
    "latest_version": "1.1.0",
    "has_updates": true,
    "update_info": {
      "version": "1.1.0",
      "release_date": "2025-01-25T00:00:00Z",
      "changelog": [
        "新增应用监控功能",
        "优化工作流执行性能",
        "修复已知安全问题"
      ],
      "download_url": "https://releases.websoft9.com/v1.1.0",
      "size": 104857600,
      "checksum": "sha256:abc123..."
    },
    "update_history": [
      {
        "version": "1.0.0",
        "updated_at": "2025-01-01T00:00:00Z",
        "status": "SUCCESS"
      }
    ]
  }
}
```

**执行系统更新**

```text
POST /api/v1/system-maintenance/updates
```

请求体参数：

| 参数名        | 数据类型 | 是否可空 | 描述                     |
| ------------- | -------- | -------- | ------------------------ |
| version       | string   | 否       | 目标版本号               |
| backup_before | boolean  | 是       | 更新前是否备份，默认true |

#### 4.4.3 安全管理

##### 角色管理

**获取角色列表**

```text
GET /api/v1/roles
```

**创建角色**

```text
POST /api/v1/roles
```

请求体参数：

| 参数名         | 数据类型 | 是否可空 | 描述                          |
| -------------- | -------- | -------- | ----------------------------- |
| name           | string   | 否       | 角色名称                      |
| code           | string   | 否       | 角色代码，唯一标识            |
| description    | string   | 是       | 角色描述                      |
| permission_ids | array    | 是       | 权限ID数组                    |
| is_system      | boolean  | 是       | 是否系统角色                  |
| sort_order     | integer  | 是       | 排序顺序，默认0               |
| status         | integer  | 是       | 状态（0-禁用，1-启用），默认1 |

**更新角色**

```text
PUT /api/v1/roles/{id}
```

**删除角色**

```text
DELETE /api/v1/roles/{id}
```

##### 权限管理

**获取权限列表**

```text
GET /api/v1/permissions
```

**获取权限树**

```text
GET /api/v1/permissions/tree
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "用户管理",
      "code": "user",
      "module": "user",
      "children": [
        {
          "id": 2,
          "name": "查看用户",
          "code": "user:read",
          "module": "user",
          "action": "read"
        },
        {
          "id": 3,
          "name": "创建用户",
          "code": "user:create",
          "module": "user",
          "action": "create"
        }
      ]
    }
  ]
}
```

##### 认证管理

**获取认证配置**

```text
GET /api/v1/auth-config
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "oauth2_enabled": true,
    "oauth2_providers": [
      {
        "name": "Google",
        "client_id": "***",
        "enabled": true
      }
    ],
    "two_factor_enabled": true,
    "two_factor_methods": ["TOTP", "EMAIL", "SMS"],
    "password_policy": {
      "min_length": 8,
      "require_uppercase": true,
      "require_lowercase": true,
      "require_numbers": true,
      "require_symbols": false
    },
    "session_config": {
      "timeout": 3600,
      "max_concurrent_sessions": 5
    }
  }
}
```

**更新认证配置**

```text
PUT /api/v1/auth-config
```

请求体参数：

| 参数名                  | 数据类型 | 是否可空 | 描述               |
| ----------------------- | -------- | -------- | ------------------ |
| oauth2_enabled          | boolean  | 是       | 是否启用OAuth2     |
| oauth2_providers        | object[] | 是       | OAuth2提供商配置   |
| two_factor_enabled      | boolean  | 是       | 是否启用双因子认证 |
| two_factor_methods      | string[] | 是       | 双因子认证方法     |
| password_policy         | object   | 是       | 密码策略配置       |
| session_timeout         | integer  | 是       | 会话超时时间（秒） |
| max_concurrent_sessions | integer  | 是       | 最大并发会话数     |

**获取API访问令牌列表**

```text
GET /api/v1/api-tokens
```

**创建API访问令牌**

```text
POST /api/v1/api-tokens
```

请求体参数：

| 参数名      | 数据类型 | 是否可空 | 描述                   |
| ----------- | -------- | -------- | ---------------------- |
| name        | string   | 否       | 令牌名称               |
| description | string   | 是       | 令牌描述               |
| scopes      | string[] | 否       | 权限范围               |
| expires_at  | string   | 是       | 过期时间，null表示永久 |

**撤销API访问令牌**

```text
DELETE /api/v1/api-tokens/{id}
```

#### 4.4.4 用户管理

**获取用户列表**

```text
GET /api/v1/users
```

查询参数：

| 参数      | 类型    | 必填 | 默认值     | 描述                             |
| --------- | ------- | ---- | ---------- | -------------------------------- |
| page      | integer | 否   | 1          | 页码                             |
| page_size | integer | 否   | 20         | 每页数量（1-100）                |
| keyword   | string  | 否   | -          | 搜索关键词（用户名、邮箱、昵称） |
| status    | integer | 否   | -          | 用户状态（0-禁用，1-启用）       |
| role_id   | integer | 否   | -          | 角色ID筛选                       |
| sort      | string  | 否   | created_at | 排序字段                         |
| order     | string  | 否   | desc       | 排序方向（asc, desc）            |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
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
        "last_login_at": "2025-07-15T10:30:00Z",
        "last_login_ip": "192.168.1.100",
        "roles": [
          {
            "id": 1,
            "name": "系统管理员",
            "code": "admin"
          }
        ],
        "created_at": "2025-07-15T10:30:00Z",
        "updated_at": "2025-07-15T10:30:00Z"
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

**获取用户详情**

```text
GET /api/v1/users/{id}
```

**创建用户**

```text
POST /api/v1/users
```

请求体参数：

| 参数名    | 数据类型 | 是否可空 | 描述                                 |
| --------- | -------- | -------- | ------------------------------------ |
| username  | string   | 否       | 用户名，3-32字符，字母数字下划线     |
| email     | string   | 否       | 邮箱地址，必须唯一                   |
| password  | string   | 否       | 密码，8-32字符，包含大小写字母和数字 |
| nickname  | string   | 是       | 用户昵称                             |
| phone     | string   | 是       | 手机号码                             |
| gender    | integer  | 是       | 性别（1-男，2-女，0-未知）           |
| signature | string   | 是       | 个性签名                             |
| timezone  | string   | 是       | 时区，默认UTC                        |
| language  | string   | 是       | 语言，默认zh-CN                      |
| role_ids  | array    | 是       | 角色ID数组                           |
| status    | integer  | 是       | 状态：0-禁用，1-启用                 |

验证规则：

- `username`: 必填，5-32字符，字母数字下划线
- `email`: 必填，有效邮箱格式，唯一
- `password`: 必填，8-32字符，包含大小写字母和数字
- `phone`: 可选，有效手机号格式

**更新用户**

```text
PUT /api/v1/users/{id}
```

**删除用户**

```text
DELETE /api/v1/users/{id}
```

**修改用户密码**

```text
PUT /api/v1/users/{id}/password
```

请求体参数：

| 参数名           | 数据类型 | 是否可空 | 描述       |
| ---------------- | -------- | -------- | ---------- |
| old_password     | string   | 否       | 原密码     |
| new_password     | string   | 否       | 新密码     |
| confirm_password | string   | 否       | 确认新密码 |

#### 4.4.5 告警通知

**获取告警规则列表**

```text
GET /api/v1/alert-rules
```

**创建告警规则**

```text
POST /api/v1/alert-rules
```

请求体参数：

| 参数名                | 数据类型 | 是否可空 | 描述                                       |
| --------------------- | -------- | -------- | ------------------------------------------ |
| name                  | string   | 否       | 告警规则名称                               |
| rule_type             | string   | 否       | 规则类型（METRIC, LOG, EVENT）             |
| target_type           | string   | 否       | 目标类型（SERVER, APP_INSTANCE, WORKFLOW） |
| target_id             | integer  | 是       | 目标ID                                     |
| metric_name           | string   | 是       | 监控指标名称                               |
| condition_expression  | string   | 否       | 条件表达式                                 |
| notification_channels | array    | 是       | 通知渠道配置                               |
| is_enabled            | boolean  | 是       | 是否启用，默认true                         |
| owner_id              | integer  | 否       | 所有者ID                                   |

**获取告警记录**

```text
GET /api/v1/alert-records
```

查询参数：

| 参数          | 类型    | 必填 | 默认值 | 描述                                |
| ------------- | ------- | ---- | ------ | ----------------------------------- |
| page          | integer | 否   | 1      | 页码                                |
| page_size     | integer | 否   | 20     | 每页数量                            |
| status        | string  | 否   | -      | 告警状态（FIRING, RESOLVED）        |
| severity      | string  | 否   | -      | 严重级别（CRITICAL, WARNING, INFO） |
| start_time    | string  | 否   | -      | 开始时间                            |
| end_time      | string  | 否   | -      | 结束时间                            |
| alert_rule_id | integer | 否   | -      | 告警规则ID                          |

**确认告警**

```text
PUT /api/v1/alert-records/{id}/acknowledge
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述     |
| ------ | -------- | -------- | -------- |
| note   | string   | 是       | 确认备注 |

**解决告警**

```text
PUT /api/v1/alert-records/{id}/resolve
```

请求体参数：

| 参数名          | 数据类型 | 是否可空 | 描述         |
| --------------- | -------- | -------- | ------------ |
| resolution_note | string   | 是       | 解决方案备注 |

#### 4.4.6 个人中心

**获取个人资料**

```text
GET /api/v1/profile
```

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "nickname": "系统管理员",
    "avatar": "https://example.com/avatars/admin.jpg",
    "phone": "13800138000",
    "gender": 1,
    "signature": "系统管理员账户",
    "timezone": "Asia/Shanghai",
    "language": "zh-CN",
    "two_factor_enabled": true,
    "notification_settings": {
      "email_notifications": true,
      "sms_notifications": false,
      "push_notifications": true,
      "marketing_emails": false
    },
    "security_settings": {
      "login_alerts": true,
      "session_timeout": 3600,
      "password_last_changed": "2025-07-15T10:30:00Z"
    },
    "last_login_at": "2025-07-15T10:30:00Z",
    "last_login_ip": "192.168.1.100",
    "created_at": "2025-07-15T10:30:00Z",
    "updated_at": "2025-07-15T10:30:00Z"
  }
}
```

**更新个人资料**

```text
PUT /api/v1/profile
```

请求体参数：

| 参数名    | 数据类型 | 是否可空 | 描述                       |
| --------- | -------- | -------- | -------------------------- |
| nickname  | string   | 是       | 昵称                       |
| avatar    | string   | 是       | 头像URL                    |
| phone     | string   | 是       | 手机号                     |
| gender    | integer  | 是       | 性别（0-未知，1-男，2-女） |
| signature | string   | 是       | 个性签名                   |
| timezone  | string   | 是       | 时区                       |
| language  | string   | 是       | 语言                       |

**上传头像**

```text
POST /api/v1/profile/avatar
```

请求体（multipart/form-data）：

```text
avatar: [头像文件]
```

响应示例：

```json
{
  "code": 200,
  "message": "头像上传成功",
  "data": {
    "avatar_url": "https://example.com/avatars/admin_new.jpg"
  }
}
```

**修改密码**

```text
PUT /api/v1/profile/password
```

请求体参数：

| 参数名           | 数据类型 | 是否可空 | 描述       |
| ---------------- | -------- | -------- | ---------- |
| old_password     | string   | 否       | 原密码     |
| new_password     | string   | 否       | 新密码     |
| confirm_password | string   | 否       | 确认新密码 |

**获取通知设置**

```text
GET /api/v1/profile/notification-settings
```

**更新通知设置**

```text
PUT /api/v1/profile/notification-settings
```

请求体参数：

| 参数名              | 数据类型 | 是否可空 | 描述     |
| ------------------- | -------- | -------- | -------- |
| email_notifications | boolean  | 是       | 邮件通知 |
| sms_notifications   | boolean  | 是       | 短信通知 |
| push_notifications  | boolean  | 是       | 推送通知 |
| marketing_emails    | boolean  | 是       | 营销邮件 |

**获取安全设置**

```text
GET /api/v1/profile/security-settings
```

**更新安全设置**

```text
PUT /api/v1/profile/security-settings
```

请求体参数：

| 参数名          | 数据类型 | 是否可空 | 描述               |
| --------------- | -------- | -------- | ------------------ |
| login_alerts    | boolean  | 是       | 登录提醒           |
| session_timeout | integer  | 是       | 会话超时时间（秒） |

**启用双因子认证**

```text
POST /api/v1/profile/two-factor/enable
```

响应示例：

```json
{
  "code": 200,
  "message": "请使用认证器扫描二维码",
  "data": {
    "qr_code": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
    "secret": "JBSWY3DPEHPK3PXP",
    "backup_codes": [
      "12345678",
      "87654321"
    ]
  }
}
```

**验证双因子认证**

```text
POST /api/v1/profile/two-factor/verify
```

请求体参数：

| 参数名 | 数据类型 | 是否可空 | 描述   |
| ------ | -------- | -------- | ------ |
| code   | string   | 否       | 验证码 |

**禁用双因子认证**

```text
POST /api/v1/profile/two-factor/disable
```

请求体参数：

| 参数名   | 数据类型 | 是否可空 | 描述     |
| -------- | -------- | -------- | -------- |
| password | string   | 否       | 当前密码 |
| code     | string   | 否       | 验证码   |

**获取登录历史**

```text
GET /api/v1/profile/login-history
```

查询参数：

| 参数      | 类型    | 必填 | 默认值 | 描述     |
| --------- | ------- | ---- | ------ | -------- |
| page      | integer | 否   | 1      | 页码     |
| page_size | integer | 否   | 20     | 每页数量 |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        "location": "北京市",
        "device": "Windows PC",
        "browser": "Chrome 120.0",
        "login_time": "2025-07-15T10:30:00Z",
        "logout_time": null,
        "status": "ACTIVE"
      }
    ]
  }
}
```

#### 4.4.7 审计日志

**获取审计日志列表**

```text
GET /api/v1/audit-logs
```

查询参数：

| 参数          | 类型    | 必填 | 默认值 | 描述     |
| ------------- | ------- | ---- | ------ | -------- |
| page          | integer | 否   | 1      | 页码     |
| page_size     | integer | 否   | 20     | 每页数量 |
| user_id       | integer | 否   | -      | 用户ID   |
| action        | string  | 否   | -      | 操作类型 |
| resource_type | string  | 否   | -      | 资源类型 |
| resource_id   | integer | 否   | -      | 资源ID   |
| start_time    | string  | 否   | -      | 开始时间 |
| end_time      | string  | 否   | -      | 结束时间 |
| ip_address    | string  | 否   | -      | IP地址   |

响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "user": {
          "id": 1,
          "username": "admin",
          "nickname": "管理员"
        },
        "action": "CREATE",
        "resource_type": "APP_INSTANCE",
        "resource_id": 123,
        "resource_name": "我的博客",
        "description": "创建应用实例：我的博客",
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        "request_data": {
          "name": "我的博客",
          "template_id": 1,
          "server_id": 1
        },
        "response_status": 200,
        "execution_time": 1500,
        "created_at": "2025-07-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

**获取审计日志详情**

```text
GET /api/v1/audit-logs/{id}
```

**获取审计统计**

```text
GET /api/v1/audit-logs/statistics
```

查询参数：

| 参数       | 类型   | 必填 | 默认值 | 描述                               |
| ---------- | ------ | ---- | ------ | ---------------------------------- |
| start_time | string | 否   | -      | 开始时间                           |
| end_time   | string | 否   | -      | 结束时间                           |
| group_by   | string | 否   | day    | 分组方式（hour, day, week, month） |

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_operations": 1250,
    "success_operations": 1200,
    "failed_operations": 50,
    "success_rate": 96.0,
    "top_users": [
      {
        "user_id": 1,
        "username": "admin",
        "operation_count": 500
      }
    ],
    "top_actions": [
      {
        "action": "READ",
        "count": 800
      },
      {
        "action": "CREATE",
        "count": 200
      }
    ],
    "timeline": [
      {
        "date": "2025-01-22",
        "count": 150
      }
    ]
  }
}
```

**导出审计日志**

```text
GET /api/v1/audit-logs/export
```

查询参数：

| 参数       | 类型    | 必填 | 默认值 | 描述                         |
| ---------- | ------- | ---- | ------ | ---------------------------- |
| format     | string  | 否   | csv    | 导出格式（csv, excel, json） |
| start_time | string  | 否   | -      | 开始时间                     |
| end_time   | string  | 否   | -      | 结束时间                     |
| user_id    | integer | 否   | -      | 用户ID                       |

响应：

```text
Content-Type: application/octet-stream
Content-Disposition: attachment; filename="audit-logs-20250122.csv"
```

## 5. 接口错误处理

### 5.1 标准错误码

| 错误码 | HTTP状态码 | 错误类型             | 描述           |
| ------ | ---------- | -------------------- | -------------- |
| 200    | 200        | SUCCESS              | 请求成功       |
| 400    | 400        | VALIDATION_ERROR     | 参数验证错误   |
| 401    | 401        | AUTHENTICATION_ERROR | 认证失败       |
| 403    | 403        | AUTHORIZATION_ERROR  | 权限不足       |
| 404    | 404        | NOT_FOUND_ERROR      | 资源不存在     |
| 409    | 409        | CONFLICT_ERROR       | 资源冲突       |
| 422    | 422        | BUSINESS_ERROR       | 业务逻辑错误   |
| 429    | 429        | RATE_LIMIT_ERROR     | 请求频率限制   |
| 500    | 500        | INTERNAL_ERROR       | 服务器内部错误 |
| 502    | 502        | GATEWAY_ERROR        | 网关错误       |
| 503    | 503        | SERVICE_UNAVAILABLE  | 服务不可用     |

### 5.2 详细错误码定义

#### 5.2.1 认证相关错误（1000-1999）

| 错误码 | 描述             |
| ------ | ---------------- |
| 1001   | 用户名或密码错误 |
| 1002   | Token已过期      |
| 1003   | Token无效        |
| 1004   | 账户已被禁用     |
| 1005   | 账户已被锁定     |
| 1006   | 密码强度不足     |
| 1007   | 验证码错误       |
| 1008   | 登录失败次数过多 |
| 1009   | 邮箱已存在       |
| 1010   | Token已使用      |

#### 5.2.2 权限相关错误（2000-2999）

| 错误码 | 描述           |
| ------ | -------------- |
| 2001   | 权限不足       |
| 2002   | 资源访问被拒绝 |
| 2003   | 操作权限不足   |
| 2004   | 角色权限不足   |
| 2005   | 资源组权限不足 |

#### 5.2.3 参数验证错误（3000-3999）

| 错误码 | 描述               |
| ------ | ------------------ |
| 3000   | 参数验证失败       |
| 3001   | 必填参数缺失       |
| 3002   | 参数格式错误       |
| 3003   | 参数值超出范围     |
| 3004   | 参数长度不符合要求 |
| 3005   | 邮箱格式错误       |
| 3006   | 手机号格式错误     |
| 3007   | URL格式错误        |
| 3008   | 日期格式错误       |
| 3009   | 邮箱未验证         |

#### 5.2.4 资源相关错误（4000-4999）

| 错误码 | 描述               |
| ------ | ------------------ |
| 4000   | 数据库记录不存在   |
| 4001   | 资源不存在         |
| 4002   | 资源已存在         |
| 4003   | 资源状态不允许操作 |
| 4004   | 资源依赖关系冲突   |
| 4005   | 资源配额不足       |
| 4006   | 资源正在使用中     |
| 4007   | 查询失败           |
| 4008   | 创建失败           |
| 4009   | 更新失败           |
| 4010   | 删除失败           |
| 4011   | 资源被禁用         |
| 4012   | 没有资源被更新     |
| 4013   | 不允许删除         |

#### 5.2.5 业务逻辑错误（5000-5999）

| 错误码 | 描述               |
| ------ | ------------------ |
| 5001   | 服务器离线无法操作 |
| 5002   | 应用部署失败       |
| 5003   | 工作流执行失败     |
| 5004   | 证书申请失败       |
| 5005   | 备份操作失败       |
| 5006   | 监控数据采集失败   |
| 5007   | 应用发布失败       |
| 5008   | 应用下线失败       |
| 5009   | 健康检查失败       |
| 5010   | 网关配置更新失败   |
| 5011   | 用户名已经存在     |
| 5012   | 权限验证无效       |

#### 5.2.6 系统相关错误（6000-6999）

| 错误码 | 描述             |
| ------ | ---------------- |
| 6001   | 系统内部错误     |
| 6002   | 缓存服务不可用   |
| 6003   | 文件系统错误     |
| 6004   | 网络连接超时     |
| 6005   | 第三方服务不可用 |
| 6006   | 系统维护中       |
| 6007   | 数据库连接失败   |

### 5.3 错误响应示例

#### 5.3.1 参数验证错误

```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": {
    "type": "VALIDATION_ERROR",
    "code": "INVALID_PARAMETER",
    "details": [
      {
        "field": "email",
        "message": "邮箱格式不正确",
        "code": "INVALID_EMAIL_FORMAT",
        "value": "invalid-email"
      },
      {
        "field": "password",
        "message": "密码长度至少8位",
        "code": "PASSWORD_TOO_SHORT",
        "value": "***"
      }
    ]
  }
}
```

#### 5.3.2 业务逻辑错误

```json
{
  "code": 422,
  "message": "应用部署失败",
  "error": {
    "type": "BUSINESS_ERROR",
    "code": "APP_DEPLOYMENT_FAILED",
    "details": [
      {
        "reason": "服务器资源不足",
        "code": "INSUFFICIENT_RESOURCES",
        "required_memory": "2Gi",
        "available_memory": "1Gi"
      }
    ]
  }
}
```

#### 5.3.3 系统错误

```json
{
  "code": 500,
  "message": "服务器内部错误",
  "error": {
    "type": "INTERNAL_ERROR",
    "code": "DATABASE_CONNECTION_FAILED",
    "details": [
      {
        "message": "数据库连接超时",
        "code": "CONNECTION_TIMEOUT"
      }
    ]
  }
}
```

## 6. 接口安全设计

### 6.1 输入验证

#### 6.1.1 参数验证规则

- **必填验证**：检查必填参数是否存在
- **类型验证**：验证参数数据类型（string, integer, boolean等）
- **格式验证**：验证邮箱、URL、日期等特殊格式
- **长度验证**：验证字符串长度范围
- **范围验证**：验证数值范围
- **正则验证**：使用正则表达式验证复杂格式

#### 6.1.2 SQL注入防护

- 使用参数化查询
- 对特殊字符进行转义
- 限制数据库用户权限
- 启用数据库审计日志

#### 6.1.3 XSS防护

- 对用户输入进行HTML转义
- 使用内容安全策略（CSP）
- 验证和过滤富文本内容
- 对输出内容进行编码

### 6.2 认证与授权

#### 6.2.1 JWT Token安全

- Token有效期设置：访问Token 24小时，刷新Token 7天
- Token签名算法：使用RS256非对称加密
- Token存储：客户端使用HttpOnly Cookie存储
- Token刷新：自动刷新机制，避免频繁登录

#### 6.2.2 权限控制

- 基于RBAC模型的细粒度权限控制
- 资源级权限验证
- 操作权限验证
- 数据权限隔离

#### 6.2.3 会话管理

- 会话超时机制
- 并发会话限制
- 异地登录检测
- 强制下线功能

### 6.3 数据加密

#### 6.3.1 传输加密

- 强制使用HTTPS（TLS 1.2+）
- 证书有效性验证
- HSTS安全头设置
- 敏感数据传输加密

#### 6.3.2 存储加密

- 密码使用bcrypt加密存储
- 敏感配置信息AES加密
- 数据库连接字符串加密
- 文件存储加密

#### 6.3.3 密钥管理

- 密钥轮换机制
- 密钥分离存储
- 密钥访问审计
- 硬件安全模块（HSM）支持（可选）
