package config

import (
	"github.com/kelseyhightower/envconfig"
)

type App struct {
	Environment string `required:"true"`
	CertDir     string `split_words:"true"`
	KeyDir      string `split_words:"true"`
}

func APP() App {
	var app App
	envconfig.MustProcess("APP", &app)

	return app
}
