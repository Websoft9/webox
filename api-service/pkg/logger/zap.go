package logger

import (
	"api-service/internal/constants"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Output type constants
const (
	OutputTypeStdout = "stdout"
	OutputTypeStderr = "stderr"
	OutputTypeFile   = "file"
)

// 常量定义
const (
	DefaultFilePermission = 0o666 // 默认文件权限
	ContextFieldCapacity  = 4     // 上下文字段容量
)

// 允许的日志文件目录前缀（安全路径）
var allowedLogPaths = []string{
	"/var/log/",
	"/tmp/",
	"./logs/",
	"./data/",
}

// OutputConfig 输出配置
type OutputConfig struct {
	Path       string
	MaxSize    int // MB
	MaxBackups int
	MaxAge     int // days
	Compress   bool
}

// ZapLogger Zap日志实现
type ZapLogger struct {
	logger      *zap.Logger
	atomicLevel zap.AtomicLevel // 支持动态级别调整
}

// NewZapLogger 创建新的Zap日志实例
func NewZapLogger(level Level, output io.Writer) Logger {
	encoderConfig := createEncoderConfig()
	atomicLevel := zap.NewAtomicLevelAt(convertToZapLevel(level))

	if output == nil {
		output = os.Stdout
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(output),
		atomicLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{
		logger:      logger,
		atomicLevel: atomicLevel,
	}
}

// NewDefaultZapLogger 创建默认Zap日志实例
func NewDefaultZapLogger() Logger {
	return NewZapLogger(InfoLevel, os.Stdout)
}

// NewZapLoggerWithConfig 根据配置创建Zap日志实例
func NewZapLoggerWithConfig(config *Config) Logger {
	output := createOutput(&OutputConfig{
		Path:       config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}, config.Output)

	return NewZapLogger(config.Level, output)
}

// NewZapLoggerWithServerConfig 根据服务器配置创建Zap日志实例
func NewZapLoggerWithServerConfig(logPath, logLevel string, maxSize, maxBackups, maxAge int, compress bool) Logger {
	level := ParseLevel(logLevel)
	output := createOutput(&OutputConfig{
		Path:       logPath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}, "")

	return NewZapLogger(level, output)
}

// createEncoderConfig 创建编码器配置
func createEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// createOutput 创建输出Writer，支持文件轮转
func createOutput(config *OutputConfig, outputType string) io.Writer {
	switch outputType {
	case OutputTypeStderr:
		return os.Stderr
	case OutputTypeStdout, "":
		// 如果明确指定为stdout或配置路径为空/stdout，返回标准输出
		if config.Path == "" || config.Path == OutputTypeStdout {
			return os.Stdout
		}
		// 如果配置了具体的文件路径，创建文件输出
		return createFileOutput(config)
	case OutputTypeFile:
		return createFileOutput(config)
	default:
		return os.Stdout
	}
}

// createFileOutput 创建文件输出，支持轮转
func createFileOutput(config *OutputConfig) io.Writer {
	if config.Path == "" {
		return os.Stdout
	}

	// 防止使用保留的输出名称作为文件名
	baseName := filepath.Base(config.Path)
	if baseName == OutputTypeStdout || baseName == OutputTypeStderr {
		return os.Stdout
	}

	// 验证文件路径安全性
	if err := validateFilePath(config.Path); err != nil {
		return os.Stdout
	}

	// 清理和规范化路径，防止路径遍历攻击
	cleanPath := filepath.Clean(config.Path)

	// 确保日志目录存在
	if err := ensureLogDir(cleanPath); err != nil {
		return os.Stdout
	}

	// 如果配置了轮转参数，使用 lumberjack
	if config.MaxSize > 0 || config.MaxBackups > 0 || config.MaxAge > 0 {
		return &lumberjack.Logger{
			Filename:   cleanPath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
			LocalTime:  true,
		}
	}

	// 普通文件输出 - 路径已通过 validateFilePath 验证和 filepath.Clean 清理
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	// #nosec G304 - 路径已经过 validateFilePath 验证并使用 filepath.Clean 清理
	if file, err := os.OpenFile(cleanPath, flags, DefaultFilePermission); err == nil {
		return file
	}

	return os.Stdout
}

// ensureLogDir 确保日志目录存在
func ensureLogDir(filename string) error {
	dir := filepath.Dir(filename)
	if dir == "." || dir == "" {
		return nil
	}

	if err := os.MkdirAll(dir, constants.DefaultLogDirPerm); err != nil {
		return fmt.Errorf("failed to create log directory %s: %w", dir, err)
	}

	return nil
}

// validateFilePath 验证文件路径安全性
func validateFilePath(filename string) error {
	if filename == "" {
		return nil
	}

	// 检查原始路径中的路径遍历字符（在清理之前）
	if strings.Contains(filename, "..") {
		return fmt.Errorf("path traversal detected: %s", filename)
	}

	// 清理路径
	cleanPath := filepath.Clean(filename)

	// 检查绝对路径权限
	if filepath.IsAbs(cleanPath) {
		for _, allowedPath := range allowedLogPaths {
			if strings.HasPrefix(cleanPath, allowedPath) {
				return nil
			}
		}
		return fmt.Errorf("absolute path not allowed: %s", filename)
	}

	// 检查相对路径安全性
	if strings.HasPrefix(cleanPath, "/") || strings.HasPrefix(cleanPath, "\\") {
		return fmt.Errorf("invalid relative path: %s", filename)
	}

	return nil
}

// convertToZapLevel 转换到Zap日志级别
func convertToZapLevel(level Level) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Debug 调试级别日志
func (z *ZapLogger) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, z.convertFields(fields...)...)
}

