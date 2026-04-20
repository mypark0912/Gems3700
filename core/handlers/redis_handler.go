package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisHandler struct {
	Client *redis.Client
}

func NewRedisHandler(client *redis.Client) *RedisHandler {
	return &RedisHandler{Client: client}
}

// GetBinaryData performs RPOP on the specified queue and returns raw bytes.
func (rh *RedisHandler) GetBinaryData(queueName string) ([]byte, error) {
	ctx := context.Background()
	// Select DB 1
	pipe := rh.Client.Pipeline()
	pipe.Do(ctx, "SELECT", 1)
	cmd := pipe.RPop(ctx, queueName)
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("GetBinaryData error: %w", err)
	}

	result, err := cmd.Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("GetBinaryData RPOP error: %w", err)
	}

	return result, nil
}

// SaveSummary serializes data as JSON and stores it in a Redis hash.
func (rh *RedisHandler) SaveSummary(key, field string, data interface{}) error {
	ctx := context.Background()

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("SaveSummary marshal error: %w", err)
	}

	pipe := rh.Client.Pipeline()
	pipe.Do(ctx, "SELECT", 1)
	pipe.HSet(ctx, key, field, string(jsonBytes))
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("SaveSummary HSET error: %w", err)
	}

	return nil
}

// GetSummary retrieves a JSON-encoded value from a Redis hash.
func (rh *RedisHandler) GetSummary(key, field string) (map[string]interface{}, error) {
	ctx := context.Background()

	pipe := rh.Client.Pipeline()
	pipe.Do(ctx, "SELECT", 1)
	cmd := pipe.HGet(ctx, key, field)
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("GetSummary error: %w", err)
	}

	val, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("GetSummary HGET error: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, fmt.Errorf("GetSummary unmarshal error: %w", err)
	}

	return result, nil
}

// Get retrieves a string value from Redis by key.
func (rh *RedisHandler) Get(ctx context.Context, key string) (string, error) {
	pipe := rh.Client.Pipeline()
	pipe.Do(ctx, "SELECT", 1)
	cmd := pipe.Get(ctx, key)
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return "", fmt.Errorf("Get error: %w", err)
	}
	result, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("Get result error: %w", err)
	}
	return result, nil
}

// Set sets a key-value pair in Redis with an optional expiration.
func (rh *RedisHandler) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	pipe := rh.Client.Pipeline()
	pipe.Do(ctx, "SELECT", 1)
	pipe.Set(ctx, key, value, expiration)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("Set error: %w", err)
	}
	return nil
}

// LPop pops an element from the left of a Redis list.
func (rh *RedisHandler) LPop(ctx context.Context, key string) (string, error) {
	pipe := rh.Client.Pipeline()
	pipe.Do(ctx, "SELECT", 1)
	cmd := pipe.LPop(ctx, key)
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return "", fmt.Errorf("LPop error: %w", err)
	}
	result, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("LPop result error: %w", err)
	}
	return result, nil
}

// HGetAll returns all fields and values of a Redis hash.
func (rh *RedisHandler) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	pipe := rh.Client.Pipeline()
	pipe.Do(ctx, "SELECT", 1)
	cmd := pipe.HGetAll(ctx, key)
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("HGetAll error: %w", err)
	}
	return cmd.Val(), nil
}

// HMSet sets multiple fields in a Redis hash.
func (rh *RedisHandler) HMSet(ctx context.Context, key string, fields map[string]interface{}) error {
	pipe := rh.Client.Pipeline()
	pipe.Do(ctx, "SELECT", 1)
	pipe.HMSet(ctx, key, fields)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("HMSet error: %w", err)
	}
	return nil
}

// HSet sets a single field in a Redis hash.
func (rh *RedisHandler) HSet(ctx context.Context, key, field string, value interface{}) error {
	pipe := rh.Client.Pipeline()
	pipe.Do(ctx, "SELECT", 1)
	pipe.HSet(ctx, key, field, value)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("HSet error: %w", err)
	}
	return nil
}

// GetQueueLength returns the length of a Redis list.
func (rh *RedisHandler) GetQueueLength(queueName string) (int64, error) {
	ctx := context.Background()

	pipe := rh.Client.Pipeline()
	pipe.Do(ctx, "SELECT", 1)
	cmd := pipe.LLen(ctx, queueName)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("GetQueueLength error: %w", err)
	}

	return cmd.Val(), nil
}
