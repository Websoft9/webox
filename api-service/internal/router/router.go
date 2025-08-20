package router

import (
	"api-service/internal/config"
	"api-service/internal/controller"
	"api-service/internal/middleware"
	"api-service/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controllers 控制器集合
type Controllers struct {
	UserController *controller.UserController
	// 可以添加更多控制器
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
	{
		// 认证相关路由（无需JWT验证）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.UserController.Register)
			auth.POST("/login", controllers.UserController.Login)
		}

		// 需要JWT认证的路由
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth(cfg))
		{
			// 用户相关路由
			users := protected.Group("/users")
			{
				// 当前用户操作
				users.GET("/profile", controllers.UserController.GetProfile)
				users.PUT("/profile", controllers.UserController.UpdateProfile)
				users.PUT("/password", controllers.UserController.ChangePassword)

				// 用户管理操作（需要管理员权限）
				users.GET("", controllers.UserController.ListUsers)                   // 获取用户列表
				users.GET("/:id", controllers.UserController.GetUser)                 // 获取单个用户
				users.PUT("/:id/status", controllers.UserController.UpdateUserStatus) // 更新用户状态
				users.DELETE("/:id", controllers.UserController.DeleteUser)           // 删除用户
			}

			// 可以在这里添加更多受保护的路由
			// applications := protected.Group("/applications")
			// {
			//     applications.POST("/", controllers.AppController.CreateApplication)
			//     applications.GET("/", controllers.AppController.ListApplications)
			//     applications.GET("/:id", controllers.AppController.GetApplication)
			// }
		}
	}

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
