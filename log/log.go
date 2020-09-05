package log

import (
	"github.com/gosmo-devs/microgateway/config"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// Init initializes the logger instance
func Init() {
	var zapLogger *zap.Logger
	var zapConfig zap.Config
	if config.Env == "development" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}
	zapConfig.DisableCaller = true
	zapLogger, _ = zapConfig.Build()
	defer zapLogger.Sync()
	logger = zapLogger.Sugar()
}

// Debug log level
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf log level with format
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Info log level
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof log level with format
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warn log level
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf log level with format
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Error log level
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf log level with format
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatal log level
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Fatalf log level with format
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
