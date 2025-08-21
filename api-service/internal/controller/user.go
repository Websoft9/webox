package controller

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/service"
	"api-service/internal/middleware"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	pkg_response "api-service/pkg/response"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
	logger      logger.Logger
}

// NewUserController 创建新的用户控制器
func NewUserController(userService service.UserService, logger logger.Logger) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

// bindAndValidateRequest 绑定并验证请求参数的通用方法
func (c *UserController) bindAndValidateRequest(ctx *gin.Context, req interface{}, action string) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		c.logger.WarnContext(ctx, action+"请求参数绑定失败", logger.ErrorField(err))
		errorMsg := middleware.T(ctx, "common.validation_failed")
		errors.HandleError(ctx, errors.NewAppErrorWithI18n(errors.CodeValidationError, errorMsg, "common.validation_failed"))
		return false
	}
	return true
}

// handleUserAuth 处理用户认证相关请求的通用方法
func (c *UserController) handleUserAuth(
	ctx *gin.Context,
	req interface{},
	action string,
	serviceFunc func(context.Context, interface{}) (interface{}, error),
	successMessage string,
) {
	if !c.bindAndValidateRequest(ctx, req, action) {
		return
	}

	result, err := serviceFunc(ctx, req)
	if err != nil {
		// 根据action类型选择日志级别
		if action == "登录" {
			c.logger.WarnContext(ctx, "用户"+action+"失败", logger.ErrorField(err))
		} else {
			c.logger.ErrorContext(ctx, "用户"+action+"失败", logger.ErrorField(err))
		}
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "用户"+action+"成功")
	pkg_response.Success(ctx, successMessage, result)
}

// Register 用户注册
func (c *UserController) Register(ctx *gin.Context) {
	var req request.UserRegisterRequest
	c.handleUserAuth(ctx, &req, "注册", func(ctx context.Context, r interface{}) (interface{}, error) {
		return c.userService.Register(ctx, r.(*request.UserRegisterRequest))
	}, middleware.T(ctx, "user.created_success"))
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var req request.UserLoginRequest
	c.handleUserAuth(ctx, &req, "登录", func(ctx context.Context, r interface{}) (interface{}, error) {
		return c.userService.Login(ctx, r.(*request.UserLoginRequest))
	}, middleware.T(ctx, "user.login_success"))
}

// GetProfile 获取用户资料
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID := c.getCurrentUserID(ctx)
	if userID == 0 {
		errors.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	profile, err := c.userService.GetProfile(ctx, userID)
	if err != nil {
		c.logger.ErrorContext(ctx, "获取用户资料失败", logger.Uint("user_id", userID), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	successMsg := middleware.T(ctx, "common.success")
	pkg_response.Success(ctx, successMsg, profile)
}

// UpdateProfile 更新用户资料
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID := c.getCurrentUserID(ctx)
	if userID == 0 {
		errors.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	var req request.UserUpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WarnContext(ctx, "更新资料请求参数绑定失败", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, "请求参数无效"))
		return
	}

	user, err := c.userService.UpdateProfile(ctx, userID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "更新用户资料失败", logger.Uint("user_id", userID), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "用户资料更新成功", logger.Uint("user_id", userID))
	pkg_response.Success(ctx, "用户资料更新成功", user)
}

// ChangePassword 修改密码
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID := c.getCurrentUserID(ctx)
	if userID == 0 {
		errors.HandleError(ctx, errors.ErrUnauthorized)
		return
	}

	var req request.UserChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WarnContext(ctx, "修改密码请求参数绑定失败", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, "请求参数无效"))
		return
	}

	err := c.userService.ChangePassword(ctx, userID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "修改密码失败", logger.Uint("user_id", userID), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "密码修改成功", logger.Uint("user_id", userID))
	pkg_response.Success(ctx, "密码修改成功", nil)
}

