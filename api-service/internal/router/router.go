package router

import (
	"api-service/internal/config"
	"api-service/internal/controller"
	serviceInterface "api-service/internal/interface/service"
	"api-service/internal/middleware"
	"api-service/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Controllers controller collection
type Controllers struct {
	UserController           *controller.UserController
	UserAuthController       *controller.UserAuthController
	I18nController           *controller.I18nController
	RolePermissionController *controller.RolePermissionController
	SecurityController       *controller.SecurityController
	HealthController         *controller.HealthController
	AuditLogController       *controller.AuditLogController
	UserProfileController    *controller.UserProfileController
	SystemConfigController   *controller.SystemConfigController

	// More controllers can be added
	// AppController  *controller.ApplicationController
}

// SetupRouter sets up router
func SetupRouter(
	controllers *Controllers,
	cfg *config.Config,
	log logger.Logger,
	permissionService serviceInterface.PermissionService,
	apiTokenService serviceInterface.APITokenService,
	auditLogService serviceInterface.AuditLogService,
) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// Health check
	setupHealthCheck(r, controllers.HealthController)

	// Swagger documentation (without middleware)
	setupSwaggerRoute(r, cfg)

	// Global middleware
	setupMiddleware(r, cfg, log, permissionService, apiTokenService, auditLogService)

	// API routes
	v1 := r.Group("/api/v1")
	setupAPIRoutes(v1, controllers)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "API endpoint not found",
			"path":    c.Request.URL.Path,
		})
	})

	return r
}

// setupMiddleware sets up global middleware
func setupMiddleware(
	r *gin.Engine,
	cfg *config.Config,
	log logger.Logger,
	permissionService serviceInterface.PermissionService,
	apiTokenService serviceInterface.APITokenService,
	auditLogService serviceInterface.AuditLogService,
) {
	r.Use(middleware.LoggerMiddleware(log))
	r.Use(middleware.CORS())
	r.Use(middleware.I18nMiddleware())
	r.Use(middleware.PermissionMiddleware(permissionService, apiTokenService, cfg, log))
	r.Use(middleware.AuditLogMiddleware(auditLogService, log))
	r.Use(middleware.ErrorHandler(log))
	r.Use(middleware.RequestValidator(log))
}

// setupHealthCheck sets up health check routes
func setupHealthCheck(r *gin.Engine, healthController *controller.HealthController) {
	if healthController == nil {
		// Fallback basic health check if controller is not available
		r.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "API service running normally",
			})
		})
		return
	}

	// Basic health endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"message":   "API service running normally",
			"timestamp": time.Now().Unix(),
		})
	})

	// Health controller endpoints
	r.GET("/ping", healthController.Ping)
	r.GET("/readiness", healthController.Readiness)
	r.GET("/liveness", healthController.Liveness)

	// Health check API group
	healthGroup := r.Group("/health")
	healthGroup.GET("/system", healthController.SystemHealth)
	healthGroup.GET("/database", healthController.DatabaseHealth)
	healthGroup.GET("/database/stats", healthController.DatabaseStats)
}

