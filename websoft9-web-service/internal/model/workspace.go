package model

import (
	"time"

	"gorm.io/gorm"
)

// ========================================
// 3.2 工作空间 (Workspace)
// ========================================

// UserFile 用户文件表
type UserFile struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	Name          string         `json:"name" gorm:"not null" binding:"required"`
	Path          string         `json:"path" gorm:"not null"`
	Type          string         `json:"type" gorm:"not null;default:FILE"` // FILE, DIRECTORY
	Size          int64          `json:"size" gorm:"default:0"`             // 字节
	MimeType      string         `json:"mime_type"`
	DownloadCount int            `json:"download_count" gorm:"default:0"`
	ParentID      *uint          `json:"parent_id"`
	Parent        *UserFile      `json:"parent" gorm:"foreignKey:ParentID"`
	Children      []UserFile     `json:"children" gorm:"foreignKey:ParentID"`
	StoragePath   string         `json:"storage_path"`
	Checksum      string         `json:"checksum"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// Workflow 工作流表
type Workflow struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"not null" binding:"required"`
	Code        string         `json:"code" gorm:"not null" binding:"required"`
	Description string         `json:"description" gorm:"type:text"`
	Definition  string         `json:"definition" gorm:"type:json;not null"` // JSON格式的工作流定义
	Status      string         `json:"status" gorm:"default:DRAFT"`          // DRAFT, ACTIVE, INACTIVE, ARCHIVED
	OwnerID     uint           `json:"owner_id" gorm:"not null"`
	Owner       User           `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Tasks []WorkflowTask `json:"tasks" gorm:"foreignKey:WorkflowID"`
}

// WorkflowTask 工作流任务表
type WorkflowTask struct {
	ID             uint           `json:"id" gorm:"primarykey"`
	Name           string         `json:"name" gorm:"not null" binding:"required"`
	WorkflowID     uint           `json:"workflow_id" gorm:"not null"`
	Workflow       Workflow       `json:"workflow" gorm:"foreignKey:WorkflowID"`
	ScheduleType   string         `json:"schedule_type" gorm:"default:MANUAL"` // MANUAL, SCHEDULE, TRIGGER
	CronExpression string         `json:"cron_expression"`
	Status         string         `json:"status" gorm:"default:DEFAULT"` // DEFAULT, ONLINE, OFFLINE
	NextRunAt      *time.Time     `json:"next_run_at"`
	LastRunAt      *time.Time     `json:"last_run_at"`
	RunCount       int            `json:"run_count" gorm:"default:0"`
	SuccessCount   int            `json:"success_count" gorm:"default:0"`
	FailureCount   int            `json:"failure_count" gorm:"default:0"`
	OwnerID        uint           `json:"owner_id" gorm:"not null"`
	Owner          User           `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Executions []WorkflowExecution `json:"executions" gorm:"foreignKey:TaskID"`
}

// WorkflowExecution 工作流执行历史表
type WorkflowExecution struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	TaskID       uint           `json:"task_id" gorm:"not null"`
	Task         WorkflowTask   `json:"task" gorm:"foreignKey:TaskID"`
	ExecutionID  string         `json:"execution_id" gorm:"uniqueIndex;not null"`
	Status       string         `json:"status" gorm:"default:PENDING"`      // PENDING, RUNNING, SUCCESS, FAILURE, STOPPED
	TriggerType  string         `json:"trigger_type" gorm:"default:MANUAL"` // MANUAL, SCHEDULE, TRIGGER
	TriggerBy    *uint          `json:"trigger_by"`
	Trigger      *User          `json:"trigger" gorm:"foreignKey:TriggerBy"`
	StartTime    *time.Time     `json:"start_time"`
	EndTime      *time.Time     `json:"end_time"`
	Duration     int            `json:"duration" gorm:"default:0"` // 秒
	ErrorMessage string         `json:"error_message" gorm:"type:text"`
	ExecutionLog string         `json:"execution_log" gorm:"type:text"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
