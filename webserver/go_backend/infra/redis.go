package infra

import (
	"context"
	"fmt"
	"log"

	"serverGO/config"

	"github.com/redis/go-redis/v9"
)

type RedisState struct {
	Client0 *redis.Client // DB 0
	Client1 *redis.Client // DB 1
	Client2 *redis.Client // DB 2
}

func NewRedisState(cfg *config.AppConfig) (*RedisState, error) {
	ctx := context.Background()

	client0 := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", cfg.RedisIP),
		DB:   0,
	})
	if err := client0.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis client ping: %w", err)
	}

	client1 := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", cfg.RedisIP),
		DB:   1,
	})

	client2 := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", cfg.RedisIP),
		DB:   2,
	})

	log.Println("Redis connected:", cfg.RedisIP)
	return &RedisState{
		Client0: client0,
		Client1: client1,
		Client2: client2,
	}, nil
}

func (r *RedisState) Close() {
	if r.Client0 != nil {
		r.Client0.Close()
	}
	if r.Client1 != nil {
		r.Client1.Close()
	}
	if r.Client2 != nil {
		r.Client2.Close()
	}
}
