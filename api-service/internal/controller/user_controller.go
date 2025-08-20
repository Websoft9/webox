package controller

import (
	"api-service/internal/dto"
	"api-service/internal/service"
	"api-service/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := c.userService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Registration failed", err.Error())
		return
	}

	response.Success(ctx, "User registered successfully", user)
}

func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	token, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	response.Success(ctx, "Login successful", gin.H{"token": token})
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized", "User ID not found")
		return
	}

	user, err := c.userService.GetProfile(userID.(uint))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "User not found", err.Error())
		return
	}

	response.Success(ctx, "Profile retrieved successfully", user)
}

// ListUsers 获取用户列表
func (c *UserController) ListUsers(ctx *gin.Context) {
	var query dto.UserListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	// 参数验证
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 20
	}
	if query.Page <= 0 {
		query.Page = 1
	}

	users, pagination, err := c.userService.ListUsers(&query)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get users", err.Error())
		return
	}

	response.Success(ctx, "Users retrieved successfully", gin.H{
		"items":      users,
		"pagination": pagination,
	})
}

// GetUser 获取用户详情
func (c *UserController) GetUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	user, err := c.userService.GetUserByID(uint(userID))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "User not found", err.Error())
		return
	}

	response.Success(ctx, "User retrieved successfully", user)
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// 验证密码复杂度
	if !c.userService.ValidatePassword(req.Password) {
		response.Error(ctx, http.StatusBadRequest, "Password validation failed",
			"Password must contain uppercase, lowercase letters and numbers")
		return
	}

	user, err := c.userService.CreateUser(&req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Failed to create user", err.Error())
		return
	}

	response.Success(ctx, "User created successfully", user)
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := c.userService.UpdateUser(uint(userID), &req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Failed to update user", err.Error())
		return
	}

	response.Success(ctx, "User updated successfully", user)
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	err = c.userService.DeleteUser(uint(userID))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Failed to delete user", err.Error())
		return
	}

	response.Success(ctx, "User deleted successfully", nil)
}

// ChangePassword 修改用户密码
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// 验证确认密码
	if req.NewPassword != req.ConfirmPassword {
		response.Error(ctx, http.StatusBadRequest, "Password confirmation failed",
			"New password and confirm password do not match")
		return
	}

	// 验证密码复杂度
	if !c.userService.ValidatePassword(req.NewPassword) {
		response.Error(ctx, http.StatusBadRequest, "Password validation failed",
			"Password must contain uppercase, lowercase letters and numbers")
		return
	}

	err = c.userService.ChangePassword(uint(userID), req.OldPassword, req.NewPassword)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Failed to change password", err.Error())
		return
	}

	response.Success(ctx, "Password changed successfully", nil)
}
