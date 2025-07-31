package service

import (
	"api-service/internal/repository"
	"api-service/pkg/auth"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Services struct {
	UserService        UserService
	ApplicationService ApplicationService
	MonitorService     MonitorService
}

func NewServices(db *gorm.DB, rdb *redis.Client, influxClient influxdb2.Client) *Services {
	// 初始化JWT认证
	jwtAuth := auth.NewJWTAuth("your-secret-key", 3600)

	// 初始化Repository
	userRepo := repository.NewUserRepository(db)
	appRepo := repository.NewApplicationRepository(db)

	// 初始化Service
	userService := NewUserService(userRepo, jwtAuth)
	appService := NewApplicationService(appRepo)
	monitorService := NewMonitorService(influxClient)

	return &Services{
		UserService:        userService,
		ApplicationService: appService,
		MonitorService:     monitorService,
	}
}
