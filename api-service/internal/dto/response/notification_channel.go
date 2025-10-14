package response

import (
	"api-service/internal/model"
	"time"
)

// NotificationChannelResponse represents the basic notification channel response
type NotificationChannelResponse struct {
	ID            uint       `json:"id"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	ChannelType   string     `json:"channel_type"`
	ChannelConfig model.JSON `json:"channel_config"`
	OwnerID       uint       `json:"owner_id"`
	Status        int8       `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
