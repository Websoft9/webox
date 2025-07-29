package database

import (
	"websoft9-web-service/internal/model"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.Server{},
		&model.Application{},
		&model.Gateway{},
	)
}
