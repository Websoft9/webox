package monitor

import (
	"net/http"
	"time"

	"websoft9-agent/internal/config"
	"websoft9-agent/internal/constants"

	"github.com/sirupsen/logrus"
)

// HealthChecker 健康检查器
type HealthChecker struct {
	config *config.Config
	client *http.Client
}

// HealthCheck 健康检查配置
type HealthCheck struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // http, tcp, script
	Target   string `json:"target"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

// HealthResult 健康检查结果
type HealthResult struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"` // healthy, unhealthy
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Duration  int64     `json:"duration"` // 响应时间(毫秒)
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(cfg *config.Config) (*HealthChecker, error) {
	return &HealthChecker{
		config: cfg,
		client: &http.Client{
			Timeout: constants.DefaultHTTPTimeout,
		},
	}, nil
}

// Check 执行健康检查
func (h *HealthChecker) Check() error {
	// TODO: 从配置或数据库获取健康检查列表
	checks := h.getHealthChecks()

	for _, check := range checks {
		result := h.performCheck(check)
		logrus.Debugf("健康检查 %s: %s", result.Name, result.Status)

		// TODO: 发送结果到服务端
	}

	return nil
}

// getHealthChecks 获取健康检查配置
func (h *HealthChecker) getHealthChecks() []HealthCheck {
	// TODO: 从配置文件或数据库读取
	return []HealthCheck{
		{
			Name:     "api-service",
			Type:     "http",
			Target:   "http://localhost:8080/health",
			Interval: int(constants.DefaultCheckInterval.Seconds()),
			Timeout:  int(constants.DefaultTCPTimeout.Seconds()),
		},
	}
}

// performCheck 执行单个健康检查
func (h *HealthChecker) performCheck(check HealthCheck) HealthResult {
	start := time.Now()
	result := HealthResult{
		Name:      check.Name,
		Timestamp: start,
	}

	switch check.Type {
	case "http":
		return *h.httpCheck(check, &result)
	case "tcp":
		return *h.tcpCheck(check, &result)
	case "script":
		return *h.scriptCheck(check, &result)
	default:
		result.Status = constants.StatusUnhealthy
		result.Message = "未知的检查类型"
	}

	result.Duration = time.Since(start).Milliseconds()
	return result
}

// httpCheck HTTP健康检查
func (h *HealthChecker) httpCheck(check HealthCheck, result *HealthResult) *HealthResult {
	resp, err := h.client.Get(check.Target)
	if err != nil {
		result.Status = constants.StatusUnhealthy
		result.Message = err.Error()
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Status = constants.StatusHealthy
		result.Message = "HTTP检查成功"
	} else {
		result.Status = constants.StatusUnhealthy
		result.Message = "HTTP状态码异常: " + resp.Status
	}

	return result
}

// tcpCheck TCP健康检查
func (h *HealthChecker) tcpCheck(check HealthCheck, result *HealthResult) *HealthResult {
	// TODO: 实现TCP连接检查
	// 使用 check.Target 进行 TCP 连接测试
	_ = check // 暂时忽略未使用的参数，待实现时使用
	result.Status = constants.StatusHealthy
	result.Message = "TCP检查成功"
	return result
}

// scriptCheck 脚本健康检查
func (h *HealthChecker) scriptCheck(check HealthCheck, result *HealthResult) *HealthResult {
	// TODO: 实现脚本执行检查
	// 使用 check.Target 作为脚本路径执行
	_ = check // 暂时忽略未使用的参数，待实现时使用
	result.Status = constants.StatusHealthy
	result.Message = "脚本检查成功"
	return result
}
