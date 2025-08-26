package security

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken hashes a token for secure storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
