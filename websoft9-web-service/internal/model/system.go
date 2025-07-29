package model

import (
	"time"

	"gorm.io/gorm"
)

// ========================================
// 3.5 平台管理 (Platform Management)
// ========================================

// UserGroup 用户组表
type UserGroup struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"not null" binding:"required"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null" binding:"required"`
	Description string         `json:"description"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	Status      int8           `json:"status" gorm:"default:1"` // 0-禁用，1-启用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// User 用户表
type User struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	GroupID      uint           `json:"group_id" gorm:"not null"`
	Group        UserGroup      `json:"group" gorm:"foreignKey:GroupID"`
	Username     string         `json:"username" gorm:"uniqueIndex;not null" binding:"required"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null" binding:"required,email"`
	PasswordHash string         `json:"-" gorm:"column:password_hash;not null"`
	Nickname     string         `json:"nickname"`
	Avatar       string         `json:"avatar"`
	Phone        string         `json:"phone"`
	Gender       int8           `json:"gender" gorm:"default:0"` // 0-未知，1-男，2-女
	Signature    string         `json:"signature"`
	Status       int8           `json:"status" gorm:"default:1"` // 0-禁用，1-启用
	LastLoginAt  *time.Time     `json:"last_login_at"`
	LastLoginIP  string         `json:"last_login_ip"`
	Timezone     string         `json:"timezone" gorm:"default:UTC"`
	Language     string         `json:"language" gorm:"default:zh-CN"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Roles []Role `json:"roles" gorm:"many2many:user_roles;"`
}

// Role 角色表
type Role struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null" binding:"required"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null" binding:"required"`
	Description string         `json:"description"`
	IsSystem    int8           `json:"is_system" gorm:"default:0"` // 是否系统角色
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	Status      int8           `json:"status" gorm:"default:1"` // 0-禁用，1-启用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Users       []User       `json:"users" gorm:"many2many:user_roles;"`
}

// Permission 权限表
type Permission struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"not null" binding:"required"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null" binding:"required"`
	Module      string         `json:"module" gorm:"not null" binding:"required"`
	Action      string         `json:"action" gorm:"not null" binding:"required"`
	Resource    string         `json:"resource"`
	Description string         `json:"description"`
	IsSystem    int8           `json:"is_system" gorm:"default:0"` // 是否系统权限
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	Status      int8           `json:"status" gorm:"default:1"` // 0-禁用，1-启用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Roles []Role `json:"roles" gorm:"many2many:role_permissions;"`
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_role"`
	RoleID    uint      `json:"role_id" gorm:"not null;uniqueIndex:idx_user_role"`
	GrantedBy *uint     `json:"granted_by"`
	GrantedAt time.Time `json:"granted_at" gorm:"default:CURRENT_TIMESTAMP"`
	Status    int8      `json:"status" gorm:"default:1"` // 0-禁用，1-启用
	CreatedAt time.Time `json:"created_at"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	RoleID       uint      `json:"role_id" gorm:"not null;uniqueIndex:idx_role_permission"`
	PermissionID uint      `json:"permission_id" gorm:"not null;uniqueIndex:idx_role_permission"`
	GrantedBy    *uint     `json:"granted_by"`
	GrantedAt    time.Time `json:"granted_at" gorm:"default:CURRENT_TIMESTAMP"`
	Status       int8      `json:"status" gorm:"default:1"` // 0-禁用，1-启用
	CreatedAt    time.Time `json:"created_at"`
}

// APIToken API访问令牌表
type APIToken struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"not null" binding:"required"`
	Token       string         `json:"token" gorm:"uniqueIndex;not null"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Scopes      string         `json:"scopes" gorm:"type:json"` // JSON格式存储权限范围
	Description string         `json:"description"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	Status      int8           `json:"status" gorm:"default:1"` // 0-禁用，1-启用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserLoginHistory 用户登录历史表
