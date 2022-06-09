package main

import (
	"github.com/gotway/gotway/internal/cache"
	cfg "github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/middleware"
	cacheMw "github.com/gotway/gotway/internal/middleware/cache"
	gatewayMw "github.com/gotway/gotway/internal/middleware/gateway"
	matchingressMw "github.com/gotway/gotway/internal/middleware/matchingress"
	kubeCtrl "github.com/gotway/gotway/pkg/kubernetes/controller"
	"github.com/gotway/gotway/pkg/log"
)

func configureMiddlewares(
	config cfg.Config,
	kubeCtrl *kubeCtrl.Controller,
	cacheController cache.Controller,
	logger log.Logger,
) []middleware.Middleware {

	middlewares := []middleware.Middleware{
		matchingressMw.New(
			kubeCtrl,
			logger.WithField("middleware", "match-service"),
		),
	}
	if config.Cache.Enabled {
		middlewares = append(middlewares,
			cacheMw.NewCacheIn(
				cacheController,
				logger.WithField("middleware", "cache-in"),
			),
		)
	}
	middlewares = append(middlewares,
		gatewayMw.New(
			gatewayMw.GatewayOptions{Timeout: config.GatewayTimeout},
			logger.WithField("middleware", "gateway"),
		),
	)
	if config.Cache.Enabled {
		middlewares = append(middlewares,
			cacheMw.NewCacheOut(
				cacheController,
				logger.WithField("middleware", "cache-out"),
			),
		)
	}

	return middlewares
}
