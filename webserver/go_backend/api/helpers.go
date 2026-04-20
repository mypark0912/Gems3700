package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/redis/go-redis/v9"
	"serverGO/infra"
)

// getRedisKey returns the Redis hash key for a channel
// e.g., channel="Main" → "Main_meter", "Main_pq", etc.
func getRedisKey(channel, item string) string {
	return fmt.Sprintf("%s_%s", channel, item)
}

// getAllRedisKeys returns all Redis key groups for a channel
func getAllRedisKeys(channel string) map[string]string {
	return map[string]string{
		"meter": getRedisKey(channel, "meter"),
		"pq":    getRedisKey(channel, "pq"),
		"alarm": getRedisKey(channel, "alarm"),
		"thd":   getRedisKey(channel, "thd"),
	}
}

// loadMeterFromRedis loads all fields from a Redis hash as map[string]float64
func loadMeterFromRedis(ctx context.Context, client *redis.Client, key string, fields []string) map[string]float64 {
	result := make(map[string]float64, len(fields))
	if len(fields) == 0 {
		// Load all
		all, err := client.HGetAll(ctx, key).Result()
		if err != nil {
			return result
		}
		for k, v := range all {
			var f float64
			fmt.Sscanf(v, "%f", &f)
			result[k] = f
		}
		return result
	}

	for _, field := range fields {
		val, err := client.HGet(ctx, key, field).Float64()
		if err == nil {
			result[field] = val
		}
	}
	return result
}

// loadJSONFromRedis reads a JSON string field from Redis hash
func loadJSONFromRedis(ctx context.Context, client *redis.Client, key, field string) (map[string]interface{}, error) {
	val, err := client.HGet(ctx, key, field).Result()
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// getChannelSetting loads channel config from Redis setup
func getChannelSetting(ctx context.Context, deps *infra.Dependencies, channel string) map[string]interface{} {
	setupJSON, err := deps.Redis.Client0.HGet(ctx, "System", "setup").Result()
	if err != nil {
		return nil
	}

	var setup map[string]interface{}
	json.Unmarshal([]byte(setupJSON), &setup)

	channels, _ := setup["channel"].([]interface{})
	for _, ch := range channels {
		chMap, ok := ch.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := chMap["channel"].(string)
		if strings.EqualFold(name, channel) {
			return chMap
		}
	}
	return nil
}

// isServiceActive checks if a systemd service is active
func isServiceActive(name string) bool {
	cmd := exec.Command("systemctl", "is-active", name)
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "active"
}
