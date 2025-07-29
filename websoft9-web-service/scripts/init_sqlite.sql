-- Websoft9 SQLite Database Initialization Script
-- Generated from Database Design Document

-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- ========================================
-- 3.1 应用商店 (App Store)
-- ========================================

-- 应用分类表
CREATE TABLE app_store_categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(32) NOT NULL UNIQUE,
    parent_id INTEGER REFERENCES app_store_categories(id),
    icon VARCHAR(255),
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用模板表
CREATE TABLE app_store_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(32) NOT NULL UNIQUE,
    category_id INTEGER NOT NULL REFERENCES app_store_categories(id),
    version VARCHAR(32) NOT NULL,
    icon VARCHAR(255),
    description TEXT,
    official_url VARCHAR(255),
    source_url VARCHAR(255),
    compose_template TEXT NOT NULL,
    download_count INTEGER DEFAULT 0,
    star_count INTEGER DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0.00,
    is_official INTEGER DEFAULT 0,
    is_featured INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用心愿单表
CREATE TABLE app_store_wishlists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    version VARCHAR(32),
    source_url VARCHAR(255),
    description TEXT,
    reward_amount DECIMAL(10,2) DEFAULT 0.00,
    priority INTEGER DEFAULT 3,
    status VARCHAR(20) DEFAULT 'PENDING',
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    vote_count INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,
    submitter_id INTEGER NOT NULL,
    completed_at DATETIME,
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用评价表
CREATE TABLE app_store_reviews (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL REFERENCES app_store_templates(id),
    user_id INTEGER NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    content TEXT,
    tags TEXT, -- JSON format
    is_helpful INTEGER DEFAULT 0,
    helpful_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用收藏表
CREATE TABLE app_store_favorites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL REFERENCES app_store_templates(id),
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(template_id, user_id)
);

-- 应用点赞表
CREATE TABLE app_store_stars (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL REFERENCES app_store_templates(id),
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(template_id, user_id)
);

-- 应用部署记录表
CREATE TABLE app_deployments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    deployment_id VARCHAR(64) NOT NULL UNIQUE,
    template_id INTEGER REFERENCES app_store_templates(id),
    app_instance_id INTEGER,
    server_id INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING',
    progress INTEGER DEFAULT 0,
    estimated_time INTEGER DEFAULT 0,
    start_time DATETIME,
    end_time DATETIME,
    error_message TEXT,
    deployment_log TEXT,
    config_data TEXT, -- JSON format
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ========================================
-- 3.2 工作空间 (Workspace)
-- ========================================

-- 用户文件表
CREATE TABLE user_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(1024) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'FILE',
    size INTEGER DEFAULT 0,
    mime_type VARCHAR(128),
    download_count INTEGER DEFAULT 0,
    parent_id INTEGER REFERENCES user_files(id),
    storage_path VARCHAR(1024),
    checksum VARCHAR(64),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 工作流表
