package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotway/gotway/internal/cache"
	cfg "github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/healthcheck"
	"github.com/gotway/gotway/internal/http"
	"github.com/gotway/gotway/internal/repository"
	kubeCtrl "github.com/gotway/gotway/pkg/kubernetes/controller"
	"github.com/gotway/gotway/pkg/kubernetes/leaderelection"
	"github.com/gotway/gotway/pkg/log"
	"github.com/gotway/gotway/pkg/metrics"
	"github.com/gotway/gotway/pkg/pprof"
)

func main() {
	config, err := cfg.GetConfig()
	if err != nil {
		panic(fmt.Errorf("error getting config %v", err))
	}
	logger := log.NewLogger(log.Fields{
		"service": "gotway",
	}, config.Env, config.LogLevel, os.Stdout)
	ctx, _ := signal.NotifyContext(context.Background(), []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT}...,
	)

	clientSets, err := getClientSets(config)
	if err != nil {
		logger.Fatal("error getting kubernetes client set ", err)
	}
	redisClient, err := getRedisClient(ctx, config)
	if err != nil {
		logger.Fatal("error getting redis client: ", err)
	}

	kubeCtrl := kubeCtrl.New(
		kubeCtrl.Options{
			Namespace:    config.Kubernetes.Namespace,
			ResyncPeriod: config.Kubernetes.ResyncPeriod,
		},
		clientSets.gotway,
		logger.WithField("type", "kubernetes"),
	)

	cacheRepo := repository.NewCacheRepoRedis(redisClient)
	cacheCtrl := cache.NewController(
		cache.Options{
			NumWorkers: config.Cache.NumWorkers,
			BufferSize: config.Cache.BufferSize,
		},
		cacheRepo,
		logger.WithField("type", "cache"),
	)
	if config.Cache.Enabled {
		go cacheCtrl.Start(ctx)
	}

	if config.Metrics.Enabled {
		m := metrics.New(
			metrics.Options{
				Path: config.Metrics.Path,
				Port: config.Metrics.Port,
			},
			logger.WithField("type", "metrics"),
		)
		go m.Start()
		defer m.Stop()
	}

	if config.PProf.Enabled {
		p := pprof.New(
			pprof.Options{Port: config.PProf.Port},
			logger.WithField("type", "pprof"),
		)
		go p.Start()
		defer p.Stop()
	}

	go func() {
		if err := kubeCtrl.Start(ctx); err != nil {
			logger.Fatalf("error starting Kubernetes controller: %v", err)
		}
	}()

	if config.HealthCheck.Enabled {
		healthCtrl := healthcheck.NewController(
			healthcheck.Options{
				CheckInterval: config.HealthCheck.Interval,
				Timeout:       config.HealthCheck.Timeout,
				NumWorkers:    config.HealthCheck.NumWorkers,
				BufferSize:    config.HealthCheck.BufferSize,
			},
			kubeCtrl,
			logger,
		)

		if config.HA.Enabled {
			identity, err := os.Hostname()
			if err != nil {
				logger.Fatalf("error getting hostname: %v", err)
			}
			instanceLogger := logger.WithField("instance", identity)

			le := leaderelection.New(kubeCtrl, clientSets.kubernetes, instanceLogger,
				leaderelection.Config{
					Identity:           identity,
					LeaseLockName:      config.HA.LeaseLockName,
					LeaseLockNamespace: config.HA.LeaseLockNamespace,
					LeaseDuration:      config.HA.LeaseDuration,
					RenewDeadline:      config.HA.RenewDeadline,
					RetryPeriod:        config.HA.RetryPeriod,
					OnStartedLeading: func(ctx context.Context) {
						instanceLogger.Info("started leading")
						healthCtrl.Start(ctx)
					},
					OnStoppedLeading: func() {
						instanceLogger.Info("stopped leading")
					},
				},
			)
			go le.Start(ctx)
		} else {
			go healthCtrl.Start(ctx)
		}
	}

	server := http.NewServer(
		http.ServerOptions{
			Port:       config.Port,
			TLSenabled: config.TLS.Enabled,
			TLScert:    config.TLS.Cert,
			TLSkey:     config.TLS.Key,
		},
		configureMiddlewares(
			config,
			kubeCtrl,
			cacheCtrl,
			logger.WithField("type", "middleware"),
		),
		kubeCtrl,
		cacheCtrl,
		logger.WithField("type", "http"),
	)
	go server.Start()
	defer server.Stop()

	<-ctx.Done()
}
