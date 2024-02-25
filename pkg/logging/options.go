package logging

import (
	"io"
)

type config struct {
	output        io.Writer
	ecsCompatible bool
	jsonEncoder   bool
}

type Option func(cfg *config)

func WithOutput(output io.Writer) Option {
	return func(cfg *config) {
		cfg.output = output
	}
}

func ECSCompatible() Option {
	return func(cfg *config) {
		cfg.ecsCompatible = true
	}
}

func ConvertToJSON() Option {
	return func(cfg *config) {
		cfg.jsonEncoder = true
	}
}
