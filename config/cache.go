package config

import "github.com/kelseyhightower/envconfig"

type Cache struct {
	Address  string `required:"true"`
	Password string `required:"true"`
	DB       int    `required:"true"`
}

func NewCache() Cache {
	var cache Cache
	envconfig.MustProcess("CACHE", &cache)

	return cache
}
