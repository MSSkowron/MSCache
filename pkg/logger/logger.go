package logger

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	l, _ := zap.NewProduction()
	logger = l.Sugar()
}

// Infof logs the provided arguments at [InfoLevel] using the provided format.
func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

// Infoln logs the provided arguments at [InfoLevel] using the provided format.
func Infoln(msg string) {
	logger.Infoln(msg)
}

// Errorf logs the provided arguments at [ErrorLevel] using the provided format.
func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

// Errorln logs the provided arguments at [ErrorLevel] using the provided format.
func Errorln(args ...any) {
	logger.Errorln(args...)
}
