package health

import (
	"context"
	"sync"
	"time"

	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/service"

	"github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/health/client"
	"github.com/gotway/gotway/pkg/log"
)

type Health struct {
	serviceController service.Controller
	clientFactory     client.Factory
	logger            log.Logger
}

// Listen checks for service health periodically
func (h *Health) Listen(ctx context.Context) {
	ticker := time.NewTicker(config.HealthCheckInterval)

	h.logger.Info("starting health check")
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("stopping health check")
			return
		case <-ticker.C:
			h.logger.Debug("checking health")
			h.updateServiceStatus()
		}
	}
}

func (h *Health) updateServiceStatus() {
	statusUpdate := h.getStatusUpdate()
	if statusUpdate == nil {
		return
	}
	for status, services := range statusUpdate.Get() {
		if err := h.serviceController.UpdateServicesStatus(status, services...); err != nil {
			h.logger.Errorf("error updating services %v status to '%s' %v", services, status, err)
		}
	}
}

func (h *Health) getStatusUpdate() *statusUpdate {
	services := h.serviceController.GetAllServiceKeys()
	statusUpdate := NewStatusUpdate()

	var wg sync.WaitGroup
	wg.Add(len(services))

	for _, serviceKey := range services {
		go func(serviceKey string) {
			defer wg.Done()

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

			client, err := h.clientFactory.GetClient(service.Type)
			if err != nil {
				h.logger.Error("error getting client ", err)
				return
			}

			if err := client.HealthCheck(healthURL); err != nil {
				if service.Status == model.ServiceStatusHealthy {
					h.logger.Infof("service %s is now idle. Cause: %v", service.Path, err)
					statusUpdate.Add(model.ServiceStatusIdle, service.Path)
				}
			} else {
				if service.Status == model.ServiceStatusIdle {
					h.logger.Infof("service %s is now healthy", service.Path)
					statusUpdate.Add(model.ServiceStatusHealthy, service.Path)
				}
			}

		}(serviceKey)
	}
	wg.Wait()

	return statusUpdate
}

func New(serviceController service.Controller, logger log.Logger) *Health {
	return &Health{
		serviceController: serviceController,
		clientFactory:     client.NewFactory(),
		logger:            logger,
	}
}
