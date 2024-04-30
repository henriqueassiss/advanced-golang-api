package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Api struct {
	Host              string        `required:"true"`
	Port              string        `required:"true"`
	Secret            string        `required:"true"`
	ReadHeaderTimeout time.Duration `split_words:"true" required:"true"`
	GracefulTimeout   time.Duration `split_words:"true" required:"true"`

	RequestLog bool `split_words:"true" required:"true"`
}

func API() Api {
	var api Api
	envconfig.MustProcess("API", &api)

	return api
}
