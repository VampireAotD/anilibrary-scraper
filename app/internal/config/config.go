package config

type Config struct {
	App  App
	Http Http
}

type Http struct {
	Addr string `env:"HTTP_ADDR" env-default:"localhost"`
	Port int    `env:"HTTP_PORT" env-required:""`
}

type App struct {
	Timezone string `env:"TIMEZONE" env-default:"Europe/Kiev"`
}
