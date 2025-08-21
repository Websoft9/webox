package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

// createFileOutput 创建文件输出，处理安全验证和文件创建
func createFileOutput(filename string) io.Writer {
	if filename == "" {
		return os.Stdout
	}

	// 验证文件路径安全性，防止路径遍历攻击
	if err := validateFilePath(filename); err != nil {
		// 如果路径不安全，回退到标准输出
		return os.Stdout
	}

	// 这里可以扩展支持文件轮转等功能
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	// #nosec G304 - 文件路径已通过 validateFilePath 函数验证安全性
	if file, err := os.OpenFile(filename, flags, DefaultFilePermission); err == nil {
		return file
	}

	// 如果文件创建失败，回退到标准输出
	return os.Stdout
}

// validateFilePath 验证文件路径安全性，防止路径遍历攻击
func validateFilePath(filename string) error {
	// 清理路径
	cleanPath := filepath.Clean(filename)

	// 检查是否包含路径遍历字符
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected in filename: %s", filename)
	}

	// 检查是否为绝对路径且在允许的目录中
	if filepath.IsAbs(cleanPath) {
		allowed := false
		for _, allowedPath := range allowedLogPaths {
			if strings.HasPrefix(cleanPath, allowedPath) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("absolute path not in allowed directories: %s", filename)
		}
	}

	// 检查相对路径是否安全
	if !filepath.IsAbs(cleanPath) {
		// 相对路径应该在当前目录或子目录中
		if strings.HasPrefix(cleanPath, "/") || strings.HasPrefix(cleanPath, "\\") {
			return fmt.Errorf("invalid relative path: %s", filename)
		}
	}

	return nil
}

// ZapLogger Zap日志实现
type ZapLogger struct {
	logger      *zap.Logger
	sugar       *zap.SugaredLogger
	atomicLevel zap.AtomicLevel // 支持动态级别调整
}

// NewZapLogger 创建新的Zap日志实例
func NewZapLogger(level Level, output io.Writer) Logger {
	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
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

	// 设置日志级别 - 使用AtomicLevel支持动态调整
	zapLevel := zapcore.InfoLevel
	switch level {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	case FatalLevel:
		zapLevel = zapcore.FatalLevel
	}

	atomicLevel := zap.NewAtomicLevelAt(zapLevel)

	// 创建核心
	if output == nil {
		output = os.Stdout
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(output),
		atomicLevel,
	)

	// 创建logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{
		logger:      logger,
		sugar:       logger.Sugar(),
		atomicLevel: atomicLevel,
	}
}

// NewDefaultZapLogger 创建默认Zap日志实例
func NewDefaultZapLogger() Logger {
	return NewZapLogger(InfoLevel, os.Stdout)
}

// NewZapLoggerWithConfig 根据配置创建Zap日志实例
func NewZapLoggerWithConfig(config *Config) Logger {
	// 确定输出目标
	var output io.Writer = os.Stdout
	switch config.Output {
	case "stderr":
		output = os.Stderr
	case "stdout", "":
		output = os.Stdout
	case "file":
		output = createFileOutput(config.Filename)
	}

	return NewZapLogger(config.Level, output)
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
	fields = append(fields, z.extractContextFields(ctx)...)
	z.Debug(msg, fields...)
}

// InfoContext 带上下文的信息日志
func (z *ZapLogger) InfoContext(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, z.extractContextFields(ctx)...)
	z.Info(msg, fields...)
}

// WarnContext 带上下文的警告日志
func (z *ZapLogger) WarnContext(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, z.extractContextFields(ctx)...)
	z.Warn(msg, fields...)
}

// ErrorContext 带上下文的错误日志
func (z *ZapLogger) ErrorContext(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, z.extractContextFields(ctx)...)
	z.Error(msg, fields...)
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
		sugar:       z.sugar,
		atomicLevel: z.atomicLevel,
	}
}

// WithContext 添加上下文
func (z *ZapLogger) WithContext(ctx context.Context) Logger {
	return z.WithFields(z.extractContextFields(ctx)...)
}

// SetLevel 设置日志级别
func (z *ZapLogger) SetLevel(level Level) {
	zapLevel := zapcore.InfoLevel
	switch level {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	case FatalLevel:
		zapLevel = zapcore.FatalLevel
	}
	z.atomicLevel.SetLevel(zapLevel)
}

// SetOutput 设置输出
func (z *ZapLogger) SetOutput(w io.Writer) {
	// Zap的输出在创建时设置，这里可以记录但不能动态修改
	// 实际项目中可能需要重新创建logger
}

// convertFields 转换字段格式 - 优化性能
func (z *ZapLogger) convertFields(fields ...Field) []zap.Field {
	if len(fields) == 0 {
		return nil
	}

	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		// 根据类型优化字段转换
		switch v := field.Value.(type) {
		case string:
			zapFields = append(zapFields, zap.String(field.Key, v))
		case int:
			zapFields = append(zapFields, zap.Int(field.Key, v))
		case int64:
			zapFields = append(zapFields, zap.Int64(field.Key, v))
		case uint:
			zapFields = append(zapFields, zap.Uint(field.Key, v))
		case float64:
			zapFields = append(zapFields, zap.Float64(field.Key, v))
		case bool:
			zapFields = append(zapFields, zap.Bool(field.Key, v))
		case error:
			if v != nil {
				zapFields = append(zapFields, zap.Error(v))
			}
		default:
			zapFields = append(zapFields, zap.Any(field.Key, v))
		}
	}
	return zapFields
}

// extractContextFields 从上下文中提取字段 - 优化版本
func (z *ZapLogger) extractContextFields(ctx context.Context) []Field {
	if ctx == nil {
		return nil
	}

	fields := make([]Field, 0, ContextFieldCapacity) // 预分配容量

	// 使用辅助函数提取各种字段
	fields = z.extractRequestID(ctx, fields)
	fields = z.extractUserInfo(ctx, fields)
	fields = z.extractTraceInfo(ctx, fields)

	return fields
}

// extractRequestID 提取请求ID
func (z *ZapLogger) extractRequestID(ctx context.Context, fields []Field) []Field {
	if requestID := ctx.Value("request_id"); requestID != nil {
		if rid, ok := requestID.(string); ok && rid != "" {
			fields = append(fields, String("request_id", rid))
		}
	}
	return fields
}

// extractUserInfo 提取用户相关信息
func (z *ZapLogger) extractUserInfo(ctx context.Context, fields []Field) []Field {
	// 提取用户ID
	if userID := ctx.Value("user_id"); userID != nil {
		fields = append(fields, Any("user_id", userID))
	}

	// 提取用户名
	if username := ctx.Value("username"); username != nil {
		if un, ok := username.(string); ok && un != "" {
			fields = append(fields, String("username", un))
		}
	}

	// 提取用户角色
	if role := ctx.Value("role"); role != nil {
		if r, ok := role.(string); ok && r != "" {
			fields = append(fields, String("role", r))
		}
	}

	return fields
}

// extractTraceInfo 提取追踪相关信息
func (z *ZapLogger) extractTraceInfo(ctx context.Context, fields []Field) []Field {
	if traceID := ctx.Value("trace_id"); traceID != nil {
		if tid, ok := traceID.(string); ok && tid != "" {
			fields = append(fields, String("trace_id", tid))
		}
	}
	return fields
}
