package health

import (
	"context"
	"time"

	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/service"

	"github.com/gotway/gotway/internal/health/client"
	"github.com/gotway/gotway/pkg/log"
)

type Options struct {
	CheckInterval time.Duration
	Timeout       time.Duration
	NumWorkers    int
	BufferSize    int
}

type Health struct {
	options           Options
	clientOptions     client.Options
	clientFactory     client.Factory
	pendingHealth     chan model.Service
	serviceController service.Controller
	logger            log.Logger
}

// Listen checks for service health periodically
func (h *Health) Listen(ctx context.Context) {
	h.logger.Info("starting health check")

	for i := 0; i < h.options.NumWorkers; i++ {
		go h.checkServices(ctx)
	}

	ticker := time.NewTicker(h.options.CheckInterval)
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("stopping health check")
			return
		case <-ticker.C:
			h.logger.Debug("checking health")
			services, err := h.serviceController.GetServices()
			if err != nil {
				h.logger.Error("error getting services ", err)
				return
			}
			for _, s := range services {
				h.pendingHealth <- s
			}
		}
	}
}

func (h *Health) checkServices(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case service := <-h.pendingHealth:
			h.updateService(service)
		}
	}
}

func (h *Health) updateService(service model.Service) {
	healthURL, err := service.HealthURL()
	if err != nil {
		h.logger.Error("error getting URL ", err)
		return
	}

	client, err := h.clientFactory.Get(service.Type, h.clientOptions)
	if err != nil {
		h.logger.Error("error getting client ", err)
		return
	}

	if err := client.HealthCheck(healthURL); err != nil {
		if service.Status == model.ServiceStatusHealthy {
			h.logger.Infof("service '%s' is now idle. Cause: %v", service.Name, err)
			service.Status = model.ServiceStatusIdle
			h.serviceController.UpsertService(service)
		}
	} else {
		if service.Status == model.ServiceStatusIdle {
			h.logger.Infof("service '%s' is now healthy", service.Name)
			service.Status = model.ServiceStatusHealthy
			h.serviceController.UpsertService(service)
		}
	}
}

func New(options Options, serviceController service.Controller, logger log.Logger) *Health {
	return &Health{
		options:           options,
		pendingHealth:     make(chan model.Service, options.BufferSize),
		serviceController: serviceController,
		clientOptions:     client.Options{Timeout: options.Timeout},
		clientFactory:     client.NewFactory(),
		logger:            logger,
	}
}
