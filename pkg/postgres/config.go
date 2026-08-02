package postgres

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"HOST" required:"true"`
	Port     string `envconfig:"PORT" default:"5432"`
	User     string `envconfig:"USER" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	DB       string `envconfig:"DB" required:"true"`
	SSLMode  string `envconfig:"SSLMODE" default:"disable"`

	MaxConns          int32         `envconfig:"MAX_CONNS" default:"10"`
	MinConns          int32         `envconfig:"MIN_CONNS" default:"2"`
	MaxConnLifetime   time.Duration `envconfig:"MAX_CONN_LIFETIME" default:"1h"`
	MaxConnIdleTime   time.Duration `envconfig:"MAX_CONN_IDLE_TIME" default:"30m"`
	HealthCheckPeriod time.Duration `envconfig:"HEALTH_CHECK_PERIOD" default:"1m"`
	ConnectTimeout    time.Duration `envconfig:"CONNECT_TIMEOUT" default:"5s"`
	OpTimeout         time.Duration `envconfig:"OP_TIMEOUT" default:"5s"`

	Addr string `ignored:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("POSTGRES", &cfg); err != nil {
		return nil, fmt.Errorf("process postgres config: %w", err)
	}

	cfg.Addr = cfg.DSN()

	return &cfg, nil
}

func NewConfigMust() *Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}

func (c *Config) DSN() string {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, c.Port),
		Path:   c.DB,
	}

	q := dsn.Query()
	q.Set("sslmode", c.SSLMode)
	dsn.RawQuery = q.Encode()

	return dsn.String()
}
