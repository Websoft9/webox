package redis

import (
	"api-service/internal/config"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       2, // Use DB 2 for testing
		},
	}

	err := Init(cfg)
	require.NoError(t, err, "Failed to initialize Redis for testing")

	// Clean test database before running tests
	err = GetClient().FlushDB(context.Background()).Err()
	require.NoError(t, err, "Failed to flush test database")
}

func TestStringOperations(t *testing.T) {
	setupTestRedis(t)
	ctx := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		key := "test_key"
		value := "test_value"

		err := Set(ctx, key, value, time.Minute*3)
		assert.NoError(t, err)

		result, err := Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("SetNX", func(t *testing.T) {
		key := "test_setnx"
		value := "test_value"

		// First SetNX should succeed
		ok, err := SetNX(ctx, key, value, time.Minute*3)
		assert.NoError(t, err)
		assert.True(t, ok)

		// Second SetNX should fail (key exists)
		ok, err = SetNX(ctx, key, "new_value", time.Minute*3)
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("Incr and IncrBy", func(t *testing.T) {
		key := "test_counter"

		result, err := Incr(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), result)

		result, err = IncrBy(ctx, key, 5)
		assert.NoError(t, err)
		assert.Equal(t, int64(6), result)
	})

	t.Run("MSet and MGet", func(t *testing.T) {
		err := MSet(ctx, "key1", "value1", "key2", "value2")
		assert.NoError(t, err)

		result, err := MGet(ctx, "key1", "key2", "nonexistent")
		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, "value1", result[0])
		assert.Equal(t, "value2", result[1])
		assert.Nil(t, result[2])
	})

	t.Run("Exists and Del", func(t *testing.T) {
		key := "test_exists"
		_ = Set(ctx, key, "value", time.Minute*3)

		count, err := Exists(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		delCount, err := Del(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), delCount)

		count, err = Exists(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestListOperations(t *testing.T) {
	setupTestRedis(t)
	ctx := context.Background()

	t.Run("Push and Pop", func(t *testing.T) {
		key := "test_list"

		// LPush
		count, err := LPush(ctx, key, "item1", "item2")
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// RPush
		count, err = RPush(ctx, key, "item3")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// LPop
		item, err := LPop(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, "item2", item)

		// RPop
		item, err = RPop(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, "item3", item)

		// LLen
		length, err := LLen(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), length)
	})

	t.Run("LRange and LIndex", func(t *testing.T) {
		key := "test_range"
		_, _ = LPush(ctx, key, "c", "b", "a")

		items, err := LRange(ctx, key, 0, -1)
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, items)

		item, err := LIndex(ctx, key, 1)
		assert.NoError(t, err)
		assert.Equal(t, "b", item)
	})
}

func TestHashOperations(t *testing.T) {
	setupTestRedis(t)
	ctx := context.Background()

	t.Run("HSet and HGet", func(t *testing.T) {
		key := "test_hash"
		field := "field1"
		value := "value1"

		count, err := HSet(ctx, key, field, value)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		result, err := HGet(ctx, key, field)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("HGetAll", func(t *testing.T) {
		key := "test_hash_all"
		_, _ = HSet(ctx, key, "field1", "value1", "field2", "value2")

		result, err := HGetAll(ctx, key)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "value1", result["field1"])
		assert.Equal(t, "value2", result["field2"])
	})

	t.Run("HExists and HDel", func(t *testing.T) {
		key := "test_hash_del"
		field := "field1"
		_, _ = HSet(ctx, key, field, "value")

		exists, err := HExists(ctx, key, field)
		assert.NoError(t, err)
		assert.True(t, exists)

		count, err := HDel(ctx, key, field)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		exists, err = HExists(ctx, key, field)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestSetOperations(t *testing.T) {
	setupTestRedis(t)
	ctx := context.Background()

	t.Run("SAdd and SMembers", func(t *testing.T) {
		key := "test_set"

		count, err := SAdd(ctx, key, "member1", "member2", "member3")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		members, err := SMembers(ctx, key)
		assert.NoError(t, err)
		assert.Len(t, members, 3)
		assert.Contains(t, members, "member1")
		assert.Contains(t, members, "member2")
		assert.Contains(t, members, "member3")
	})

	t.Run("SIsMember and SCard", func(t *testing.T) {
		key := "test_set_member"
		_, _ = SAdd(ctx, key, "member1", "member2")

		isMember, err := SIsMember(ctx, key, "member1")
		assert.NoError(t, err)
		assert.True(t, isMember)

		isMember, err = SIsMember(ctx, key, "nonexistent")
		assert.NoError(t, err)
		assert.False(t, isMember)

		count, err := SCard(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})

	t.Run("SRem", func(t *testing.T) {
		key := "test_set_rem"
		_, _ = SAdd(ctx, key, "member1", "member2")

		count, err := SRem(ctx, key, "member1")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		isMember, err := SIsMember(ctx, key, "member1")
		assert.NoError(t, err)
		assert.False(t, isMember)
	})
}

func TestMessageQueueOperations(t *testing.T) {
	setupTestRedis(t)
	ctx := context.Background()

	t.Run("Publish", func(t *testing.T) {
		channel := "test_channel"
		message := "test_message"

		// Since we don't have subscribers, this should return 0
		count, err := Publish(ctx, channel, message)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("BLPop", func(t *testing.T) {
		key := "test_queue"
		_, _ = LPush(ctx, key, "message1", "message2")

		result, err := BLPop(ctx, time.Second, key)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, key, result[0])
		assert.Equal(t, "message2", result[1])
	})
}

func TestExpiration(t *testing.T) {
	setupTestRedis(t)
	ctx := context.Background()

	t.Run("TTL and Expire", func(t *testing.T) {
		key := "test_ttl"
		_ = Set(ctx, key, "value", time.Second*10)

		ttl, err := TTL(ctx, key)
		assert.NoError(t, err)
		assert.True(t, ttl > 0)

		err = Expire(ctx, key, time.Second*20)
		assert.NoError(t, err)

		newTTL, err := TTL(ctx, key)
		assert.NoError(t, err)
		assert.True(t, newTTL > ttl)
	})
}
