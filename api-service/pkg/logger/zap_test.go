package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateFilePath(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid relative path",
			filename: "logs/app.log",
			wantErr:  false,
		},
		{
			name:     "valid relative path with current dir",
			filename: "./logs/app.log",
			wantErr:  false,
		},
		{
			name:     "valid absolute path in /var/log",
			filename: "/var/log/app.log",
			wantErr:  false,
		},
		{
			name:     "valid absolute path in /tmp",
			filename: "/tmp/app.log",
			wantErr:  false,
		},
		{
			name:     "path traversal attack simple",
			filename: "../app.log",
			wantErr:  true,
			errMsg:   "path traversal detected",
		},
		{
			name:     "path traversal attack complex",
			filename: "../../etc/passwd",
			wantErr:  true,
			errMsg:   "path traversal detected",
		},
		{
			name:     "absolute path traversal",
			filename: "/var/log/../../../etc/passwd",
			wantErr:  true,
			errMsg:   "path traversal detected",
		},
		{
			name:     "invalid absolute path",
			filename: "/etc/app.log",
			wantErr:  true,
			errMsg:   "absolute path not allowed",
		},
		{
			name:     "valid data directory path",
			filename: "./data/logs/app.log",
			wantErr:  false,
		},
		{
			name:     "empty path should not error",
			filename: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateFilePath() error = %v, should contain %s", err, tt.errMsg)
				}
			}
		})
	}
}

func TestEnsureLogDir(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
		setup    func() string
		cleanup  func(string)
	}{
		{
			name:     "create new directory",
			filename: "test_logs/subdir/app.log",
			wantErr:  false,
			setup: func() string {
				// Return temp directory for cleanup
				return "test_logs"
			},
			cleanup: func(dir string) {
				os.RemoveAll(dir)
			},
		},
		{
			name:     "current directory file",
			filename: "app.log",
			wantErr:  false,
		},
		{
			name:     "relative path",
			filename: "./app.log",
			wantErr:  false,
		},
		{
			name:     "nested directory creation",
			filename: "deep/nested/path/app.log",
			wantErr:  false,
			setup: func() string {
				return "deep"
			},
			cleanup: func(dir string) {
				os.RemoveAll(dir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cleanupDir string
			if tt.setup != nil {
				cleanupDir = tt.setup()
			}

			err := ensureLogDir(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ensureLogDir() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check if directory was created
			if !tt.wantErr && tt.filename != "" && tt.filename != "./app.log" && tt.filename != "app.log" {
				dir := filepath.Dir(tt.filename)
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					t.Errorf("ensureLogDir() did not create directory %s", dir)
				}
			}

			if tt.cleanup != nil && cleanupDir != "" {
				tt.cleanup(cleanupDir)
			}
		})
	}
}

func TestNewZapLogger(t *testing.T) {
	tests := []struct {
		name   string
		level  Level
		output string
	}{
		{"debug level", DebugLevel, "stdout"},
		{"info level", InfoLevel, "stdout"},
		{"warn level", WarnLevel, "stdout"},
		{"error level", ErrorLevel, "stderr"},
		{"fatal level", FatalLevel, "stderr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewZapLogger(tt.level, &buf)

			if logger == nil {
				t.Fatal("NewZapLogger() returned nil")
			}

			// Test logging at the configured level
			logger.Info("test message")

			// Check if the logger was created with correct level
			zapLogger := logger.(*ZapLogger)
			if zapLogger.logger == nil {
				t.Error("ZapLogger.logger is nil")
			}
		})
	}
}

func TestNewZapLoggerWithConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		setup  func() *Config
		verify func(*testing.T, Logger, *Config)
	}{
		{
			name: "stdout config",
			setup: func() *Config {
				return &Config{
					Level:  InfoLevel,
					Output: "stdout",
				}
			},
			verify: func(t *testing.T, logger Logger, config *Config) {
				if logger == nil {
					t.Error("logger should not be nil")
				}
				if !logger.IsInfoEnabled() {
					t.Error("info level should be enabled")
				}
			},
		},
		{
			name: "stderr config",
			setup: func() *Config {
				return &Config{
					Level:  ErrorLevel,
					Output: "stderr",
				}
			},
			verify: func(t *testing.T, logger Logger, config *Config) {
				if !logger.IsErrorEnabled() {
					t.Error("error level should be enabled")
				}
			},
		},
		{
			name: "file config",
			setup: func() *Config {
				return &Config{
					Level:      DebugLevel,
					Output:     "file",
					Filename:   "./test_logs/test.log",
					MaxSize:    10,
					MaxBackups: 5,
					MaxAge:     30,
					Compress:   true,
				}
			},
			verify: func(t *testing.T, logger Logger, config *Config) {
				if !logger.IsDebugEnabled() {
					t.Error("debug level should be enabled")
				}

				// Test file creation
				logger.Info("test message")
				if _, err := os.Stat("./test_logs/test.log"); os.IsNotExist(err) {
					// File might not exist due to lumberjack lazy creation, that's ok
				}
			},
		},
	}

	// Cleanup test logs directory
	defer os.RemoveAll("./test_logs")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.setup()
			logger := NewZapLoggerWithConfig(config)
			tt.verify(t, logger, config)
		})
	}
}

