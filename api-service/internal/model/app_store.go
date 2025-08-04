package model

import (
	"time"

	"gorm.io/gorm"
)

// ========================================
// 3.1 应用商店 (App Store)
// ========================================

// AppStoreCategory 应用分类表
type AppStoreCategory struct {
	ID          uint               `json:"id" gorm:"primarykey"`
	Name        string             `json:"name" gorm:"not null" binding:"required"`
	Code        string             `json:"code" gorm:"uniqueIndex;not null" binding:"required"`
	ParentID    *uint              `json:"parent_id"`
	Parent      *AppStoreCategory  `json:"parent" gorm:"foreignKey:ParentID"`
	Children    []AppStoreCategory `json:"children" gorm:"foreignKey:ParentID"`
	Icon        string             `json:"icon"`
	Description string             `json:"description" gorm:"type:text"`
	SortOrder   int                `json:"sort_order" gorm:"default:0"`
	Status      int8               `json:"status" gorm:"default:1"` // 0-禁用，1-启用
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `json:"-" gorm:"index"`

	// 关联关系
	Templates []AppStoreTemplate `json:"templates" gorm:"foreignKey:CategoryID"`
}

// AppStoreTemplate 应用模板表
type AppStoreTemplate struct {
	ID              uint             `json:"id" gorm:"primarykey"`
	Name            string           `json:"name" gorm:"not null" binding:"required"`
	Code            string           `json:"code" gorm:"uniqueIndex;not null" binding:"required"`
	CategoryID      uint             `json:"category_id" gorm:"not null"`
	Category        AppStoreCategory `json:"category" gorm:"foreignKey:CategoryID"`
	Version         string           `json:"version" gorm:"not null"`
	Icon            string           `json:"icon"`
	Description     string           `json:"description" gorm:"type:text"`
	OfficialURL     string           `json:"official_url"`
	SourceURL       string           `json:"source_url"`
	ComposeTemplate string           `json:"compose_template" gorm:"type:text;not null"`
	DownloadCount   int              `json:"download_count" gorm:"default:0"`
	StarCount       int              `json:"star_count" gorm:"default:0"`
	Rating          float64          `json:"rating" gorm:"type:decimal(3,2);default:0.00"`
	IsOfficial      int8             `json:"is_official" gorm:"default:0"`
	IsFeatured      int8             `json:"is_featured" gorm:"default:0"`
	Status          int8             `json:"status" gorm:"default:1"` // 0-下架，1-上架
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `json:"-" gorm:"index"`

	// 关联关系
	Reviews     []AppStoreReview   `json:"reviews" gorm:"foreignKey:TemplateID"`
	Favorites   []AppStoreFavorite `json:"favorites" gorm:"foreignKey:TemplateID"`
	Stars       []AppStoreStar     `json:"stars" gorm:"foreignKey:TemplateID"`
	Downloads   []AppStoreDownload `json:"downloads" gorm:"foreignKey:TemplateID"`
	Deployments []AppDeployment    `json:"deployments" gorm:"foreignKey:TemplateID"`
}

// AppStoreWishlist 应用心愿单表
type AppStoreWishlist struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	Name         string         `json:"name" gorm:"not null" binding:"required"`
	Version      string         `json:"version"`
	SourceURL    string         `json:"source_url"`
	Description  string         `json:"description" gorm:"type:text"`
	RewardAmount float64        `json:"reward_amount" gorm:"type:decimal(10,2);default:0.00"`
	Priority     int8           `json:"priority" gorm:"default:3"`     // 1-高，2-中，3-低
	Status       string         `json:"status" gorm:"default:PENDING"` // PENDING, IN_PROGRESS, COMPLETED, EXPIRED
	ViewCount    int            `json:"view_count" gorm:"default:0"`
	LikeCount    int            `json:"like_count" gorm:"default:0"`
	VoteCount    int            `json:"vote_count" gorm:"default:0"`
	CommentCount int            `json:"comment_count" gorm:"default:0"`
	SubmitterID  uint           `json:"submitter_id" gorm:"not null"`
	Submitter    User           `json:"submitter" gorm:"foreignKey:SubmitterID"`
	CompletedAt  *time.Time     `json:"completed_at"`
	ExpiresAt    *time.Time     `json:"expires_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Comments []AppStoreWishlistComment `json:"comments" gorm:"foreignKey:WishlistID"`
	Votes    []AppStoreWishlistVote    `json:"votes" gorm:"foreignKey:WishlistID"`
	Likes    []AppStoreWishlistLike    `json:"likes" gorm:"foreignKey:WishlistID"`
}

// AppStoreReview 应用评价表
type AppStoreReview struct {
	ID           uint             `json:"id" gorm:"primarykey"`
	TemplateID   uint             `json:"template_id" gorm:"not null"`
	Template     AppStoreTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	UserID       uint             `json:"user_id" gorm:"not null"`
	User         User             `json:"user" gorm:"foreignKey:UserID"`
	Rating       int8             `json:"rating" gorm:"not null"` // 1-5分
	Content      string           `json:"content" gorm:"type:text"`
	Tags         string           `json:"tags" gorm:"type:json"`
	IsHelpful    int8             `json:"is_helpful" gorm:"default:0"`
	HelpfulCount int              `json:"helpful_count" gorm:"default:0"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `json:"-" gorm:"index"`
}

