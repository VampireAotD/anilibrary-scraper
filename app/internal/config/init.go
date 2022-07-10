package config

import (
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

func New(path string) (*Config, error) {
	var cfg Config

	path = filepath.Clean(path)

	err := cleanenv.ReadConfig(path, &cfg)

	if err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
