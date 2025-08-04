#!/bin/bash

# Websoft9 提交消息格式检查脚本
# 基于 Conventional Commits 规范
# https://www.conventionalcommits.org/

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 帮助信息
show_help() {
    echo "Usage: $0 [OPTIONS] <commit-message>"
    echo ""
    echo "检查提交消息是否符合 Conventional Commits 规范"
    echo ""
    echo "Options:"
    echo "  -h, --help     显示帮助信息"
    echo "  -f, --file     从文件读取提交消息"
    echo "  -v, --verbose  显示详细信息"
    echo ""
    echo "Examples:"
    echo "  $0 'feat(auth): add JWT token refresh mechanism'"
    echo "  $0 -f .git/COMMIT_EDITMSG"
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
}

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
    local verbose="$2"
    
    if [ -z "$message" ]; then
        log_error "提交消息不能为空"
        return 1
    fi
    
    # 移除前后空白字符
    message=$(echo "$message" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    
    if [ "$verbose" = "true" ]; then
        log_info "检查提交消息: '$message'"
    fi
    
    # 定义允许的提交类型
    local valid_types="feat|fix|docs|style|refactor|test|chore|perf|ci"
    
    # Conventional Commits 正则表达式
    # 格式: <type>[optional scope]: <description>
    local regex="^(${valid_types})(\([a-zA-Z0-9_-]+\))?: .{1,50}$"
    
    if [[ ! "$message" =~ $regex ]]; then
        log_error "提交消息格式不正确"
        echo ""
        echo "正确格式: <type>[optional scope]: <description>"
        echo ""
        echo "示例:"
        echo "  feat(auth): add JWT token refresh mechanism"
        echo "  fix(api): handle null pointer in user service"
        echo "  docs(readme): update installation instructions"
        echo ""
        echo "要求:"
        echo "  - 类型必须是: ${valid_types//|/, }"
        echo "  - 描述长度: 1-50 字符"
        echo "  - 格式: 类型(可选范围): 描述"
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
    
    if [ "$verbose" = "true" ]; then
        log_info "类型: $type"
        if [ -n "$scope" ]; then
            log_info "范围: $scope"
        fi
        log_info "描述: $description"
    fi
    
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

# 主函数
main() {
    local commit_message=""
    local from_file=false
    local verbose=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -f|--file)
                from_file=true
                shift
                if [[ $# -gt 0 ]]; then
                    commit_message="$1"
                else
                    log_error "选项 -f 需要指定文件路径"
                    exit 1
                fi
                ;;
            -v|--verbose)
                verbose=true
                ;;
            -*)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
            *)
                if [ -z "$commit_message" ]; then
                    commit_message="$1"
                else
                    log_error "只能指定一个提交消息"
                    exit 1
                fi
                ;;
        esac
        shift
    done
    
    # 如果从文件读取
    if [ "$from_file" = true ]; then
        if [ ! -f "$commit_message" ]; then
            log_error "文件不存在: $commit_message"
            exit 1
        fi
        commit_message=$(head -n 1 "$commit_message")
    fi
    
    # 如果没有提供提交消息，尝试从标准输入读取
    if [ -z "$commit_message" ]; then
        if [ -t 0 ]; then
            log_error "请提供提交消息"
            show_help
            exit 1
        else
            commit_message=$(head -n 1)
        fi
    fi
    
    # 检查提交消息
    if check_commit_message "$commit_message" "$verbose"; then
        exit 0
    else
        exit 1
    fi
}

# 如果脚本被直接执行
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi