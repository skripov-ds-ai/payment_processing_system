package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	zap.Logger
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.Info(fmt.Sprintf(msg, args...))
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.Error(fmt.Sprintf(msg, args...))
}
