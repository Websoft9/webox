package database

import (
	"websoft9-web-service/internal/model"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// 用户权限相关
		&model.User{},
		&model.Role{},
		&model.Permission{},

		// 资源管理相关
		&model.Server{},
		&model.Application{},
		&model.Gateway{},
	)
}
