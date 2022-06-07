package main

import (
	"context"

	cfg "github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/healthcheck"
	kubeCtrl "github.com/gotway/gotway/pkg/kubernetes/controller"
	"github.com/gotway/gotway/pkg/log"
)

type leader struct {
	config     cfg.Config
	kubeCtrl   *kubeCtrl.Controller
	healthCtrl *healthcheck.Controller
	logger     log.Logger
}

func (l *leader) start(ctx context.Context) {
	l.healthCtrl.Start(ctx)
}

func (l *leader) hasFeaturesEnabled() bool {
	return l.config.HealthCheck.Enabled
}

func newLeader(config cfg.Config, kubeCtrl *kubeCtrl.Controller, logger log.Logger) *leader {
	return &leader{
		config:   config,
		kubeCtrl: kubeCtrl,
		healthCtrl: healthcheck.NewController(
			healthcheck.Options{
				CheckInterval: config.HealthCheck.Interval,
				Timeout:       config.HealthCheck.Timeout,
				NumWorkers:    config.HealthCheck.NumWorkers,
				BufferSize:    config.HealthCheck.BufferSize,
			},
			kubeCtrl,
			logger,
		),
		logger: logger,
	}
}
