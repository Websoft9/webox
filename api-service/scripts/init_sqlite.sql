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
-- 3.1 Platform Home
-- ========================================

-- App shortcut navigation table
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
-- 3.2 Project Management
-- ========================================

-- Projects table
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

-- Project members table
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

-- Project environment variables table
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

-- Project files table
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

-- Project activity logs table
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

-- Workflows table
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

-- Workflow tasks table
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

-- Workflow execution history table
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

-- Resource groups table
CREATE TABLE IF NOT EXISTS resource_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(32) NOT NULL UNIQUE,
    description TEXT,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    sort_order INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1, -- 0-disabled, 1-enabled
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Database connections table
CREATE TABLE IF NOT EXISTS database_connections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(64) NOT NULL,
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

-- Servers table
CREATE TABLE IF NOT EXISTS servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(64) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    host VARCHAR(255) NOT NULL,
    internal_ip VARCHAR(45),
    ipv6_address VARCHAR(45),
    ssh_port INTEGER DEFAULT 22,
    os_distro VARCHAR(32),
    os_version VARCHAR(64),
    kernel_version VARCHAR(64),
    cpu_cores INTEGER DEFAULT 0,
    memory_total INTEGER DEFAULT 0,
    disk_total INTEGER DEFAULT 0,
    architecture VARCHAR(16),
    resource_group_id INTEGER,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- Server agents table
CREATE TABLE IF NOT EXISTS server_agents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    server_id INTEGER NOT NULL,
    agent_id VARCHAR(64) NOT NULL UNIQUE,
    deployment_type VARCHAR(20) DEFAULT 'docker',
    container_id VARCHAR(64),
    container_name VARCHAR(128),
    service_name VARCHAR(64),
    binary_path VARCHAR(255),
    config_path VARCHAR(255),
    agent_ip VARCHAR(45),
    agent_port INTEGER DEFAULT 8080,
    version VARCHAR(32),
    pull_mode INTEGER DEFAULT 1,
    pull_interval INTEGER DEFAULT 30,
    last_heartbeat_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(server_id)
);

-- Application instances table
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

-- SSL certificates table
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

-- Secret key management table
CREATE TABLE IF NOT EXISTS secret_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    key_type VARCHAR(20) NOT NULL, -- API_KEY, DATABASE, SSH, CERTIFICATE, CUSTOM
    description TEXT,
    custom_fields TEXT, -- JSON format
    expires_at DATETIME,
    resource_group_id INTEGER REFERENCES resource_groups(id) ON DELETE SET NULL,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- User secret access table (many-to-many relationship between users and secret_keys)
CREATE TABLE IF NOT EXISTS user_secret (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    secret_key_id INTEGER NOT NULL REFERENCES secret_keys(id) ON DELETE CASCADE,
    granted_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, secret_key_id)
);

-- Application gateways table
CREATE TABLE IF NOT EXISTS app_gateways (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(64) NOT NULL,
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

-- Application gateway publishing table
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

-- Application gateway access control rules table
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
-- 3.3 App Store
-- ========================================

-- App store categories table
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

-- App store templates table
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

-- App store wishlists table
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

-- App store reviews table
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

-- App favorites table
CREATE TABLE IF NOT EXISTS app_store_favorites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL REFERENCES app_store_templates(id),
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(template_id, user_id)
);

-- App stars table
CREATE TABLE IF NOT EXISTS app_store_stars (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL REFERENCES app_store_templates(id),
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(template_id, user_id)
);

-- App reports table
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

