package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func New() *Logger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}

	zapLogger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
	}
}

func NewDevelopment() *Logger {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
	}
}
