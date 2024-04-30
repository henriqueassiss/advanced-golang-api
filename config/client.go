package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Client struct {
	BaseUrl string `required:"true" split_words:"true"`
}

func NewClient() Client {
	var client Client
	envconfig.MustProcess("CLIENT", &client)

	return client
}
