package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

func New() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
