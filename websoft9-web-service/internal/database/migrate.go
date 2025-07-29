package database

import (
	"websoft9-web-service/internal/model"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// 3.5 平台管理 (Platform Management)
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

		// 3.1 应用商店 (App Store)
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

		// 3.2 工作空间 (Workspace)
		&model.UserFile{},
		&model.Workflow{},
		&model.WorkflowTask{},
		&model.WorkflowExecution{},

		// 3.3 资源管理 (Resource Management)
		&model.ResourceGroup{},
		&model.DatabaseConnection{},
		&model.Server{},
		&model.ServerAgent{},
		&model.AppInstance{},
		&model.Application{}, // 兼容现有代码

		// 3.4 安全管控 (Security Management)
		&model.SSLCertificate{},
		&model.SecretKey{},
		&model.AppGateway{},
		&model.AppGatewayPublish{},
		&model.AppGatewayAccessRule{},
		&model.AuditLog{},
		&model.Gateway{}, // 兼容现有代码
	)
}
