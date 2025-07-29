-- Websoft9 MySQL Database Initialization Script
-- Generated from Database Design Document

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ========================================
-- 3.1 应用商店 (App Store)
-- ========================================

-- 应用分类表
CREATE TABLE `app_store_categories` (
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
CREATE TABLE `app_store_templates` (
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
CREATE TABLE `app_store_wishlists` (
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
CREATE TABLE `app_store_reviews` (
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
CREATE TABLE `app_store_favorites` (
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
CREATE TABLE `app_store_stars` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT '应用模板ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_template_user` (`template_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_stars_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用点赞表';

-- 应用部署记录表
CREATE TABLE `app_deployments` (
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

-- ========================================
-- 3.2 工作空间 (Workspace)
-- ========================================

-- 用户文件表
CREATE TABLE `user_files` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `name` VARCHAR(255) NOT NULL COMMENT '文件名称',
    `path` VARCHAR(1024) NOT NULL COMMENT '文件路径',
    `type` ENUM('FILE', 'DIRECTORY') NOT NULL DEFAULT 'FILE' COMMENT '类型：FILE, DIRECTORY',
    `size` BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小(字节)',
    `mime_type` VARCHAR(128) NULL COMMENT 'MIME类型',
    `download_count` INT NOT NULL DEFAULT 0 COMMENT '下载次数',
    `parent_id` BIGINT UNSIGNED NULL COMMENT '父目录ID',
    `storage_path` VARCHAR(1024) NULL COMMENT '存储路径',
    `checksum` VARCHAR(64) NULL COMMENT '文件校验和',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_type` (`type`),
    CONSTRAINT `fk_user_files_parent` FOREIGN KEY (`parent_id`) REFERENCES `user_files` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户文件表';

-- 工作流表
CREATE TABLE `workflows` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '工作流名称',
    `code` VARCHAR(32) NOT NULL COMMENT '工作流编码',
    `description` TEXT NULL COMMENT '工作流描述',
    `definition` JSON NOT NULL COMMENT '工作流定义',
    `status` ENUM('DRAFT', 'ACTIVE', 'INACTIVE', 'ARCHIVED') NOT NULL DEFAULT 'DRAFT' COMMENT '状态',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作流表';

-- 工作流任务表
CREATE TABLE `workflow_tasks` (
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
CREATE TABLE `workflow_executions` (
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

-- ========================================
-- 3.3 资源管理 (Resource Management)
-- ========================================

-- 资源组表
CREATE TABLE `resource_groups` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '资源组名称',
    `code` VARCHAR(32) NOT NULL UNIQUE COMMENT '资源组编码',
    `description` TEXT NULL COMMENT '资源组描述',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='资源组表';

-- 数据库连接表
CREATE TABLE `database_connections` (
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
CREATE TABLE `servers` (
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
    `status` ENUM('UNKNOWN', 'RUNNING', 'STOPPED', 'ERROR') NOT NULL DEFAULT 'UNKNOWN' COMMENT '服务器状态',
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
    CONSTRAINT `fk_servers_resource_group` FOREIGN KEY (`resource_group_id`) REFERENCES `resource_groups` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务器表';

-- 客户端表
CREATE TABLE `server_agents` (
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
CREATE TABLE `app_instances` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
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
    `resource_group_id` BIGINT UNSIGNED NULL COMMENT '资源组ID',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_server_id` (`server_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_resource_group_id` (`resource_group_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_app_instances_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_app_instances_server` FOREIGN KEY (`server_id`) REFERENCES `servers` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_app_instances_resource_group` FOREIGN KEY (`resource_group_id`) REFERENCES `resource_groups` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用实例表';

-- ========================================
-- 3.4 安全管控 (Security Management)
-- ========================================

-- SSL证书表
CREATE TABLE `ssl_certificates` (
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
CREATE TABLE `secret_keys` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '密钥名称',
    `key_type` ENUM('API_KEY', 'DATABASE', 'SSH', 'CERTIFICATE', 'CUSTOM') NOT NULL COMMENT '密钥类型',
    `encrypted_value` TEXT NOT NULL COMMENT '加密值',
    `description` TEXT NULL COMMENT '描述',
    `custom_fields` JSON NULL COMMENT '自定义字段',
    `authorized_users` JSON NULL COMMENT '授权用户',
    `authorized_groups` JSON NULL COMMENT '授权用户组',
    `expires_at` DATETIME NULL COMMENT '过期时间',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_key_type` (`key_type`),
    KEY `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='密钥管理表';

-- 应用网关表
CREATE TABLE `app_gateways` (
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
CREATE TABLE `app_gateways_publishes` (
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

-- 审计日志表
CREATE TABLE `audit_logs` (
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
-- 3.5 平台管理 (Platform Management)
-- ========================================

-- 用户组表
CREATE TABLE `user_groups` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '用户组名称',
    `code` VARCHAR(32) NOT NULL UNIQUE COMMENT '用户组编码',
    `description` TEXT NULL COMMENT '用户组描述',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户组表';

-- 用户表
CREATE TABLE `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `group_id` BIGINT UNSIGNED NOT NULL COMMENT '用户组ID',
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
    KEY `idx_group_id` (`group_id`),
    KEY `idx_status` (`status`),
    KEY `idx_email` (`email`),
    CONSTRAINT `fk_users_group` FOREIGN KEY (`group_id`) REFERENCES `user_groups` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 角色表
CREATE TABLE `roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL UNIQUE COMMENT '角色名称',
    `code` VARCHAR(32) NOT NULL UNIQUE COMMENT '角色编码',
    `description` TEXT NULL COMMENT '角色描述',
    `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统角色',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 权限表
CREATE TABLE `permissions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '权限名称',
    `code` VARCHAR(64) NOT NULL UNIQUE COMMENT '权限编码',
    `module` VARCHAR(32) NOT NULL COMMENT '模块名称',
    `action` VARCHAR(32) NOT NULL COMMENT '操作名称',
    `resource` VARCHAR(64) NULL COMMENT '资源标识',
    `description` TEXT NULL COMMENT '权限描述',
    `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统权限',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT '排序',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_module` (`module`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 用户角色关联表
CREATE TABLE `user_roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    `granted_by` BIGINT UNSIGNED NULL COMMENT '授权人ID',
    `granted_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '授权时间',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
    KEY `idx_role_id` (`role_id`),
    KEY `idx_granted_by` (`granted_by`),
    CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 角色权限关联表
CREATE TABLE `role_permissions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    `permission_id` BIGINT UNSIGNED NOT NULL COMMENT '权限ID',
    `granted_by` BIGINT UNSIGNED NULL COMMENT '授权人ID',
    `granted_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '授权时间',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
    KEY `idx_permission_id` (`permission_id`),
    KEY `idx_granted_by` (`granted_by`),
    CONSTRAINT `fk_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- 系统配置表
CREATE TABLE `system_configs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `config_key` VARCHAR(64) NOT NULL UNIQUE COMMENT '配置键',
    `config_value` TEXT NULL COMMENT '配置值',
    `config_type` ENUM('STRING', 'INTEGER', 'BOOLEAN', 'JSON', 'TEXT') NOT NULL DEFAULT 'STRING' COMMENT '配置类型',
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
    KEY `idx_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 告警规则表
CREATE TABLE `alert_rules` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT '规则名称',
    `rule_type` ENUM('THRESHOLD', 'ANOMALY', 'CUSTOM') NOT NULL COMMENT '规则类型',
    `target_type` ENUM('SERVER', 'APPLICATION', 'DATABASE', 'GATEWAY') NOT NULL COMMENT '目标类型',
    `target_id` BIGINT UNSIGNED NULL COMMENT '目标ID',
    `metric_name` VARCHAR(64) NULL COMMENT '指标名称',
    `condition_expression` TEXT NOT NULL COMMENT '条件表达式',
    `notification_channels` JSON NULL COMMENT '通知渠道',
    `is_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '所有者ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_target_type` (`target_type`),
    KEY `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警规则表';

-- 告警记录表
CREATE TABLE `alert_records` (
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
    CONSTRAINT `fk_alert_records_rule` FOREIGN KEY (`alert_rule_id`) REFERENCES `alert_rules` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警记录表';

-- ========================================
-- 初始数据插入
-- ========================================

-- 插入默认用户组
INSERT INTO `user_groups` (`name`, `code`, `description`) VALUES 
('管理员组', 'admin', '系统管理员用户组'),
('普通用户组', 'user', '普通用户组');

-- 插入默认角色
INSERT INTO `roles` (`name`, `code`, `description`, `is_system`) VALUES 
('超级管理员', 'super_admin', '系统超级管理员', 1),
('管理员', 'admin', '系统管理员', 1),
('开发者', 'developer', '开发者角色', 1),
('运维人员', 'operator', '运维人员角色', 1),
('普通用户', 'user', '普通用户角色', 1);

-- 插入默认权限
INSERT INTO `permissions` (`name`, `code`, `module`, `action`, `is_system`) VALUES 
('用户管理', 'user.manage', 'user', 'manage', 1),
('角色管理', 'role.manage', 'role', 'manage', 1),
('权限管理', 'permission.manage', 'permission', 'manage', 1),
('应用管理', 'app.manage', 'app', 'manage', 1),
('服务器管理', 'server.manage', 'server', 'manage', 1),
('监控查看', 'monitor.view', 'monitor', 'view', 1),
('系统配置', 'system.config', 'system', 'config', 1);

-- 插入默认应用分类
INSERT INTO `app_store_categories` (`name`, `code`, `description`) VALUES 
('Web服务', 'web', 'Web服务器和相关应用'),
('数据库', 'database', '各种数据库系统'),
('开发工具', 'development', '开发和构建工具'),
('监控工具', 'monitoring', '系统监控和日志工具'),
('安全工具', 'security', '安全防护工具'),
('其他', 'others', '其他应用');

-- 插入系统配置
INSERT INTO `system_configs` (`config_key`, `config_value`, `config_type`, `category`, `description`) VALUES 
('system.name', 'Websoft9', 'STRING', 'basic', '系统名称'),
('system.version', '1.0.0', 'STRING', 'basic', '系统版本'),
('system.timezone', 'Asia/Shanghai', 'STRING', 'basic', '系统时区'),
('system.language', 'zh-CN', 'STRING', 'basic', '系统语言'),
('security.password_min_length', '6', 'INTEGER', 'security', '密码最小长度'),
('security.session_timeout', '3600', 'INTEGER', 'security', '会话超时时间(秒)');

SET FOREIGN_KEY_CHECKS = 1;