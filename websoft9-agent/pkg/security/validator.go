package security

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// CommandValidator 命令验证器
type CommandValidator struct {
	allowedCommands map[string]bool
	allowedServices map[string]bool
}

// NewCommandValidator 创建命令验证器
func NewCommandValidator() *CommandValidator {
	return &CommandValidator{
		allowedCommands: map[string]bool{
			"systemctl": true,
			"docker":    true,
			"ps":        true,
			"top":       true,
			"df":        true,
			"free":      true,
			"uptime":    true,
			"whoami":    true,
			"id":        true,
			"pwd":       true,
			"ls":        true,
			"cat":       true,
			"grep":      true,
			"awk":       true,
			"sed":       true,
			"tail":      true,
			"head":      true,
		},
		allowedServices: map[string]bool{
			"nginx":     true,
			"apache2":   true,
			"mysql":     true,
			"redis":     true,
			"docker":    true,
			"ssh":       true,
			"cron":      true,
			"rsyslog":   true,
			"firewalld": true,
			"iptables":  true,
		},
	}
}

// ValidateCommand 验证命令是否安全
func (v *CommandValidator) ValidateCommand(command string) error {
	if command == "" {
		return fmt.Errorf("命令不能为空")
	}

	// 检查危险字符
	dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "{", "}", "[", "]", "<", ">", "\"", "'"}
	for _, char := range dangerousChars {
		if strings.Contains(command, char) {
			return fmt.Errorf("命令包含危险字符: %s", char)
		}
	}

	// 解析命令的第一个部分（实际的命令）
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("无效的命令格式")
	}

	baseCommand := parts[0]

	// 检查命令是否在白名单中
	if !v.allowedCommands[baseCommand] {
		return fmt.Errorf("命令 '%s' 不在允许列表中", baseCommand)
	}

	// 对于 systemctl 命令，进行额外验证
	if baseCommand == "systemctl" && len(parts) >= 3 {
		action := parts[1]
		serviceName := parts[2]

		if err := v.ValidateSystemctlAction(action); err != nil {
			return err
		}

		if err := v.ValidateServiceName(serviceName); err != nil {
			return err
		}
	}

	return nil
}

// ValidateSystemctlAction 验证 systemctl 操作
func (v *CommandValidator) ValidateSystemctlAction(action string) error {
	allowedActions := map[string]bool{
		"start":   true,
		"stop":    true,
		"restart": true,
		"status":  true,
		"enable":  true,
		"disable": true,
		"reload":  true,
	}

	if !allowedActions[action] {
		return fmt.Errorf("systemctl 操作 '%s' 不被允许", action)
	}

	return nil
}

// ValidateServiceName 验证服务名称
func (v *CommandValidator) ValidateServiceName(serviceName string) error {
	if serviceName == "" {
		return fmt.Errorf("服务名称不能为空")
	}

	// 服务名只能包含字母、数字、连字符和下划线
	validServiceName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validServiceName.MatchString(serviceName) {
		return fmt.Errorf("服务名称格式无效: %s", serviceName)
	}

	// 检查服务名长度
	if len(serviceName) > 64 {
		return fmt.Errorf("服务名称过长: %s", serviceName)
	}

	// 可选：检查服务是否在允许列表中（根据需要启用）
	// if !v.allowedServices[serviceName] {
	//     return fmt.Errorf("服务 '%s' 不在允许列表中", serviceName)
	// }

	return nil
}

// PathValidator 路径验证器
type PathValidator struct {
	allowedDirs []string
}

// NewPathValidator 创建路径验证器
func NewPathValidator() *PathValidator {
	return &PathValidator{
		allowedDirs: []string{
			"/etc/websoft9",
			"/var/lib/websoft9",
			"/opt/websoft9",
			"/tmp/websoft9",
		},
	}
}

// ValidateConfigPath 验证配置文件路径
func (v *PathValidator) ValidateConfigPath(configPath string) error {
	if configPath == "" {
		return fmt.Errorf("配置文件路径不能为空")
	}

	// 清理路径，解析符号链接和相对路径
	cleanPath, err := filepath.Abs(filepath.Clean(configPath))
	if err != nil {
		return fmt.Errorf("无法解析配置文件路径: %v", err)
	}

	// 检查路径是否在允许的目录内
	allowed := false
	for _, allowedDir := range v.allowedDirs {
		if strings.HasPrefix(cleanPath, allowedDir) {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("配置文件路径不在允许的目录内: %s", cleanPath)
	}

	// 检查文件扩展名
	ext := filepath.Ext(cleanPath)
	allowedExts := map[string]bool{
		".yaml": true,
		".yml":  true,
		".json": true,
		".toml": true,
		".conf": true,
	}

	if !allowedExts[ext] {
		return fmt.Errorf("不支持的配置文件格式: %s", ext)
	}

	return nil
}

// SanitizeInput 清理输入字符串
func SanitizeInput(input string) string {
	// 移除控制字符
	input = regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(input, "")

	// 限制长度
	if len(input) > 1024 {
		input = input[:1024]
	}

	// 去除首尾空白
	input = strings.TrimSpace(input)

	return input
}

// ValidateTaskParams 验证任务参数
func ValidateTaskParams(params map[string]interface{}) error {
	for key, value := range params {
		// 验证键名
		if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(key) {
			return fmt.Errorf("无效的参数键名: %s", key)
		}

		// 验证值的类型和内容
		switch v := value.(type) {
		case string:
			if len(v) > 4096 {
				return fmt.Errorf("参数 %s 的值过长", key)
			}
		case map[string]interface{}:
			// 递归验证嵌套参数
			if err := ValidateTaskParams(v); err != nil {
				return err
			}
		}
	}

	return nil
}
