package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Homyakadze14/RecipeSite/config"
	"github.com/redis/go-redis/v9"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = 5 * time.Second
)

func New(cfg config.Redis) (*redis.Client, error) {
	connAttempts := _defaultConnAttempts
	connTimeout := _defaultConnTimeout

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.ADDRESS,
		Password: cfg.PASSWORD,
		DB:       0,
	})

	var err error

	for connAttempts > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), _defaultConnTimeout)
		defer cancel()

		if err = client.Ping(ctx).Err(); err == nil {
			break
		}

		log.Printf("Redis is trying to connect, attempts left: %d", connAttempts)

		time.Sleep(connTimeout)

		connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("redis - NewRedis - connAttempts == 0: %w", err)
	}

	return client, nil
}
