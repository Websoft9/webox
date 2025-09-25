package email

import (
	"api-service/internal/config"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"reflect"
	"regexp"
	"strings"
)

// EmailService defines the interface for email operations
// It provides methods for sending different types of emails with internationalization support
type EmailService interface {
	// SendVerificationEmail sends an email verification message to the user
	// Parameters:
	//   - ctx: context for request tracing and cancellation
	//   - expires: Expiry time in minutes
	//   - to: recipient email address
	//   - token: verification token to be included in the email
	//   - baseURL: base URL for constructing verification links
	//   - lang: language code for internationalization
	SendVerificationEmail(ctx context.Context, expires int, to, token, baseURL, lang string) error

	// SendPasswordResetEmail sends a password reset email to the user
	// Parameters are similar to SendVerificationEmail but for password reset functionality
	SendPasswordResetEmail(ctx context.Context, expires int, to, token, baseURL, lang string) error

	// SendEmail sends a generic email with custom subject and body
	// This is the core method used by other specialized email methods
	SendEmail(ctx context.Context, to, subject, body string) error
}

// EmailTemplateData holds the data structure for email templates
// It contains URLs and other dynamic content that will be injected into email templates
type EmailTemplateData struct {
	VerifyURL string // URL for email verification
	ResetURL  string // URL for password reset
	Expires   int    // Expiry time in minutes
}

// emailService implements the EmailService interface
// It handles SMTP configuration, template rendering, and email delivery
type emailService struct {
	config *config.EmailConfigMain // SMTP configuration settings
	logger logger.Logger           // Logger for tracking email operations
}

// NewEmailService creates a new instance of the email service
// It initializes the service with the provided configuration, logger, and i18n support
// Returns an EmailService interface implementation
func NewEmailService(cfg *config.Config, logger logger.Logger) EmailService {
	return &emailService{
		config: &cfg.Email,
		logger: logger,
	}
}

// SendVerificationEmail sends an email verification message to the specified recipient
// It constructs a verification URL using the provided token and base URL,
// then renders the email template with internationalized content
func (s *emailService) SendVerificationEmail(ctx context.Context, expires int, to, token, baseURL, lang string) error {
	// Construct the verification URL with the token parameter for API endpoint
	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", baseURL, token)

	// Get localized email subject and template content
	subject := i18n.T("email.verification_subject", lang)
	templateContent := i18n.T("email.verification_template", lang)

	// Prepare template data with the verification URL
	templateData := EmailTemplateData{
		VerifyURL: verifyURL,
		Expires:   expires,
	}

	// Render the email template with the provided data
	body := s.renderTemplate(templateContent, templateData)

	// Send the rendered email
	return s.SendEmail(ctx, to, subject, body)
}

// SendPasswordResetEmail sends a password reset email to the specified recipient
// It constructs a reset URL using the provided token and base URL,
// then renders the email template with internationalized content
func (s *emailService) SendPasswordResetEmail(ctx context.Context, expires int, to, token, baseURL, lang string) error {
	// Construct the password reset URL with the token parameter for API endpoint
	resetURL := fmt.Sprintf("%s/api/v1/auth/reset-password?token=%s", baseURL, token)

	// Get localized email subject and template content
	subject := i18n.T("email.password_reset_subject", lang)
	templateContent := i18n.T("email.password_reset_template", lang)

	// Prepare template data with the reset URL
	templateData := EmailTemplateData{
		ResetURL: resetURL,
		Expires:  expires,
	}

	// Render the email template with the provided data
	body := s.renderTemplate(templateContent, templateData)

	// Send the rendered email
	return s.SendEmail(ctx, to, subject, body)
}

// renderTemplate renders an email template with the provided data
// It uses custom variable replacement to support ${variable} format
// Returns the rendered HTML content as a string
func (s *emailService) renderTemplate(templateContent string, data EmailTemplateData) string {
	// Define regex pattern for ${variable} format
	re := regexp.MustCompile(`\$\{([^}]+)\}`)

	// Replace all ${variable} placeholders with corresponding data values
	result := re.ReplaceAllStringFunc(templateContent, func(match string) string {
		// Extract variable name from ${varName}
		varName := strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")

		// Get field value from data struct using reflection
		value := s.getFieldValue(data, varName)
		if value == "" {
			// Return original placeholder if field not found
			s.logger.Warn("Template variable not found", logger.String("variable", varName))
			return match
		}

		return value
	})

	return result
}

