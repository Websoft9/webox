package utils

import (
	"testing"
)

// TestIsPrivateIP_CorrectTest 正确的测试用例
func TestIsPrivateIP_CorrectTest(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{
			name:     "IPv4 private 10.x.x.x",
			ip:       "10.0.0.1",
			expected: true,
		},
		{
			name:     "IPv4 private 192.168.x.x",
			ip:       "192.168.1.1",
			expected: true,
		},
		{
			name:     "IPv4 public",
			ip:       "8.8.8.8",
			expected: false,
		},
		{
			name:     "IPv4 loopback",
			ip:       "127.0.0.1",
			expected: true,
		},
		{
			name:     "invalid IP",
			ip:       "invalid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPrivateIP(tt.ip)
			if result != tt.expected {
				t.Errorf("IsPrivateIP(%s) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}

// TestIsPrivateIP_FailingTest 故意失败的测试用例
func TestIsPrivateIP_FailingTest(t *testing.T) {
	// 这个测试故意设置错误的期望值，导致测试失败
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{
			name:     "IPv4 private 10.x.x.x - wrong expectation",
			ip:       "10.0.0.1",
			expected: false, // 故意设置错误：私有IP应该返回true，但这里期望false
		},
		{
			name:     "IPv4 public - wrong expectation",
			ip:       "8.8.8.8",
			expected: true, // 故意设置错误：公共IP应该返回false，但这里期望true
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPrivateIP(tt.ip)
			if result != tt.expected {
				t.Errorf("IsPrivateIP(%s) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}
