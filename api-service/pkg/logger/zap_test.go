package logger

import (
	"testing"
)

func TestValidateFilePath(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid relative path",
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
			name:     "path traversal attack",
			filename: "../../../etc/passwd",
			wantErr:  true,
		},
		{
			name:     "path traversal in absolute path",
			filename: "/var/log/../../../etc/passwd",
			wantErr:  true,
		},
		{
			name:     "invalid absolute path",
			filename: "/etc/app.log",
			wantErr:  true,
		},
		{
			name:     "invalid relative path with leading slash",
			filename: "/logs/app.log",
			wantErr:  true,
		},
		{
			name:     "valid data directory path",
			filename: "./data/logs/app.log",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFilePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewZapLoggerWithConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "stdout output",
			config: &Config{
				Level:  InfoLevel,
				Output: "stdout",
			},
		},
		{
			name: "stderr output",
			config: &Config{
				Level:  ErrorLevel,
				Output: "stderr",
			},
		},
		{
			name: "safe file output",
			config: &Config{
				Level:    DebugLevel,
				Output:   "file",
				Filename: "./logs/test.log",
			},
		},
		{
			name: "unsafe file output - should fallback to stdout",
			config: &Config{
				Level:    InfoLevel,
				Output:   "file",
				Filename: "../../../etc/passwd",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewZapLoggerWithConfig(tt.config)
			if logger == nil {
				t.Error("NewZapLoggerWithConfig() returned nil")
			}

			// Test basic logging functionality
			logger.Info("test message")
			logger.Debug("debug message")
			logger.Warn("warning message")
		})
	}
}

func TestZapLoggerLevels(t *testing.T) {
	logger := NewDefaultZapLogger()

	// Test level checking methods
	if !logger.IsInfoEnabled() {
		t.Error("Info level should be enabled by default")
	}

	// Test logging methods
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warning message")
	logger.Error("error message")

	// Test with fields
	logger.Info("message with fields",
		String("key1", "value1"),
		Int("key2", 42),
		Bool("key3", true),
	)
}

func TestZapLoggerWithFields(t *testing.T) {
	logger := NewDefaultZapLogger()

	// Test WithFields
	loggerWithFields := logger.WithFields(
		String("service", "test"),
		String("version", "1.0.0"),
	)

	loggerWithFields.Info("test message")

	// Test field conversion
	zapLogger := logger.(*ZapLogger)
	fields := []Field{
		String("string_field", "test"),
		Int("int_field", 123),
		Bool("bool_field", true),
		Float64("float_field", 3.14),
		ErrorField(nil), // nil error should be handled gracefully
	}

	zapFields := zapLogger.convertFields(fields...)
	if len(zapFields) != 5 { // All fields should be converted, including nil error
		t.Errorf("Expected 5 fields, got %d", len(zapFields))
	}
}
