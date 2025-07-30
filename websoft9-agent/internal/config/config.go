package config

import (
	"os"

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
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// 设置默认值
	setDefaults(&config)

	return &config, nil
}

func setDefaults(config *Config) {
	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 9090
	}
	if config.Redis.Host == "" {
		config.Redis.Host = "localhost"
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
