package main

import (
	"api-service/internal/config"
	"api-service/internal/controller"
	"api-service/internal/model"
	"api-service/internal/repository"
	"api-service/internal/router"
	"api-service/internal/service"
	"api-service/pkg/auth"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
	"log"
)

func main() {
	// 1. 初始化日志系统
	zapLogger := logger.NewDefaultZapLogger()
	logger.SetDefault(zapLogger)

	zapLogger.Info("应用程序启动")

	// 2. 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	zapLogger.Info("配置加载成功")

	// 3. 初始化数据库
	db, err := utils.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	zapLogger.Info("数据库连接成功")

	// 4. 数据库迁移
	if migrateErr := db.AutoMigrate(&model.User{}); migrateErr != nil {
		log.Fatal("Failed to migrate database:", migrateErr)
	}
	zapLogger.Info("数据库迁移完成")

	// 5. 初始化Redis
	_, err = utils.InitRedis(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}
	zapLogger.Info("Redis连接成功")

	// 6. 初始化InfluxDB
	_, err = utils.InitInfluxDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize InfluxDB:", err)
	}
	zapLogger.Info("InfluxDB连接成功")

	// 7. 初始化JWT认证
	jwtAuth := auth.NewJWTAuth(cfg.JWT.Secret, cfg.JWT.ExpireTime)

	// 8. 初始化Repository层
	userRepo := repository.NewUserRepository(db)

	// 9. 初始化Service层
	userService := service.NewUserService(userRepo, jwtAuth, zapLogger)

	// 10. 初始化Controller层
	userController := controller.NewUserController(userService, zapLogger)

	// 11. 初始化路由
	r := router.SetupRouter(&router.Controllers{
		UserController: userController,
	}, cfg, zapLogger)

	// 12. 启动服务器
	zapLogger.Info("服务器启动", logger.String("port", cfg.Server.Port))
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