// ListUsers 获取用户列表（管理员功能）
func (c *UserController) ListUsers(ctx *gin.Context) {
	// 检查管理员权限
	if err := c.checkAdminPermission(ctx); err != nil {
		errors.HandleError(ctx, err)
		return
	}

	var req request.UserListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.WarnContext(ctx, "用户列表请求参数绑定失败", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, "请求参数无效"))
		return
	}

	users, total, err := c.userService.ListUsers(ctx, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "获取用户列表失败", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	// 构建分页响应
	paginationResp := &struct {
		*request.UserListRequest
		Total int64                      `json:"total"`
		Users *response.UserListResponse `json:"users"`
	}{
		UserListRequest: &req,
		Total:           total,
		Users:           users,
	}

	pkg_response.Success(ctx, "获取用户列表成功", paginationResp)
}

// GetUser 获取单个用户信息（管理员功能）
func (c *UserController) GetUser(ctx *gin.Context) {
	// 检查管理员权限
	if err := c.checkAdminPermission(ctx); err != nil {
		errors.HandleError(ctx, err)
		return
	}

	userID, err := c.getIDFromPath(ctx, "id")
	if err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, "无效的用户ID"))
		return
	}

	user, err := c.userService.GetUser(ctx, userID)
	if err != nil {
		c.logger.ErrorContext(ctx, "获取用户信息失败", logger.Uint("target_user_id", userID), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	pkg_response.Success(ctx, "获取用户信息成功", user)
}

// UpdateUserStatus 更新用户状态（管理员功能）
func (c *UserController) UpdateUserStatus(ctx *gin.Context) {
	// 检查管理员权限
	if err := c.checkAdminPermission(ctx); err != nil {
		errors.HandleError(ctx, err)
		return
	}

	userID, err := c.getIDFromPath(ctx, "id")
	if err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, "无效的用户ID"))
		return
	}

	var req request.UserUpdateStatusRequest
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		c.logger.WarnContext(ctx, "更新用户状态请求参数绑定失败", logger.ErrorField(bindErr))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, "请求参数无效"))
		return
	}

	updateErr := c.userService.UpdateUserStatus(ctx, userID, &req)
	if updateErr != nil {
		c.logger.ErrorContext(ctx, "更新用户状态失败", logger.Uint("target_user_id", userID), logger.ErrorField(updateErr))
		errors.HandleError(ctx, updateErr)
		return
	}

	c.logger.InfoContext(ctx, "用户状态更新成功", logger.Uint("target_user_id", userID), logger.String("new_status", req.Status))
	pkg_response.Success(ctx, "用户状态更新成功", nil)
}

// DeleteUser 删除用户（管理员功能）
func (c *UserController) DeleteUser(ctx *gin.Context) {
	// 检查管理员权限
	if err := c.checkAdminPermission(ctx); err != nil {
		errors.HandleError(ctx, err)
		return
	}

	userID, err := c.getIDFromPath(ctx, "id")
	if err != nil {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, "无效的用户ID"))
		return
	}

	// 不能删除自己
	currentUserID := c.getCurrentUserID(ctx)
	if currentUserID == userID {
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationError, "不能删除自己的账号"))
		return
	}

	err = c.userService.DeleteUser(ctx, userID)
	if err != nil {
		c.logger.ErrorContext(ctx, "删除用户失败", logger.Uint("target_user_id", userID), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "用户删除成功", logger.Uint("target_user_id", userID))
	pkg_response.Success(ctx, "用户删除成功", nil)
}

// getCurrentUserID 从上下文中获取当前用户ID
func (c *UserController) getCurrentUserID(ctx *gin.Context) uint {
	if userID, exists := ctx.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}

// checkAdminPermission 检查管理员权限
func (c *UserController) checkAdminPermission(ctx *gin.Context) error {
	userID := c.getCurrentUserID(ctx)
	if userID == 0 {
		return errors.ErrUnauthorized
	}

	// 检查用户访问权限
	return c.userService.ValidateUserAccess(ctx, userID, "user:manage")
}

// getIDFromPath 从路径参数中获取ID
func (c *UserController) getIDFromPath(ctx *gin.Context, paramName string) (uint, error) {
	idStr := ctx.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