-- App download records table
CREATE TABLE IF NOT EXISTS app_store_downloads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL REFERENCES app_store_templates(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    ip_address VARCHAR(45),
    user_agent VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- App wishlist comments table
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

-- App wishlist votes table
CREATE TABLE IF NOT EXISTS app_store_wishlist_votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wishlist_id INTEGER NOT NULL REFERENCES app_store_wishlists(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(wishlist_id, user_id)
);

-- App wishlist likes table
CREATE TABLE IF NOT EXISTS app_store_wishlist_likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wishlist_id INTEGER NOT NULL REFERENCES app_store_wishlists(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(wishlist_id, user_id)
);

-- App wishlist reports table
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

-- App deployment records table
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
-- 3.4 Platform Management
-- ========================================

-- Users table
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

-- Roles table
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

-- Modules table
CREATE TABLE IF NOT EXISTS modules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(64) NOT NULL UNIQUE,
    description TEXT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Permissions table
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

-- User roles association table
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

-- Role permissions association table
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

-- API access tokens table
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

-- User two-factor authentication table
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

-- User profile configuration table
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

-- User login history table
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

-- Webhook configuration table
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

-- Webhook execution log table
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

-- Platform update records table
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

-- Container cluster nodes table
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

-- Container images table
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

-- System configuration table
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

-- Notification channels table
CREATE TABLE IF NOT EXISTS notification_channels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    channel_type VARCHAR(20) NOT NULL CHECK (channel_type IN ('EMAIL', 'WEBHOOK')),
    channel_config TEXT NOT NULL, -- JSON format
    owner_id INTEGER NOT NULL REFERENCES users(id),
    status INTEGER DEFAULT 1, --  0-disabled, 1-enabled
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Alert rules table
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

-- Alert records table
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

-- Notification templates table
CREATE TABLE IF NOT EXISTS notification_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    template_type VARCHAR(20) NOT NULL,
    subject VARCHAR(255),
    content TEXT NOT NULL,
    is_system INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Notification records table
CREATE TABLE IF NOT EXISTS notification_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER REFERENCES notification_templates(id),
    channel_type VARCHAR(20) NOT NULL,
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

-- Messages table
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

-- User notification records table
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

-- Audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    username VARCHAR(64),
    action VARCHAR(32) NOT NULL,
    module VARCHAR(32) NOT NULL,
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

-- Tag management
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(128) NOT NULL, -- Tag name
    color VARCHAR(16), -- Tag color
    description TEXT, -- Tag description
    created_by INTEGER DEFAULT 0, -- Creator user ID
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(name)
);

-- Tag references
CREATE TABLE IF NOT EXISTS taggings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tag_id INTEGER NOT NULL, -- Tag ID
    resource_code VARCHAR(64) NOT NULL, -- Resource code
    created_by INTEGER DEFAULT 0, -- Association creator
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Secret references table
CREATE TABLE IF NOT EXISTS secret_references (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    secret_id INTEGER NOT NULL,
    resource_code VARCHAR(64) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (secret_id) REFERENCES secret_keys(id) ON DELETE CASCADE
);

-- ========================================
-- Index creation
-- ========================================

-- Project related indexes
CREATE INDEX IF NOT EXISTS idx_projects_owner_status ON projects(owner_id, status);
CREATE INDEX IF NOT EXISTS idx_project_members_project ON project_members(project_id);
CREATE INDEX IF NOT EXISTS idx_project_members_user ON project_members(user_id);
CREATE INDEX IF NOT EXISTS idx_project_environments_project ON project_environments(project_id);
CREATE INDEX IF NOT EXISTS idx_project_files_project ON project_files(project_id);
CREATE INDEX IF NOT EXISTS idx_project_files_parent ON project_files(parent_id);
CREATE INDEX IF NOT EXISTS idx_project_activities_project ON project_activities(project_id);

-- Notification channels related indexes
CREATE INDEX IF NOT EXISTS idx_notification_channels_type ON notification_channels(channel_type);
CREATE INDEX IF NOT EXISTS idx_notification_channels_owner ON notification_channels(owner_id);
CREATE INDEX IF NOT EXISTS idx_notification_channels_code ON notification_channels(code);

-- User permissions related indexes
CREATE INDEX IF NOT EXISTS idx_user_roles_user ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission ON role_permissions(permission_code);
CREATE INDEX IF NOT EXISTS idx_module_code ON modules (code);

