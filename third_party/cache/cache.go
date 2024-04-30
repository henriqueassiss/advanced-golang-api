package cache

import (
	"github.com/henriqueassiss/advanced-golang-api/config"

	"github.com/redis/go-redis/v9"
)

func New(cfg config.Cache) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}
