package logger

import (
	"go.uber.org/zap"
)

type NopLogger struct {
	l *zap.Logger
}

func NewNopLogger(l *zap.Logger) *NopLogger {
	return &NopLogger{l: l}
}

func (z *NopLogger) Debug(msg string, args ...Field) {
	z.l.Debug(msg, z.toZapFields(args)...)
}

func (z *NopLogger) Info(msg string, args ...Field) {
	z.l.Info(msg, z.toZapFields(args)...)
}

func (z *NopLogger) Warn(msg string, args ...Field) {
	z.l.Warn(msg, z.toZapFields(args)...)
}

func (z *NopLogger) Error(msg string, args ...Field) {
	z.l.Error(msg, z.toZapFields(args)...)
}

func (z *NopLogger) toZapFields(args []Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
	}
	return res
}