CREATE TABLE workflows (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(32) NOT NULL,
    description TEXT,
    definition TEXT NOT NULL, -- JSON format
    status VARCHAR(20) DEFAULT 'DRAFT',
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 工作流任务表
CREATE TABLE workflow_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    workflow_id INTEGER NOT NULL REFERENCES workflows(id),
    schedule_type VARCHAR(20) DEFAULT 'MANUAL',
    cron_expression VARCHAR(100),
    status VARCHAR(20) DEFAULT 'DEFAULT',
    next_run_at DATETIME,
    last_run_at DATETIME,
    run_count INTEGER DEFAULT 0,
    success_count INTEGER DEFAULT 0,
    failure_count INTEGER DEFAULT 0,
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 工作流执行历史表
CREATE TABLE workflow_executions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL REFERENCES workflow_tasks(id),
    execution_id VARCHAR(64) NOT NULL UNIQUE,
    status VARCHAR(20) DEFAULT 'PENDING',
    trigger_type VARCHAR(20) DEFAULT 'MANUAL',
    trigger_by INTEGER,
    start_time DATETIME,
    end_time DATETIME,
    duration INTEGER DEFAULT 0,
    error_message TEXT,
    execution_log TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ========================================
-- 3.3 资源管理 (Resource Management)
-- ========================================

-- 数据库连接表
CREATE TABLE database_connections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    db_type VARCHAR(20) NOT NULL,
    host VARCHAR(255) NOT NULL,
    port INTEGER NOT NULL,
    database VARCHAR(64),
    username VARCHAR(64),
    password VARCHAR(255),
    ssl_enabled INTEGER DEFAULT 0,
    connection_timeout INTEGER DEFAULT 30,
    max_connections INTEGER DEFAULT 10,
    status VARCHAR(20) DEFAULT 'CONNECTED',
    version VARCHAR(32),
    charset VARCHAR(32),
    description TEXT,
    last_connected_at DATETIME,
    owner_id INTEGER NOT NULL,
    resource_group_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 服务器表
CREATE TABLE servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    internal_ip VARCHAR(45),
    ssh_port INTEGER DEFAULT 22,
    os_type VARCHAR(32) NOT NULL,
    os_version VARCHAR(64),
    kernel_version VARCHAR(64),
    cpu_cores INTEGER DEFAULT 0,
    memory_total INTEGER DEFAULT 0,
    disk_total INTEGER DEFAULT 0,
    architecture VARCHAR(16),
    status VARCHAR(20) DEFAULT 'UNKNOWN',
    last_heartbeat_at DATETIME,
    resource_group_id INTEGER,
    owner_id INTEGER NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 客户端表
CREATE TABLE server_agents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    server_id INTEGER NOT NULL REFERENCES servers(id),
    container_id VARCHAR(64),
    agent_ip VARCHAR(45),
    agent_port INTEGER DEFAULT 22,
    version VARCHAR(32),
    status VARCHAR(20) DEFAULT 'UNKNOWN',
    last_heartbeat_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用实例表
CREATE TABLE app_instances (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    template_id INTEGER NOT NULL,
    server_id INTEGER NOT NULL REFERENCES servers(id),
    container_id VARCHAR(64),
    container_name VARCHAR(64),
    image_name VARCHAR(255),
    image_tag VARCHAR(100),
    status VARCHAR(20) DEFAULT 'DEFAULT',
    started_at DATETIME,
    stopped_at DATETIME,
    resource_group_id INTEGER,
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ========================================
-- 3.4 安全管控 (Security Management)
-- ========================================

-- SSL证书表
CREATE TABLE ssl_certificates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    domain VARCHAR(255) NOT NULL,
    certificate_type VARCHAR(20) DEFAULT 'LETS_ENCRYPT',
    certificate_data TEXT NOT NULL,
    private_key_data TEXT NOT NULL,
    certificate_chain TEXT,
    issuer VARCHAR(255),
    subject VARCHAR(255),
    serial_number VARCHAR(64),
    not_before DATETIME,
    not_after DATETIME,
    auto_renew INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'PENDING',
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 密钥管理表
CREATE TABLE secret_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    key_type VARCHAR(20) NOT NULL,
    encrypted_value TEXT NOT NULL,
    description TEXT,
    custom_fields TEXT, -- JSON format
    authorized_users TEXT, -- JSON format
    authorized_groups TEXT, -- JSON format
    expires_at DATETIME,
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用网关表
CREATE TABLE app_gateways (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    server_id INTEGER NOT NULL REFERENCES servers(id),
    container_id VARCHAR(64),
    description TEXT,
    status VARCHAR(20) DEFAULT 'DEFAULT',
    started_at DATETIME,
    stopped_at DATETIME,
    resource_group_id INTEGER,
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用网关发布表
CREATE TABLE app_gateways_publishes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_instance_id INTEGER NOT NULL,
    app_gateway_id INTEGER NOT NULL REFERENCES app_gateways(id),
    service_domain VARCHAR(255) NOT NULL,
    service_port INTEGER DEFAULT 8080,
    alert_rule_id INTEGER,
    limit_rules TEXT,
    health_check_enabled INTEGER DEFAULT 1,
    audit_log_enabled INTEGER DEFAULT 1,
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 审计日志表
CREATE TABLE audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    username VARCHAR(64),
    action VARCHAR(32) NOT NULL,
    module VARCHAR(32) NOT NULL,
    resource_type VARCHAR(32),
    resource_id INTEGER,
    resource_name VARCHAR(64),
    description TEXT,
    ip_address VARCHAR(45),
    user_agent VARCHAR(255),
    request_method VARCHAR(10),
    request_url VARCHAR(255),
    request_params TEXT, -- JSON format
    response_status INTEGER,
    response_time INTEGER,
    success INTEGER DEFAULT 1,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ========================================
-- 3.5 平台管理 (Platform Management)
-- ========================================

-- 用户组表
CREATE TABLE user_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(32) NOT NULL UNIQUE,
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 用户表
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL REFERENCES user_groups(id),
    username VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(64),
    avatar VARCHAR(255),
    phone VARCHAR(20),
    gender INTEGER DEFAULT 0,
    signature VARCHAR(255),
    status INTEGER DEFAULT 1,
    last_login_at DATETIME,
    last_login_ip VARCHAR(45),
    timezone VARCHAR(64) DEFAULT 'UTC',
    language VARCHAR(10) DEFAULT 'zh-CN',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 角色表
CREATE TABLE roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL UNIQUE,
    code VARCHAR(32) NOT NULL UNIQUE,
    description TEXT,
    is_system INTEGER DEFAULT 0,
    sort_order INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 权限表
CREATE TABLE permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(64) NOT NULL UNIQUE,
    module VARCHAR(32) NOT NULL,
    action VARCHAR(32) NOT NULL,
    resource VARCHAR(64),
    description TEXT,
    is_system INTEGER DEFAULT 0,
    sort_order INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 用户角色关联表
CREATE TABLE user_roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id),
    role_id INTEGER NOT NULL REFERENCES roles(id),
    granted_by INTEGER,
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id)
);

