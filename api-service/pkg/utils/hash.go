package utils

import (
	"crypto/sha256"
	"fmt"
)

// SHA256Hash 使用 SHA256 算法生成哈希值
// 替代不安全的 MD5 算法
func SHA256Hash(text string) string {
	hash := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", hash)
}

// MD5Hash 已废弃：使用 SHA256Hash 替代
// Deprecated: MD5 is cryptographically broken and should not be used for security purposes.
// Use SHA256Hash instead.
func MD5Hash(text string) string {
	// 为了向后兼容，使用 SHA256 替代 MD5
	return SHA256Hash(text)
}
