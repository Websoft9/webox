#!/bin/bash

# Websoft9 提交消息格式检查脚本
# 基于 Conventional Commits 规范
# https://www.conventionalcommits.org/
# 
# 这是 Git commit-msg hook，会在提交时自动运行
# 基于 scripts/check-commit-message.sh 实现

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_error() {
    echo -e "${RED}❌ ERROR: $1${NC}" >&2
}

log_success() {
    echo -e "${GREEN}✅ SUCCESS: $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  WARNING: $1${NC}"
}

log_info() {
    echo -e "${BLUE}ℹ️  INFO: $1${NC}"
}

# 检查提交消息格式
check_commit_message() {
    local message="$1"
    
    if [ -z "$message" ]; then
        log_error "提交消息不能为空"
        return 1
    fi
    
    # 移除前后空白字符
    message=$(echo "$message" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    
    # 跳过合并提交和 revert 提交
    if [[ "$message" =~ ^Merge\ branch|^Revert\ |^Merge\ pull\ request ]]; then
        log_info "跳过合并/回滚提交检查"
        return 0
    fi
    
    # 定义允许的提交类型
    local valid_types="feat|fix|docs|style|refactor|test|chore|perf|ci"
    
    # Conventional Commits 正则表达式
    # 格式: <type>[optional scope]: <description>
    local regex="^(${valid_types})(\([a-zA-Z0-9_-]+\))?: .{1,100}$"
    
    if [[ ! "$message" =~ $regex ]]; then
        log_error "提交消息格式不正确"
        echo ""
        echo "正确格式: <type>[optional scope]: <description>"
        echo ""
        echo "示例:"
        echo "  feat(auth): add JWT token refresh mechanism"
        echo "  fix(api): handle null pointer in user service"
        echo "  docs(readme): update installation instructions"
        echo "  refactor(i18n): add internationalization support"
        echo ""
        echo "要求:"
        echo "  - 类型必须是: ${valid_types//|/, }"
        echo "  - 描述长度: 1-100 字符"
        echo "  - 格式: 类型(可选范围): 描述"
        echo ""
        echo "支持的提交类型："
        echo "  feat     - 新功能"
        echo "  fix      - Bug 修复"
        echo "  docs     - 文档更新"
        echo "  style    - 代码格式调整"
        echo "  refactor - 代码重构"
        echo "  test     - 测试相关"
        echo "  chore    - 构建过程或辅助工具的变动"
        echo "  perf     - 性能优化"
        echo "  ci       - CI/CD 相关"
        echo ""
        return 1
    fi
    
    # 提取类型和描述
    local type=$(echo "$message" | sed -E "s/^(${valid_types})(\([^)]+\))?: .*/\1/")
    local scope=""
    local description=""
    
    # 提取范围（如果存在）
    if echo "$message" | grep -q '('; then
        scope=$(echo "$message" | sed -E 's/^[^(]*\(([^)]+)\).*/\1/')
    fi
    
    description=$(echo "$message" | sed -E "s/^(${valid_types})(\([^)]+\))?: (.*)/\3/")
    
    # 检查描述是否以小写字母开头
    if [[ "$description" =~ ^[A-Z] ]]; then
        log_warning "建议描述以小写字母开头"
    fi
    
    # 检查描述是否以句号结尾
    if [[ "$description" =~ \.$$ ]]; then
        log_warning "描述不应以句号结尾"
    fi
    
    # 检查是否包含常见的不规范词汇
    local bad_words=("fixed" "added" "updated" "changed")
    for word in "${bad_words[@]}"; do
        if [[ "$description" =~ ^$word ]]; then
            log_warning "建议使用动词原形而不是过去式: '$word' -> '${word%ed}'"
        fi
    done
    
    log_success "提交消息格式正确"
    return 0
}

# 主函数 - Git hook 专用
main() {
    local commit_msg_file="$1"
    
    # 检查是否提供了提交消息文件
    if [ -z "$commit_msg_file" ]; then
        log_error "缺少提交消息文件参数"
        exit 1
    fi
    
    # 检查文件是否存在
    if [ ! -f "$commit_msg_file" ]; then
        log_error "提交消息文件不存在: $commit_msg_file"
        exit 1
    fi
    
    # 读取提交消息（只读第一行）
    local commit_message=$(head -n 1 "$commit_msg_file")
    
    # 检查提交消息
    if check_commit_message "$commit_message"; then
        exit 0
    else
        echo ""
        log_error "提交被拒绝。请修改提交消息后重试。"
        echo ""
        echo "可以使用以下命令修改提交消息："
        echo "  git commit --amend -m \"new message\""
        echo ""
        exit 1
    fi
}

# 执行主函数，传入 Git 提供的参数
main "$@"