package main

import (
	"fmt"
	"log"
	"os"
	"websoft9-web-service/internal/config"
	"websoft9-web-service/internal/database"
	"websoft9-web-service/internal/router"
	"websoft9-web-service/internal/service"

	"github.com/go-redis/redis/v9"
	"golang.org/x/net/context"
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

	// 从环境变量读取 Redis 配置
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" || redisPort == "" {
		fmt.Println("REDIS_HOST or REDIS_PORT is not set")
		os.Exit(1)
	}

	// 构建 Redis 地址
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// 测试连接
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connected to Redis successfully!")

	// 初始化InfluxDB
	influxClient, err := database.InitInfluxDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize InfluxDB:", err)
	}

	// 初始化服务
	services := service.NewServices(db, rdb, influxClient)

	// 初始化路由
	r := router.SetupRouter(services)

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
