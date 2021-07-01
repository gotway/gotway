package service

import (
	"net/http"

	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/proxy"
	"github.com/gotway/gotway/internal/repository"
	"github.com/gotway/gotway/pkg/log"
)

type Controller interface {
	CreateService(service model.Service) error
	GetServices() ([]model.Service, error)
	GetService(key string) (model.Service, error)
	DeleteService(key string) error
	UpsertService(service model.Service) error
	ReverseProxy(
		w http.ResponseWriter,
		r *http.Request,
		service model.Service,
		handler proxy.ResponseHandler,
	) error
}

type BasicController struct {
	serviceRepo repository.ServiceRepo
	logger      log.Logger
}

// CreateService creates a new service
func (c BasicController) CreateService(service model.Service) error {
	return c.serviceRepo.Create(service)
}

// GetServices gets services paginated
func (c BasicController) GetServices() ([]model.Service, error) {
	return c.serviceRepo.GetAll()
}

// GetService gets a service
func (c BasicController) GetService(key string) (model.Service, error) {
	return c.serviceRepo.Get(key)
}

// DeleteService deletes a service
func (c BasicController) DeleteService(key string) error {
	return c.serviceRepo.Delete(key)
}

// UpdateServiceStatus updates the status of a service
func (c BasicController) UpsertService(service model.Service) error {
	return c.serviceRepo.Upsert(service)
}

// ReverseProxy forwards traffic to a service
func (c BasicController) ReverseProxy(
	w http.ResponseWriter,
	r *http.Request,
	service model.Service,
	handler proxy.ResponseHandler,
) error {
	p, err := proxy.New(service, handler, c.logger.WithField("type", "proxy"))
	if err != nil {
		return err
	}
	return p.ReverseProxy(w, r)
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func NewController(
	serviceRepo repository.ServiceRepo,
	logger log.Logger,
) Controller {
	return BasicController{serviceRepo, logger}
}
