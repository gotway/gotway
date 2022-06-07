package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gotway/gotway/internal/cache"
	cfg "github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/http"
	"github.com/gotway/gotway/internal/middleware"
	cacheMw "github.com/gotway/gotway/internal/middleware/cache"
	gatewayMw "github.com/gotway/gotway/internal/middleware/gateway"
	matchingressMw "github.com/gotway/gotway/internal/middleware/matchingress"
	"github.com/gotway/gotway/internal/repository"
	kubeCtrl "github.com/gotway/gotway/pkg/kubernetes/controller"
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

func getClientSet(
	config cfg.Config,
) (*clientsetv1alpha1.Clientset, error) {
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
	return clientsetv1alpha1.NewForConfig(restConfig)
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

	clientSet, err := getClientSet(config)
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
		clientSet,
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
		if err := kubeCtrl.Run(ctx); err != nil {
			logger.Fatalf("error starting Kubernetes controller: %v", err)
		}
	}()

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

	leader := newLeader(config, kubeCtrl, logger)
	if leader.hasFeaturesEnabled() {
		leader.start(ctx)
	}

	<-ctx.Done()
}
