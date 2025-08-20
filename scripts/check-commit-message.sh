#!/bin/bash

# 提交消息格式检查脚本
# 用于验证提交消息是否符合 Conventional Commits 规范

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 使用说明
usage() {
    echo "用法: $0 <commit-message>"
    echo ""
    echo "检查提交消息是否符合 Conventional Commits 规范"
    echo ""
    echo "格式: <type>[optional scope]: <description>"
    echo ""
    echo "支持的类型:"
    echo "  feat     - 新功能"
    echo "  fix      - 修复 bug"
    echo "  docs     - 文档更新"
    echo "  style    - 代码格式调整"
    echo "  refactor - 代码重构"
    echo "  test     - 测试相关"
    echo "  chore    - 构建过程或辅助工具的变动"
    echo "  perf     - 性能优化"
    echo "  ci       - CI/CD 相关"
    echo ""
    echo "示例:"
    echo "  feat(auth): add JWT token refresh mechanism"
    echo "  fix(api): handle null pointer in user service"
    echo "  docs(readme): update installation instructions"
}

# 检查参数
if [ $# -eq 0 ]; then
    echo -e "${RED}错误: 缺少提交消息参数${NC}"
    usage
    exit 1
fi

if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    usage
    exit 0
fi

commit_message="$1"

# 定义正则表达式
# 支持的类型
types="feat|fix|docs|style|refactor|test|chore|perf|ci"

# 完整的正则表达式
# ^(type)(\(scope\))?: description
commit_regex="^(${types})(\([a-zA-Z0-9_-]+\))?: .{1,72}$"

# 检查提交消息格式
if echo "$commit_message" | grep -qE "$commit_regex"; then
    echo -e "${GREEN}✅ 提交消息格式正确${NC}"
    echo "消息: $commit_message"

    # 提取类型和作用域
    type=$(echo "$commit_message" | sed -E "s/^(feat|fix|docs|style|refactor|test|chore|perf|ci)(\([^)]+\))?: .*/\1/")
    scope=""
    description=""

    if echo "$commit_message" | grep -qE "^(feat|fix|docs|style|refactor|test|chore|perf|ci)\([^)]+\): "; then
        # 有作用域的情况
        scope=$(echo "$commit_message" | sed -E "s/^(feat|fix|docs|style|refactor|test|chore|perf|ci)\(([^)]+)\): .*/\2/")
        description=$(echo "$commit_message" | sed -E "s/^(feat|fix|docs|style|refactor|test|chore|perf|ci)\([^)]+\): (.*)/\2/")
    else
        # 没有作用域的情况
        description=$(echo "$commit_message" | sed -E "s/^(feat|fix|docs|style|refactor|test|chore|perf|ci): (.*)/\2/")
    fi

    echo "类型: $type"
    if [ -n "$scope" ]; then
        echo "作用域: $scope"
    fi
    echo "描述: $description"

    exit 0
else
    echo -e "${RED}❌ 提交消息格式不正确${NC}"
    echo "消息: $commit_message"
    echo ""
    echo -e "${YELLOW}正确格式:${NC} <type>[optional scope]: <description>"
    echo ""
    echo -e "${YELLOW}支持的类型:${NC}"
    echo "  feat, fix, docs, style, refactor, test, chore, perf, ci"
    echo ""
    echo -e "${YELLOW}示例:${NC}"
    echo "  feat(auth): add JWT token refresh mechanism"
    echo "  fix(api): handle null pointer in user service"
    echo "  docs: update README installation guide"
    echo ""
    echo -e "${YELLOW}常见问题:${NC}"

    # 分析常见错误
    if echo "$commit_message" | grep -qE "^[A-Z]"; then
        echo "  - 类型应该使用小写字母"
    fi

    if echo "$commit_message" | grep -qE "^(${types})(\([^)]+\))?[^:]"; then
        echo "  - 类型后面应该有冒号和空格"
    fi

    if echo "$commit_message" | grep -qE "^(${types})(\([^)]+\))?: [A-Z]"; then
        echo "  - 描述应该以小写字母开头"
    fi

    if [ ${#commit_message} -gt 100 ]; then
        echo "  - 提交消息过长，建议不超过 100 个字符"
    fi

    if [ ${#commit_message} -lt 5 ]; then
        echo "  - 提交消息过短，应该提供有意义的描述"
    fi

    # 检查是否使用了不支持的类型
    possible_type=$(echo "$commit_message" | sed -E 's/^([a-zA-Z]+).*/\1/')
    if ! echo "$possible_type" | grep -qE "^(${types})$"; then
        echo "  - 不支持的类型 '$possible_type'，请使用: feat, fix, docs, style, refactor, test, chore, perf, ci"
    fi

    exit 1
fi