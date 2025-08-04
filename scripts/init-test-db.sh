#!/bin/bash

# 测试数据库初始化脚本
# 支持 SQLite 数据库初始化和种子数据

set -e

# 默认参数
DB_TYPE="sqlite"
DB_PATH="./test.db"
SEED_DATA=false

# 解析命令行参数
while [[ $# -gt 0 ]]; do
  case $1 in
    --type)
      DB_TYPE="$2"
      shift 2
      ;;
    --database)
      DB_PATH="$2"
      shift 2
      ;;
    --seed)
      SEED_DATA=true
      shift
      ;;
    *)
      echo "未知参数: $1"
      exit 1
      ;;
  esac
done

echo "初始化测试数据库..."
echo "数据库类型: $DB_TYPE"
echo "数据库路径: $DB_PATH"
echo "种子数据: $SEED_DATA"

if [ "$DB_TYPE" = "sqlite" ]; then
    # 删除现有数据库文件
    if [ -f "$DB_PATH" ]; then
        echo "删除现有数据库文件: $DB_PATH"
        rm -f "$DB_PATH"
    fi
    
    # 创建数据库目录
    DB_DIR=$(dirname "$DB_PATH")
    if [ ! -d "$DB_DIR" ]; then
        mkdir -p "$DB_DIR"
    fi
    
    echo "创建 SQLite 数据库: $DB_PATH"
    
    # 创建基本表结构（示例）
    sqlite3 "$DB_PATH" << 'EOF'
-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 应用表
CREATE TABLE IF NOT EXISTS applications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image VARCHAR(255),
    status VARCHAR(20) DEFAULT 'stopped',
    user_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 配置表
CREATE TABLE IF NOT EXISTS configurations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
EOF
    
    if [ "$SEED_DATA" = true ]; then
        echo "插入种子数据..."
        sqlite3 "$DB_PATH" << 'EOF'
-- 插入测试用户
INSERT INTO users (username, email, password_hash, role) VALUES 
('admin', 'admin@websoft9.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjPeGvGzluGoSvXz2XZnQ9vO0rO9dO', 'admin'),
('testuser', 'test@websoft9.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjPeGvGzluGoSvXz2XZnQ9vO0rO9dO', 'user');

-- 插入测试应用
INSERT INTO applications (name, description, image, status, user_id) VALUES 
('nginx', 'Nginx Web Server', 'nginx:latest', 'running', 1),
('mysql', 'MySQL Database', 'mysql:8.0', 'stopped', 1),
('redis', 'Redis Cache', 'redis:7.0', 'running', 2);

-- 插入系统配置
INSERT INTO configurations (key, value, description) VALUES 
('system.name', 'Websoft9 Test', '系统名称'),
('system.version', '1.0.0', '系统版本'),
('docker.registry', 'docker.io', 'Docker 镜像仓库');
EOF
    fi
    
    echo "✅ SQLite 数据库初始化完成"
    
    # 显示数据库信息
    echo "数据库表:"
    sqlite3 "$DB_PATH" ".tables"
    
    if [ "$SEED_DATA" = true ]; then
        echo "用户数量: $(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM users;")"
        echo "应用数量: $(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM applications;")"
    fi
    
else
    echo "❌ 不支持的数据库类型: $DB_TYPE"
    echo "目前只支持 sqlite"
    exit 1
fi

echo "数据库初始化完成！"