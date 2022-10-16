package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Env string

const (
	Local Env = "local"
)

type Config struct {
	App    App
	HTTP   HTTP
	Redis  Redis
	Jaeger Jaeger
}

func New() (Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

type HTTP struct {
	Addr string `env:"HTTP_ADDR" env-default:"localhost"`
	Port int    `env:"HTTP_PORT" env-required:""`
}

type App struct {
	Timezone string `env:"TIMEZONE" env-default:"Europe/Kiev"`
	Env      Env    `env:"APP_ENV" env-default:"production"`
	Name     string `env:"APP_NAME" env-default:"anilibrary-scraper"`
}

type Redis struct {
	Address     string        `env:"REDIS_ADDRESS" env-required:""`
	Password    string        `env:"REDIS_PASSWORD" env-default:"default"`
	PoolTimeout time.Duration `env:"REDIS_POOL_TIMEOUT" env-default:"5s"`
	PoolSize    int           `env:"REDIS_POOL_SIZE" env-default:"-1"`
	IdleSize    int           `env:"REDIS_IDLE_SIZE" env-default:"-1"`
}

type Jaeger struct {
	TraceEndpoint string `env:"JAEGER_TRACE_ENDPOINT"`
}
