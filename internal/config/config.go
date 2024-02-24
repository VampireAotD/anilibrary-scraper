package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/fx"
)

type Config struct {
	fx.Out

	App   App
	JWT   JWT
	Kafka Kafka
	Redis Redis
	HTTP  HTTP
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
	Addr           string `env:"HTTP_ADDR" env-default:"localhost"`
	MainPort       int    `env:"HTTP_MAIN_PORT" env-required:""`
	MonitoringPort int    `env:"HTTP_MONITORING_PORT" env-required:""`
}

type App struct {
	Env  string `env:"APP_ENV" env-default:"local"`
	Name string `env:"APP_NAME" env-default:"anilibrary-scraper"`
}

type Redis struct {
	Address     string        `env:"REDIS_ADDRESS" env-required:""`
	Password    string        `env:"REDIS_PASSWORD" env-required:""`
	PoolTimeout time.Duration `env:"REDIS_POOL_TIMEOUT" env-default:"5s"`
}

type Kafka struct {
	Address   string `env:"KAFKA_ADDRESS" env-required:""`
	Username  string `env:"KAFKA_CLIENT_USERS" env-required:""`
	Password  string `env:"KAFKA_CLIENT_PASSWORDS" env-required:""`
	Topic     string `env:"KAFKA_TOPIC" env-default:"scraper_topic"`
	Partition int    `env:"KAFKA_PARTITION" env-default:"0"`
}

type JWT struct {
	Secret []byte `env:"JWT_SECRET" env-required:""`
}
