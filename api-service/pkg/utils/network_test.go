package utils

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRealIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		headers  map[string]string
		expected string
	}{
		{
			name: "X-Real-IP header",
			headers: map[string]string{
				"X-Real-IP": "192.168.1.100",
			},
			expected: "192.168.1.100",
		},
		{
			name: "X-Forwarded-For single IP",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.195",
			},
			expected: "203.0.113.195",
		},
		{
			name: "X-Forwarded-For multiple IPs",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.195, 70.41.3.18, 150.172.238.178",
			},
			expected: "203.0.113.195",
		},
		{
			name: "X-Forwarded-For with spaces",
			headers: map[string]string{
				"X-Forwarded-For": "  203.0.113.195  ,  70.41.3.18  ",
			},
			expected: "203.0.113.195",
		},
		{
			name: "CF-Connecting-IP header",
			headers: map[string]string{
				"CF-Connecting-IP": "198.51.100.101",
			},
			expected: "198.51.100.101",
		},
		{
			name: "True-Client-IP header",
			headers: map[string]string{
				"True-Client-IP": "198.51.100.102",
			},
			expected: "198.51.100.102",
		},
		{
			name: "X-Forwarded header",
			headers: map[string]string{
				"X-Forwarded": "198.51.100.103",
			},
			expected: "198.51.100.103",
		},
		{
			name: "Forwarded header RFC 7239",
			headers: map[string]string{
				"Forwarded": "for=192.0.2.43;proto=http;by=203.0.113.43",
			},
			expected: "192.0.2.43",
		},
		{
			name: "Forwarded header with quotes",
			headers: map[string]string{
				"Forwarded": `for="192.0.2.44";proto=https`,
			},
			expected: "192.0.2.44",
		},
		{
			name: "Priority test - X-Real-IP wins",
			headers: map[string]string{
				"X-Real-IP":        "192.168.1.100",
				"X-Forwarded-For":  "203.0.113.195",
				"CF-Connecting-IP": "198.51.100.101",
			},
			expected: "192.168.1.100",
		},
		{
			name: "Invalid IP ignored",
			headers: map[string]string{
				"X-Real-IP":       "invalid-ip",
				"X-Forwarded-For": "203.0.113.195",
			},
			expected: "203.0.113.195",
		},
		{
			name: "Unknown value ignored",
			headers: map[string]string{
				"X-Real-IP":       "unknown",
				"X-Forwarded-For": "203.0.113.195",
			},
			expected: "203.0.113.195",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin context with request
			c, _ := gin.CreateTestContext(nil)
			req, _ := http.NewRequest("GET", "/", nil)

			// Set headers
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			c.Request = req

			// Test GetRealIP
			result := GetRealIP(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		// IPv4 private ranges
		{"IPv4 10.x.x.x", "10.0.0.1", true},
		{"IPv4 172.16.x.x", "172.16.0.1", true},
		{"IPv4 172.31.x.x", "172.31.255.255", true},
		{"IPv4 192.168.x.x", "192.168.1.1", true},
		{"IPv4 127.x.x.x", "127.0.0.1", true},

		// IPv4 public IPs
		{"IPv4 public 8.8.8.8", "8.8.8.8", false},
		{"IPv4 public 203.0.113.1", "203.0.113.1", false},
		{"IPv4 172.15.x.x (not private)", "172.15.0.1", false},
		{"IPv4 172.32.x.x (not private)", "172.32.0.1", false},

		// IPv6 addresses
		{"IPv6 loopback", "::1", true},
		{"IPv6 link-local", "fe80::1", true},
		{"IPv6 unique local fc00", "fc00::1", true},
		{"IPv6 unique local fd00", "fd00::1", true},
		{"IPv6 public", "2001:db8::1", false},

		// Invalid IPs
		{"Invalid IP", "invalid-ip", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPrivateIP(tt.ip)
			assert.Equal(t, tt.expected, result, "IP: %s", tt.ip)
		})
	}
}
