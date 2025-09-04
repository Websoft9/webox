package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

// UserAuthService user authentication business logic interface
type UserAuthService interface {
	// User registration and authentication
	Register(ctx context.Context, req *request.UserRegisterRequest) (*response.UserResponse, error)
	Login(ctx context.Context, req *request.UserLoginRequest, clientIP string) (*response.UserLoginResponse, error)
	Logout(ctx context.Context, token string) error

	// Password reset
	ForgotPassword(ctx context.Context, req *request.ForgotPasswordRequest) error
	ValidateResetToken(ctx context.Context, token string) error
	ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error

	// Email verification
	VerifyEmail(ctx context.Context, req *request.VerifyEmailRequest) error
	ResendVerificationEmail(ctx context.Context, req *request.ResendVerificationRequest) error

	// OAuth2 login
	OAuth2Login(ctx context.Context, req *request.OAuth2LoginRequest, clientIP string) (*response.UserLoginResponse, error)
}
