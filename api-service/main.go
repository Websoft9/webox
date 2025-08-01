package main

import (
	"api-service/internal/config"
	"api-service/internal/database"
	"api-service/internal/router"
	"api-service/internal/service"
	"log"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 自动迁移数据库表结构
	// if err := database.AutoMigrate(db); err != nil {
	// 	log.Fatal("Failed to migrate database:", err)
	// }

	// 初始化Redis
	rdb, err := database.InitRedis(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// 初始化InfluxDB
	influxClient, err := database.InitInfluxDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize InfluxDB:", err)
	}

	// 初始化服务
	services := service.NewServices(db, rdb, influxClient, cfg)

	// 初始化路由
	r := router.SetupRouter(services, cfg)

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
