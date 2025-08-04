#!/bin/bash

# Websoft9 项目测试执行脚本
# 执行所有测试并生成报告

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
TEST_TIMEOUT=${TEST_TIMEOUT:-10m}
COVERAGE_THRESHOLD=${COVERAGE_THRESHOLD:-80}
REPORTS_DIR=${REPORTS_DIR:-./reports}
VERBOSE=${VERBOSE:-false}

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

log_step() {
    echo -e "${BLUE}🔄 $1${NC}"
}

# 帮助信息
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "执行 Websoft9 项目的所有测试"
    echo ""
    echo "Options:"
    echo "  -h, --help              显示帮助信息"
    echo "  -v, --verbose           显示详细输出"
    echo "  -c, --coverage          生成覆盖率报告"
    echo "  -u, --unit              只运行单元测试"
    echo "  -i, --integration       只运行集成测试"
    echo "  -t, --timeout DURATION  测试超时时间 (默认: ${TEST_TIMEOUT})"
    echo "  --threshold PERCENT     覆盖率阈值 (默认: ${COVERAGE_THRESHOLD}%)"
    echo "  --reports-dir DIR       报告输出目录 (默认: ${REPORTS_DIR})"
    echo ""
    echo "Environment Variables:"
    echo "  TEST_TIMEOUT           测试超时时间"
    echo "  COVERAGE_THRESHOLD     覆盖率阈值"
    echo "  REPORTS_DIR           报告输出目录"
    echo "  VERBOSE               详细输出模式"
    echo ""
    echo "Examples:"
    echo "  $0                     # 运行所有测试"
    echo "  $0 -c                  # 运行测试并生成覆盖率报告"
    echo "  $0 -u -v               # 只运行单元测试，显示详细输出"
    echo "  $0 --threshold 90      # 设置覆盖率阈值为 90%"
}

# 检查依赖
check_dependencies() {
    log_step "检查依赖..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装或不在 PATH 中"
        exit 1
    fi
    
    local go_version=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Go 版本: $go_version"
    
    # 检查是否在 Go 模块目录中
    if [ ! -f "go.mod" ]; then
        log_error "当前目录不是 Go 模块根目录"
        exit 1
    fi
    
    log_success "依赖检查完成"
}

# 创建报告目录
setup_reports_dir() {
    if [ ! -d "$REPORTS_DIR" ]; then
        mkdir -p "$REPORTS_DIR"
        log_info "创建报告目录: $REPORTS_DIR"
    fi
}

# 运行单元测试
run_unit_tests() {
    log_step "运行单元测试..."
    
    local test_args="-v -timeout=$TEST_TIMEOUT"
    local coverage_args=""
    
    if [ "$GENERATE_COVERAGE" = "true" ]; then
        coverage_args="-coverprofile=$REPORTS_DIR/coverage.out -covermode=atomic"
    fi
    
    if [ "$VERBOSE" = "true" ]; then
        test_args="$test_args -v"
    fi
    
    # 查找所有包含测试的包
    local test_packages=$(go list ./... | grep -v vendor | grep -v /build/)
    
    if [ -z "$test_packages" ]; then
        log_warning "未找到测试文件"
        return 0
    fi
    
    log_info "测试包数量: $(echo "$test_packages" | wc -l)"
    
    # 运行测试
    if go test $test_args $coverage_args $test_packages 2>&1 | tee "$REPORTS_DIR/unit-tests.log"; then
        log_success "单元测试通过"
        return 0
    else
        log_error "单元测试失败"
        return 1
    fi
}

# 运行集成测试
run_integration_tests() {
    log_step "运行集成测试..."
    
    # 查找集成测试文件
    local integration_tests=$(find . -name "*_integration_test.go" -o -name "*_test.go" -path "*/integration/*")
    
    if [ -z "$integration_tests" ]; then
        log_info "未找到集成测试文件"
        return 0
    fi
    
    log_info "集成测试文件数量: $(echo "$integration_tests" | wc -l)"
    
    # 设置集成测试环境变量
    export INTEGRATION_TEST=true
    
    # 运行集成测试
    local test_args="-v -timeout=$TEST_TIMEOUT -tags=integration"
    
    if [ "$VERBOSE" = "true" ]; then
        test_args="$test_args -v"
    fi
    
    if go test $test_args ./... 2>&1 | tee "$REPORTS_DIR/integration-tests.log"; then
        log_success "集成测试通过"
        return 0
    else
        log_error "集成测试失败"
        return 1
    fi
}

# 生成覆盖率报告
generate_coverage_report() {
    if [ ! -f "$REPORTS_DIR/coverage.out" ]; then
        log_warning "未找到覆盖率数据文件"
        return 0
    fi
    
    log_step "生成覆盖率报告..."
    
    # 生成 HTML 报告
    go tool cover -html="$REPORTS_DIR/coverage.out" -o "$REPORTS_DIR/coverage.html"
    log_info "HTML 覆盖率报告: $REPORTS_DIR/coverage.html"
    
    # 生成文本报告
    go tool cover -func="$REPORTS_DIR/coverage.out" > "$REPORTS_DIR/coverage.txt"
    log_info "文本覆盖率报告: $REPORTS_DIR/coverage.txt"
    
    # 计算总覆盖率
    local total_coverage=$(go tool cover -func="$REPORTS_DIR/coverage.out" | grep "total:" | awk '{print $3}' | sed 's/%//')
    
    if [ -n "$total_coverage" ]; then
        log_info "总覆盖率: ${total_coverage}%"
        
        # 检查覆盖率阈值
        if (( $(echo "$total_coverage >= $COVERAGE_THRESHOLD" | bc -l) )); then
            log_success "覆盖率达到阈值 (${total_coverage}% >= ${COVERAGE_THRESHOLD}%)"
        else
            log_warning "覆盖率未达到阈值 (${total_coverage}% < ${COVERAGE_THRESHOLD}%)"
            return 1
        fi
    else
        log_warning "无法计算总覆盖率"
    fi
    
    return 0
}

