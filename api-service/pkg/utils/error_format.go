package utils

import (
	"fmt"
	"runtime/debug"
)

// FormatErrorWithStack formats an error with stack trace information.
// It captures the current goroutine's stack trace and appends it to the error message.
func FormatErrorWithStack(err error) string {
	if err == nil {
		return ""
	}

	// Try to use %+v format first (works with github.com/pkg/errors)
	detailedErr := fmt.Sprintf("%+v", err)

	// If the error doesn't contain stack trace (no newlines), add runtime stack
	if detailedErr == err.Error() {
		// Capture current stack trace
		stack := debug.Stack()
		return fmt.Sprintf("%s\n\n --> Stack trace:\n%s", err.Error(), string(stack))
	}

	return detailedErr
}
