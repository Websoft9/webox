# Websoft9 功能详细设计说明书 V1.1

**目录**

- [Websoft9 功能详细设计说明书 V1.1](#websoft9-功能详细设计说明书-v11)
  - [1. 引言](#1-引言)
  - [2. License功能设计](#2-license功能设计)
  - [3. License接口设计](#3-license接口设计)
  - [4. License数据库设计](#4-license数据库设计)
  - [5. 附录](#5-附录)

## 1. 引言

本文档为 **Websoft9架构升级** 的详细设计文档，本文档的编写目的在于明确 **Websoft9架构升级** 需求的开发途径以及应用方法，通过此文档，为此 **Websoft9架构升级** 需求的维护提供清晰、详细的设计，为下阶段开发工作的开展起到指导作用。

本文档的预期读者是 **Websoft9架构升级** 需求相关的业务人员、开发人员以及系统运维人员。

## 2. License功能设计

支持管理员管理平台License：

- 【License查看】
  - 当前License信息：授权用户数、功能模块、有效期。
  - License状态：有效/即将过期/已过期。

- 【License更新】
  - License文件上传：支持.lic格式文件上传。
  - License验证：验证License文件的有效性和完整性。
  - License激活：激活新的License并更新系统配置。
  - 备份恢复：支持License的备份和恢复。

- 【License监控】
  - 到期提醒：License即将到期时发送提醒通知。
  - 使用监控：监控License使用情况，防止超限使用。

## 3. License接口设计

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

## 4. License数据库设计

**系统配置表（system_configs）**

| 字段名        | 类型            | 约束               | 默认值                                        | 注释     |
| ------------- | --------------- | ------------------ | --------------------------------------------- | -------- |
| id            | BIGINT UNSIGNED | PK, AUTO_INCREMENT | -                                             | 配置ID   |
| config_key    | VARCHAR(64)     | NOT NULL, UNIQUE   | -                                             | 配置键   |
| config_value  | TEXT            | -                  | NULL                                          | 配置值   |
| config_type   | ENUM            | -                  | 'STRING'                                      | 配置类型 |
| category      | VARCHAR(32)     | NOT NULL           | -                                             | 配置分类 |
| description   | TEXT            | -                  | NULL                                          | 配置描述 |
| is_readonly   | TINYINT(1)      | -                  | 0                                             | 是否只读 |
| is_encrypted  | TINYINT(1)      | -                  | 0                                             | 是否加密 |
| default_value | TEXT            | -                  | NULL                                          | 默认值   |
| sort_order    | INT             | -                  | 0                                             | 排序     |
| created_at    | DATETIME        | NOT NULL           | CURRENT_TIMESTAMP                             | 创建时间 |
| updated_at    | DATETIME        | NOT NULL           | CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 更新时间 |

## 5. 附录

<!-- 这里写功能的所包含的枚举值、字典等定义，或者相关的参考文档引用 -->
