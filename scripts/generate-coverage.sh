#!/bin/bash

# Websoft9 项目覆盖率报告生成脚本
# 生成详细的测试覆盖率报告

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
COVERAGE_FILE=${COVERAGE_FILE:-coverage.out}
REPORTS_DIR=${REPORTS_DIR:-./reports}
COVERAGE_THRESHOLD=${COVERAGE_THRESHOLD:-80}
EXCLUDE_PATTERNS=${EXCLUDE_PATTERNS:-"vendor,build,mocks,*.pb.go,*_mock.go"}

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
    echo "生成 Websoft9 项目的测试覆盖率报告"
    echo ""
    echo "Options:"
    echo "  -h, --help              显示帮助信息"
    echo "  -f, --file FILE         覆盖率数据文件 (默认: ${COVERAGE_FILE})"
    echo "  -o, --output DIR        输出目录 (默认: ${REPORTS_DIR})"
    echo "  -t, --threshold PERCENT 覆盖率阈值 (默认: ${COVERAGE_THRESHOLD}%)"
    echo "  -e, --exclude PATTERNS  排除模式，逗号分隔 (默认: ${EXCLUDE_PATTERNS})"
    echo "  --html                  生成 HTML 报告"
    echo "  --json                  生成 JSON 报告"
    echo "  --xml                   生成 XML 报告 (Cobertura 格式)"
    echo "  --lcov                  生成 LCOV 报告"
    echo "  --all                   生成所有格式的报告"
    echo ""
    echo "Environment Variables:"
    echo "  COVERAGE_FILE          覆盖率数据文件路径"
    echo "  REPORTS_DIR           报告输出目录"
    echo "  COVERAGE_THRESHOLD    覆盖率阈值"
    echo "  EXCLUDE_PATTERNS      排除模式"
    echo ""
    echo "Examples:"
    echo "  $0                     # 生成基本覆盖率报告"
    echo "  $0 --html              # 生成 HTML 报告"
    echo "  $0 --all               # 生成所有格式的报告"
    echo "  $0 -t 90 --html        # 设置阈值为 90% 并生成 HTML 报告"
}

# 检查依赖
check_dependencies() {
    log_step "检查依赖..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装或不在 PATH 中"
        exit 1
    fi
    
    # 检查是否安装了额外的工具
    local missing_tools=()
    
    if [ "$GENERATE_XML" = "true" ] && ! command -v gocover-cobertura &> /dev/null; then
        missing_tools+=("gocover-cobertura")
    fi
    
    if [ "$GENERATE_LCOV" = "true" ] && ! command -v gcov2lcov &> /dev/null; then
        missing_tools+=("gcov2lcov")
    fi
    
    if [ ${#missing_tools[@]} -gt 0 ]; then
        log_warning "以下工具未安装，将跳过相应格式的报告生成:"
        for tool in "${missing_tools[@]}"; do
            echo "  - $tool"
        done
        echo ""
        echo "安装命令:"
        for tool in "${missing_tools[@]}"; do
            case $tool in
                gocover-cobertura)
                    echo "  go install github.com/boumenot/gocover-cobertura@latest"
                    ;;
                gcov2lcov)
                    echo "  go install github.com/jandelgado/gcov2lcov@latest"
                    ;;
            esac
        done
    fi
    
    log_success "依赖检查完成"
}

# 创建输出目录
setup_output_dir() {
    if [ ! -d "$REPORTS_DIR" ]; then
        mkdir -p "$REPORTS_DIR"
        log_info "创建输出目录: $REPORTS_DIR"
    fi
}

# 检查覆盖率文件
check_coverage_file() {
    if [ ! -f "$COVERAGE_FILE" ]; then
        log_error "覆盖率文件不存在: $COVERAGE_FILE"
        log_info "请先运行测试生成覆盖率数据:"
        log_info "  go test -coverprofile=$COVERAGE_FILE ./..."
        exit 1
    fi
    
    log_info "使用覆盖率文件: $COVERAGE_FILE"
}

# 过滤覆盖率数据
filter_coverage_data() {
    log_step "过滤覆盖率数据..."
    
    local filtered_file="$REPORTS_DIR/coverage-filtered.out"
    
    # 复制原始文件
    cp "$COVERAGE_FILE" "$filtered_file"
    
    # 应用排除模式
    IFS=',' read -ra PATTERNS <<< "$EXCLUDE_PATTERNS"
    for pattern in "${PATTERNS[@]}"; do
        pattern=$(echo "$pattern" | xargs) # 去除空白字符
        if [ -n "$pattern" ]; then
            log_info "排除模式: $pattern"
            grep -v "$pattern" "$filtered_file" > "$filtered_file.tmp" && mv "$filtered_file.tmp" "$filtered_file"
        fi
    done
    
    COVERAGE_FILE="$filtered_file"
    log_success "覆盖率数据过滤完成"
}

