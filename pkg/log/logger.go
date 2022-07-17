package log

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(production bool, logLevel, appName string) (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return nil, err
	}
	var config zap.Config
	if production {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	config.EncoderConfig.TimeKey = "@timestamp"
	config.EncoderConfig.MessageKey = "message"
	config.Level = level

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	if logger == nil {
		return nil, errors.New("error init logger")
	}
	logger = logger.With(zap.String("app", appName))
	return logger, nil
}
