package logging

import (
	"io"
)

type config struct {
	output        io.Writer
	level         Level
	asDefault     bool
	jsonEncoder   bool
	ecsCompatible bool
}

type Option func(cfg *config)

// WithOutput will set the output of the logger.
func WithOutput(output io.Writer) Option {
	return func(cfg *config) {
		cfg.output = output
	}
}

// WithLevel will set the level of the logger based on the string value.
// If the value is not valid or level is not supported - DebugLevel will be set.
func WithLevel(level string) Option {
	return func(cfg *config) {
		var lvl Level

		if err := lvl.UnmarshalText([]byte(level)); err != nil {
			lvl = DebugLevel
		}

		cfg.level = lvl
	}
}

// AsDefault will set configured logger as default which will be be used by Get.
func AsDefault() Option {
	return func(cfg *config) {
		cfg.asDefault = true
	}
}

// ConvertToJSON will enable json encoder for logging.
func ConvertToJSON() Option {
	return func(cfg *config) {
		cfg.jsonEncoder = true
	}
}

// ECSCompatible will enable ecs compatible scheme for logging.
func ECSCompatible() Option {
	return func(cfg *config) {
		cfg.ecsCompatible = true
	}
}