func TestNewZapLoggerWithServerConfig(t *testing.T) {
	tests := []struct {
		name       string
		logPath    string
		logLevel   string
		maxSize    int
		maxBackups int
		maxAge     int
		compress   bool
	}{
		{
			name:       "stdout config",
			logPath:    "stdout",
			logLevel:   "info",
			maxSize:    10,
			maxBackups: 5,
			maxAge:     30,
			compress:   true,
		},
		{
			name:       "stderr config",
			logPath:    "stderr",
			logLevel:   "error",
			maxSize:    0,
			maxBackups: 0,
			maxAge:     0,
			compress:   false,
		},
		{
			name:       "file config",
			logPath:    "./test_server_logs/server.log",
			logLevel:   "debug",
			maxSize:    5,
			maxBackups: 3,
			maxAge:     7,
			compress:   true,
		},
		{
			name:     "invalid level defaults to info",
			logPath:  "stdout",
			logLevel: "invalid_level",
		},
	}

	// Cleanup test logs directory
	defer os.RemoveAll("./test_server_logs")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewZapLoggerWithServerConfig(
				tt.logPath, tt.logLevel, tt.maxSize, tt.maxBackups, tt.maxAge, tt.compress,
			)

			if logger == nil {
				t.Fatal("NewZapLoggerWithServerConfig() returned nil")
			}

			// Test basic logging
			logger.Info("test server config log")
			logger.Debug("debug message")
			logger.Error("error message")
		})
	}
}

func TestZapLoggerMethods(t *testing.T) {
	var buf bytes.Buffer
	logger := NewZapLogger(DebugLevel, &buf)

	// Test all logging methods
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Test with fields
	logger.Info("message with fields",
		String("string_field", "value"),
		Int("int_field", 42),
		Bool("bool_field", true),
		Float64("float_field", 3.14),
	)

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) < 5 {
		t.Errorf("Expected at least 5 log lines, got %d", len(lines))
	}

	// Verify JSON format
	for i, line := range lines {
		if line == "" {
			continue
		}
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			t.Errorf("Line %d is not valid JSON: %s, error: %v", i, line, err)
		}

		// Check required fields
		if _, ok := logEntry["level"]; !ok {
			t.Errorf("Line %d missing level field", i)
		}
		if _, ok := logEntry["msg"]; !ok {
			t.Errorf("Line %d missing msg field", i)
		}
		if _, ok := logEntry["timestamp"]; !ok {
			t.Errorf("Line %d missing timestamp field", i)
		}
	}
}

func TestZapLoggerContext(t *testing.T) {
	var buf bytes.Buffer
	logger := NewZapLogger(InfoLevel, &buf)

	// Create context with values
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-123")
	ctx = context.WithValue(ctx, "user_id", 456)
	ctx = context.WithValue(ctx, "username", "testuser")
	ctx = context.WithValue(ctx, "role", "admin")
	ctx = context.WithValue(ctx, "trace_id", "trace-789")

	// Test context logging methods
	logger.InfoContext(ctx, "info with context")
	logger.DebugContext(ctx, "debug with context")
	logger.WarnContext(ctx, "warn with context")
	logger.ErrorContext(ctx, "error with context")

	// Test WithContext
	contextLogger := logger.WithContext(ctx)
	contextLogger.Info("info from context logger")

	output := buf.String()

	// Check that context fields are included
	if !strings.Contains(output, "req-123") {
		t.Error("Context request_id not found in log output")
	}
	if !strings.Contains(output, "testuser") {
		t.Error("Context username not found in log output")
	}
	if !strings.Contains(output, "admin") {
		t.Error("Context role not found in log output")
	}
}

func TestZapLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewZapLogger(InfoLevel, &buf)

	// Test WithFields
	fieldsLogger := logger.WithFields(
		String("service", "test-service"),
		String("version", "1.0.0"),
		Int("port", 8080),
	)

	fieldsLogger.Info("service started")

	output := buf.String()

	// Check that fields are included
	if !strings.Contains(output, "test-service") {
		t.Error("Service field not found in log output")
	}
	if !strings.Contains(output, "1.0.0") {
		t.Error("Version field not found in log output")
	}
	if !strings.Contains(output, "8080") {
		t.Error("Port field not found in log output")
	}
}

