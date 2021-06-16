package health

import (
	"context"
	"sync"
	"time"

	"github.com/gotway/gotway/internal/controller"
	"github.com/gotway/gotway/internal/model"

	"github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/health/client"
	"github.com/gotway/gotway/pkg/log"
)

type Health struct {
	serviceController controller.ServiceController

	logger log.Logger
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
	setToHealthy, setToIdle := h.getServicesToChangeStatus()

	for _, service := range setToHealthy {
		err := h.serviceController.UpdateServiceStatus(service, model.ServiceStatusHealthy)
		if err != nil {
			h.logger.Error(err)
		}
	}
	for _, service := range setToIdle {
		err := h.serviceController.UpdateServiceStatus(service, model.ServiceStatusIdle)
		if err != nil {
			h.logger.Error(err)
		}
	}
}

func (h *Health) getServicesToChangeStatus() (setToHealthy []string, setToIdle []string) {
	var healthyServices []string
	var idleServices []string

	services := h.serviceController.GetAllServiceKeys()
	var wg sync.WaitGroup
	for _, serviceKey := range services {
		wg.Add(1)

		go func(serviceKey string) {
			defer wg.Done()

			service, err := h.serviceController.GetService(serviceKey)
			if err != nil {
				h.logger.Errorf("unable to get service with key '%s'", serviceKey, err)
				return
			}
			client, err := client.New(service)
			if err != nil {
				h.logger.Error("error creating client", err)
				return
			}

			err = client.HealthCheck()
			if err != nil {
				if service.Status == model.ServiceStatusHealthy {
					h.logger.Infof("service %s is now idle. Cause: %v", service.Path, err)
					idleServices = append(idleServices, service.Path)
				}
			} else {
				if service.Status == model.ServiceStatusIdle {
					h.logger.Infof("Service %s is now healthy", service.Path)
					healthyServices = append(healthyServices, service.Path)
				}
			}

		}(serviceKey)
	}
	wg.Wait()
	return healthyServices, idleServices
}

func New(serviceController controller.ServiceController, logger log.Logger) *Health {
	return &Health{serviceController, logger}
}
