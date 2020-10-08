package controller

import (
	"net/http"

	"github.com/gotway/gotway/core"
	"github.com/gotway/gotway/model"
	"github.com/gotway/gotway/proxy"
)

// ServiceController controller
type ServiceController struct {
	serviceRepository model.ServiceRepositoryI
}

func newServiceController(serviceRepository model.ServiceRepositoryI) ServiceController {
	return ServiceController{
		serviceRepository: serviceRepository,
	}
}

// GetServices get services paginated
func (c ServiceController) GetServices(offset, limit int) (core.ServicePage, error) {
	keys := c.GetAllServiceKeys()
	if len(keys) == 0 || offset > len(keys) {
		return core.ServicePage{}, core.ErrServiceNotFound
	}

	lowerIndex := offset
	upperIndex := min(offset+limit, len(keys))
	slicedKeys := keys[lowerIndex:upperIndex]
	if len(slicedKeys) == 0 {
		return core.ServicePage{}, core.ErrServiceNotFound
	}

	services, err := c.serviceRepository.GetServices(slicedKeys...)
	if err != nil {
		return core.ServicePage{}, err
	}

	return core.ServicePage{Services: services, TotalCount: len(keys)}, nil
}

// GetAllServiceKeys retrieves all service keys
func (c ServiceController) GetAllServiceKeys() []string {
	return c.serviceRepository.GetAllServiceKeys()
}

// RegisterService adds a new service
func (c ServiceController) RegisterService(serviceDetail core.ServiceDetail) error {
	return c.serviceRepository.StoreService(serviceDetail)
}

// GetService gets a service
func (c ServiceController) GetService(key string) (core.Service, error) {
	return c.serviceRepository.GetService(key)
}

// GetServiceDetail gets a service with extra info
func (c ServiceController) GetServiceDetail(key string) (core.ServiceDetail, error) {
	return c.serviceRepository.GetServiceDetail(key)
}

// DeleteService deletes a service
func (c ServiceController) DeleteService(key string) error {
	return c.serviceRepository.DeleteService(key)
}

// UpdateServiceStatus updates the status of a service
func (c ServiceController) UpdateServiceStatus(key string, status core.ServiceStatus) error {
	return c.serviceRepository.UpdateServiceStatus(key, status)
}

// ReverseProxy forwards traffic to a service
func (c ServiceController) ReverseProxy(w http.ResponseWriter, r *http.Request, service core.Service) error {
	p, err := proxy.NewProxy(service, Cache.HandleResponse)
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
