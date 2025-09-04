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

var (
	client *redis.Client
	once   sync.Once
)

const (
	defaultConnectTimeout = 5 * time.Second
	defaultReadTimeout    = 3 * time.Second
	defaultWriteTimeout   = 3 * time.Second
	defaultPoolSize       = 10
	defaultPoolTimeout    = 30 * time.Second
)

func Init(cfg *config.Config) error {
	var initErr error
	once.Do(func() {
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

		ctx, cancel := context.WithTimeout(context.Background(), defaultConnectTimeout)
		defer cancel()

		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			initErr = fmt.Errorf("failed to connect to Redis: %w", err)
			return
		}

		client = rdb
		// Only log if default logger is available
		if defaultLogger := logger.GetDefault(); defaultLogger != nil {
			defaultLogger.Info("Redis client initialized successfully")
		}
	})

	return initErr
}

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

func Close() error {
	if client == nil {
		return nil
	}
	return client.Close()
}

func Ping(ctx context.Context) (string, error) {
	return GetClient().Ping(ctx).Result()
}

// String operations
func Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return GetClient().Set(ctx, key, value, expiration).Err()
}

func Get(ctx context.Context, key string) (string, error) {
	return GetClient().Get(ctx, key).Result()
}

func GetDel(ctx context.Context, key string) (string, error) {
	return GetClient().GetDel(ctx, key).Result()
}

func Exists(ctx context.Context, keys ...string) (int64, error) {
	return GetClient().Exists(ctx, keys...).Result()
}

func Del(ctx context.Context, keys ...string) (int64, error) {
	return GetClient().Del(ctx, keys...).Result()
}

func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return GetClient().Expire(ctx, key, expiration).Err()
}

func TTL(ctx context.Context, key string) (time.Duration, error) {
	return GetClient().TTL(ctx, key).Result()
}

func Incr(ctx context.Context, key string) (int64, error) {
	return GetClient().Incr(ctx, key).Result()
}

func IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return GetClient().IncrBy(ctx, key, value).Result()
}

func SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	return GetClient().SetNX(ctx, key, value, expiration).Result()
}

func MGet(ctx context.Context, keys ...string) ([]any, error) {
	return GetClient().MGet(ctx, keys...).Result()
}

func MSet(ctx context.Context, values ...any) error {
	return GetClient().MSet(ctx, values...).Err()
}

// List operations
func LPush(ctx context.Context, key string, values ...any) (int64, error) {
	return GetClient().LPush(ctx, key, values...).Result()
}

func RPush(ctx context.Context, key string, values ...any) (int64, error) {
	return GetClient().RPush(ctx, key, values...).Result()
}

func LPop(ctx context.Context, key string) (string, error) {
	return GetClient().LPop(ctx, key).Result()
}

func RPop(ctx context.Context, key string) (string, error) {
	return GetClient().RPop(ctx, key).Result()
}

func LLen(ctx context.Context, key string) (int64, error) {
	return GetClient().LLen(ctx, key).Result()
}

func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return GetClient().LRange(ctx, key, start, stop).Result()
}

func LIndex(ctx context.Context, key string, index int64) (string, error) {
	return GetClient().LIndex(ctx, key, index).Result()
}

func LSet(ctx context.Context, key string, index int64, value any) error {
	return GetClient().LSet(ctx, key, index, value).Err()
}

func LTrim(ctx context.Context, key string, start, stop int64) error {
	return GetClient().LTrim(ctx, key, start, stop).Err()
}

// Hash operations
func HSet(ctx context.Context, key string, values ...any) (int64, error) {
	return GetClient().HSet(ctx, key, values...).Result()
}

func HGet(ctx context.Context, key, field string) (string, error) {
	return GetClient().HGet(ctx, key, field).Result()
}

func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return GetClient().HGetAll(ctx, key).Result()
}

func HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return GetClient().HDel(ctx, key, fields...).Result()
}

func HExists(ctx context.Context, key, field string) (bool, error) {
	return GetClient().HExists(ctx, key, field).Result()
}

func HLen(ctx context.Context, key string) (int64, error) {
	return GetClient().HLen(ctx, key).Result()
}

func HKeys(ctx context.Context, key string) ([]string, error) {
	return GetClient().HKeys(ctx, key).Result()
}

func HVals(ctx context.Context, key string) ([]string, error) {
	return GetClient().HVals(ctx, key).Result()
}

func HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	return GetClient().HIncrBy(ctx, key, field, incr).Result()
}

func HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	return GetClient().HMGet(ctx, key, fields...).Result()
}

// Set operations
func SAdd(ctx context.Context, key string, members ...any) (int64, error) {
	return GetClient().SAdd(ctx, key, members...).Result()
}

func SMembers(ctx context.Context, key string) ([]string, error) {
	return GetClient().SMembers(ctx, key).Result()
}

func SIsMember(ctx context.Context, key string, member any) (bool, error) {
	return GetClient().SIsMember(ctx, key, member).Result()
}

func SCard(ctx context.Context, key string) (int64, error) {
	return GetClient().SCard(ctx, key).Result()
}

func SRem(ctx context.Context, key string, members ...any) (int64, error) {
	return GetClient().SRem(ctx, key, members...).Result()
}

func SPop(ctx context.Context, key string) (string, error) {
	return GetClient().SPop(ctx, key).Result()
}

func SRandMember(ctx context.Context, key string) (string, error) {
	return GetClient().SRandMember(ctx, key).Result()
}

func SDiff(ctx context.Context, keys ...string) ([]string, error) {
	return GetClient().SDiff(ctx, keys...).Result()
}

func SInter(ctx context.Context, keys ...string) ([]string, error) {
	return GetClient().SInter(ctx, keys...).Result()
}

func SUnion(ctx context.Context, keys ...string) ([]string, error) {
	return GetClient().SUnion(ctx, keys...).Result()
}

// Message Queue operations
func Publish(ctx context.Context, channel string, message any) (int64, error) {
	return GetClient().Publish(ctx, channel, message).Result()
}

func Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return GetClient().Subscribe(ctx, channels...)
}

func PSubscribe(ctx context.Context, patterns ...string) *redis.PubSub {
	return GetClient().PSubscribe(ctx, patterns...)
}

func BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return GetClient().BLPop(ctx, timeout, keys...).Result()
}

func BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return GetClient().BRPop(ctx, timeout, keys...).Result()
}

func RPoplPush(ctx context.Context, source, destination string) (string, error) {
	return GetClient().RPopLPush(ctx, source, destination).Result()
}

func BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (string, error) {
	return GetClient().BRPopLPush(ctx, source, destination, timeout).Result()
}
