package controller

import (
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

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=32"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=32"`
}

func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
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
	var req LoginRequest
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

func (c *UserController) ListUsers(ctx *gin.Context) {
	var params service.UserQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.Sort == "" {
		params.Sort = "created_at"
	}
	if params.Order == "" {
		params.Order = "desc"
	}

	result, err := c.userService.ListUsersWithFilter(&params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get users", err.Error())
		return
	}

	response.Success(ctx, "success", result)
}

// GetUser 获取用户详情
func (c *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "User not found", err.Error())
		return
	}

	response.Success(ctx, "success", user)
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req service.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
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
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var req service.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := c.userService.UpdateUser(uint(id), &req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Failed to update user", err.Error())
		return
	}

	response.Success(ctx, "User updated successfully", user)
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	err = c.userService.DeleteUser(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Failed to delete user", err.Error())
		return
	}

	response.Success(ctx, "User deleted successfully", nil)
}

// ChangeUserPassword 修改用户密码
func (c *UserController) ChangeUserPassword(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// 验证新密码和确认密码是否一致
	if req.NewPassword != req.ConfirmPassword {
		response.Error(ctx, http.StatusBadRequest, "Password mismatch", "New password and confirm password do not match")
		return
	}

	err = c.userService.ChangePassword(uint(id), req.OldPassword, req.NewPassword)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Failed to change password", err.Error())
		return
	}

	response.Success(ctx, "Password changed successfully", nil)
}
