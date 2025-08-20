package config

import (
	"api-service/internal/constants"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	InfluxDB InfluxDBConfig `mapstructure:"influxdb"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	GRPC     GRPCConfig     `mapstructure:"grpc"`
	MCP      MCPConfig      `mapstructure:"mcp"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type InfluxDBConfig struct {
	URL    string `mapstructure:"url"`
	Token  string `mapstructure:"token"`
	Org    string `mapstructure:"org"`
	Bucket string `mapstructure:"bucket"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"`
}

type GRPCConfig struct {
	Port string `mapstructure:"port"`
}

// MCPConfig defines MCP (Model Context Protocol) server configuration
type MCPConfig struct {
	Enabled bool                   `mapstructure:"enabled"`
	Servers map[string]MCPServer   `mapstructure:"servers"`
	Config  MCPGlobalConfig        `mapstructure:"config"`
}

// MCPServer defines configuration for individual MCP servers
type MCPServer struct {
	Type        string            `mapstructure:"type"`         // stdio, sse, websocket
	Command     string            `mapstructure:"command"`      // Command to run
	Args        []string          `mapstructure:"args"`         // Command arguments
	Env         map[string]string `mapstructure:"env"`          // Environment variables
	WorkingDir  string            `mapstructure:"working_dir"`  // Working directory
	Enabled     bool              `mapstructure:"enabled"`      // Whether server is enabled
	Gallery     bool              `mapstructure:"gallery"`      // Whether to show in gallery
	Description string            `mapstructure:"description"`  // Server description
	AutoInstall bool              `mapstructure:"auto_install"` // Auto install dependencies
}

// MCPGlobalConfig defines global MCP configuration
type MCPGlobalConfig struct {
	LogLevel       string `mapstructure:"log_level"`        // debug, info, warn, error
	Timeout        int    `mapstructure:"timeout"`          // Connection timeout in seconds
	RetryAttempts  int    `mapstructure:"retry_attempts"`   // Number of retry attempts
	RetryDelay     int    `mapstructure:"retry_delay"`      // Delay between retries in seconds
	ConfigPath     string `mapstructure:"config_path"`      // Path to MCP config file
	AutoGenerate   bool   `mapstructure:"auto_generate"`    // Auto-generate config if missing
	NodeOptions    string `mapstructure:"node_options"`     // Node.js options for npm servers
	PythonPath     string `mapstructure:"python_path"`      // Python path for Python servers
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func setDefaults() {
	viper.SetDefault("server.port", constants.DefaultPort)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.path", "./data/websoft9.db")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("jwt.secret", "change-this-secret-key-in-production")
	viper.SetDefault("jwt.expire_time", constants.DefaultJWTExpireTime)
	viper.SetDefault("grpc.port", "9090")

	// MCP defaults
	viper.SetDefault("mcp.enabled", true)
	viper.SetDefault("mcp.config.log_level", "info")
	viper.SetDefault("mcp.config.timeout", 30)
	viper.SetDefault("mcp.config.retry_attempts", 3)
	viper.SetDefault("mcp.config.retry_delay", 5)
	viper.SetDefault("mcp.config.config_path", "./configs/mcp.json")
	viper.SetDefault("mcp.config.auto_generate", true)
	viper.SetDefault("mcp.config.node_options", "--max-old-space-size=2048")
	viper.SetDefault("mcp.config.python_path", "python3")

	// Default MCP servers
	setDefaultMCPServers()
}

// setDefaultMCPServers sets up default MCP server configurations
func setDefaultMCPServers() {
	// Filesystem server - for file system access
	viper.SetDefault("mcp.servers.filesystem.type", "stdio")
	viper.SetDefault("mcp.servers.filesystem.command", "npx")
	viper.SetDefault("mcp.servers.filesystem.args", []string{"-y", "@modelcontextprotocol/server-filesystem"})
	viper.SetDefault("mcp.servers.filesystem.enabled", true)
	viper.SetDefault("mcp.servers.filesystem.description", "File system access for reading and writing files")
	viper.SetDefault("mcp.servers.filesystem.auto_install", true)
	viper.SetDefault("mcp.servers.filesystem.env", map[string]string{
		"NODE_OPTIONS": "--max-old-space-size=2048",
	})

	// Memory server - for persistent memory across conversations
	viper.SetDefault("mcp.servers.memory.type", "stdio")
	viper.SetDefault("mcp.servers.memory.command", "npx")
	viper.SetDefault("mcp.servers.memory.args", []string{"-y", "@modelcontextprotocol/server-memory"})
	viper.SetDefault("mcp.servers.memory.enabled", true)
	viper.SetDefault("mcp.servers.memory.description", "Persistent memory for storing information across conversations")
	viper.SetDefault("mcp.servers.memory.auto_install", true)

	// Context7 server - for enhanced context management
	viper.SetDefault("mcp.servers.context7.type", "stdio")
	viper.SetDefault("mcp.servers.context7.command", "npx")
	viper.SetDefault("mcp.servers.context7.args", []string{"-y", "@upstash/context7-mcp@latest"})
	viper.SetDefault("mcp.servers.context7.enabled", true)
	viper.SetDefault("mcp.servers.context7.gallery", true)
	viper.SetDefault("mcp.servers.context7.description", "Enhanced context management with Upstash integration")
	viper.SetDefault("mcp.servers.context7.auto_install", true)

	// Git server - for version control operations
	viper.SetDefault("mcp.servers.git.type", "stdio")
	viper.SetDefault("mcp.servers.git.command", "npx")
	viper.SetDefault("mcp.servers.git.args", []string{"-y", "@modelcontextprotocol/server-git"})
	viper.SetDefault("mcp.servers.git.enabled", true)
	viper.SetDefault("mcp.servers.git.description", "Git version control operations")
	viper.SetDefault("mcp.servers.git.auto_install", true)

	// Database server - for database operations
	viper.SetDefault("mcp.servers.database.type", "stdio")
	viper.SetDefault("mcp.servers.database.command", "npx")
	viper.SetDefault("mcp.servers.database.args", []string{"-y", "@modelcontextprotocol/server-sqlite"})
	viper.SetDefault("mcp.servers.database.enabled", true)
	viper.SetDefault("mcp.servers.database.description", "SQLite database operations")
	viper.SetDefault("mcp.servers.database.auto_install", true)
}
