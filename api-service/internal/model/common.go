package model

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
)

// JSON custom JSON type
type JSON map[string]interface{}

// Scan implements sql.Scanner interface for GORM
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot convert %T to JSON", value)
	}

	if len(bytes) == 0 {
		*j = make(map[string]interface{})
		return nil
	}

	return json.Unmarshal(bytes, j)
}

// Value implements driver.Valuer interface for GORM
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

// GenerateCode generates a random code with prefix and specified length
// The code format is: {prefix}_{random_string}
// The random string consists of lowercase letters and digits
func GenerateCode(prefix string, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, charsetLen)
		result[i] = charset[num.Int64()]
	}

	return fmt.Sprintf("%s_%s", prefix, string(result))
}
