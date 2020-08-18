package log

import (
	"github.com/gosmo-devs/microgateway/config"
	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

func initZap() zapLogger {
	var zapConfig zap.Config
	if config.Env == "development" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}
	zapConfig.DisableCaller = true

	var logger *zap.Logger
	logger, _ = zapConfig.Build()
	defer logger.Sync()

	return zapLogger{logger.Sugar()}
}

func (z zapLogger) Debug(args ...interface{}) {
	z.logger.Debug(args...)
}

func (z zapLogger) Debugf(format string, args ...interface{}) {
	z.logger.Debugf(format, args...)
}

func (z zapLogger) Info(args ...interface{}) {
	z.logger.Info(args...)
}

func (z zapLogger) Infof(format string, args ...interface{}) {
	z.logger.Infof(format, args...)
}

func (z zapLogger) Warn(args ...interface{}) {
	z.logger.Warn(args...)
}

func (z zapLogger) Warnf(format string, args ...interface{}) {
	z.logger.Warnf(format, args...)
}

func (z zapLogger) Error(args ...interface{}) {
	z.logger.Error(args...)
}

func (z zapLogger) Errorf(format string, args ...interface{}) {
	z.logger.Errorf(format, args...)
}

func (z zapLogger) Fatal(args ...interface{}) {
	z.logger.Fatal(args...)
}

func (z zapLogger) Fatalf(format string, args ...interface{}) {
	z.logger.Fatalf(format, args...)
}
