package logger

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestInfof(t *testing.T) {
	zapLogger := zaptest.NewLogger(t, zaptest.WrapOptions(
		zap.Hooks(func(e zapcore.Entry) error {
			assert.Equal(t, e.Level, zap.InfoLevel)
			return nil
		})))
	l := NewLogger(zapLogger)
	l.Infof("info message")
}

func TestErrorf(t *testing.T) {
	zapLogger := zaptest.NewLogger(t, zaptest.WrapOptions(
		zap.Hooks(func(e zapcore.Entry) error {
			assert.Equal(t, e.Level, zap.ErrorLevel)
			return nil
		})))
	l := NewLogger(zapLogger)
	l.Errorf("error message")
}
