package handlers

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisPoolMu sync.Mutex
	redisPool   = make(map[string]*redis.Client)
)

func redisPoolKey(ip string, db int) string {
	return fmt.Sprintf("%s:%d", ip, db)
}

func GetRedisClient(ip string, db int) *redis.Client {
	key := redisPoolKey(ip, db)

	redisPoolMu.Lock()
	defer redisPoolMu.Unlock()

	if client, ok := redisPool[key]; ok {
		return client
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", ip),
		DB:       db,
		PoolSize: 7,
	})

	redisPool[key] = client
	return client
}

func GetRedisBinaryClient(ip string, db int) *redis.Client {
	return GetRedisClient(ip, db)
}

func CleanupRedisPools() {
	redisPoolMu.Lock()
	defer redisPoolMu.Unlock()

	for key, client := range redisPool {
		_ = client.Close()
		delete(redisPool, key)
	}
}

func PingRedis(client *redis.Client) error {
	return client.Ping(context.Background()).Err()
}
