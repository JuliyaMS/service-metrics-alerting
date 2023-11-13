package logger

import (
	"go.uber.org/zap"
)

var Logger = NewLogger()

func NewLogger() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger.Sugar()
}
