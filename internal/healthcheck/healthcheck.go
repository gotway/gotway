package healthcheck

import (
	"context"
	"fmt"
	"net/url"
	"time"

	kubernetesCtrl "github.com/gotway/gotway/pkg/kubernetes/controller"
	crdv1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1"
	"github.com/gotway/gotway/pkg/log"
)

type Options struct {
	CheckInterval time.Duration
	Timeout       time.Duration
	NumWorkers    int
	BufferSize    int
}

type Controller struct {
	options       Options
	client        client
	pendingHealth chan crdv1alpha1.IngressHTTP
	kubeCtrl      *kubernetesCtrl.Controller
	logger        log.Logger
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
			services, err := c.kubeCtrl.ListIngresses()
			if err != nil {
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
			c.updateService(ctx, service)
		}
	}
}

func (c *Controller) updateService(ctx context.Context, ingress crdv1alpha1.IngressHTTP) {
	healthURL, err := getHealthUrl(ingress)
	if err != nil {
		c.logger.Error("error getting health url ", err)
		return
	}
	if err := c.client.healthCheck(healthURL); err != nil {
		if ingress.Status.IsServiceHealthy {
			c.logger.Infof("service '%s' is idle: %v", ingress.Spec.Service.Name, err)
			if err := c.kubeCtrl.UpdateIngressStatus(ctx, ingress, false); err != nil {
				c.logger.Errorf("error updating service '%s': %v", ingress.Spec.Service.Name, err)
			}
		}
	} else {
		if !ingress.Status.IsServiceHealthy {
			c.logger.Infof("service '%s' is healthy", ingress.Spec.Service.Name)
			if err := c.kubeCtrl.UpdateIngressStatus(ctx, ingress, true); err != nil {
				c.logger.Errorf("error updating service '%s' status: %v", ingress.Spec.Service.Name, err)
			}
		}
	}
}

func getHealthUrl(ingress crdv1alpha1.IngressHTTP) (*url.URL, error) {
	healthPath := ingress.Spec.Service.HealthPath
	if healthPath == "" {
		healthPath = "/health"
	}
	return url.Parse(fmt.Sprintf("%s%s", ingress.Spec.Service.URL, healthPath))
}

func NewController(
	options Options,
	kubeCtrl *kubernetesCtrl.Controller,
	logger log.Logger,
) *Controller {
	return &Controller{
		options:       options,
		client:        newClient(clientOptions{timeout: options.Timeout}),
		pendingHealth: make(chan crdv1alpha1.IngressHTTP, options.BufferSize),
		kubeCtrl:      kubeCtrl,
		logger:        logger,
	}
}