// AppStoreFavorite 应用收藏表
type AppStoreFavorite struct {
	ID         uint             `json:"id" gorm:"primarykey"`
	TemplateID uint             `json:"template_id" gorm:"not null"`
	Template   AppStoreTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	UserID     uint             `json:"user_id" gorm:"not null"`
	User       User             `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt  time.Time        `json:"created_at"`
}

// AppStoreStar 应用点赞表
type AppStoreStar struct {
	ID         uint             `json:"id" gorm:"primarykey"`
	TemplateID uint             `json:"template_id" gorm:"not null"`
	Template   AppStoreTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	UserID     uint             `json:"user_id" gorm:"not null"`
	User       User             `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt  time.Time        `json:"created_at"`
}

// AppStoreReport 应用举报表
type AppStoreReport struct {
	ID         uint             `json:"id" gorm:"primarykey"`
	TemplateID uint             `json:"template_id" gorm:"not null"`
	Template   AppStoreTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	UserID     uint             `json:"user_id" gorm:"not null"`
	User       User             `json:"user" gorm:"foreignKey:UserID"`
	Reason     string           `json:"reason" gorm:"not null"` // SPAM, INAPPROPRIATE, COPYRIGHT
	Detail     string           `json:"detail" gorm:"type:text"`
	Status     string           `json:"status" gorm:"default:PENDING"` // PENDING, PROCESSING, RESOLVED, REJECTED
	HandledBy  *uint            `json:"handled_by"`
	Handler    *User            `json:"handler" gorm:"foreignKey:HandledBy"`
	HandledAt  *time.Time       `json:"handled_at"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	DeletedAt  gorm.DeletedAt   `json:"-" gorm:"index"`
}

// AppStoreDownload 应用下载记录表
type AppStoreDownload struct {
	ID         uint             `json:"id" gorm:"primarykey"`
	TemplateID uint             `json:"template_id" gorm:"not null"`
	Template   AppStoreTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	UserID     uint             `json:"user_id" gorm:"not null"`
	User       User             `json:"user" gorm:"foreignKey:UserID"`
	IPAddress  string           `json:"ip_address"`
	UserAgent  string           `json:"user_agent"`
	CreatedAt  time.Time        `json:"created_at"`
}

// AppStoreWishlistComment 应用心愿单评论表
type AppStoreWishlistComment struct {
	ID         uint                      `json:"id" gorm:"primarykey"`
	WishlistID uint                      `json:"wishlist_id" gorm:"not null"`
	Wishlist   AppStoreWishlist          `json:"wishlist" gorm:"foreignKey:WishlistID"`
	UserID     uint                      `json:"user_id" gorm:"not null"`
	User       User                      `json:"user" gorm:"foreignKey:UserID"`
	ParentID   *uint                     `json:"parent_id"`
	Parent     *AppStoreWishlistComment  `json:"parent" gorm:"foreignKey:ParentID"`
	Children   []AppStoreWishlistComment `json:"children" gorm:"foreignKey:ParentID"`
	Content    string                    `json:"content" gorm:"type:text;not null"`
	LikeCount  int                       `json:"like_count" gorm:"default:0"`
	CreatedAt  time.Time                 `json:"created_at"`
	UpdatedAt  time.Time                 `json:"updated_at"`
	DeletedAt  gorm.DeletedAt            `json:"-" gorm:"index"`
}

// AppStoreWishlistVote 应用心愿单投票表
type AppStoreWishlistVote struct {
	ID         uint             `json:"id" gorm:"primarykey"`
	WishlistID uint             `json:"wishlist_id" gorm:"not null"`
	Wishlist   AppStoreWishlist `json:"wishlist" gorm:"foreignKey:WishlistID"`
	UserID     uint             `json:"user_id" gorm:"not null"`
	User       User             `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt  time.Time        `json:"created_at"`
}