type UserLoginHistory struct {
	ID         uint       `json:"id" gorm:"primarykey"`
	UserID     uint       `json:"user_id" gorm:"not null"`
	User       User       `json:"user" gorm:"foreignKey:UserID"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	Location   string     `json:"location"`
	Device     string     `json:"device"`
	Browser    string     `json:"browser"`
	LoginTime  time.Time  `json:"login_time" gorm:"default:CURRENT_TIMESTAMP"`
	LogoutTime *time.Time `json:"logout_time"`
	Status     string     `json:"status" gorm:"default:ACTIVE"` // ACTIVE, EXPIRED, LOGOUT
	SessionID  string     `json:"session_id"`
	CreatedAt  time.Time  `json:"created_at"`
}

// SystemConfig 系统配置表
type SystemConfig struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	ConfigKey       string         `json:"config_key" gorm:"uniqueIndex;not null" binding:"required"`
	ConfigValue     string         `json:"config_value" gorm:"type:text"`
	ConfigType      string         `json:"config_type" gorm:"default:STRING"` // STRING, INTEGER, BOOLEAN, JSON, TEXT
	Category        string         `json:"category" gorm:"not null" binding:"required"`
	Description     string         `json:"description" gorm:"type:text"`
	IsReadonly      int8           `json:"is_readonly" gorm:"default:0"`
	IsEncrypted     int8           `json:"is_encrypted" gorm:"default:0"`
	DefaultValue    string         `json:"default_value" gorm:"type:text"`
	ValidationRules string         `json:"validation_rules" gorm:"type:json"`
	SortOrder       int            `json:"sort_order" gorm:"default:0"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// AlertRule 告警规则表
type AlertRule struct {
	ID                   uint           `json:"id" gorm:"primarykey"`
	Name                 string         `json:"name" gorm:"not null" binding:"required"`
	RuleType             string         `json:"rule_type" gorm:"not null"`   // THRESHOLD, ANOMALY, CUSTOM
	TargetType           string         `json:"target_type" gorm:"not null"` // SERVER, APPLICATION, DATABASE, GATEWAY
	TargetID             *uint          `json:"target_id"`
	MetricName           string         `json:"metric_name"`
	ConditionExpression  string         `json:"condition_expression" gorm:"type:text;not null"`
	NotificationChannels string         `json:"notification_channels" gorm:"type:json"`
	IsEnabled            int8           `json:"is_enabled" gorm:"default:1"`
	OwnerID              uint           `json:"owner_id" gorm:"not null"`
	Owner                User           `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`
}

// AlertRecord 告警记录表
type AlertRecord struct {
	ID                   uint           `json:"id" gorm:"primarykey"`
	AlertRuleID          uint           `json:"alert_rule_id" gorm:"not null"`
	AlertRule            AlertRule      `json:"alert_rule" gorm:"foreignKey:AlertRuleID"`
	AlertID              string         `json:"alert_id" gorm:"uniqueIndex;not null"`
	Title                string         `json:"title" gorm:"not null"`
	Description          string         `json:"description" gorm:"type:text"`
	Status               string         `json:"status" gorm:"default:FIRING"` // FIRING, RESOLVED, ACKNOWLEDGED
	FiredAt              time.Time      `json:"fired_at" gorm:"default:CURRENT_TIMESTAMP"`
	ResolvedAt           *time.Time     `json:"resolved_at"`
	AcknowledgedAt       *time.Time     `json:"acknowledged_at"`
	AcknowledgedBy       *uint          `json:"acknowledged_by"`
	ResolutionNote       string         `json:"resolution_note" gorm:"type:text"`
	NotificationSent     int8           `json:"notification_sent" gorm:"default:0"`
	NotificationChannels string         `json:"notification_channels" gorm:"type:json"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`
}

// NotificationTemplate 通知模板表
type NotificationTemplate struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Name      string         `json:"name" gorm:"not null" binding:"required"`
	Type      string         `json:"type" gorm:"not null"` // EMAIL, SMS, WEBHOOK, PUSH
	Subject   string         `json:"subject"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Variables string         `json:"variables" gorm:"type:json"`
	IsSystem  int8           `json:"is_system" gorm:"default:0"`
	Status    int8           `json:"status" gorm:"default:1"` // 0-禁用，1-启用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// NotificationRecord 通知记录表
type NotificationRecord struct {
	ID            uint                  `json:"id" gorm:"primarykey"`
	TemplateID    *uint                 `json:"template_id"`
	Template      *NotificationTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	Type          string                `json:"type" gorm:"not null"` // EMAIL, SMS, WEBHOOK, PUSH
	Recipient     string                `json:"recipient" gorm:"not null"`
	Subject       string                `json:"subject"`
	Content       string                `json:"content" gorm:"type:text;not null"`
	Status        string                `json:"status" gorm:"default:PENDING"` // PENDING, SENT, FAILED
	SentAt        *time.Time            `json:"sent_at"`
	ErrorMsg      string                `json:"error_msg" gorm:"type:text"`
	RetryCount    int                   `json:"retry_count" gorm:"default:0"`
	ReferenceID   string                `json:"reference_id"`
	ReferenceType string                `json:"reference_type"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	DeletedAt     gorm.DeletedAt        `json:"-" gorm:"index"`
}