-- Alert notification related indexes
CREATE INDEX IF NOT EXISTS idx_alert_rules_owner ON alert_rules(owner_id);
CREATE INDEX IF NOT EXISTS idx_alert_records_rule ON alert_records(alert_rule_id);
CREATE INDEX IF NOT EXISTS idx_alert_records_status ON alert_records(status);
CREATE INDEX IF NOT EXISTS idx_notifications_target_status ON notifications(target_type, status);
CREATE INDEX IF NOT EXISTS idx_user_notifications_user_read ON user_notifications(user_id, is_read);

-- Time related indexes
CREATE INDEX IF NOT EXISTS idx_project_activities_created ON project_activities(created_at);
CREATE INDEX IF NOT EXISTS idx_workflow_executions_created ON workflow_executions(created_at);
CREATE INDEX IF NOT EXISTS idx_user_login_history_login_time ON user_login_history(login_time);

-- App store related indexes
CREATE INDEX IF NOT EXISTS idx_app_store_templates_category ON app_store_templates(category_id);
CREATE INDEX IF NOT EXISTS idx_app_store_templates_status ON app_store_templates(status);
CREATE INDEX IF NOT EXISTS idx_app_store_reviews_template ON app_store_reviews(template_id);
CREATE INDEX IF NOT EXISTS idx_app_store_reviews_user ON app_store_reviews(user_id);

-- Workflow related indexes
CREATE INDEX IF NOT EXISTS idx_workflows_project ON workflows(project_id);
CREATE INDEX IF NOT EXISTS idx_workflows_owner ON workflows(owner_id);
CREATE INDEX IF NOT EXISTS idx_workflow_tasks_workflow ON workflow_tasks(workflow_id);
CREATE INDEX IF NOT EXISTS idx_workflow_executions_task ON workflow_executions(task_id);

-- Resource management related indexes
CREATE INDEX IF NOT EXISTS idx_servers_name ON servers(name);
CREATE INDEX IF NOT EXISTS idx_servers_host ON servers(host);
CREATE INDEX IF NOT EXISTS idx_servers_owner ON servers(owner_id);
CREATE INDEX IF NOT EXISTS idx_servers_deleted_at ON servers(deleted_at);
CREATE INDEX IF NOT EXISTS idx_server_agents_server ON server_agents(server_id);
CREATE INDEX IF NOT EXISTS idx_server_agents_last_heartbeat ON server_agents(last_heartbeat_at);
CREATE INDEX IF NOT EXISTS idx_server_agents_deployment_type ON server_agents(deployment_type);
CREATE INDEX IF NOT EXISTS idx_app_instances_server ON app_instances(server_id);
CREATE INDEX IF NOT EXISTS idx_app_instances_template ON app_instances(template_id);

-- secret_references indexes
CREATE INDEX IF NOT EXISTS idx_secret_references_secret_id ON secret_references(secret_id);
CREATE INDEX IF NOT EXISTS idx_secret_references_resource_code ON secret_references(resource_code);
CREATE UNIQUE INDEX IF NOT EXISTS idx_secret_references_secret_resource ON secret_references(secret_id, resource_code);

-- Audit log indexes
CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at);

-- Tag management indexes
CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);
CREATE INDEX IF NOT EXISTS idx_tags_created_by ON tags(created_by);
CREATE INDEX IF NOT EXISTS idx_taggings_tag_id ON taggings(tag_id);
CREATE INDEX IF NOT EXISTS idx_taggings_resource_code ON taggings(resource_code);

-- ========================================
-- Create triggers for automatic updated_at field update
-- ========================================

-- Project related triggers
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

-- Create trigger for notification_channels updated_at
CREATE TRIGGER IF NOT EXISTS update_notification_channels_updated_at
    AFTER UPDATE ON notification_channels
    BEGIN
        UPDATE notification_channels SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- App store related triggers
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

-- User related triggers
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

CREATE TRIGGER IF NOT EXISTS modules_updated_at
AFTER UPDATE ON modules
FOR EACH ROW
BEGIN
    UPDATE modules SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Server related triggers
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

-- App instance related triggers
CREATE TRIGGER IF NOT EXISTS update_app_instances_updated_at
    AFTER UPDATE ON app_instances
    BEGIN
        UPDATE app_instances SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- Workflow related triggers
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

-- Other triggers
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
