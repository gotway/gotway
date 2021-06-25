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
	serviceChan       chan string
	serviceController service.Controller
	clientFactory     client.Factory
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
			for _, s := range h.serviceController.GetAllServiceKeys() {
				h.serviceChan <- s
			}
		}
	}
}

func (h *Health) checkServices(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case serviceKey := <-h.serviceChan:
			h.updateService(serviceKey)
		}
	}
}

func (h *Health) updateService(serviceKey string) {
	service, err := h.serviceController.GetService(serviceKey)
	if err != nil {
		h.logger.Errorf("unable to get service with key '%s'", serviceKey, err)
		return
	}

	healthURL, err := service.HealthURL()
	if err != nil {
		h.logger.Error("error getting URL ", err)
		return
	}

	client, err := h.clientFactory.GetClient(service.Type, h.clientOptions)
	if err != nil {
		h.logger.Error("error getting client ", err)
		return
	}

	if err := client.HealthCheck(healthURL); err != nil {
		if service.Status == model.ServiceStatusHealthy {
			h.logger.Infof("service %s is now idle. Cause: %v", service.Path, err)
			h.serviceController.UpdateServicesStatus(model.ServiceStatusIdle, service.Path)
		}
	} else {
		if service.Status == model.ServiceStatusIdle {
			h.logger.Infof("service %s is now healthy", service.Path)
			h.serviceController.UpdateServicesStatus(model.ServiceStatusHealthy, service.Path)
		}
	}
}

func New(options Options, serviceController service.Controller, logger log.Logger) *Health {
	return &Health{
		options:           options,
		serviceChan:       make(chan string, options.BufferSize),
		serviceController: serviceController,
		clientOptions:     client.Options{Timeout: options.Timeout},
		clientFactory:     client.NewFactory(),
		logger:            logger,
	}
}
