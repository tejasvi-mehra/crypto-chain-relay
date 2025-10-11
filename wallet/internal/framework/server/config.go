package server

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type ServiceServerConfig struct {
	Host         string        `env:"HOST" envDefault:"127.0.0.1"`
	Port         int           `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"15s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"15s"`
}

func NewServiceServerConfig() (*ServiceServerConfig, error) {
	c := &ServiceServerConfig{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *ServiceServerConfig) ToServerConfig() *Config {
	return &Config{
		Name:         "service",
		Host:         c.Host,
		Port:         c.Port,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		GracePeriod:  5 * time.Second,
	}
}
