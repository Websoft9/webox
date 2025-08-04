package model

import (
	"time"

	"gorm.io/gorm"
)

// ========================================
// 3.4 安全管控 (Security Management)
// ========================================

// SSLCertificate SSL证书表
type SSLCertificate struct {
	ID     uint   `json:"id" gorm:"primarykey"`
	Name   string `json:"name" gorm:"not null" binding:"required"`
	Domain string `json:"domain" gorm:"not null" binding:"required"`
	// CertificateType 证书类型: LETS_ENCRYPT, COMMERCIAL, SELF_SIGNED
	CertificateType  string         `json:"certificate_type" gorm:"default:LETS_ENCRYPT"`
	CertificateData  string         `json:"certificate_data" gorm:"type:text;not null"`
	PrivateKeyData   string         `json:"private_key_data" gorm:"type:text;not null"`
	CertificateChain string         `json:"certificate_chain" gorm:"type:text"`
	Issuer           string         `json:"issuer"`
	Subject          string         `json:"subject"`
	SerialNumber     string         `json:"serial_number"`
	NotBefore        *time.Time     `json:"not_before"`
	NotAfter         *time.Time     `json:"not_after"`
	AutoRenew        int8           `json:"auto_renew" gorm:"default:0"`
	Status           string         `json:"status" gorm:"default:PENDING"` // PENDING, VALID, EXPIRED, REVOKED
	OwnerID          uint           `json:"owner_id" gorm:"not null"`
	Owner            User           `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// SecretKey 密钥管理表
type SecretKey struct {
	ID               uint           `json:"id" gorm:"primarykey"`
	Name             string         `json:"name" gorm:"not null" binding:"required"`
	KeyType          string         `json:"key_type" gorm:"not null"` // API_KEY, DATABASE, SSH, CERTIFICATE, CUSTOM
	EncryptedValue   string         `json:"encrypted_value" gorm:"type:text;not null"`
	Description      string         `json:"description" gorm:"type:text"`
	CustomFields     string         `json:"custom_fields" gorm:"type:json"`
	AuthorizedUsers  string         `json:"authorized_users" gorm:"type:json"`
	AuthorizedGroups string         `json:"authorized_groups" gorm:"type:json"`
	ExpiresAt        *time.Time     `json:"expires_at"`
	OwnerID          uint           `json:"owner_id" gorm:"not null"`
	Owner            User           `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// AppGateway 应用网关表
type AppGateway struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	Name            string         `json:"name" gorm:"not null" binding:"required"`
	ServerID        uint           `json:"server_id" gorm:"not null"`
	Server          Server         `json:"server" gorm:"foreignKey:ServerID"`
	ContainerID     string         `json:"container_id"`
	Description     string         `json:"description" gorm:"type:text"`
	Status          string         `json:"status" gorm:"default:DEFAULT"` // DEFAULT, RUNNING, STOPPED, ERROR
	StartedAt       *time.Time     `json:"started_at"`
	StoppedAt       *time.Time     `json:"stopped_at"`
	ResourceGroupID *uint          `json:"resource_group_id"`
	ResourceGroup   *ResourceGroup `json:"resource_group" gorm:"foreignKey:ResourceGroupID"`
	OwnerID         uint           `json:"owner_id" gorm:"not null"`
	Owner           User           `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Publishes   []AppGatewayPublish    `json:"publishes" gorm:"foreignKey:AppGatewayID"`
	AccessRules []AppGatewayAccessRule `json:"access_rules" gorm:"foreignKey:GatewayID"`
}

// AppGatewayPublish 应用网关发布表
type AppGatewayPublish struct {
	ID                 uint           `json:"id" gorm:"primarykey"`
	AppInstanceID      uint           `json:"app_instance_id" gorm:"not null"`
	AppInstance        AppInstance    `json:"app_instance" gorm:"foreignKey:AppInstanceID"`
	AppGatewayID       uint           `json:"app_gateway_id" gorm:"not null"`
	AppGateway         AppGateway     `json:"app_gateway" gorm:"foreignKey:AppGatewayID"`
	ServiceDomain      string         `json:"service_domain" gorm:"not null"`
	ServicePort        int            `json:"service_port" gorm:"default:8080"`
	AlertRuleID        *uint          `json:"alert_rule_id"`
	AlertRule          *AlertRule     `json:"alert_rule" gorm:"foreignKey:AlertRuleID"`
	LimitRules         string         `json:"limit_rules" gorm:"type:text"`
	HealthCheckEnabled int8           `json:"health_check_enabled" gorm:"default:1"`
	AuditLogEnabled    int8           `json:"audit_log_enabled" gorm:"default:1"`
	OwnerID            uint           `json:"owner_id" gorm:"not null"`
	Owner              User           `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
}

// AppGatewayAccessRule 应用网关访问控制规则表
type AppGatewayAccessRule struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	GatewayID  uint           `json:"gateway_id" gorm:"not null"`
	Gateway    AppGateway     `json:"gateway" gorm:"foreignKey:GatewayID"`
	RuleName   string         `json:"rule_name" gorm:"not null" binding:"required"`
	RuleType   string         `json:"rule_type" gorm:"not null"` // IP_WHITELIST, IP_BLACKLIST, RATE_LIMIT
	LimitCount int            `json:"limit_count" gorm:"not null"`
	TimeWindow int            `json:"time_window" gorm:"not null"` // 秒
	TargetPath string         `json:"target_path"`
	TargetIP   string         `json:"target_ip"`
	Action     string         `json:"action" gorm:"default:BLOCK"` // BLOCK, ALLOW, REDIRECT
	Status     int8           `json:"status" gorm:"default:1"`     // 0-禁用，1-启用
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// AuditLog 审计日志表
type AuditLog struct {
	ID             uint      `json:"id" gorm:"primarykey"`
	UserID         *uint     `json:"user_id"`
	User           *User     `json:"user" gorm:"foreignKey:UserID"`
	Username       string    `json:"username"`
	Action         string    `json:"action" gorm:"not null"`
	Module         string    `json:"module" gorm:"not null"`
	ResourceType   string    `json:"resource_type"`
	ResourceID     *uint     `json:"resource_id"`
	ResourceName   string    `json:"resource_name"`
	Description    string    `json:"description" gorm:"type:text"`
	IPAddress      string    `json:"ip_address"`
	UserAgent      string    `json:"user_agent"`
	RequestMethod  string    `json:"request_method"`
	RequestURL     string    `json:"request_url"`
	RequestParams  string    `json:"request_params" gorm:"type:json"`
	ResponseStatus *int      `json:"response_status"`
	ResponseTime   *int      `json:"response_time"` // 毫秒
	Success        int8      `json:"success" gorm:"default:1"`
	ErrorMessage   string    `json:"error_message" gorm:"type:text"`
	CreatedAt      time.Time `json:"created_at"`
}

// Gateway 网关表 (兼容现有代码)
type Gateway struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	Name          string         `json:"name" gorm:"not null"`
	Domain        string         `json:"domain" gorm:"uniqueIndex;not null"`
	ApplicationID uint           `json:"application_id"`
	Application   Application    `json:"application" gorm:"foreignKey:ApplicationID"`
	SSLEnabled    bool           `json:"ssl_enabled" gorm:"default:false"`
	SSLCert       string         `json:"ssl_cert"`
	AccessRules   string         `json:"access_rules" gorm:"type:text"`
	Status        string         `json:"status" gorm:"default:active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
