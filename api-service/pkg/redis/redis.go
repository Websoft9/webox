// Package redis provides a Redis client wrapper with common operations
// It implements singleton pattern for Redis client management
package redis

import (
	"api-service/internal/config"
	"api-service/pkg/logger"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Global variables for Redis client management
var (
	client *redis.Client // Redis client instance (singleton)
	once   sync.Once     // Ensures single initialization
)

// Default timeout and connection pool configuration
const (
	defaultConnectTimeout = 5 * time.Second  // Connection establishment timeout
	defaultReadTimeout    = 3 * time.Second  // Read operation timeout
	defaultWriteTimeout   = 3 * time.Second  // Write operation timeout
	defaultPoolSize       = 10               // Maximum number of connections in pool
	defaultPoolTimeout    = 30 * time.Second // Pool operation timeout
)

// Init initializes the Redis client with the given configuration
// Uses singleton pattern to ensure only one client instance exists
// Returns error if connection fails
func Init(cfg *config.Config) error {
	var initErr error
	once.Do(func() {
		// Create Redis client with configuration
		rdb := redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			Password:     cfg.Redis.Password,
			DB:           cfg.Redis.DB,
			DialTimeout:  defaultConnectTimeout,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			PoolSize:     defaultPoolSize,
			PoolTimeout:  defaultPoolTimeout,
		})

		// Test connection with ping
		ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
		defer cancel()

		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			initErr = fmt.Errorf("failed to connect to Redis: %w", err)
			return
		}

		// Connection successful, store client instance
		client = rdb
		// Only log if default logger is available
		if defaultLogger := logger.GetDefault(); defaultLogger != nil {
			defaultLogger.Info("Redis client initialized successfully")
		}
	})

	return initErr
}

// GetClient returns the Redis client instance
// Panics if client is not initialized - call Init() first
func GetClient() *redis.Client {
	if client == nil {
		// Only log if default logger is available
		if defaultLogger := logger.GetDefault(); defaultLogger != nil {
			defaultLogger.Error("Redis client not initialized. Call Init() first")
		}
		panic("Redis client not initialized")
	}
	return client
}

// Close closes the Redis client connection
// Returns nil if client is not initialized
func Close() error {
	if client == nil {
		return nil
	}
	return client.Close()
}

// Ping tests the connection to Redis server
// Returns "PONG" if connection is healthy
func Ping(ctx context.Context) (string, error) {
	return GetClient().Ping(ctx).Result()
}

// String operations - Redis string data type operations

// Set stores a key-value pair with optional expiration
// If expiration is 0, key will not expire
func Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return GetClient().Set(ctx, key, value, expiration).Err()
}

// Get retrieves the value of a key
// Returns redis.Nil error if key doesn't exist
func Get(ctx context.Context, key string) (string, error) {
	return GetClient().Get(ctx, key).Result()
}

// GetDel atomically gets and deletes a key
// Returns the value before deletion
func GetDel(ctx context.Context, key string) (string, error) {
	return GetClient().GetDel(ctx, key).Result()
}

// Exists checks if keys exist
// Returns the number of existing keys
func Exists(ctx context.Context, keys ...string) (int64, error) {
	return GetClient().Exists(ctx, keys...).Result()
}

// Del deletes one or more keys
// Returns the number of deleted keys
func Del(ctx context.Context, keys ...string) (int64, error) {
	return GetClient().Del(ctx, keys...).Result()
}

// Expire sets a key's time to live
// Key will be deleted after expiration time
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return GetClient().Expire(ctx, key, expiration).Err()
}

// TTL returns the remaining time to live of a key
// Returns -1 if key has no expiration, -2 if key doesn't exist
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return GetClient().TTL(ctx, key).Result()
}

// Incr increments the number stored at key by 1
// Returns the value after increment
func Incr(ctx context.Context, key string) (int64, error) {
	return GetClient().Incr(ctx, key).Result()
}

// IncrBy increments the number stored at key by increment
// Returns the value after increment
func IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return GetClient().IncrBy(ctx, key, value).Result()
}

