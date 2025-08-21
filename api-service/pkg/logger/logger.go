package logger

import (
	"context"
	"io"
	"time"
)

// Level 日志级别
type Level int

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// 日志级别字符串常量
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
)

// String 返回日志级别的字符串表示
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return LevelDebug
	case InfoLevel:
		return LevelInfo
	case WarnLevel:
		return LevelWarn
	case ErrorLevel:
		return LevelError
	case FatalLevel:
		return LevelFatal
	default:
		return "unknown"
	}
}

// ParseLevel 从字符串解析日志级别
func ParseLevel(s string) Level {
	switch s {
	case LevelDebug:
		return DebugLevel
	case LevelInfo:
		return InfoLevel
	case LevelWarn:
		return WarnLevel
	case LevelError:
		return ErrorLevel
	case LevelFatal:
		return FatalLevel
	default:
		return InfoLevel
	}
}

// Config 日志配置
type Config struct {
	Level      Level  `json:"level" yaml:"level"`             // 日志级别
	Format     string `json:"format" yaml:"format"`           // 日志格式: json, console
	Output     string `json:"output" yaml:"output"`           // 输出目标: stdout, stderr, file
	Filename   string `json:"filename" yaml:"filename"`       // 文件名(当output为file时)
	MaxSize    int    `json:"max_size" yaml:"max_size"`       // 最大文件大小(MB)
	MaxBackups int    `json:"max_backups" yaml:"max_backups"` // 保留文件数
	MaxAge     int    `json:"max_age" yaml:"max_age"`         // 保留天数
	Compress   bool   `json:"compress" yaml:"compress"`       // 是否压缩
}

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
}

// Logger 日志接口
type Logger interface {
	// 基础日志方法
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	// 带上下文的日志方法
	DebugContext(ctx context.Context, msg string, fields ...Field)
	InfoContext(ctx context.Context, msg string, fields ...Field)
	WarnContext(ctx context.Context, msg string, fields ...Field)
	ErrorContext(ctx context.Context, msg string, fields ...Field)

	// 级别检查方法 - 性能优化
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarnEnabled() bool
	IsErrorEnabled() bool

	// 配置方法
	WithFields(fields ...Field) Logger
	WithContext(ctx context.Context) Logger
	SetLevel(level Level)
	SetOutput(w io.Writer)
}

// String 创建字符串字段
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int 创建整数字段
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64 创建64位整数字段
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Uint 创建无符号整数字段
func Uint(key string, value uint) Field {
	return Field{Key: key, Value: value}
}

// Float64 创建浮点数字段
func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

// Bool 创建布尔字段
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Duration 创建时间间隔字段
func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value.String()}
}

// Time 创建时间字段
func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value}
}

// ErrorField 创建错误字段
func ErrorField(err error) Field {
	if err == nil {
		return Field{Key: "error", Value: nil}
	}
	return Field{Key: "error", Value: err.Error()}
}

// Any 创建任意类型字段
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// 全局日志实例
var defaultLogger Logger

// SetDefault 设置默认日志实例
func SetDefault(l Logger) {
	defaultLogger = l
}

// GetDefault 获取默认日志实例
func GetDefault() Logger {
	return defaultLogger
}

// 便捷方法，使用默认日志实例
func Debug(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Fatal(msg, fields...)
	}
}
