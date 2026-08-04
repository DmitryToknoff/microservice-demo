package redis

import (
	"fmt"
	"net"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"HOST" required:"true"`
	Port     string `envconfig:"PORT" default:"6379"`
	Password string `envconfig:"PASSWORD" default:""`
	DB       int    `envconfig:"DB" default:"0"`

	PoolSize     int           `envconfig:"POOL_SIZE" default:"10"`
	MinIdleConns int           `envconfig:"MIN_IDLE_CONNS" default:"5"`
	DialTimeout  time.Duration `envconfig:"DIAL_TIMEOUT" default:"5s"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"3s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"3s"`

	Addr string `ignored:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("REDIS", &cfg); err != nil {
		return nil, fmt.Errorf("process redis config: %w", err)
	}

	cfg.Addr = net.JoinHostPort(cfg.Host, cfg.Port)
	return &cfg, nil
}

func MustNewConfig() *Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