// SetNX sets key to value only if key doesn't exist (SET if Not eXists)
// Returns true if key was set, false if key already exists
func SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	return GetClient().SetNX(ctx, key, value, expiration).Result()
}

// MGet returns the values of all specified keys
// Returns slice with nil values for non-existing keys
func MGet(ctx context.Context, keys ...string) ([]any, error) {
	return GetClient().MGet(ctx, keys...).Result()
}

// MSet sets multiple key-value pairs atomically
// Values should be provided as alternating key-value pairs
func MSet(ctx context.Context, values ...any) error {
	return GetClient().MSet(ctx, values...).Err()
}

// List operations - Redis list data type operations

// LPush inserts values at the head of the list
// Returns the length of the list after operation
func LPush(ctx context.Context, key string, values ...any) (int64, error) {
	return GetClient().LPush(ctx, key, values...).Result()
}

// RPush inserts values at the tail of the list
// Returns the length of the list after operation
func RPush(ctx context.Context, key string, values ...any) (int64, error) {
	return GetClient().RPush(ctx, key, values...).Result()
}

// LPop removes and returns the first element from the list
// Returns redis.Nil error if list is empty
func LPop(ctx context.Context, key string) (string, error) {
	return GetClient().LPop(ctx, key).Result()
}

// RPop removes and returns the last element from the list
// Returns redis.Nil error if list is empty
func RPop(ctx context.Context, key string) (string, error) {
	return GetClient().RPop(ctx, key).Result()
}

// LLen returns the length of the list
// Returns 0 if key doesn't exist
func LLen(ctx context.Context, key string) (int64, error) {
	return GetClient().LLen(ctx, key).Result()
}

// LRange returns a range of elements from the list
// start and stop are zero-based indexes, negative values count from the end
func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return GetClient().LRange(ctx, key, start, stop).Result()
}

// LIndex returns the element at the specified index
// Negative indexes count from the end (-1 is the last element)
func LIndex(ctx context.Context, key string, index int64) (string, error) {
	return GetClient().LIndex(ctx, key, index).Result()
}

// LSet sets the value of an element at the specified index
// Returns error if index is out of range
func LSet(ctx context.Context, key string, index int64, value any) error {
	return GetClient().LSet(ctx, key, index, value).Err()
}

// LTrim trims the list to contain only elements within the specified range
// Elements outside the range are removed
func LTrim(ctx context.Context, key string, start, stop int64) error {
	return GetClient().LTrim(ctx, key, start, stop).Err()
}

// Hash operations - Redis hash data type operations

// HSet sets field-value pairs in the hash
// Returns the number of fields that were added
func HSet(ctx context.Context, key string, values ...any) (int64, error) {
	return GetClient().HSet(ctx, key, values...).Result()
}

// HGet returns the value associated with field in the hash
// Returns redis.Nil error if field doesn't exist
func HGet(ctx context.Context, key, field string) (string, error) {
	return GetClient().HGet(ctx, key, field).Result()
}

// HGetAll returns all fields and values in the hash
// Returns empty map if key doesn't exist
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return GetClient().HGetAll(ctx, key).Result()
}

// HDel deletes one or more hash fields
// Returns the number of fields that were removed
func HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return GetClient().HDel(ctx, key, fields...).Result()
}

// HExists checks if a field exists in the hash
// Returns true if field exists, false otherwise
func HExists(ctx context.Context, key, field string) (bool, error) {
	return GetClient().HExists(ctx, key, field).Result()
}

// HLen returns the number of fields in the hash
// Returns 0 if key doesn't exist
func HLen(ctx context.Context, key string) (int64, error) {
	return GetClient().HLen(ctx, key).Result()
}

// HKeys returns all field names in the hash
// Returns empty slice if key doesn't exist
func HKeys(ctx context.Context, key string) ([]string, error) {
	return GetClient().HKeys(ctx, key).Result()
}

// HVals returns all values in the hash
// Returns empty slice if key doesn't exist
func HVals(ctx context.Context, key string) ([]string, error) {
	return GetClient().HVals(ctx, key).Result()
}

