package main

import (
	"context"
	"fmt"

	goRedis "github.com/go-redis/redis/v8"
	cfg "github.com/gotway/gotway/internal/config"
	clientsetv1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1/apis/clientset/versioned"
	"github.com/gotway/gotway/pkg/redis"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type clientSets struct {
	gotway     *clientsetv1alpha1.Clientset
	kubernetes *kubernetes.Clientset
}

func getClientSets(config cfg.Config) (*clientSets, error) {
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

	gotwayClientSet, err := clientsetv1alpha1.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("error getting gotway clientset: %v", err)
	}
	kubeClientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("error getting kubernetes clientset: %v", err)
	}

	return &clientSets{
		gotway:     gotwayClientSet,
		kubernetes: kubeClientSet,
	}, nil
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
