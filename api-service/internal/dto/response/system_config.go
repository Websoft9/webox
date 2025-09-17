package response

// SystemConfigsResponse represents a single system configuration item
type SystemConfigsResponse struct {
	ID           uint   `json:"id" example:"1"`
	ConfigKey    string `json:"config_key" example:"site_name"`
	ConfigValue  string `json:"config_value" example:"My Website"`
	ConfigType   string `json:"config_type" example:"STRING"`
	Category     string `json:"category" example:"General"`
	Description  string `json:"description" example:"The name of the website"`
	IsReadonly   bool   `json:"is_readonly" example:"false"`
	IsEncrypted  bool   `json:"is_encrypted" example:"false"`
	DefaultValue string `json:"default_value" example:"Default Site Name"`
	SortOrder    int    `json:"sort_order" example:"1"`
}

// ListSystemConfigsResponse represents a paginated list of system configurations
type ListSystemConfigsResponse struct {
	Items []SystemConfigsResponse `json:"items"`
}