# 生成测试报告摘要
generate_test_summary() {
    log_step "生成测试报告摘要..."
    
    local summary_file="$REPORTS_DIR/test-summary.md"
    
    cat > "$summary_file" << EOF
# 测试报告摘要

**生成时间:** $(date)
**Go 版本:** $(go version | awk '{print $3}')
**测试超时:** $TEST_TIMEOUT
**覆盖率阈值:** ${COVERAGE_THRESHOLD}%

## 测试结果

EOF
    
    # 单元测试结果
    if [ -f "$REPORTS_DIR/unit-tests.log" ]; then
        local unit_test_result=$(tail -n 5 "$REPORTS_DIR/unit-tests.log" | grep -E "(PASS|FAIL)" | tail -n 1)
        echo "### 单元测试" >> "$summary_file"
        echo "- 状态: $unit_test_result" >> "$summary_file"
        echo "- 日志: [unit-tests.log](./unit-tests.log)" >> "$summary_file"
        echo "" >> "$summary_file"
    fi
    
    # 集成测试结果
    if [ -f "$REPORTS_DIR/integration-tests.log" ]; then
        local integration_test_result=$(tail -n 5 "$REPORTS_DIR/integration-tests.log" | grep -E "(PASS|FAIL)" | tail -n 1)
        echo "### 集成测试" >> "$summary_file"
        echo "- 状态: $integration_test_result" >> "$summary_file"
        echo "- 日志: [integration-tests.log](./integration-tests.log)" >> "$summary_file"
        echo "" >> "$summary_file"
    fi
    
    # 覆盖率报告
    if [ -f "$REPORTS_DIR/coverage.txt" ]; then
        local total_coverage=$(grep "total:" "$REPORTS_DIR/coverage.txt" | awk '{print $3}')
        echo "### 代码覆盖率" >> "$summary_file"
        echo "- 总覆盖率: $total_coverage" >> "$summary_file"
        echo "- HTML 报告: [coverage.html](./coverage.html)" >> "$summary_file"
        echo "- 详细报告: [coverage.txt](./coverage.txt)" >> "$summary_file"
        echo "" >> "$summary_file"
    fi
    
    log_info "测试报告摘要: $summary_file"
}

# 清理函数
cleanup() {
    log_info "清理临时文件..."
    # 清理可能的临时文件
    find . -name "*.test" -delete 2>/dev/null || true
}

# 主函数
main() {
    local run_unit_tests=true
    local run_integration_tests=true
    local generate_coverage=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--verbose)
                VERBOSE=true
                ;;
            -c|--coverage)
                generate_coverage=true
                ;;
            -u|--unit)
                run_integration_tests=false
                ;;
            -i|--integration)
                run_unit_tests=false
                ;;
            -t|--timeout)
                shift
                if [[ $# -gt 0 ]]; then
                    TEST_TIMEOUT="$1"
                else
                    log_error "选项 --timeout 需要指定时间"
                    exit 1
                fi
                ;;
            --threshold)
                shift
                if [[ $# -gt 0 ]]; then
                    COVERAGE_THRESHOLD="$1"
                else
                    log_error "选项 --threshold 需要指定百分比"
                    exit 1
                fi
                ;;
            --reports-dir)
                shift
                if [[ $# -gt 0 ]]; then
                    REPORTS_DIR="$1"
                else
                    log_error "选项 --reports-dir 需要指定目录"
                    exit 1
                fi
                ;;
            -*)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
        shift
    done
    
    # 设置清理陷阱
    trap cleanup EXIT
    
    # 检查依赖
    check_dependencies
    
    # 设置报告目录
    setup_reports_dir
    
    # 设置覆盖率生成
    GENERATE_COVERAGE=$generate_coverage
    
    local exit_code=0
    
    # 运行测试
    if [ "$run_unit_tests" = true ]; then
        if ! run_unit_tests; then
            exit_code=1
        fi
    fi
    
    if [ "$run_integration_tests" = true ]; then
        if ! run_integration_tests; then
            exit_code=1
        fi
    fi
    
    # 生成覆盖率报告
    if [ "$generate_coverage" = true ]; then
        if ! generate_coverage_report; then
            # 覆盖率不达标不算测试失败，只是警告
            log_warning "覆盖率检查未通过，但不影响测试结果"
        fi
    fi
    
    # 生成测试报告摘要
    generate_test_summary
    
    if [ $exit_code -eq 0 ]; then
        log_success "所有测试完成"
    else
        log_error "部分测试失败"
    fi
    
    exit $exit_code
}

# 如果脚本被直接执行
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi