package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotway/gotway/internal/cache"
	"github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/health"
	"github.com/gotway/gotway/internal/http"
	"github.com/gotway/gotway/internal/repository"
	"github.com/gotway/gotway/internal/service"
	"github.com/gotway/gotway/pkg/log"
	"github.com/gotway/gotway/pkg/redis"

	goRedis "github.com/go-redis/redis/v8"
)

type stoppable interface {
	Stop()
}

func handleExit(
	logger log.Logger,
	sigs <-chan os.Signal,
	cancel context.CancelFunc,
	stoppables ...stoppable,
) {
	sig := <-sigs
	logger.Infof("received signal %s", sig.String())
	cancel()
	for _, s := range stoppables {
		s.Stop()
	}
}

func main() {
	signals := make(chan os.Signal)
	signal.Notify(
		signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)

	ctx, cancel := context.WithCancel(context.Background())

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

	options := http.ServerOptions{
		Port:       config.Port,
		TLSenabled: config.TLS,
		TLScert:    config.TLScert,
		TLSkey:     config.TLSkey,
	}
	s := http.NewServer(
		options,
		cacheController,
		serviceController,
		logger.WithField("type", "http"),
	)
	go s.Start()

	health := health.New(serviceController, logger.WithField("type", "health"))
	go health.Listen(ctx)

	handleExit(logger, signals, cancel, s)
}
