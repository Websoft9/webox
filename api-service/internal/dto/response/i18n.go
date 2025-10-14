package response

// LanguageInfo represents language information
type LanguageInfo struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	NativeName string `json:"native_name"`
	Region     string `json:"region"`
	ISOCode    string `json:"iso_code"`
	Supported  bool   `json:"supported"`
	IsDefault  bool   `json:"is_default"`
}

// SupportedLanguagesResponse represents the response for supported languages
type SupportedLanguagesResponse struct {
	Languages       []LanguageInfo `json:"languages"`
	DefaultLanguage string         `json:"default_language"`
}
