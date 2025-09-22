package response

import "time"

// NotificationRecordResponse represents the notification record response
type NotificationRecordResponse struct {
	ID            uint       `json:"id"`
	TemplateID    *uint      `json:"template_id"`
	ChannelType   string     `json:"channel_type"`
	Recipient     string     `json:"recipient"`
	Subject       *string    `json:"subject"`
	Content       string     `json:"content"`
	Status        string     `json:"status"`
	SentAt        *time.Time `json:"sent_at"`
	ErrorMsg      *string    `json:"error_msg"`
	RetryCount    int        `json:"retry_count"`
	ReferenceID   *string    `json:"reference_id"`
	ReferenceType *string    `json:"reference_type"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// NotificationRecordListResponse notification record list response
type NotificationRecordListResponse struct {
	Page       int                          `json:"page"`        // Current page number
	PageSize   int                          `json:"page_size"`   // Items per page
	Total      int64                        `json:"total"`       // Total records
	TotalPages int                          `json:"total_pages"` // Total pages
	Items      []NotificationRecordResponse `json:"items"`       // Notification record list
}

// NotificationRecordDetailResponse represents the detailed notification record response
type NotificationRecordDetailResponse struct {
	NotificationRecordResponse
	TemplateName *string `json:"template_name,omitempty"`
	UserName     *string `json:"user_name,omitempty"`
}
