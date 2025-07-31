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