// getFieldValue extracts field value from struct using reflection
// Supports case-insensitive field matching for template variables
func (s *emailService) getFieldValue(data EmailTemplateData, fieldName string) string {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	// Clean up field name by trimming spaces and other whitespace characters
	fieldName = strings.TrimSpace(fieldName)

	// Iterate through struct fields to find matching field name
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		// Case-insensitive field name matching
		if strings.EqualFold(field.Name, fieldName) {
			fieldValue := v.Field(i)

			// Convert field value to string based on its type
			switch fieldValue.Kind() {
			case reflect.String:
				return fieldValue.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return fmt.Sprintf("%d", fieldValue.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return fmt.Sprintf("%d", fieldValue.Uint())
			case reflect.Float32, reflect.Float64:
				return fmt.Sprintf("%g", fieldValue.Float())
			case reflect.Bool:
				return fmt.Sprintf("%t", fieldValue.Bool())
			default:
				// For other types, use the default string representation
				return fmt.Sprintf("%v", fieldValue.Interface())
			}
		}
	}

	return ""
}

// SendEmail sends an email with the specified subject and body to the recipient
// It handles both TLS and non-TLS SMTP connections based on configuration
// This is the core method that performs the actual email delivery
func (s *emailService) SendEmail(ctx context.Context, to, subject, body string) error {
	s.logger.InfoContext(ctx, "Sending email",
		logger.String("to", to),
		logger.String("subject", subject))

	// Build the email message with proper headers and formatting
	msg := s.buildMessage(s.config.SMTP.From, to, subject, body)

	// Create SMTP authentication using configured credentials
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)

	// Construct the SMTP server address
	addr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)

	// Send email using appropriate method based on TLS configuration
	var err error
	if s.config.SMTP.UseTLS {
		// Use TLS connection for secure email delivery
		err = s.sendWithTLS(addr, auth, s.config.SMTP.From, []string{to}, msg)
	} else {
		// Use standard SMTP connection
		err = smtp.SendMail(addr, auth, s.config.SMTP.From, []string{to}, msg)
	}

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to send email",
			logger.String("to", to),
			logger.ErrorField(err))
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.InfoContext(ctx, "Email sent successfully", logger.String("to", to))
	return nil
}

// sendWithTLS sends email using a TLS-encrypted SMTP connection
// This method provides secure email delivery by establishing a TLS connection
// before sending the email content
func (s *emailService) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Create TLS configuration with security settings
	tlsConfig := &tls.Config{
		ServerName: s.config.SMTP.Host,
		MinVersion: tls.VersionTLS12, // Enforce minimum TLS version 1.2 for security
	}

	// Establish TLS connection to the SMTP server
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create SMTP client using the TLS connection
	client, err := smtp.NewClient(conn, s.config.SMTP.Host)
	if err != nil {
		return err
	}
	defer func() {
		if quitErr := client.Quit(); quitErr != nil {
			// Log the error but don't return it as it's in defer
			// The main operation has already completed
			_ = quitErr // Explicitly ignore the error
		}
	}()

	// Authenticate with the SMTP server if authentication is provided
	if auth != nil {
		if authErr := client.Auth(auth); authErr != nil {
			return authErr
		}
	}

	// Set the sender address
	if mailErr := client.Mail(from); mailErr != nil {
		return mailErr
	}

	// Set recipient addresses (can handle multiple recipients)
	for _, addr := range to {
		if rcptErr := client.Rcpt(addr); rcptErr != nil {
			return rcptErr
		}
	}

	// Send the email message content
	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write(msg)
	return err
}

// buildMessage constructs a properly formatted email message
// It creates the email headers and body according to RFC 5322 standards
// The message includes MIME headers for HTML content with UTF-8 encoding
func (s *emailService) buildMessage(from, to, subject, body string) []byte {
	// Build email headers following RFC 5322 format
	msg := fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)

	// Add MIME headers for HTML content support
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"

	// Separate headers from body with empty line
	msg += "\r\n"

	// Add the email body content
	msg += body

	return []byte(msg)
}
