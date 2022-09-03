package config

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
	Host     string `env:"REDIS_HOSTNAME" env-default:"0.0.0.0"`
	Password string `env:"REDIS_PASSWORD" env-default:"default"`
	Port     int    `env:"REDIS_PORT" env-default:"6379"`
}
