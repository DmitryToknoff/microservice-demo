package logger

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Level  string `envconfig:"LEVEL" default:"info"`
	Folder string `envconfig:"FOLDER" default:"./logs"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("LOGGER", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func MustNewConfig() *Config {
	cfg, err := NewConfig()

	if err != nil {
		panic(err)
	}

	return cfg
}
