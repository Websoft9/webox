package response

import (
	"api-service/internal/model"
	"time"
)

// NotificationChannelResponse represents the basic notification channel response
type NotificationChannelResponse struct {
	ID            uint                    `json:"id"`
	Code          string                  `json:"code"`
	Name          string                  `json:"name"`
	Description   *string                 `json:"description,omitempty"`
	ChannelType   string                  `json:"channel_type"`
	ChannelConfig model.JSONChannelConfig `json:"channel_config"`
	OwnerID       uint                    `json:"owner_id"`
	Status        int8                    `json:"status"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

// NotificationChannelListResponse represents the paginated list response
type NotificationChannelListResponse struct {
	Page       int                           `json:"page"`
	PageSize   int                           `json:"page_size"`
	Total      int64                         `json:"total"`
	TotalPages int                           `json:"total_pages"`
	Items      []NotificationChannelResponse `json:"items"`
}

// NotificationChannelDetailResponse represents the detailed notification channel response
type NotificationChannelDetailResponse struct {
	NotificationChannelResponse
	OwnerName *string `json:"owner_name,omitempty"`
}

// TestChannelResponse represents the response for channel testing
type TestChannelResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
