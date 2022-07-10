package config

type Config struct {
	Http Http
}

type Http struct {
	Addr string `env:"HTTP_ADDR" env-default:"localhost"`
	Port int    `env:"HTTP_PORT" env-required:""`
}
