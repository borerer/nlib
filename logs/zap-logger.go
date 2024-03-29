package logs

import (
	"github.com/borerer/nlib/configs"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func Init(config *configs.ServerConfig) error {
	var err error
	var level zap.AtomicLevel
	level, err = zap.ParseAtomicLevel(config.LogLevel)
	if err != nil {
		return err
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = level
	logger, err = cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	return nil
}

func GetZapLogger() *zap.Logger {
	return logger
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
