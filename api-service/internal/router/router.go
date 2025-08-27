package router

import (
	"api-service/internal/config"
	"api-service/internal/controller"
	"api-service/internal/middleware"
	"api-service/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controllers controller collection
type Controllers struct {
	UserController       *controller.UserController
	I18nController       *controller.I18nController
	RoleController       *controller.RoleController
	PermissionController *controller.PermissionController
	APITokenController   *controller.APITokenController
	TwoFactorController  *controller.TwoFactorController
	AuthConfigController *controller.AuthConfigController
	// More controllers can be added
	// AppController  *controller.ApplicationController
}

// SetupRouter sets up router
func SetupRouter(controllers *Controllers, cfg *config.Config, log logger.Logger) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// Global middleware
	setupMiddleware(r, log)

	// Health check
	setupHealthCheck(r)

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
func setupMiddleware(r *gin.Engine, log logger.Logger) {
	r.Use(middleware.LoggerMiddleware(log))
	r.Use(middleware.CORS())
	r.Use(middleware.I18nMiddleware())
	r.Use(middleware.ErrorHandler(log))
	r.Use(middleware.RequestValidator(log))
}

// setupHealthCheck sets up health check routes
func setupHealthCheck(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "API service running normally",
		})
	})
}

// setupAPIRoutes sets up API routes
func setupAPIRoutes(v1 *gin.RouterGroup, controllers *Controllers) {
	setupUserRoutes(v1, controllers.UserController)
	setupI18nRoutes(v1, controllers.I18nController)

	// Routes requiring JWT authentication
	protected := v1.Group("/")
	// TODO: JWT middleware will be added here when config is available

	setupProtectedUserRoutes(protected, controllers.UserController)
	setupProtectedRoleRoutes(protected, controllers.RoleController)
	setupProtectedPermissionRoutes(protected, controllers.PermissionController)
	setupProtectedAPITokenRoutes(protected, controllers.APITokenController)
	setupProtectedTwoFactorRoutes(protected, controllers.TwoFactorController)
	setupProtectedAuthConfigRoutes(protected, controllers.AuthConfigController)
}

// setupUserRoutes sets up user related routes
func setupUserRoutes(v1 *gin.RouterGroup, userController *controller.UserController) {
	// Authentication related routes (no JWT verification required)
	auth := v1.Group("/auth")
	auth.POST("/register", userController.Register)
	auth.POST("/login", userController.Login)
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

// setupProtectedUserRoutes sets up user management routes
func setupProtectedUserRoutes(protected *gin.RouterGroup, userController *controller.UserController) {
	users := protected.Group("/users")
	// Current user operations
	users.GET("/profile", userController.GetProfile)
	users.PUT("/profile", userController.UpdateProfile)
	users.PUT("/password", userController.ChangePassword)

	// User management operations (admin permissions required)
	users.POST("", userController.CreateUser)
	users.GET("", userController.ListUsers)
	users.GET("/:id", userController.GetUser)
	users.PUT("/:id", userController.UpdateUser)
	users.PUT("/:id/status", userController.UpdateUserStatus)
	users.PUT("/:id/password", userController.UpdateUserPassword)
	users.DELETE("/:id", userController.DeleteUser)
}

// setupProtectedRoleRoutes sets up role management routes
func setupProtectedRoleRoutes(protected *gin.RouterGroup, roleController *controller.RoleController) {
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

// setupProtectedPermissionRoutes sets up permission management routes
func setupProtectedPermissionRoutes(protected *gin.RouterGroup, permissionController *controller.PermissionController) {
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

// setupProtectedAPITokenRoutes sets up API token management routes
func setupProtectedAPITokenRoutes(protected *gin.RouterGroup, apiTokenController *controller.APITokenController) {
	if apiTokenController == nil {
		return
	}

	apiTokens := protected.Group("/api-tokens")
	apiTokens.GET("", apiTokenController.ListAPITokens)
	apiTokens.POST("", apiTokenController.CreateAPIToken)
	apiTokens.GET("/:id", apiTokenController.GetAPIToken)
	apiTokens.PUT("/:id", apiTokenController.UpdateAPIToken)
	apiTokens.DELETE("/:id", apiTokenController.RevokeAPIToken)
	apiTokens.POST("/:id/refresh", apiTokenController.RefreshAPIToken)
	apiTokens.DELETE("", apiTokenController.BatchRevokeAPITokens)
}

// setupProtectedTwoFactorRoutes sets up two-factor authentication routes
func setupProtectedTwoFactorRoutes(protected *gin.RouterGroup, twoFactorController *controller.TwoFactorController) {
	if twoFactorController == nil {
		return
	}

	// Two-factor authentication routes under users
	users := protected.Group("/users")
	twoFactor := users.Group("/:user_id/two-factor")
	twoFactor.GET("", twoFactorController.GetTwoFactorStatus)
	twoFactor.POST("/enable", twoFactorController.EnableTOTP)
	twoFactor.POST("/confirm", twoFactorController.ConfirmTOTP)
	twoFactor.POST("/disable", twoFactorController.DisableTwoFactor)
	twoFactor.POST("/verify", twoFactorController.VerifyTwoFactor)
	twoFactor.POST("/totp/generate", twoFactorController.GenerateTOTPSecret)
}

// setupProtectedAuthConfigRoutes sets up authentication config routes
func setupProtectedAuthConfigRoutes(protected *gin.RouterGroup, authConfigController *controller.AuthConfigController) {
	if authConfigController == nil {
		return
	}

	// Authentication config routes
	authConfig := protected.Group("/auth-config")
	authConfig.GET("", authConfigController.GetAuthConfig)
	authConfig.PUT("", authConfigController.UpdateAuthConfig)
	authConfig.GET("/oauth2-providers", authConfigController.GetOAuth2Providers)
}