-- 角色权限关联表
CREATE TABLE role_permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL REFERENCES roles(id),
    permission_id INTEGER NOT NULL REFERENCES permissions(id),
    granted_by INTEGER,
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, permission_id)
);

-- 资源组表
CREATE TABLE resource_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(32) NOT NULL UNIQUE,
    description TEXT,
    owner_id INTEGER NOT NULL,
    sort_order INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 系统配置表
CREATE TABLE system_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    config_key VARCHAR(64) NOT NULL UNIQUE,
    config_value TEXT,
    config_type VARCHAR(20) DEFAULT 'STRING',
    category VARCHAR(32) NOT NULL,
    description TEXT,
    is_readonly INTEGER DEFAULT 0,
    is_encrypted INTEGER DEFAULT 0,
    default_value TEXT,
    validation_rules TEXT, -- JSON format
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 告警规则表
CREATE TABLE alert_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    rule_type VARCHAR(20) NOT NULL,
    target_type VARCHAR(20) NOT NULL,
    target_id INTEGER,
    metric_name VARCHAR(64),
    condition_expression TEXT NOT NULL,
    notification_channels TEXT, -- JSON format
    is_enabled INTEGER DEFAULT 1,
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 告警记录表
CREATE TABLE alert_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alert_rule_id INTEGER NOT NULL REFERENCES alert_rules(id),
    alert_id VARCHAR(64) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'FIRING',
    fired_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    resolved_at DATETIME,
    acknowledged_at DATETIME,
    acknowledged_by INTEGER,
    resolution_note TEXT,
    notification_sent INTEGER DEFAULT 0,
    notification_channels TEXT, -- JSON format
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 通知模板表
CREATE TABLE notification_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    type VARCHAR(20) NOT NULL,
    subject VARCHAR(255),
    content TEXT NOT NULL,
    variables TEXT, -- JSON format
    is_system INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 通知记录表
CREATE TABLE notification_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER,
    type VARCHAR(20) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    subject VARCHAR(255),
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING',
    sent_at DATETIME,
    error_msg TEXT,
    retry_count INTEGER DEFAULT 0,
    reference_id VARCHAR(64),
    reference_type VARCHAR(32),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ========================================
