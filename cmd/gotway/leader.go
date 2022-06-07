package main

import (
	"context"

	cfg "github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/healthcheck"
	kubeCtrl "github.com/gotway/gotway/pkg/kubernetes/controller"
	"github.com/gotway/gotway/pkg/log"
)

type leader struct {
	config   cfg.Config
	kubeCtrl *kubeCtrl.Controller
	logger   log.Logger
}

func (l *leader) run(ctx context.Context) {
	healthCtrl := healthcheck.NewController(
		healthcheck.Options{
			CheckInterval: l.config.HealthCheck.Interval,
			Timeout:       l.config.HealthCheck.Timeout,
			NumWorkers:    l.config.HealthCheck.NumWorkers,
			BufferSize:    l.config.HealthCheck.BufferSize,
		},
		l.kubeCtrl,
		l.logger,
	)
	if l.config.HealthCheck.Enabled {
		healthCtrl.Start(ctx)
	}
}

func newLeader(config cfg.Config, kubeCtrl *kubeCtrl.Controller, logger log.Logger) *leader {
	return &leader{
		config:   config,
		kubeCtrl: kubeCtrl,
		logger:   logger,
	}
}