// AppStoreWishlistLike 应用心愿单点赞表
type AppStoreWishlistLike struct {
	ID         uint             `json:"id" gorm:"primarykey"`
	WishlistID uint             `json:"wishlist_id" gorm:"not null"`
	Wishlist   AppStoreWishlist `json:"wishlist" gorm:"foreignKey:WishlistID"`
	UserID     uint             `json:"user_id" gorm:"not null"`
	User       User             `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt  time.Time        `json:"created_at"`
}

// AppStoreWishlistReport 应用心愿单举报表
type AppStoreWishlistReport struct {
	ID         uint             `json:"id" gorm:"primarykey"`
	WishlistID uint             `json:"wishlist_id" gorm:"not null"`
	Wishlist   AppStoreWishlist `json:"wishlist" gorm:"foreignKey:WishlistID"`
	UserID     uint             `json:"user_id" gorm:"not null"`
	User       User             `json:"user" gorm:"foreignKey:UserID"`
	Reason     string           `json:"reason" gorm:"not null"` // SPAM, INAPPROPRIATE, DUPLICATE
	Detail     string           `json:"detail" gorm:"type:text"`
	Status     string           `json:"status" gorm:"default:PENDING"` // PENDING, PROCESSING, RESOLVED, REJECTED
	HandledBy  *uint            `json:"handled_by"`
	Handler    *User            `json:"handler" gorm:"foreignKey:HandledBy"`
	HandledAt  *time.Time       `json:"handled_at"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	DeletedAt  gorm.DeletedAt   `json:"-" gorm:"index"`
}

// AppDeployment 应用部署记录表
type AppDeployment struct {
	ID            uint              `json:"id" gorm:"primarykey"`
	DeploymentID  string            `json:"deployment_id" gorm:"uniqueIndex;not null"`
	TemplateID    *uint             `json:"template_id"`
	Template      *AppStoreTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	AppInstanceID *uint             `json:"app_instance_id"`
	AppInstance   *AppInstance      `json:"app_instance" gorm:"foreignKey:AppInstanceID"`
	ServerID      uint              `json:"server_id" gorm:"not null"`
	Server        Server            `json:"server" gorm:"foreignKey:ServerID"`
	Status        string            `json:"status" gorm:"default:PENDING"` // PENDING, RUNNING, SUCCESS, FAILED, CANCELED
	Progress      int8              `json:"progress" gorm:"default:0"`
	EstimatedTime int               `json:"estimated_time" gorm:"default:0"` // 秒
	StartTime     *time.Time        `json:"start_time"`
	EndTime       *time.Time        `json:"end_time"`
	ErrorMessage  string            `json:"error_message" gorm:"type:text"`
	DeploymentLog string            `json:"deployment_log" gorm:"type:text"`
	ConfigData    string            `json:"config_data" gorm:"type:json"`
	OwnerID       uint              `json:"owner_id" gorm:"not null"`
	Owner         User              `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `json:"-" gorm:"index"`
}

// AppShortcut 应用快捷导航表
type AppShortcut struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	AppInstanceID uint           `json:"app_instance_id" gorm:"not null"`
	AppInstance   AppInstance    `json:"app_instance" gorm:"foreignKey:AppInstanceID"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Icon          string         `json:"icon"`
	SortOrder     int            `json:"sort_order" gorm:"default:0"`
	AccessCount   int            `json:"access_count" gorm:"default:0"`
	LastAccessed  *time.Time     `json:"last_accessed"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