// HIncrBy increments the number stored at field in the hash by increment
// Returns the value after increment
func HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	return GetClient().HIncrBy(ctx, key, field, incr).Result()
}

// HMGet returns the values associated with the specified fields
// Returns slice with nil values for non-existing fields
func HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	return GetClient().HMGet(ctx, key, fields...).Result()
}

// Set operations - Redis set data type operations

// SAdd adds one or more members to the set
// Returns the number of members that were added (excluding existing members)
func SAdd(ctx context.Context, key string, members ...any) (int64, error) {
	return GetClient().SAdd(ctx, key, members...).Result()
}

// SMembers returns all members of the set
// Returns empty slice if key doesn't exist
func SMembers(ctx context.Context, key string) ([]string, error) {
	return GetClient().SMembers(ctx, key).Result()
}

// SIsMember checks if member is in the set
// Returns true if member exists, false otherwise
func SIsMember(ctx context.Context, key string, member any) (bool, error) {
	return GetClient().SIsMember(ctx, key, member).Result()
}

// SCard returns the number of members in the set
// Returns 0 if key doesn't exist
func SCard(ctx context.Context, key string) (int64, error) {
	return GetClient().SCard(ctx, key).Result()
}

// SRem removes one or more members from the set
// Returns the number of members that were removed
func SRem(ctx context.Context, key string, members ...any) (int64, error) {
	return GetClient().SRem(ctx, key, members...).Result()
}

// SPop removes and returns a random member from the set
// Returns redis.Nil error if set is empty
func SPop(ctx context.Context, key string) (string, error) {
	return GetClient().SPop(ctx, key).Result()
}

// SRandMember returns a random member from the set without removing it
// Returns redis.Nil error if set is empty
func SRandMember(ctx context.Context, key string) (string, error) {
	return GetClient().SRandMember(ctx, key).Result()
}

// SDiff returns the members of the set resulting from the difference between the first set and all successive sets
// Returns empty slice if no difference exists
func SDiff(ctx context.Context, keys ...string) ([]string, error) {
	return GetClient().SDiff(ctx, keys...).Result()
}

// SInter returns the members of the set resulting from the intersection of all given sets
// Returns empty slice if no intersection exists
func SInter(ctx context.Context, keys ...string) ([]string, error) {
	return GetClient().SInter(ctx, keys...).Result()
}

// SUnion returns the members of the set resulting from the union of all given sets
// Returns empty slice if all sets are empty
func SUnion(ctx context.Context, keys ...string) ([]string, error) {
	return GetClient().SUnion(ctx, keys...).Result()
}

// Message Queue operations - Redis pub/sub and blocking list operations

// Publish posts a message to the given channel
// Returns the number of clients that received the message
func Publish(ctx context.Context, channel string, message any) (int64, error) {
	return GetClient().Publish(ctx, channel, message).Result()
}

// Subscribe subscribes to the given channels
// Returns a PubSub instance for receiving messages
func Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return GetClient().Subscribe(ctx, channels...)
}

// PSubscribe subscribes to channels matching the given patterns
// Returns a PubSub instance for receiving messages
func PSubscribe(ctx context.Context, patterns ...string) *redis.PubSub {
	return GetClient().PSubscribe(ctx, patterns...)
}

// BLPop blocks until an element is available to pop from the head of any of the given lists
// Returns the key and element, or timeout if no element is available within timeout
func BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return GetClient().BLPop(ctx, timeout, keys...).Result()
}

// BRPop blocks until an element is available to pop from the tail of any of the given lists
// Returns the key and element, or timeout if no element is available within timeout
func BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return GetClient().BRPop(ctx, timeout, keys...).Result()
}

// RPoplPush atomically removes the last element from source list and pushes it to destination list
// Returns the element that was moved
func RPoplPush(ctx context.Context, source, destination string) (string, error) {
	return GetClient().RPopLPush(ctx, source, destination).Result()
}

// BRPopLPush blocks until an element is available to pop from the tail of source list and push to destination list
// Returns the element that was moved, or timeout if no element is available within timeout
func BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (string, error) {
	return GetClient().BRPopLPush(ctx, source, destination, timeout).Result()
}
