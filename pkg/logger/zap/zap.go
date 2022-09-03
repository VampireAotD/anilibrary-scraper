package zap

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const timeEncoderLayout string = "02/01/2006 15:04:05"

type Logger struct {
	*zap.Logger
}

func NewLogger(console io.Writer, files ...io.Writer) Logger {
	pe := zap.NewProductionEncoderConfig()

	// file
	pe.TimeKey = "time"
	pe.EncodeTime = zapcore.TimeEncoderOfLayout(timeEncoderLayout)
	pe.CallerKey = "file"
	pe.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(level.CapitalString())
	}
	pe.MessageKey = "log"
	fileEncoder := zapcore.NewJSONEncoder(pe)

	// console
	pe.EncodeTime = zapcore.TimeEncoderOfLayout(timeEncoderLayout)
	pe.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString("|")
		encoder.AppendString(level.CapitalString())
		encoder.AppendString("|")
	}
	pe.ConsoleSeparator = " "
	pe.EncodeName = func(s string, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(s)
		encoder.AppendString("|")
	}
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	cores := make([]zapcore.Core, len(files)+1, len(files)+1)

	// console
	cores[0] = zapcore.NewCore(consoleEncoder,
		zapcore.AddSync(console),
		zap.DebugLevel,
	)

	for i := range files {
		core := zapcore.NewCore(fileEncoder,
			zapcore.AddSync(files[i]),
			zap.DebugLevel,
		)

		cores[i+1] = core
	}

	return Logger{zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
	)}
}
