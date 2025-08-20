package logger

import (
	"context"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger Zap日志实现
type ZapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
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

	// 设置日志级别
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

	// 创建核心
	if output == nil {
		output = os.Stdout
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(output),
		zapLevel,
	)

	// 创建logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

// NewDefaultZapLogger 创建默认Zap日志实例
func NewDefaultZapLogger() Logger {
	return NewZapLogger(InfoLevel, os.Stdout)
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

// WithFields 添加字段
func (z *ZapLogger) WithFields(fields ...Field) Logger {
	return &ZapLogger{
		logger: z.logger.With(z.convertFields(fields...)...),
		sugar:  z.logger.Sugar(),
	}
}

// WithContext 添加上下文
func (z *ZapLogger) WithContext(ctx context.Context) Logger {
	return z.WithFields(z.extractContextFields(ctx)...)
}

// SetLevel 设置日志级别
func (z *ZapLogger) SetLevel(level Level) {
	// Zap的级别在创建时设置，这里可以记录但不能动态修改
	// 实际项目中可能需要使用zap.AtomicLevel来支持动态修改
}

// SetOutput 设置输出
func (z *ZapLogger) SetOutput(w io.Writer) {
	// Zap的输出在创建时设置，这里可以记录但不能动态修改
	// 实际项目中可能需要重新创建logger
}

// convertFields 转换字段格式
func (z *ZapLogger) convertFields(fields ...Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

// extractContextFields 从上下文中提取字段
func (z *ZapLogger) extractContextFields(ctx context.Context) []Field {
	var fields []Field

	// 可以从上下文中提取请求ID、用户ID等信息
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields = append(fields, String("request_id", requestID.(string)))
	}

	if userID := ctx.Value("user_id"); userID != nil {
		fields = append(fields, Any("user_id", userID))
	}

	return fields
}