# 生成基本覆盖率报告
generate_basic_report() {
    log_step "生成基本覆盖率报告..."
    
    local output_file="$REPORTS_DIR/coverage.txt"
    
    go tool cover -func="$COVERAGE_FILE" > "$output_file"
    
    log_info "基本覆盖率报告: $output_file"
    
    # 显示覆盖率摘要
    echo ""
    echo "=== 覆盖率摘要 ==="
    tail -n 10 "$output_file"
    echo ""
}

# 生成 HTML 报告
generate_html_report() {
    if [ "$GENERATE_HTML" != "true" ]; then
        return 0
    fi
    
    log_step "生成 HTML 覆盖率报告..."
    
    local output_file="$REPORTS_DIR/coverage.html"
    
    go tool cover -html="$COVERAGE_FILE" -o "$output_file"
    
    log_success "HTML 覆盖率报告: $output_file"
}

# 生成 JSON 报告
generate_json_report() {
    if [ "$GENERATE_JSON" != "true" ]; then
        return 0
    fi
    
    log_step "生成 JSON 覆盖率报告..."
    
    local output_file="$REPORTS_DIR/coverage.json"
    
    # 使用 go tool cover 生成 JSON 格式的报告
    # 注意：这需要 Go 1.20+ 版本
    if go tool cover -func="$COVERAGE_FILE" -o json > "$output_file" 2>/dev/null; then
        log_success "JSON 覆盖率报告: $output_file"
    else
        log_warning "当前 Go 版本不支持 JSON 格式输出，跳过 JSON 报告生成"
    fi
}

# 生成 XML 报告 (Cobertura 格式)
generate_xml_report() {
    if [ "$GENERATE_XML" != "true" ]; then
        return 0
    fi
    
    log_step "生成 XML 覆盖率报告 (Cobertura 格式)..."
    
    if ! command -v gocover-cobertura &> /dev/null; then
        log_warning "gocover-cobertura 未安装，跳过 XML 报告生成"
        return 0
    fi
    
    local output_file="$REPORTS_DIR/coverage.xml"
    
    gocover-cobertura < "$COVERAGE_FILE" > "$output_file"
    
    log_success "XML 覆盖率报告: $output_file"
}

# 生成 LCOV 报告
generate_lcov_report() {
    if [ "$GENERATE_LCOV" != "true" ]; then
        return 0
    fi
    
    log_step "生成 LCOV 覆盖率报告..."
    
    if ! command -v gcov2lcov &> /dev/null; then
        log_warning "gcov2lcov 未安装，跳过 LCOV 报告生成"
        return 0
    fi
    
    local output_file="$REPORTS_DIR/coverage.lcov"
    
    gcov2lcov -infile "$COVERAGE_FILE" -outfile "$output_file"
    
    log_success "LCOV 覆盖率报告: $output_file"
}