// setupSwaggerRoute sets up swagger documentation route
func setupSwaggerRoute(r *gin.Engine, cfg *config.Config) {
	// Only expose swagger in development mode
	if cfg.Server.Mode == "debug" || cfg.Server.Mode == "test" {
		// Add swagger route without middleware
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

// setupAPIRoutes sets up API routes
func setupAPIRoutes(v1 *gin.RouterGroup, controllers *Controllers) {
	// Routes requiring JWT authentication
	protected := v1.Group("/")
	setupUserAuthRoutes(v1, controllers.UserAuthController)
	setupI18nRoutes(v1, controllers.I18nController)
	setupUserRoutes(protected, controllers.UserController)
	setupRoleRoutes(protected, controllers.RolePermissionController)
	setupPermissionRoutes(protected, controllers.RolePermissionController)
	setupAPITokenRoutes(protected, controllers.SecurityController)
	setupTwoFactorRoutes(protected, controllers.SecurityController)
	setupAuthConfigRoutes(protected, controllers.SecurityController)
	setupAuditLogRoutes(protected, controllers.AuditLogController)
	setupUserProfileRoutes(protected, controllers.UserProfileController)
	setupSystemConfigRoutes(protected, controllers.SystemConfigController)
}

// setupUserRoutes sets up user related routes
func setupUserAuthRoutes(v1 *gin.RouterGroup, userAuthController *controller.UserAuthController) {
	// Authentication related routes (no JWT verification required)
	auth := v1.Group("/auth")
	auth.POST("/register", userAuthController.Register)
	auth.POST("/login", userAuthController.Login)
	auth.POST("/logout", userAuthController.Logout)
	auth.POST("/forgot-password", userAuthController.ForgotPassword)
	auth.GET("/reset-password", userAuthController.ShowResetPasswordForm)
	auth.POST("/reset-password", userAuthController.ResetPassword)
	auth.GET("/verify-email", userAuthController.VerifyEmail)
	auth.POST("/resend-verification", userAuthController.ResendVerificationEmail)
	auth.POST("/oauth2/login", userAuthController.OAuth2Login)
}

// setupI18nRoutes sets up internationalization routes
func setupI18nRoutes(v1 *gin.RouterGroup, i18nController *controller.I18nController) {
	if i18nController == nil {
		return
	}

	// i18n related routes (no JWT verification required)
	i18nGroup := v1.Group("/i18n")
	i18nGroup.GET("/languages", i18nController.GetLanguages)
	i18nGroup.GET("/translations/:lang", i18nController.GetTranslations)
	i18nGroup.GET("/test", i18nController.TestI18n)
}

// setupUserRoutes sets up user management routes
func setupUserRoutes(protected *gin.RouterGroup, userController *controller.UserController) {
	users := protected.Group("/users")

	// User management operations (admin permissions required)
	users.POST("", userController.CreateUser)
	users.GET("", userController.ListUsers)
	users.GET("/:id", userController.GetUser)
	users.PUT("/:id", userController.UpdateUser)
	users.PUT("/:id/status", userController.UpdateUserStatus)
	users.PUT("/:id/password", userController.UpdateUserPassword)
	users.DELETE("/:id", userController.DeleteUser)
}

// setupRoleRoutes sets up role management routes
func setupRoleRoutes(protected *gin.RouterGroup, roleController *controller.RolePermissionController) {
	if roleController == nil {
		return
	}

	roles := protected.Group("/roles")
	roles.GET("", roleController.ListRoles)
	roles.POST("", roleController.CreateRole)
	roles.GET("/:id", roleController.GetRole)
	roles.PUT("/:id", roleController.UpdateRole)
	roles.DELETE("/:id", roleController.DeleteRole)

	// Role permission management
	roles.POST("/:id/permissions", roleController.AssignPermissions)
	roles.DELETE("/:id/permissions", roleController.RemovePermissions)
	roles.GET("/:id/users", roleController.GetRoleUsers)
}

// setupPermissionRoutes sets up permission management routes
func setupPermissionRoutes(protected *gin.RouterGroup, permissionController *controller.RolePermissionController) {
	if permissionController == nil {
		return
	}

	permissions := protected.Group("/permissions")
	permissions.GET("", permissionController.ListPermissions)
	permissions.GET("/tree", permissionController.GetPermissionTree)
	permissions.POST("", permissionController.CreatePermission)
	permissions.GET("/:id", permissionController.GetPermission)
	permissions.PUT("/:id", permissionController.UpdatePermission)
	permissions.DELETE("/:id", permissionController.DeletePermission)
	permissions.GET("/:id/roles", permissionController.GetPermissionRoles)
}

// setupAPITokenRoutes sets up API token management routes
func setupAPITokenRoutes(protected *gin.RouterGroup, securityController *controller.SecurityController) {
	if securityController == nil {
		return
	}

	apiTokens := protected.Group("/api-tokens")
	apiTokens.POST("/revoke", securityController.RevokeAPIToken)
	apiTokens.GET("/refresh", securityController.RefreshAPIToken)
}

// setupTwoFactorRoutes sets up two-factor authentication routes
func setupTwoFactorRoutes(protected *gin.RouterGroup, securityController *controller.SecurityController) {
	if securityController == nil {
		return
	}

	// Two-factor authentication routes under users
	twoFactor := protected.Group("/two-factor")
	twoFactor.GET("", securityController.GetTwoFactorStatus)
	twoFactor.POST("/enable", securityController.EnableTOTP)
	twoFactor.POST("/confirm", securityController.ConfirmTOTP)
	twoFactor.POST("/disable", securityController.DisableTwoFactor)
	twoFactor.POST("/verify", securityController.VerifyTwoFactor)
	twoFactor.POST("/totp/generate", securityController.GenerateTOTPSecret)
}

// setupAuthConfigRoutes sets up authentication config routes
func setupAuthConfigRoutes(protected *gin.RouterGroup, securityController *controller.SecurityController) {
	if securityController == nil {
		return
	}

	// Authentication config routes
	authConfig := protected.Group("/auth-config")
	authConfig.GET("", securityController.GetAuthConfig)
	authConfig.PUT("", securityController.UpdateAuthConfig)
	authConfig.GET("/oauth2-providers", securityController.GetOAuth2Providers)
}

// setupAuditLogRoutes sets up audit log routes
func setupAuditLogRoutes(protected *gin.RouterGroup, auditLogController *controller.AuditLogController) {
	if auditLogController == nil {
		return
	}

	// Audit log routes
	auditLogs := protected.Group("/audit-logs")
	auditLogs.GET("", auditLogController.ListAuditLogs)
	auditLogs.GET("/:id", auditLogController.GetAuditLog)
	auditLogs.GET("/statistics", auditLogController.GetAuditLogStatistics)
	auditLogs.GET("/export", auditLogController.ExportAuditLogs)
}

// setupUserProfileRoutes sets up user profile routes
func setupUserProfileRoutes(protected *gin.RouterGroup, userProfileController *controller.UserProfileController) {
	if userProfileController == nil {
		return
	}

	profile := protected.Group("/profile")
	profile.GET("", userProfileController.GetProfile)
	profile.PUT("", userProfileController.UpdateProfile)
	profile.PUT("/password", userProfileController.ChangePassword)
	profile.GET("/login-history", userProfileController.GetLoginHistories)
	profile.GET("/notification-settings", userProfileController.GetNotificationSettings)
	profile.PUT("/notification-settings", userProfileController.UpdateNotificationSettings)
	profile.GET("/security-settings", userProfileController.GetSecuritySettings)
	profile.PUT("/security-settings", userProfileController.UpdateSecuritySettings)
}

// setupSystemConfigRoutes sets up system configuration routes
func setupSystemConfigRoutes(protected *gin.RouterGroup, systemConfigController *controller.SystemConfigController) {
	if systemConfigController == nil {
		return
	}

	// System configuration routes
	systemConfigs := protected.Group("/system-configs")
	systemConfigs.GET("", systemConfigController.ListSystemConfigs)
	systemConfigs.PUT("", systemConfigController.BatchUpdateSystemConfigs)
	systemConfigs.GET("/basic", systemConfigController.ListBasicConfigs)
	systemConfigs.GET("/security", systemConfigController.ListSecurityConfigs)
	systemConfigs.GET("/email", systemConfigController.ListEmailConfigs)
	systemConfigs.POST("/smtp/test", systemConfigController.TestSMTP)
}
