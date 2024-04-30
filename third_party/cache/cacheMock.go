package cache

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func NewMock(t *testing.T) *redis.Client {
	cacheMock := miniredis.RunT(t)

	return redis.NewClient(&redis.Options{
		Addr: cacheMock.Addr(),
	})
}
