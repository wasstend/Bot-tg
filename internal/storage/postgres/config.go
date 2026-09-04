package postgres

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Username string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" default:"5432"`
	Database string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" default:"30s"`
}

func NewMust() Config {
	var config Config

	if err := envconfig.Process("POSTGRES", &config); err != nil {
		panic("failed to parse postgres config: " + err.Error())
	}

	return config
}