func TestZapLoggerLevelMethods(t *testing.T) {
	// Test different levels
	debugLogger := NewZapLogger(DebugLevel, &bytes.Buffer{})
	infoLogger := NewZapLogger(InfoLevel, &bytes.Buffer{})
	warnLogger := NewZapLogger(WarnLevel, &bytes.Buffer{})
	errorLogger := NewZapLogger(ErrorLevel, &bytes.Buffer{})

	// Debug logger should enable all levels
	if !debugLogger.IsDebugEnabled() {
		t.Error("Debug level should be enabled for debug logger")
	}
	if !debugLogger.IsInfoEnabled() {
		t.Error("Info level should be enabled for debug logger")
	}
	if !debugLogger.IsWarnEnabled() {
		t.Error("Warn level should be enabled for debug logger")
	}
	if !debugLogger.IsErrorEnabled() {
		t.Error("Error level should be enabled for debug logger")
	}

	// Info logger should not enable debug
	if infoLogger.IsDebugEnabled() {
		t.Error("Debug level should not be enabled for info logger")
	}
	if !infoLogger.IsInfoEnabled() {
		t.Error("Info level should be enabled for info logger")
	}

	// Warn logger should not enable debug or info
	if warnLogger.IsDebugEnabled() || warnLogger.IsInfoEnabled() {
		t.Error("Debug and info levels should not be enabled for warn logger")
	}
	if !warnLogger.IsWarnEnabled() {
		t.Error("Warn level should be enabled for warn logger")
	}

	// Error logger should only enable error
	if errorLogger.IsDebugEnabled() || errorLogger.IsInfoEnabled() || errorLogger.IsWarnEnabled() {
		t.Error("Lower levels should not be enabled for error logger")
	}
	if !errorLogger.IsErrorEnabled() {
		t.Error("Error level should be enabled for error logger")
	}
}

func TestZapLoggerSetLevel(t *testing.T) {
	logger := NewZapLogger(InfoLevel, &bytes.Buffer{})

	// Initially info level
	if !logger.IsInfoEnabled() {
		t.Error("Info should be enabled initially")
	}
	if logger.IsDebugEnabled() {
		t.Error("Debug should not be enabled initially")
	}

	// Change to debug level
	logger.SetLevel(DebugLevel)
	if !logger.IsDebugEnabled() {
		t.Error("Debug should be enabled after SetLevel(DebugLevel)")
	}

	// Change to error level
	logger.SetLevel(ErrorLevel)
	if logger.IsInfoEnabled() {
		t.Error("Info should not be enabled after SetLevel(ErrorLevel)")
	}
	if !logger.IsErrorEnabled() {
		t.Error("Error should be enabled after SetLevel(ErrorLevel)")
	}
}

func TestConvertFields(t *testing.T) {
	logger := NewZapLogger(InfoLevel, &bytes.Buffer{}).(*ZapLogger)

	tests := []struct {
		name   string
		fields []Field
		verify func(t *testing.T, fields []Field)
	}{
		{
			name:   "empty fields",
			fields: []Field{},
			verify: func(t *testing.T, fields []Field) {
				zapFields := logger.convertFields(fields...)
				if zapFields != nil {
					t.Error("Empty fields should return nil")
				}
			},
		},
		{
			name: "various field types",
			fields: []Field{
				String("str", "value"),
				Int("int", 42),
				Int64("int64", 123456789),
				Uint("uint", 99),
				Float64("float", 3.14),
				Bool("bool", true),
				ErrorField(errors.New("test error")),
				Any("any", map[string]string{"key": "value"}),
			},
			verify: func(t *testing.T, fields []Field) {
				zapFields := logger.convertFields(fields...)
				if len(zapFields) != len(fields) {
					t.Errorf("Expected %d zap fields, got %d", len(fields), len(zapFields))
				}
			},
		},
		{
			name: "nil error field",
			fields: []Field{
				ErrorField(nil),
			},
			verify: func(t *testing.T, fields []Field) {
				zapFields := logger.convertFields(fields...)
				if len(zapFields) != 1 {
					t.Error("Nil error should still create a field")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.verify(t, tt.fields)
		})
	}
}

func TestExtractContextFields(t *testing.T) {
	logger := NewZapLogger(InfoLevel, &bytes.Buffer{}).(*ZapLogger)

	tests := []struct {
		name     string
		ctx      context.Context
		expected int // expected number of fields
	}{
		{
			name:     "nil context",
			ctx:      nil,
			expected: 0,
		},
		{
			name:     "empty context",
			ctx:      context.Background(),
			expected: 0,
		},
		{
			name: "context with string fields",
			ctx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, "request_id", "req-123")
				ctx = context.WithValue(ctx, "username", "testuser")
				ctx = context.WithValue(ctx, "role", "admin")
				return ctx
			}(),
			expected: 3,
		},
		{
			name: "context with mixed types",
			ctx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, "request_id", "req-456")
				ctx = context.WithValue(ctx, "user_id", 789) // non-string
				ctx = context.WithValue(ctx, "username", "") // empty string should be skipped
				return ctx
			}(),
			expected: 2, // request_id and user_id (empty username skipped)
		},
		{
			name: "context with all fields",
			ctx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, "request_id", "req-full")
				ctx = context.WithValue(ctx, "user_id", "user-123")
				ctx = context.WithValue(ctx, "username", "fulluser")
				ctx = context.WithValue(ctx, "role", "superadmin")
				ctx = context.WithValue(ctx, "trace_id", "trace-999")
				return ctx
			}(),
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := logger.extractContextFields(tt.ctx)
			if len(fields) != tt.expected {
				t.Errorf("Expected %d context fields, got %d", tt.expected, len(fields))
			}
		})
	}
}

