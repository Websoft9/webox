package logger

import (
	"context"
	"io"
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

// Uint 创建无符号整数字段
func Uint(key string, value uint) Field {
	return Field{Key: key, Value: value}
}

// ErrorField 创建错误字段
func ErrorField(err error) Field {
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

func ErrorLog(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Fatal(msg, fields...)
	}
}
