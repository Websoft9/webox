package request

import "api-service/internal/dto"

// GetNotificationRecordListRequest represents the request for getting notification record list
type GetNotificationRecordListRequest struct {
	dto.BaseListRequest
	ChannelType string `form:"channel_type" binding:"omitempty,oneof=EMAIL WEBHOOK INTERNAL" json:"channel_type"`
	Status      string `form:"status" binding:"omitempty,oneof=PENDING SENT FAILED RETRY" json:"status"`
	Recipient   string `form:"recipient" json:"recipient"`
	SentStart   string `form:"sent_start" binding:"omitempty,datetime=2006-01-02 15:04:05" json:"sent_start"`
	SentEnd     string `form:"sent_end" binding:"omitempty,datetime=2006-01-02 15:04:05" json:"sent_end"`
}
