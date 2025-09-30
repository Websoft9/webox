package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// Language constants
const (
	LangEnUS = "en-US"
	LangZhCN = "zh-CN"
)

// Bundle holds the i18n bundle
var Bundle *i18n.Bundle

// SupportedLanguages contains all supported languages (will be initialized from config)
var SupportedLanguages []string

// DefaultLanguage is the fallback language (will be initialized from config)
var DefaultLanguage string

// getProjectRoot returns the absolute path to the project root directory
// It searches for go.mod file starting from the current file's directory
func getProjectRoot() (string, error) {
	// Get the directory of the current source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	// Start from the directory containing this file
	dir := filepath.Dir(filename)

	// Walk up the directory tree to find go.mod
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			// Found go.mod, this is the project root
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding go.mod
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find project root (go.mod not found)")
}

// Init initializes the i18n bundle with default configuration
func Init() error {
	//TODO: 从`config.yaml`配置文件初始化
	return InitWithConfig(LangEnUS, []string{LangEnUS, LangZhCN})
}

// InitWithConfig initializes the i18n bundle with custom configuration
// This function loads language files from configs/lang/ directory
func InitWithConfig(defaultLang string, supportedLangs []string) error {
	// Set configuration
	DefaultLanguage = defaultLang
	SupportedLanguages = supportedLangs

	Bundle = i18n.NewBundle(language.AmericanEnglish)
	Bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// Get project root directory
	projectRoot, err := getProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to get project root: %w", err)
	}

	// Load all locale files from configs/lang/ directory using absolute path
	for _, lang := range SupportedLanguages {
		filename := filepath.Join(projectRoot, "configs", "lang", fmt.Sprintf("%s.yaml", lang))
		if err := loadLanguageFile(filename); err != nil {
			return fmt.Errorf("failed to load language file %s: %w", filename, err)
		}
	}

	return nil
}

// loadLanguageFile loads a language file from the filesystem
func loadLanguageFile(filename string) error {
	// Clean the file path to prevent path traversal
	cleanPath := filepath.Clean(filename)

	// Validate that the path is absolute and doesn't contain path traversal attempts
	if !filepath.IsAbs(cleanPath) {
		return fmt.Errorf("language file path must be absolute: %s", cleanPath)
	}

	// Additional security check: ensure the path doesn't contain ".."
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid language file path (contains ..): %s", cleanPath)
	}

	// Check if file exists
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return fmt.Errorf("language file does not exist: %s", cleanPath)
	}

	// Read file content with cleaned path
	// #nosec G304 - Path is validated and cleaned above, and only loads from trusted configs/lang directory
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to read language file: %w", err)
	}

	// Parse the file
	_, err = Bundle.ParseMessageFileBytes(data, filepath.Base(cleanPath))
	if err != nil {
		return fmt.Errorf("failed to parse language file: %w", err)
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

// NormalizeLanguage normalizes language codes to standard format (exported version)
func NormalizeLanguage(lang string) string {
	return normalizeLanguage(lang)
}

// normalizeLanguage normalizes language codes to standard format
func normalizeLanguage(lang string) string {
	if lang == "" {
		return DefaultLanguage
	}

	// Convert to lowercase for comparison
	lang = strings.ToLower(strings.TrimSpace(lang))

	// Handle common variations and normalize to standard format
	switch {
	case lang == "en" || lang == "en-us" || strings.HasPrefix(lang, "en-"):
		return LangEnUS
	case lang == "zh" || lang == "zh-cn" || lang == "zh-hans" ||
		lang == "zh-tw" || lang == "zh-hant" || strings.HasPrefix(lang, "zh-"):
		return LangZhCN
	default:
		// Return as-is for other languages, but check if it's in our supported list
		return lang
	}
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
	case LangEnUS:
		info["name"] = EnglishName
		info["native_name"] = EnglishName
		info["region"] = "United States"
		info["iso_code"] = LangEnUS
	case LangZhCN:
		info["name"] = ChineseName
		info["native_name"] = "中文"
		info["region"] = "China"
		info["iso_code"] = LangZhCN
	default:
		info["name"] = lang
		info["native_name"] = lang
		info["iso_code"] = lang
	}

	return info
}

// globalInstance holds the global I18n instance
var globalInstance *I18n

// GetInstance returns the global I18n instance
func GetInstance() *I18n {
	if globalInstance == nil {
		globalInstance = NewI18n()
	}
	return globalInstance
}

// I18n provides internationalization functionality
type I18n struct {
	defaultLang string
}

// NewI18n creates a new I18n instance
func NewI18n() *I18n {
	return &I18n{
		defaultLang: DefaultLanguage,
	}
}
