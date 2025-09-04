package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

// UserAuthService 用户认证业务逻辑接口
type UserAuthService interface {
	// 用户注册与认证
	Register(ctx context.Context, req *request.UserRegisterRequest) (*response.UserResponse, error)
	Login(ctx context.Context, req *request.UserLoginRequest, clientIP string) (*response.UserLoginResponse, error)
	Logout(ctx context.Context, token string) error

	// 密码重置
	ForgotPassword(ctx context.Context, req *request.ForgotPasswordRequest) error
	ValidateResetToken(ctx context.Context, token string) error
	ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error

	// 邮箱验证
	VerifyEmail(ctx context.Context, req *request.VerifyEmailRequest) error
	ResendVerificationEmail(ctx context.Context, req *request.ResendVerificationRequest) error

	// OAuth2登录
	OAuth2Login(ctx context.Context, req *request.OAuth2LoginRequest, clientIP string) (*response.UserLoginResponse, error)
}
