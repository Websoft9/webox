package config

import (
	"fmt"
	"os"
	"websoft9-agent/internal/constants"
	"websoft9-agent/pkg/security"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Config Agent 配置结构
type Config struct {
	Server ServerConfig `yaml:"server"`
	Redis  RedisConfig  `yaml:"redis"`
	Log    LogConfig    `yaml:"log"`
	Agent  AgentConfig  `yaml:"agent"`
}

// ServerConfig 服务端配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	TLS  bool   `yaml:"tls"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

// AgentConfig Agent 配置
type AgentConfig struct {
	ID                string `yaml:"id"`
	HeartbeatInterval int    `yaml:"heartbeat_interval"`
	MonitorInterval   int    `yaml:"monitor_interval"`
	WorkDir           string `yaml:"work_dir"`
}

// Load 加载配置文件
func Load(configFile string) (*Config, error) {
	// 验证配置文件路径安全性
	pathValidator := security.NewPathValidator()
	if err := pathValidator.ValidateConfigPath(configFile); err != nil {
		logrus.WithFields(logrus.Fields{
			"config_file": configFile,
			"error":       err.Error(),
		}).Error("Security audit: invalid config file path")
		return nil, fmt.Errorf("配置文件路径验证失败: %v", err)
	}

	// 记录安全审计日志
	logrus.WithFields(logrus.Fields{
		"config_file": configFile,
		"action":      "config_file_load",
	}).Info("Security audit: loading config file")

	// #nosec G304 - Config file path is validated by security.PathValidator
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 设置默认值
	setDefaults(&config)

	// 验证配置内容
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	return &config, nil
}

func setDefaults(config *Config) {
	if config.Server.Host == "" {
		config.Server.Host = constants.DefaultHost
	}
	if config.Server.Port == 0 {
		config.Server.Port = 9090
	}
	if config.Redis.Host == "" {
		config.Redis.Host = constants.DefaultHost
	}
	if config.Redis.Port == 0 {
		config.Redis.Port = 6379
	}
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}
	if config.Agent.HeartbeatInterval == 0 {
		config.Agent.HeartbeatInterval = 30
	}
	if config.Agent.MonitorInterval == 0 {
		config.Agent.MonitorInterval = 60
	}
	if config.Agent.WorkDir == "" {
		config.Agent.WorkDir = "/var/lib/websoft9/agent"
	}
}

// validateConfig 验证配置内容
func validateConfig(config *Config) error {
	// 验证服务器配置
	if config.Server.Host == "" {
		return fmt.Errorf("服务器地址不能为空")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("服务器端口无效: %d", config.Server.Port)
	}

	// 验证 Redis 配置
	if config.Redis.Host == "" {
		return fmt.Errorf("Redis 地址不能为空")
	}
	if config.Redis.Port <= 0 || config.Redis.Port > 65535 {
		return fmt.Errorf("Redis 端口无效: %d", config.Redis.Port)
	}

	// 验证 Agent 配置
	if config.Agent.ID == "" {
		return fmt.Errorf("Agent ID 不能为空")
	}
	if config.Agent.HeartbeatInterval <= 0 {
		return fmt.Errorf("心跳间隔必须大于 0")
	}
	if config.Agent.MonitorInterval <= 0 {
		return fmt.Errorf("监控间隔必须大于 0")
	}

	// 验证日志级别
	validLogLevels := map[string]bool{
		"trace": true,
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
		"panic": true,
	}
	if !validLogLevels[config.Log.Level] {
		return fmt.Errorf("无效的日志级别: %s", config.Log.Level)
	}

	return nil
}
