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
	// More controllers can be added
	// AppController  *controller.ApplicationController
}

// SetupRouter 设置路由
func SetupRouter(controllers *Controllers, cfg *config.Config, log logger.Logger) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 全局中间件
	r.Use(middleware.LoggerMiddleware(log))
	r.Use(middleware.CORS())
	r.Use(middleware.I18nMiddleware())
	r.Use(middleware.ErrorHandler(log))
	r.Use(middleware.RequestValidator(log))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "API服务运行正常",
		})
	})

	// API版本1路由组
	v1 := r.Group("/api/v1")

	// 认证相关路由（无需JWT验证）
	auth := v1.Group("/auth")
	auth.POST("/register", controllers.UserController.Register)
	auth.POST("/login", controllers.UserController.Login)

	// i18n相关路由（无需JWT验证）
	if controllers.I18nController != nil {
		i18nGroup := v1.Group("/i18n")
		i18nGroup.GET("/languages", controllers.I18nController.GetLanguages)
		i18nGroup.GET("/translations/:lang", controllers.I18nController.GetTranslations)
		i18nGroup.GET("/test", controllers.I18nController.TestI18n)
	}

	// 需要JWT认证的路由
	protected := v1.Group("/")
	protected.Use(middleware.JWTAuth(cfg))

	// 用户相关路由
	users := protected.Group("/users")
	// 当前用户操作
	users.GET("/profile", controllers.UserController.GetProfile)
	users.PUT("/profile", controllers.UserController.UpdateProfile)
	users.PUT("/password", controllers.UserController.ChangePassword)

	// 用户管理操作（需要管理员权限）
	users.POST("", controllers.UserController.CreateUser)                     // 创建用户
	users.GET("", controllers.UserController.ListUsers)                       // 获取用户列表
	users.GET("/:id", controllers.UserController.GetUser)                     // 获取单个用户
	users.PUT("/:id", controllers.UserController.UpdateUser)                  // 更新用户信息
	users.PUT("/:id/status", controllers.UserController.UpdateUserStatus)     // 更新用户状态
	users.PUT("/:id/password", controllers.UserController.UpdateUserPassword) // 管理员修改用户密码
	users.DELETE("/:id", controllers.UserController.DeleteUser)               // 删除用户

	// 角色管理路由
	if controllers.RoleController != nil {
		roles := protected.Group("/roles")
		roles.GET("", controllers.RoleController.ListRoles)
		roles.POST("", controllers.RoleController.CreateRole)
		roles.GET("/:id", controllers.RoleController.GetRole)
		roles.PUT("/:id", controllers.RoleController.UpdateRole)
		roles.DELETE("/:id", controllers.RoleController.DeleteRole)

		// 角色权限管理
		roles.POST("/:id/permissions", controllers.RoleController.AssignPermissions)
		roles.DELETE("/:id/permissions", controllers.RoleController.RemovePermissions)

		// 角色用户管理
		roles.GET("/:id/users", controllers.RoleController.GetRoleUsers)
	}

	// Permission management routes
	if controllers.PermissionController != nil {
		permissions := protected.Group("/permissions")
		permissions.GET("", controllers.PermissionController.ListPermissions)
		permissions.POST("", controllers.PermissionController.CreatePermission)
		permissions.GET("/tree", controllers.PermissionController.GetPermissionTree)
		permissions.GET("/:id", controllers.PermissionController.GetPermission)
		permissions.PUT("/:id", controllers.PermissionController.UpdatePermission)
		permissions.DELETE("/:id", controllers.PermissionController.DeletePermission)

		// Permission role management
		permissions.GET("/:id/roles", controllers.PermissionController.GetPermissionRoles)
	}

	// API Token management routes
	if controllers.APITokenController != nil {
		apiTokens := protected.Group("/api-tokens")
		apiTokens.GET("", controllers.APITokenController.ListAPITokens)
		apiTokens.POST("", controllers.APITokenController.CreateAPIToken)
		apiTokens.GET("/:id", controllers.APITokenController.GetAPIToken)
		apiTokens.PUT("/:id", controllers.APITokenController.UpdateAPIToken)
		apiTokens.POST("/:id/revoke", controllers.APITokenController.RevokeAPIToken)
		apiTokens.POST("/:id/refresh", controllers.APITokenController.RefreshAPIToken)

		// API Token validation (no JWT required)
		v1.POST("/api-tokens/validate", controllers.APITokenController.ValidateAPIToken)
	}

	// Two-factor authentication routes
	if controllers.TwoFactorController != nil {
		twoFactor := protected.Group("/2fa")

		// TOTP management
		twoFactor.POST("/totp/enable", controllers.TwoFactorController.EnableTOTP)
		twoFactor.POST("/totp/confirm", controllers.TwoFactorController.ConfirmTOTP)
		twoFactor.POST("/totp/disable", controllers.TwoFactorController.DisableTOTP)

		// Email 2FA management
		twoFactor.POST("/email/enable", controllers.TwoFactorController.EnableEmailTwoFactor)
		twoFactor.POST("/email/disable", controllers.TwoFactorController.DisableEmailTwoFactor)
		twoFactor.POST("/email/send-code", controllers.TwoFactorController.SendEmailCode)

		// General 2FA operations
		twoFactor.GET("/status", controllers.TwoFactorController.GetTwoFactorStatus)
		twoFactor.POST("/backup-codes", controllers.TwoFactorController.GenerateBackupCodes)

		// 2FA verification (no JWT required for login flow)
		v1.POST("/2fa/verify", controllers.TwoFactorController.VerifyTwoFactor)
	}

	// TODO: 应用相关路由将在后续版本中实现

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "接口不存在",
			"path":    c.Request.URL.Path,
		})
	})

	return r
}
