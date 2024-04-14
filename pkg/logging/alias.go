package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
)

var (
	String  = zap.String
	Error   = zap.Error
	Any     = zap.Any
	Int     = zap.Int
	Bool    = zap.Bool
	Float64 = zap.Float64
)

type (
	Logger = zap.Logger
	Level  = zapcore.Level
)
