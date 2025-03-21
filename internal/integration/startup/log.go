package startup

import (
	"go.uber.org/zap"
	"webook/pkg/logger"
)

func InitLogger() logger.LoggerV1 {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewNopLogger(l)
}
