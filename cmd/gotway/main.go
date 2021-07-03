package main

import (
	"context"
	"os"

	"github.com/gotway/gotway/internal/cache"
	"github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/health"
	"github.com/gotway/gotway/internal/http"
	"github.com/gotway/gotway/internal/middleware"
	cacheMw "github.com/gotway/gotway/internal/middleware/cache"
	gatewayMw "github.com/gotway/gotway/internal/middleware/gateway"
	matchserviceMw "github.com/gotway/gotway/internal/middleware/matchservice"
	"github.com/gotway/gotway/internal/repository"
	"github.com/gotway/gotway/internal/service"
	gs "github.com/gotway/gotway/pkg/graceful_shutdown"
	"github.com/gotway/gotway/pkg/log"
	"github.com/gotway/gotway/pkg/metrics"
	"github.com/gotway/gotway/pkg/pprof"
	"github.com/gotway/gotway/pkg/redis"

	goRedis "github.com/go-redis/redis/v8"
)

func configureMiddlewares(
	cacheController cache.Controller,
	serviceController service.Controller,
	logger log.Logger,
) []middleware.Middleware {

	return []middleware.Middleware{
		matchserviceMw.New(
			serviceController,
			logger.WithField("subtype", "match-service"),
		),
		cacheMw.NewCacheIn(
			cacheController,
			logger.WithField("subtype", "cache-in"),
		),
		gatewayMw.New(
			gatewayMw.GatewayOptions{Timeout: config.GatewayTimeout},
			logger.WithField("subtype", "gateway"),
		),
		cacheMw.NewCacheOut(
			cacheController,
			logger.WithField("subtype", "cache-out"),
		),
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	shutdownHooks := []gs.ShutdownHook{}

	logger := log.NewLogger(log.Fields{
		"service": "gotway",
	}, config.Env, config.LogLevel, os.Stdout)

	opts, err := goRedis.ParseURL(config.RedisUrl)
	if err != nil {
		logger.Fatal("error getting redis options ", err)
	}
	client := goRedis.NewClient(opts)
	defer client.Close()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Fatal("error connecting to redis")
	}
	redisClient := redis.New(client)
	logger.Info("connected to redis")

	if config.Metrics {
		m := metrics.New(
			metrics.Options{
				Path: config.MetricsPath,
				Port: config.MetricsPort,
			},
			logger.WithField("type", "metrics"),
		)
		go m.Start()
		shutdownHooks = append(shutdownHooks, m.Stop)
	}

	if config.PProf {
		p := pprof.New(
			pprof.Options{Port: config.PProfPort},
			logger.WithField("type", "pprof"),
		)
		go p.Start()
		shutdownHooks = append(shutdownHooks, p.Stop)
	}

	serviceRepo := repository.NewServiceRepoRedis(redisClient)
	cacheRepo := repository.NewCacheRepoRedis(redisClient)

	serviceController := service.NewController(
		serviceRepo,
		logger.WithField("type", "service-ctrl"),
	)
	cacheController := cache.NewController(
		cacheRepo,
		serviceRepo,
		logger.WithField("type", "cache-ctrl"),
	)
	go cacheController.ListenResponses(ctx)

	server := http.NewServer(
		http.ServerOptions{
			Port:       config.Port,
			TLSenabled: config.TLS,
			TLScert:    config.TLScert,
			TLSkey:     config.TLSkey,
		},
		configureMiddlewares(
			cacheController,
			serviceController,
			logger.WithField("type", "middleware"),
		),
		cacheController,
		serviceController,
		logger.WithField("type", "http"),
	)
	go server.Start()
	shutdownHooks = append(shutdownHooks, server.Stop)

	health := health.New(
		health.Options{
			CheckInterval: config.HealthCheckInterval,
			Timeout:       config.HealthCheckTimeout,
			NumWorkers:    config.HealthNumWorkers,
			BufferSize:    config.HealthBufferSize,
		},
		serviceController,
		logger.WithField("type", "health"),
	)
	go health.Listen(ctx)

	gs.GracefulShutdown(logger, cancel, shutdownHooks...)
}
