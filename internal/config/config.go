package config

import "time"

type Config struct {
	App   App
	HTTP  HTTP
	Redis Redis
}

type HTTP struct {
	Addr string `env:"HTTP_ADDR" env-default:"localhost"`
	Port int    `env:"HTTP_PORT" env-required:""`
}

type App struct {
	Timezone string `env:"TIMEZONE" env-default:"Europe/Kiev"`
}

type Redis struct {
	Address     string        `env:"REDIS_ADDRESS" env-required:""`
	Password    string        `env:"REDIS_PASSWORD" env-default:"default"`
	PoolTimeout time.Duration `env:"REDIS_POOL_TIMEOUT" env-default:"5s"`
	PoolSize    int           `env:"REDIS_POOL_SIZE" env-default:"-1"`
	IdleSize    int           `env:"REDIS_IDLE_SIZE" env-default:"-1"`
}