-- 索引创建
-- ========================================

-- 应用商店相关索引
CREATE INDEX idx_app_store_templates_category ON app_store_templates(category_id);
CREATE INDEX idx_app_store_templates_status ON app_store_templates(status);
CREATE INDEX idx_app_store_reviews_template ON app_store_reviews(template_id);
CREATE INDEX idx_app_store_reviews_user ON app_store_reviews(user_id);

-- 工作流相关索引
CREATE INDEX idx_workflows_owner ON workflows(owner_id);
CREATE INDEX idx_workflow_tasks_workflow ON workflow_tasks(workflow_id);
CREATE INDEX idx_workflow_executions_task ON workflow_executions(task_id);

-- 资源管理相关索引
CREATE INDEX idx_servers_owner ON servers(owner_id);
CREATE INDEX idx_servers_status ON servers(status);
CREATE INDEX idx_app_instances_server ON app_instances(server_id);
CREATE INDEX idx_app_instances_template ON app_instances(template_id);

-- 用户权限相关索引
CREATE INDEX idx_users_group ON users(group_id);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role_id);
CREATE INDEX idx_role_permissions_role ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission ON role_permissions(permission_id);

-- 审计日志索引
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);

-- ========================================
-- 初始数据插入
-- ========================================

-- 插入默认用户组
INSERT INTO user_groups (name, code, description) VALUES 
('管理员组', 'admin', '系统管理员用户组'),
('普通用户组', 'user', '普通用户组');

-- 插入默认角色
INSERT INTO roles (name, code, description, is_system) VALUES 
('超级管理员', 'super_admin', '系统超级管理员', 1),
('管理员', 'admin', '系统管理员', 1),
('开发者', 'developer', '开发者角色', 1),
('运维人员', 'operator', '运维人员角色', 1),
('普通用户', 'user', '普通用户角色', 1);

-- 插入默认权限
INSERT INTO permissions (name, code, module, action, is_system) VALUES 
('用户管理', 'user.manage', 'user', 'manage', 1),
('角色管理', 'role.manage', 'role', 'manage', 1),
('权限管理', 'permission.manage', 'permission', 'manage', 1),
('应用管理', 'app.manage', 'app', 'manage', 1),
('服务器管理', 'server.manage', 'server', 'manage', 1),
('监控查看', 'monitor.view', 'monitor', 'view', 1),
('系统配置', 'system.config', 'system', 'config', 1);

-- 插入默认应用分类
INSERT INTO app_store_categories (name, code, description) VALUES 
('Web服务', 'web', 'Web服务器和相关应用'),
('数据库', 'database', '各种数据库系统'),
('开发工具', 'development', '开发和构建工具'),
('监控工具', 'monitoring', '系统监控和日志工具'),
('安全工具', 'security', '安全防护工具'),
('其他', 'others', '其他应用');

-- 插入系统配置
INSERT INTO system_configs (config_key, config_value, config_type, category, description) VALUES 
('system.name', 'Websoft9', 'STRING', 'basic', '系统名称'),
('system.version', '1.0.0', 'STRING', 'basic', '系统版本'),
('system.timezone', 'Asia/Shanghai', 'STRING', 'basic', '系统时区'),
('system.language', 'zh-CN', 'STRING', 'basic', '系统语言'),
('security.password_min_length', '6', 'INTEGER', 'security', '密码最小长度'),
('security.session_timeout', '3600', 'INTEGER', 'security', '会话超时时间(秒)');

-- 创建触发器用于自动更新 updated_at 字段
CREATE TRIGGER update_app_store_categories_updated_at 
    AFTER UPDATE ON app_store_categories
    BEGIN
        UPDATE app_store_categories SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_app_store_templates_updated_at 
    AFTER UPDATE ON app_store_templates
    BEGIN
        UPDATE app_store_templates SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_users_updated_at 
    AFTER UPDATE ON users
    BEGIN
        UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_servers_updated_at 
    AFTER UPDATE ON servers
    BEGIN
        UPDATE servers SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;