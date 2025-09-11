-- Websoft9 SQLite Database Initialization Script V1.1
-- Generated from Database Design Document V1.1
-- Compatible with SQLite 3.35+

-- Enable foreign key constraints and other optimizations
-- Temporarily disable foreign keys for data insertion
PRAGMA foreign_keys = OFF;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = 10000;
PRAGMA temp_store = memory;

-- ========================================
-- 3.1 平台主页 (Platform Home)
-- ========================================

-- 应用快捷导航表
CREATE TABLE IF NOT EXISTS app_shortcuts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_instance_id INTEGER NOT NULL,
    name VARCHAR(64),
    description VARCHAR(255),
    icon VARCHAR(255),
    sort_order INTEGER DEFAULT 0,
    access_count INTEGER DEFAULT 0,
    last_accessed DATETIME,
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (app_instance_id) REFERENCES app_instances(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


-- ========================================
-- 3.2 项目管理 (Project Management)
-- ========================================

-- 项目表
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    identifier VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    tags TEXT, -- JSON format
    icon VARCHAR(255),
    status VARCHAR(20) DEFAULT 'NORMAL',
    owner_id INTEGER NOT NULL,
    default_resource_group VARCHAR(50) DEFAULT 'default',
    default_timezone VARCHAR(50) DEFAULT 'Asia/Shanghai',
    log_retention_days INTEGER DEFAULT 30,
    backup_strategy VARCHAR(20) DEFAULT 'daily',
    access_control VARCHAR(50) DEFAULT 'members',
    api_access INTEGER DEFAULT 1,
    audit_enabled INTEGER DEFAULT 1,
    last_activity_at DATETIME,
    archived_at DATETIME,
    archived_by INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (archived_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 项目成员表
CREATE TABLE IF NOT EXISTS project_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role VARCHAR(50) NOT NULL,
    is_admin INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    invited_by INTEGER,
    invited_at DATETIME,
    joined_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (invited_by) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(project_id, user_id)
);

-- 项目环境变量表
CREATE TABLE IF NOT EXISTS project_environments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    value TEXT,
    type VARCHAR(20) DEFAULT 'normal',
    description VARCHAR(500),
    scope VARCHAR(20) DEFAULT 'global',
    scope_target VARCHAR(100),
    is_encrypted INTEGER DEFAULT 0,
    created_by INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT
);

