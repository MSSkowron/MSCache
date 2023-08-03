package logger

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	l, _ := zap.NewProduction()
	logger = l.Sugar()
}

func Info(args ...any) {
	logger.Info(args...)
}

func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

func Infoln(msg string) {
	logger.Infoln(msg)
}

func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

func Errorln(args ...any) {
	logger.Errorln(args...)
}
