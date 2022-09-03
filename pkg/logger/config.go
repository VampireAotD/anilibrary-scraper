package logger

import (
	"io"
	"os"
)

type Config struct {
	ConsoleOutput io.Writer
	LogFile       *os.File
}
