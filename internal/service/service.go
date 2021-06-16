package service

import (
	"net/http"

	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/proxy"
	"github.com/gotway/gotway/internal/repository"
	"github.com/gotway/gotway/pkg/log"
)

type Controller interface {
	GetServices(offset, limit int) (model.ServicePage, error)
	GetAllServiceKeys() []string
	RegisterService(serviceDetail model.ServiceDetail) error
	GetService(key string) (model.Service, error)
	GetServiceDetail(key string) (model.ServiceDetail, error)
	DeleteService(key string) error
	UpdateServiceStatus(key string, status model.ServiceStatus) error
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

// GetServices get services paginated
func (c BasicController) GetServices(offset, limit int) (model.ServicePage, error) {
	keys := c.GetAllServiceKeys()
	if len(keys) == 0 || offset > len(keys) {
		return model.ServicePage{}, model.ErrServiceNotFound
	}

	lowerIndex := offset
	upperIndex := min(offset+limit, len(keys))
	slicedKeys := keys[lowerIndex:upperIndex]
	if len(slicedKeys) == 0 {
		return model.ServicePage{}, model.ErrServiceNotFound
	}

	services, err := c.serviceRepo.GetServices(slicedKeys...)
	if err != nil {
		return model.ServicePage{}, err
	}

	return model.ServicePage{Services: services, TotalCount: len(keys)}, nil
}

// GetAllServiceKeys retrieves all service keys
func (c BasicController) GetAllServiceKeys() []string {
	return c.serviceRepo.GetAllServiceKeys()
}

// RegisterService adds a new service
func (c BasicController) RegisterService(serviceDetail model.ServiceDetail) error {
	return c.serviceRepo.StoreService(serviceDetail)
}

// GetService gets a service
func (c BasicController) GetService(key string) (model.Service, error) {
	return c.serviceRepo.GetService(key)
}

// GetServiceDetail gets a service with extra info
func (c BasicController) GetServiceDetail(key string) (model.ServiceDetail, error) {
	return c.serviceRepo.GetServiceDetail(key)
}

// DeleteService deletes a service
func (c BasicController) DeleteService(key string) error {
	return c.serviceRepo.DeleteService(key)
}

// UpdateServiceStatus updates the status of a service
func (c BasicController) UpdateServiceStatus(key string, status model.ServiceStatus) error {
	return c.serviceRepo.UpdateServiceStatus(key, status)
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
