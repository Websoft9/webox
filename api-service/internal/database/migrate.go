package database

import (
	"api-service/internal/model"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// 系统管理相关表
		&model.UserGroup{},
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
		&model.APIToken{},
		&model.UserLoginHistory{},
		&model.SystemConfig{},
		&model.AlertRule{},
		&model.AlertRecord{},
		&model.NotificationTemplate{},
		&model.NotificationRecord{},
		&model.Notification{},
		&model.UserNotification{},
		&model.WebhookConfig{},
		&model.WebhookLog{},
		&model.RepositoryConfig{},
		&model.PlatformUpdate{},
		&model.DockerSwarmNode{},
		&model.DockerImage{},
		&model.StorageQuota{},

		// 应用商店相关表
		&model.AppStoreCategory{},
		&model.AppStoreTemplate{},
		&model.AppStoreWishlist{},
		&model.AppStoreReview{},
		&model.AppStoreFavorite{},
		&model.AppStoreStar{},
		&model.AppStoreReport{},
		&model.AppStoreDownload{},
		&model.AppStoreWishlistComment{},
		&model.AppStoreWishlistVote{},
		&model.AppStoreWishlistLike{},
		&model.AppStoreWishlistReport{},
		&model.AppDeployment{},
		&model.AppShortcut{},

		// 资源管理相关表
		&model.ResourceGroup{},
		&model.DatabaseConnection{},
		&model.Server{},
		&model.ServerAgent{},
		&model.AppInstance{},
		&model.Application{},
	)
}
