package response

// NotificationTemplateResponse represents a notification template in response
type NotificationTemplateResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	TemplateType string  `json:"template_type"`
	Subject      *string `json:"subject"`
	IsSystem     int     `json:"is_system"`
	Status       int     `json:"status"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// NotificationTemplateDetailResponse represents detailed notification template information
type NotificationTemplateDetailResponse struct {
	NotificationTemplateResponse
	Content string `json:"content"`
}
