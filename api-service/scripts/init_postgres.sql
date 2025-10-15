-- Websoft9 PostgreSQL Database Initialization Script V1.1
-- Generated from Database Design Document V1.1
-- Compatible with PostgreSQL 12+

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

-- ========================================
-- 3.1 Platform Home (Platform Home)
-- ========================================

-- App shortcuts table
CREATE TABLE app_shortcuts (
    id BIGSERIAL PRIMARY KEY,
    app_instance_id BIGINT NOT NULL,
    name VARCHAR(64),
    description VARCHAR(255),
    icon VARCHAR(255),
    sort_order INTEGER NOT NULL DEFAULT 0,
    access_count INTEGER NOT NULL DEFAULT 0,
    last_accessed TIMESTAMP,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_app_shortcuts_app_instance ON app_shortcuts(app_instance_id);
CREATE INDEX idx_app_shortcuts_user ON app_shortcuts(user_id);
CREATE INDEX idx_app_shortcuts_sort_order ON app_shortcuts(sort_order);

COMMENT ON TABLE app_shortcuts IS 'Application shortcut navigation table';

-- ========================================
-- 3.2 Project Management (Project Management)
-- ========================================

-- Projects table
CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    identifier VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    tags JSONB,
    icon VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'NORMAL' CHECK (status IN ('NORMAL', 'ARCHIVED', 'DELETED')),
    owner_id BIGINT NOT NULL,
    default_resource_group VARCHAR(50) NOT NULL DEFAULT 'default',
    default_timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Shanghai',
    log_retention_days INTEGER NOT NULL DEFAULT 30,
    backup_strategy VARCHAR(20) NOT NULL DEFAULT 'daily' CHECK (backup_strategy IN ('daily', 'weekly', 'monthly', 'disabled')),
    access_control VARCHAR(50) NOT NULL DEFAULT 'members',
    api_access BOOLEAN NOT NULL DEFAULT TRUE,
    audit_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    last_activity_at TIMESTAMP,
    archived_at TIMESTAMP,
    archived_by BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_projects_owner ON projects(owner_id);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_archived_by ON projects(archived_by);

COMMENT ON TABLE projects IS 'Projects table';

-- Continue with the rest of the tables from the MySQL script,
-- converting MySQL-specific syntax to PostgreSQL syntax

-- Users table (simplified for demo)
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(64),
    avatar VARCHAR(255),
    phone VARCHAR(20),
    gender SMALLINT NOT NULL DEFAULT 0 CHECK (gender IN (0, 1, 2)),
    signature VARCHAR(255),
    status SMALLINT NOT NULL DEFAULT 1 CHECK (status IN (0, 1)),
    last_login_at TIMESTAMP,
    last_login_ip INET,
    timezone VARCHAR(64) NOT NULL DEFAULT 'UTC',
    language VARCHAR(10) NOT NULL DEFAULT 'zh-CN',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_email ON users(email);

COMMENT ON TABLE users IS 'Users table';

-- Roles table
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE,
    code VARCHAR(32) NOT NULL UNIQUE,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 1 CHECK (status IN (-1, 0, 1)),
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_roles_status ON roles(status);
CREATE INDEX idx_roles_created_by ON roles(created_by);
CREATE INDEX idx_roles_updated_by ON roles(updated_by);

COMMENT ON TABLE roles IS 'Roles table';

-- Permissions table
CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT,
    scope VARCHAR(64) NOT NULL,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(64) NOT NULL UNIQUE,
    module VARCHAR(32) NOT NULL,
    action VARCHAR(32) NOT NULL,
    resource VARCHAR(64),
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    is_menu BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 1 CHECK (status IN (-1, 0, 1)),
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_permissions_parent ON permissions(parent_id);
CREATE INDEX idx_permissions_scope ON permissions(scope);
CREATE INDEX idx_permissions_module ON permissions(module);
CREATE INDEX idx_permissions_action ON permissions(action);
CREATE INDEX idx_permissions_status ON permissions(status);
CREATE INDEX idx_permissions_created_by ON permissions(created_by);
CREATE INDEX idx_permissions_updated_by ON permissions(updated_by);

COMMENT ON TABLE permissions IS 'Permissions table';

-- User roles association table
CREATE TABLE user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    granted_by BIGINT,
    granted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    status SMALLINT NOT NULL DEFAULT 1 CHECK (status IN (-1, 0, 1)),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id)
);

CREATE INDEX idx_user_roles_role ON user_roles(role_id);
CREATE INDEX idx_user_roles_granted_by ON user_roles(granted_by);
CREATE INDEX idx_user_roles_status ON user_roles(status);

COMMENT ON TABLE user_roles IS 'User roles association table';

-- Role permissions association table
CREATE TABLE role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    granted_by BIGINT,
    granted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status SMALLINT NOT NULL DEFAULT 1 CHECK (status IN (-1, 0, 1)),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, permission_id)
);

