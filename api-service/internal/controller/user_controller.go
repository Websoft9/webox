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

// Request structs
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,max=64"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=128"`
	Nickname string `json:"nickname" binding:"max=100"`
	Phone    string `json:"phone" binding:"max=20"`
	Gender   int8   `json:"gender" binding:"min=0,max=2"`
	GroupID  uint   `json:"group_id" binding:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,max=64"`
	Email    string `json:"email" binding:"omitempty,email,max=255"`
	Nickname string `json:"nickname" binding:"omitempty,max=100"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Gender   int8   `json:"gender" binding:"omitempty,min=0,max=2"`
	Avatar   string `json:"avatar" binding:"omitempty,max=255"`
	GroupID  uint   `json:"group_id" binding:"omitempty"`
	Status   int8   `json:"status" binding:"omitempty,min=0,max=1"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128"`
}

func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数验证失败", err.Error())
		return
	}

	user, err := c.userService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "用户注册失败", err.Error())
		return
	}

	response.Success(ctx, "success", user)
}

func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数验证失败", err.Error())
		return
	}

	token, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, "登录失败", err.Error())
		return
	}

	response.Success(ctx, "success", gin.H{"token": token})
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "认证失败", "用户ID未找到")
		return
	}

	user, err := c.userService.GetProfile(userID.(uint))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "用户不存在", err.Error())
		return
	}

	response.Success(ctx, "success", user)
}

// ListUsers - GET /api/v1/users - 用户列表查询（支持分页、搜索、筛选）
func (c *UserController) ListUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")
	groupID := ctx.Query("group_id")

	users, total, err := c.userService.ListUsers(page, pageSize, keyword, status, groupID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "获取用户列表失败", err.Error())
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	hasNext := int64(page) < totalPages
	hasPrev := page > 1

	response.Success(ctx, "success", gin.H{
		"items": users,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    hasNext,
			"has_prev":    hasPrev,
		},
	})
}

// GetUserByID - GET /api/v1/users/{id} - 用户详情查询
func (c *UserController) GetUserByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误", "无效的用户ID")
		return
	}

	user, err := c.userService.GetProfile(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "用户不存在", err.Error())
		return
	}

	response.Success(ctx, "success", user)
}

// CreateUser - POST /api/v1/users - 用户创建
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数验证失败", err.Error())
		return
	}

	// 转换为service层的请求结构
	serviceReq := &service.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Gender:   req.Gender,
		GroupID:  req.GroupID,
	}

	user, err := c.userService.CreateUser(serviceReq)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "创建用户失败", err.Error())
		return
	}

	response.Success(ctx, "success", user)
}

// UpdateUser - PUT /api/v1/users/{id} - 用户更新
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误", "无效的用户ID")
		return
	}

	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数验证失败", err.Error())
		return
	}

	// 转换为service层的请求结构
	serviceReq := &service.UpdateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Gender:   req.Gender,
		Avatar:   req.Avatar,
		GroupID:  req.GroupID,
		Status:   req.Status,
	}

	user, err := c.userService.UpdateUser(uint(id), serviceReq)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "更新用户失败", err.Error())
		return
	}

	response.Success(ctx, "success", user)
}

// DeleteUser - DELETE /api/v1/users/{id} - 用户删除
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误", "无效的用户ID")
		return
	}

	err = c.userService.DeleteUser(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "删除用户失败", err.Error())
		return
	}

	response.Success(ctx, "success", gin.H{"id": id})
}

// ChangePassword - PUT /api/v1/users/{id}/password - 密码修改
func (c *UserController) ChangePassword(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数错误", "无效的用户ID")
		return
	}

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "参数验证失败", err.Error())
		return
	}

	err = c.userService.ChangePassword(uint(id), req.OldPassword, req.NewPassword)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "密码修改失败", err.Error())
		return
	}

	response.Success(ctx, "success", gin.H{"message": "密码修改成功"})
}