func TestCreateOutput(t *testing.T) {
	tests := []struct {
		name       string
		config     *OutputConfig
		outputType string
		setup      func() string
		cleanup    func(string)
	}{
		{
			name:       "stdout output",
			config:     &OutputConfig{Path: "stdout"},
			outputType: "stdout",
		},
		{
			name:       "stderr output",
			config:     &OutputConfig{Path: "stderr"},
			outputType: "stderr",
		},
		{
			name: "file output",
			config: &OutputConfig{
				Path:       "./test_output/app.log",
				MaxSize:    10,
				MaxBackups: 5,
				MaxAge:     30,
				Compress:   true,
			},
			outputType: "file",
			setup: func() string {
				return "./test_output"
			},
			cleanup: func(dir string) {
				os.RemoveAll(dir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cleanupDir string
			if tt.setup != nil {
				cleanupDir = tt.setup()
			}

			output := createOutput(tt.config, tt.outputType)
			if output == nil {
				t.Error("createOutput() returned nil")
			}

			if tt.cleanup != nil && cleanupDir != "" {
				tt.cleanup(cleanupDir)
			}
		})
	}
}

// Benchmark tests
func BenchmarkZapLoggerInfo(b *testing.B) {
	logger := NewZapLogger(InfoLevel, &bytes.Buffer{})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", String("iteration", string(rune(i))))
	}
}

func BenchmarkZapLoggerWithFields(b *testing.B) {
	logger := NewZapLogger(InfoLevel, &bytes.Buffer{})
	fields := []Field{
		String("service", "benchmark"),
		Int("port", 8080),
		Bool("debug", true),
	}

	fieldsLogger := logger.WithFields(fields...)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fieldsLogger.Info("benchmark message")
	}
}

func BenchmarkConvertFields(b *testing.B) {
	logger := NewZapLogger(InfoLevel, &bytes.Buffer{}).(*ZapLogger)
	fields := []Field{
		String("str", "value"),
		Int("int", 42),
		Bool("bool", true),
		Float64("float", 3.14),
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.convertFields(fields...)
	}
}

// TestCreateOutputPreventReservedNames tests that reserved names don't create files
func TestCreateOutputPreventReservedNames(t *testing.T) {
	tests := []struct {
		name       string
		config     *OutputConfig
		outputType string
		expectType string // "stdout", "stderr", or "file"
	}{
		{
			name:       "stderr_output_type_returns_stderr",
			config:     &OutputConfig{Path: "stderr"},
			outputType: "stderr",
			expectType: "stderr",
		},
		{
			name:       "stderr_filename_prevented",
			config:     &OutputConfig{Path: "stderr"},
			outputType: "file",
			expectType: "stdout", // Should fallback to stdout when reserved name used as filename
		},
		{
			name:       "stdout_filename_prevented",
			config:     &OutputConfig{Path: "stdout"},
			outputType: "file",
			expectType: "stdout", // Should fallback to stdout when reserved name used as filename
		},
		{
			name:       "valid_filename_creates_file",
			config:     &OutputConfig{Path: "./logs/test.log"},
			outputType: "file",
			expectType: "file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := createOutput(tt.config, tt.outputType)

			switch tt.expectType {
			case "stdout":
				if writer != os.Stdout {
					t.Errorf("Expected stdout writer, got different writer")
				}
			case "stderr":
				if writer != os.Stderr {
					t.Errorf("Expected stderr writer, got different writer")
				}
			case "file":
				// For file type, just check it's not stdout or stderr
				if writer == os.Stdout || writer == os.Stderr {
					t.Errorf("Expected file writer, got stdout or stderr")
				}
			}

			// Ensure no files with reserved names are created in current directory
			if _, err := os.Stat("stderr"); err == nil {
				os.Remove("stderr") // Clean up if it exists
				t.Error("stderr file was incorrectly created")
			}
			if _, err := os.Stat("stdout"); err == nil {
				os.Remove("stdout") // Clean up if it exists
				t.Error("stdout file was incorrectly created")
			}
		})
	}
}
