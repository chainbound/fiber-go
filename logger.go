package client

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger(level string) *zap.SugaredLogger {
	var logger *zap.Logger
	zapLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}

	logger = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.Lock(os.Stdout),
		zapLevel,
	))

	return logger.Sugar()
}
