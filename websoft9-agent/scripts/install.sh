#!/bin/bash

# Websoft9 Agent 安装脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否为 root 用户
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "此脚本需要 root 权限运行"
        exit 1
    fi
}

# 检查系统要求
check_requirements() {
    log_info "检查系统要求..."
    
    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    # 检查 Docker 服务状态
    if ! systemctl is-active --quiet docker; then
        log_warn "Docker 服务未运行，正在启动..."
        systemctl start docker
        systemctl enable docker
    fi
    
    log_info "系统要求检查完成"
}

# 创建目录结构
create_directories() {
    log_info "创建目录结构..."
    
    mkdir -p /etc/websoft9
    mkdir -p /var/lib/websoft9/agent
    mkdir -p /var/log/websoft9
    
    log_info "目录创建完成"
}

# 安装 Agent
install_agent() {
    log_info "安装 Websoft9 Agent..."
    
    # 复制二进制文件
    if [[ -f "./build/websoft9-agent" ]]; then
        cp ./build/websoft9-agent /usr/local/bin/
        chmod +x /usr/local/bin/websoft9-agent
    else
        log_error "找不到构建的二进制文件，请先运行 make build"
        exit 1
    fi
    
    # 复制配置文件
    if [[ -f "./configs/agent.yaml" ]]; then
        cp ./configs/agent.yaml /etc/websoft9/
    else
        log_error "找不到配置文件"
        exit 1
    fi
    
    log_info "Agent 安装完成"
}

# 创建 systemd 服务
create_service() {
    log_info "创建 systemd 服务..."
    
    cat > /etc/systemd/system/websoft9-agent.service << EOF
[Unit]
Description=Websoft9 Agent
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/websoft9-agent -config /etc/websoft9/agent.yaml
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable websoft9-agent
    
    log_info "systemd 服务创建完成"
}

# 启动服务
start_service() {
    log_info "启动 Websoft9 Agent 服务..."
    
    systemctl start websoft9-agent
    
    # 检查服务状态
    if systemctl is-active --quiet websoft9-agent; then
        log_info "Websoft9 Agent 服务启动成功"
    else
        log_error "Websoft9 Agent 服务启动失败"
        systemctl status websoft9-agent
        exit 1
    fi
}

# 显示状态
show_status() {
    log_info "安装完成！"
    echo
    echo "服务状态:"
    systemctl status websoft9-agent --no-pager -l
    echo
    echo "查看日志:"
    echo "  journalctl -u websoft9-agent -f"
    echo
    echo "管理服务:"
    echo "  systemctl start websoft9-agent"
    echo "  systemctl stop websoft9-agent"
    echo "  systemctl restart websoft9-agent"
    echo "  systemctl status websoft9-agent"
}

# 主函数
main() {
    log_info "开始安装 Websoft9 Agent..."
    
    check_root
    check_requirements
    create_directories
    install_agent
    create_service
    start_service
    show_status
}

# 运行主函数
main "$@"