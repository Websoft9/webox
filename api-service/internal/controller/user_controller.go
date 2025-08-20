package controller

import (
	"api-service/internal/repository"
	"api-service/internal/service"
	"api-service/internal/validator"
	"api-service/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService   service.UserService
	userValidator *validator.UserValidator
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService:   userService,
		userValidator: validator.NewUserValidator(),
	}
}

// GetProfile 获取用户资料
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

// ListUsers 用户列表（旧版本，保持兼容性）
func (c *UserController) ListUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	users, total, err := c.userService.ListUsers(page, pageSize)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get users", err.Error())
		return
	}

	response.Success(ctx, "Users retrieved successfully", gin.H{
		"users":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ========================================
// 按照API设计说明书实现的新接口
// ========================================

// GetUsers GET /api/v1/users - 用户列表查询（支持分页、搜索、筛选）
func (c *UserController) GetUsers(ctx *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	// 构建筛选条件
	filters := repository.UserFilters{
		Keyword: ctx.Query("keyword"),
		Sort:    ctx.DefaultQuery("sort", "created_at"),
		Order:   ctx.DefaultQuery("order", "desc"),
	}

	// 状态筛选
	if statusStr := ctx.Query("status"); statusStr != "" {
		if status, err := strconv.ParseInt(statusStr, 10, 8); err == nil {
			statusInt8 := int8(status)
			filters.Status = &statusInt8
		}
	}

	// 角色筛选
	if roleIDStr := ctx.Query("role_id"); roleIDStr != "" {
		if roleID, err := strconv.ParseUint(roleIDStr, 10, 32); err == nil {
			roleIDUint := uint(roleID)
			filters.RoleID = &roleIDUint
		}
	}

	// 用户组筛选
	if groupIDStr := ctx.Query("group_id"); groupIDStr != "" {
		if groupID, err := strconv.ParseUint(groupIDStr, 10, 32); err == nil {
			groupIDUint := uint(groupID)
			filters.GroupID = &groupIDUint
		}
	}

	result, err := c.userService.GetUsersList(page, pageSize, filters)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get users", err.Error())
		return
	}

	response.Success(ctx, "success", result)
}

// GetUserDetail GET /api/v1/users/{id} - 用户详情查询
func (c *UserController) GetUserDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	result, err := c.userService.GetUserDetail(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			response.Error(ctx, http.StatusNotFound, "User not found", err.Error())
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "Failed to get user detail", err.Error())
		return
	}

	response.Success(ctx, "success", result)
}

// CreateUser POST /api/v1/users - 用户创建
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req service.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// 验证请求参数
	if validationErrors := c.userValidator.ValidateCreateUser(req); len(validationErrors) > 0 {
		response.ValidationError(ctx, validationErrors)
		return
	}

	user, err := c.userService.CreateUser(req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			statusCode = http.StatusConflict
		}
		response.Error(ctx, statusCode, "Failed to create user", err.Error())
		return
	}

	response.Success(ctx, "User created successfully", user)
}

// UpdateUser PUT /api/v1/users/{id} - 用户更新
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

	// 验证请求参数
	if validationErrors := c.userValidator.ValidateUpdateUser(req); len(validationErrors) > 0 {
		response.ValidationError(ctx, validationErrors)
		return
	}

	err = c.userService.UpdateUser(uint(id), req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "username already exists" || err.Error() == "email already exists" {
			statusCode = http.StatusConflict
		}
		response.Error(ctx, statusCode, "Failed to update user", err.Error())
		return
	}

	response.Success(ctx, "User updated successfully", nil)
}

// DeleteUser DELETE /api/v1/users/{id} - 用户删除
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	err = c.userService.DeleteUser(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			response.Error(ctx, http.StatusNotFound, "User not found", err.Error())
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	response.Success(ctx, "User deleted successfully", nil)
}

// ChangePassword PUT /api/v1/users/{id}/password - 密码修改
func (c *UserController) ChangePassword(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var req service.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// 验证请求参数
	if validationErrors := c.userValidator.ValidateChangePassword(req); len(validationErrors) > 0 {
		response.ValidationError(ctx, validationErrors)
		return
	}

	err = c.userService.ChangeUserPassword(uint(id), req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "old password is incorrect" {
			statusCode = http.StatusUnauthorized
		}
		response.Error(ctx, statusCode, "Failed to change password", err.Error())
		return
	}

	response.Success(ctx, "Password changed successfully", nil)
}