CREATE INDEX idx_role_permissions_permission ON role_permissions(permission_id);
CREATE INDEX idx_role_permissions_granted_by ON role_permissions(granted_by);
CREATE INDEX idx_role_permissions_status ON role_permissions(status);

COMMENT ON TABLE role_permissions IS 'Role permissions association table';

-- Secret references table
CREATE TABLE IF NOT EXISTS secret_references (
    id BIGSERIAL PRIMARY KEY,
    secret_id BIGINT NOT NULL,
    resource_code VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_secret_references_secret FOREIGN KEY (secret_id) REFERENCES secret_keys(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_secret_references_secret_id ON secret_references(secret_id);
CREATE INDEX IF NOT EXISTS idx_secret_references_resource_code ON secret_references(resource_code);
CREATE UNIQUE INDEX IF NOT EXISTS idx_secret_references_secret_resource ON secret_references(secret_id, resource_code);


-- Add foreign key constraints
ALTER TABLE app_shortcuts ADD CONSTRAINT fk_app_shortcuts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE projects ADD CONSTRAINT fk_projects_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE RESTRICT;
ALTER TABLE projects ADD CONSTRAINT fk_projects_archived_by FOREIGN KEY (archived_by) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE permissions ADD CONSTRAINT fk_permissions_parent FOREIGN KEY (parent_id) REFERENCES permissions(id) ON DELETE SET NULL;
ALTER TABLE permissions ADD CONSTRAINT fk_permissions_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE permissions ADD CONSTRAINT fk_permissions_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE roles ADD CONSTRAINT fk_roles_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE roles ADD CONSTRAINT fk_roles_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;
ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_granted_by FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;
ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE;
ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_granted_by FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL;

-- Insert default roles
INSERT INTO roles (name, code, description, is_system, status) VALUES
('Super Administrator', 'super_admin', 'System super administrator with all permissions', TRUE, 1),
('System Administrator', 'admin', 'System administrator responsible for platform management', TRUE, 1),
('Project Administrator', 'project_admin', 'Project administrator responsible for project management', TRUE, 1),
('Developer', 'developer', 'Developer role responsible for application development', TRUE, 1),
('Operator', 'operator', 'Operator role responsible for system operations', TRUE, 1),
('Regular User', 'user', 'Regular user role with basic permissions', TRUE, 1);

-- Insert default permissions (sample)
INSERT INTO permissions (scope, name, code, module, action, description, is_system, status) VALUES
-- Platform home permissions
('platform', 'Project Overview Dashboard', 'b2fabec8-f7fe-480d-883c-5e0a355b8c69', 'project_overview', '*', 'All permissions for project overview dashboard', TRUE, 1),
('platform', 'Project Overview Dashboard Query', 'e0aec98a-101f-494a-b5ea-ada5cb5071d1', 'project_overview', 'query', 'Query permission for project overview dashboard', TRUE, 1),
('platform', 'Monitor Overview Dashboard', 'b61077db-d098-42d0-ade9-472038167183', 'monitor_overview', '*', 'All permissions for monitor overview dashboard', TRUE, 1),
('platform', 'Monitor Overview Dashboard Query', '91ab578c-1c7e-4f56-9708-4d00f4ad3c2d', 'monitor_overview', 'query', 'Query permission for monitor overview dashboard', TRUE, 1),

-- Project management permissions
('project', 'Project Dashboard', 'dee6721a-870a-413c-8b6c-2fb82b153479', 'dashboard', '*', 'All permissions for project dashboard', TRUE, 1),
('project', 'Project Dashboard Query', 'b0335079-0f4e-49a2-91ed-3946c0f4d13a', 'dashboard', 'query', 'Query permission for project dashboard', TRUE, 1),

-- User management permissions
('platform', 'User Management', '1adb3cfc-ddd8-42ab-9c0f-3b2881184852', 'user', '*', 'All permissions for user management', TRUE, 1),
('platform', 'User Management Query', 'a4e88f35-af15-4fe0-b30b-ec2f516e9fc3', 'user', 'query', 'Query permission for user management', TRUE, 1),
('platform', 'User Management Create', '01aaa408-953e-4c5a-8cb0-b507322a4205', 'user', 'create', 'Create permission for user management', TRUE, 1),
('platform', 'User Management Update', 'a176c0fc-24ba-42f6-bb8f-1f80a8a2fe29', 'user', 'update', 'Update permission for user management', TRUE, 1),
('platform', 'User Management Delete', '9644b60c-b8ea-4351-b672-e7b656bb33c2', 'user', 'delete', 'Delete permission for user management', TRUE, 1);

-- Create functions for automatic timestamp updates
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at columns
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_permissions_updated_at BEFORE UPDATE ON permissions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_app_shortcuts_updated_at BEFORE UPDATE ON app_shortcuts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_roles_updated_at BEFORE UPDATE ON user_roles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_role_permissions_updated_at BEFORE UPDATE ON role_permissions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
