package logging

import "io"

type config struct {
	output        io.Writer
	logFiles      []io.Writer
	ecsCompatible bool
}

type Option func(cfg *config)

func WithOutput(output io.Writer) Option {
	return func(cfg *config) {
		cfg.output = output
	}
}

func WithLogFiles(files ...io.Writer) Option {
	return func(cfg *config) {
		cfg.logFiles = append(cfg.logFiles, files...)
	}
}

func ECSCompatible() Option {
	return func(cfg *config) {
		cfg.ecsCompatible = true
	}
}
