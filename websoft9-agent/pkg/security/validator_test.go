package security

import (
	"testing"
)

func TestCommandValidator_ValidateCommand(t *testing.T) {
	validator := NewCommandValidator()

	tests := []struct {
		name    string
		command string
		wantErr bool
	}{
		{
			name:    "valid systemctl command",
			command: "systemctl status nginx",
			wantErr: false,
		},
		{
			name:    "valid docker command",
			command: "docker ps",
			wantErr: false,
		},
		{
			name:    "command injection attempt with semicolon",
			command: "ls; rm -rf /",
			wantErr: true,
		},
		{
			name:    "command injection attempt with pipe",
			command: "cat /etc/passwd | grep root",
			wantErr: true,
		},
		{
			name:    "command injection attempt with backticks",
			command: "echo `whoami`",
			wantErr: true,
		},
		{
			name:    "unauthorized command",
			command: "rm -rf /",
			wantErr: true,
		},
		{
			name:    "empty command",
			command: "",
			wantErr: true,
		},
		{
			name:    "command with dangerous characters",
			command: "ls && echo 'hacked'",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandValidator_ValidateSystemctlAction(t *testing.T) {
	validator := NewCommandValidator()

	tests := []struct {
		name    string
		action  string
		wantErr bool
	}{
		{
			name:    "valid start action",
			action:  "start",
			wantErr: false,
		},
		{
			name:    "valid stop action",
			action:  "stop",
			wantErr: false,
		},
		{
			name:    "valid restart action",
			action:  "restart",
			wantErr: false,
		},
		{
			name:    "valid status action",
			action:  "status",
			wantErr: false,
		},
		{
			name:    "invalid action",
			action:  "destroy",
			wantErr: true,
		},
		{
			name:    "empty action",
			action:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateSystemctlAction(tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSystemctlAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandValidator_ValidateServiceName(t *testing.T) {
	validator := NewCommandValidator()

	tests := []struct {
		name        string
		serviceName string
		wantErr     bool
	}{
		{
			name:        "valid service name",
			serviceName: "nginx",
			wantErr:     false,
		},
		{
			name:        "valid service name with hyphen",
			serviceName: "my-service",
			wantErr:     false,
		},
		{
			name:        "valid service name with underscore",
			serviceName: "my_service",
			wantErr:     false,
		},
		{
			name:        "invalid service name with special chars",
			serviceName: "nginx; rm -rf /",
			wantErr:     true,
		},
		{
			name:        "invalid service name with spaces",
			serviceName: "my service",
			wantErr:     true,
		},
		{
			name:        "empty service name",
			serviceName: "",
			wantErr:     true,
		},
		{
			name:        "too long service name",
			serviceName: "this-is-a-very-long-service-name-that-exceeds-the-maximum-allowed-length-limit",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateServiceName(tt.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateServiceName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPathValidator_ValidateConfigPath(t *testing.T) {
	validator := NewPathValidator()

	tests := []struct {
		name       string
		configPath string
		wantErr    bool
	}{
		{
			name:       "valid config path",
			configPath: "/etc/websoft9/config.yaml",
			wantErr:    false,
		},
		{
			name:       "valid config path with json",
			configPath: "/var/lib/websoft9/config.json",
			wantErr:    false,
		},
		{
			name:       "path traversal attempt",
			configPath: "/etc/websoft9/../../../etc/passwd",
			wantErr:    true,
		},
		{
			name:       "unauthorized directory",
			configPath: "/etc/shadow",
			wantErr:    true,
		},
		{
			name:       "invalid file extension",
			configPath: "/etc/websoft9/config.exe",
			wantErr:    true,
		},
		{
			name:       "empty path",
			configPath: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateConfigPath(tt.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal input",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "input with control characters",
			input:    "hello\x00\x1f world",
			expected: "hello world",
		},
		{
			name:     "input with leading/trailing spaces",
			input:    "  hello world  ",
			expected: "hello world",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "input with only spaces",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeInput() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateTaskParams(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid params",
			params: map[string]interface{}{
				"service": "nginx",
				"action":  "start",
			},
			wantErr: false,
		},
		{
			name: "invalid key name",
			params: map[string]interface{}{
				"service@name": "nginx",
			},
			wantErr: true,
		},
		{
			name: "value too long",
			params: map[string]interface{}{
				"data": string(make([]byte, 5000)),
			},
			wantErr: true,
		},
		{
			name: "nested params",
			params: map[string]interface{}{
				"config": map[string]interface{}{
					"host": "localhost",
					"port": 8080,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid nested key",
			params: map[string]interface{}{
				"config": map[string]interface{}{
					"host@name": "localhost",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTaskParams(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTaskParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