// Info 信息级别日志
func (z *ZapLogger) Info(msg string, fields ...Field) {
	z.logger.Info(msg, z.convertFields(fields...)...)
}

// Warn 警告级别日志
func (z *ZapLogger) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, z.convertFields(fields...)...)
}

// Error 错误级别日志
func (z *ZapLogger) Error(msg string, fields ...Field) {
	z.logger.Error(msg, z.convertFields(fields...)...)
}

// Fatal 致命级别日志
func (z *ZapLogger) Fatal(msg string, fields ...Field) {
	z.logger.Fatal(msg, z.convertFields(fields...)...)
}

// DebugContext 带上下文的调试日志
func (z *ZapLogger) DebugContext(ctx context.Context, msg string, fields ...Field) {
	z.Debug(msg, append(fields, z.extractContextFields(ctx)...)...)
}

// InfoContext 带上下文的信息日志
func (z *ZapLogger) InfoContext(ctx context.Context, msg string, fields ...Field) {
	z.Info(msg, append(fields, z.extractContextFields(ctx)...)...)
}

// WarnContext 带上下文的警告日志
func (z *ZapLogger) WarnContext(ctx context.Context, msg string, fields ...Field) {
	z.Warn(msg, append(fields, z.extractContextFields(ctx)...)...)
}

// ErrorContext 带上下文的错误日志
func (z *ZapLogger) ErrorContext(ctx context.Context, msg string, fields ...Field) {
	z.Error(msg, append(fields, z.extractContextFields(ctx)...)...)
}

// IsDebugEnabled 检查是否启用调试日志
func (z *ZapLogger) IsDebugEnabled() bool {
	return z.logger.Core().Enabled(zapcore.DebugLevel)
}

// IsInfoEnabled 检查是否启用信息日志
func (z *ZapLogger) IsInfoEnabled() bool {
	return z.logger.Core().Enabled(zapcore.InfoLevel)
}

// IsWarnEnabled 检查是否启用警告日志
func (z *ZapLogger) IsWarnEnabled() bool {
	return z.logger.Core().Enabled(zapcore.WarnLevel)
}

// IsErrorEnabled 检查是否启用错误日志
func (z *ZapLogger) IsErrorEnabled() bool {
	return z.logger.Core().Enabled(zapcore.ErrorLevel)
}

// WithFields 添加字段
func (z *ZapLogger) WithFields(fields ...Field) Logger {
	if len(fields) == 0 {
		return z
	}

	return &ZapLogger{
		logger:      z.logger.With(z.convertFields(fields...)...),
		atomicLevel: z.atomicLevel,
	}
}

// WithContext 添加上下文
func (z *ZapLogger) WithContext(ctx context.Context) Logger {
	return z.WithFields(z.extractContextFields(ctx)...)
}

// SetLevel 设置日志级别
func (z *ZapLogger) SetLevel(level Level) {
	z.atomicLevel.SetLevel(convertToZapLevel(level))
}

// SetOutput 设置输出（Zap不支持动态修改输出）
func (z *ZapLogger) SetOutput(w io.Writer) {
	// Zap的输出在创建时设置，无法动态修改
}

// convertFields 转换字段格式
func (z *ZapLogger) convertFields(fields ...Field) []zap.Field {
	if len(fields) == 0 {
		return nil
	}

	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, z.convertField(field))
	}
	return zapFields
}

// convertField 转换单个字段
func (z *ZapLogger) convertField(field Field) zap.Field {
	switch v := field.Value.(type) {
	case string:
		return zap.String(field.Key, v)
	case int:
		return zap.Int(field.Key, v)
	case int64:
		return zap.Int64(field.Key, v)
	case uint:
		return zap.Uint(field.Key, v)
	case float64:
		return zap.Float64(field.Key, v)
	case bool:
		return zap.Bool(field.Key, v)
	case error:
		if v != nil {
			return zap.Error(v)
		}
		return zap.String(field.Key, "")
	default:
		return zap.Any(field.Key, v)
	}
}

// extractContextFields 从上下文中提取字段
func (z *ZapLogger) extractContextFields(ctx context.Context) []Field {
	if ctx == nil {
		return nil
	}

	fields := make([]Field, 0, ContextFieldCapacity)

	// 定义要提取的字段名和对应的键
	contextKeys := map[string]string{
		"request_id": "request_id",
		"user_id":    "user_id",
		"username":   "username",
		"role":       "role",
		"trace_id":   "trace_id",
	}

	for ctxKey, fieldKey := range contextKeys {
		if value := ctx.Value(ctxKey); value != nil {
			if str, ok := value.(string); ok && str != "" {
				fields = append(fields, String(fieldKey, str))
			} else if ctxKey == "user_id" {
				// user_id 可能不是字符串类型
				fields = append(fields, Any(fieldKey, value))
			}
		}
	}

	return fields
}
