package i18n

import (
	"embed"
	"fmt"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed locales/*.yaml
var localeFS embed.FS

// Bundle holds the i18n bundle
var Bundle *i18n.Bundle

// SupportedLanguages contains all supported languages
var SupportedLanguages = []string{"en", "zh"}

// DefaultLanguage is the fallback language
const DefaultLanguage = "en"

// Init initializes the i18n bundle
func Init() error {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// Load all locale files from embedded filesystem
	for _, lang := range SupportedLanguages {
		filename := fmt.Sprintf("locales/%s.yaml", lang)
		data, err := localeFS.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read locale file %s: %w", filename, err)
		}

		_, err = Bundle.ParseMessageFileBytes(data, filename)
		if err != nil {
			return fmt.Errorf("failed to parse locale file %s: %w", filename, err)
		}
	}

	return nil
}

// GetLocalizer returns a localizer for the given language
func GetLocalizer(lang string) *i18n.Localizer {
	if lang == "" {
		lang = DefaultLanguage
	}

	// Normalize language code
	lang = normalizeLanguage(lang)

	// Validate language is supported
	if !isLanguageSupported(lang) {
		lang = DefaultLanguage
	}

	return i18n.NewLocalizer(Bundle, lang)
}

// T translates a message key with the given language
func T(key, lang string, templateData ...map[string]interface{}) string {
	localizer := GetLocalizer(lang)

	config := &i18n.LocalizeConfig{
		MessageID: key,
	}

	// Add template data if provided
	if len(templateData) > 0 && templateData[0] != nil {
		config.TemplateData = templateData[0]
	}

	message, err := localizer.Localize(config)
	if err != nil {
		// Return the key if translation fails
		return key
	}

	return message
}

// TWithPlural translates a message key with plural support
func TWithPlural(key, lang string, count int, templateData ...map[string]interface{}) string {
	localizer := GetLocalizer(lang)

	config := &i18n.LocalizeConfig{
		MessageID:   key,
		PluralCount: count,
	}

	// Add template data if provided
	if len(templateData) > 0 && templateData[0] != nil {
		config.TemplateData = templateData[0]
	}

	message, err := localizer.Localize(config)
	if err != nil {
		// Return the key if translation fails
		return key
	}

	return message
}

// normalizeLanguage normalizes language codes (e.g., "zh-CN" -> "zh")
func normalizeLanguage(lang string) string {
	if lang == "" {
		return DefaultLanguage
	}

	// Extract main language code
	parts := strings.Split(lang, "-")
	if len(parts) > 0 {
		return strings.ToLower(parts[0])
	}

	return strings.ToLower(lang)
}

// isLanguageSupported checks if a language is supported
func isLanguageSupported(lang string) bool {
	for _, supported := range SupportedLanguages {
		if supported == lang {
			return true
		}
	}
	return false
}

// GetSupportedLanguages returns all supported languages
func GetSupportedLanguages() []string {
	return SupportedLanguages
}

// DetectLanguageFromHeader detects language from Accept-Language header
func DetectLanguageFromHeader(acceptLang string) string {
	if acceptLang == "" {
		return DefaultLanguage
	}

	// Parse Accept-Language header
	// Format: "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7"
	languages := strings.Split(acceptLang, ",")

	for _, lang := range languages {
		// Remove quality factor if present
		lang = strings.Split(lang, ";")[0]
		lang = strings.TrimSpace(lang)

		// Normalize and check if supported
		normalized := normalizeLanguage(lang)
		if isLanguageSupported(normalized) {
			return normalized
		}
	}

	return DefaultLanguage
}

// MustT is like T but panics if the bundle is not initialized
func MustT(key, lang string, templateData ...map[string]interface{}) string {
	if Bundle == nil {
		panic("i18n bundle not initialized. Call i18n.Init() first")
	}
	return T(key, lang, templateData...)
}

// GetAvailableKeys returns all available message keys for debugging
func GetAvailableKeys() map[string][]string {
	if Bundle == nil {
		return nil
	}

	result := make(map[string][]string)

	for _, lang := range SupportedLanguages {
		// This is a simplified version - in real implementation,
		// you might want to iterate through actual message files
		result[lang] = []string{
			"user.not_found",
			"user.created_success",
			"auth.token_invalid",
			"common.success",
			"error.database_error",
		}
	}

	return result
}

// Language names
const (
	EnglishName = "English"
	ChineseName = "Chinese"
)

// GetLanguageInfo returns information about a language
func GetLanguageInfo(lang string) map[string]interface{} {
	lang = normalizeLanguage(lang)

	info := map[string]interface{}{
		"code":      lang,
		"supported": isLanguageSupported(lang),
		"default":   lang == DefaultLanguage,
	}

	// Add language names
	switch lang {
	case "en":
		info["name"] = EnglishName
		info["native_name"] = EnglishName
	case "zh":
		info["name"] = ChineseName
		info["native_name"] = "中文"
	default:
		info["name"] = lang
		info["native_name"] = lang
	}

	return info
}
