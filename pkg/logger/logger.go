package logger

import (
	"anilibrary-request-parser/pkg/logger/zap"
)

func New(cfg Config) (zap.Logger, error) {
	return zap.NewLogger(cfg.ConsoleOutput, cfg.LogFile), nil
}
