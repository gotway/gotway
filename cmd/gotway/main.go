package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotway/gotway/internal/cache"
	cfg "github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/health"
	"github.com/gotway/gotway/internal/http"
	"github.com/gotway/gotway/internal/middleware"
	cacheMw "github.com/gotway/gotway/internal/middleware/cache"
	gatewayMw "github.com/gotway/gotway/internal/middleware/gateway"
	matchserviceMw "github.com/gotway/gotway/internal/middleware/matchservice"
	"github.com/gotway/gotway/internal/repository"
	"github.com/gotway/gotway/internal/service"
	kubernetesCtrl "github.com/gotway/gotway/pkg/kubernetes/controller"
	clientsetv1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1/apis/clientset/versioned"
	"github.com/gotway/gotway/pkg/log"
	"github.com/gotway/gotway/pkg/metrics"
	"github.com/gotway/gotway/pkg/pprof"
	"github.com/gotway/gotway/pkg/redis"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	goRedis "github.com/go-redis/redis/v8"
)

func configureMiddlewares(
	config cfg.Config,
	cacheController cache.Controller,
	serviceController service.Controller,
	logger log.Logger,
) []middleware.Middleware {

	middlewares := []middleware.Middleware{
		matchserviceMw.New(
			serviceController,
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

func getKubeClientSet(config cfg.Config) (*clientsetv1alpha1.Clientset, error) {
	var restConfig *rest.Config
	var err error
	if config.Kubernetes.KubeConfig != "" {
		restConfig, err = clientcmd.BuildConfigFromFlags("", config.Kubernetes.KubeConfig)
	} else {
		restConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}
	clientSet, err := clientsetv1alpha1.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

func getRedisClient(ctx context.Context, config cfg.Config) (redis.Cmdable, error) {
	opts, err := goRedis.ParseURL(config.RedisUrl)
	if err != nil {
		return nil, fmt.Errorf("error getting redis options %v", err)
	}
	client := goRedis.NewClient(opts)
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("error connecting to redis %v", err)
	}
	return redis.New(client), nil
}

func main() {
	config, err := cfg.GetConfig()
	if err != nil {
		panic(fmt.Errorf("error getting config %v", err))
	}
	ctx, _ := signal.NotifyContext(context.Background(), []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT}...,
	)

	logger := log.NewLogger(log.Fields{
		"service": "gotway",
	}, config.Env, config.LogLevel, os.Stdout)

	clientSet, err := getKubeClientSet(config)
	if err != nil {
		logger.Fatal("error getting kubernetes client set ", err)
	}
	redisClient, err := getRedisClient(ctx, config)
	if err != nil {
		logger.Fatal("error getting redis client: ", err)
	}

	kubeCtrl := kubernetesCtrl.New(
		kubernetesCtrl.Options{
			Namespace: config.Kubernetes.Namespace,
		},
		clientSet,
		logger.WithField("type", "kubernetes"),
	)
	go kubeCtrl.Run(ctx)

	serviceRepo := repository.NewServiceRepoRedis(redisClient)
	cacheRepo := repository.NewCacheRepoRedis(redisClient)

	serviceController := service.NewController(
		serviceRepo,
		logger.WithField("type", "service-ctrl"),
	)
	cacheController := cache.NewController(
		cache.Options{
			NumWorkers: config.Cache.NumWorkers,
			BufferSize: config.Cache.BufferSize,
		},
		cacheRepo,
		logger.WithField("type", "cache"),
	)
	if config.Cache.Enabled {
		go cacheController.ListenResponses(ctx)
	}

	if config.HealthCheck.Enabled {
		health := health.New(
			health.Options{
				CheckInterval: config.HealthCheck.Interval,
				Timeout:       config.HealthCheck.Timeout,
				NumWorkers:    config.HealthCheck.NumWorkers,
				BufferSize:    config.HealthCheck.BufferSize,
			},
			serviceController,
			logger.WithField("type", "health"),
		)
		go health.Listen(ctx)
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

	server := http.NewServer(
		http.ServerOptions{
			Port:       config.Port,
			TLSenabled: config.TLS.Enabled,
			TLScert:    config.TLS.Cert,
			TLSkey:     config.TLS.Key,
		},
		configureMiddlewares(
			config,
			cacheController,
			serviceController,
			logger.WithField("type", "middleware"),
		),
		cacheController,
		serviceController,
		logger.WithField("type", "http"),
	)
	go server.Start()
	defer server.Stop()

	<-ctx.Done()
}