# 计算覆盖率统计
calculate_coverage_stats() {
    log_step "计算覆盖率统计..."
    
    local stats_file="$REPORTS_DIR/coverage-stats.json"
    
    # 从覆盖率文件中提取统计信息
    local total_coverage=$(go tool cover -func="$COVERAGE_FILE" | grep "total:" | awk '{print $3}' | sed 's/%//')
    local total_lines=$(go tool cover -func="$COVERAGE_FILE" | grep -v "total:" | awk '{sum += $2} END {print sum}')
    local covered_lines=$(go tool cover -func="$COVERAGE_FILE" | grep -v "total:" | awk '{sum += $3} END {print sum}')
    
    # 按包统计覆盖率
    local package_stats=$(go tool cover -func="$COVERAGE_FILE" | grep -v "total:" | awk '
    {
        package = $1
        gsub(/\/[^\/]*$/, "", package)  # 提取包名
        lines[package] += $2
        covered[package] += $3
    }
    END {
        for (pkg in lines) {
            if (lines[pkg] > 0) {
                coverage = (covered[pkg] / lines[pkg]) * 100
                printf "%s:%.1f ", pkg, coverage
            }
        }
    }')
    
    # 生成 JSON 统计报告
    cat > "$stats_file" << EOF
{
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "total_coverage": ${total_coverage:-0},
  "total_lines": ${total_lines:-0},
  "covered_lines": ${covered_lines:-0},
  "threshold": $COVERAGE_THRESHOLD,
  "threshold_met": $([ "${total_coverage:-0}" -ge "$COVERAGE_THRESHOLD" ] && echo "true" || echo "false"),
  "packages": {
EOF
    
    # 添加包级别的统计
    local first=true
    for stat in $package_stats; do
        local pkg=$(echo "$stat" | cut -d: -f1)
        local coverage=$(echo "$stat" | cut -d: -f2)
        
        if [ "$first" = true ]; then
            first=false
        else
            echo "," >> "$stats_file"
        fi
        
        echo -n "    \"$pkg\": $coverage" >> "$stats_file"
    done
    
    echo "" >> "$stats_file"
    echo "  }" >> "$stats_file"
    echo "}" >> "$stats_file"
    
    log_info "覆盖率统计: $stats_file"
    
    # 显示统计摘要
    echo ""
    echo "=== 覆盖率统计 ==="
    echo "总覆盖率: ${total_coverage:-0}%"
    echo "总行数: ${total_lines:-0}"
    echo "覆盖行数: ${covered_lines:-0}"
    echo "阈值: ${COVERAGE_THRESHOLD}%"
    echo "是否达标: $([ "${total_coverage:-0}" -ge "$COVERAGE_THRESHOLD" ] && echo "是" || echo "否")"
    echo ""
    
    # 返回是否达到阈值
    if [ "${total_coverage:-0}" -ge "$COVERAGE_THRESHOLD" ]; then
        return 0
    else
        return 1
    fi
}

# 生成覆盖率徽章
generate_coverage_badge() {
    log_step "生成覆盖率徽章..."
    
    local total_coverage=$(go tool cover -func="$COVERAGE_FILE" | grep "total:" | awk '{print $3}' | sed 's/%//')
    local badge_file="$REPORTS_DIR/coverage-badge.svg"
    
    # 确定徽章颜色
    local color="red"
    if [ "${total_coverage:-0}" -ge 90 ]; then
        color="brightgreen"
    elif [ "${total_coverage:-0}" -ge 80 ]; then
        color="green"
    elif [ "${total_coverage:-0}" -ge 70 ]; then
        color="yellow"
    elif [ "${total_coverage:-0}" -ge 60 ]; then
        color="orange"
    fi
    
    # 生成 SVG 徽章（简单版本）
    cat > "$badge_file" << EOF
<svg xmlns="http://www.w3.org/2000/svg" width="104" height="20">
  <linearGradient id="b" x2="0" y2="100%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <mask id="a">
    <rect width="104" height="20" rx="3" fill="#fff"/>
  </mask>
  <g mask="url(#a)">
    <path fill="#555" d="M0 0h63v20H0z"/>
    <path fill="$color" d="M63 0h41v20H63z"/>
    <path fill="url(#b)" d="M0 0h104v20H0z"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
    <text x="31.5" y="15" fill="#010101" fill-opacity=".3">coverage</text>
    <text x="31.5" y="14">coverage</text>
    <text x="82.5" y="15" fill="#010101" fill-opacity=".3">${total_coverage:-0}%</text>
    <text x="82.5" y="14">${total_coverage:-0}%</text>
  </g>
</svg>
EOF
    
    log_info "覆盖率徽章: $badge_file"
}

# 生成报告索引
generate_report_index() {
    log_step "生成报告索引..."
    
    local index_file="$REPORTS_DIR/index.html"
    local total_coverage=$(go tool cover -func="$COVERAGE_FILE" | grep "total:" | awk '{print $3}' | sed 's/%//')
    
    cat > "$index_file" << EOF
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Websoft9 覆盖率报告</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { text-align: center; margin-bottom: 40px; }
        .coverage-badge { font-size: 24px; padding: 10px 20px; border-radius: 5px; color: white; }
        .coverage-high { background-color: #28a745; }
        .coverage-medium { background-color: #ffc107; color: black; }
        .coverage-low { background-color: #dc3545; }
        .reports { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
        .report-card { border: 1px solid #ddd; border-radius: 8px; padding: 20px; }
        .report-card h3 { margin-top: 0; }
        .report-link { display: inline-block; padding: 8px 16px; background-color: #007bff; color: white; text-decoration: none; border-radius: 4px; }
        .report-link:hover { background-color: #0056b3; }
        .timestamp { text-align: center; color: #666; margin-top: 40px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Websoft9 覆盖率报告</h1>
        <div class="coverage-badge $([ "${total_coverage:-0}" -ge 80 ] && echo "coverage-high" || ([ "${total_coverage:-0}" -ge 60 ] && echo "coverage-medium" || echo "coverage-low"))">
            总覆盖率: ${total_coverage:-0}%
        </div>
    </div>
    
    <div class="reports">
EOF
    
    # 添加可用的报告链接
    if [ -f "$REPORTS_DIR/coverage.html" ]; then
        cat >> "$index_file" << EOF
        <div class="report-card">
            <h3>HTML 报告</h3>
            <p>交互式的 HTML 覆盖率报告，可以查看每个文件的详细覆盖情况。</p>
            <a href="coverage.html" class="report-link">查看 HTML 报告</a>
        </div>
EOF
    fi
    
    if [ -f "$REPORTS_DIR/coverage.txt" ]; then
        cat >> "$index_file" << EOF
        <div class="report-card">
            <h3>文本报告</h3>
            <p>简洁的文本格式覆盖率报告，适合命令行查看。</p>
            <a href="coverage.txt" class="report-link">查看文本报告</a>
        </div>
EOF
    fi
    
    if [ -f "$REPORTS_DIR/coverage.json" ]; then
        cat >> "$index_file" << EOF
        <div class="report-card">
            <h3>JSON 报告</h3>
            <p>机器可读的 JSON 格式报告，适合程序处理。</p>
            <a href="coverage.json" class="report-link">查看 JSON 报告</a>
        </div>
EOF
    fi
    
    if [ -f "$REPORTS_DIR/coverage.xml" ]; then
        cat >> "$index_file" << EOF
        <div class="report-card">
            <h3>XML 报告 (Cobertura)</h3>
            <p>Cobertura 格式的 XML 报告，兼容多种 CI/CD 工具。</p>
            <a href="coverage.xml" class="report-link">查看 XML 报告</a>
        </div>
EOF
    fi
    
    if [ -f "$REPORTS_DIR/coverage-stats.json" ]; then
        cat >> "$index_file" << EOF
        <div class="report-card">
            <h3>统计报告</h3>
            <p>详细的覆盖率统计信息，包括包级别的覆盖率。</p>
            <a href="coverage-stats.json" class="report-link">查看统计报告</a>
        </div>
EOF
    fi
    
    cat >> "$index_file" << EOF
    </div>
    
    <div class="timestamp">
        报告生成时间: $(date)
    </div>
</body>
</html>
EOF
    
    log_success "报告索引: $index_file"
}

# 主函数
main() {
    local generate_html=false
    local generate_json=false
    local generate_xml=false
    local generate_lcov=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -f|--file)
                shift
                if [[ $# -gt 0 ]]; then
                    COVERAGE_FILE="$1"
                else
                    log_error "选项 -f 需要指定文件路径"
                    exit 1
                fi
                ;;
            -o|--output)
                shift
                if [[ $# -gt 0 ]]; then
                    REPORTS_DIR="$1"
                else
                    log_error "选项 -o 需要指定目录路径"
                    exit 1
                fi
                ;;
            -t|--threshold)
                shift
                if [[ $# -gt 0 ]]; then
                    COVERAGE_THRESHOLD="$1"
                else
                    log_error "选项 -t 需要指定阈值"
                    exit 1
                fi
                ;;
            -e|--exclude)
                shift
                if [[ $# -gt 0 ]]; then
                    EXCLUDE_PATTERNS="$1"
                else
                    log_error "选项 -e 需要指定排除模式"
                    exit 1
                fi
                ;;
            --html)
                generate_html=true
                ;;
            --json)
                generate_json=true
                ;;
            --xml)
                generate_xml=true
                ;;
            --lcov)
                generate_lcov=true
                ;;
            --all)
                generate_html=true
                generate_json=true
                generate_xml=true
                generate_lcov=true
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
    
    # 设置全局变量
    GENERATE_HTML=$generate_html
    GENERATE_JSON=$generate_json
    GENERATE_XML=$generate_xml
    GENERATE_LCOV=$generate_lcov
    
    # 检查依赖
    check_dependencies
    
    # 设置输出目录
    setup_output_dir
    
    # 检查覆盖率文件
    check_coverage_file
    
    # 过滤覆盖率数据
    filter_coverage_data
    
    # 生成各种格式的报告
    generate_basic_report
    generate_html_report
    generate_json_report
    generate_xml_report
    generate_lcov_report
    
    # 计算统计信息
    local coverage_ok=true
    if ! calculate_coverage_stats; then
        coverage_ok=false
    fi
    
    # 生成徽章和索引
    generate_coverage_badge
    generate_report_index
    
    # 总结
    echo ""
    echo "=== 报告生成完成 ==="
    echo "输出目录: $REPORTS_DIR"
    echo "报告索引: $REPORTS_DIR/index.html"
    echo ""
    
    if [ "$coverage_ok" = true ]; then
        log_success "覆盖率达到要求 (≥ ${COVERAGE_THRESHOLD}%)"
        exit 0
    else
        log_warning "覆盖率未达到要求 (< ${COVERAGE_THRESHOLD}%)"
        exit 1
    fi
}

# 如果脚本被直接执行
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi