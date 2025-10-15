-- Websoft9 MySQL Database Initialization Script V1.1
-- Generated from Database Design Document V1.1
-- Compatible with MySQL 8.0+

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
SET sql_mode = 'STRICT_TRANS_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO';

-- ========================================
-- 3.1 Platform Home
-- ========================================

-- App shortcut navigation table
CREATE TABLE IF NOT EXISTS `app_shortcuts` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `app_instance_id` BIGINT UNSIGNED NOT NULL COMMENT 'App instance ID',
    `name` VARCHAR(64) NULL COMMENT 'Custom name',
    `description` VARCHAR(255) NULL COMMENT 'Custom description',
    `icon` VARCHAR(255) NULL COMMENT 'Custom icon',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT 'Sort order',
    `access_count` INT NOT NULL DEFAULT 0 COMMENT 'Access count',
    `last_accessed` DATETIME NULL COMMENT 'Last access time',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_app_instance_id` (`app_instance_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App shortcut navigation table';

-- ========================================
-- 3.2 Project Management
-- ========================================

-- Projects table
CREATE TABLE IF NOT EXISTS `projects` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(100) NOT NULL COMMENT 'Project name',
    `identifier` VARCHAR(50) NOT NULL UNIQUE COMMENT 'Project identifier',
    `description` TEXT NULL COMMENT 'Project description',
    `tags` JSON NULL COMMENT 'Project tags',
    `icon` VARCHAR(255) NULL COMMENT 'Project icon URL',
    `status` ENUM('NORMAL', 'ARCHIVED', 'DELETED') NOT NULL DEFAULT 'NORMAL' COMMENT 'Project status',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Project owner ID',
    `default_resource_group` VARCHAR(50) NOT NULL DEFAULT 'default' COMMENT 'Default resource group',
    `default_timezone` VARCHAR(50) NOT NULL DEFAULT 'Asia/Shanghai' COMMENT 'Default timezone',
    `log_retention_days` INT NOT NULL DEFAULT 30 COMMENT 'Log retention days',
    `backup_strategy` ENUM('daily', 'weekly', 'monthly', 'disabled') NOT NULL DEFAULT 'daily' COMMENT 'Backup strategy',
    `access_control` VARCHAR(50) NOT NULL DEFAULT 'members' COMMENT 'Access control',
    `api_access` BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'API access permission',
    `audit_enabled` BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'Audit log enabled',
    `last_activity_at` DATETIME NULL COMMENT 'Last activity time',
    `archived_at` DATETIME NULL COMMENT 'Archive time',
    `archived_by` BIGINT UNSIGNED NULL COMMENT 'Archived by user ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    `deleted_at` DATETIME NULL COMMENT 'Soft delete time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_identifier` (`identifier`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    KEY `idx_archived_by` (`archived_by`),
    CONSTRAINT `fk_projects_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_projects_archived_by` FOREIGN KEY (`archived_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Projects table';

-- Project members table
CREATE TABLE IF NOT EXISTS `project_members` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT 'Project ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `role` VARCHAR(50) NOT NULL COMMENT 'Role: admin, developer, operator, viewer',
    `is_admin` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether project admin',
    `status` ENUM('active', 'inactive', 'pending') NOT NULL DEFAULT 'active' COMMENT 'Member status',
    `invited_by` BIGINT UNSIGNED NULL COMMENT 'Inviter ID',
    `invited_at` DATETIME NULL COMMENT 'Invitation time',
    `joined_at` DATETIME NULL COMMENT 'Join time',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_project_user` (`project_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_role` (`role`),
    KEY `idx_status` (`status`),
    KEY `idx_invited_by` (`invited_by`),
    CONSTRAINT `fk_project_members_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_project_members_invited_by` FOREIGN KEY (`invited_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Project members table';

-- Project environment variables table
CREATE TABLE IF NOT EXISTS `project_environments` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT 'Project ID',
    `name` VARCHAR(100) NOT NULL COMMENT 'Variable name',
    `value` TEXT NULL COMMENT 'Variable value',
    `type` ENUM('normal', 'sensitive') NOT NULL DEFAULT 'normal' COMMENT 'Variable type',
    `description` VARCHAR(500) NULL COMMENT 'Variable description',
    `scope` ENUM('global', 'app', 'workflow') NOT NULL DEFAULT 'global' COMMENT 'Scope',
    `scope_target` VARCHAR(100) NULL COMMENT 'Scope target',
    `is_encrypted` BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Whether encrypted storage',
    `created_by` BIGINT UNSIGNED NOT NULL COMMENT 'Creator ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_project_name_scope` (`project_id`, `name`, `scope`, `scope_target`),
    KEY `idx_created_by` (`created_by`),
    KEY `idx_type` (`type`),
    KEY `idx_scope` (`scope`),
    CONSTRAINT `fk_project_environments_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_project_environments_created_by` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Project environment variables table';

-- Project files table
CREATE TABLE IF NOT EXISTS `project_files` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT 'Project ID',
    `name` VARCHAR(255) NOT NULL COMMENT 'File name',
    `path` VARCHAR(1024) NOT NULL COMMENT 'File path',
    `type` ENUM('FILE', 'DIRECTORY') NOT NULL DEFAULT 'FILE' COMMENT 'Type: FILE, DIRECTORY',
    `size` BIGINT NOT NULL DEFAULT 0 COMMENT 'File size (bytes)',
    `mime_type` VARCHAR(128) NULL COMMENT 'MIME type',
    `download_count` INT NOT NULL DEFAULT 0 COMMENT 'Download count',
    `parent_id` BIGINT UNSIGNED NULL COMMENT 'Parent directory ID',
    `storage_path` VARCHAR(1024) NULL COMMENT 'Storage path',
    `checksum` VARCHAR(64) NULL COMMENT 'File checksum',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_project_id` (`project_id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_type` (`type`),
    KEY `idx_path` (`path`(255)),
    CONSTRAINT `fk_project_files_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_project_files_parent` FOREIGN KEY (`parent_id`) REFERENCES `project_files` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_project_files_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Project files table';

-- Project activity logs table
CREATE TABLE IF NOT EXISTS `project_activities` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT 'Project ID',
    `user_id` BIGINT UNSIGNED NULL COMMENT 'Operation user ID',
    `username` VARCHAR(64) NULL COMMENT 'Operation username',
    `action` VARCHAR(50) NOT NULL COMMENT 'Action',
    `resource_type` VARCHAR(50) NULL COMMENT 'Resource type',
    `resource_id` BIGINT UNSIGNED NULL COMMENT 'Resource ID',
    `resource_name` VARCHAR(255) NULL COMMENT 'Resource name',
    `description` TEXT NULL COMMENT 'Operation description',
    `metadata` JSON NULL COMMENT 'Operation metadata',
    `ip_address` VARCHAR(45) NULL COMMENT 'Operation IP address',
    `user_agent` VARCHAR(500) NULL COMMENT 'User agent',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    PRIMARY KEY (`id`),
    KEY `idx_project_id` (`project_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_action` (`action`),
    KEY `idx_resource_type` (`resource_type`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_project_activities_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_project_activities_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Project activity logs table';

-- Workflows table
CREATE TABLE IF NOT EXISTS `workflows` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Workflow name',
    `code` VARCHAR(32) NOT NULL COMMENT 'Workflow code',
    `description` TEXT NULL COMMENT 'Workflow description',
    `definition` JSON NOT NULL COMMENT 'Workflow definition',
    `status` ENUM('DRAFT', 'ACTIVE', 'INACTIVE', 'ARCHIVED') NOT NULL DEFAULT 'DRAFT' COMMENT 'Status',
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT 'Project ID',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_project_id` (`project_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_workflows_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_workflows_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Workflows table';

-- Workflow tasks table
CREATE TABLE IF NOT EXISTS `workflow_tasks` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Task name',
    `workflow_id` BIGINT UNSIGNED NOT NULL COMMENT 'Workflow ID',
    `schedule_type` ENUM('MANUAL', 'SCHEDULE', 'TRIGGER') NOT NULL DEFAULT 'MANUAL' COMMENT 'Schedule type',
    `cron_expression` VARCHAR(100) NULL COMMENT 'Cron expression',
    `status` ENUM('DEFAULT', 'ONLINE', 'OFFLINE') NOT NULL DEFAULT 'DEFAULT' COMMENT 'Task status',
    `next_run_at` DATETIME NULL COMMENT 'Next run time',
    `last_run_at` DATETIME NULL COMMENT 'Last run time',
    `run_count` INT NOT NULL DEFAULT 0 COMMENT 'Run count',
    `success_count` INT NOT NULL DEFAULT 0 COMMENT 'Success count',
    `failure_count` INT NOT NULL DEFAULT 0 COMMENT 'Failure count',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_workflow_id` (`workflow_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_workflow_tasks_workflow` FOREIGN KEY (`workflow_id`) REFERENCES `workflows` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Workflow tasks table';

-- Workflow execution history table
CREATE TABLE IF NOT EXISTS `workflow_executions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `task_id` BIGINT UNSIGNED NOT NULL COMMENT 'Task ID',
    `execution_id` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Execution unique ID',
    `status` ENUM('PENDING', 'RUNNING', 'SUCCESS', 'FAILURE', 'STOPPED') NOT NULL DEFAULT 'PENDING' COMMENT 'Execution status',
    `trigger_type` ENUM('MANUAL', 'SCHEDULE', 'TRIGGER') NOT NULL DEFAULT 'MANUAL' COMMENT 'Trigger type',
    `trigger_by` BIGINT UNSIGNED NULL COMMENT 'Triggered by ID',
    `start_time` DATETIME NULL COMMENT 'Start time',
    `end_time` DATETIME NULL COMMENT 'End time',
    `duration` INT NOT NULL DEFAULT 0 COMMENT 'Execution duration (seconds)',
    `error_message` TEXT NULL COMMENT 'Error message',
    `execution_log` TEXT NULL COMMENT 'Execution log',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_task_id` (`task_id`),
    KEY `idx_status` (`status`),
    KEY `idx_trigger_by` (`trigger_by`),
    CONSTRAINT `fk_workflow_executions_task` FOREIGN KEY (`task_id`) REFERENCES `workflow_tasks` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Workflow execution history table';

-- Resource groups table
CREATE TABLE IF NOT EXISTS `resource_groups` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT 'Project ID',
    `name` VARCHAR(64) NOT NULL COMMENT 'Resource group name',
    `code` VARCHAR(32) NOT NULL COMMENT 'Resource group code',
    `description` TEXT NULL COMMENT 'Resource group description',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT 'Sort order',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Status: 0-disabled, 1-enabled',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_project_id` (`project_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_resource_groups_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_resource_groups_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Resource groups table';

-- Database connections table
CREATE TABLE IF NOT EXISTS `database_connections` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Connection name',
    `db_type` ENUM('mysql', 'postgresql', 'redis', 'mongodb') NOT NULL COMMENT 'Database type',
    `host` VARCHAR(255) NOT NULL COMMENT 'Host address',
    `port` INT NOT NULL COMMENT 'Port number',
    `database` VARCHAR(64) NULL COMMENT 'Database name',
    `username` VARCHAR(64) NULL COMMENT 'Username',
    `password` VARCHAR(255) NULL COMMENT 'Encrypted password',
    `ssl_enabled` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether SSL enabled',
    `connection_timeout` INT NOT NULL DEFAULT 30 COMMENT 'Connection timeout (seconds)',
    `max_connections` INT NOT NULL DEFAULT 10 COMMENT 'Max connections',
    `status` ENUM('CONNECTED', 'DISCONNECTED', 'ERROR') NOT NULL DEFAULT 'CONNECTED' COMMENT 'Connection status',
    `version` VARCHAR(32) NULL COMMENT 'Database version',
    `charset` VARCHAR(32) NULL COMMENT 'Character set',
    `description` TEXT NULL COMMENT 'Description info',
    `last_connected_at` DATETIME NULL COMMENT 'Last connected time',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `resource_group_id` BIGINT UNSIGNED NULL COMMENT 'Resource group ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_resource_group_id` (`resource_group_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_db_connections_resource_group` FOREIGN KEY (`resource_group_id`) REFERENCES `resource_groups` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Database connections table';

-- Servers table
CREATE TABLE IF NOT EXISTS `servers` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Server name',
    `hostname` VARCHAR(255) NOT NULL COMMENT 'Hostname',
    `host` VARCHAR(255) NOT NULL COMMENT 'Host address (IP or domain) for SSH/Agent priority connection',
    `internal_ip` VARCHAR(45) NULL COMMENT 'Internal IP address',
    `ipv6_address` VARCHAR(45) NULL COMMENT 'IPv6 address for future expansion',
    `ssh_port` INT NOT NULL DEFAULT 22 COMMENT 'SSH port',
    `ssh_credential_id` VARCHAR(64) NULL COMMENT 'SSH credential ID from key management',
    `os_distro` VARCHAR(32) NULL COMMENT 'Operating system distribution',
    `os_version` VARCHAR(64) NULL COMMENT 'Operating system version',
    `kernel_version` VARCHAR(64) NULL COMMENT 'Kernel version',
    `cpu_cores` INT NOT NULL DEFAULT 0 COMMENT 'CPU cores count',
    `memory_total` BIGINT NOT NULL DEFAULT 0 COMMENT 'Total memory (MB)',
    `disk_total` BIGINT NOT NULL DEFAULT 0 COMMENT 'Total disk space (MB)',
    `architecture` VARCHAR(16) NULL COMMENT 'System architecture',
    `resource_group_id` BIGINT UNSIGNED NULL COMMENT 'Resource group ID',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner user ID',
    `description` TEXT NULL COMMENT 'Server description',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    `deleted_at` DATETIME NULL COMMENT 'Soft delete time',
    PRIMARY KEY (`id`),
    KEY `idx_name` (`name`),
    KEY `idx_host` (`host`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_resource_group_id` (`resource_group_id`),
    KEY `idx_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_servers_resource_group` FOREIGN KEY (`resource_group_id`) REFERENCES `resource_groups` (`id`) ON DELETE SET NULL,
    CONSTRAINT `fk_servers_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Servers table';

-- Server agents table
CREATE TABLE IF NOT EXISTS `server_agents` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `server_id` BIGINT UNSIGNED NOT NULL COMMENT 'Server ID',
    `agent_id` VARCHAR(64) NOT NULL COMMENT 'Agent unique identifier',
    `deployment_type` ENUM('docker', 'systemd') DEFAULT 'docker' COMMENT 'Agent deployment type: docker container or systemd service',
    `container_id` VARCHAR(64) NULL COMMENT 'Container ID (for Docker deployment)',
    `container_name` VARCHAR(128) NULL COMMENT 'Container name (for Docker deployment)',
    `service_name` VARCHAR(64) NULL COMMENT 'Service name (for systemd deployment)',
    `binary_path` VARCHAR(255) NULL COMMENT 'Binary file path (for systemd deployment)',
    `config_path` VARCHAR(255) NULL COMMENT 'Configuration file path',
    `agent_ip` VARCHAR(45) NULL COMMENT 'Agent IP address',
    `agent_port` INT NOT NULL DEFAULT 8080 COMMENT 'Agent communication port',
    `version` VARCHAR(32) NULL COMMENT 'Agent version',
    `pull_mode` BOOLEAN DEFAULT TRUE COMMENT 'Whether Pull mode is enabled',
    `pull_interval` INT DEFAULT 30 COMMENT 'Pull interval in seconds',
    `last_heartbeat_at` DATETIME NULL COMMENT 'Last heartbeat time',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_server_id` (`server_id`),
    UNIQUE KEY `uk_agent_id` (`agent_id`),
    KEY `idx_server_id` (`server_id`),
    KEY `idx_last_heartbeat` (`last_heartbeat_at`),
    KEY `idx_deployment_type` (`deployment_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Server agents table';

-- Application instances table
CREATE TABLE IF NOT EXISTS `app_instances` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `project_id` BIGINT UNSIGNED NOT NULL COMMENT 'Project ID',
    `name` VARCHAR(64) NOT NULL COMMENT 'App instance name',
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT 'App template ID',
    `server_id` BIGINT UNSIGNED NOT NULL COMMENT 'Server ID',
    `container_id` VARCHAR(64) NULL COMMENT 'Container ID',
    `container_name` VARCHAR(64) NULL COMMENT 'Container name',
    `image_name` VARCHAR(255) NULL COMMENT 'Image name',
    `image_tag` VARCHAR(100) NULL COMMENT 'Image tag',
    `status` ENUM('DEFAULT', 'DEPLOYMENT', 'RUNNING', 'PAUSED', 'STOPPED', 'UPDATE') NOT NULL DEFAULT 'DEFAULT' COMMENT 'App status',
    `started_at` DATETIME NULL COMMENT 'Start time',
    `stopped_at` DATETIME NULL COMMENT 'Stop time',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Application instances table';

-- SSL certificates table
CREATE TABLE IF NOT EXISTS `ssl_certificates` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Certificate name',
    `domain` VARCHAR(255) NOT NULL COMMENT 'Domain name',
    `certificate_type` ENUM('LETS_ENCRYPT', 'COMMERCIAL', 'SELF_SIGNED') NOT NULL DEFAULT 'LETS_ENCRYPT' COMMENT 'Certificate type',
    `certificate_data` TEXT NOT NULL COMMENT 'Certificate content',
    `private_key_data` TEXT NOT NULL COMMENT 'Private key content',
    `certificate_chain` TEXT NULL COMMENT 'Certificate chain',
    `issuer` VARCHAR(255) NULL COMMENT 'Issuer',
    `subject` VARCHAR(255) NULL COMMENT 'Subject',
    `serial_number` VARCHAR(64) NULL COMMENT 'Serial number',
    `not_before` DATETIME NULL COMMENT 'Effective time',
    `not_after` DATETIME NULL COMMENT 'Expiration time',
    `auto_renew` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Auto renew',
    `status` ENUM('PENDING', 'VALID', 'EXPIRED', 'REVOKED') NOT NULL DEFAULT 'PENDING' COMMENT 'Status',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_domain` (`domain`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    KEY `idx_not_after` (`not_after`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='SSL certificates table';

-- Secret key management table
CREATE TABLE IF NOT EXISTS `secret_keys` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Secret key name',
    `key_type` ENUM('SECRET_KEY', 'ACCOUNT', 'FILE') NOT NULL COMMENT 'Secret key type',
    `description` TEXT NULL COMMENT 'Description',
    `custom_fields` JSON NULL COMMENT 'Custom fields',
    `expires_at` DATETIME NULL COMMENT 'Expiration time',
    `resource_group_id` BIGINT UNSIGNED NULL COMMENT 'Resource group ID',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_resource_group_id` (`resource_group_id`),
    KEY `idx_key_type` (`key_type`),
    KEY `idx_expires_at` (`expires_at`),
    CONSTRAINT `fk_secret_keys_resource_group` FOREIGN KEY (`resource_group_id`) REFERENCES `resource_groups` (`id`) ON DELETE SET NULL,
    CONSTRAINT `fk_secret_keys_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Secret key management table';

-- Application gateways table
CREATE TABLE IF NOT EXISTS `app_gateways` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Gateway name',
    `server_id` BIGINT UNSIGNED NOT NULL COMMENT 'Server ID',
    `container_id` VARCHAR(64) NULL COMMENT 'Gateway container ID',
    `description` TEXT NULL COMMENT 'Gateway description',
    `status` ENUM('DEFAULT', 'RUNNING', 'STOPPED', 'ERROR') NOT NULL DEFAULT 'DEFAULT' COMMENT 'Gateway status',
    `started_at` DATETIME NULL COMMENT 'Start time',
    `stopped_at` DATETIME NULL COMMENT 'Stop time',
    `resource_group_id` BIGINT UNSIGNED NULL COMMENT 'Resource group ID',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_server_id` (`server_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_resource_group_id` (`resource_group_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_app_gateways_server` FOREIGN KEY (`server_id`) REFERENCES `servers` (`id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_app_gateways_resource_group` FOREIGN KEY (`resource_group_id`) REFERENCES `resource_groups` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Application gateways table';

-- Application gateway publishing table
CREATE TABLE IF NOT EXISTS `app_gateways_publishes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `app_instance_id` BIGINT UNSIGNED NOT NULL COMMENT 'App instance ID',
    `app_gateway_id` BIGINT UNSIGNED NOT NULL COMMENT 'App gateway ID',
    `service_domain` VARCHAR(255) NOT NULL COMMENT 'Domain (service name)',
    `service_port` INT NOT NULL DEFAULT 8080 COMMENT 'Service port',
    `alert_rule_id` BIGINT UNSIGNED NULL COMMENT 'Monitor alert rule ID',
    `limit_rules` TEXT NULL COMMENT 'Access control policy',
    `health_check_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Health check enabled',
    `audit_log_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Audit log enabled',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_app_instance_id` (`app_instance_id`),
    KEY `idx_app_gateway_id` (`app_gateway_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_service_domain` (`service_domain`),
    CONSTRAINT `fk_gateway_publishes_instance` FOREIGN KEY (`app_instance_id`) REFERENCES `app_instances` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_gateway_publishes_gateway` FOREIGN KEY (`app_gateway_id`) REFERENCES `app_gateways` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Application gateway publishing table';

-- Application gateway access control rules table
CREATE TABLE IF NOT EXISTS `app_gateway_access_rules` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `gateway_id` BIGINT UNSIGNED NOT NULL COMMENT 'Gateway ID',
    `rule_name` VARCHAR(64) NOT NULL COMMENT 'Rule name',
    `rule_type` ENUM('IP_WHITELIST', 'IP_BLACKLIST', 'RATE_LIMIT') NOT NULL COMMENT 'Rule type',
    `limit_count` INT NOT NULL COMMENT 'Limit count',
    `time_window` INT NOT NULL COMMENT 'Time window (seconds)',
    `target_path` VARCHAR(255) NULL COMMENT 'Target path',
    `target_ip` VARCHAR(45) NULL COMMENT 'Target IP',
    `action` ENUM('BLOCK', 'ALLOW', 'REDIRECT') NOT NULL DEFAULT 'BLOCK' COMMENT 'Trigger action',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Status: 0-disabled, 1-enabled',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_gateway_id` (`gateway_id`),
    KEY `idx_rule_type` (`rule_type`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_gateway_access_rules_gateway` FOREIGN KEY (`gateway_id`) REFERENCES `app_gateways` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Application gateway access control rules table';

ALTER TABLE `app_gateways_publishes` ADD CONSTRAINT `fk_gateway_publishes_alert_rule` FOREIGN KEY (`alert_rule_id`) REFERENCES `alert_rules` (`id`) ON DELETE SET NULL;
ALTER TABLE `database_connections` ADD CONSTRAINT `fk_database_connections_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT;
ALTER TABLE `ssl_certificates` ADD CONSTRAINT `fk_ssl_certificates_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT;
ALTER TABLE `app_gateways` ADD CONSTRAINT `fk_app_gateways_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT;
ALTER TABLE `app_gateways_publishes` ADD CONSTRAINT `fk_gateway_publishes_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT;

-- ========================================
-- 3.3 App Store
-- ========================================

-- App store categories table
CREATE TABLE IF NOT EXISTS `app_store_categories` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Category name',
    `code` VARCHAR(32) NOT NULL UNIQUE COMMENT 'Category code',
    `parent_id` BIGINT UNSIGNED NULL COMMENT 'Parent category ID',
    `icon` VARCHAR(255) NULL COMMENT 'Icon URL',
    `description` TEXT NULL COMMENT 'Category description',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT 'Sort order',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Status: 0-disabled, 1-enabled',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_app_categories_parent` FOREIGN KEY (`parent_id`) REFERENCES `app_store_categories` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App store categories table';

-- App store templates table
CREATE TABLE IF NOT EXISTS `app_store_templates` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'App name',
    `code` VARCHAR(32) NOT NULL UNIQUE COMMENT 'App code',
    `category_id` BIGINT UNSIGNED NOT NULL COMMENT 'Category ID',
    `version` VARCHAR(32) NOT NULL COMMENT 'Version number',
    `icon` VARCHAR(255) NULL COMMENT 'Icon URL',
    `description` TEXT NULL COMMENT 'App description',
    `official_url` VARCHAR(255) NULL COMMENT 'Official website',
    `source_url` VARCHAR(255) NULL COMMENT 'Source code URL',
    `compose_template` TEXT NOT NULL COMMENT 'Docker Compose template',
    `download_count` INT NOT NULL DEFAULT 0 COMMENT 'Download count',
    `star_count` INT NOT NULL DEFAULT 0 COMMENT 'Star count',
    `rating` DECIMAL(3,2) NOT NULL DEFAULT 0.00 COMMENT 'Rating',
    `is_official` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether official app',
    `is_featured` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether featured app',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Status: 0-delisted, 1-listed',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_status` (`status`),
    KEY `idx_is_featured` (`is_featured`),
    KEY `idx_download_count` (`download_count`),
    CONSTRAINT `fk_app_templates_category` FOREIGN KEY (`category_id`) REFERENCES `app_store_categories` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App store templates table';

-- App store wishlists table
CREATE TABLE IF NOT EXISTS `app_store_wishlists` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'App name',
    `version` VARCHAR(32) NULL COMMENT 'Version number',
    `source_url` VARCHAR(255) NULL COMMENT 'Source URL',
    `description` TEXT NULL COMMENT 'Requirement description',
    `reward_amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT 'Reward amount',
    `priority` TINYINT(1) NOT NULL DEFAULT 3 COMMENT 'Priority: 1-high, 2-medium, 3-low',
    `status` ENUM('PENDING', 'IN_PROGRESS', 'COMPLETED', 'EXPIRED') NOT NULL DEFAULT 'PENDING' COMMENT 'Status',
    `view_count` INT NOT NULL DEFAULT 0 COMMENT 'View count',
    `like_count` INT NOT NULL DEFAULT 0 COMMENT 'Star count',
    `vote_count` INT NOT NULL DEFAULT 0 COMMENT 'Vote count',
    `comment_count` INT NOT NULL DEFAULT 0 COMMENT 'Comment count',
    `submitter_id` BIGINT UNSIGNED NOT NULL COMMENT 'Submitter ID',
    `completed_at` DATETIME NULL COMMENT 'Completion time',
    `expires_at` DATETIME NULL COMMENT 'Expiration time',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_submitter_id` (`submitter_id`),
    KEY `idx_status` (`status`),
    KEY `idx_priority` (`priority`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App store wishlists table';

-- App store reviews table
CREATE TABLE IF NOT EXISTS `app_store_reviews` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT 'App template ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `rating` TINYINT(1) NOT NULL COMMENT 'Rating (1-5 points)',
    `content` TEXT NULL COMMENT 'Review content',
    `tags` JSON NULL COMMENT 'Review tags',
    `is_helpful` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether helpful',
    `helpful_count` INT NOT NULL DEFAULT 0 COMMENT 'Helpful count',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_rating` (`rating`),
    CONSTRAINT `fk_reviews_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE CASCADE,
    CONSTRAINT `chk_rating` CHECK (`rating` >= 1 AND `rating` <= 5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App reviews table';

-- App favorites table
CREATE TABLE IF NOT EXISTS `app_store_favorites` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT 'App template ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_template_user` (`template_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_favorites_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App favorites table';

-- App stars table
CREATE TABLE IF NOT EXISTS `app_store_stars` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT 'App template ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_template_user` (`template_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_stars_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App star table';

-- App reports table
CREATE TABLE IF NOT EXISTS `app_store_reports` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT 'App template ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'Reporter user ID',
    `reason` ENUM('SPAM', 'INAPPROPRIATE', 'COPYRIGHT') NOT NULL COMMENT 'Report reason',
    `detail` TEXT NULL COMMENT 'Detail description',
    `status` ENUM('PENDING', 'PROCESSING', 'RESOLVED', 'REJECTED') NOT NULL DEFAULT 'PENDING' COMMENT 'Processing status',
    `handled_by` BIGINT UNSIGNED NULL COMMENT 'Handler ID',
    `handled_at` DATETIME NULL COMMENT 'Processing time',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    KEY `idx_handled_by` (`handled_by`),
    CONSTRAINT `fk_app_reports_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_app_reports_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_app_reports_handled_by` FOREIGN KEY (`handled_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App report table';

-- App download records table
CREATE TABLE IF NOT EXISTS `app_store_downloads` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NOT NULL COMMENT 'App template ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `ip_address` VARCHAR(45) NULL COMMENT 'IP address',
    `user_agent` VARCHAR(255) NULL COMMENT 'User agent',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_app_downloads_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_app_downloads_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App download records table';

-- App wishlist comments table
CREATE TABLE IF NOT EXISTS `app_store_wishlist_comments` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `wishlist_id` BIGINT UNSIGNED NOT NULL COMMENT 'Wishlist ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `parent_id` BIGINT UNSIGNED NULL COMMENT 'Parent comment ID',
    `content` TEXT NOT NULL COMMENT 'Comment content',
    `like_count` INT NOT NULL DEFAULT 0 COMMENT 'Star count',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_wishlist_id` (`wishlist_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_parent_id` (`parent_id`),
    CONSTRAINT `fk_wishlist_comments_wishlist` FOREIGN KEY (`wishlist_id`) REFERENCES `app_store_wishlists` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_wishlist_comments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_wishlist_comments_parent` FOREIGN KEY (`parent_id`) REFERENCES `app_store_wishlist_comments` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App wishlist comments table';

-- App wishlist votes table
CREATE TABLE IF NOT EXISTS `app_store_wishlist_votes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `wishlist_id` BIGINT UNSIGNED NOT NULL COMMENT 'Wishlist ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_wishlist_user` (`wishlist_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_wishlist_votes_wishlist` FOREIGN KEY (`wishlist_id`) REFERENCES `app_store_wishlists` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_wishlist_votes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App wishlist votes table';

-- App wishlist likes table
CREATE TABLE IF NOT EXISTS `app_store_wishlist_likes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `wishlist_id` BIGINT UNSIGNED NOT NULL COMMENT 'Wishlist ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_wishlist_user` (`wishlist_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_wishlist_likes_wishlist` FOREIGN KEY (`wishlist_id`) REFERENCES `app_store_wishlists` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_wishlist_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App wishlist likes table';

-- App wishlist reports table
CREATE TABLE IF NOT EXISTS `app_store_wishlist_reports` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `wishlist_id` BIGINT UNSIGNED NOT NULL COMMENT 'Wishlist ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'Reporter user ID',
    `reason` ENUM('SPAM', 'INAPPROPRIATE', 'DUPLICATE') NOT NULL COMMENT 'Report reason',
    `detail` TEXT NULL COMMENT 'Detail description',
    `status` ENUM('PENDING', 'PROCESSING', 'RESOLVED', 'REJECTED') NOT NULL DEFAULT 'PENDING' COMMENT 'Processing status',
    `handled_by` BIGINT UNSIGNED NULL COMMENT 'Handler ID',
    `handled_at` DATETIME NULL COMMENT 'Processing time',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_wishlist_id` (`wishlist_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    KEY `idx_handled_by` (`handled_by`),
    CONSTRAINT `fk_wishlist_reports_wishlist` FOREIGN KEY (`wishlist_id`) REFERENCES `app_store_wishlists` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_wishlist_reports_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_wishlist_reports_handled_by` FOREIGN KEY (`handled_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App wishlist reports table';

-- App deployment records table
CREATE TABLE IF NOT EXISTS `app_deployments` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `deployment_id` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Deployment unique ID',
    `template_id` BIGINT UNSIGNED NULL COMMENT 'App template ID',
    `app_instance_id` BIGINT UNSIGNED NULL COMMENT 'App instance ID',
    `server_id` BIGINT UNSIGNED NOT NULL COMMENT 'Server ID',
    `status` ENUM('PENDING', 'RUNNING', 'SUCCESS', 'FAILED', 'CANCELLED') NOT NULL DEFAULT 'PENDING' COMMENT 'Deployment status',
    `progress` TINYINT NOT NULL DEFAULT 0 COMMENT 'Deployment progress',
    `estimated_time` INT NOT NULL DEFAULT 0 COMMENT 'Estimated time (seconds)',
    `start_time` DATETIME NULL COMMENT 'Start time',
    `end_time` DATETIME NULL COMMENT 'End time',
    `error_message` TEXT NULL COMMENT 'Error message',
    `deployment_log` TEXT NULL COMMENT 'Deployment log',
    `config_data` JSON NULL COMMENT 'Deployment config',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_server_id` (`server_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_deployments_template` FOREIGN KEY (`template_id`) REFERENCES `app_store_templates` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App deployment records table';

ALTER TABLE `app_shortcuts` ADD CONSTRAINT `fk_app_shortcuts_app_instance` FOREIGN KEY (`app_instance_id`) REFERENCES `app_instances` (`id`) ON DELETE CASCADE;
ALTER TABLE `app_shortcuts` ADD CONSTRAINT `fk_app_shortcuts_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
ALTER TABLE `app_store_wishlists` ADD CONSTRAINT `fk_app_wishlists_submitter` FOREIGN KEY (`submitter_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

-- ========================================
-- 3.4 Platform Management
-- ========================================

-- System configuration table
CREATE TABLE IF NOT EXISTS `system_configs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `config_key` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Configuration key',
    `config_value` TEXT NULL COMMENT 'Configuration value',
    `config_type` ENUM('STRING', 'INTEGER', 'BOOLEAN', 'JSON', 'FLOAT') NOT NULL DEFAULT 'STRING' COMMENT 'Configuration type',
    `category` VARCHAR(32) NOT NULL COMMENT 'Configuration category',
    `description` TEXT NULL COMMENT 'Configuration description',
    `is_readonly` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether read-only',
    `is_encrypted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether encrypted',
    `default_value` TEXT NULL COMMENT 'Default value',
    `validation_rules` JSON NULL COMMENT 'Validation rules',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT 'Sort order',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_category` (`category`),
    KEY `idx_config_type` (`config_type`),
    KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='System configuration table';

-- Service configuration parameters table
CREATE TABLE IF NOT EXISTS `service_configs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Configuration code (unique identifier)',
    `config_key` VARCHAR(64) NOT NULL COMMENT 'Configuration key',
    `config_value` TEXT NULL COMMENT 'Configuration value',
    `config_type` ENUM('STRING', 'INTEGER', 'BOOLEAN', 'JSON', 'FLOAT') NOT NULL DEFAULT 'STRING' COMMENT 'Configuration type',
    `category` VARCHAR(32) NOT NULL COMMENT 'Configuration category',
    `description` TEXT NULL COMMENT 'Configuration description',
    `is_readonly` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether read-only',
    `is_encrypted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether encrypted',
    `default_value` TEXT NULL COMMENT 'Default value',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT 'Sort order',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_category` (`category`),
    KEY `idx_config_type` (`config_type`),
    KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Service configuration parameters table';

-- Webhook configuration table
CREATE TABLE IF NOT EXISTS `webhook_configs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Webhook name',
    `url` VARCHAR(255) NOT NULL COMMENT 'Callback URL',
    `secret` VARCHAR(255) NULL COMMENT 'Signature secret',
    `events` JSON NOT NULL COMMENT 'Listen events',
    `headers` JSON NULL COMMENT 'Custom headers',
    `timeout` INT NOT NULL DEFAULT 30 COMMENT 'Timeout (seconds)',
    `retry_count` INT NOT NULL DEFAULT 3 COMMENT 'Retry count',
    `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Whether enabled',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_is_active` (`is_active`),
    CONSTRAINT `fk_webhook_configs_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Webhook configuration table';

-- Webhook execution log table
CREATE TABLE IF NOT EXISTS `webhook_logs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `webhook_id` BIGINT UNSIGNED NOT NULL COMMENT 'Webhook ID',
    `event_type` VARCHAR(32) NOT NULL COMMENT 'Event type',
    `payload` JSON NOT NULL COMMENT 'Request payload',
    `request_id` VARCHAR(64) NOT NULL COMMENT 'Request ID',
    `status_code` INT NULL COMMENT 'Response status code',
    `response_body` TEXT NULL COMMENT 'Response body',
    `response_time` INT NOT NULL DEFAULT 0 COMMENT 'Response time (milliseconds)',
    `retry_count` INT NOT NULL DEFAULT 0 COMMENT 'Retry count',
    `success` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether successful',
    `error_message` TEXT NULL COMMENT 'Error message',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    PRIMARY KEY (`id`),
    KEY `idx_webhook_id` (`webhook_id`),
    KEY `idx_event_type` (`event_type`),
    KEY `idx_success` (`success`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_webhook_logs_webhook` FOREIGN KEY (`webhook_id`) REFERENCES `webhook_configs` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Webhook execution log table';

-- Repository configuration table
CREATE TABLE IF NOT EXISTS `repository_configs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Configuration name',
    `type` ENUM('docker', 'apt', 'yum', 'npm', 'pip') NOT NULL COMMENT 'Repository type',
    `url` VARCHAR(255) NOT NULL COMMENT 'Repository URL',
    `username` VARCHAR(64) NULL COMMENT 'Username',
    `password` VARCHAR(255) NULL COMMENT 'Password',
    `is_default` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether default',
    `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether system config',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Status: 0-disabled, 1-enabled',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Repository configuration table';

-- Platform update records table
CREATE TABLE IF NOT EXISTS `platform_updates` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `version` VARCHAR(32) NOT NULL COMMENT 'Version number',
    `changelog` TEXT NULL COMMENT 'Changelog',
    `download_url` VARCHAR(255) NULL COMMENT 'Download URL',
    `file_size` BIGINT NOT NULL DEFAULT 0 COMMENT 'File size',
    `checksum` VARCHAR(64) NULL COMMENT 'File checksum',
    `status` ENUM('PENDING', 'DOWNLOADING', 'INSTALLING', 'SUCCESS', 'FAILED') NOT NULL DEFAULT 'PENDING' COMMENT 'Update status',
    `started_at` DATETIME NULL COMMENT 'Start time',
    `completed_at` DATETIME NULL COMMENT 'Completion time',
    `error_message` TEXT NULL COMMENT 'Error message',
    `backup_path` VARCHAR(1024) NULL COMMENT 'Backup path',
    `updated_by` BIGINT UNSIGNED NULL COMMENT 'Updated by ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_version` (`version`),
    KEY `idx_status` (`status`),
    KEY `idx_updated_by` (`updated_by`),
    CONSTRAINT `fk_platform_updates_updated_by` FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Platform update records table';

-- Container cluster nodes table
CREATE TABLE IF NOT EXISTS `docker_swarm_nodes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `node_id` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Docker node ID',
    `hostname` VARCHAR(255) NOT NULL COMMENT 'Host name',
    `ip_address` VARCHAR(45) NOT NULL COMMENT 'IP address',
    `role` ENUM('manager', 'worker') NOT NULL COMMENT 'Node role',
    `status` ENUM('ready', 'down', 'unknown') NOT NULL DEFAULT 'ready' COMMENT 'Node status',
    `availability` ENUM('active', 'pause', 'drain') NOT NULL DEFAULT 'active' COMMENT 'Availability',
    `engine_version` VARCHAR(32) NULL COMMENT 'Engine version',
    `labels` JSON NULL COMMENT 'Node labels',
    `resources` JSON NULL COMMENT 'Resource info',
    `joined_at` DATETIME NULL COMMENT 'Join time',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_role` (`role`),
    KEY `idx_status` (`status`),
    KEY `idx_availability` (`availability`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Container cluster nodes table';

-- Container images table
CREATE TABLE IF NOT EXISTS `docker_images` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `image_id` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Docker image ID',
    `repository` VARCHAR(255) NOT NULL COMMENT 'Image repository',
    `tag` VARCHAR(100) NOT NULL COMMENT 'Image tag',
    `digest` VARCHAR(128) NULL COMMENT 'Image digest',
    `size` BIGINT NOT NULL DEFAULT 0 COMMENT 'Image size',
    `architecture` VARCHAR(32) NULL COMMENT 'Architecture',
    `os` VARCHAR(32) NULL COMMENT 'Operating system',
    `status` ENUM('AVAILABLE', 'PULLING', 'ERROR') NOT NULL DEFAULT 'AVAILABLE' COMMENT 'Image status',
    `usage_count` INT NOT NULL DEFAULT 0 COMMENT 'Usage count',
    `last_used_at` DATETIME NULL COMMENT 'Last used time',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_repository_tag` (`repository`, `tag`),
    KEY `idx_status` (`status`),
    KEY `idx_usage_count` (`usage_count`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Container images table';

-- Users table
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Username',
    `email` VARCHAR(255) NOT NULL UNIQUE COMMENT 'Email',
    `password_hash` VARCHAR(255) NOT NULL COMMENT 'Password hash',
    `nickname` VARCHAR(64) NULL COMMENT 'Nickname',
    `avatar` VARCHAR(255) NULL COMMENT 'Avatar URL',
    `phone` VARCHAR(20) NULL COMMENT 'Phone number',
    `gender` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Gender: 0-unknown, 1-male, 2-female',
    `signature` VARCHAR(255) NULL COMMENT 'Personal signature',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Status: 0-disabled, 1-enabled',
    `last_login_at` DATETIME NULL COMMENT 'Last login time',
    `last_login_ip` VARCHAR(45) NULL COMMENT 'Last login IP',
    `timezone` VARCHAR(64) NOT NULL DEFAULT 'UTC' COMMENT 'Timezone',
    `language` VARCHAR(10) NOT NULL DEFAULT 'zh-CN' COMMENT 'Language',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`),
    KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Users table';

-- Modules table
CREATE TABLE IF NOT EXISTS `modules` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'module name',
    `code` VARCHAR(64) NOT NULL UNIQUE COMMENT 'module code',
    `description` TEXT NULL COMMENT 'module description',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_module_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='modules table';

-- Permissions table
CREATE TABLE IF NOT EXISTS `permissions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `parent_code` VARCHAR(64) NULL COMMENT 'Parent permission code',
    `scope` VARCHAR(64) NOT NULL COMMENT 'Permission scope: platform, project',
    `name` VARCHAR(64) NOT NULL COMMENT 'Permission name',
    `code` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Permission code',
    `module` VARCHAR(32) NOT NULL COMMENT 'Module name',
    `action` VARCHAR(32) NOT NULL COMMENT 'Action name',
    `resource` VARCHAR(64) NULL COMMENT 'Resource identifier (URI)',
    `element` VARCHAR(64) NULL COMMENT 'UI element ID',
    `description` TEXT NULL COMMENT 'Permission description',
    `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether system permission',
    `is_menu` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether menu permission',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT 'Sort order',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Status: -1-deleted, 0-disabled, 1-enabled',
    `created_by` BIGINT UNSIGNED NULL COMMENT 'Creator ID',
    `updated_by` BIGINT UNSIGNED NULL COMMENT 'Updated by ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Permissions table';

-- Roles table
CREATE TABLE IF NOT EXISTS `roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Role name',
    `code` VARCHAR(32) NOT NULL UNIQUE COMMENT 'Role code',
    `description` TEXT NULL COMMENT 'Role description',
    `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether system role',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT 'Sort order',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Status: -1-deleted, 0-disabled, 1-enabled',
    `created_by` BIGINT UNSIGNED NULL COMMENT 'Creator ID',
    `updated_by` BIGINT UNSIGNED NULL COMMENT 'Updated by ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_by` (`created_by`),
    KEY `idx_updated_by` (`updated_by`),
    CONSTRAINT `fk_roles_created_by` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
    CONSTRAINT `fk_roles_updated_by` FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Roles table';

-- User roles association table
CREATE TABLE IF NOT EXISTS `user_roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `role_id` BIGINT UNSIGNED NOT NULL COMMENT 'Role ID',
    `granted_by` BIGINT UNSIGNED NULL COMMENT 'Granted by ID',
    `granted_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Grant time',
    `expires_at` DATETIME NULL COMMENT 'Expiration time',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Status: -1-deleted, 0-disabled, 1-enabled',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
    KEY `idx_role_id` (`role_id`),
    KEY `idx_granted_by` (`granted_by`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_user_roles_granted_by` FOREIGN KEY (`granted_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User roles association table';

-- Role permissions association table
CREATE TABLE IF NOT EXISTS `role_permissions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `role_id` BIGINT UNSIGNED NOT NULL COMMENT 'Role ID',
    `permission_code` VARCHAR(64) NOT NULL COMMENT 'Permission code',
    `granted_by` BIGINT UNSIGNED NULL COMMENT 'Granted by ID',
    `granted_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Grant time',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Status: -1-deleted, 0-disabled, 1-enabled',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_permission` (`role_id`, `permission_code`),
    KEY `idx_permission_code` (`permission_code`),
    KEY `idx_granted_by` (`granted_by`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_code`) REFERENCES `permissions` (`code`) ON DELETE CASCADE,
    CONSTRAINT `fk_role_permissions_granted_by` FOREIGN KEY (`granted_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Role permissions association table';

-- API access tokens table
CREATE TABLE IF NOT EXISTS `api_tokens` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Token name',
    `token` VARCHAR(255) NOT NULL UNIQUE COMMENT 'Token value (encrypted storage)',
    `token_hash` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Token hash value',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `scopes` JSON NULL COMMENT 'Permission scope',
    `description` TEXT NULL COMMENT 'Token description',
    `last_used_at` DATETIME NULL COMMENT 'Last used time',
    `last_used_ip` VARCHAR(45) NULL COMMENT 'Last used IP',
    `expires_at` DATETIME NULL COMMENT 'Expiration time',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_expires_at` (`expires_at`),
    CONSTRAINT `fk_api_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API access tokens table';

-- User two-factor authentication table
CREATE TABLE IF NOT EXISTS `user_two_factor` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `method` VARCHAR(32) NOT NULL COMMENT 'Authentication method (TOTP, EMAIL)',
    `secret` VARCHAR(255) NULL COMMENT 'Secret (encrypted storage)',
    `backup_codes` JSON NULL COMMENT 'Backup codes',
    `email` VARCHAR(255) NULL COMMENT 'Email',
    `enabled` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether enabled',
    `verified_at` DATETIME NULL COMMENT 'Verification time',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_method` (`user_id`, `method`),
    KEY `idx_enabled` (`enabled`),
    CONSTRAINT `fk_user_two_factor_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User two-factor authentication table';

-- User profile configuration table
CREATE TABLE IF NOT EXISTS `user_profile` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `config_key` VARCHAR(64) NOT NULL COMMENT 'Configuration key',
    `config_value` TEXT NULL COMMENT 'Configuration value',
    `description` TEXT NULL COMMENT 'Configuration description',
    `is_readonly` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether read-only',
    `is_encrypted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether encrypted',
    `default_value` TEXT NULL COMMENT 'Default value',
    `sort_order` INT NOT NULL DEFAULT 0 COMMENT 'Sort order',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_config_key` (`user_id`, `config_key`),
    KEY `idx_sort_order` (`sort_order`),
    CONSTRAINT `fk_user_profile_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User profile configuration table';

-- User login history table
CREATE TABLE IF NOT EXISTS `user_login_history` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `ip_address` VARCHAR(45) NULL COMMENT 'IP address',
    `user_agent` VARCHAR(255) NULL COMMENT 'User agent',
    `location` VARCHAR(100) NULL COMMENT 'Login location',
    `device` VARCHAR(100) NULL COMMENT 'Device info',
    `browser` VARCHAR(100) NULL COMMENT 'Browser info',
    `login_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Login time',
    `logout_time` DATETIME NULL COMMENT 'Logout time',
    `status` ENUM('ACTIVE', 'EXPIRED', 'LOGOUT') NOT NULL DEFAULT 'ACTIVE' COMMENT 'Session status',
    `session_id` VARCHAR(128) NULL COMMENT 'Session ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_login_time` (`login_time`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_user_login_history_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User login history table';

-- Alert rules table
CREATE TABLE IF NOT EXISTS `alert_rules` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Rule name',
    `rule_type` ENUM('SYSTEM', 'APPLICATION', 'CUSTOM') NOT NULL COMMENT 'Rule type',
    `target_type` ENUM('SERVER', 'APPLICATION', 'SERVICE') NOT NULL COMMENT 'Target type',
    `target_id` BIGINT UNSIGNED NULL COMMENT 'Target ID',
    `metric_name` VARCHAR(64) NULL COMMENT 'Metric name',
    `condition_expression` TEXT NOT NULL COMMENT 'Condition expression',
    `is_enabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Whether enabled',
    `owner_id` BIGINT UNSIGNED NOT NULL COMMENT 'Owner ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_rule_type` (`rule_type`),
    KEY `idx_target_type` (`target_type`),
    KEY `idx_target_id` (`target_id`),
    KEY `idx_owner_id` (`owner_id`),
    KEY `idx_is_enabled` (`is_enabled`),
    CONSTRAINT `fk_alert_rules_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Alert rules table';

-- Alert records table
CREATE TABLE IF NOT EXISTS `alert_records` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `alert_rule_id` BIGINT UNSIGNED NOT NULL COMMENT 'Alert rule ID',
    `alert_id` VARCHAR(64) NOT NULL UNIQUE COMMENT 'Alert ID',
    `title` VARCHAR(255) NOT NULL COMMENT 'Alert title',
    `description` TEXT NULL COMMENT 'Alert description',
    `status` ENUM('FIRING', 'RESOLVED', 'ACKNOWLEDGED') NOT NULL DEFAULT 'FIRING' COMMENT 'Status',
    `fired_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Fired time',
    `resolved_at` DATETIME NULL COMMENT 'Resolved time',
    `acknowledged_at` DATETIME NULL COMMENT 'Acknowledged time',
    `acknowledged_by` BIGINT UNSIGNED NULL COMMENT 'Acknowledged by ID',
    `resolution_note` TEXT NULL COMMENT 'Resolution note',
    `notification_sent` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Notification sent',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_alert_rule_id` (`alert_rule_id`),
    KEY `idx_status` (`status`),
    KEY `idx_fired_at` (`fired_at`),
    KEY `idx_acknowledged_by` (`acknowledged_by`),
    CONSTRAINT `fk_alert_records_rule` FOREIGN KEY (`alert_rule_id`) REFERENCES `alert_rules` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_alert_records_acknowledged_by` FOREIGN KEY (`acknowledged_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Alert records table';

-- Notification templates table
CREATE TABLE IF NOT EXISTS `notification_templates` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL COMMENT 'Template name',
    `type` ENUM('EMAIL', 'SMS', 'WEBHOOK', 'PUSH') NOT NULL COMMENT 'Notification type',
    `subject` VARCHAR(255) NULL COMMENT 'Notification subject',
    `content` TEXT NOT NULL COMMENT 'Notification content template',
    `variables` JSON NULL COMMENT 'Template variables',
    `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether system template',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Status: 0-disabled, 1-enabled',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Notification templates table';

-- Notification records table
CREATE TABLE IF NOT EXISTS `notification_records` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `template_id` BIGINT UNSIGNED NULL COMMENT 'Template ID',
    `type` ENUM('EMAIL', 'SMS', 'WEBHOOK', 'PUSH') NOT NULL COMMENT 'Notification type',
    `recipient` VARCHAR(255) NOT NULL COMMENT 'Recipient',
    `subject` VARCHAR(255) NULL COMMENT 'Notification subject',
    `content` TEXT NOT NULL COMMENT 'Notification content',
    `status` ENUM('PENDING', 'SENDING', 'SUCCESS', 'FAILED') NOT NULL DEFAULT 'PENDING' COMMENT 'Send status',
    `sent_at` DATETIME NULL COMMENT 'Send time',
    `error_msg` TEXT NULL COMMENT 'Error message',
    `retry_count` INT NOT NULL DEFAULT 0 COMMENT 'Retry count',
    `reference_id` VARCHAR(64) NULL COMMENT 'Reference ID',
    `reference_type` VARCHAR(32) NULL COMMENT 'Reference type',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`),
    KEY `idx_reference` (`reference_type`, `reference_id`),
    CONSTRAINT `fk_notification_records_template` FOREIGN KEY (`template_id`) REFERENCES `notification_templates` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Notification records table';

-- Messages table
CREATE TABLE IF NOT EXISTS `notifications` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(255) NOT NULL COMMENT 'Message title',
    `content` TEXT NOT NULL COMMENT 'Message content',
    `type` ENUM('SYSTEM', 'ALERT', 'TASK', 'CUSTOM') NOT NULL COMMENT 'Message type',
    `level` ENUM('INFO', 'WARNING', 'ERROR', 'SUCCESS') NOT NULL DEFAULT 'INFO' COMMENT 'Message level',
    `sender_id` BIGINT UNSIGNED NULL COMMENT 'Sender ID',
    `target_type` ENUM('USER', 'GROUP', 'ALL') NOT NULL COMMENT 'Target type',
    `target_ids` JSON NULL COMMENT 'Target ID list',
    `channels` JSON NULL COMMENT 'Send channels',
    `template_id` BIGINT UNSIGNED NULL COMMENT 'Message template ID',
    `variables` JSON NULL COMMENT 'Template variables',
    `sent_at` DATETIME NULL COMMENT 'Send time',
    `status` ENUM('PENDING', 'SENDING', 'SUCCESS', 'FAILED') NOT NULL DEFAULT 'PENDING' COMMENT 'Send status',
    `error_msg` TEXT NULL COMMENT 'Error message',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    PRIMARY KEY (`id`),
    KEY `idx_type` (`type`),
    KEY `idx_level` (`level`),
    KEY `idx_sender_id` (`sender_id`),
    KEY `idx_target_type` (`target_type`),
    KEY `idx_status` (`status`),
    KEY `idx_template_id` (`template_id`),
    CONSTRAINT `fk_notifications_sender` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`) ON DELETE SET NULL,
    CONSTRAINT `fk_notifications_template` FOREIGN KEY (`template_id`) REFERENCES `notification_templates` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Messages table';

-- User notification records table
CREATE TABLE IF NOT EXISTS `user_notifications` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `notification_id` BIGINT UNSIGNED NOT NULL COMMENT 'Notification ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
    `is_read` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether read',
    `read_at` DATETIME NULL COMMENT 'Read time',
    `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Whether deleted',
    `deleted_at` DATETIME NULL COMMENT 'Delete time',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_notification_user` (`notification_id`, `user_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_is_read` (`is_read`),
    KEY `idx_is_deleted` (`is_deleted`),
    CONSTRAINT `fk_user_notifications_notification` FOREIGN KEY (`notification_id`) REFERENCES `notifications` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_user_notifications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User notification records table';

-- Audit logs table
CREATE TABLE IF NOT EXISTS `audit_logs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NULL COMMENT 'User ID',
    `username` VARCHAR(64) NULL COMMENT 'Username',
    `action` VARCHAR(32) NOT NULL COMMENT 'Action',
    `module` VARCHAR(32) NOT NULL COMMENT 'Module name',
    `resource_type` VARCHAR(32) NULL COMMENT 'Resource type',
    `resource_id` BIGINT UNSIGNED NULL COMMENT 'Resource ID',
    `resource_name` VARCHAR(64) NULL COMMENT 'Resource name',
    `description` TEXT NULL COMMENT 'Operation description',
    `ip_address` VARCHAR(45) NULL COMMENT 'IP address',
    `user_agent` VARCHAR(255) NULL COMMENT 'User agent',
    `request_method` VARCHAR(10) NULL COMMENT 'Request method',
    `request_url` VARCHAR(255) NULL COMMENT 'Request URL',
    `request_params` JSON NULL COMMENT 'Request parameters',
    `response_status` INT NULL COMMENT 'Response status',
    `response_time` INT NULL COMMENT 'Response time (milliseconds)',
    `success` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Whether successful',
    `error_message` TEXT NULL COMMENT 'Error message',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_action` (`action`),
    KEY `idx_module` (`module`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_success` (`success`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Audit logs table';

-- Tag management
CREATE TABLE tags (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL COMMENT 'Tag name',
    color VARCHAR(16) COMMENT 'Tag color',
    description TEXT COMMENT 'Tag description',
    created_by BIGINT DEFAULT 0 COMMENT 'Creator user ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_name (name),
    INDEX idx_created_by (created_by),
    UNIQUE KEY ux_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Tag management table';

-- Tag references
CREATE TABLE taggings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tag_id BIGINT NOT NULL COMMENT 'Tag ID',
    resource_id BIGINT NOT NULL COMMENT 'Resource ID',
    created_by BIGINT DEFAULT 0 COMMENT 'Association creator',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tag_id (tag_id),
    INDEX idx_resource_id (resource_id),
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Tag-Resource association table';

-- Secret references table
CREATE TABLE IF NOT EXISTS `secret_references` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `secret_id` BIGINT UNSIGNED NOT NULL,
    `resource_code` VARCHAR(64) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_secret_references_secret_id` (`secret_id`),
    KEY `idx_secret_references_resource_code` (`resource_code`),
    UNIQUE KEY `idx_secret_references_secret_resource` (`secret_id`, `resource_code`),
    CONSTRAINT `fk_secret_references_secret` FOREIGN KEY (`secret_id`) REFERENCES `secret_keys` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Secret references table';
-- ========================================
-- Index optimization
-- ========================================

-- Add composite indexes for frequently queried fields
CREATE INDEX `idx_projects_owner_status` ON `projects` (`owner_id`, `status`);
CREATE INDEX `idx_app_instances_project_status` ON `app_instances` (`project_id`, `status`);
CREATE INDEX `idx_app_instances_server_status` ON `app_instances` (`server_id`, `status`);
CREATE INDEX `idx_servers_owner_status` ON `servers` (`owner_id`, `status`);
CREATE INDEX `idx_workflows_project_status` ON `workflows` (`project_id`, `status`);
CREATE INDEX `idx_workflows_owner_status` ON `workflows` (`owner_id`, `status`);
CREATE INDEX `idx_alert_records_rule_status` ON `alert_records` (`alert_rule_id`, `status`);
CREATE INDEX `idx_notifications_target_status` ON `notifications` (`target_type`, `status`);
CREATE INDEX `idx_user_notifications_user_read` ON `user_notifications` (`user_id`, `is_read`);

-- Add indexes for time range queries
CREATE INDEX `idx_audit_logs_created_user` ON `audit_logs` (`created_at`, `user_id`);
CREATE INDEX `idx_project_activities_created_project` ON `project_activities` (`created_at`, `project_id`);
CREATE INDEX `idx_workflow_executions_created_status` ON `workflow_executions` (`created_at`, `status`);

-- Add indexes for full-text search (if needed)
-- ALTER TABLE `app_store_templates` ADD FULLTEXT(`name`, `description`);
-- ALTER TABLE `app_store_wishlists` ADD FULLTEXT(`name`, `description`);

SET FOREIGN_KEY_CHECKS = 1;
