package router

import (
	"api-service/internal/controller"
	"api-service/internal/middleware"
	"api-service/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(services *service.Services) *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 初始化控制器
	userController := controller.NewUserController(services.UserService)
	appController := controller.NewApplicationController(services.ApplicationService)

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}

		// 需要认证的路由
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			// 用户相关路由
			users := protected.Group("/users")
			{
				users.GET("/profile", userController.GetProfile)
				users.GET("/", userController.ListUsers)
			}

			// 应用相关路由
			applications := protected.Group("/applications")
			{
				applications.POST("/", appController.CreateApplication)
				applications.GET("/", appController.ListApplications)
				applications.GET("/:id", appController.GetApplication)
				applications.POST("/:id/deploy", appController.DeployApplication)
				applications.POST("/:id/stop", appController.StopApplication)
				applications.POST("/:id/restart", appController.RestartApplication)
			}

			// 监控相关路由
			monitoring := protected.Group("/monitoring")
			{
				monitoring.GET("/servers/:id/metrics", func(c *gin.Context) {
					// TODO: 实现服务器监控数据获取
				})
				monitoring.GET("/applications/:id/metrics", func(c *gin.Context) {
					// TODO: 实现应用监控数据获取
				})
			}

			// 网关相关路由
			gateway := protected.Group("/gateway")
			{
				gateway.GET("/", func(c *gin.Context) {
					// TODO: 实现网关列表获取
				})
				gateway.POST("/", func(c *gin.Context) {
					// TODO: 实现网关创建
				})
			}
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
