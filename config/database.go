package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Driver                 string        `required:"true"`
	Host                   string        `required:"true"`
	Port                   uint16        `required:"true"`
	Name                   string        `required:"true"`
	User                   string        `required:"true"`
	Password               string        `required:"true"`
	SslMode                string        `split_words:"true" required:"true"`
	MaxConnectionPool      int           `split_words:"true" required:"true"`
	MaxIdleConnections     int           `split_words:"true" required:"true"`
	ConnectionsMaxLifeTime time.Duration `split_words:"true" required:"true"`
}

func DataStore() Database {
	var db Database
	envconfig.MustProcess("DB", &db)

	return db
}
