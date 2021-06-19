package main

import (
	"context"
	"os"

	"github.com/gotway/gotway/internal/cache"
	"github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/health"
	"github.com/gotway/gotway/internal/http"
	"github.com/gotway/gotway/internal/repository"
	"github.com/gotway/gotway/internal/service"
	gs "github.com/gotway/gotway/pkg/graceful_shutdown"
	"github.com/gotway/gotway/pkg/log"
	"github.com/gotway/gotway/pkg/metrics"
	"github.com/gotway/gotway/pkg/redis"

	goRedis "github.com/go-redis/redis/v8"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	stoppables := []gs.Stoppable{}

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
		metricsOptions := metrics.MetricsOptions{
			Path: config.MetricsPath,
			Port: config.MetricsPort,
		}
		m := metrics.New(metricsOptions, logger.WithField("type", "metrics"))
		go m.Start()
		stoppables = append(stoppables, m)
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

	httpOptions := http.ServerOptions{
		Port:       config.Port,
		TLSenabled: config.TLS,
		TLScert:    config.TLScert,
		TLSkey:     config.TLSkey,
	}
	s := http.NewServer(
		httpOptions,
		cacheController,
		serviceController,
		logger.WithField("type", "http"),
	)
	go s.Start()
	stoppables = append(stoppables, s)

	health := health.New(serviceController, logger.WithField("type", "health"))
	go health.Listen(ctx)

	gs.GracefulShutdown(cancel, stoppables...)
}