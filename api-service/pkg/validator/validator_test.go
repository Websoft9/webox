package validator

import (
	"api-service/pkg/errors"
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
		errCode  int
	}{
		{
			name:     "valid username",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "valid username with numbers",
			username: "user123",
			wantErr:  false,
		},
		{
			name:     "valid username with underscore",
			username: "test_user",
			wantErr:  false,
		},
		{
			name:     "empty username",
			username: "",
			wantErr:  true,
			errCode:  errors.CodeInvalidParameterFormat,
		},
		{
			name:     "username too short",
			username: "ab",
			wantErr:  true,
			errCode:  errors.CodeInvalidParameterFormat,
		},
		{
			name:     "username too long",
			username: "verylongusernamethatexceedslimit",
			wantErr:  true,
			errCode:  errors.CodeInvalidParameterFormat,
		},
		{
			name:     "username with invalid characters",
			username: "user@name",
			wantErr:  true,
			errCode:  errors.CodeInvalidParameterFormat,
		},
		{
			name:     "username starting with number",
			username: "123user",
			wantErr:  true,
			errCode:  errors.CodeInvalidParameterFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errCode != 0 {
				if appErr, ok := err.(*errors.AppError); ok {
					if appErr.Code != tt.errCode {
						t.Errorf("ValidateUsername() error code = %v, want %v", appErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("ValidateUsername() error is not AppError type")
				}
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
		errCode  int
	}{
		{
			name:     "valid strong password",
			password: "StrongP@ss123",
			wantErr:  false,
		},
		{
			name:     "valid password with minimum requirements",
			password: "Pass123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
			errCode:  errors.CodeRequiredParameterMissing,
		},
		{
			name:     "password too short",
			password: "12345",
			wantErr:  true,
			errCode:  errors.CodeRequiredParameterMissing,
		},
		{
			name:     "password too long",
			password: "verylongpasswordthatexceedsfiftycharlimitandshouldfail",
			wantErr:  true,
			errCode:  errors.CodeRequiredParameterMissing,
		},
		{
			name:     "weak password - only lowercase",
			password: "password",
			wantErr:  true,
			errCode:  errors.CodePasswordTooWeak,
		},
		{
			name:     "weak password - only numbers",
			password: "123456",
			wantErr:  true,
			errCode:  errors.CodePasswordTooWeak,
		},
		{
			name:     "password with uppercase and lowercase",
			password: "Password",
			wantErr:  false,
		},
		{
			name:     "password with lowercase and numbers",
			password: "password123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errCode != 0 {
				if appErr, ok := err.(*errors.AppError); ok {
					if appErr.Code != tt.errCode {
						t.Errorf("ValidatePassword() error code = %v, want %v", appErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("ValidatePassword() error is not AppError type")
				}
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
		errCode int
	}{
		{
			name:    "valid email",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "valid email with plus",
			email:   "user+tag@example.com",
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
			errCode: errors.CodeInvalidEmailFormat,
		},
		{
			name:    "invalid email format - no @",
			email:   "userexample.com",
			wantErr: true,
			errCode: errors.CodeInvalidEmailFormat,
		},
		{
			name:    "invalid email format - no domain",
			email:   "user@",
			wantErr: true,
			errCode: errors.CodeInvalidEmailFormat,
		},
		{
			name:    "invalid email format - no TLD",
			email:   "user@example",
			wantErr: true,
			errCode: errors.CodeInvalidEmailFormat,
		},
		{
			name:    "invalid email format - multiple @",
			email:   "user@@example.com",
			wantErr: true,
			errCode: errors.CodeInvalidEmailFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errCode != 0 {
				if appErr, ok := err.(*errors.AppError); ok {
					if appErr.Code != tt.errCode {
						t.Errorf("ValidateEmail() error code = %v, want %v", appErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("ValidateEmail() error is not AppError type")
				}
			}
		})
	}
}

// 基准测试
func BenchmarkValidateUsername(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateUsername("testuser123")
	}
}

func BenchmarkValidatePassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidatePassword("StrongPassword123")
	}
}

func BenchmarkValidateEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateEmail("user@example.com")
	}
}
