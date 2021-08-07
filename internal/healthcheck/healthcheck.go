package healthcheck

import (
	"context"
	"time"

	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/service"

	"github.com/gotway/gotway/pkg/log"
)

type Options struct {
	CheckInterval time.Duration
	Timeout       time.Duration
	NumWorkers    int
	BufferSize    int
}

type Controller struct {
	options           Options
	client            client
	pendingHealth     chan model.Service
	serviceController service.Controller
	logger            log.Logger
}

// Start checks for service health periodically
func (c *Controller) Start(ctx context.Context) {
	c.logger.Info("starting health check")

	for i := 0; i < c.options.NumWorkers; i++ {
		go c.checkServices(ctx)
	}

	ticker := time.NewTicker(c.options.CheckInterval)
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("stopping health check")
			return
		case <-ticker.C:
			c.logger.Debug("checking health")
			services, err := c.serviceController.GetServices()
			if err != nil && err != model.ErrServiceNotFound {
				c.logger.Error("error getting services ", err)
				continue
			}
			for _, s := range services {
				c.pendingHealth <- s
			}
		}
	}
}

func (c *Controller) checkServices(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case service := <-c.pendingHealth:
			c.updateService(service)
		}
	}
}

func (c *Controller) updateService(service model.Service) {
	healthURL, err := service.HealthURL()
	if err != nil {
		c.logger.Error("error getting URL ", err)
		return
	}

	if err := c.client.healthCheck(healthURL); err != nil {
		if service.Status == model.ServiceStatusHealthy {
			c.logger.Infof("service '%s' is now idle. Cause: %v", service.ID, err)
			service.Status = model.ServiceStatusIdle
			c.serviceController.UpsertService(service)
		}
	} else {
		if service.Status == model.ServiceStatusIdle {
			c.logger.Infof("service '%s' is now healthy", service.ID)
			service.Status = model.ServiceStatusHealthy
			c.serviceController.UpsertService(service)
		}
	}
}

func NewController(options Options, serviceController service.Controller, logger log.Logger) *Controller {
	return &Controller{
		options:           options,
		client:            newClient(clientOptions{timeout: options.Timeout}),
		pendingHealth:     make(chan model.Service, options.BufferSize),
		serviceController: serviceController,
		logger:            logger,
	}
}
