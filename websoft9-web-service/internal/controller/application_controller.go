package controller

import (
	"net/http"
	"strconv"
	"websoft9-web-service/internal/model"
	"websoft9-web-service/internal/service"
	"websoft9-web-service/pkg/response"

	"github.com/gin-gonic/gin"
)

type ApplicationController struct {
	appService service.ApplicationService
}

func NewApplicationController(appService service.ApplicationService) *ApplicationController {
	return &ApplicationController{
		appService: appService,
	}
}

type CreateApplicationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Version     string `json:"version"`
	ServerID    uint   `json:"server_id" binding:"required"`
	Config      string `json:"config"`
}

func (c *ApplicationController) CreateApplication(ctx *gin.Context) {
	var req CreateApplicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	app := &model.Application{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Version:     req.Version,
		ServerID:    req.ServerID,
		Config:      req.Config,
		Status:      "created",
	}

	if err := c.appService.CreateApplication(app); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to create application", err.Error())
		return
	}

	response.Success(ctx, "Application created successfully", app)
}

func (c *ApplicationController) GetApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid application ID", err.Error())
		return
	}

	app, err := c.appService.GetApplication(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "Application not found", err.Error())
		return
	}

	response.Success(ctx, "Application retrieved successfully", app)
}

func (c *ApplicationController) ListApplications(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	apps, total, err := c.appService.ListApplications(page, pageSize)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get applications", err.Error())
		return
	}

	response.Success(ctx, "Applications retrieved successfully", gin.H{
		"applications": apps,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
	})
}

func (c *ApplicationController) DeployApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid application ID", err.Error())
		return
	}

	if err := c.appService.DeployApplication(uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to deploy application", err.Error())
		return
	}

	response.Success(ctx, "Application deployed successfully", nil)
}

func (c *ApplicationController) StopApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid application ID", err.Error())
		return
	}

	if err := c.appService.StopApplication(uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to stop application", err.Error())
		return
	}

	response.Success(ctx, "Application stopped successfully", nil)
}

func (c *ApplicationController) RestartApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid application ID", err.Error())
		return
	}

	if err := c.appService.RestartApplication(uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to restart application", err.Error())
		return
	}

	response.Success(ctx, "Application restarted successfully", nil)
}