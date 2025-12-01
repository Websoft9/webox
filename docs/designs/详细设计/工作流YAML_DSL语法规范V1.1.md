# 工作流 YAML DSL 语法规范 V1.1

**说明**：本文档定义 Websoft9 工作流的 YAML DSL（Domain Specific Language）语法规范，为工作流编排提供标准化的定义语言。

## 文档元信息

- **负责人**：Websoft9 team
- **审核**：Websoft9 team
- **创建日期**：2025-12-01
- **版本**：V1.1
- **更新日期**：2025-12-01
- **参考文档**：工作流功能设计V1.4.md

---

## 目录

- [工作流 YAML DSL 语法规范 V1.1](#工作流-yaml-dsl-语法规范-v11)
  - [文档元信息](#文档元信息)
  - [目录](#目录)
  - [1. 概述](#1-概述)
    - [1.1 设计目标与原则](#11-设计目标与原则)
    - [1.2 与 GitHub Actions 语法的兼容性说明](#12-与-github-actions-语法的兼容性说明)
    - [1.3 YAML 基础语法要求](#13-yaml-基础语法要求)
    - [1.4 文件命名与存储规范](#14-文件命名与存储规范)
  - [2. 工作流顶级结构](#2-工作流顶级结构)
    - [2.1 完整结构概览](#21-完整结构概览)
    - [2.2 字段优先级与解析顺序](#22-字段优先级与解析顺序)
  - [3. 工作流元数据](#3-工作流元数据)
    - [3.1 name - 工作流名称](#31-name---工作流名称)
    - [3.2 version - 版本号](#32-version---版本号)
    - [3.3 description - 工作流描述](#33-description---工作流描述)
    - [3.4 author - 作者信息](#34-author---作者信息)
    - [3.5 tags - 标签列表](#35-tags---标签列表)
    - [3.6 category - 分类](#36-category---分类)
  - [4. 输入参数](#4-输入参数)
    - [4.1 variables - 变量定义](#41-variables---变量定义)
    - [4.2 变量名规范](#42-变量名规范)
    - [4.3 数据类型](#43-数据类型)
    - [4.4 验证规则](#44-验证规则)
    - [4.5 完整示例](#45-完整示例)
  - [5. 全局配置](#5-全局配置)
    - [5.1 env - 全局环境变量](#51-env---全局环境变量)
    - [5.2 defaults - 默认配置](#52-defaults---默认配置)
      - [5.2.1 defaults.run - 运行时默认配置](#521-defaultsrun---运行时默认配置)
      - [5.2.2 defaults.retry - 重试默认配置](#522-defaultsretry---重试默认配置)
      - [5.2.3 defaults.timeout - 超时默认配置](#523-defaultstimeout---超时默认配置)
    - [5.3 concurrency - 并发控制](#53-concurrency---并发控制)
      - [5.3.1 concurrency.group - 并发组](#531-concurrencygroup---并发组)
      - [5.3.2 concurrency.cancel-in-progress - 取消正在进行的执行](#532-concurrencycancel-in-progress---取消正在进行的执行)
  - [6. 步骤定义](#6-步骤定义)
    - [6.1 steps - 步骤列表](#61-steps---步骤列表)
    - [6.2 步骤通用字段](#62-步骤通用字段)
    - [6.3 关键字段说明](#63-关键字段说明)
  - [7. 上下文](#7-上下文)
    - [7.1 上下文概述](#71-上下文概述)
    - [7.2 常用上下文示例](#72-常用上下文示例)
  - [8. 表达式](#8-表达式)
    - [8.1 表达式语法](#81-表达式语法)
    - [8.2 运算符](#82-运算符)
    - [8.3 内置函数](#83-内置函数)
  - [9. 内置组件](#9-内置组件)
    - [9.1 组件概述](#91-组件概述)
    - [9.2 SHELL - Shell命令执行](#92-shell---shell命令执行)
    - [9.3 HTTP - HTTP请求](#93-http---http请求)
    - [9.4 EMAIL - 邮件发送](#94-email---邮件发送)
    - [9.5 COMPOSE - 应用部署](#95-compose---应用部署)
    - [9.6 FILE\_UPLOAD - 文件上传](#96-file_upload---文件上传)
    - [9.7 FILE\_DOWNLOAD - 文件下载](#97-file_download---文件下载)
    - [9.8 S3\_UPLOAD - S3对象上传](#98-s3_upload---s3对象上传)
    - [9.9 S3\_DOWNLOAD - S3对象下载](#99-s3_download---s3对象下载)
  - [10. 输出管理](#10-输出管理)
    - [10.1 步骤输出定义](#101-步骤输出定义)
    - [10.2 输出引用](#102-输出引用)
    - [10.3 输出大小限制](#103-输出大小限制)
  - [11. 自定义组件](#11-自定义组件)
    - [11.1 组件引用语法](#111-组件引用语法)
    - [11.2 本地组件引用](#112-本地组件引用)
    - [11.3 远程组件引用](#113-远程组件引用)
    - [11.4 组件版本管理](#114-组件版本管理)
  - [12. 错误处理](#12-错误处理)
    - [12.1 步骤级错误处理](#121-步骤级错误处理)
    - [12.2 工作流级错误处理](#122-工作流级错误处理)
  - [13. 安全规范](#13-安全规范)
    - [13.1 密钥引用规范](#131-密钥引用规范)
      - [13.1.1 基本引用语法](#1311-基本引用语法)
      - [13.1.2 密钥类型和字段访问](#1312-密钥类型和字段访问)
      - [13.1.3 密钥引用位置](#1313-密钥引用位置)
    - [13.2 敏感数据处理](#132-敏感数据处理)
      - [13.2.1 日志脱敏](#1321-日志脱敏)
      - [13.2.2 环境变量保护](#1322-环境变量保护)
      - [13.2.3 密钥安全规范](#1323-密钥安全规范)
    - [13.3 权限最小化原则](#133-权限最小化原则)
      - [13.3.1 文件权限](#1331-文件权限)
      - [13.3.2 命令执行限制](#1332-命令执行限制)
      - [13.3.3 网络访问控制](#1333-网络访问控制)
  - [14. 触发器配置](#14-触发器配置)
    - [14.1 触发器语法规划](#141-触发器语法规划)
    - [14.2 定时触发配置](#142-定时触发配置)
    - [14.3 Webhook 触发配置](#143-webhook-触发配置)
    - [14.4 手动触发配置](#144-手动触发配置)
  - [15. GitHub Action 组件](#15-github-action-组件)
    - [15.1 组件类型](#151-组件类型)
    - [15.2 引用语法](#152-引用语法)
    - [15.3 常用 Action 示例](#153-常用-action-示例)
    - [15.4 与自定义组件混合使用](#154-与自定义组件混合使用)
    - [15.5 输入输出处理](#155-输入输出处理)
    - [15.6 环境变量共享](#156-环境变量共享)
    - [15.7 限制和约束](#157-限制和约束)
  - [16. 语法限制与约束](#16-语法限制与约束)
    - [16.1 文件大小限制](#161-文件大小限制)
    - [16.2 嵌套深度限制](#162-嵌套深度限制)
    - [16.3 字符串长度限制](#163-字符串长度限制)
    - [16.4 数量限制](#164-数量限制)
    - [16.5 执行时间限制](#165-执行时间限制)
    - [16.6 输出大小限制](#166-输出大小限制)
  - [17. 完整示例](#17-完整示例)
    - [17.1 应用自动部署工作流](#171-应用自动部署工作流)
    - [17.2 数据库定时备份工作流](#172-数据库定时备份工作流)
  - [18. 版本兼容性](#18-版本兼容性)
    - [18.1 DSL 版本声明](#181-dsl-版本声明)
    - [18.2 版本迁移指南](#182-版本迁移指南)
    - [18.3 废弃语法说明](#183-废弃语法说明)
  - [19. 附录](#19-附录)
    - [19.1 YAML Schema 定义](#191-yaml-schema-定义)
    - [19.2 语法速查表](#192-语法速查表)
    - [19.3 保留关键字列表](#193-保留关键字列表)
    - [19.4 错误代码参考](#194-错误代码参考)
    - [19.5 与 GitHub Actions 语法对照表](#195-与-github-actions-语法对照表)
  - [参考文档](#参考文档)
  - [变更记录](#变更记录)

---

## 1. 概述

### 1.1 设计目标与原则

**设计目标**：

- **简洁易读**：使用 YAML 格式，语法简洁，易于理解和维护
- **功能完整**：支持复杂的工作流编排，包括条件分支、重试、并行执行等
- **类型安全**：提供严格的类型定义和验证机制
- **可扩展性**：支持自定义组件和插件扩展
- **兼容性**：参考 GitHub Actions 语法，降低学习成本

**设计原则**：

- **声明式**：描述"做什么"而不是"怎么做"
- **可组合**：组件可以灵活组合和复用
- **可测试**：支持工作流的单元测试和集成测试
- **可观测**：提供完整的执行日志和监控指标

### 1.2 与 GitHub Actions 语法的兼容性说明

Websoft9 工作流 DSL 参考了 GitHub Actions 的语法设计，但根据平台特性进行了调整：

**相似之处**：

- 使用 YAML 格式定义工作流
- 支持 `${{ }}` 表达式语法
- 支持 `steps` 步骤定义
- 支持 `env` 环境变量
- 支持 `if` 条件执行
- 支持 `outputs` 输出定义

**差异之处**：

| 特性     | GitHub Actions        | Websoft9 Workflow      |
| -------- | --------------------- | ---------------------- |
| 触发器   | `on` 字段             | 通过任务调度配置       |
| 作业     | `jobs` 多作业并行     | `steps` 顺序执行       |
| 运行环境 | `runs-on` 指定 Runner | `target` 指定服务器    |
| 组件引用 | `uses: actions/xxx`   | `type: COMPONENT_TYPE` |
| 密钥引用 | `secrets.XXX`         | `secrets.XXX`          |

### 1.3 YAML 基础语法要求

**版本和编码**：

- YAML 版本：1.2
- 文件编码：UTF-8（无 BOM）
- 文件扩展名：`.yaml` 或 `.yml`

**缩进规则**：

- 使用 **2 个空格** 进行缩进
- **禁止使用 Tab 字符**
- 同级元素必须保持相同缩进

**数据类型示例**：

```yaml
# 字符串
name: "工作流名称"
description: 工作流描述

# 数字
timeout: 300
retry_count: 3

# 布尔值
enabled: true
debug: false

# 列表
tags:
  - deployment
  - production

# 对象
config:
  key1: value1
  key2: value2

# 多行字符串（保留换行）
script: |
  #!/bin/bash
  echo "Line 1"
  echo "Line 2"

# 多行字符串（折叠换行）
description: >
  这是一个很长的描述，
  会被折叠成一行。
```

### 1.4 文件命名与存储规范

**文件命名规范**：

- 文件名格式：`{workflow_code}.yaml`
- 使用小写字母、数字和连字符
- 示例：`app-deployment.yaml`、`data-backup.yaml`

**存储位置**：

- 工作流定义存储在项目空间中
- 路径：`/project/{project_id}/workflows/{workflow_code}.yaml`
- 数据库中存储 YAML 内容的字符串形式

---

## 2. 工作流顶级结构

### 2.1 完整结构概览

```yaml
# ============ 工作流元数据（必填） ============
name: string                    # 工作流名称
version: string                 # 版本号
description: string             # 工作流描述

# ============ 可选元数据 ============
author: string                  # 作者
tags: array                     # 标签列表
category: string                # 分类

# ============ 输入参数（可选） ============
variables:
  - name: string                # 变量名
    type: string                # 数据类型
    value: any                  # 默认值
    required: boolean           # 是否必填
    description: string         # 变量描述
    validation: object          # 验证规则

# ============ 全局配置（可选） ============
env:                            # 全局环境变量
  KEY: value

# ============ 步骤定义（必填） ============
steps:
  - id: string                  # 步骤唯一标识
    target: string              # 目标服务器
    name: string                # 步骤名称
    type: string                # 组件类型
    if: string                  # 条件表达式
    timeout: integer            # 超时时间（秒）
    retry: object               # 重试配置
    with: object                # 组件配置参数
    env: object                 # 步骤环境变量
    outputs: object             # 输出变量定义
```

### 2.2 字段优先级与解析顺序

**解析顺序**：

1. 元数据解析：解析 `name`、`version`、`description` 等元数据
2. 变量定义解析：解析 `variables` 定义，构建变量上下文
3. 全局配置解析：解析 `env` 全局环境变量
4. 步骤定义解析：解析 `steps` 列表，构建执行计划
5. 依赖关系分析：分析步骤间的依赖关系
6. 表达式预计算：预计算常量表达式

**字段优先级**：

- 步骤级配置 > 全局配置
- 步骤 `env` > 全局 `env`
- 步骤 `timeout` > 默认 `timeout`

---

## 3. 工作流元数据

### 3.1 name - 工作流名称

**说明**：工作流的显示名称，用于界面展示和日志记录。

**类型**：`string`

**约束**：

- 必填
- 长度：1-64 字符
- 支持中文、英文、数字、空格和常用符号

**示例**：

```yaml
name: "应用自动部署"
name: "Database Backup Workflow"
```

### 3.2 version - 版本号

**说明**：工作流的版本号，遵循语义化版本规范。

**类型**：`string`

**约束**：

- 必填
- 格式：`major.minor.patch`（如 `1.0.0`）
- 遵循 Semantic Versioning 2.0.0

**示例**：

```yaml
version: "1.0.0"
version: "2.1.3"
```

### 3.3 description - 工作流描述

**说明**：工作流的详细描述，说明工作流的用途和功能。

**类型**：`string`

**约束**：

- 必填
- 长度：1-500 字符
- 支持多行文本

**示例**：

```yaml
description: "自动化应用部署工作流，从Git仓库拉取代码，构建Docker镜像，部署到生产环境"

description: >
  这是一个数据库备份工作流，
  每天凌晨2点自动执行，
  备份数据库并上传到S3存储。
```

### 3.4 author - 作者信息

**说明**：工作流的创建者或维护者信息。

**类型**：`string`

**约束**：

- 可选
- 长度：1-64 字符

**示例**：

```yaml
author: "DevOps Team"
author: "张三 <zhangsan@example.com>"
```

### 3.5 tags - 标签列表

**说明**：工作流的分类标签，用于搜索和过滤。

**类型**：`array<string>`

**约束**：

- 可选
- 每个标签长度：1-32 字符
- 最多 10 个标签

**示例**：

```yaml
tags:
  - deployment
  - production
  - docker
  - automation
```

### 3.6 category - 分类

**说明**：工作流的业务分类。

**类型**：`string`

**约束**：

- 可选
- 预定义分类：`deployment`、`backup`、`monitoring`、`maintenance`、`testing`、`other`

**示例**：

```yaml
category: deployment
category: backup
```

---

## 4. 输入参数

### 4.1 variables - 变量定义

**说明**：定义工作流的输入参数，支持类型验证和默认值。

**类型**：`array<object>`

**结构**：

```yaml
variables:
  - name: string                # 变量名（必填）
    type: string                # 数据类型（必填）
    value: any                  # 默认值（可选）
    required: boolean           # 是否必填（可选，默认false）
    description: string         # 变量描述（可选）
    validation: object          # 验证规则（可选）
```

### 4.2 变量名规范

**约束**：

- 使用大写字母、数字和下划线
- 必须以字母开头
- 长度：1-64 字符
- 示例：`APP_NAME`、`SERVER_ID`、`ENVIRONMENT`

### 4.3 数据类型

**支持的类型**：

| 类型      | 说明                 | 示例值               |
| --------- | -------------------- | -------------------- |
| `string`  | 字符串               | `"myapp"`            |
| `number`  | 数字（整数或浮点数） | `123`、`3.14`        |
| `boolean` | 布尔值               | `true`、`false`      |
| `object`  | 对象                 | `{"key": "value"}`   |
| `array`   | 数组                 | `["item1", "item2"]` |

### 4.4 验证规则

**validation 对象结构**：

```yaml
validation:
  pattern: string               # 正则表达式（string类型）
  min: number                   # 最小值/最小长度
  max: number                   # 最大值/最大长度
  enum: array                   # 枚举值列表
```

### 4.5 完整示例

```yaml
variables:
  # 字符串类型，带正则验证
  - name: APP_NAME
    type: string
    value: "myapp"
    required: true
    description: "应用名称"
    validation:
      pattern: "^[a-z0-9-]+$"
      min: 1
      max: 32

  # 数字类型，带范围验证
  - name: SERVER_ID
    type: number
    value: 1
    required: true
    description: "目标服务器ID"
    validation:
      min: 1
      max: 1000

  # 枚举类型
  - name: ENVIRONMENT
    type: string
    value: "production"
    required: true
    description: "部署环境"
    validation:
      enum:
        - development
        - staging
        - production

  # 布尔类型
  - name: DEBUG_MODE
    type: boolean
    value: false
    required: false
    description: "是否启用调试模式"
```

---

## 5. 全局配置

### 5.1 env - 全局环境变量

**说明**：定义全局环境变量，所有步骤都可以访问。

**类型**：`object`

**约束**：

- 可选
- 键名：大写字母、数字和下划线
- 值：字符串类型

**示例**：

```yaml
env:
  DOCKER_REGISTRY: "registry.example.com"
  NOTIFICATION_EMAIL: "ops@example.com"
  LOG_LEVEL: "INFO"
```

**使用方式**：

```yaml
steps:
  - id: build_image
    type: SHELL
    with:
      command: |
        docker build -t ${{ env.DOCKER_REGISTRY }}/myapp:latest .
```

### 5.2 defaults - 默认配置

**说明**：定义工作流的默认配置，可被步骤级配置覆盖。

**类型**：`object`（可选）

#### 5.2.1 defaults.run - 运行时默认配置

**说明**：定义步骤执行的默认配置。

**结构**：

```yaml
defaults:
  run:
    shell: /bin/bash           # 默认Shell类型
    working_directory: /tmp    # 默认工作目录
```

#### 5.2.2 defaults.retry - 重试默认配置

**说明**：定义步骤失败时的默认重试策略。

**结构**：

```yaml
defaults:
  retry:
    max_attempts: 3
    initial_interval: 1
    backoff_coefficient: 2.0
    maximum_interval: 60
```

#### 5.2.3 defaults.timeout - 超时默认配置

**说明**：定义步骤执行的默认超时时间。

**结构**：

```yaml
defaults:
  timeout: 300  # 默认超时时间（秒）
```

**完整示例**：

```yaml
defaults:
  run:
    shell: /bin/bash
    working_directory: /opt/apps
  retry:
    max_attempts: 3
    initial_interval: 1
    backoff_coefficient: 2.0
  timeout: 600
```

### 5.3 concurrency - 并发控制

**说明**：控制工作流的并发执行行为（未来版本支持）。

**类型**：`object`（可选）

#### 5.3.1 concurrency.group - 并发组

**说明**：定义并发组，同一组内的工作流不会并发执行。

**示例**：

```yaml
concurrency:
  group: deployment-${{ variables.ENVIRONMENT }}
```

#### 5.3.2 concurrency.cancel-in-progress - 取消正在进行的执行

**说明**：当新的执行开始时，是否取消同组内正在进行的执行。

**示例**：

```yaml
concurrency:
  group: deployment-${{ variables.ENVIRONMENT }}
  cancel_in_progress: true
```

**注意**：并发控制功能在 V1.0 版本中暂不支持，将在未来版本中实现。

---

## 6. 步骤定义

### 6.1 steps - 步骤列表

**说明**：定义工作流的执行步骤，按顺序执行。

**类型**：`array<object>`

**约束**：

- 必填
- 至少包含 1 个步骤
- 最多 100 个步骤

### 6.2 步骤通用字段

**完整结构**：

```yaml
steps:
  - id: string                  # 步骤唯一标识（必填）
    target: string              # 目标服务器（必填）
    name: string                # 步骤名称（必填）
    type: string                # 组件类型（必填）
    if: string                  # 条件表达式（可选）
    timeout: integer            # 超时时间（秒）（可选，默认300）
    retry: object               # 重试配置（可选）
    with: object                # 组件配置参数（必填）
    env: object                 # 步骤环境变量（可选）
    outputs: object             # 输出变量定义（可选）
```

### 6.3 关键字段说明

**id - 步骤唯一标识**：

- 使用小写字母、数字和下划线
- 必须以字母开头
- 长度：1-64 字符
- 工作流内唯一

**target - 目标服务器**：

```yaml
# 固定服务器
target: "server1"

# 变量引用
target: ${{ variables.SERVER_ID }}
```

**type - 组件类型**：

支持的组件类型：`SHELL`、`HTTP`、`EMAIL`、`FILE_UPLOAD`、`FILE_DOWNLOAD`、`COMPOSE`、`S3_UPLOAD`、`S3_DOWNLOAD`

**if - 条件执行**：

```yaml
# 基于变量条件
if: ${{ variables.ENVIRONMENT == 'production' }}

# 基于步骤状态
if: ${{ success() }}
if: ${{ failure() }}
if: ${{ always() }}
```

**retry - 重试配置**：

```yaml
retry:
  max_attempts: 3               # 最大重试次数（1-10）
  initial_interval: 1           # 初始重试间隔（秒）
  backoff_coefficient: 2.0      # 退避系数
  maximum_interval: 60          # 最大重试间隔（秒）
```

---

## 7. 上下文

### 7.1 上下文概述

**说明**：上下文是工作流执行过程中可以访问的数据集合，通过 `${{ context.property }}` 语法访问。

**支持的上下文**：

| 上下文      | 说明           | 示例                                   |
| ----------- | -------------- | -------------------------------------- |
| `variables` | 输入参数       | `${{ variables.APP_NAME }}`            |
| `env`       | 环境变量       | `${{ env.DOCKER_REGISTRY }}`           |
| `secrets`   | 密钥凭据       | `${{ secrets.API_TOKEN }}`             |
| `workflow`  | 工作流信息     | `${{ workflow.id }}`                   |
| `execution` | 执行实例信息   | `${{ execution.id }}`                  |
| `steps`     | 步骤信息与输出 | `${{ steps.build.outputs.image_tag }}` |
| `project`   | 项目信息       | `${{ project.id }}`                    |

### 7.2 常用上下文示例

**variables 上下文**：

```yaml
variables:
  - name: APP_NAME
    type: string
    value: "myapp"

steps:
  - id: deploy
    type: COMPOSE
    with:
      app_name: ${{ variables.APP_NAME }}
```

**secrets 上下文**：

```yaml
steps:
  - id: clone_repo
    type: SHELL
    env:
      GIT_TOKEN: ${{ secrets.GIT_TOKEN }}
    with:
      command: git clone https://github.com/example/repo.git
```

**steps 上下文**：

```yaml
steps:
  - id: build
    type: SHELL
    with:
      command: |
        echo "BUILD_ID=12345" >> $GITHUB_OUTPUT
    outputs:
      build_id:
        value: ${{ steps.build.result.BUILD_ID }}

  - id: deploy
    type: COMPOSE
    with:
      image_tag: ${{ steps.build.outputs.build_id }}
```

---

## 8. 表达式

### 8.1 表达式语法

**说明**：表达式使用 `${{ <expression> }}` 语法，用于动态计算值。

**基本语法**：

```yaml
# 简单引用
${{ variables.APP_NAME }}

# 属性访问
${{ steps.build.outputs.image_tag }}

# 函数调用
${{ contains(variables.TAGS, 'production') }}

# 运算符
${{ variables.COUNT > 10 }}

# 复杂表达式
${{ variables.ENV == 'prod' && steps.build.status == 'success' }}
```

### 8.2 运算符

**比较运算符**：`==`、`!=`、`<`、`>`、`<=`、`>=`

**逻辑运算符**：`&&`、`||`、`!`

**示例**：

```yaml
if: ${{ variables.ENVIRONMENT == 'production' }}
if: ${{ variables.ENV == 'prod' && steps.build.status == 'success' }}
```

### 8.3 内置函数

**字符串函数**：

| 函数                      | 说明                   | 示例                                        |
| ------------------------- | ---------------------- | ------------------------------------------- |
| `contains(str, substr)`   | 检查是否包含子字符串   | `${{ contains(variables.TEXT, 'hello') }}`  |
| `startsWith(str, prefix)` | 检查是否以指定前缀开头 | `${{ startsWith(variables.FILE, 'app-') }}` |
| `endsWith(str, suffix)`   | 检查是否以指定后缀结尾 | `${{ endsWith(variables.FILE, '.yaml') }}`  |
| `toUpper(str)`            | 转换为大写             | `${{ toUpper(variables.TEXT) }}`            |
| `toLower(str)`            | 转换为小写             | `${{ toLower(variables.TEXT) }}`            |

**类型转换函数**：

| 函数            | 说明             | 示例                                  |
| --------------- | ---------------- | ------------------------------------- |
| `toJSON(obj)`   | 转换为JSON字符串 | `${{ toJSON(variables.CONFIG) }}`     |
| `fromJSON(str)` | 从JSON字符串解析 | `${{ fromJSON(variables.JSON_STR) }}` |

**状态检查函数**：

| 函数        | 说明                       | 使用场景         |
| ----------- | -------------------------- | ---------------- |
| `success()` | 检查前面的步骤是否全部成功 | 成功后执行的步骤 |
| `failure()` | 检查前面的步骤是否有失败   | 失败后执行的步骤 |
| `always()`  | 总是返回true               | 清理步骤         |

**示例**：

```yaml
steps:
  - id: build
    type: SHELL
    with:
      command: ./build.sh

  - id: notify_success
    type: EMAIL
    if: ${{ success() }}
    with:
      subject: "构建成功"

  - id: cleanup
    type: SHELL
    if: ${{ always() }}
    with:
      command: rm -rf /tmp/build
```

---

## 9. 内置组件

### 9.1 组件概述

**说明**：组件是工作流步骤的执行单元，每个组件实现特定的功能。

**组件分类**：

- **应用组件**：SHELL、HTTP、EMAIL、FILE_UPLOAD、FILE_DOWNLOAD、COMPOSE
- **云服务组件**：S3_UPLOAD、S3_DOWNLOAD

### 9.2 SHELL - Shell命令执行

**说明**：在目标服务器上执行Shell命令或脚本。

**配置参数（with）**：

| 参数                | 类型   | 必填 | 说明                 |
| ------------------- | ------ | ---- | -------------------- |
| `command`           | string | 是   | Shell命令或脚本      |
| `working_directory` | string | 否   | 工作目录（默认/tmp） |
| `server_id`         | number | 是   | 目标服务器ID         |

**输出结果（result）**：

| 字段        | 类型   | 说明     |
| ----------- | ------ | -------- |
| `exit_code` | number | 退出码   |
| `stdout`    | string | 标准输出 |
| `stderr`    | string | 标准错误 |

**示例**：

```yaml
- id: run_script
  name: "执行部署脚本"
  target: "server1"
  type: SHELL
  timeout: 600
  with:
    command: |
      #!/bin/bash
      set -e
      echo "Starting deployment..."
      ./deploy.sh
      echo "DEPLOY_TIME=$(date +%Y%m%d%H%M%S)" >> $GITHUB_OUTPUT
    working_directory: /opt/apps
    server_id: ${{ variables.SERVER_ID }}
  env:
    APP_NAME: ${{ variables.APP_NAME }}
  outputs:
    deploy_time:
      value: ${{ steps.run_script.result.DEPLOY_TIME }}
```

### 9.3 HTTP - HTTP请求

**配置参数（with）**：

| 参数      | 类型   | 必填 | 说明                              |
| --------- | ------ | ---- | --------------------------------- |
| `method`  | string | 是   | HTTP方法（GET/POST/PUT/DELETE等） |
| `url`     | string | 是   | 请求URL                           |
| `headers` | object | 否   | 请求头                            |
| `body`    | string | 否   | 请求体                            |

**示例**：

```yaml
- id: api_call
  name: "调用部署API"
  target: "server1"
  type: HTTP
  with:
    method: POST
    url: https://api.example.com/deploy
    headers:
      Content-Type: application/json
      Authorization: Bearer ${{ secrets.API_TOKEN }}
    body: |
      {
        "app": "${{ variables.APP_NAME }}",
        "version": "${{ variables.VERSION }}"
      }
```

### 9.4 EMAIL - 邮件发送

**配置参数（with）**：

| 参数            | 类型   | 必填 | 说明           |
| --------------- | ------ | ---- | -------------- |
| `smtp_server`   | string | 是   | SMTP服务器地址 |
| `smtp_port`     | number | 是   | SMTP端口       |
| `smtp_user`     | string | 是   | SMTP用户名     |
| `smtp_password` | string | 是   | SMTP密码       |
| `from`          | string | 是   | 发件人邮箱     |
| `to`            | array  | 是   | 收件人邮箱列表 |
| `subject`       | string | 是   | 邮件主题       |
| `body`          | string | 是   | 邮件正文       |

**示例**：

```yaml
- id: send_notification
  name: "发送部署通知"
  target: "server1"
  type: EMAIL
  if: ${{ success() }}
  with:
    smtp_server: smtp.example.com
    smtp_port: 587
    smtp_user: ${{ secrets.SMTP_USER }}
    smtp_password: ${{ secrets.SMTP_PASSWORD }}
    from: noreply@example.com
    to:
      - admin@example.com
      - ops@example.com
    subject: "部署成功: ${{ variables.APP_NAME }}"
    body: |
      应用 ${{ variables.APP_NAME }} 已成功部署。
      执行ID: ${{ execution.id }}
```

### 9.5 COMPOSE - 应用部署

**配置参数（with）**：

| 参数          | 类型    | 必填 | 说明                     |
| ------------- | ------- | ---- | ------------------------ |
| `app_name`    | string  | 是   | 应用名称                 |
| `server_id`   | number  | 是   | 目标服务器ID             |
| `pull_images` | boolean | 否   | 是否拉取镜像（默认true） |

**示例**：

```yaml
- id: deploy_app
  name: "部署应用"
  target: "server1"
  type: COMPOSE
  with:
    app_name: ${{ variables.APP_NAME }}
    server_id: ${{ variables.SERVER_ID }}
    pull_images: true
```

### 9.6 FILE_UPLOAD - 文件上传

**说明**：上传文件到目标服务器。

**配置参数（with）**：

| 参数               | 类型    | 必填 | 说明                  |
| ------------------ | ------- | ---- | --------------------- |
| `source_path`      | string  | 是   | 源文件路径            |
| `destination_path` | string  | 是   | 目标文件路径          |
| `server_id`        | number  | 是   | 目标服务器ID          |
| `overwrite`        | boolean | 否   | 是否覆盖（默认false） |
| `permissions`      | string  | 否   | 文件权限（如"0644"）  |

**示例**：

```yaml
- id: upload_config
  name: "上传配置文件"
  target: "server1"
  type: FILE_UPLOAD
  with:
    source_path: /tmp/config.yaml
    destination_path: /etc/app/config.yaml
    server_id: ${{ variables.SERVER_ID }}
    overwrite: true
    permissions: "0644"
```

### 9.7 FILE_DOWNLOAD - 文件下载

**说明**：从URL下载文件到目标服务器。

**配置参数（with）**：

| 参数               | 类型    | 必填 | 说明                    |
| ------------------ | ------- | ---- | ----------------------- |
| `source_url`       | string  | 是   | 源文件URL               |
| `destination_path` | string  | 是   | 目标文件路径            |
| `server_id`        | number  | 是   | 目标服务器ID            |
| `verify_ssl`       | boolean | 否   | 是否验证SSL（默认true） |

**示例**：

```yaml
- id: download_package
  name: "下载安装包"
  target: "server1"
  type: FILE_DOWNLOAD
  with:
    source_url: https://releases.example.com/myapp.tar.gz
    destination_path: /tmp/myapp.tar.gz
    server_id: ${{ variables.SERVER_ID }}
    verify_ssl: true
```

### 9.8 S3_UPLOAD - S3对象上传

**说明**：上传文件到S3兼容的对象存储。

**配置参数（with）**：

| 参数              | 类型   | 必填 | 说明       |
| ----------------- | ------ | ---- | ---------- |
| `access_key`      | string | 是   | Access Key |
| `secret_key`      | string | 是   | Secret Key |
| `region`          | string | 是   | 区域       |
| `bucket`          | string | 是   | 存储桶名称 |
| `source_path`     | string | 是   | 源文件路径 |
| `destination_key` | string | 是   | 目标对象键 |

**示例**：

```yaml
- id: upload_backup
  name: "上传备份到S3"
  target: "server1"
  type: S3_UPLOAD
  with:
    access_key: ${{ secrets.AWS_ACCESS_KEY }}
    secret_key: ${{ secrets.AWS_SECRET_KEY }}
    region: us-east-1
    bucket: myapp-backups
    source_path: /tmp/backup.tar.gz
    destination_key: backups/backup-${{ execution.id }}.tar.gz
```

### 9.9 S3_DOWNLOAD - S3对象下载

**说明**：从S3兼容的对象存储下载文件。

**配置参数（with）**：

| 参数               | 类型   | 必填 | 说明         |
| ------------------ | ------ | ---- | ------------ |
| `access_key`       | string | 是   | Access Key   |
| `secret_key`       | string | 是   | Secret Key   |
| `region`           | string | 是   | 区域         |
| `bucket`           | string | 是   | 存储桶名称   |
| `source_key`       | string | 是   | 源对象键     |
| `destination_path` | string | 是   | 目标文件路径 |

**示例**：

```yaml
- id: download_backup
  name: "从S3下载备份"
  target: "server1"
  type: S3_DOWNLOAD
  with:
    access_key: ${{ secrets.AWS_ACCESS_KEY }}
    secret_key: ${{ secrets.AWS_SECRET_KEY }}
    region: us-east-1
    bucket: myapp-backups
    source_key: backups/latest.tar.gz
    destination_path: /tmp/restore.tar.gz
```

---

## 10. 输出管理

### 10.1 步骤输出定义

**说明**：步骤可以定义输出变量，供后续步骤使用。

**定义方式**：

```yaml
outputs:
  <output_name>:
    value: string               # 输出值（表达式）
    description: string         # 输出描述（可选）
```

**输出来源**：

1. **组件执行结果**：通过 `steps.<step_id>.result.<field>` 访问
2. **环境变量输出**：在Shell命令中使用 `echo "KEY=value" >> $GITHUB_OUTPUT`

**示例**：

```yaml
steps:
  - id: build
    type: SHELL
    with:
      command: |
        BUILD_ID=$(date +%Y%m%d%H%M%S)
        echo "BUILD_ID=${BUILD_ID}" >> $GITHUB_OUTPUT
    outputs:
      build_id:
        value: ${{ steps.build.result.BUILD_ID }}
        description: "构建ID"

  - id: deploy
    type: COMPOSE
    with:
      app_name: myapp
      image_tag: ${{ steps.build.outputs.build_id }}
```

### 10.2 输出引用

**引用语法**：

```yaml
${{ steps.<step_id>.outputs.<output_name> }}
```

### 10.3 输出大小限制

- 单个输出值：最大 1MB
- 单个步骤总输出：最大 10MB
- 工作流总输出：最大 100MB

---

## 11. 自定义组件

### 11.1 组件引用语法

**说明**：自定义组件功能在 V1.0 版本中暂不支持，将在未来版本中实现。

**规划的引用语法**：

```yaml
steps:
  - id: custom_step
    name: "使用自定义组件"
    target: "server1"
    type: CUSTOM
    uses: ./components/my-component@v1.0.0
    with:
      param1: value1
      param2: value2
```

### 11.2 本地组件引用

**说明**：引用项目空间中的自定义组件（未来版本）。

**语法**：

```yaml
uses: ./components/component-name@version
```

### 11.3 远程组件引用

**说明**：引用远程仓库中的自定义组件（未来版本）。

**语法**：

```yaml
uses: github.com/org/repo/component@version
```

### 11.4 组件版本管理

**说明**：自定义组件支持语义化版本管理（未来版本）。

**版本指定方式**：

- 精确版本：`@1.0.0`
- 主版本：`@v1`
- 最新版本：`@latest`

**注意**：自定义组件功能将在后续版本中实现，当前版本请使用内置组件。

---

## 12. 错误处理

### 12.1 步骤级错误处理

**重试机制**：

```yaml
- id: api_call
  type: HTTP
  retry:
    max_attempts: 3
    initial_interval: 1
    backoff_coefficient: 2.0
    maximum_interval: 60
  with:
    url: https://api.example.com/deploy
```

**条件执行**：

```yaml
- id: cleanup_on_failure
  type: SHELL
  if: ${{ failure() }}
  with:
    command: ./cleanup.sh
```

### 12.2 工作流级错误处理

**失败通知**：

```yaml
steps:
  - id: main_task
    type: SHELL
    with:
      command: ./main-task.sh

  - id: notify_failure
    type: EMAIL
    if: ${{ failure() }}
    with:
      subject: "工作流执行失败"
      body: "错误信息: ${{ execution.error_message }}"
```

**清理步骤**：

```yaml
- id: cleanup
  type: SHELL
  if: ${{ always() }}
  with:
    command: |
      rm -rf /tmp/build
      docker system prune -f
```

---

## 13. 安全规范

### 13.1 密钥引用规范

#### 13.1.1 基本引用语法

**使用 secrets 上下文**：

```yaml
steps:
  - id: clone_repo
    type: SHELL
    env:
      GIT_TOKEN: ${{ secrets.GIT_TOKEN }}
    with:
      command: git clone https://$GIT_TOKEN@github.com/example/repo.git
```

**禁止明文存储**：

```yaml
# ❌ 错误：明文密码
env:
  DB_PASSWORD: "mypassword123"

# ✅ 正确：使用secrets
env:
  DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
```

#### 13.1.2 密钥类型和字段访问

**支持的密钥类型**：

| 密钥类型            | 字段定义                                   | 使用场景                  |
| ------------------- | ------------------------------------------ | ------------------------- |
| `api_key`           | `api_key`: string                          | API调用认证               |
| `username_password` | `username`: string `password`: string      | 数据库、SSH、SMTP等       |
| `access_key_secret` | `access_key`: string `secret_key`: string  | 云服务认证（AWS、阿里云） |
| `ssh_key`           | `private_key`: string `passphrase`: string | SSH连接                   |
| `oauth_token`       | `token`: string `refresh_token`: string    | OAuth认证                 |
| `certificate`       | `cert`: string `key`: string               | TLS/SSL证书               |

**多字段密钥访问**：

```yaml
steps:
  - id: backup_database
    type: SHELL
    env:
      DB_USER: ${{ secrets.DB_CREDENTIAL.username }}
      DB_PASSWORD: ${{ secrets.DB_CREDENTIAL.password }}
      AWS_ACCESS_KEY_ID: ${{ secrets.AWS_CREDENTIAL.access_key }}
      AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_CREDENTIAL.secret_key }}
    with:
      command: |
        mysqldump -u $DB_USER -p$DB_PASSWORD mydb > backup.sql
        aws s3 cp backup.sql s3://my-bucket/
```

#### 13.1.3 密钥引用位置

**1. 环境变量中引用**：

```yaml
steps:
  - id: deploy
    type: SHELL
    env:
      API_TOKEN: ${{ secrets.API_TOKEN }}
    with:
      command: ./deploy.sh
```

**2. 组件参数中引用**：

```yaml
steps:
  - id: api_call
    type: HTTP
    with:
      method: POST
      url: https://api.example.com/deploy
      headers:
        Authorization: Bearer ${{ secrets.API_TOKEN }}
```

**3. 命令字符串中引用**：

```yaml
steps:
  - id: ssh_deploy
    type: SHELL
    env:
      SSH_KEY: ${{ secrets.SSH_KEY.private_key }}
    with:
      command: |
        ssh -i $SSH_KEY user@server 'bash -s' << 'EOF'
          cd /opt/app && git pull
        EOF
```

### 13.2 敏感数据处理

#### 13.2.1 日志脱敏

**自动脱敏**：

- 凭据值自动在日志中脱敏
- 显示为 `***`
- 错误信息中不包含密钥原文

**示例**：

```yaml
steps:
  - id: connect_db
    type: SHELL
    env:
      DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
    with:
      command: |
        echo "Connecting to database..."
        mysql -p$DB_PASSWORD -e "SELECT 1"
```

日志输出：

```
Connecting to database...
mysql -p*** -e "SELECT 1"
```

#### 13.2.2 环境变量保护

**通过环境变量传递密钥**：

```yaml
steps:
  - id: deploy
    type: SHELL
    env:
      API_KEY: ${{ secrets.API_KEY }}
      DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
    with:
      command: |
        # 凭据通过环境变量传递，不在命令中显示
        ./deploy.sh
```

**禁止在命令中直接使用密钥**：

```yaml
# ❌ 错误：密钥直接暴露在命令中
steps:
  - id: bad_example
    type: SHELL
    with:
      command: curl -H "Authorization: Bearer ${{ secrets.API_TOKEN }}" https://api.example.com

# ✅ 正确：通过环境变量传递
steps:
  - id: good_example
    type: SHELL
    env:
      API_TOKEN: ${{ secrets.API_TOKEN }}
    with:
      command: curl -H "Authorization: Bearer $API_TOKEN" https://api.example.com
```

#### 13.2.3 密钥安全规范

**传输安全**：

- 所有密钥通过 TLS 1.3 加密传输
- 支持信封加密（二次加密）
- 使用临时会话密钥

**存储安全**：

- 密钥仅在 Worker 进程内存中存在
- 严禁写入磁盘、日志或临时文件
- 执行完成后立即清除内存

**使用安全**：

- 每次执行从凭据管理系统获取最新密钥
- 密钥仅在当前执行实例中有效
- 不缓存密钥数据
- 一次性使用原则

**管理安全**：

- 工作流定义中仅保存密钥引用（名称）
- 不保存密钥原文
- 基于 RBAC 的访问控制
- 支持密钥轮换和过期管理

### 13.3 权限最小化原则

#### 13.3.1 文件权限

**设置合适的文件权限**：

```yaml
steps:
  - id: create_config
    type: FILE_UPLOAD
    with:
      source_path: /tmp/config.yaml
      destination_path: /etc/app/config.yaml
      permissions: "0600"  # 仅所有者可读写
```

**权限建议**：

| 文件类型 | 推荐权限 | 说明                 |
| -------- | -------- | -------------------- |
| 配置文件 | 0600     | 仅所有者可读写       |
| 脚本文件 | 0700     | 仅所有者可读写执行   |
| 数据文件 | 0644     | 所有者可写，其他只读 |
| 密钥文件 | 0400     | 仅所有者只读         |

#### 13.3.2 命令执行限制

**禁止的危险命令**：

- `rm -rf /`：删除根目录
- `dd if=/dev/zero of=/dev/sda`：覆盖磁盘
- `:(){ :|:& };:`：Fork炸弹
- `chmod -R 777 /`：修改根目录权限

**限制访问的敏感目录**：

- `/etc`：系统配置目录
- `/sys`：系统内核接口
- `/proc`：进程信息
- `/boot`：启动文件

**使用受限用户**：

```yaml
steps:
  - id: safe_execution
    type: SHELL
    with:
      command: |
        # 使用非root用户执行
        su - appuser -c './app-script.sh'
```

#### 13.3.3 网络访问控制

**限制外部访问**：

```yaml
steps:
  - id: api_call
    type: HTTP
    with:
      method: POST
      url: https://internal-api.example.com/endpoint  # 仅访问内部API
      headers:
        Authorization: Bearer ${{ secrets.API_TOKEN }}
```

**使用白名单**：

- 仅允许访问预定义的域名和IP
- 禁止访问内网敏感服务
- 限制端口访问范围

---

## 14. 触发器配置

### 14.1 触发器语法规划

**字段**：

```yaml
name: "应用自动部署"
version: "1.0.0"
description: "自动化应用部署工作流"

# 触发器配置
triggers:
  - type: schedule      # 定时触发
    name: "每日凌晨部署"
    cron: "0 2 * * *"
    timezone: "Asia/Shanghai"
    enabled: true

  - type: webhook       # Webhook触发
    name: "代码推送触发"
    filters:
      - field: "event"
        operator: "equals"
        value: "push"
    variable_mapping:
      COMMIT_HASH: "$.head_commit.id"
    enabled: true

variables:
  - name: APP_NAME
    type: string
    value: "myapp"

steps:
  - id: deploy
    name: "部署应用"
    type: COMPOSE
    with:
      app_name: ${{ variables.APP_NAME }}
```

### 14.2 定时触发配置

**Cron 表达式格式**（5 字段）：

```
┌───────────── 分钟 (0 - 59)
│ ┌───────────── 小时 (0 - 23)
│ │ ┌───────────── 日 (1 - 31)
│ │ │ ┌───────────── 月 (1 - 12)
│ │ │ │ ┌───────────── 星期 (0 - 7)
│ │ │ │ │
* * * * *
```

**配置示例**：

```yaml
triggers:
  - type: schedule
    name: "每日备份"
    cron: "0 2 * * *"           # 每天凌晨2点
    timezone: "Asia/Shanghai"   # 时区
    enabled: true               # 是否启用
    missed_policy: "skip"       # 错过执行策略
```

**常用 Cron 表达式**：

| 表达式         | 说明             |
| -------------- | ---------------- |
| `0 2 * * *`    | 每天凌晨2点执行  |
| `*/30 * * * *` | 每30分钟执行一次 |
| `0 */2 * * *`  | 每2小时执行一次  |
| `0 0 * * 0`    | 每周日凌晨执行   |
| `0 0 1 * *`    | 每月1号凌晨执行  |

### 14.3 Webhook 触发配置

**基本配置**：

```yaml
triggers:
  - type: webhook
    name: "GitHub推送触发"
    enabled: true
    filters:
      - field: "event"
        operator: "equals"
        value: "push"
      - field: "branch"
        operator: "equals"
        value: "main"
    variable_mapping:
      REPO_NAME: "$.repository.name"
      COMMIT_HASH: "$.head_commit.id"
      AUTHOR: "$.head_commit.author.email"
```

**事件过滤操作符**：

| 操作符        | 说明      | 示例                                   |
| ------------- | --------- | -------------------------------------- |
| `equals`      | 等于      | `"branch" equals "main"`               |
| `not_equals`  | 不等于    | `"event" not_equals "delete"`          |
| `contains`    | 包含      | `"message" contains "hotfix"`          |
| `starts_with` | 以...开头 | `"branch" starts_with "release/"`      |
| `ends_with`   | 以...结尾 | `"file" ends_with ".yaml"`             |
| `in`          | 在列表中  | `"branch" in ["main", "develop"]`      |
| `regex`       | 正则匹配  | `"branch" regex "^release/v\d+\.\d+$"` |

**变量映射（JSONPath）**：

```yaml
variable_mapping:
  APP_NAME: "$.repository"          # 根级字段
  BRANCH: "$.ref"                   # 根级字段
  COMMIT_ID: "$.commits[0].id"     # 数组第一个元素
  ALL_TAGS: "$.tags[*].name"       # 数组所有元素
```

### 14.4 手动触发配置

手动触发不需要在 YAML 中配置，通过 API 或 Web 界面触发即可。

**API 触发示例**：

```bash
POST /api/v1/workflows/{id}/execute

{
  "input_variables": {
    "APP_NAME": "test-app",
    "ENVIRONMENT": "production"
  },
  "async": true
}
```

## 15. GitHub Action 组件

### 15.1 组件类型

**GITHUB_ACTION 组件类型**：

用于在 Websoft9 工作流中运行 GitHub Actions 组件。

**支持的 Action 类型**：

| Action 类型      | 说明                 | 示例                   |
| ---------------- | -------------------- | ---------------------- |
| NodeJS Action    | 使用 Node.js 运行    | `actions/checkout@v4`  |
| Docker Action    | 使用 Docker 容器运行 | `docker://alpine:3.18` |
| Composite Action | 组合多个步骤         | 自定义组合 Action      |

### 15.2 引用语法

**基本语法**：

```yaml
steps:
  - id: step_id
    name: "步骤名称"
    type: GITHUB_ACTION
    with:
      uses: owner/repo@version
      with:
        input_param: value
```

**完整示例**：

```yaml
steps:
  - id: checkout_code
    name: "检出代码"
    type: GITHUB_ACTION
    with:
      uses: actions/checkout@v4
      with:
        repository: example/myapp
        ref: main
        token: ${{ secrets.GITHUB_TOKEN }}
```

### 15.3 常用 Action 示例

**代码管理**：

```yaml
# 检出代码
- id: checkout
  type: GITHUB_ACTION
  with:
    uses: actions/checkout@v4
    with:
      repository: example/myapp
      ref: main

# 缓存依赖
- id: cache
  type: GITHUB_ACTION
  with:
    uses: actions/cache@v4
    with:
      path: ~/.npm
      key: npm-${{ hashFiles('package-lock.json') }}
```

**构建工具**：

```yaml
# 设置 Node.js
- id: setup_node
  type: GITHUB_ACTION
  with:
    uses: actions/setup-node@v4
    with:
      node-version: '20'
      cache: 'npm'

# 设置 Python
- id: setup_python
  type: GITHUB_ACTION
  with:
    uses: actions/setup-python@v5
    with:
      python-version: '3.11'
```

**制品文件管理**：

```yaml
# 上传构建制品文件
- id: upload_artifact
  type: GITHUB_ACTION
  with:
    uses: actions/upload-artifact@v4
    with:
      name: build-output
      path: dist/

# 下载构建制品文件
- id: download_artifact
  type: GITHUB_ACTION
  with:
    uses: actions/download-artifact@v4
    with:
      name: build-output
      path: dist/
```

### 15.4 与自定义组件混合使用

**完整工作流示例**：

```yaml
name: "应用构建和部署"
version: "1.0.0"
description: "使用 GitHub Actions 构建，使用自定义组件部署"

variables:
  - name: APP_NAME
    type: string
    value: "myapp"

steps:
  # 1. 使用 GitHub Action 检出代码
  - id: checkout
    name: "检出代码"
    type: GITHUB_ACTION
    with:
      uses: actions/checkout@v4
      with:
        repository: example/myapp
        ref: main
        token: ${{ secrets.GITHUB_TOKEN }}

  # 2. 使用 GitHub Action 设置 Node.js
  - id: setup_node
    name: "设置 Node.js"
    type: GITHUB_ACTION
    with:
      uses: actions/setup-node@v4
      with:
        node-version: '20'
        cache: 'npm'

  # 3. 使用自定义 SHELL 组件构建
  - id: build
    name: "构建应用"
    type: SHELL
    with:
      command: |
        npm install
        npm run build
      server_id: 1

  # 4. 使用自定义 COMPOSE 组件部署
  - id: deploy
    name: "部署应用"
    type: COMPOSE
    with:
      app_name: ${{ variables.APP_NAME }}
      server_id: 1
```

### 15.5 输入输出处理

**输入参数映射**：

```yaml
# Action 定义的输入参数
inputs:
  who-to-greet:
    description: 'Who to greet'
    required: true
    default: 'World'

# Websoft9 工作流中传递参数
steps:
  - id: greet
    type: GITHUB_ACTION
    with:
      uses: actions/hello-world-action@v1
      with:
        who-to-greet: ${{ variables.USER_NAME }}
```

**输出变量引用**：

```yaml
steps:
  - id: action_step
    type: GITHUB_ACTION
    with:
      uses: actions/some-action@v1
    outputs:
      result_data:
        value: ${{ steps.action_step.outputs.data }}

  # 在后续步骤中使用输出
  - id: use_output
    type: SHELL
    with:
      command: |
        echo "Action result: ${{ steps.action_step.outputs.result_data }}"
```

### 15.6 环境变量共享

**工作流级别环境变量**：

```yaml
env:
  NODE_ENV: production
  API_URL: https://api.example.com

steps:
  # GitHub Action可以访问
  - id: action_step
    type: GITHUB_ACTION
    with:
      uses: actions/some-action@v1

  # 自定义组件也可以访问
  - id: shell_step
    type: SHELL
    with:
      command: |
        echo "NODE_ENV: $NODE_ENV"
```

### 15.7 限制和约束

**支持的特性**：

- ✅ 基本输入输出
- ✅ 环境变量
- ✅ 工作目录共享
- ✅ 密钥引用
- ✅ 条件执行
- ✅ 重试机制

**不支持的特性**：

- ❌ GitHub 特定上下文
- ❌ GitHub API
- ❌ Artifact 持久化存（仅支持临时存储）
- ❌ Matrix 策略
- ❌ Reusable workflows

**性能考虑**：

- NodeJS Action 启动快（< 1秒）
- Docker Action 启动慢（5-30秒）
- 建议优先使用 NodeJS Action
- 大型 Docker 镜像建议预拉取

---

## 16. 语法限制与约束

### 16.1 文件大小限制

- YAML文件大小：最大 1MB
- 超过限制时建议拆分为多个工作流

### 16.2 嵌套深度限制

- YAML嵌套深度：最大 10 层
- 表达式嵌套深度：最大 5 层

### 16.3 字符串长度限制

| 字段       | 最大长度   |
| ---------- | ---------- |
| 工作流名称 | 64 字符    |
| 步骤名称   | 64 字符    |
| 变量名     | 64 字符    |
| 描述       | 500 字符   |
| 命令/脚本  | 10000 字符 |

### 16.4 数量限制

| 项目                   | 最大数量 |
| ---------------------- | -------- |
| 步骤数量               | 100      |
| 变量数量               | 50       |
| 环境变量数量           | 100      |
| 输出变量数量（每步骤） | 20       |

### 16.5 执行时间限制

- 单个步骤超时：默认 300 秒，最大 86400 秒（24小时）
- 工作流总超时：默认 1 小时，最大 24 小时

### 16.6 输出大小限制

- 单个输出值：最大 1MB
- 单个步骤总输出：最大 10MB
- 工作流总输出：最大 100MB

---

## 17. 完整示例

### 17.1 应用自动部署工作流

```yaml
name: "应用自动部署"
version: "1.0.0"
description: "从Git仓库拉取代码，构建Docker镜像，部署到生产环境"
author: "DevOps Team"
tags:
  - deployment
  - docker
  - automation
category: deployment

variables:
  - name: REPO_URL
    type: string
    value: "https://github.com/example/myapp.git"
    required: true
    description: "Git仓库地址"

  - name: BRANCH
    type: string
    value: "main"
    required: true
    description: "Git分支"

  - name: APP_NAME
    type: string
    value: "myapp"
    required: true
    description: "应用名称"
    validation:
      pattern: "^[a-z0-9-]+$"
      min: 1
      max: 32

  - name: SERVER_ID
    type: number
    value: 1
    required: true
    description: "目标服务器ID"

env:
  DOCKER_REGISTRY: "registry.example.com"
  NOTIFICATION_EMAIL: "ops@example.com"

steps:
  - id: clone_repo
    name: "拉取代码"
    target: "server1"
    type: SHELL
    timeout: 300
    retry:
      max_attempts: 3
      initial_interval: 5
      backoff_coefficient: 2.0
    with:
      command: |
        #!/bin/bash
        set -e
        rm -rf /tmp/build/${{ variables.APP_NAME }}
        git clone -b ${{ variables.BRANCH }} ${{ variables.REPO_URL }} /tmp/build/${{ variables.APP_NAME }}
        cd /tmp/build/${{ variables.APP_NAME }}
        echo "COMMIT_HASH=$(git rev-parse HEAD)" >> $GITHUB_OUTPUT
      working_directory: /tmp
      server_id: ${{ variables.SERVER_ID }}
    env:
      GIT_TOKEN: ${{ secrets.GIT_TOKEN }}
    outputs:
      commit_hash:
        value: ${{ steps.clone_repo.result.COMMIT_HASH }}
        description: "Git提交哈希值"

  - id: build_image
    name: "构建Docker镜像"
    target: "server1"
    type: SHELL
    timeout: 600
    with:
      command: |
        #!/bin/bash
        set -e
        cd /tmp/build/${{ variables.APP_NAME }}
        docker build -t ${{ env.DOCKER_REGISTRY }}/${{ variables.APP_NAME }}:${{ steps.clone_repo.outputs.commit_hash }} .
        docker push ${{ env.DOCKER_REGISTRY }}/${{ variables.APP_NAME }}:${{ steps.clone_repo.outputs.commit_hash }}
      server_id: ${{ variables.SERVER_ID }}
    env:
      DOCKER_USERNAME: ${{ secrets.DOCKER_USERNAME }}
      DOCKER_PASSWORD: ${{ secrets.DOCKER_PASSWORD }}

  - id: deploy_app
    name: "部署应用"
    target: "server1"
    type: COMPOSE
    timeout: 300
    with:
      app_name: ${{ variables.APP_NAME }}
      server_id: ${{ variables.SERVER_ID }}
      pull_images: true
      recreate: true
    env:
      IMAGE_TAG: ${{ steps.clone_repo.outputs.commit_hash }}

  - id: health_check
    name: "健康检查"
    target: "server1"
    type: HTTP
    timeout: 60
    retry:
      max_attempts: 5
      initial_interval: 10
      backoff_coefficient: 1.5
    with:
      method: GET
      url: https://${{ variables.APP_NAME }}.example.com/health
      expected_status: 200

  - id: notify_success
    name: "发送成功通知"
    target: "server1"
    type: EMAIL
    if: ${{ success() }}
    with:
      smtp_server: smtp.example.com
      smtp_port: 587
      smtp_user: ${{ secrets.SMTP_USER }}
      smtp_password: ${{ secrets.SMTP_PASSWORD }}
      from: noreply@example.com
      to:
        - ${{ env.NOTIFICATION_EMAIL }}
      subject: "✅ 部署成功: ${{ variables.APP_NAME }}"
      body: |
        应用 ${{ variables.APP_NAME }} 已成功部署。

        提交哈希: ${{ steps.clone_repo.outputs.commit_hash }}
        部署时间: ${{ execution.start_time }}
        执行者: ${{ execution.trigger_by }}

  - id: notify_failure
    name: "发送失败通知"
    target: "server1"
    type: EMAIL
    if: ${{ failure() }}
    with:
      smtp_server: smtp.example.com
      smtp_port: 587
      smtp_user: ${{ secrets.SMTP_USER }}
      smtp_password: ${{ secrets.SMTP_PASSWORD }}
      from: noreply@example.com
      to:
        - ${{ env.NOTIFICATION_EMAIL }}
      subject: "❌ 部署失败: ${{ variables.APP_NAME }}"
      body: |
        应用 ${{ variables.APP_NAME }} 部署失败。

        错误信息: ${{ execution.error_message }}
        部署时间: ${{ execution.start_time }}

  - id: cleanup
    name: "清理临时文件"
    target: "server1"
    type: SHELL
    if: ${{ always() }}
    with:
      command: |
        rm -rf /tmp/build/${{ variables.APP_NAME }}
      server_id: ${{ variables.SERVER_ID }}
```

### 17.2 数据库定时备份工作流

```yaml
name: "数据库定时备份"
version: "1.0.0"
description: "每天凌晨2点自动备份数据库并上传到S3"
category: backup

variables:
  - name: DB_HOST
    type: string
    value: "localhost"
    required: true

  - name: DB_NAME
    type: string
    value: "myapp"
    required: true

  - name: BACKUP_RETENTION_DAYS
    type: number
    value: 30
    required: true
    description: "备份保留天数"

steps:
  - id: backup_database
    name: "备份数据库"
    target: "server1"
    type: SHELL
    timeout: 1800
    with:
      command: |
        #!/bin/bash
        set -e
        BACKUP_FILE="/tmp/backup-${{ variables.DB_NAME }}-$(date +%Y%m%d%H%M%S).sql.gz"
        mysqldump -h ${{ variables.DB_HOST }} \
                  -u ${{ secrets.DB_USER }} \
                  -p${{ secrets.DB_PASSWORD }} \
                  ${{ variables.DB_NAME }} | gzip > $BACKUP_FILE
        echo "BACKUP_FILE=$BACKUP_FILE" >> $GITHUB_OUTPUT
        echo "BACKUP_SIZE=$(du -h $BACKUP_FILE | cut -f1)" >> $GITHUB_OUTPUT
      server_id: 1
    outputs:
      backup_file:
        value: ${{ steps.backup_database.result.BACKUP_FILE }}
      backup_size:
        value: ${{ steps.backup_database.result.BACKUP_SIZE }}

  - id: upload_to_s3
    name: "上传到S3"
    target: "server1"
    type: S3_UPLOAD
    with:
      access_key: ${{ secrets.AWS_ACCESS_KEY }}
      secret_key: ${{ secrets.AWS_SECRET_KEY }}
      region: us-east-1
      bucket: database-backups
      source_path: ${{ steps.backup_database.outputs.backup_file }}
      destination_key: ${{ variables.DB_NAME }}/backup-$(date +%Y%m%d%H%M%S).sql.gz

  - id: cleanup_old_backups
    name: "清理过期备份"
    target: "server1"
    type: SHELL
    with:
      command: |
        find /tmp -name "backup-${{ variables.DB_NAME }}-*.sql.gz" -mtime +${{ variables.BACKUP_RETENTION_DAYS }} -delete
      server_id: 1

  - id: notify
    name: "发送通知"
    target: "server1"
    type: EMAIL
    if: ${{ always() }}
    with:
      smtp_server: smtp.example.com
      smtp_port: 587
      smtp_user: ${{ secrets.SMTP_USER }}
      smtp_password: ${{ secrets.SMTP_PASSWORD }}
      from: noreply@example.com
      to:
        - dba@example.com
      subject: "数据库备份 - ${{ variables.DB_NAME }}"
      body: |
        数据库 ${{ variables.DB_NAME }} 备份完成。

        备份文件: ${{ steps.backup_database.outputs.backup_file }}
        备份大小: ${{ steps.backup_database.outputs.backup_size }}
        备份时间: ${{ execution.start_time }}
        执行状态: ${{ execution.status }}
```

---

## 18. 版本兼容性

### 18.1 DSL 版本声明

当前 DSL 版本：**V1.0**

工作流定义中的 `version` 字段表示工作流自身的版本，不是 DSL 版本。

### 18.2 版本迁移指南

**未来版本升级时的兼容性策略**：

- **向后兼容**：新版本 DSL 将保持对旧版本的兼容
- **废弃通知**：废弃的语法将提前至少一个大版本通知
- **迁移工具**：提供自动化迁移工具辅助升级

### 18.3 废弃语法说明

当前版本（V1.0）无废弃语法。

---

## 19. 附录

### 19.1 YAML Schema 定义

完整的 JSON Schema 定义请参考：`schemas/workflow-dsl-v1.0.json`

### 19.2 语法速查表

**工作流结构**：

```yaml
name: string
version: string
description: string
variables: array
env: object
steps: array
```

**步骤结构**：

```yaml
- id: string
  target: string
  name: string
  type: string
  if: string
  timeout: integer
  retry: object
  with: object
  env: object
  outputs: object
```

**表达式语法**：

```yaml
${{ variables.NAME }}
${{ env.KEY }}
${{ secrets.TOKEN }}
${{ steps.id.outputs.name }}
${{ success() }}
${{ failure() }}
${{ always() }}
```

### 19.3 保留关键字列表

以下关键字为系统保留，不可用作变量名或步骤ID：

- `workflow`
- `execution`
- `project`
- `server`
- `system`
- `internal`

### 19.4 错误代码参考

| 错误代码                 | 说明               | 解决方案                     |
| ------------------------ | ------------------ | ---------------------------- |
| `YAML_PARSE_ERROR`       | YAML 语法错误      | 检查 YAML 格式，确保缩进正确 |
| `INVALID_VARIABLE_NAME`  | 变量名不符合规范   | 使用大写字母、数字和下划线   |
| `INVALID_STEP_ID`        | 步骤ID不符合规范   | 使用小写字母、数字和下划线   |
| `DUPLICATE_STEP_ID`      | 步骤ID重复         | 确保每个步骤ID唯一           |
| `UNKNOWN_COMPONENT_TYPE` | 未知的组件类型     | 检查组件类型是否正确         |
| `INVALID_EXPRESSION`     | 表达式语法错误     | 检查表达式语法               |
| `CIRCULAR_DEPENDENCY`    | 步骤间存在循环依赖 | 检查步骤依赖关系             |

### 19.5 与 GitHub Actions 语法对照表

| 功能       | GitHub Actions | Websoft9 Workflow   |
| ---------- | -------------- | ------------------- |
| 工作流名称 | `name`         | `name`              |
| 触发器     | `on`           | 通过任务调度配置    |
| 作业       | `jobs`         | `steps`（顺序执行） |
| 步骤       | `steps`        | `steps`             |
| 运行环境   | `runs-on`      | `target`            |
| 环境变量   | `env`          | `env`               |
| 密钥       | `secrets`      | `secrets`           |
| 条件执行   | `if`           | `if`                |
| 输出       | `outputs`      | `outputs`           |
| 表达式     | `${{ }}`       | `${{ }}`            |
| 组件引用   | `uses`         | `type`              |

---

## 参考文档

- [工作流功能设计V1.4](./工作流功能设计V1.3.md)
- [Temporal 官方文档](https://docs.temporal.io)
- [GitHub Actions 语法](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)
- [YAML 1.2 规范](https://yaml.org/spec/1.2/spec.html)
- [Semantic Versioning 2.0.0](https://semver.org/)

---

## 变更记录

| 版本 | 日期       | 变更人        | 变更内容                                                                                                                                                                                                                                                                                       |
| ---- | ---------- | ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| V1.0 | 2025-12-01 | Websoft9 team | 初始版本，定义完整的 YAML DSL 语法规范                                                                                                                                                                                                                                                         |
| V1.1 | 2025-12-01 | Websoft9 team | 补充完善以下内容：1. 扩展第13章"安全规范"，增加密钥类型、字段访问、引用位置、安全规范等详细说明2. 新增第13A章"触发器配置"（未来版本），说明定时触发、Webhook触发的YAML语法规划3. 新增第13B章"GitHub Action 组件"，说明GITHUB_ACTION组件类型的引用语法和使用方式4. 与工作流功能设计V1.4保持一致 |

---

**文档结束**
