package request

import "api-service/internal/dto/common"

// GetNotificationRecordListRequest represents the request for getting notification record list
type GetNotificationRecordListRequest struct {
	common.BaseListRequest
	ChannelType string `form:"channel_type" binding:"omitempty,oneof=EMAIL WEBHOOK INTERNAL" json:"channel_type"`
	Status      string `form:"status" binding:"omitempty,oneof=PENDING SENT FAILED RETRY" json:"status"`
	Recipient   string `form:"recipient" json:"recipient"`
	SentStart   string `form:"sent_start" binding:"omitempty" json:"sent_start"`
	SentEnd     string `form:"sent_end" binding:"omitempty" json:"sent_end"`
}
