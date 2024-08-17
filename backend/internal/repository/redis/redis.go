package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Homyakadze14/RecipeSite/internal/common"
	"github.com/redis/go-redis/v9"
)

const (
	_defaultExperationTime = 5 * time.Minute
)

type RedisRepo struct {
	redis *redis.Client
}

func NewRedisRepository(redis *redis.Client) *RedisRepo {
	return &RedisRepo{redis}
}

func (r *RedisRepo) Set(ctx context.Context, key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("RedisRepo - Set - json.Marshal: %w", err)
	}

	err = r.redis.Set(ctx, key, p, _defaultExperationTime).Err()
	if err != nil {
		return fmt.Errorf("RedisRepo - Set - r.redis.Set: %w", err)
	}

	return nil
}

func (r *RedisRepo) Get(ctx context.Context, key string, dest interface{}) error {
	var value []byte
	err := r.redis.Get(ctx, key).Scan(&value)
	if err != nil {
		if err == redis.Nil {
			return common.ErrCacheKeyNotFound
		}
		return fmt.Errorf("RedisRepo - Get - r.redis.Get: %w", err)
	}
	return json.Unmarshal(value, dest)
}

func (r *RedisRepo) Del(ctx context.Context, key string) (res int64, err error) {
	res, err = r.redis.Del(ctx, key).Result()
	if err != nil {
		return res, fmt.Errorf("RedisRepo - Del - r.redis.Del: %w", err)
	}
	return res, nil
}
