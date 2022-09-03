package logger

import (
	"log"
	"os"
	"path/filepath"

	zapcore "anilibrary-request-parser/pkg/logger/zap"
	"go.uber.org/zap"
)

type Zap struct {
	Logger *zap.Logger
	File   *os.File
}

func New(path string) (*Zap, error) {
	var zapLogger Zap

	path = filepath.Clean(path)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		log.Fatalf("error while creating file %s", err)
	}

	logger := zapcore.NewLogger(os.Stdout, file)

	zapLogger.Logger = logger
	zapLogger.File = file

	return &zapLogger, nil
}
