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

// Constant definitions
const (
	DefaultFilePermission = 0o666 // Default file permission
	ContextFieldCapacity  = 4     // Context field capacity
)

// Allowed log file directory prefixes (safe paths)
var allowedLogPaths = []string{
	"/var/log/",
	"/tmp/",
	"./logs/",
	"./data/",
}

// OutputConfig output configuration
type OutputConfig struct {
	Path       string
	MaxSize    int // MB
	MaxBackups int
	MaxAge     int // days
	Compress   bool
}

// ZapLogger Zap log implementation
type ZapLogger struct {
	logger      *zap.Logger
	atomicLevel zap.AtomicLevel // Support dynamic level adjustment
}

// NewZapLogger create new Zap logger instance
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

// NewDefaultZapLogger create default Zap logger instance
func NewDefaultZapLogger() Logger {
	return NewZapLogger(InfoLevel, os.Stdout)
}

// NewZapLoggerWithConfig create Zap logger instance from configuration
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

// NewZapLoggerWithServerConfig create Zap logger instance from server configuration
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

// createEncoderConfig create encoder configuration
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

// createOutput create output Writer with file rotation support
func createOutput(config *OutputConfig, outputType string) io.Writer {
	switch outputType {
	case OutputTypeStderr:
		return os.Stderr
	case OutputTypeStdout, "":
		// If explicitly specified as stdout or config path is empty/stdout, return stdout
		if config.Path == "" || config.Path == OutputTypeStdout {
			return os.Stdout
		}
		// If specific file path is configured, create file output
		return createFileOutput(config)
	case OutputTypeFile:
		return createFileOutput(config)
	default:
		return os.Stdout
	}
}

// createFileOutput create file output with rotation support
func createFileOutput(config *OutputConfig) io.Writer {
	if config.Path == "" {
		return os.Stdout
	}

	// Prevent using reserved output names as filenames
	baseName := filepath.Base(config.Path)
	if baseName == OutputTypeStdout || baseName == OutputTypeStderr {
		return os.Stdout
	}

	// Validate file path security
	if err := validateFilePath(config.Path); err != nil {
		return os.Stdout
	}

	// Clean and normalize path to prevent path traversal attacks
	cleanPath := filepath.Clean(config.Path)

	// Ensure log directory exists
	if err := ensureLogDir(cleanPath); err != nil {
		return os.Stdout
	}

	// If rotation parameters are configured, use lumberjack
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

	// Regular file output - path verified by validateFilePath and cleaned by filepath.Clean
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	// #nosec G304 - 路径已经过 validateFilePath 验证并使用 filepath.Clean 清理
	if file, err := os.OpenFile(cleanPath, flags, DefaultFilePermission); err == nil {
		return file
	}

	return os.Stdout
}

// ensureLogDir ensure log directory exists
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

// validateFilePath validate file path security
func validateFilePath(filename string) error {
	if filename == "" {
		return nil
	}

	// Check path traversal characters in original path (before cleaning)
	if strings.Contains(filename, "..") {
		return fmt.Errorf("path traversal detected: %s", filename)
	}

	// Clean path
	cleanPath := filepath.Clean(filename)

	// Check absolute path permissions
	if filepath.IsAbs(cleanPath) {
		for _, allowedPath := range allowedLogPaths {
			if strings.HasPrefix(cleanPath, allowedPath) {
				return nil
			}
		}
		return fmt.Errorf("absolute path not allowed: %s", filename)
	}

	// Check relative path security
	if strings.HasPrefix(cleanPath, "/") || strings.HasPrefix(cleanPath, "\\") {
		return fmt.Errorf("invalid relative path: %s", filename)
	}

	return nil
}

// convertToZapLevel convert to Zap log level
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

// Debug debug level logging
func (z *ZapLogger) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, z.convertFields(fields...)...)
}

// Info info level logging
func (z *ZapLogger) Info(msg string, fields ...Field) {
	z.logger.Info(msg, z.convertFields(fields...)...)
}

// Warn warning level logging
func (z *ZapLogger) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, z.convertFields(fields...)...)
}

// Error error level logging
func (z *ZapLogger) Error(msg string, fields ...Field) {
	z.logger.Error(msg, z.convertFields(fields...)...)
}

// Fatal fatal level logging
func (z *ZapLogger) Fatal(msg string, fields ...Field) {
	z.logger.Fatal(msg, z.convertFields(fields...)...)
}

// DebugContext debug logging with context
func (z *ZapLogger) DebugContext(ctx context.Context, msg string, fields ...Field) {
	z.Debug(msg, append(fields, z.extractContextFields(ctx)...)...)
}

// InfoContext info logging with context
func (z *ZapLogger) InfoContext(ctx context.Context, msg string, fields ...Field) {
	z.Info(msg, append(fields, z.extractContextFields(ctx)...)...)
}

// WarnContext warning logging with context
func (z *ZapLogger) WarnContext(ctx context.Context, msg string, fields ...Field) {
	z.Warn(msg, append(fields, z.extractContextFields(ctx)...)...)
}

// ErrorContext error logging with context
func (z *ZapLogger) ErrorContext(ctx context.Context, msg string, fields ...Field) {
	z.Error(msg, append(fields, z.extractContextFields(ctx)...)...)
}

// IsDebugEnabled check if debug logging is enabled
func (z *ZapLogger) IsDebugEnabled() bool {
	return z.logger.Core().Enabled(zapcore.DebugLevel)
}

// IsInfoEnabled check if info logging is enabled
func (z *ZapLogger) IsInfoEnabled() bool {
	return z.logger.Core().Enabled(zapcore.InfoLevel)
}

// IsWarnEnabled check if warning logging is enabled
func (z *ZapLogger) IsWarnEnabled() bool {
	return z.logger.Core().Enabled(zapcore.WarnLevel)
}

// IsErrorEnabled check if error logging is enabled
func (z *ZapLogger) IsErrorEnabled() bool {
	return z.logger.Core().Enabled(zapcore.ErrorLevel)
}

// WithFields add fields
func (z *ZapLogger) WithFields(fields ...Field) Logger {
	if len(fields) == 0 {
		return z
	}

	return &ZapLogger{
		logger:      z.logger.With(z.convertFields(fields...)...),
		atomicLevel: z.atomicLevel,
	}
}

// WithContext add context
func (z *ZapLogger) WithContext(ctx context.Context) Logger {
	return z.WithFields(z.extractContextFields(ctx)...)
}

// SetLevel set logging level
func (z *ZapLogger) SetLevel(level Level) {
	z.atomicLevel.SetLevel(convertToZapLevel(level))
}

// SetOutput set output (Zap does not support dynamic output modification)
func (z *ZapLogger) SetOutput(w io.Writer) {
	// Zap output is set at creation time and cannot be dynamically modified
}

// convertFields convert field format
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

// convertField convert single field
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

// extractContextFields extract fields from context
func (z *ZapLogger) extractContextFields(ctx context.Context) []Field {
	if ctx == nil {
		return nil
	}

	fields := make([]Field, 0, ContextFieldCapacity)

	// Define field names to extract and their corresponding keys
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
				// user_id might not be string type
				fields = append(fields, Any(fieldKey, value))
			}
		}
	}

	return fields
}
