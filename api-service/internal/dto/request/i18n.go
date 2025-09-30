package request

// SwitchLanguageRequest represents a request to switch user language
type SwitchLanguageRequest struct {
	Language string `json:"language" binding:"required" validate:"required,min=2,max=10"`
}
