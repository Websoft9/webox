package config

import (
	"api-service/internal/constants"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig    `mapstructure:"server"`
	Database DatabaseConfig  `mapstructure:"database"`
	Redis    RedisConfig     `mapstructure:"redis"`
	InfluxDB InfluxDBConfig  `mapstructure:"influxdb"`
	JWT      JWTConfig       `mapstructure:"jwt"`
	GRPC     GRPCConfig      `mapstructure:"grpc"`
	I18n     I18nConfig      `mapstructure:"i18n"`
	Email    EmailConfigMain `mapstructure:"email"`
	App      AppConfig       `mapstructure:"app"`
}

type ServerConfig struct {
	Port   string             `mapstructure:"port"`
	Mode   string             `mapstructure:"mode"`
	Log    LogConfig          `mapstructure:"log"`
	Config ServerConfigParams `mapstructure:"config"`
}

// ServerConfigParams HTTP服务器配置参数
type ServerConfigParams struct {
	ReadTimeout       int `mapstructure:"read_timeout"`        // 读取超时时间(秒)
	WriteTimeout      int `mapstructure:"write_timeout"`       // 写入超时时间(秒)
	IdleTimeout       int `mapstructure:"idle_timeout"`        // 空闲超时时间(秒)
	ReadHeaderTimeout int `mapstructure:"read_header_timeout"` // 读取头部超时时间(秒)
	MaxHeaderBytes    int `mapstructure:"max_header_bytes"`    // 最大头部字节数
}

// LogConfig 日志配置
type LogConfig struct {
	LogLevel      string `mapstructure:"log_leve"` // 注意：保持与配置文件中的拼写一致
	LogPath       string `mapstructure:"log_path"`
	LogMaxSize    int    `mapstructure:"log_max_size"` // MB
	LogMaxBackups int    `mapstructure:"log_max_backups"`
	LogMaxAge     int    `mapstructure:"log_max_age"` // days
	LogCompress   bool   `mapstructure:"log_compress"`
}

type DatabaseConfig struct {
	// Database type: sqlite, mysql, postgres
	Type string `mapstructure:"type"`

	// SQLite specific
	Path string `mapstructure:"path"`

	// MySQL/PostgreSQL specific
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`

	// Connection pool settings
	MaxIdleConns    int `mapstructure:"max_idle_conns"`
	MaxOpenConns    int `mapstructure:"max_open_conns"`
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`

	// SSL/TLS settings for MySQL/PostgreSQL
	SSLMode   string `mapstructure:"ssl_mode"`
	SSLCert   string `mapstructure:"ssl_cert"`
	SSLKey    string `mapstructure:"ssl_key"`
	SSLRootCA string `mapstructure:"ssl_rootca"`

	// Connection timeout in seconds
	ConnectTimeout int `mapstructure:"connect_timeout"`

	// Charset for MySQL
	Charset string `mapstructure:"charset"`

	// Timezone for MySQL
	Timezone string `mapstructure:"timezone"`
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

type I18nConfig struct {
	DefaultLanguage    string   `mapstructure:"default_language"`
	SupportedLanguages []string `mapstructure:"supported_languages"`
}

type EmailConfigMain struct {
	SMTP SMTPConfigMain `mapstructure:"smtp"`
}

type SMTPConfigMain struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
	UseTLS   bool   `mapstructure:"use_tls"`
}

type AppConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Name    string `mapstructure:"name"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvPrefix("WEBSOFT9")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables to config keys
	_ = viper.BindEnv("database.type", "WEBSOFT9_DB_TYPE")
	_ = viper.BindEnv("database.host", "WEBSOFT9_DB_HOST")
	_ = viper.BindEnv("database.port", "WEBSOFT9_DB_PORT")
	_ = viper.BindEnv("database.database", "WEBSOFT9_DB_NAME")
	_ = viper.BindEnv("database.username", "WEBSOFT9_DB_USER")
	_ = viper.BindEnv("database.password", "WEBSOFT9_DB_PASSWORD")
	_ = viper.BindEnv("database.path", "WEBSOFT9_DB_PATH")
	_ = viper.BindEnv("database.ssl_mode", "WEBSOFT9_DB_SSL_MODE")
	_ = viper.BindEnv("database.max_idle_conns", "WEBSOFT9_DB_MAX_IDLE_CONNS")
	_ = viper.BindEnv("database.max_open_conns", "WEBSOFT9_DB_MAX_OPEN_CONNS")
	_ = viper.BindEnv("database.conn_max_lifetime", "WEBSOFT9_DB_CONN_MAX_LIFETIME")

	_ = viper.BindEnv("redis.host", "WEBSOFT9_REDIS_HOST")
	_ = viper.BindEnv("redis.port", "WEBSOFT9_REDIS_PORT")
	_ = viper.BindEnv("redis.password", "WEBSOFT9_REDIS_PASSWORD")
	_ = viper.BindEnv("redis.db", "WEBSOFT9_REDIS_DB")

	_ = viper.BindEnv("jwt.secret", "WEBSOFT9_JWT_SECRET")
	_ = viper.BindEnv("jwt.expire_time", "WEBSOFT9_JWT_EXPIRE_TIME")

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

	// Server config defaults
	viper.SetDefault("server.config.read_timeout", int(constants.DefaultReadTimeout.Seconds()))
	viper.SetDefault("server.config.write_timeout", int(constants.DefaultWriteTimeout.Seconds()))
	viper.SetDefault("server.config.idle_timeout", int(constants.DefaultIdleTimeout.Seconds()))
	viper.SetDefault("server.config.read_header_timeout", int(constants.DefaultReadHeaderTimeout.Seconds()))
	viper.SetDefault("server.config.max_header_bytes", constants.DefaultMaxHeaderBytes)

	// Log defaults
	viper.SetDefault("server.log.log_leve", constants.DefaultLogLevel)
	viper.SetDefault("server.log.log_path", "./logs/websoft9.log")
	viper.SetDefault("server.log.log_max_size", constants.DefaultLogMaxSize)
	viper.SetDefault("server.log.log_max_backups", constants.DefaultLogMaxBackups)
	viper.SetDefault("server.log.log_max_age", constants.DefaultLogMaxAge)
	viper.SetDefault("server.log.log_compress", true)

	// Database defaults
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.path", "./data/websoft9.db")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", constants.DefaultMySQLPort)
	viper.SetDefault("database.database", "websoft9")
	viper.SetDefault("database.username", "websoft9")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.max_idle_conns", constants.DefaultMaxIdleConns)
	viper.SetDefault("database.max_open_conns", constants.DefaultMaxOpenConns)
	viper.SetDefault("database.conn_max_lifetime", constants.DefaultConnMaxLifetime)
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.connect_timeout", constants.DefaultConnectTimeout)
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.timezone", "Local")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.db", 0)

	// JWT defaults
	viper.SetDefault("jwt.secret", "change-this-secret-key-in-production")
	viper.SetDefault("jwt.expire_time", constants.DefaultJWTExpireTime)

	// gRPC defaults
	viper.SetDefault("grpc.port", "9090")

	// i18n defaults
	viper.SetDefault("i18n.default_language", "en-US")
	viper.SetDefault("i18n.supported_languages", []string{"en-US", "zh-CN"})

	// Email defaults
	viper.SetDefault("email.smtp.host", "smtp.gmail.com")
	viper.SetDefault("email.smtp.port", constants.SMTPDefaultPort)
	viper.SetDefault("email.smtp.username", "")
	viper.SetDefault("email.smtp.password", "")
	viper.SetDefault("email.smtp.from", "noreply@websoft9.com")
	viper.SetDefault("email.smtp.use_tls", true)

	// App defaults
	viper.SetDefault("app.base_url", "http://localhost:3000")
	viper.SetDefault("app.name", "Websoft9")
}
