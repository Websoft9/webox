-- Websoft9 MySQL Database Initialization Script V1.1
-- Generated from Database Design Document V1.1
-- Compatible with MySQL 8.0+

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
SET sql_mode = 'STRICT_TRANS_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO';

-- ========================================
-- 3.1 平台主页 (Platform Home)
-- ========================================

-- 应用快捷导航表
CREATE TABLE IF NOT EXISTS `app_shortcuts` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `app_instance_id` BIGINT UNSIGNED NOT NULL COMMENT '应用实例ID',
    `name` VARCHAR(64) NULL COMMENT '自定义名称',
    `description` VARCHAR(255) NULL COMMENT '自定义描述',
    `icon` VARCHAR(255) NULL COMMENT '自定义图标',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序顺序',
    `access_count` INT NOT NULL DEFAULT 0 COMMENT '访问次数',
    `last_accessed` DATETIME NULL COMMENT '最后访问时间',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_app_instance_id` (`app_instance_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用快捷导航表';

-- ========================================
-- 3.2 项目管理 (Project Management)
-- ========================================

-- 项目表
CREATE TABLE IF NOT EXISTS `projects` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(100) NOT NULL COMMENT '项目名称',
    `identifier` VARCHAR(50) NOT NULL UNIQUE COMMENT '项目标识符',
    `description` TEXT NULL COMMENT '项目描述',
    `tags` JSON NULL COMMENT '项目标签',
    `icon` VARCHAR(255) NULL COMMENT '项目图标URL',
    `status` ENUM('NORMAL', 'ARCHIVED', 'DELETED') NOT NULL DEFAULT 'NORMAL' COMMENT '项目状态',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '项目所有者ID',
    `default_resource_group` VARCHAR(50) NOT NULL DEFAULT 'default' COMMENT '默认资源组',
    `default_timezone` VARCHAR(50) NOT NULL DEFAULT 'Asia/Shanghai' COMMENT '默认时区',
    `log_retention_days` INT NOT NULL DEFAULT 30 COMMENT '日志保留天数',
    `backup_strategy` ENUM('daily', 'weekly', 'monthly', 'disabled') NOT NULL DEFAULT 'daily' COMMENT '备份策略',
    `access_control` VARCHAR(50) NOT NULL DEFAULT 'members' COMMENT '访问控制',
    `api_access` BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'API访问权限',
    `audit_enabled` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '审计日志启用',
    `last_activity_at` DATETIME NULL COMMENT '最后活跃时间',
    `archived_at` DATETIME NULL COMMENT '归档时间',
    `archived_by` BIGINT UNSIGNED NULL COMMENT '归档操作人ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_identifier` (`identifier`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    KEY `idx_archived_by` (`archived_by`),
    CONSTRAINT `fk_projects_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_projects_archived_by` FOREIGN KEY (`archived_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目表';

-- 项目成员表
CREATE TABLE IF NOT EXISTS `project_members` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT '项目ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `role` VARCHAR(50) NOT NULL COMMENT '角色：admin, developer, operator, viewer',
    `is_admin` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否项目管理员',
    `status` ENUM('active', 'inactive', 'pending') NOT NULL DEFAULT 'active' COMMENT '成员状态',
    `invited_by` BIGINT UNSIGNED NULL COMMENT '邀请人ID',
    `invited_at` DATETIME NULL COMMENT '邀请时间',
    `joined_at` DATETIME NULL COMMENT '加入时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_project_user` (`project_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_role` (`role`),
    KEY `idx_status` (`status`),
    KEY `idx_invited_by` (`invited_by`),
    CONSTRAINT `fk_project_members_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_project_members_invited_by` FOREIGN KEY (`invited_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目成员表';

-- 项目环境变量表
CREATE TABLE IF NOT EXISTS `project_environments` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT '项目ID',
    `name` VARCHAR(100) NOT NULL COMMENT '变量名称',
    `value` TEXT NULL COMMENT '变量值',
    `type` ENUM('normal', 'sensitive') NOT NULL DEFAULT 'normal' COMMENT '变量类型',
    `description` VARCHAR(500) NULL COMMENT '变量描述',
    `scope` ENUM('global', 'app', 'workflow') NOT NULL DEFAULT 'global' COMMENT '作用范围',
    `scope_target` VARCHAR(100) NULL COMMENT '作用目标',
    `is_encrypted` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否加密存储',
    `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_project_name_scope` (`project_id`, `name`, `scope`, `scope_target`),
    KEY `idx_created_by` (`created_by`),
    KEY `idx_type` (`type`),
    KEY `idx_scope` (`scope`),
    CONSTRAINT `fk_project_environments_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_project_environments_created_by` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目环境变量表';

-- 项目文件表
CREATE TABLE IF NOT EXISTS `project_files` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT '项目ID',
    `name` VARCHAR(255) NOT NULL COMMENT '文件名称',
    `path` VARCHAR(1024) NOT NULL COMMENT '文件路径',
    `type` ENUM('FILE', 'DIRECTORY') NOT NULL DEFAULT 'FILE' COMMENT '类型：FILE, DIRECTORY',
    `size` BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小(字节)',
    `mime_type` VARCHAR(128) NULL COMMENT 'MIME类型',
    `download_count` INT NOT NULL DEFAULT 0 COMMENT '下载次数',
    `parent_id` BIGINT UNSIGNED NULL COMMENT '父目录ID',
    `storage_path` VARCHAR(1024) NULL COMMENT '存储路径',
    `checksum` VARCHAR(64) NULL COMMENT '文件校验和',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_project_id` (`project_id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_type` (`type`),
    KEY `idx_path` (`path`(255)),
    CONSTRAINT `fk_project_files_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_project_files_parent` FOREIGN KEY (`parent_id`) REFERENCES `project_files` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_project_files_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目文件表';

-- 项目活动记录表
CREATE TABLE IF NOT EXISTS `project_activities` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT '项目ID',
    `user_id` BIGINT UNSIGNED NULL COMMENT '操作用户ID',
    `username` VARCHAR(64) NULL COMMENT '操作用户名',
    `action` VARCHAR(50) NOT NULL COMMENT '操作动作',
    `resource_type` VARCHAR(50) NULL COMMENT '资源类型',
    `resource_id` BIGINT UNSIGNED NULL COMMENT '资源ID',
    `resource_name` VARCHAR(255) NULL COMMENT '资源名称',
    `description` TEXT NULL COMMENT '操作描述',
    `metadata` JSON NULL COMMENT '操作元数据',
    `ip_address` VARCHAR(45) NULL COMMENT '操作IP地址',
    `user_agent` VARCHAR(500) NULL COMMENT '用户代理',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_project_id` (`project_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_action` (`action`),
    KEY `idx_resource_type` (`resource_type`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_project_activities_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_project_activities_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目活动记录表';

-- 工作流表
CREATE TABLE IF NOT EXISTS `workflows` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '工作流名称',
    `code` VARCHAR(32) NOT NULL COMMENT '工作流编码',
    `description` TEXT NULL COMMENT '工作流描述',
    `definition` JSON NOT NULL COMMENT '工作流定义',
    `status` ENUM('DRAFT', 'ACTIVE', 'INACTIVE', 'ARCHIVED') NOT NULL DEFAULT 'DRAFT' COMMENT '状态',
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT '项目ID',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_project_id` (`project_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_workflows_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_workflows_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作流表';

-- 工作流任务表
CREATE TABLE IF NOT EXISTS `workflow_tasks` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '任务名称',
    `workflow_id` BIGINT UNSIGNED NOT NULL COMMENT '工作流ID',
    `schedule_type` ENUM('MANUAL', 'SCHEDULE', 'TRIGGER') NOT NULL DEFAULT 'MANUAL' COMMENT '调度类型',
    `cron_expression` VARCHAR(100) NULL COMMENT 'Cron表达式',
    `status` ENUM('DEFAULT', 'ONLINE', 'OFFLINE') NOT NULL DEFAULT 'DEFAULT' COMMENT '任务状态',
    `next_run_at` DATETIME NULL COMMENT '下次运行时间',
    `last_run_at` DATETIME NULL COMMENT '上次运行时间',
    `run_count` INT NOT NULL DEFAULT 0 COMMENT '运行次数',
    `success_count` INT NOT NULL DEFAULT 0 COMMENT '成功次数',
    `failure_count` INT NOT NULL DEFAULT 0 COMMENT '失败次数',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_workflow_id` (`workflow_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_workflow_tasks_workflow` FOREIGN KEY (`workflow_id`) REFERENCES `workflows` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作流任务表';

-- 工作流执行历史表
CREATE TABLE IF NOT EXISTS `workflow_executions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `task_id` BIGINT UNSIGNED NOT NULL COMMENT '任务ID',
    `execution_id` VARCHAR(64) NOT NULL UNIQUE COMMENT '执行唯一标识',
    `status` ENUM('PENDING', 'RUNNING', 'SUCCESS', 'FAILURE', 'STOPPED') NOT NULL DEFAULT 'PENDING' COMMENT '执行状态',
    `trigger_type` ENUM('MANUAL', 'SCHEDULE', 'TRIGGER') NOT NULL DEFAULT 'MANUAL' COMMENT '触发类型',
    `trigger_by` BIGINT UNSIGNED NULL COMMENT '触发者ID',
    `start_time` DATETIME NULL COMMENT '开始时间',
    `end_time` DATETIME NULL COMMENT '结束时间',
    `duration` INT NOT NULL DEFAULT 0 COMMENT '执行时长(秒)',
    `error_message` TEXT NULL COMMENT '错误信息',
    `execution_log` TEXT NULL COMMENT '执行日志',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_task_id` (`task_id`),
    KEY `idx_status` (`status`),
    KEY `idx_trigger_by` (`trigger_by`),
    CONSTRAINT `fk_workflow_executions_task` FOREIGN KEY (`task_id`) REFERENCES `workflow_tasks` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作流执行历史表';

-- 资源组表
CREATE TABLE IF NOT EXISTS `resource_groups` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT '项目ID',
    `name` VARCHAR(64) NOT NULL COMMENT '资源组名称',
    `code` VARCHAR(32) NOT NULL COMMENT '资源组编码',
    `description` TEXT NULL COMMENT '资源组描述',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_project_id` (`project_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_resource_groups_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_resource_groups_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='资源组表';

-- 数据库连接表
CREATE TABLE IF NOT EXISTS `database_connections` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '连接名称',
    `db_type` ENUM('mysql', 'postgresql', 'redis', 'mongodb') NOT NULL COMMENT '数据库类型',
    `host` VARCHAR(255) NOT NULL COMMENT '主机地址',
    `port` INT NOT NULL COMMENT '端口号',
    `database` VARCHAR(64) NULL COMMENT '数据库名称',
    `username` VARCHAR(64) NULL COMMENT '用户名',
    `password` VARCHAR(255) NULL COMMENT '加密密码',
    `ssl_enabled` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否启用SSL',
    `connection_timeout` INT NOT NULL DEFAULT 30 COMMENT '连接超时时间(秒)',
    `max_connections` INT NOT NULL DEFAULT 10 COMMENT '最大连接数',
    `status` ENUM('CONNECTED', 'DISCONNECTED', 'ERROR') NOT NULL DEFAULT 'CONNECTED' COMMENT '连接状态',
    `version` VARCHAR(32) NULL COMMENT '数据库版本',
    `charset` VARCHAR(32) NULL COMMENT '字符集',
    `description` TEXT NULL COMMENT '描述信息',
    `last_connected_at` DATETIME NULL COMMENT '最后连接时间',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `resource_group_id` BIGINT UNSIGNED NULL COMMENT '资源组ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_resource_group_id` (`resource_group_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_db_connections_resource_group` FOREIGN KEY (`resource_group_id`) REFERENCES `resource_groups` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据库连接表';

-- 服务器表
CREATE TABLE IF NOT EXISTS `servers` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '服务器名称',
    `hostname` VARCHAR(255) NOT NULL COMMENT '主机名',
    `ip_address` VARCHAR(45) NOT NULL COMMENT 'IP地址',
    `internal_ip` VARCHAR(45) NULL COMMENT '内网IP',
    `ssh_port` INT NOT NULL DEFAULT 22 COMMENT 'SSH端口',
    `os_type` VARCHAR(32) NOT NULL COMMENT '操作系统类型',
    `os_version` VARCHAR(64) NULL COMMENT '操作系统版本',
    `kernel_version` VARCHAR(64) NULL COMMENT '内核版本',
    `cpu_cores` INT NOT NULL DEFAULT 0 COMMENT 'CPU核心数',
    `memory_total` BIGINT NOT NULL DEFAULT 0 COMMENT '总内存(MB)',
    `disk_total` BIGINT NOT NULL DEFAULT 0 COMMENT '总磁盘空间(MB)',
    `architecture` VARCHAR(16) NULL COMMENT '系统架构',
    `status` ENUM('UNKNOWN', 'RUNNING', 'STOPPED') NOT NULL DEFAULT 'UNKNOWN' COMMENT '服务器状态',
    `last_heartbeat_at` DATETIME NULL COMMENT '最后心跳时间',
    `resource_group_id` BIGINT UNSIGNED NULL COMMENT '资源组ID',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `description` TEXT NULL COMMENT '服务器描述',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_resource_group_id` (`resource_group_id`),
    KEY `idx_status` (`status`),
    KEY `idx_ip_address` (`ip_address`),
    CONSTRAINT `fk_servers_resource_group` FOREIGN KEY (`resource_group_id`) REFERENCES `resource_groups` (`id`) ON DELETE SET NULL,
    CONSTRAINT `fk_servers_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务器表';

-- 客户端表
CREATE TABLE IF NOT EXISTS `server_agents` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `server_id` BIGINT UNSIGNED NOT NULL COMMENT '服务器ID',
    `container_id` VARCHAR(64) NULL COMMENT '容器ID',
    `agent_ip` VARCHAR(45) NULL COMMENT '客户端IP',
    `agent_port` INT NOT NULL DEFAULT 22 COMMENT '客户端端口',
    `version` VARCHAR(32) NULL COMMENT '客户端版本',
    `status` ENUM('UNKNOWN', 'ONLINE', 'OFFLINE') NOT NULL DEFAULT 'UNKNOWN' COMMENT '客户端状态',
    `last_heartbeat_at` DATETIME NULL COMMENT '最后心跳时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_server_id` (`server_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_server_agents_server` FOREIGN KEY (`server_id`) REFERENCES `servers` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户端表';

-- 应用实例表
CREATE TABLE IF NOT EXISTS `app_instances` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT '项目ID',
    `name` VARCHAR(64) NOT NULL COMMENT '应用实例名称',
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT '应用模板ID',
    `server_id` BIGINT UNSIGNED NOT NULL COMMENT '服务器ID',
    `container_id` VARCHAR(64) NULL COMMENT '容器ID',
    `container_name` VARCHAR(64) NULL COMMENT '容器名称',
    `image_name` VARCHAR(255) NULL COMMENT '镜像名称',
    `image_tag` VARCHAR(100) NULL COMMENT '镜像标签',
    `status` ENUM('DEFAULT', 'DEPLOYMENT', 'RUNNING', 'PAUSED', 'STOPPED', 'UPDATE') NOT NULL DEFAULT 'DEFAULT' COMMENT '应用状态',
    `started_at` DATETIME NULL COMMENT '启动时间',
    `stopped_at` DATETIME NULL COMMENT '停止时间',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_project_id` (`project_id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_server_id` (`server_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_app_instances_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_app_instances_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_app_instances_server` FOREIGN KEY (`server_id`) REFERENCES `servers` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_app_instances_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用实例表';

-- SSL证书表
CREATE TABLE IF NOT EXISTS `ssl_certificates` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '证书名称',
    `domain` VARCHAR(255) NOT NULL COMMENT '域名',
    `certificate_type` ENUM('LETS_ENCRYPT', 'COMMERCIAL', 'SELF_SIGNED') NOT NULL DEFAULT 'LETS_ENCRYPT' COMMENT '证书类型',
    `certificate_data` TEXT NOT NULL COMMENT '证书内容',
    `private_key_data` TEXT NOT NULL COMMENT '私钥内容',
    `certificate_chain` TEXT NULL COMMENT '证书链',
    `issuer` VARCHAR(255) NULL COMMENT '颁发者',
    `subject` VARCHAR(255) NULL COMMENT '主题',
    `serial_number` VARCHAR(64) NULL COMMENT '序列号',
    `not_before` DATETIME NULL COMMENT '生效时间',
    `not_after` DATETIME NULL COMMENT '过期时间',
    `auto_renew` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '自动续期',
    `status` ENUM('PENDING', 'VALID', 'EXPIRED', 'REVOKED') NOT NULL DEFAULT 'PENDING' COMMENT '状态',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_domain` (`domain`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    KEY `idx_not_after` (`not_after`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SSL证书表';

-- 密钥管理表
CREATE TABLE IF NOT EXISTS `secret_keys` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '密钥名称',
    `key_type` ENUM('API_KEY', 'DATABASE', 'SSH', 'CERTIFICATE', 'CUSTOM') NOT NULL COMMENT '密钥类型',
    `encrypted_value` TEXT NOT NULL COMMENT '加密值',
    `description` TEXT NULL COMMENT '描述',
    `custom_fields` JSON NULL COMMENT '自定义字段',
    `expires_at` DATETIME NULL COMMENT '过期时间',
    `resource_group_id` BIGINT UNSIGNED NULL COMMENT '资源组ID',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_resource_group_id` (`resource_group_id`),
    KEY `idx_key_type` (`key_type`),
    KEY `idx_expires_at` (`expires_at`),
    CONSTRAINT `fk_secret_keys_resource_group` FOREIGN KEY (`resource_group_id`) REFERENCES `resource_groups` (`id`) ON DELETE SET NULL,
    CONSTRAINT `fk_secret_keys_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='密钥管理表';

-- 应用网关表
CREATE TABLE IF NOT EXISTS `app_gateways` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '网关名称',
    `server_id` BIGINT UNSIGNED NOT NULL COMMENT '服务器ID',
    `container_id` VARCHAR(64) NULL COMMENT '网关容器ID',
    `description` TEXT NULL COMMENT '网关描述',
    `status` ENUM('DEFAULT', 'RUNNING', 'STOPPED', 'ERROR') NOT NULL DEFAULT 'DEFAULT' COMMENT '网关状态',
    `started_at` DATETIME NULL COMMENT '启动时间',
    `stopped_at` DATETIME NULL COMMENT '停止时间',
    `resource_group_id` BIGINT UNSIGNED NULL COMMENT '资源组ID',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_server_id` (`server_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_resource_group_id` (`resource_group_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_app_gateways_server` FOREIGN KEY (`server_id`) REFERENCES `servers` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_app_gateways_resource_group` FOREIGN KEY (`resource_group_id`) REFERENCES `resource_groups` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用网关表';

-- 应用网关发布表
CREATE TABLE IF NOT EXISTS `app_gateways_publishes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `app_instance_id` BIGINT UNSIGNED NOT NULL COMMENT '应用实例ID',
    `app_gateway_id` BIGINT UNSIGNED NOT NULL COMMENT '应用网关ID',
    `service_domain` VARCHAR(255) NOT NULL COMMENT '域名(服务名称)',
    `service_port` INT NOT NULL DEFAULT 8080 COMMENT '服务端口',
    `alert_rule_id` BIGINT UNSIGNED NULL COMMENT '监控告警规则ID',
    `limit_rules` TEXT NULL COMMENT '访问控制策略',
    `health_check_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '健康检查开启',
    `audit_log_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '审计日志开启',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_app_instance_id` (`app_instance_id`),
    KEY `idx_app_gateway_id` (`app_gateway_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_service_domain` (`service_domain`),
    CONSTRAINT `fk_gateway_publishes_instance` FOREIGN KEY (`app_instance_id`) REFERENCES `app_instances` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_gateway_publishes_gateway` FOREIGN KEY (`app_gateway_id`) REFERENCES `app_gateways` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用网关发布表';

-- 应用网关访问控制规则表
CREATE TABLE IF NOT EXISTS `app_gateway_access_rules` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `gateway_id` BIGINT UNSIGNED NOT NULL COMMENT '网关ID',
    `rule_name` VARCHAR(64) NOT NULL COMMENT '规则名称',
    `rule_type` ENUM('IP_WHITELIST', 'IP_BLACKLIST', 'RATE_LIMIT') NOT NULL COMMENT '规则类型',
    `limit_count` INT NOT NULL COMMENT '限制次数',
    `time_window` INT NOT NULL COMMENT '时间窗口(秒)',
    `target_path` VARCHAR(255) NULL COMMENT '目标路径',
    `target_ip` VARCHAR(45) NULL COMMENT '目标IP',
    `action` ENUM('BLOCK', 'ALLOW', 'REDIRECT') NOT NULL DEFAULT 'BLOCK' COMMENT '触发动作',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_gateway_id` (`gateway_id`),
    KEY `idx_rule_type` (`rule_type`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_gateway_access_rules_gateway` FOREIGN KEY (`gateway_id`) REFERENCES `app_gateways` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用网关访问控制规则表';

ALTER TABLE `app_gateways_publishes` ADD CONSTRAINT `fk_gateway_publishes_alert_rule` FOREIGN KEY (`alert_rule_id`) REFERENCES `alert_rules` (`id`) ON DELETE SET NULL;
ALTER TABLE `database_connections` ADD CONSTRAINT `fk_database_connections_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT;
ALTER TABLE `ssl_certificates` ADD CONSTRAINT `fk_ssl_certificates_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT;
ALTER TABLE `app_gateways` ADD CONSTRAINT `fk_app_gateways_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT;
ALTER TABLE `app_gateways_publishes` ADD CONSTRAINT `fk_gateway_publishes_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT;

-- ========================================
-- 3.3 应用市场 (App Store)
-- ========================================

-- 应用分类表
CREATE TABLE IF NOT EXISTS `app_store_categories` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '分类名称',
    `code` VARCHAR(32) NOT NULL UNIQUE COMMENT '分类编码',
    `parent_id` BIGINT UNSIGNED NULL COMMENT '父分类ID',
    `icon` VARCHAR(255) NULL COMMENT '图标URL',
    `description` TEXT NULL COMMENT '分类描述',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_app_categories_parent` FOREIGN KEY (`parent_id`) REFERENCES `app_store_categories` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用分类表';

-- 应用模板表
CREATE TABLE IF NOT EXISTS `app_store_templates` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '应用名称',
    `code` VARCHAR(32) NOT NULL UNIQUE COMMENT '应用编码',
    `category_id` BIGINT UNSIGNED NOT NULL COMMENT '分类ID',
    `version` VARCHAR(32) NOT NULL COMMENT '版本号',
    `icon` VARCHAR(255) NULL COMMENT '图标URL',
    `description` TEXT NULL COMMENT '应用描述',
    `official_url` VARCHAR(255) NULL COMMENT '官方网站',
    `source_url` VARCHAR(255) NULL COMMENT '源码地址',
    `compose_template` TEXT NOT NULL COMMENT 'Docker Compose模板',
    `download_count` INT NOT NULL DEFAULT 0 COMMENT '下载次数',
    `star_count` INT NOT NULL DEFAULT 0 COMMENT '点赞数',
    `rating` DECIMAL(3,2) NOT NULL DEFAULT 0.00 COMMENT '评分',
    `is_official` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否官方应用',
    `is_featured` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否推荐应用',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-下架，1-上架',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_status` (`status`),
    KEY `idx_is_featured` (`is_featured`),
    KEY `idx_download_count` (`download_count`),
    CONSTRAINT `fk_app_templates_category` FOREIGN KEY (`category_id`) REFERENCES `app_store_categories` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用模板表';

-- 应用心愿单表
CREATE TABLE IF NOT EXISTS `app_store_wishlists` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '应用名称',
    `version` VARCHAR(32) NULL COMMENT '版本号',
    `source_url` VARCHAR(255) NULL COMMENT '来源地址',
    `description` TEXT NULL COMMENT '需求描述',
    `reward_amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '悬赏金额',
    `priority` TINYINT(1) NOT NULL DEFAULT 3 COMMENT '优先级：1-高，2-中，3-低',
    `status` ENUM('PENDING', 'IN_PROGRESS', 'COMPLETED', 'EXPIRED') NOT NULL DEFAULT 'PENDING' COMMENT '状态',
    `view_count` INT NOT NULL DEFAULT 0 COMMENT '浏览数',
    `like_count` INT NOT NULL DEFAULT 0 COMMENT '点赞数',
    `vote_count` INT NOT NULL DEFAULT 0 COMMENT '投票数',
    `comment_count` INT NOT NULL DEFAULT 0 COMMENT '评论数',
    `submitter_id` BIGINT UNSIGNED NOT NULL COMMENT '提交者ID',
    `completed_at` DATETIME NULL COMMENT '完成时间',
    `expires_at` DATETIME NULL COMMENT '过期时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_submitter_id` (`submitter_id`),
    KEY `idx_status` (`status`),
    KEY `idx_priority` (`priority`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用心愿单表';

-- 应用评价表
CREATE TABLE IF NOT EXISTS `app_store_reviews` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT '应用模板ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `rating` TINYINT(1) NOT NULL COMMENT '评分（1-5分）',
    `content` TEXT NULL COMMENT '评价内容',
    `tags` JSON NULL COMMENT '评价标签',
    `is_helpful` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否有帮助',
    `helpful_count` INT NOT NULL DEFAULT 0 COMMENT '有帮助数量',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_rating` (`rating`),
    CONSTRAINT `fk_reviews_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE CASCADE,
    CONSTRAINT `chk_rating` CHECK (`rating` >= 1 AND `rating` <= 5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用评价表';

-- 应用收藏表
CREATE TABLE IF NOT EXISTS `app_store_favorites` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT '应用模板ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_template_user` (`template_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_favorites_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用收藏表';

-- 应用点赞表
CREATE TABLE IF NOT EXISTS `app_store_stars` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT '应用模板ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_template_user` (`template_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_stars_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用点赞表';

-- 应用举报表
CREATE TABLE IF NOT EXISTS `app_store_reports` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT '应用模板ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '举报用户ID',
    `reason` ENUM('SPAM', 'INAPPROPRIATE', 'COPYRIGHT') NOT NULL COMMENT '举报原因',
    `detail` TEXT NULL COMMENT '详细说明',
    `status` ENUM('PENDING', 'PROCESSING', 'RESOLVED', 'REJECTED') NOT NULL DEFAULT 'PENDING' COMMENT '处理状态',
    `handled_by` BIGINT UNSIGNED NULL COMMENT '处理人ID',
    `handled_at` DATETIME NULL COMMENT '处理时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    KEY `idx_handled_by` (`handled_by`),
    CONSTRAINT `fk_app_reports_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_app_reports_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_app_reports_handled_by` FOREIGN KEY (`handled_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用举报表';

-- 应用下载记录表
CREATE TABLE IF NOT EXISTS `app_store_downloads` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT '应用模板ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `ip_address` VARCHAR(45) NULL COMMENT 'IP地址',
    `user_agent` VARCHAR(255) NULL COMMENT '用户代理',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_app_downloads_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_app_downloads_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用下载记录表';

-- 应用心愿单评论表
CREATE TABLE IF NOT EXISTS `app_store_wishlist_comments` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `wishlist_id` BIGINT UNSIGNED NOT NULL COMMENT '心愿单ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `parent_id` BIGINT UNSIGNED NULL COMMENT '父评论ID',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `like_count` INT NOT NULL DEFAULT 0 COMMENT '点赞数',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_wishlist_id` (`wishlist_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_parent_id` (`parent_id`),
    CONSTRAINT `fk_wishlist_comments_wishlist` FOREIGN KEY (`wishlist_id`) REFERENCES `app_store_wishlists` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_wishlist_comments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_wishlist_comments_parent` FOREIGN KEY (`parent_id`) REFERENCES `app_store_wishlist_comments` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用心愿单评论表';

-- 应用心愿单投票表
CREATE TABLE IF NOT EXISTS `app_store_wishlist_votes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `wishlist_id` BIGINT UNSIGNED NOT NULL COMMENT '心愿单ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_wishlist_user` (`wishlist_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_wishlist_votes_wishlist` FOREIGN KEY (`wishlist_id`) REFERENCES `app_store_wishlists` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_wishlist_votes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用心愿单投票表';

-- 应用心愿单点赞表
CREATE TABLE IF NOT EXISTS `app_store_wishlist_likes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `wishlist_id` BIGINT UNSIGNED NOT NULL COMMENT '心愿单ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_wishlist_user` (`wishlist_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_wishlist_likes_wishlist` FOREIGN KEY (`wishlist_id`) REFERENCES `app_store_wishlists` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_wishlist_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用心愿单点赞表';

-- 应用心愿单举报表
CREATE TABLE IF NOT EXISTS `app_store_wishlist_reports` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `wishlist_id` BIGINT UNSIGNED NOT NULL COMMENT '心愿单ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '举报用户ID',
    `reason` ENUM('SPAM', 'INAPPROPRIATE', 'DUPLICATE') NOT NULL COMMENT '举报原因',
    `detail` TEXT NULL COMMENT '详细说明',
    `status` ENUM('PENDING', 'PROCESSING', 'RESOLVED', 'REJECTED') NOT NULL DEFAULT 'PENDING' COMMENT '处理状态',
    `handled_by` BIGINT UNSIGNED NULL COMMENT '处理人ID',
    `handled_at` DATETIME NULL COMMENT '处理时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_wishlist_id` (`wishlist_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    KEY `idx_handled_by` (`handled_by`),
    CONSTRAINT `fk_wishlist_reports_wishlist` FOREIGN KEY (`wishlist_id`) REFERENCES `app_store_wishlists` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_wishlist_reports_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_wishlist_reports_handled_by` FOREIGN KEY (`handled_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用心愿单举报表';

-- 应用部署记录表
CREATE TABLE IF NOT EXISTS `app_deployments` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `deployment_id` VARCHAR(64) NOT NULL UNIQUE COMMENT '部署唯一标识',
    `template_id` BIGINT UNSIGNED NULL COMMENT '应用模板ID',
    `app_instance_id` BIGINT UNSIGNED NULL COMMENT '应用实例ID',
    `server_id` BIGINT UNSIGNED NOT NULL COMMENT '服务器ID',
    `status` ENUM('PENDING', 'RUNNING', 'SUCCESS', 'FAILED', 'CANCELLED') NOT NULL DEFAULT 'PENDING' COMMENT '部署状态',
    `progress` TINYINT NOT NULL DEFAULT 0 COMMENT '部署进度',
    `estimated_time` INT NOT NULL DEFAULT 0 COMMENT '预计时间(秒)',
    `start_time` DATETIME NULL COMMENT '开始时间',
    `end_time` DATETIME NULL COMMENT '结束时间',
    `error_message` TEXT NULL COMMENT '错误信息',
    `deployment_log` TEXT NULL COMMENT '部署日志',
    `config_data` JSON NULL COMMENT '部署配置',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_server_id` (`server_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_deployments_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用部署记录表';

ALTER TABLE `app_shortcuts` ADD CONSTRAINT `fk_app_shortcuts_app_instance` FOREIGN KEY (`app_instance_id`) REFERENCES `app_instances` (`id`) ON DELETE CASCADE;
ALTER TABLE `app_shortcuts` ADD CONSTRAINT `fk_app_shortcuts_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
ALTER TABLE `app_store_wishlists` ADD CONSTRAINT `fk_app_wishlists_submitter` FOREIGN KEY (`submitter_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

-- ========================================
-- 3.4 平台管理 (Platform Management)
-- ========================================

-- 系统配置表
CREATE TABLE IF NOT EXISTS `system_configs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `config_key` VARCHAR(64) NOT NULL UNIQUE COMMENT '配置键',
    `config_value` TEXT NULL COMMENT '配置值',
    `config_type` ENUM('STRING', 'INTEGER', 'BOOLEAN', 'JSON', 'FLOAT') NOT NULL DEFAULT 'STRING' COMMENT '配置类型',
    `category` VARCHAR(32) NOT NULL COMMENT '配置分类',
    `description` TEXT NULL COMMENT '配置描述',
    `is_readonly` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否只读',
    `is_encrypted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否加密',
    `default_value` TEXT NULL COMMENT '默认值',
    `validation_rules` JSON NULL COMMENT '验证规则',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_category` (`category`),
    KEY `idx_config_type` (`config_type`),
    KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- Webhook配置表
CREATE TABLE IF NOT EXISTS `webhook_configs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Webhook名称',
    `url` VARCHAR(255) NOT NULL COMMENT '回调URL',
    `secret` VARCHAR(255) NULL COMMENT '签名密钥',
    `events` JSON NOT NULL COMMENT '监听事件',
    `headers` JSON NULL COMMENT '自定义请求头',
    `timeout` INT NOT NULL DEFAULT 30 COMMENT '超时时间(秒)',
    `retry_count` INT NOT NULL DEFAULT 3 COMMENT '重试次数',
    `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_is_active` (`is_active`),
    CONSTRAINT `fk_webhook_configs_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Webhook配置表';

-- Webhook执行日志表
CREATE TABLE IF NOT EXISTS `webhook_logs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `webhook_id` BIGINT UNSIGNED NOT NULL COMMENT 'WebhookID',
    `event_type` VARCHAR(32) NOT NULL COMMENT '事件类型',
    `payload` JSON NOT NULL COMMENT '请求载荷',
    `request_id` VARCHAR(64) NOT NULL COMMENT '请求ID',
    `status_code` INT NULL COMMENT '响应状态码',
    `response_body` TEXT NULL COMMENT '响应内容',
    `response_time` INT NOT NULL DEFAULT 0 COMMENT '响应时间(毫秒)',
    `retry_count` INT NOT NULL DEFAULT 0 COMMENT '重试次数',
    `success` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否成功',
    `error_message` TEXT NULL COMMENT '错误信息',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_webhook_id` (`webhook_id`),
    KEY `idx_event_type` (`event_type`),
    KEY `idx_success` (`success`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_webhook_logs_webhook` FOREIGN KEY (`webhook_id`) REFERENCES `webhook_configs` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Webhook执行日志表';

-- 软件源配置表
CREATE TABLE IF NOT EXISTS `repository_configs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '配置名称',
    `type` ENUM('docker', 'apt', 'yum', 'npm', 'pip') NOT NULL COMMENT '仓库类型',
    `url` VARCHAR(255) NOT NULL COMMENT '仓库地址',
    `username` VARCHAR(64) NULL COMMENT '用户名',
    `password` VARCHAR(255) NULL COMMENT '密码',
    `is_default` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否默认',
    `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统配置',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='软件源配置表';

-- 平台更新记录表
CREATE TABLE IF NOT EXISTS `platform_updates` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `version` VARCHAR(32) NOT NULL COMMENT '版本号',
    `changelog` TEXT NULL COMMENT '更新日志',
    `download_url` VARCHAR(255) NULL COMMENT '下载地址',
    `file_size` BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小',
    `checksum` VARCHAR(64) NULL COMMENT '文件校验和',
    `status` ENUM('PENDING', 'DOWNLOADING', 'INSTALLING', 'SUCCESS', 'FAILED') NOT NULL DEFAULT 'PENDING' COMMENT '更新状态',
    `started_at` DATETIME NULL COMMENT '开始时间',
    `completed_at` DATETIME NULL COMMENT '完成时间',
    `error_message` TEXT NULL COMMENT '错误信息',
    `backup_path` VARCHAR(1024) NULL COMMENT '备份路径',
    `updated_by` BIGINT UNSIGNED NULL COMMENT '更新人ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_version` (`version`),
    KEY `idx_status` (`status`),
    KEY `idx_updated_by` (`updated_by`),
    CONSTRAINT `fk_platform_updates_updated_by` FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='平台更新记录表';

-- 容器集群节点表
CREATE TABLE IF NOT EXISTS `docker_swarm_nodes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `node_id` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Docker节点ID',
    `hostname` VARCHAR(255) NOT NULL COMMENT '主机名',
    `ip_address` VARCHAR(45) NOT NULL COMMENT 'IP地址',
    `role` ENUM('manager', 'worker') NOT NULL COMMENT '节点角色',
    `status` ENUM('ready', 'down', 'unknown') NOT NULL DEFAULT 'ready' COMMENT '节点状态',
    `availability` ENUM('active', 'pause', 'drain') NOT NULL DEFAULT 'active' COMMENT '可用性',
    `engine_version` VARCHAR(32) NULL COMMENT '引擎版本',
    `labels` JSON NULL COMMENT '节点标签',
    `resources` JSON NULL COMMENT '资源信息',
    `joined_at` DATETIME NULL COMMENT '加入时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_role` (`role`),
    KEY `idx_status` (`status`),
    KEY `idx_availability` (`availability`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='容器集群节点表';

-- 容器镜像表
CREATE TABLE IF NOT EXISTS `docker_images` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `image_id` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Docker镜像ID',
    `repository` VARCHAR(255) NOT NULL COMMENT '镜像仓库',
    `tag` VARCHAR(100) NOT NULL COMMENT '镜像标签',
    `digest` VARCHAR(128) NULL COMMENT '镜像摘要',
    `size` BIGINT NOT NULL DEFAULT 0 COMMENT '镜像大小',
    `architecture` VARCHAR(32) NULL COMMENT '架构',
    `os` VARCHAR(32) NULL COMMENT '操作系统',
    `status` ENUM('AVAILABLE', 'PULLING', 'ERROR') NOT NULL DEFAULT 'AVAILABLE' COMMENT '镜像状态',
    `usage_count` INT NOT NULL DEFAULT 0 COMMENT '使用次数',
    `last_used_at` DATETIME NULL COMMENT '最后使用时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_repository_tag` (`repository`, `tag`),
    KEY `idx_status` (`status`),
    KEY `idx_usage_count` (`usage_count`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='容器镜像表';

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(64) NOT NULL UNIQUE COMMENT '用户名',
    `email` VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
    `nickname` VARCHAR(64) NULL COMMENT '昵称',
    `avatar` VARCHAR(255) NULL COMMENT '头像URL',
    `phone` VARCHAR(20) NULL COMMENT '手机号',
    `gender` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女',
    `signature` VARCHAR(255) NULL COMMENT '个性签名',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `last_login_at` DATETIME NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(45) NULL COMMENT '最后登录IP',
    `timezone` VARCHAR(64) NOT NULL DEFAULT 'UTC' COMMENT '时区',
    `language` VARCHAR(10) NOT NULL DEFAULT 'zh-CN' COMMENT '语言',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`),
    KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 权限表
CREATE TABLE IF NOT EXISTS `permissions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `parent_code` VARCHAR(64) NULL COMMENT '父权限code',
    `scope` VARCHAR(64) NOT NULL COMMENT '权限域：platform, project',
    `name` VARCHAR(64) NOT NULL COMMENT '权限名称',
    `code` VARCHAR(64) NOT NULL UNIQUE COMMENT '权限编码',
    `module` VARCHAR(32) NOT NULL COMMENT '模块名称',
    `action` VARCHAR(32) NOT NULL COMMENT '操作名称',
    `resource` VARCHAR(64) NULL COMMENT '资源标识（URI）',
    `element` VARCHAR(64) NULL COMMENT 'UI元素ID',
    `description` TEXT NULL COMMENT '权限描述',
    `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统权限',
    `is_menu` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否菜单权限',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：-1-删除，0-禁用，1-启用',
    `created_by` BIGINT UNSIGNED NULL COMMENT '创建人ID',
    `updated_by` BIGINT UNSIGNED NULL COMMENT '更新人ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_parent_code` (`parent_code`),
    KEY `idx_scope` (`scope`),
    KEY `idx_module` (`module`),
    KEY `idx_action` (`action`),
    KEY `idx_status` (`status`),
    KEY `idx_created_by` (`created_by`),
    KEY `idx_updated_by` (`updated_by`),
    CONSTRAINT `fk_permissions_parent` FOREIGN KEY (`parent_code`) REFERENCES `permissions` (`code`) ON DELETE SET NULL,
    CONSTRAINT `fk_permissions_created_by` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
    CONSTRAINT `fk_permissions_updated_by` FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 角色表
CREATE TABLE IF NOT EXISTS `roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL UNIQUE COMMENT '角色名称',
    `code` VARCHAR(32) NOT NULL UNIQUE COMMENT '角色编码',
    `description` TEXT NULL COMMENT '角色描述',
    `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统角色',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：-1-删除，0-禁用，1-启用',
    `created_by` BIGINT UNSIGNED NULL COMMENT '创建人ID',
    `updated_by` BIGINT UNSIGNED NULL COMMENT '更新人ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_by` (`created_by`),
    KEY `idx_updated_by` (`updated_by`),
    CONSTRAINT `fk_roles_created_by` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
    CONSTRAINT `fk_roles_updated_by` FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS `user_roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    `granted_by` BIGINT UNSIGNED NULL COMMENT '授权人ID',
    `granted_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '授权时间',
    `expires_at` DATETIME NULL COMMENT '过期时间',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：-1-删除，0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
    KEY `idx_role_id` (`role_id`),
    KEY `idx_granted_by` (`granted_by`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_user_roles_granted_by` FOREIGN KEY (`granted_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS `role_permissions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    `permission_code` VARCHAR(64) NOT NULL COMMENT '权限code',
    `granted_by` BIGINT UNSIGNED NULL COMMENT '授权人ID',
    `granted_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '授权时间',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：-1-删除，0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_permission` (`role_id`, `permission_code`),
    KEY `idx_permission_code` (`permission_code`),
    KEY `idx_granted_by` (`granted_by`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_code`) REFERENCES `permissions` (`code`) ON DELETE CASCADE,
    CONSTRAINT `fk_role_permissions_granted_by` FOREIGN KEY (`granted_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- API访问令牌表
CREATE TABLE IF NOT EXISTS `api_tokens` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '令牌名称',
    `token` VARCHAR(255) NOT NULL UNIQUE COMMENT '令牌值（加密存储）',
    `token_hash` VARCHAR(64) NOT NULL UNIQUE COMMENT '令牌哈希值',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `scopes` JSON NULL COMMENT '权限范围',
    `description` TEXT NULL COMMENT '令牌描述',
    `last_used_at` DATETIME NULL COMMENT '最后使用时间',
    `last_used_ip` VARCHAR(45) NULL COMMENT '最后使用IP',
    `expires_at` DATETIME NULL COMMENT '过期时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_expires_at` (`expires_at`),
    CONSTRAINT `fk_api_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API访问令牌表';

-- 用户双因子认证表
CREATE TABLE IF NOT EXISTS `user_two_factor` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `method` VARCHAR(32) NOT NULL COMMENT '认证方法（TOTP、EMAIL）',
    `secret` VARCHAR(255) NULL COMMENT '密钥（加密存储）',
    `backup_codes` JSON NULL COMMENT '备用码',
    `email` VARCHAR(255) NULL COMMENT '邮箱',
    `enabled` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否启用',
    `verified_at` DATETIME NULL COMMENT '验证时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_method` (`user_id`, `method`),
    KEY `idx_enabled` (`enabled`),
    CONSTRAINT `fk_user_two_factor_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户双因子认证表';

-- 用户个人中心配置表
CREATE TABLE IF NOT EXISTS `user_profile` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `config_key` VARCHAR(64) NOT NULL COMMENT '配置键',
    `config_value` TEXT NULL COMMENT '配置值',
    `description` TEXT NULL COMMENT '配置描述',
    `is_readonly` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否只读',
    `is_encrypted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否加密',
    `default_value` TEXT NULL COMMENT '默认值',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_config_key` (`user_id`, `config_key`),
    KEY `idx_sort_order` (`sort_order`),
    CONSTRAINT `fk_user_profile_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户个人中心配置表';

-- 用户登录历史表
CREATE TABLE IF NOT EXISTS `user_login_history` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `ip_address` VARCHAR(45) NULL COMMENT 'IP地址',
    `user_agent` VARCHAR(255) NULL COMMENT '用户代理',
    `location` VARCHAR(100) NULL COMMENT '登录地点',
    `device` VARCHAR(100) NULL COMMENT '设备信息',
    `browser` VARCHAR(100) NULL COMMENT '浏览器信息',
    `login_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    `logout_time` DATETIME NULL COMMENT '登出时间',
    `status` ENUM('ACTIVE', 'EXPIRED', 'LOGOUT') NOT NULL DEFAULT 'ACTIVE' COMMENT '会话状态',
    `session_id` VARCHAR(128) NULL COMMENT '会话ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_login_time` (`login_time`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_user_login_history_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户登录历史表';

-- 告警规则表
CREATE TABLE IF NOT EXISTS `alert_rules` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '规则名称',
    `rule_type` ENUM('SYSTEM', 'APPLICATION', 'CUSTOM') NOT NULL COMMENT '规则类型',
    `target_type` ENUM('SERVER', 'APPLICATION', 'SERVICE') NOT NULL COMMENT '目标类型',
    `target_id` BIGINT UNSIGNED NULL COMMENT '目标ID',
    `metric_name` VARCHAR(64) NULL COMMENT '指标名称',
    `condition_expression` TEXT NOT NULL COMMENT '条件表达式',
    `notification_channels` JSON NULL COMMENT '通知渠道',
    `is_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_rule_type` (`rule_type`),
    KEY `idx_target_type` (`target_type`),
    KEY `idx_target_id` (`target_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_is_enabled` (`is_enabled`),
    CONSTRAINT `fk_alert_rules_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警规则表';

-- 告警记录表
CREATE TABLE IF NOT EXISTS `alert_records` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `alert_rule_id` BIGINT UNSIGNED NOT NULL COMMENT '告警规则ID',
    `alert_id` VARCHAR(64) NOT NULL UNIQUE COMMENT '告警ID',
    `title` VARCHAR(255) NOT NULL COMMENT '告警标题',
    `description` TEXT NULL COMMENT '告警描述',
    `status` ENUM('FIRING', 'RESOLVED', 'ACKNOWLEDGED') NOT NULL DEFAULT 'FIRING' COMMENT '状态',
    `fired_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '触发时间',
    `resolved_at` DATETIME NULL COMMENT '解决时间',
    `acknowledged_at` DATETIME NULL COMMENT '确认时间',
    `acknowledged_by` BIGINT UNSIGNED NULL COMMENT '确认人ID',
    `resolution_note` TEXT NULL COMMENT '解决说明',
    `notification_sent` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '通知已发送',
    `notification_channels` JSON NULL COMMENT '通知渠道',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_alert_rule_id` (`alert_rule_id`),
    KEY `idx_status` (`status`),
    KEY `idx_fired_at` (`fired_at`),
    KEY `idx_acknowledged_by` (`acknowledged_by`),
    CONSTRAINT `fk_alert_records_rule` FOREIGN KEY (`alert_rule_id`) REFERENCES `alert_rules` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_alert_records_acknowledged_by` FOREIGN KEY (`acknowledged_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警记录表';

-- 通知模板表
CREATE TABLE IF NOT EXISTS `notification_templates` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '模板名称',
    `type` ENUM('EMAIL', 'SMS', 'WEBHOOK', 'PUSH') NOT NULL COMMENT '通知类型',
    `subject` VARCHAR(255) NULL COMMENT '通知主题',
    `content` TEXT NOT NULL COMMENT '通知内容模板',
    `variables` JSON NULL COMMENT '模板变量',
    `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统模板',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知模板表';

-- 通知记录表
CREATE TABLE IF NOT EXISTS `notification_records` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NULL COMMENT '模板ID',
    `type` ENUM('EMAIL', 'SMS', 'WEBHOOK', 'PUSH') NOT NULL COMMENT '通知类型',
    `recipient` VARCHAR(255) NOT NULL COMMENT '接收者',
    `subject` VARCHAR(255) NULL COMMENT '通知主题',
    `content` TEXT NOT NULL COMMENT '通知内容',
    `status` ENUM('PENDING', 'SENDING', 'SUCCESS', 'FAILED') NOT NULL DEFAULT 'PENDING' COMMENT '发送状态',
    `sent_at` DATETIME NULL COMMENT '发送时间',
    `error_msg` TEXT NULL COMMENT '错误信息',
    `retry_count` INT NOT NULL DEFAULT 0 COMMENT '重试次数',
    `reference_id` VARCHAR(64) NULL COMMENT '关联ID',
    `reference_type` VARCHAR(32) NULL COMMENT '关联类型',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`),
    KEY `idx_reference` (`reference_type`, `reference_id`),
    CONSTRAINT `fk_notification_records_template` FOREIGN KEY (`template_id`) REFERENCES `notification_templates` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知记录表';

-- 通知消息表
CREATE TABLE IF NOT EXISTS `notifications` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(255) NOT NULL COMMENT '消息标题',
    `content` TEXT NOT NULL COMMENT '消息内容',
    `type` ENUM('SYSTEM', 'ALERT', 'TASK', 'CUSTOM') NOT NULL COMMENT '消息类型',
    `level` ENUM('INFO', 'WARNING', 'ERROR', 'SUCCESS') NOT NULL DEFAULT 'INFO' COMMENT '消息级别',
    `sender_id` BIGINT UNSIGNED NULL COMMENT '发送者ID',
    `target_type` ENUM('USER', 'GROUP', 'ALL') NOT NULL COMMENT '目标类型',
    `target_ids` JSON NULL COMMENT '目标ID列表',
    `channels` JSON NULL COMMENT '发送渠道',
    `template_id` BIGINT UNSIGNED NULL COMMENT '消息模板ID',
    `variables` JSON NULL COMMENT '模板变量',
    `sent_at` DATETIME NULL COMMENT '发送时间',
    `status` ENUM('PENDING', 'SENDING', 'SUCCESS', 'FAILED') NOT NULL DEFAULT 'PENDING' COMMENT '发送状态',
    `error_msg` TEXT NULL COMMENT '错误信息',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_type` (`type`),
    KEY `idx_level` (`level`),
    KEY `idx_sender_id` (`sender_id`),
    KEY `idx_target_type` (`target_type`),
    KEY `idx_status` (`status`),
    KEY `idx_template_id` (`template_id`),
    CONSTRAINT `fk_notifications_sender` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`) ON DELETE SET NULL,
    CONSTRAINT `fk_notifications_template` FOREIGN KEY (`template_id`) REFERENCES `notification_templates` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知消息表';

-- 用户通知记录表
CREATE TABLE IF NOT EXISTS `user_notifications` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `notification_id` BIGINT UNSIGNED NOT NULL COMMENT '通知ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `is_read` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已读',
    `read_at` DATETIME NULL COMMENT '阅读时间',
    `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
    `deleted_at` DATETIME NULL COMMENT '删除时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_notification_user` (`notification_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_is_read` (`is_read`),
    KEY `idx_is_deleted` (`is_deleted`),
    CONSTRAINT `fk_user_notifications_notification` FOREIGN KEY (`notification_id`) REFERENCES `notifications` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_user_notifications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户通知记录表';

-- 审计日志表
CREATE TABLE IF NOT EXISTS `audit_logs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NULL COMMENT '用户ID',
    `username` VARCHAR(64) NULL COMMENT '用户名',
    `action` VARCHAR(32) NOT NULL COMMENT '操作动作',
    `module` VARCHAR(32) NOT NULL COMMENT '模块名称',
    `resource_type` VARCHAR(32) NULL COMMENT '资源类型',
    `resource_id` BIGINT UNSIGNED NULL COMMENT '资源ID',
    `resource_name` VARCHAR(64) NULL COMMENT '资源名称',
    `description` TEXT NULL COMMENT '操作描述',
    `ip_address` VARCHAR(45) NULL COMMENT 'IP地址',
    `user_agent` VARCHAR(255) NULL COMMENT '用户代理',
    `request_method` VARCHAR(10) NULL COMMENT '请求方法',
    `request_url` VARCHAR(255) NULL COMMENT '请求URL',
    `request_params` JSON NULL COMMENT '请求参数',
    `response_status` INT NULL COMMENT '响应状态',
    `response_time` INT NULL COMMENT '响应时间(毫秒)',
    `success` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否成功',
    `error_message` TEXT NULL COMMENT '错误信息',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_action` (`action`),
    KEY `idx_module` (`module`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_success` (`success`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';


-- ========================================
-- 初始数据插入
-- ========================================

-- 插入默认角色
INSERT INTO `roles` (`name`, `code`, `description`, `is_system`, `status`) VALUES
('超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 1, 1),
('系统管理员', 'admin', '系统管理员，负责平台管理', 1, 1),
('项目管理员', 'project_admin', '项目管理员，负责项目管理', 1, 1),
('开发者', 'developer', '开发者角色，负责应用开发', 1, 1),
('运维人员', 'operator', '运维人员角色，负责系统运维', 1, 1),
('普通用户', 'user', '普通用户角色，基础权限', 1, 1);

-- 插入默认权限
INSERT INTO `permissions` (`parent_code`,`scope`,`name`,`code`,`module`,`action`,`resource`,`element`,`description`,`is_system`,`is_menu`,`sort_order`,`status`,`created_by`,`updated_by`,`created_at`,`updated_at`) VALUES
	 (NULL,'platform','Websoft9','e06e84c7-fe48-4d0d-9dbb-cbe65cf5d936','platform','*',NULL,NULL,'Default permissions root node',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e06e84c7-fe48-4d0d-9dbb-cbe65cf5d936','platform','主页','16aa0ce0-ce2a-4075-8aee-5f9ccb472819','home','*',NULL,NULL,'平台主页',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('16aa0ce0-ce2a-4075-8aee-5f9ccb472819','platform','项目总览看板','b2fabec8-f7fe-480d-883c-5e0a355b8c69','project_overview','*',NULL,NULL,'项目总览仪表盘全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b2fabec8-f7fe-480d-883c-5e0a355b8c69','platform','项目总览看板查询','e0aec98a-101f-494a-b5ea-ada5cb5071d1','project_overview','query',NULL,NULL,'项目总览仪表盘查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('16aa0ce0-ce2a-4075-8aee-5f9ccb472819','platform','监控总览看板','b61077db-d098-42d0-ade9-472038167183','monitor_overview','*',NULL,NULL,'监控总览仪表盘全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b61077db-d098-42d0-ade9-472038167183','platform','监控总览看板查询','91ab578c-1c7e-4f56-9708-4d00f4ad3c2d','monitor_overview','query',NULL,NULL,'监控总览仪表盘查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('16aa0ce0-ce2a-4075-8aee-5f9ccb472819','platform','应用快捷导航','8422b16f-b03f-4458-9159-3b78e571a88c','app_navigation','*',NULL,NULL,'应用快捷访问导航全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('8422b16f-b03f-4458-9159-3b78e571a88c','platform','应用快捷导航查询','9f99aee4-fffc-4ad6-b49c-7bf2567352cc','app_navigation','query',NULL,NULL,'应用快捷访问导航查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('8422b16f-b03f-4458-9159-3b78e571a88c','platform','应用快捷导航创建','a52729f5-7300-456c-a90a-1f7e4bd4d0db','app_navigation','create',NULL,NULL,'应用快捷访问导航创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('8422b16f-b03f-4458-9159-3b78e571a88c','platform','应用快捷导航修改','3ea52579-58a8-463e-9a46-0dbec048fb07','app_navigation','update',NULL,NULL,'应用快捷访问导航修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('8422b16f-b03f-4458-9159-3b78e571a88c','platform','应用快捷导航删除','c3393a62-4526-4f8a-bdd9-d6012f6d590f','app_navigation','delete',NULL,NULL,'应用快捷访问导航删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e06e84c7-fe48-4d0d-9dbb-cbe65cf5d936','project','个人空间','08d5484f-ee91-47ba-99f9-ae90c8d5684d','personal_folder','*',NULL,NULL,'个人空间',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('08d5484f-ee91-47ba-99f9-ae90c8d5684d','project','个人空间查询','af309229-e1c0-4053-b8c3-e35ef0b89107','personal_folder','query',NULL,NULL,'个人文件夹查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('08d5484f-ee91-47ba-99f9-ae90c8d5684d','project','个人空间创建','593b5fb7-4146-4f17-9948-bf2827c4cd94','personal_folder','create',NULL,NULL,'个人文件夹创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('08d5484f-ee91-47ba-99f9-ae90c8d5684d','project','个人空间修改','67e63e8c-17b8-49d7-b6cd-6618614d3d5c','personal_folder','update',NULL,NULL,'个人文件夹修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('08d5484f-ee91-47ba-99f9-ae90c8d5684d','project','个人空间删除','d8e80678-a03a-42c3-99d8-29ed4b9032ee','personal_folder','delete',NULL,NULL,'个人文件夹删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('08d5484f-ee91-47ba-99f9-ae90c8d5684d','project','个人空间上传','f47d800e-d504-4de4-946c-60add0123f8d','personal_folder','create',NULL,NULL,'个人文件夹上传权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('08d5484f-ee91-47ba-99f9-ae90c8d5684d','project','个人空间下载','89a1752f-7fda-4949-a241-f5704f713eeb','personal_folder','query',NULL,NULL,'个人文件夹下载权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e06e84c7-fe48-4d0d-9dbb-cbe65cf5d936','project','项目','e1dbef1f-22e8-4b53-9e20-80a83c6b697b','project','*',NULL,NULL,'项目',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e1dbef1f-22e8-4b53-9e20-80a83c6b697b','project','项目仪表盘','dee6721a-870a-413c-8b6c-2fb82b153479','project_dashboard','*',NULL,NULL,'项目仪表盘全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('dee6721a-870a-413c-8b6c-2fb82b153479','project','项目仪表盘查询','b0335079-0f4e-49a2-91ed-3946c0f4d13a','project_dashboard','query',NULL,NULL,'项目仪表盘查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('dee6721a-870a-413c-8b6c-2fb82b153479','project','监控看板查询','222dd995-dbb0-4af3-9ecc-ec372a7fe744','project_dashboard','query',NULL,NULL,'项目监控看板查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('dee6721a-870a-413c-8b6c-2fb82b153479','project','任务看板查询','68e2bbed-27ec-4d6c-95b4-71978b522596','project_dashboard','query',NULL,NULL,'项目任务看板查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('dee6721a-870a-413c-8b6c-2fb82b153479','project','资源看板查询','df53f526-9ae3-4ea8-a137-dea54e9e0306','project_dashboard','query',NULL,NULL,'项目资源看板查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e1dbef1f-22e8-4b53-9e20-80a83c6b697b','project','项目空间','93d6d1db-b814-43ed-864a-4f485af455bd','project_folder','*',NULL,NULL,'项目文件夹管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('93d6d1db-b814-43ed-864a-4f485af455bd','project','项目空间查询','ecf7b9a6-3a73-4189-a3a5-b00fcc404936','project_folder','query',NULL,NULL,'项目文件夹查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('93d6d1db-b814-43ed-864a-4f485af455bd','project','项目空间创建','049c318d-3593-4998-843b-7abff0e47e17','project_folder','create',NULL,NULL,'项目文件夹创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('93d6d1db-b814-43ed-864a-4f485af455bd','project','项目空间修改','375e6ef7-70d1-4f5d-98ce-affc9732a220','project_folder','update',NULL,NULL,'项目文件夹修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('93d6d1db-b814-43ed-864a-4f485af455bd','project','项目空间删除','f0a85d28-e7b9-45db-8f74-f3263fd3b02a','project_folder','delete',NULL,NULL,'项目文件夹删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('93d6d1db-b814-43ed-864a-4f485af455bd','project','项目空间上传','47b481e8-39a8-469f-9d83-a930d1053108','project_folder','create',NULL,NULL,'项目文件夹上传权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('93d6d1db-b814-43ed-864a-4f485af455bd','project','项目空间下载','d246514e-8b6f-4d5d-b154-123768a656df','project_folder','query',NULL,NULL,'项目文件夹下载权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e1dbef1f-22e8-4b53-9e20-80a83c6b697b','project','应用管理','5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','application','*',NULL,NULL,'应用管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','project','应用查询','4ac0b046-dd27-4fb3-acba-f9618237061c','application','query',NULL,NULL,'应用查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','project','应用创建','53122f15-2c25-4aa6-a2e1-ea3690d8c2f5','application','create',NULL,NULL,'应用创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','project','应用修改','56beb0c8-783e-4c87-8e48-048bd13b7d31','application','update',NULL,NULL,'应用修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','project','应用删除','a01c53a1-c58e-4bc3-8c82-6228a358193f','application','delete',NULL,NULL,'应用删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','project','应用安装','ea61a581-62c5-4206-a48a-58c6d468c6bb','application','create',NULL,NULL,'应用安装权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','project','应用卸载','7b1e1414-4ca6-46f1-bc50-4d94c9cebd02','application','delete',NULL,NULL,'应用卸载权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','project','应用启动','6b68510c-bfdf-4c04-897f-9a15d080a8ed','application','update',NULL,NULL,'应用启动权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','project','应用停止','c420843c-3984-4041-b990-a7bd19a91271','application','update',NULL,NULL,'应用停止权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','project','应用重启','0af2937b-d661-427b-93ed-7f1c769cfdda','application','update',NULL,NULL,'应用重启权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','project','应用暂停','8a8db5d8-46b3-4b0c-8e00-0c9e04dd5504','application','update',NULL,NULL,'应用暂停权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','project','应用克隆','70c2228f-45d9-4ca9-98e5-e587b9249711','application','create',NULL,NULL,'应用克隆权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564','project','应用发布','4f52c57e-17c7-463d-b3a0-79d319c5f860','application','create',NULL,NULL,'应用发布权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e1dbef1f-22e8-4b53-9e20-80a83c6b697b','project','工作流管理','fd6f3c5e-748a-4b8a-beb1-dc73fd3e3ff8','workflow','*',NULL,NULL,'工作流管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('fd6f3c5e-748a-4b8a-beb1-dc73fd3e3ff8','project','工作流查询','4ae783ba-05f9-419c-acf3-d2efa7fab21a','workflow','query',NULL,NULL,'工作流查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('fd6f3c5e-748a-4b8a-beb1-dc73fd3e3ff8','project','工作流创建','626a0c9c-93d7-438c-91dc-e1d5624c038c','workflow','create',NULL,NULL,'工作流创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('fd6f3c5e-748a-4b8a-beb1-dc73fd3e3ff8','project','工作流修改','44a0c82d-db53-40e0-9aa4-678f3607e5f6','workflow','update',NULL,NULL,'工作流修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('fd6f3c5e-748a-4b8a-beb1-dc73fd3e3ff8','project','工作流删除','00d7960a-df41-413d-87be-9dd0d02fcd59','workflow','delete',NULL,NULL,'工作流删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('fd6f3c5e-748a-4b8a-beb1-dc73fd3e3ff8','project','工作流执行','3c2d5b0d-b95d-49b4-943c-05fd7307c00e','workflow','create',NULL,NULL,'工作流执行权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('fd6f3c5e-748a-4b8a-beb1-dc73fd3e3ff8','project','工作流上线','579d2a25-643b-494b-8fb7-526c78b48425','workflow','update',NULL,NULL,'工作流上线权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('fd6f3c5e-748a-4b8a-beb1-dc73fd3e3ff8','project','工作流下线','c91054e7-dd09-4056-83f2-5c2d4e274ed9','workflow','update',NULL,NULL,'工作流下线权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e1dbef1f-22e8-4b53-9e20-80a83c6b697b','project','任务管理','cfda7904-5e54-4049-8bcf-261327b95194','job','*',NULL,NULL,'任务管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cfda7904-5e54-4049-8bcf-261327b95194','project','任务查询','9cf8b7d4-1c71-464e-86ce-af0281cb659d','job','query',NULL,NULL,'任务查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cfda7904-5e54-4049-8bcf-261327b95194','project','任务创建','21055d16-de3b-48ec-9ebf-694a52c994c9','job','create',NULL,NULL,'任务创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cfda7904-5e54-4049-8bcf-261327b95194','project','任务修改','70231840-1cea-4804-adb8-43cb02a9f346','job','update',NULL,NULL,'任务修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cfda7904-5e54-4049-8bcf-261327b95194','project','任务删除','15395549-aa2b-4797-8339-37170cacff8b','job','delete',NULL,NULL,'任务删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cfda7904-5e54-4049-8bcf-261327b95194','project','任务执行','d7b0d62e-518d-422c-9b92-b89d838c3468','job','create',NULL,NULL,'任务执行权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cfda7904-5e54-4049-8bcf-261327b95194','project','任务启用','d43b3275-84e1-45a1-966b-216d423a5022','job','update',NULL,NULL,'任务启用权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cfda7904-5e54-4049-8bcf-261327b95194','project','任务禁用','8ea00e66-107b-4781-ad97-5161d47595dc','job','update',NULL,NULL,'任务禁用权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e1dbef1f-22e8-4b53-9e20-80a83c6b697b','project','资源','b8535fc4-8da4-4184-b184-55a3090e45fb','project_resource','*',NULL,NULL,'项目资源管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b8535fc4-8da4-4184-b184-55a3090e45fb','project','资源组管理','c1bae555-b464-444a-8f9c-5fc7a133b674','resource_group','*',NULL,NULL,'资源组管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('c1bae555-b464-444a-8f9c-5fc7a133b674','project','资源组查询','6d9e764a-49e0-42f5-9fb8-3efc2cd40e70','resource_group','query',NULL,NULL,'资源组查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('c1bae555-b464-444a-8f9c-5fc7a133b674','project','资源组创建','2c844549-d528-4e05-9e22-39f920bb5bca','resource_group','create',NULL,NULL,'资源组创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('c1bae555-b464-444a-8f9c-5fc7a133b674','project','资源组修改','3415c393-d433-4fed-b958-fa3581114669','resource_group','update',NULL,NULL,'资源组修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('c1bae555-b464-444a-8f9c-5fc7a133b674','project','资源组删除','7beabed6-c57d-4185-8edf-ea64c098d53f','resource_group','delete',NULL,NULL,'资源组删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b8535fc4-8da4-4184-b184-55a3090e45fb','project','服务器管理','b943ea96-4d7d-4a68-bf15-02a1ddded533','server','*',NULL,NULL,'服务器管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b943ea96-4d7d-4a68-bf15-02a1ddded533','project','服务器查询','bbc990d0-b47e-4ddd-9d60-6941e6498386','server','query',NULL,NULL,'服务器查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b943ea96-4d7d-4a68-bf15-02a1ddded533','project','服务器创建','d1f06bce-29c6-4419-8a07-7738d9825f1a','server','create',NULL,NULL,'服务器创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b943ea96-4d7d-4a68-bf15-02a1ddded533','project','服务器修改','0af412da-f92f-4c98-8e87-ee6c5558b8c1','server','update',NULL,NULL,'服务器修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b943ea96-4d7d-4a68-bf15-02a1ddded533','project','服务器删除','2e42c632-7c91-4e49-acde-70b188da7d1c','server','delete',NULL,NULL,'服务器删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b943ea96-4d7d-4a68-bf15-02a1ddded533','project','服务器启动','711e966e-8a9d-4a13-afec-03776f0f4ebc','server','update',NULL,NULL,'服务器启动权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b943ea96-4d7d-4a68-bf15-02a1ddded533','project','服务器停止','c7580746-6e75-4fc2-a923-e2f6157d5c1d','server','update',NULL,NULL,'服务器停止权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b943ea96-4d7d-4a68-bf15-02a1ddded533','project','服务器重启','03ac0ca1-e756-4001-bde0-39c17b59281d','server','update',NULL,NULL,'服务器重启权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b943ea96-4d7d-4a68-bf15-02a1ddded533','project','服务器关机','adb4e670-8522-4840-a75e-bc2a9643c11f','server','update',NULL,NULL,'服务器关机权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b943ea96-4d7d-4a68-bf15-02a1ddded533','project','服务器终端','a0a528ec-6938-48cf-8307-8575da72ebc3','server','query',NULL,NULL,'服务器终端权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b943ea96-4d7d-4a68-bf15-02a1ddded533','project','服务器文件','64618046-70f3-4971-aeef-bff6313669c8','server','create',NULL,NULL,'服务器文件管理权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b943ea96-4d7d-4a68-bf15-02a1ddded533','project','服务器系统服务','6ab395aa-2ff9-44ed-aa5b-2f5c09d44475','server','query',NULL,NULL,'服务器系统服务权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b8535fc4-8da4-4184-b184-55a3090e45fb','project','密钥管理','9606f68e-7aa7-452e-a3af-2d21e243971b','secret','*',NULL,NULL,'密钥管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('9606f68e-7aa7-452e-a3af-2d21e243971b','project','密钥查询','651efc85-4833-46b8-9d41-12158b6988c6','secret','query',NULL,NULL,'密钥查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('9606f68e-7aa7-452e-a3af-2d21e243971b','project','密钥创建','461d7acc-d38f-4927-a74a-d2232b1cb9d8','secret','create',NULL,NULL,'密钥创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('9606f68e-7aa7-452e-a3af-2d21e243971b','project','密钥修改','5efcc84d-1f23-45f6-83ef-6cca8498b456','secret','update',NULL,NULL,'密钥修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('9606f68e-7aa7-452e-a3af-2d21e243971b','project','密钥删除','1cb1d050-9623-498d-abcd-8d9c4b88bc55','secret','delete',NULL,NULL,'密钥删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b8535fc4-8da4-4184-b184-55a3090e45fb','project','数据库管理','218f52bf-e23b-42ed-9e52-11a4685f7574','database','*',NULL,NULL,'数据库管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('218f52bf-e23b-42ed-9e52-11a4685f7574','project','数据库查询','019ffd00-5a4a-48a0-b939-694ac139e91d','database','query',NULL,NULL,'数据库查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('218f52bf-e23b-42ed-9e52-11a4685f7574','project','数据库创建','2db086d7-aaea-4eeb-99f6-2753ad69e604','database','create',NULL,NULL,'数据库创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('218f52bf-e23b-42ed-9e52-11a4685f7574','project','数据库修改','f4704220-3059-4684-a861-8845d5344f01','database','update',NULL,NULL,'数据库修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('218f52bf-e23b-42ed-9e52-11a4685f7574','project','数据库删除','f4070c8a-8bcc-49ef-b480-b91424a79aec','database','delete',NULL,NULL,'数据库删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b8535fc4-8da4-4184-b184-55a3090e45fb','project','应用网关管理','cc20089c-8b1a-43c9-8862-9a1ce00ab629','gateway','*',NULL,NULL,'应用网关管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cc20089c-8b1a-43c9-8862-9a1ce00ab629','project','应用网关查询','2372a133-912e-4835-ba53-e7cc27912052','gateway','query',NULL,NULL,'应用网关查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cc20089c-8b1a-43c9-8862-9a1ce00ab629','project','应用网关创建','d4c26f66-b9f5-49a7-972e-530bd96a95ba','gateway','create',NULL,NULL,'应用网关创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cc20089c-8b1a-43c9-8862-9a1ce00ab629','project','应用网关修改','837ccd3e-9c94-4348-9ec1-7c512da11acf','gateway','update',NULL,NULL,'应用网关修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cc20089c-8b1a-43c9-8862-9a1ce00ab629','project','应用网关删除','53872b5c-36e4-4e50-96a6-ab1bcaf75761','gateway','delete',NULL,NULL,'应用网关删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cc20089c-8b1a-43c9-8862-9a1ce00ab629','project','应用网关启动','69815cb0-cf91-47be-ab14-d01f932cfc88','gateway','update',NULL,NULL,'应用网关启动权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cc20089c-8b1a-43c9-8862-9a1ce00ab629','project','应用网关停止','656e6f83-0bde-4f7c-81ae-f6171827a2b5','gateway','update',NULL,NULL,'应用网关停止权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b8535fc4-8da4-4184-b184-55a3090e45fb','project','证书管理','7d833f2a-521c-4cf2-83de-f76284dc229a','certificate','*',NULL,NULL,'证书管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('7d833f2a-521c-4cf2-83de-f76284dc229a','project','证书查询','d04fa821-c5f0-49db-b0ff-160cba0f26a6','certificate','query',NULL,NULL,'证书查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('7d833f2a-521c-4cf2-83de-f76284dc229a','project','证书创建','76adb0bc-472d-4da6-900a-72443ca925eb','certificate','create',NULL,NULL,'证书创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('7d833f2a-521c-4cf2-83de-f76284dc229a','project','证书修改','9caeef89-180c-4488-9179-31b7495bb08b','certificate','update',NULL,NULL,'证书修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('7d833f2a-521c-4cf2-83de-f76284dc229a','project','证书删除','3cb9563f-eadb-4d22-9312-84531e11d803','certificate','delete',NULL,NULL,'证书删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('b8535fc4-8da4-4184-b184-55a3090e45fb','project','云资源管理','5a09b3cf-aa47-4f72-ae2c-8c8216474119','cloud_resource','*',NULL,NULL,'云资源管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a09b3cf-aa47-4f72-ae2c-8c8216474119','project','云资源查询','b7fbc10d-1c5e-416f-ae45-909abfb336a3','cloud_resource','query',NULL,NULL,'云资源查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a09b3cf-aa47-4f72-ae2c-8c8216474119','project','云资源创建','8f932526-2fdc-409c-a923-13f1cd23cec2','cloud_resource','create',NULL,NULL,'云资源创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a09b3cf-aa47-4f72-ae2c-8c8216474119','project','云资源修改','38c1bdfc-df30-4a46-b67d-eba5c03f1279','cloud_resource','update',NULL,NULL,'云资源修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('5a09b3cf-aa47-4f72-ae2c-8c8216474119','project','云资源删除','63451d9c-9d64-47c7-9df5-7a7b91bdeee3','cloud_resource','delete',NULL,NULL,'云资源删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e1dbef1f-22e8-4b53-9e20-80a83c6b697b','project','项目团队管理','78b516f8-69e5-432b-b4cd-45ce29306ee5','project_team','*',NULL,NULL,'项目团队管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('78b516f8-69e5-432b-b4cd-45ce29306ee5','project','项目团队查询','f435ce92-d2ca-4974-b405-d19240f8efea','project_team','query',NULL,NULL,'项目团队查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('78b516f8-69e5-432b-b4cd-45ce29306ee5','project','项目团队创建','956b258f-67a1-4d69-96c8-7aa9e2bc91ff','project_team','create',NULL,NULL,'项目团队创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('78b516f8-69e5-432b-b4cd-45ce29306ee5','project','项目团队修改','8541ef92-1adf-4436-9e38-92679a6a692e','project_team','update',NULL,NULL,'项目团队修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('78b516f8-69e5-432b-b4cd-45ce29306ee5','project','项目团队删除','5b0866ee-2b5e-4af7-8983-377bffe40a8c','project_team','delete',NULL,NULL,'项目团队删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e1dbef1f-22e8-4b53-9e20-80a83c6b697b','project','项目设置','72eed345-8696-418b-a11c-6a5c1ba6c25d','project_setting','*',NULL,NULL,'项目设置全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('72eed345-8696-418b-a11c-6a5c1ba6c25d','project','项目设置查询','aea39fc1-e9e2-48b2-8004-889ceecca98b','project_setting','query',NULL,NULL,'项目设置查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('72eed345-8696-418b-a11c-6a5c1ba6c25d','project','项目设置修改','3e653cd0-fc11-4284-9847-cd1180f0b607','project_setting','update',NULL,NULL,'项目设置修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e06e84c7-fe48-4d0d-9dbb-cbe65cf5d936','platform','应用','523386f8-aba4-45a2-9f6e-45dc2d4de481','apps','*',NULL,NULL,'应用',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('523386f8-aba4-45a2-9f6e-45dc2d4de481','platform','应用市场','493b3dca-fe70-44e9-a57e-56247b98ae02','marketplace','*',NULL,NULL,'应用市场全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('493b3dca-fe70-44e9-a57e-56247b98ae02','platform','应用市场查询','17aadb85-a87b-4dfa-b0c9-1a6e58ed36dc','marketplace','query',NULL,NULL,'应用市场查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('493b3dca-fe70-44e9-a57e-56247b98ae02','platform','应用市场创建','77bdea9e-ef36-4e0b-8f34-0a85535ee6e1','marketplace','create',NULL,NULL,'应用市场创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('493b3dca-fe70-44e9-a57e-56247b98ae02','platform','应用市场修改','e15d1dce-0056-43f3-82f6-81f5c85fc516','marketplace','update',NULL,NULL,'应用市场修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('493b3dca-fe70-44e9-a57e-56247b98ae02','platform','应用市场删除','9596604b-8c1b-4b1e-90a6-3b83f8b55ec0','marketplace','delete',NULL,NULL,'应用市场删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('493b3dca-fe70-44e9-a57e-56247b98ae02','platform','应用市场下载','c59de836-3de9-4494-98ed-b8e74d698cb5','marketplace','query',NULL,NULL,'应用市场下载权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('523386f8-aba4-45a2-9f6e-45dc2d4de481','platform','应用心愿单','d705e552-6684-4651-a4dd-066e0fe4c044','wishlist','*',NULL,NULL,'应用心愿单全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('d705e552-6684-4651-a4dd-066e0fe4c044','platform','应用心愿单查询','225f4393-804a-4453-b83a-c1dd0c476a94','wishlist','query',NULL,NULL,'应用心愿单查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('d705e552-6684-4651-a4dd-066e0fe4c044','platform','应用心愿单创建','520bb09c-3ed5-43c3-a188-f3d6d19279a9','wishlist','create',NULL,NULL,'应用心愿单创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('d705e552-6684-4651-a4dd-066e0fe4c044','platform','应用心愿单修改','ca0fa161-5578-4225-8bb6-0ad29f3149ee','wishlist','update',NULL,NULL,'应用心愿单修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('d705e552-6684-4651-a4dd-066e0fe4c044','platform','应用心愿单删除','cb76753a-6452-4c9a-8583-26571829148a','wishlist','delete',NULL,NULL,'应用心愿单删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e06e84c7-fe48-4d0d-9dbb-cbe65cf5d936','platform','管理员设置','dee24e8d-d29a-474d-b889-9ec4a8296a04','admin_setting','*',NULL,NULL,'管理员设置',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('dee24e8d-d29a-474d-b889-9ec4a8296a04','platform','管理员设置查询','7c366f6c-b342-449b-97e3-f94529c5b105','admin_setting','query',NULL,NULL,'管理员设置查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('dee24e8d-d29a-474d-b889-9ec4a8296a04','platform','管理员设置修改','3fa5ca8a-174a-44ab-872b-d6df77ca0d65','admin_setting','update',NULL,NULL,'管理员设置修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('e06e84c7-fe48-4d0d-9dbb-cbe65cf5d936','platform','平台管理','777614c1-7f84-4e3e-95ee-edab4f5a47bc','platform_setting','*',NULL,NULL,'平台管理',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('777614c1-7f84-4e3e-95ee-edab4f5a47bc','platform','项目管理','30e70df2-6701-43d4-ae9e-f869ac63e50f','project_manage','*',NULL,NULL,'项目管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('30e70df2-6701-43d4-ae9e-f869ac63e50f','platform','项目管理查询','1fd95b32-2d00-44e2-89fc-1dc6eaa8062a','project_manage','query',NULL,NULL,'项目管理查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('30e70df2-6701-43d4-ae9e-f869ac63e50f','platform','项目管理创建','8d1b0d1b-9b66-49ea-bf09-ddd595a57b85','project_manage','create',NULL,NULL,'项目管理创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('30e70df2-6701-43d4-ae9e-f869ac63e50f','platform','项目管理修改','9a20ca40-408b-4f57-a1a6-1c1d3a53a845','project_manage','update',NULL,NULL,'项目管理修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('30e70df2-6701-43d4-ae9e-f869ac63e50f','platform','项目管理删除','c25194e2-0262-43e8-a7d6-49b14764edbf','project_manage','delete',NULL,NULL,'项目管理删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('30e70df2-6701-43d4-ae9e-f869ac63e50f','platform','项目管理归档','324b7674-3f44-4116-942d-d219040c9c4e','project_manage','create',NULL,NULL,'项目管理归档权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('30e70df2-6701-43d4-ae9e-f869ac63e50f','platform','项目管理恢复','790f1c9c-127a-4b2a-97e9-d98101bf36b5','project_manage','create',NULL,NULL,'项目管理恢复权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('777614c1-7f84-4e3e-95ee-edab4f5a47bc','platform','安全管理','cba5a624-afa7-4fff-a1f1-f6f791e5589a','security','*',NULL,NULL,'安全管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cba5a624-afa7-4fff-a1f1-f6f791e5589a','platform','角色管理','f9ea8516-3350-44cc-adb5-10cf0e5f8e70','role','*','/roles',NULL,'角色管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('f9ea8516-3350-44cc-adb5-10cf0e5f8e70','platform','角色管理查询','3493ff86-ea83-4acd-99d5-c82e0da36b69','role','query','/roles',NULL,'角色管理查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('f9ea8516-3350-44cc-adb5-10cf0e5f8e70','platform','角色管理创建','f99358ba-b8b6-4408-a1c8-c805d917fb6f','role','create','/roles',NULL,'角色管理创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('f9ea8516-3350-44cc-adb5-10cf0e5f8e70','platform','角色管理修改','30c36491-1e77-416b-b056-f860821dd093','role','update','/roles/*',NULL,'角色管理修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('f9ea8516-3350-44cc-adb5-10cf0e5f8e70','platform','角色管理删除','5c9cc887-2d4c-47d7-99cd-ee55056b2bac','role','delete','/roles/*',NULL,'角色管理删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cba5a624-afa7-4fff-a1f1-f6f791e5589a','platform','权限管理','1de209de-9182-4514-8016-4fc3401b7968','permission','*','/permissions',NULL,'权限管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1de209de-9182-4514-8016-4fc3401b7968','platform','权限管理查询','095cba59-25e0-4d6d-8739-4c537ade0b71','permission','query','/permissions',NULL,'权限管理查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1de209de-9182-4514-8016-4fc3401b7968','platform','权限管理创建','4db0477f-ef96-40f3-9853-3cb2df4b0901','permission','create','/permissions',NULL,'权限管理创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1de209de-9182-4514-8016-4fc3401b7968','platform','权限管理修改','cfcf331f-a33d-4521-a7f7-f0fbece23fec','permission','update','/permissions/*',NULL,'权限管理修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1de209de-9182-4514-8016-4fc3401b7968','platform','权限管理删除','9dbd8b47-d475-4a89-b619-3ea8e087fea6','permission','delete','/permissions/*',NULL,'权限管理删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('cba5a624-afa7-4fff-a1f1-f6f791e5589a','platform','认证管理','a5b06786-9af8-4e8e-bd0c-5c4bf6f74dc1','auth','*','/auth-config',NULL,'认证管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('a5b06786-9af8-4e8e-bd0c-5c4bf6f74dc1','platform','认证管理查询','ff1f1828-c355-45fa-8b97-e3ebaeeae022','auth','query','/auth-config',NULL,'认证管理查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('a5b06786-9af8-4e8e-bd0c-5c4bf6f74dc1','platform','认证管理创建','b15817ce-73d4-40d4-924f-b4cc93e15d12','auth','create','/auth-config',NULL,'认证管理创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('a5b06786-9af8-4e8e-bd0c-5c4bf6f74dc1','platform','认证管理修改','64d5ec25-2510-4ad2-8988-4dc326c92f2e','auth','update','/auth-config',NULL,'认证管理修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('a5b06786-9af8-4e8e-bd0c-5c4bf6f74dc1','platform','认证管理删除','8b518edc-6ee0-4d3d-8af7-da30e901e04d','auth','delete','/auth-config',NULL,'认证管理删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('777614c1-7f84-4e3e-95ee-edab4f5a47bc','platform','用户管理','1adb3cfc-ddd8-42ab-9c0f-3b2881184852','user','*','/users',NULL,'用户管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1adb3cfc-ddd8-42ab-9c0f-3b2881184852','platform','用户管理查询','a4e88f35-af15-4fe0-b30b-ec2f516e9fc3','user','query','/users',NULL,'用户管理查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1adb3cfc-ddd8-42ab-9c0f-3b2881184852','platform','用户管理创建','01aaa408-953e-4c5a-8cb0-b507322a4205','user','create','/users',NULL,'用户管理创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1adb3cfc-ddd8-42ab-9c0f-3b2881184852','platform','用户管理修改','a176c0fc-24ba-42f6-bb8f-1f80a8a2fe29','user','update','/users/*',NULL,'用户管理修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1adb3cfc-ddd8-42ab-9c0f-3b2881184852','platform','用户管理删除','9644b60c-b8ea-4351-b672-e7b656bb33c2','user','delete','/users/*',NULL,'用户管理删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1adb3cfc-ddd8-42ab-9c0f-3b2881184852','platform','用户管理启用','167cfa0f-89ac-44a6-ab9c-fac2ec13ff62','user','update','/users/*/status',NULL,'用户管理启用权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1adb3cfc-ddd8-42ab-9c0f-3b2881184852','platform','用户管理禁用','d1e3a654-78a9-4e78-905c-92f44e13ebc6','user','update','/users/*/status',NULL,'用户管理禁用权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('777614c1-7f84-4e3e-95ee-edab4f5a47bc','platform','告警通知','ef4eab25-7246-43b2-aa11-5a5863d51de2','notification','*',NULL,NULL,'告警通知管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('ef4eab25-7246-43b2-aa11-5a5863d51de2','platform','告警通知查询','021d7475-a9c6-4b06-9178-1e86e89cab88','notification','query',NULL,NULL,'告警通知查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('ef4eab25-7246-43b2-aa11-5a5863d51de2','platform','告警通知创建','f2b881a4-a53c-49ed-a980-132854d7dfd1','notification','create',NULL,NULL,'告警通知创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('ef4eab25-7246-43b2-aa11-5a5863d51de2','platform','告警通知修改','d92c16b0-9a5b-4052-86ca-2973708add9f','notification','update',NULL,NULL,'告警通知修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('ef4eab25-7246-43b2-aa11-5a5863d51de2','platform','告警通知删除','51aa6a7d-66bc-4551-b04b-7a4bdf8e9b3e','notification','delete',NULL,NULL,'告警通知删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('777614c1-7f84-4e3e-95ee-edab4f5a47bc','platform','个人中心','208ebe4c-74f3-48ea-9819-9aa9a8212158','profile','*',NULL,NULL,'个人中心全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('208ebe4c-74f3-48ea-9819-9aa9a8212158','platform','个人中心查询','1e4b883c-1455-44a0-a61b-c98fa68238c5','profile','query',NULL,NULL,'个人中心查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('208ebe4c-74f3-48ea-9819-9aa9a8212158','platform','个人中心修改','6bb78f8b-bbaa-48ec-b408-1d903415881b','profile','update',NULL,NULL,'个人中心修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('777614c1-7f84-4e3e-95ee-edab4f5a47bc','platform','审计日志','556f99a8-6626-4200-9f80-6abfa68c18e2','audit_log','*','/audit-logs',NULL,'审计日志全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('556f99a8-6626-4200-9f80-6abfa68c18e2','platform','审计日志查询','2a65e848-d451-469d-9d8c-76c7b527cec5','audit_log','query','/audit-logs',NULL,'审计日志查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('556f99a8-6626-4200-9f80-6abfa68c18e2','platform','审计日志导出','e706486b-32a9-4aec-be78-b9d480b571ab','audit_log','query','/audit-logs/export',NULL,'审计日志导出权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19');

-- 插入默认用户及角色权限
INSERT INTO `users` (`id`,`username`,`email`,`password_hash`,`nickname`,`avatar`,`phone`,`gender`,`signature`,`status`,`last_login_at`,`last_login_ip`,`timezone`,`language`,`created_at`,`updated_at`) VALUES
	 (1,'admin','admin@websoft9.com','d1a7b27aa60359a6033b046831168286dddc5a83268d4d483d30a09983f9c944','Manager','','',0,'Websoft9 manager',1,NULL,'','Asia/Shanghai','zh-CN','2025-09-01 12:11:19','2025-09-01 12:11:19');

INSERT INTO `user_roles` (`id`,`user_id`,`role_id`,`granted_by`,`granted_at`,`expires_at`,`status`,`created_at`,`updated_at`) VALUES
	 (1,1,1,1,'2025-09-01 12:11:19',NULL,1,'2025-09-01 12:11:19','2025-09-01 12:11:19');

INSERT INTO `role_permissions` (`role_id`,`permission_code`,`granted_by`,`granted_at`,`status`,`created_at`,`updated_at`) VALUES
         (1, 'e06e84c7-fe48-4d0d-9dbb-cbe65cf5d936',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '16aa0ce0-ce2a-4075-8aee-5f9ccb472819',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'b2fabec8-f7fe-480d-883c-5e0a355b8c69',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'e0aec98a-101f-494a-b5ea-ada5cb5071d1',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'b61077db-d098-42d0-ade9-472038167183',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '91ab578c-1c7e-4f56-9708-4d00f4ad3c2d',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '8422b16f-b03f-4458-9159-3b78e571a88c',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '9f99aee4-fffc-4ad6-b49c-7bf2567352cc',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'a52729f5-7300-456c-a90a-1f7e4bd4d0db',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '3ea52579-58a8-463e-9a46-0dbec048fb07',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'c3393a62-4526-4f8a-bdd9-d6012f6d590f',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '08d5484f-ee91-47ba-99f9-ae90c8d5684d',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'af309229-e1c0-4053-b8c3-e35ef0b89107',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '593b5fb7-4146-4f17-9948-bf2827c4cd94',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '67e63e8c-17b8-49d7-b6cd-6618614d3d5c',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'd8e80678-a03a-42c3-99d8-29ed4b9032ee',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'f47d800e-d504-4de4-946c-60add0123f8d',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '89a1752f-7fda-4949-a241-f5704f713eeb',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'e1dbef1f-22e8-4b53-9e20-80a83c6b697b',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'dee6721a-870a-413c-8b6c-2fb82b153479',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'b0335079-0f4e-49a2-91ed-3946c0f4d13a',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '222dd995-dbb0-4af3-9ecc-ec372a7fe744',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '68e2bbed-27ec-4d6c-95b4-71978b522596',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'df53f526-9ae3-4ea8-a137-dea54e9e0306',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '93d6d1db-b814-43ed-864a-4f485af455bd',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'ecf7b9a6-3a73-4189-a3a5-b00fcc404936',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '049c318d-3593-4998-843b-7abff0e47e17',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '375e6ef7-70d1-4f5d-98ce-affc9732a220',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'f0a85d28-e7b9-45db-8f74-f3263fd3b02a',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '47b481e8-39a8-469f-9d83-a930d1053108',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'd246514e-8b6f-4d5d-b154-123768a656df',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '5a94af8e-0d2b-4c8b-8f64-5fb3f8a18564',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '4ac0b046-dd27-4fb3-acba-f9618237061c',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '53122f15-2c25-4aa6-a2e1-ea3690d8c2f5',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '56beb0c8-783e-4c87-8e48-048bd13b7d31',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'a01c53a1-c58e-4bc3-8c82-6228a358193f',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'ea61a581-62c5-4206-a48a-58c6d468c6bb',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '7b1e1414-4ca6-46f1-bc50-4d94c9cebd02',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '6b68510c-bfdf-4c04-897f-9a15d080a8ed',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'c420843c-3984-4041-b990-a7bd19a91271',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '0af2937b-d661-427b-93ed-7f1c769cfdda',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '8a8db5d8-46b3-4b0c-8e00-0c9e04dd5504',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '70c2228f-45d9-4ca9-98e5-e587b9249711',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '4f52c57e-17c7-463d-b3a0-79d319c5f860',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'fd6f3c5e-748a-4b8a-beb1-dc73fd3e3ff8',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '4ae783ba-05f9-419c-acf3-d2efa7fab21a',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '626a0c9c-93d7-438c-91dc-e1d5624c038c',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '44a0c82d-db53-40e0-9aa4-678f3607e5f6',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '00d7960a-df41-413d-87be-9dd0d02fcd59',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '3c2d5b0d-b95d-49b4-943c-05fd7307c00e',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '579d2a25-643b-494b-8fb7-526c78b48425',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'c91054e7-dd09-4056-83f2-5c2d4e274ed9',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'cfda7904-5e54-4049-8bcf-261327b95194',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '9cf8b7d4-1c71-464e-86ce-af0281cb659d',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '21055d16-de3b-48ec-9ebf-694a52c994c9',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '70231840-1cea-4804-adb8-43cb02a9f346',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '15395549-aa2b-4797-8339-37170cacff8b',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'd7b0d62e-518d-422c-9b92-b89d838c3468',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'd43b3275-84e1-45a1-966b-216d423a5022',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '8ea00e66-107b-4781-ad97-5161d47595dc',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'b8535fc4-8da4-4184-b184-55a3090e45fb',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'c1bae555-b464-444a-8f9c-5fc7a133b674',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '6d9e764a-49e0-42f5-9fb8-3efc2cd40e70',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '2c844549-d528-4e05-9e22-39f920bb5bca',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '3415c393-d433-4fed-b958-fa3581114669',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '7beabed6-c57d-4185-8edf-ea64c098d53f',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'b943ea96-4d7d-4a68-bf15-02a1ddded533',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'bbc990d0-b47e-4ddd-9d60-6941e6498386',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'd1f06bce-29c6-4419-8a07-7738d9825f1a',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '0af412da-f92f-4c98-8e87-ee6c5558b8c1',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '2e42c632-7c91-4e49-acde-70b188da7d1c',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '711e966e-8a9d-4a13-afec-03776f0f4ebc',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'c7580746-6e75-4fc2-a923-e2f6157d5c1d',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '03ac0ca1-e756-4001-bde0-39c17b59281d',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'adb4e670-8522-4840-a75e-bc2a9643c11f',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'a0a528ec-6938-48cf-8307-8575da72ebc3',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '64618046-70f3-4971-aeef-bff6313669c8',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '6ab395aa-2ff9-44ed-aa5b-2f5c09d44475',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '9606f68e-7aa7-452e-a3af-2d21e243971b',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '651efc85-4833-46b8-9d41-12158b6988c6',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '461d7acc-d38f-4927-a74a-d2232b1cb9d8',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '5efcc84d-1f23-45f6-83ef-6cca8498b456',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '1cb1d050-9623-498d-abcd-8d9c4b88bc55',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '218f52bf-e23b-42ed-9e52-11a4685f7574',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '019ffd00-5a4a-48a0-b939-694ac139e91d',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '2db086d7-aaea-4eeb-99f6-2753ad69e604',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'f4704220-3059-4684-a861-8845d5344f01',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'f4070c8a-8bcc-49ef-b480-b91424a79aec',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'cc20089c-8b1a-43c9-8862-9a1ce00ab629',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '2372a133-912e-4835-ba53-e7cc27912052',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'd4c26f66-b9f5-49a7-972e-530bd96a95ba',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '837ccd3e-9c94-4348-9ec1-7c512da11acf',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '53872b5c-36e4-4e50-96a6-ab1bcaf75761',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '69815cb0-cf91-47be-ab14-d01f932cfc88',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '656e6f83-0bde-4f7c-81ae-f6171827a2b5',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '7d833f2a-521c-4cf2-83de-f76284dc229a',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'd04fa821-c5f0-49db-b0ff-160cba0f26a6',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '76adb0bc-472d-4da6-900a-72443ca925eb',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '9caeef89-180c-4488-9179-31b7495bb08b',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '3cb9563f-eadb-4d22-9312-84531e11d803',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '5a09b3cf-aa47-4f72-ae2c-8c8216474119',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'b7fbc10d-1c5e-416f-ae45-909abfb336a3',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '8f932526-2fdc-409c-a923-13f1cd23cec2',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '38c1bdfc-df30-4a46-b67d-eba5c03f1279',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '63451d9c-9d64-47c7-9df5-7a7b91bdeee3',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '78b516f8-69e5-432b-b4cd-45ce29306ee5',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'f435ce92-d2ca-4974-b405-d19240f8efea',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '956b258f-67a1-4d69-96c8-7aa9e2bc91ff',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '8541ef92-1adf-4436-9e38-92679a6a692e',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '5b0866ee-2b5e-4af7-8983-377bffe40a8c',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '72eed345-8696-418b-a11c-6a5c1ba6c25d',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'aea39fc1-e9e2-48b2-8004-889ceecca98b',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '3e653cd0-fc11-4284-9847-cd1180f0b607',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '523386f8-aba4-45a2-9f6e-45dc2d4de481',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '493b3dca-fe70-44e9-a57e-56247b98ae02',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '17aadb85-a87b-4dfa-b0c9-1a6e58ed36dc',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '77bdea9e-ef36-4e0b-8f34-0a85535ee6e1',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'e15d1dce-0056-43f3-82f6-81f5c85fc516',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '9596604b-8c1b-4b1e-90a6-3b83f8b55ec0',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'c59de836-3de9-4494-98ed-b8e74d698cb5',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'd705e552-6684-4651-a4dd-066e0fe4c044',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '225f4393-804a-4453-b83a-c1dd0c476a94',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '520bb09c-3ed5-43c3-a188-f3d6d19279a9',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'ca0fa161-5578-4225-8bb6-0ad29f3149ee',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'cb76753a-6452-4c9a-8583-26571829148a',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'dee24e8d-d29a-474d-b889-9ec4a8296a04',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '7c366f6c-b342-449b-97e3-f94529c5b105',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '3fa5ca8a-174a-44ab-872b-d6df77ca0d65',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '777614c1-7f84-4e3e-95ee-edab4f5a47bc',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '30e70df2-6701-43d4-ae9e-f869ac63e50f',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '1fd95b32-2d00-44e2-89fc-1dc6eaa8062a',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '8d1b0d1b-9b66-49ea-bf09-ddd595a57b85',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '9a20ca40-408b-4f57-a1a6-1c1d3a53a845',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'c25194e2-0262-43e8-a7d6-49b14764edbf',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '324b7674-3f44-4116-942d-d219040c9c4e',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '790f1c9c-127a-4b2a-97e9-d98101bf36b5',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'cba5a624-afa7-4fff-a1f1-f6f791e5589a',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'f9ea8516-3350-44cc-adb5-10cf0e5f8e70',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '3493ff86-ea83-4acd-99d5-c82e0da36b69',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'f99358ba-b8b6-4408-a1c8-c805d917fb6f',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '30c36491-1e77-416b-b056-f860821dd093',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '5c9cc887-2d4c-47d7-99cd-ee55056b2bac',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '1de209de-9182-4514-8016-4fc3401b7968',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '095cba59-25e0-4d6d-8739-4c537ade0b71',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '4db0477f-ef96-40f3-9853-3cb2df4b0901',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'cfcf331f-a33d-4521-a7f7-f0fbece23fec',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '9dbd8b47-d475-4a89-b619-3ea8e087fea6',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'a5b06786-9af8-4e8e-bd0c-5c4bf6f74dc1',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'ff1f1828-c355-45fa-8b97-e3ebaeeae022',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'b15817ce-73d4-40d4-924f-b4cc93e15d12',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '64d5ec25-2510-4ad2-8988-4dc326c92f2e',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '8b518edc-6ee0-4d3d-8af7-da30e901e04d',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '1adb3cfc-ddd8-42ab-9c0f-3b2881184852',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'a4e88f35-af15-4fe0-b30b-ec2f516e9fc3',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '01aaa408-953e-4c5a-8cb0-b507322a4205',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'a176c0fc-24ba-42f6-bb8f-1f80a8a2fe29',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '9644b60c-b8ea-4351-b672-e7b656bb33c2',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '167cfa0f-89ac-44a6-ab9c-fac2ec13ff62',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'd1e3a654-78a9-4e78-905c-92f44e13ebc6',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'ef4eab25-7246-43b2-aa11-5a5863d51de2',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '021d7475-a9c6-4b06-9178-1e86e89cab88',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'f2b881a4-a53c-49ed-a980-132854d7dfd1',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'd92c16b0-9a5b-4052-86ca-2973708add9f',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '51aa6a7d-66bc-4551-b04b-7a4bdf8e9b3e',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '208ebe4c-74f3-48ea-9819-9aa9a8212158',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '1e4b883c-1455-44a0-a61b-c98fa68238c5',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '6bb78f8b-bbaa-48ec-b408-1d903415881b',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '556f99a8-6626-4200-9f80-6abfa68c18e2',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '2a65e848-d451-469d-9d8c-76c7b527cec5',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'e706486b-32a9-4aec-be78-b9d480b571ab',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19');

-- 插入默认应用分类
INSERT INTO `app_store_categories` (`name`, `code`, `description`, `sort_order`, `status`) VALUES
('Web服务', 'web', 'Web服务器和相关应用', 1, 1),
('数据库', 'database', '各种数据库系统', 2, 1),
('开发工具', 'development', '开发和构建工具', 3, 1),
('监控工具', 'monitoring', '系统监控和日志工具', 4, 1),
('安全工具', 'security', '安全防护工具', 5, 1),
('存储服务', 'storage', '文件存储和对象存储服务', 6, 1),
('消息队列', 'message_queue', '消息队列和流处理服务', 7, 1),
('容器编排', 'orchestration', '容器编排和管理工具', 8, 1),
('其他', 'others', '其他应用', 99, 1);

-- 插入系统配置
INSERT INTO `system_configs` (`config_key`, `config_value`, `config_type`, `category`, `description`, `is_readonly`, `sort_order`) VALUES
-- 基础配置
('system.name', 'Websoft9', 'STRING', 'basic', '系统名称', 0, 1),
('system.version', '1.1.0', 'STRING', 'basic', '系统版本', 1, 2),
('system.timezone', 'Asia/Shanghai', 'STRING', 'basic', '系统时区', 0, 3),
('system.language', 'zh-CN', 'STRING', 'basic', '系统默认语言', 0, 4),
('system.logo', '', 'STRING', 'basic', '系统Logo URL', 0, 5),
('system.favicon', '', 'STRING', 'basic', '系统Favicon URL', 0, 6),

-- 安全配置
('security.password_min_length', '8', 'INTEGER', 'security', '密码最小长度', 0, 10),
('security.password_complexity', 'true', 'BOOLEAN', 'security', '是否启用密码复杂度检查', 0, 11),
('security.session_timeout', '3600', 'INTEGER', 'security', '会话超时时间(秒)', 0, 12),
('security.max_login_attempts', '5', 'INTEGER', 'security', '最大登录尝试次数', 0, 13),
('security.lockout_duration', '300', 'INTEGER', 'security', '账户锁定时间(秒)', 0, 14),
('security.jwt_secret', '', 'STRING', 'security', 'JWT密钥', 1, 15),
('security.jwt_expire_hours', '24', 'INTEGER', 'security', 'JWT过期时间(小时)', 0, 16),

-- 邮件配置
('email.smtp_host', '', 'STRING', 'email', 'SMTP服务器地址', 0, 20),
('email.smtp_port', '587', 'INTEGER', 'email', 'SMTP端口', 0, 21),
('email.smtp_username', '', 'STRING', 'email', 'SMTP用户名', 0, 22),
('email.smtp_password', '', 'STRING', 'email', 'SMTP密码', 1, 23),
('email.smtp_encryption', 'tls', 'STRING', 'email', 'SMTP加密方式', 0, 24),
('email.from_address', '', 'STRING', 'email', '发件人邮箱', 0, 25),
('email.from_name', 'Websoft9', 'STRING', 'email', '发件人名称', 0, 26),

-- 存储配置
('storage.default_driver', 'local', 'STRING', 'storage', '默认存储驱动', 0, 30),
('storage.max_file_size', '100', 'INTEGER', 'storage', '最大文件大小(MB)', 0, 31),
('storage.allowed_extensions', 'jpg,jpeg,png,gif,pdf,doc,docx,xls,xlsx,ppt,pptx,txt,zip,tar,gz', 'STRING', 'storage', '允许的文件扩展名', 0, 32),

-- 监控配置
('monitor.metrics_retention_days', '30', 'INTEGER', 'monitor', '监控数据保留天数', 0, 40),
('monitor.alert_check_interval', '60', 'INTEGER', 'monitor', '告警检查间隔(秒)', 0, 41),
('monitor.default_alert_channels', '["email"]', 'JSON', 'monitor', '默认告警通道', 0, 42);

-- 插入默认通知模板
INSERT INTO `notification_templates` (`name`, `type`, `subject`, `content`, `variables`, `is_system`, `status`) VALUES
('用户注册通知', 'EMAIL', '欢迎注册 {{system_name}}', '亲爱的 {{username}}，\n\n欢迎注册 {{system_name}}！\n\n您的账户已成功创建，现在可以开始使用我们的服务了。\n\n如有任何问题，请联系我们的支持团队。\n\n祝您使用愉快！\n\n{{system_name}} 团队', '["username", "system_name"]', 1, 1),
('密码重置通知', 'EMAIL', '{{system_name}} 密码重置', '亲爱的 {{username}}，\n\n您的密码已成功重置。\n\n如果这不是您的操作，请立即联系我们的支持团队。\n\n{{system_name}} 团队', '["username", "system_name"]', 1, 1),
('系统告警通知', 'EMAIL', '{{system_name}} 系统告警', '告警标题：{{alert_title}}\n告警描述：{{alert_description}}\n触发时间：{{fired_at}}\n告警级别：{{alert_level}}\n\n请及时处理。', '["alert_title", "alert_description", "fired_at", "alert_level", "system_name"]', 1, 1),
('应用部署成功', 'EMAIL', '应用部署成功通知', '亲爱的 {{username}}，\n\n您的应用 {{app_name}} 已成功部署到服务器 {{server_name}}。\n\n访问地址：{{app_url}}\n部署时间：{{deployed_at}}\n\n{{system_name}} 团队', '["username", "app_name", "server_name", "app_url", "deployed_at", "system_name"]', 1, 1),
('应用部署失败', 'EMAIL', '应用部署失败通知', '亲爱的 {{username}}，\n\n您的应用 {{app_name}} 部署失败。\n\n错误信息：{{error_message}}\n失败时间：{{failed_at}}\n\n请检查配置后重试。\n\n{{system_name}} 团队', '["username", "app_name", "error_message", "failed_at", "system_name"]', 1, 1);

-- ========================================
-- 索引优化
-- ========================================

-- 为经常查询的字段添加复合索引
CREATE INDEX `idx_projects_owner_status` ON `projects` (`owner_id`, `status`);
CREATE INDEX `idx_app_instances_project_status` ON `app_instances` (`project_id`, `status`);
CREATE INDEX `idx_app_instances_server_status` ON `app_instances` (`server_id`, `status`);
CREATE INDEX `idx_servers_owner_status` ON `servers` (`owner_id`, `status`);
CREATE INDEX `idx_workflows_project_status` ON `workflows` (`project_id`, `status`);
CREATE INDEX `idx_workflows_owner_status` ON `workflows` (`owner_id`, `status`);
CREATE INDEX `idx_alert_records_rule_status` ON `alert_records` (`alert_rule_id`, `status`);
CREATE INDEX `idx_notifications_target_status` ON `notifications` (`target_type`, `status`);
CREATE INDEX `idx_user_notifications_user_read` ON `user_notifications` (`user_id`, `is_read`);

-- 为时间范围查询添加索引
CREATE INDEX `idx_audit_logs_created_user` ON `audit_logs` (`created_at`, `user_id`);
CREATE INDEX `idx_project_activities_created_project` ON `project_activities` (`created_at`, `project_id`);
CREATE INDEX `idx_workflow_executions_created_status` ON `workflow_executions` (`created_at`, `status`);

-- 为全文搜索添加索引（如果需要）
-- ALTER TABLE `app_store_templates` ADD FULLTEXT(`name`, `description`);
-- ALTER TABLE `app_store_wishlists` ADD FULLTEXT(`name`, `description`);

SET FOREIGN_KEY_CHECKS = 1;
