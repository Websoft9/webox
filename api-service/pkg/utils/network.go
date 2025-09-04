package utils

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	unknownValue = "unknown"

	// IPv4 private ranges
	ipv4Class10  = 10  // 10.0.0.0/8
	ipv4Class172 = 172 // 172.16.0.0/12
	ipv4Class192 = 192 // 192.168.0.0/16
	ipv4Loopback = 127 // 127.0.0.0/8

	// IPv6 private ranges
	ipv6LinkLocalMask = 0xc0 // fe80::/10
	ipv6LinkLocalByte = 0x80
)

// GetRealIP extracts the real client IP from the Gin context
// It checks various headers in order of priority to find the real client IP
func GetRealIP(c *gin.Context) string {
	// Define priority order of headers to check
	headerPriority := []string{
		"X-Real-IP",
		"X-Forwarded-For",
		"CF-Connecting-IP",
		"True-Client-IP",
		"X-Forwarded",
	}

	// Check standard headers first
	for _, header := range headerPriority {
		if ip := extractIPFromHeader(c, header); ip != "" {
			return ip
		}
	}

	// Check Forwarded header (RFC 7239) separately due to different format
	if ip := extractIPFromForwardedHeader(c); ip != "" {
		return ip
	}

	// Fall back to Gin's ClientIP() which uses RemoteAddr
	return c.ClientIP()
}

// extractIPFromHeader extracts IP from standard headers
func extractIPFromHeader(c *gin.Context, headerName string) string {
	headerValue := c.GetHeader(headerName)
	if headerValue == "" || headerValue == unknownValue {
		return ""
	}

	// Handle X-Forwarded-For which can contain multiple IPs
	if headerName == "X-Forwarded-For" {
		return getFirstValidIP(headerValue)
	}

	// For other headers, validate and return single IP
	if parsedIP := net.ParseIP(headerValue); parsedIP != nil {
		return headerValue
	}

	return ""
}

// getFirstValidIP extracts the first valid IP from comma-separated list
func getFirstValidIP(ips string) string {
	ipList := strings.Split(ips, ",")
	for _, ip := range ipList {
		cleanIP := strings.TrimSpace(ip)
		if cleanIP != "" && cleanIP != unknownValue {
			if parsedIP := net.ParseIP(cleanIP); parsedIP != nil {
				return cleanIP
			}
		}
	}
	return ""
}

// extractIPFromForwardedHeader extracts IP from Forwarded header (RFC 7239)
func extractIPFromForwardedHeader(c *gin.Context) string {
	forwarded := c.GetHeader("Forwarded")
	if forwarded == "" || forwarded == unknownValue {
		return ""
	}

	// Parse the Forwarded header format: for=192.0.2.43;proto=http;by=203.0.113.43
	parts := strings.Split(forwarded, ";")
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if strings.HasPrefix(strings.ToLower(trimmedPart), "for=") {
			ip := strings.TrimSpace(strings.TrimPrefix(trimmedPart, "for="))
			ip = strings.Trim(ip, `"`) // Remove quotes if present
			if parsedIP := net.ParseIP(ip); parsedIP != nil {
				return ip
			}
		}
	}
	return ""
}

// IsPrivateIP checks if the given IP address is a private IP
func IsPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	return isPrivateIPv4(parsedIP) || isPrivateIPv6(parsedIP)
}

// isPrivateIPv4 checks IPv4 private ranges
func isPrivateIPv4(parsedIP net.IP) bool {
	// Check for IPv4 private ranges
	if ipv4 := parsedIP.To4(); ipv4 != nil {
		return isIPv4InPrivateRange(ipv4)
	}
	return false
}

// isIPv4InPrivateRange checks if IPv4 is in private ranges
func isIPv4InPrivateRange(ipv4 net.IP) bool {
	// 10.0.0.0/8
	if ipv4[0] == ipv4Class10 {
		return true
	}
	// 172.16.0.0/12
	if ipv4[0] == ipv4Class172 && ipv4[1] >= 16 && ipv4[1] <= 31 {
		return true
	}
	// 192.168.0.0/16
	if ipv4[0] == ipv4Class192 && ipv4[1] == 168 {
		return true
	}
	// 127.0.0.0/8 (loopback)
	if ipv4[0] == ipv4Loopback {
		return true
	}
	return false
}

// isPrivateIPv6 checks IPv6 private ranges
func isPrivateIPv6(parsedIP net.IP) bool {
	if ipv6 := parsedIP.To16(); ipv6 != nil {
		// ::1 (loopback)
		if parsedIP.Equal(net.IPv6loopback) {
			return true
		}
		// fc00::/7 (unique local)
		if ipv6[0] == 0xfc || ipv6[0] == 0xfd {
			return true
		}
		// fe80::/10 (link-local)
		if ipv6[0] == 0xfe && (ipv6[1]&ipv6LinkLocalMask) == ipv6LinkLocalByte {
			return true
		}
	}
	return false
}