// Notification 通知消息表
type Notification struct {
	ID         uint                  `json:"id" gorm:"primarykey"`
	Title      string                `json:"title" gorm:"not null"`
	Content    string                `json:"content" gorm:"type:text;not null"`
	Type       string                `json:"type" gorm:"not null"`      // SYSTEM, USER, ALERT
	Level      string                `json:"level" gorm:"default:INFO"` // INFO, WARNING, ERROR
	SenderID   *uint                 `json:"sender_id"`
	Sender     *User                 `json:"sender" gorm:"foreignKey:SenderID"`
	TargetType string                `json:"target_type" gorm:"not null"` // USER, GROUP, ALL
	TargetIDs  string                `json:"target_ids" gorm:"type:json"`
	Channels   string                `json:"channels" gorm:"type:json"`
	TemplateID *uint                 `json:"template_id"`
	Template   *NotificationTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	Variables  string                `json:"variables" gorm:"type:json"`
	SentAt     *time.Time            `json:"sent_at"`
	Status     string                `json:"status" gorm:"default:PENDING"` // PENDING, SENT, FAILED
	ErrorMsg   string                `json:"error_msg" gorm:"type:text"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
	DeletedAt  gorm.DeletedAt        `json:"-" gorm:"index"`
}

// UserNotification 用户通知记录表
type UserNotification struct {
	ID             uint         `json:"id" gorm:"primarykey"`
	NotificationID uint         `json:"notification_id" gorm:"not null"`
	Notification   Notification `json:"notification" gorm:"foreignKey:NotificationID"`
	UserID         uint         `json:"user_id" gorm:"not null"`
	User           User         `json:"user" gorm:"foreignKey:UserID"`
	IsRead         int8         `json:"is_read" gorm:"default:0"`
	ReadAt         *time.Time   `json:"read_at"`
	IsDeleted      int8         `json:"is_deleted" gorm:"default:0"`
	DeletedAt      *time.Time   `json:"deleted_at"`
	CreatedAt      time.Time    `json:"created_at"`
}

// WebhookConfig Webhook配置表
type WebhookConfig struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	Name       string         `json:"name" gorm:"not null" binding:"required"`
	URL        string         `json:"url" gorm:"not null" binding:"required,url"`
	Secret     string         `json:"secret"`
	Events     string         `json:"events" gorm:"type:json;not null"`
	Headers    string         `json:"headers" gorm:"type:json"`
	Timeout    int            `json:"timeout" gorm:"default:30"`
	RetryCount int            `json:"retry_count" gorm:"default:3"`
	IsActive   int8           `json:"is_active" gorm:"default:1"`
	OwnerID    uint           `json:"owner_id" gorm:"not null"`
	Owner      User           `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// WebhookLog Webhook执行日志表
type WebhookLog struct {
	ID           uint          `json:"id" gorm:"primarykey"`
	WebhookID    uint          `json:"webhook_id" gorm:"not null"`
	Webhook      WebhookConfig `json:"webhook" gorm:"foreignKey:WebhookID"`
	EventType    string        `json:"event_type" gorm:"not null"`
	Payload      string        `json:"payload" gorm:"type:json;not null"`
	RequestID    string        `json:"request_id" gorm:"not null"`
	StatusCode   *int          `json:"status_code"`
	ResponseBody string        `json:"response_body" gorm:"type:text"`
	ResponseTime int           `json:"response_time" gorm:"default:0"` // 毫秒
	RetryCount   int           `json:"retry_count" gorm:"default:0"`
	Success      int8          `json:"success" gorm:"default:0"`
	ErrorMessage string        `json:"error_message" gorm:"type:text"`
	CreatedAt    time.Time     `json:"created_at"`
}

// RepositoryConfig 软件源配置表
type RepositoryConfig struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Name      string         `json:"name" gorm:"not null" binding:"required"`
	Type      string         `json:"type" gorm:"not null"` // DOCKER, APT, YUM, NPM
	URL       string         `json:"url" gorm:"not null" binding:"required,url"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	IsDefault int8           `json:"is_default" gorm:"default:0"`
	IsSystem  int8           `json:"is_system" gorm:"default:0"`
	Status    int8           `json:"status" gorm:"default:1"` // 0-禁用，1-启用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// PlatformUpdate 平台更新记录表
type PlatformUpdate struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	Version      string         `json:"version" gorm:"not null"`
	Changelog    string         `json:"changelog" gorm:"type:text"`
	DownloadURL  string         `json:"download_url"`
	FileSize     int64          `json:"file_size" gorm:"default:0"`
	Checksum     string         `json:"checksum"`
	Status       string         `json:"status" gorm:"default:PENDING"` // PENDING, DOWNLOADING, INSTALLING, SUCCESS, FAILED
	StartedAt    *time.Time     `json:"started_at"`
	CompletedAt  *time.Time     `json:"completed_at"`
	ErrorMessage string         `json:"error_message" gorm:"type:text"`
	BackupPath   string         `json:"backup_path"`
	UpdatedBy    *uint          `json:"updated_by"`
	Updater      *User          `json:"updater" gorm:"foreignKey:UpdatedBy"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// DockerSwarmNode 容器集群节点表
type DockerSwarmNode struct {
	ID            uint       `json:"id" gorm:"primarykey"`
	NodeID        string     `json:"node_id" gorm:"uniqueIndex;not null"`
	Hostname      string     `json:"hostname" gorm:"not null"`
	IPAddress     string     `json:"ip_address" gorm:"not null"`
	Role          string     `json:"role" gorm:"not null"`               // MANAGER, WORKER
	Status        string     `json:"status" gorm:"default:READY"`        // READY, DOWN, UNKNOWN
	Availability  string     `json:"availability" gorm:"default:ACTIVE"` // ACTIVE, PAUSE, DRAIN
	EngineVersion string     `json:"engine_version"`
	Labels        string     `json:"labels" gorm:"type:json"`
	Resources     string     `json:"resources" gorm:"type:json"`
	JoinedAt      *time.Time `json:"joined_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// DockerImage 容器镜像表
type DockerImage struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	ImageID      string         `json:"image_id" gorm:"uniqueIndex;not null"`
	Repository   string         `json:"repository" gorm:"not null"`
	Tag          string         `json:"tag" gorm:"not null"`
	Digest       string         `json:"digest"`
	Size         int64          `json:"size" gorm:"default:0"`
	Architecture string         `json:"architecture"`
	OS           string         `json:"os"`
	Status       string         `json:"status" gorm:"default:AVAILABLE"` // AVAILABLE, PULLING, ERROR
	UsageCount   int            `json:"usage_count" gorm:"default:0"`
	LastUsedAt   *time.Time     `json:"last_used_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// StorageQuota 存储配额表
type StorageQuota struct {
	ID               uint           `json:"id" gorm:"primarykey"`
	UserID           uint           `json:"user_id" gorm:"not null"`
	User             User           `json:"user" gorm:"foreignKey:UserID"`
	QuotaType        string         `json:"quota_type" gorm:"not null"`  // WORKSPACE, DATABASE, LOG
	TotalQuota       int64          `json:"total_quota" gorm:"not null"` // 字节
	UsedQuota        int64          `json:"used_quota" gorm:"default:0"`
	WarningThreshold float64        `json:"warning_threshold" gorm:"default:80.00"`
	Status           int8           `json:"status" gorm:"default:1"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}
