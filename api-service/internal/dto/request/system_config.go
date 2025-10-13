package request

// ListSystemConfigsRequest request for querying system configurations
type ListSystemConfigsRequest struct {
	Category string `form:"category"`
	Keyword  string `form:"keyword"`
}

// UpdateSystemConfigRequest request for updating a single system configuration
type UpdateSystemConfigRequest struct {
	ConfigKey   string  `json:"config_key" validate:"required"`
	ConfigValue string  `json:"config_value" validate:"required"`
	ConfigType  string  `json:"config_type" validate:"required,oneof=STRING BOOLEAN INTEGER JSON FLOAT"`
	Category    string  `json:"category" validate:"required"`
	Description *string `json:"description,omitempty"`
	SortOrder   *int    `json:"sort_order,omitempty"`
}

// BatchUpdateSystemConfigsRequest request for batch updating system configurations
type BatchUpdateSystemConfigsRequest struct {
	Configs []UpdateSystemConfigRequest `json:"configs" validate:"required,dive"`
}

// TestSMTPRequest request for testing SMTP configuration
type TestSMTPRequest struct {
	TestEmail string `json:"test_email" validate:"required,email"`
}

// ListSystemConfigsFilter filter conditions for listing system configurations
type ListSystemConfigsFilter struct {
	Category string `json:"category"`
	Keyword  string `json:"keyword"`
}
