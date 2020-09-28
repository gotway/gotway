package log

import (
	"github.com/gosmo-devs/microgateway/config"
	"go.uber.org/zap"
)

func initZap() LoggerI {
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

	return logger.Sugar()
}
