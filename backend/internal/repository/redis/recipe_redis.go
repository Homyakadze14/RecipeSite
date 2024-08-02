package redisrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/redis/go-redis/v9"
)

const (
	_defaultExperationTime = 5 * time.Minute
)

var (
	ErrRedisKeyNotFound = errors.New("key not found")
)

type RecipeRedisRepo struct {
	redis *redis.Client
}

func NewRecipeRedisRepository(redis *redis.Client) *RecipeRedisRepo {
	return &RecipeRedisRepo{redis}
}

func (r *RecipeRedisRepo) Set(ctx context.Context, key string, recipe entities.FullRecipe) error {
	err := r.redis.Set(ctx, key, recipe, _defaultExperationTime).Err()

	if err != nil {
		return fmt.Errorf("RecipeRedisRepo - Set - r.redis.Set: %w", err)
	}

	return nil
}

func (r *RecipeRedisRepo) Get(ctx context.Context, key string) (recipe entities.FullRecipe, err error) {
	err = r.redis.Get(ctx, key).Scan(&recipe)

	if err != nil {
		if err == redis.Nil {
			return recipe, ErrRedisKeyNotFound
		}
		return recipe, fmt.Errorf("RecipeRedisRepo - Get - r.redis.Get: %w", err)
	}

	return recipe, nil
}

func (r *RecipeRedisRepo) Del(ctx context.Context, key string) (res int64, err error) {
	res, err = r.redis.Del(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return res, ErrRedisKeyNotFound
		}
		return res, fmt.Errorf("RecipeRedisRepo - Del - r.redis.Del: %w", err)
	}

	return res, nil
}
