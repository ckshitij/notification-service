package redis

import (
	"context"
	"time"

	"github.com/ckshitij/notify-srv/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewClient(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return rdb, nil
}
