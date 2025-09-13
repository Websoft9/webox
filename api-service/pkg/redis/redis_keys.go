package redis

import (
	"fmt"
)

// Redis token storage constants
const (
	RedisNilError = "redis: nil"
)

// Redis keys constants
const (
	REDIS_PREFIX_AUTH   = "AUTH"
	REDIS_PREFIX_TOKEN  = "TOKEN"
	REDIS_PREFIX_LOCKER = "LOCKER"
	REDIS_KEYS_JOINER   = ":"

	// #nosec G101 -- This is not a credential, just a Redis key prefix
	RK_AUTH_TOKEN     = "AUTH:TOKEN:%s"
	RK_EMAIL_VERILOCK = "LOCKER:REGISTER:%s" // Prefix for email verification lock in Redis

	// Redis key constants for login security
	RK_LOGIN_ATTEMPTS = "AUTH:LOGIN:ATTEMPTS:%s" // Login attempts counter key
	RK_LOGIN_LOCKER   = "LOCKER:LOGIN:%s"        // Login lockout key
)

func FormatRedisKey(redisPrefix, key string) string {
	return fmt.Sprintf(redisPrefix, key)
}