-- 项目文件表
CREATE TABLE IF NOT EXISTS project_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(1024) NOT NULL,
    type VARCHAR(20) DEFAULT 'FILE',
    size INTEGER DEFAULT 0,
    mime_type VARCHAR(128),
    download_count INTEGER DEFAULT 0,
    parent_id INTEGER,
    storage_path VARCHAR(1024),
    checksum VARCHAR(64),
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES project_files(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

-- 项目活动记录表
CREATE TABLE IF NOT EXISTS project_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    user_id INTEGER,
    username VARCHAR(64),
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50),
    resource_id INTEGER,
    resource_name VARCHAR(255),
    description TEXT,
    metadata TEXT, -- JSON format
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 工作流表
CREATE TABLE IF NOT EXISTS workflows (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(32) NOT NULL,
    description TEXT,
    definition TEXT NOT NULL, -- JSON format
    status VARCHAR(20) DEFAULT 'DRAFT',
    project_id INTEGER NOT NULL REFERENCES projects(id),
    owner_id INTEGER NOT NULL REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 工作流任务表
CREATE TABLE IF NOT EXISTS workflow_tasks (
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
CREATE TABLE IF NOT EXISTS workflow_executions (
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

-- 资源组表
CREATE TABLE IF NOT EXISTS resource_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(32) NOT NULL UNIQUE,
    description TEXT,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    sort_order INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1, -- 0-禁用，1-启用
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 数据库连接表
CREATE TABLE IF NOT EXISTS database_connections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    db_type VARCHAR(20) NOT NULL, -- mysql, postgresql, redis
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
    resource_group_id INTEGER REFERENCES resource_groups(id) ON DELETE SET NULL,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 服务器表
CREATE TABLE IF NOT EXISTS servers (
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
    resource_group_id INTEGER REFERENCES resource_groups(id),
    owner_id INTEGER NOT NULL REFERENCES users(id),
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 客户端表
CREATE TABLE IF NOT EXISTS server_agents (
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
CREATE TABLE IF NOT EXISTS app_instances (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL REFERENCES projects(id),
    name VARCHAR(64) NOT NULL,
    template_id INTEGER NOT NULL REFERENCES app_store_templates(id),
    server_id INTEGER NOT NULL REFERENCES servers(id),
    container_id VARCHAR(64),
    container_name VARCHAR(64),
    image_name VARCHAR(255),
    image_tag VARCHAR(100),
    status VARCHAR(20) DEFAULT 'DEFAULT',
    started_at DATETIME,
    stopped_at DATETIME,
    owner_id INTEGER NOT NULL REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- SSL证书表
CREATE TABLE IF NOT EXISTS ssl_certificates (
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
CREATE TABLE IF NOT EXISTS secret_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    key_type VARCHAR(20) NOT NULL, -- API_KEY, DATABASE, SSH, CERTIFICATE, CUSTOM
    encrypted_value TEXT NOT NULL,
    description TEXT,
    custom_fields TEXT, -- JSON format
    expires_at DATETIME,
    resource_group_id INTEGER REFERENCES resource_groups(id) ON DELETE SET NULL,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用网关表
CREATE TABLE IF NOT EXISTS app_gateways (
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
CREATE TABLE IF NOT EXISTS app_gateways_publishes (
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

-- 应用网关访问控制规则表
CREATE TABLE IF NOT EXISTS app_gateway_access_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    gateway_id INTEGER NOT NULL REFERENCES app_gateways(id),
    rule_name VARCHAR(64) NOT NULL,
    rule_type VARCHAR(20) NOT NULL,
    limit_count INTEGER NOT NULL,
    time_window INTEGER NOT NULL,
    target_path VARCHAR(255),
    target_ip VARCHAR(45),
    action VARCHAR(20) DEFAULT 'BLOCK',
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ========================================
-- 3.3 应用市场 (App Store)
-- ========================================

-- 应用分类表
CREATE TABLE IF NOT EXISTS app_store_categories (
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
CREATE TABLE IF NOT EXISTS app_store_templates (
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
CREATE TABLE IF NOT EXISTS app_store_wishlists (
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
CREATE TABLE IF NOT EXISTS app_store_reviews (
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
CREATE TABLE IF NOT EXISTS app_store_favorites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL REFERENCES app_store_templates(id),
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(template_id, user_id)
);

-- 应用点赞表
CREATE TABLE IF NOT EXISTS app_store_stars (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL REFERENCES app_store_templates(id),
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(template_id, user_id)
);

-- 应用举报表
CREATE TABLE IF NOT EXISTS app_store_reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL REFERENCES app_store_templates(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    reason VARCHAR(20) NOT NULL,
    detail TEXT,
    status VARCHAR(20) DEFAULT 'PENDING',
    handled_by INTEGER REFERENCES users(id),
    handled_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用下载记录表
CREATE TABLE IF NOT EXISTS app_store_downloads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL REFERENCES app_store_templates(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    ip_address VARCHAR(45),
    user_agent VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用心愿单评论表
CREATE TABLE IF NOT EXISTS app_store_wishlist_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wishlist_id INTEGER NOT NULL REFERENCES app_store_wishlists(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    parent_id INTEGER REFERENCES app_store_wishlist_comments(id),
    content TEXT NOT NULL,
    like_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用心愿单投票表
CREATE TABLE IF NOT EXISTS app_store_wishlist_votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wishlist_id INTEGER NOT NULL REFERENCES app_store_wishlists(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(wishlist_id, user_id)
);

-- 应用心愿单点赞表
CREATE TABLE IF NOT EXISTS app_store_wishlist_likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wishlist_id INTEGER NOT NULL REFERENCES app_store_wishlists(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(wishlist_id, user_id)
);

-- 应用心愿单举报表
CREATE TABLE IF NOT EXISTS app_store_wishlist_reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wishlist_id INTEGER NOT NULL REFERENCES app_store_wishlists(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    reason VARCHAR(20) NOT NULL,
    detail TEXT,
    status VARCHAR(20) DEFAULT 'PENDING',
    handled_by INTEGER REFERENCES users(id),
    handled_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用部署记录表
CREATE TABLE IF NOT EXISTS app_deployments (
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
-- 3.4 平台管理 (Platform Management)
-- ========================================

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
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
CREATE TABLE IF NOT EXISTS roles (
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
CREATE TABLE IF NOT EXISTS permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_code VARCHAR(64) REFERENCES permissions(code),
    scope VARCHAR(64) NOT NULL,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(64) NOT NULL UNIQUE,
    module VARCHAR(32) NOT NULL,
    action VARCHAR(32) NOT NULL,
    resource VARCHAR(64),
    element VARCHAR(64),
    description TEXT,
    is_system INTEGER DEFAULT 0,
    is_menu INTEGER DEFAULT 0,
    sort_order INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    created_by INTEGER REFERENCES users(id),
    updated_by INTEGER REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id),
    role_id INTEGER NOT NULL REFERENCES roles(id),
    granted_by INTEGER REFERENCES users(id),
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id)
);

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL REFERENCES roles(id),
    permission_code VARCHAR(64) NOT NULL REFERENCES permissions(code),
    granted_by INTEGER,
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, permission_code)
);

-- API访问令牌表
CREATE TABLE IF NOT EXISTS api_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    scopes TEXT, -- JSON format
    description TEXT,
    last_used_at DATETIME,
    last_used_ip VARCHAR(45),
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 用户双因子认证表
CREATE TABLE IF NOT EXISTS user_two_factor (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id),
    method VARCHAR(32) NOT NULL,
    secret VARCHAR(255),
    backup_codes TEXT, -- JSON format
    email VARCHAR(255),
    enabled INTEGER DEFAULT 0,
    verified_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, method)
);

-- 用户个人中心配置表
CREATE TABLE IF NOT EXISTS user_profile (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id),
    category VARCHAR(64) DEFAULT 'general',
    config_key VARCHAR(64) NOT NULL,
    config_value TEXT,
    description TEXT,
    is_readonly INTEGER DEFAULT 0,
    is_encrypted INTEGER DEFAULT 0,
    default_value TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, config_key)
);

-- 用户登录历史表
CREATE TABLE IF NOT EXISTS user_login_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id),
    ip_address VARCHAR(45),
    user_agent VARCHAR(255),
    location VARCHAR(100),
    device VARCHAR(100),
    browser VARCHAR(100),
    login_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    logout_time DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Webhook配置表
CREATE TABLE IF NOT EXISTS webhook_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    url VARCHAR(255) NOT NULL,
    secret VARCHAR(255),
    events TEXT NOT NULL, -- JSON format
    headers TEXT, -- JSON format
    timeout INTEGER DEFAULT 30,
    retry_count INTEGER DEFAULT 3,
    is_active INTEGER DEFAULT 1,
    owner_id INTEGER NOT NULL REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Webhook执行日志表
CREATE TABLE IF NOT EXISTS webhook_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    webhook_id INTEGER NOT NULL REFERENCES webhook_configs(id),
    event_type VARCHAR(32) NOT NULL,
    payload TEXT NOT NULL, -- JSON format
    request_id VARCHAR(64) NOT NULL,
    status_code INTEGER,
    response_body TEXT,
    response_time INTEGER DEFAULT 0,
    retry_count INTEGER DEFAULT 0,
    success INTEGER DEFAULT 0,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 平台更新记录表
CREATE TABLE IF NOT EXISTS platform_updates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version VARCHAR(32) NOT NULL,
    changelog TEXT,
    download_url VARCHAR(255),
    file_size INTEGER DEFAULT 0,
    checksum VARCHAR(64),
    status VARCHAR(20) DEFAULT 'PENDING',
    started_at DATETIME,
    completed_at DATETIME,
    error_message TEXT,
    backup_path VARCHAR(1024),
    updated_by INTEGER REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 容器集群节点表
CREATE TABLE IF NOT EXISTS docker_swarm_nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id VARCHAR(64) NOT NULL UNIQUE,
    hostname VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    role VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'ready',
    availability VARCHAR(20) DEFAULT 'active',
    engine_version VARCHAR(32),
    labels TEXT, -- JSON format
    resources TEXT, -- JSON format
    joined_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 容器镜像表
CREATE TABLE IF NOT EXISTS docker_images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    image_id VARCHAR(64) NOT NULL UNIQUE,
    repository VARCHAR(255) NOT NULL,
    tag VARCHAR(100) NOT NULL,
    digest VARCHAR(128),
    size INTEGER DEFAULT 0,
    architecture VARCHAR(32),
    os VARCHAR(32),
    status VARCHAR(20) DEFAULT 'AVAILABLE',
    usage_count INTEGER DEFAULT 0,
    last_used_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
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
CREATE TABLE IF NOT EXISTS alert_rules (
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
CREATE TABLE IF NOT EXISTS alert_records (
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
CREATE TABLE IF NOT EXISTS notification_templates (
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
CREATE TABLE IF NOT EXISTS notification_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER REFERENCES notification_templates(id),
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

-- 通知消息表
CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(20) NOT NULL,
    level VARCHAR(20) DEFAULT 'INFO',
    sender_id INTEGER REFERENCES users(id),
    target_type VARCHAR(20) NOT NULL,
    target_ids TEXT, -- JSON format
    channels TEXT, -- JSON format
    template_id INTEGER REFERENCES notification_templates(id),
    variables TEXT, -- JSON format
    sent_at DATETIME,
    status VARCHAR(20) DEFAULT 'PENDING',
    error_msg TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 用户通知记录表
CREATE TABLE IF NOT EXISTS user_notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    notification_id INTEGER NOT NULL REFERENCES notifications(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    is_read INTEGER DEFAULT 0,
    read_at DATETIME,
    is_deleted INTEGER DEFAULT 0,
    deleted_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(notification_id, user_id)
);

-- 审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
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
-- 初始数据插入
-- ========================================

-- 插入默认角色
INSERT OR IGNORE INTO `roles` (`name`, `code`, `description`, `is_system`, `status`) VALUES
('超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 1, 1),
('系统管理员', 'admin', '系统管理员，负责平台管理', 1, 1),
('项目管理员', 'project_admin', '项目管理员，负责项目管理', 1, 1),
('开发者', 'developer', '开发者角色，负责应用开发', 1, 1),
('运维人员', 'operator', '运维人员角色，负责系统运维', 1, 1),
('普通用户', 'user', '普通用户角色，基础权限', 1, 1);

-- 插入默认权限
INSERT OR IGNORE INTO `permissions` (`parent_code`,`scope`,`name`,`code`,`module`,`action`,`resource`,`element`,`description`,`is_system`,`is_menu`,`sort_order`,`status`,`created_by`,`updated_by`,`created_at`,`updated_at`) VALUES
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
	 ('e06e84c7-fe48-4d0d-9dbb-cbe65cf5d936','platform','个人空间','08d5484f-ee91-47ba-99f9-ae90c8d5684d','personal_folder','*',NULL,NULL,'个人空间',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('08d5484f-ee91-47ba-99f9-ae90c8d5684d','platform','个人空间查询','af309229-e1c0-4053-b8c3-e35ef0b89107','personal_folder','query',NULL,NULL,'个人文件夹查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('08d5484f-ee91-47ba-99f9-ae90c8d5684d','platform','个人空间创建','593b5fb7-4146-4f17-9948-bf2827c4cd94','personal_folder','create',NULL,NULL,'个人文件夹创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('08d5484f-ee91-47ba-99f9-ae90c8d5684d','platform','个人空间修改','67e63e8c-17b8-49d7-b6cd-6618614d3d5c','personal_folder','update',NULL,NULL,'个人文件夹修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('08d5484f-ee91-47ba-99f9-ae90c8d5684d','platform','个人空间删除','d8e80678-a03a-42c3-99d8-29ed4b9032ee','personal_folder','delete',NULL,NULL,'个人文件夹删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('08d5484f-ee91-47ba-99f9-ae90c8d5684d','platform','个人空间上传','f47d800e-d504-4de4-946c-60add0123f8d','personal_folder','create',NULL,NULL,'个人文件夹上传权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('08d5484f-ee91-47ba-99f9-ae90c8d5684d','platform','个人空间下载','89a1752f-7fda-4949-a241-f5704f713eeb','personal_folder','query',NULL,NULL,'个人文件夹下载权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
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
	 ('f9ea8516-3350-44cc-adb5-10cf0e5f8e70','platform','角色详情查询','8394ae5f-c269-4a9a-8c4f-ffeb32b0a703','role','query','/roles/*',NULL,'角色详情查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
   ('cba5a624-afa7-4fff-a1f1-f6f791e5589a','platform','权限管理','1de209de-9182-4514-8016-4fc3401b7968','permission','*','/permissions',NULL,'权限管理全部权限',1,1,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1de209de-9182-4514-8016-4fc3401b7968','platform','权限管理查询','095cba59-25e0-4d6d-8739-4c537ade0b71','permission','query','/permissions',NULL,'权限管理查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1de209de-9182-4514-8016-4fc3401b7968','platform','权限管理创建','4db0477f-ef96-40f3-9853-3cb2df4b0901','permission','create','/permissions',NULL,'权限管理创建权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1de209de-9182-4514-8016-4fc3401b7968','platform','权限管理修改','cfcf331f-a33d-4521-a7f7-f0fbece23fec','permission','update','/permissions/*',NULL,'权限管理修改权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1de209de-9182-4514-8016-4fc3401b7968','platform','权限管理删除','9dbd8b47-d475-4a89-b619-3ea8e087fea6','permission','delete','/permissions/*',NULL,'权限管理删除权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
   ('1de209de-9182-4514-8016-4fc3401b7968','platform','权限详情查询','358f754c-e66f-44f1-bf6c-142ae3610d89','permission','query','/permissions/*',NULL,'权限详情查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
	 ('1de209de-9182-4514-8016-4fc3401b7968','platform','权限树查询','2b7afd04-9020-4ef8-ab2f-5a62a4b49fcf','permission','query','/permissions/tree',NULL,'权限树查询权限',1,0,0,1,1,1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
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
INSERT OR IGNORE INTO `users` (`id`,`username`,`email`,`password_hash`,`nickname`,`avatar`,`phone`,`gender`,`signature`,`status`,`last_login_at`,`last_login_ip`,`timezone`,`language`,`created_at`,`updated_at`) VALUES
	 (1,'admin','admin@websoft9.com','d1a7b27aa60359a6033b046831168286dddc5a83268d4d483d30a09983f9c944','Manager','','',0,'Websoft9 manager',1,NULL,'','Asia/Shanghai','zh-CN','2025-09-01 12:11:19','2025-09-01 12:11:19');

INSERT OR IGNORE INTO `user_roles` (`id`,`user_id`,`role_id`,`granted_by`,`granted_at`,`expires_at`,`status`,`created_at`,`updated_at`) VALUES
	 (1,1,1,1,'2025-09-01 12:11:19',NULL,1,'2025-09-01 12:11:19','2025-09-01 12:11:19');

INSERT OR IGNORE INTO `role_permissions` (`role_id`,`permission_code`,`granted_by`,`granted_at`,`status`,`created_at`,`updated_at`) VALUES
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
         (1, '8394ae5f-c269-4a9a-8c4f-ffeb32b0a703',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '1de209de-9182-4514-8016-4fc3401b7968',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '095cba59-25e0-4d6d-8739-4c537ade0b71',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '4db0477f-ef96-40f3-9853-3cb2df4b0901',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, 'cfcf331f-a33d-4521-a7f7-f0fbece23fec',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '9dbd8b47-d475-4a89-b619-3ea8e087fea6',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '358f754c-e66f-44f1-bf6c-142ae3610d89',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
         (1, '2b7afd04-9020-4ef8-ab2f-5a62a4b49fcf',1,'2025-09-01 12:11:19',1,'2025-09-01 12:11:19','2025-09-01 12:11:19'),
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
INSERT OR IGNORE INTO `app_store_categories` (`name`, `code`, `description`, `sort_order`, `status`) VALUES
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
INSERT OR IGNORE INTO `system_configs` (`config_key`, `config_value`, `config_type`, `category`, `description`, `is_readonly`, `sort_order`) VALUES
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
INSERT OR IGNORE INTO `notification_templates` (`name`, `type`, `subject`, `content`, `variables`, `is_system`, `status`) VALUES
('用户注册通知', 'EMAIL', '欢迎注册 {{system_name}}', '亲爱的 {{username}}，\n\n欢迎注册 {{system_name}}！\n\n您的账户已成功创建，现在可以开始使用我们的服务了。\n\n如有任何问题，请联系我们的支持团队。\n\n祝您使用愉快！\n\n{{system_name}} 团队', '["username", "system_name"]', 1, 1),
('密码重置通知', 'EMAIL', '{{system_name}} 密码重置', '亲爱的 {{username}}，\n\n您的密码已成功重置。\n\n如果这不是您的操作，请立即联系我们的支持团队。\n\n{{system_name}} 团队', '["username", "system_name"]', 1, 1),
('系统告警通知', 'EMAIL', '{{system_name}} 系统告警', '告警标题：{{alert_title}}\n告警描述：{{alert_description}}\n触发时间：{{fired_at}}\n告警级别：{{alert_level}}\n\n请及时处理。', '["alert_title", "alert_description", "fired_at", "alert_level", "system_name"]', 1, 1),
('应用部署成功', 'EMAIL', '应用部署成功通知', '亲爱的 {{username}}，\n\n您的应用 {{app_name}} 已成功部署到服务器 {{server_name}}。\n\n访问地址：{{app_url}}\n部署时间：{{deployed_at}}\n\n{{system_name}} 团队', '["username", "app_name", "server_name", "app_url", "deployed_at", "system_name"]', 1, 1),
('应用部署失败', 'EMAIL', '应用部署失败通知', '亲爱的 {{username}}，\n\n您的应用 {{app_name}} 部署失败。\n\n错误信息：{{error_message}}\n失败时间：{{failed_at}}\n\n请检查配置后重试。\n\n{{system_name}} 团队', '["username", "app_name", "error_message", "failed_at", "system_name"]', 1, 1);

-- ========================================
-- 索引创建
-- ========================================

-- 项目相关索引
CREATE INDEX IF NOT EXISTS idx_projects_owner_status ON projects(owner_id, status);
CREATE INDEX IF NOT EXISTS idx_project_members_project ON project_members(project_id);
CREATE INDEX IF NOT EXISTS idx_project_members_user ON project_members(user_id);
CREATE INDEX IF NOT EXISTS idx_project_environments_project ON project_environments(project_id);
CREATE INDEX IF NOT EXISTS idx_project_files_project ON project_files(project_id);
CREATE INDEX IF NOT EXISTS idx_project_files_parent ON project_files(parent_id);
CREATE INDEX IF NOT EXISTS idx_project_activities_project ON project_activities(project_id);

-- 用户权限相关索引
CREATE INDEX IF NOT EXISTS idx_user_roles_user ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission ON role_permissions(permission_code);

-- 告警通知相关索引
CREATE INDEX IF NOT EXISTS idx_alert_rules_owner ON alert_rules(owner_id);
CREATE INDEX IF NOT EXISTS idx_alert_records_rule ON alert_records(alert_rule_id);
CREATE INDEX IF NOT EXISTS idx_alert_records_status ON alert_records(status);
CREATE INDEX IF NOT EXISTS idx_notifications_target_status ON notifications(target_type, status);
CREATE INDEX IF NOT EXISTS idx_user_notifications_user_read ON user_notifications(user_id, is_read);

-- 时间相关索引
CREATE INDEX IF NOT EXISTS idx_project_activities_created ON project_activities(created_at);
CREATE INDEX IF NOT EXISTS idx_workflow_executions_created ON workflow_executions(created_at);
CREATE INDEX IF NOT EXISTS idx_user_login_history_login_time ON user_login_history(login_time);

-- 应用商店相关索引
CREATE INDEX IF NOT EXISTS idx_app_store_templates_category ON app_store_templates(category_id);
CREATE INDEX IF NOT EXISTS idx_app_store_templates_status ON app_store_templates(status);
CREATE INDEX IF NOT EXISTS idx_app_store_reviews_template ON app_store_reviews(template_id);
CREATE INDEX IF NOT EXISTS idx_app_store_reviews_user ON app_store_reviews(user_id);

-- 工作流相关索引
CREATE INDEX IF NOT EXISTS idx_workflows_project ON workflows(project_id);
CREATE INDEX IF NOT EXISTS idx_workflows_owner ON workflows(owner_id);
CREATE INDEX IF NOT EXISTS idx_workflow_tasks_workflow ON workflow_tasks(workflow_id);
CREATE INDEX IF NOT EXISTS idx_workflow_executions_task ON workflow_executions(task_id);

-- 资源管理相关索引
CREATE INDEX IF NOT EXISTS idx_servers_owner ON servers(owner_id);
CREATE INDEX IF NOT EXISTS idx_servers_status ON servers(status);
CREATE INDEX IF NOT EXISTS idx_app_instances_server ON app_instances(server_id);
CREATE INDEX IF NOT EXISTS idx_app_instances_template ON app_instances(template_id);

-- 审计日志索引
CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at);

-- ========================================
-- 创建触发器用于自动更新 updated_at 字段
-- ========================================

-- 项目相关触发器
CREATE TRIGGER IF NOT EXISTS update_projects_updated_at
    AFTER UPDATE ON projects
    BEGIN
        UPDATE projects SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_project_members_updated_at
    AFTER UPDATE ON project_members
    BEGIN
        UPDATE project_members SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_project_environments_updated_at
    AFTER UPDATE ON project_environments
    BEGIN
        UPDATE project_environments SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_project_files_updated_at
    AFTER UPDATE ON project_files
    BEGIN
        UPDATE project_files SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- 应用商店相关触发器
CREATE TRIGGER IF NOT EXISTS update_app_store_categories_updated_at
    AFTER UPDATE ON app_store_categories
    BEGIN
        UPDATE app_store_categories SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_app_store_templates_updated_at
    AFTER UPDATE ON app_store_templates
    BEGIN
        UPDATE app_store_templates SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_app_store_wishlists_updated_at
    AFTER UPDATE ON app_store_wishlists
    BEGIN
        UPDATE app_store_wishlists SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- 用户相关触发器
CREATE TRIGGER IF NOT EXISTS update_users_updated_at
    AFTER UPDATE ON users
    BEGIN
        UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_roles_updated_at
    AFTER UPDATE ON roles
    BEGIN
        UPDATE roles SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_permissions_updated_at
    AFTER UPDATE ON permissions
    BEGIN
        UPDATE permissions SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- 服务器相关触发器
CREATE TRIGGER IF NOT EXISTS update_servers_updated_at
    AFTER UPDATE ON servers
    BEGIN
        UPDATE servers SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_server_agents_updated_at
    AFTER UPDATE ON server_agents
    BEGIN
        UPDATE server_agents SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- 应用实例相关触发器
CREATE TRIGGER IF NOT EXISTS update_app_instances_updated_at
    AFTER UPDATE ON app_instances
    BEGIN
        UPDATE app_instances SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- 工作流相关触发器
CREATE TRIGGER IF NOT EXISTS update_workflows_updated_at
    AFTER UPDATE ON workflows
    BEGIN
        UPDATE workflows SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_workflow_tasks_updated_at
    AFTER UPDATE ON workflow_tasks
    BEGIN
        UPDATE workflow_tasks SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_workflow_executions_updated_at
    AFTER UPDATE ON workflow_executions
    BEGIN
        UPDATE workflow_executions SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- 其他触发器
CREATE TRIGGER IF NOT EXISTS update_resource_groups_updated_at
    AFTER UPDATE ON resource_groups
    BEGIN
        UPDATE resource_groups SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_system_configs_updated_at
    AFTER UPDATE ON system_configs
    BEGIN
        UPDATE system_configs SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_alert_rules_updated_at
    AFTER UPDATE ON alert_rules
    BEGIN
        UPDATE alert_rules SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_alert_records_updated_at
    AFTER UPDATE ON alert_records
    BEGIN
        UPDATE alert_records SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_notification_templates_updated_at
    AFTER UPDATE ON notification_templates
    BEGIN
        UPDATE notification_templates SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_notification_records_updated_at
    AFTER UPDATE ON notification_records
    BEGIN
        UPDATE notification_records SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER IF NOT EXISTS update_notifications_updated_at
    AFTER UPDATE ON notifications
    BEGIN
        UPDATE notifications SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- Re-enable foreign key constraints after data insertion
PRAGMA foreign_keys = ON;
