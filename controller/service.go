package controller

import (
	"errors"
	"net/http"

	"github.com/gotway/gotway/core"
	"github.com/gotway/gotway/model"
	"github.com/gotway/gotway/proxy"
)

// ServiceController controller
type ServiceController struct {
	serviceRepository     model.ServiceRepositoryI
	cacheConfigRepository model.CacheConfigRepositoryI
}

// NewServiceController creates a new service controller
func NewServiceController(serviceRepository model.ServiceRepositoryI,
	cacheConfigRepository model.CacheConfigRepositoryI) ServiceController {
	return ServiceController{
		serviceRepository:     serviceRepository,
		cacheConfigRepository: cacheConfigRepository,
	}
}

// GetServices get services paginated
func (c ServiceController) GetServices(offset, limit int) (core.ServicePage, error) {
	keys := c.serviceRepository.GetAllServiceKeys()
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
	err := c.serviceRepository.StoreService(serviceDetail.Service)
	if err != nil {
		return err
	}
	if !serviceDetail.Cache.IsEmpty() {
		return c.cacheConfigRepository.StoreConfig(serviceDetail.Cache, serviceDetail.Service.Path)
	}
	return nil
}

// GetService gets a service
func (c ServiceController) GetService(key string) (core.Service, error) {
	service, err := c.serviceRepository.GetService(key)
	if err != nil {
		return core.Service{}, err
	}
	return service, nil
}

// GetServiceDetail gets a service
func (c ServiceController) GetServiceDetail(key string) (core.ServiceDetail, error) {
	service, err := c.GetService(key)
	if err != nil {
		return core.ServiceDetail{}, err
	}

	config, err := c.cacheConfigRepository.GetConfig(key)
	if err != nil {
		if !errors.Is(err, core.ErrCacheNotFound) {
			return core.ServiceDetail{}, err
		}
		config = core.DefaultCacheConfig
	}

	return core.ServiceDetail{
		Service: service,
		Cache:   config,
	}, nil
}

// DeleteService deletes a service
func (c ServiceController) DeleteService(key string) error {
	err := c.serviceRepository.ExistService(key)
	if err != nil {
		return err
	}
	err = c.serviceRepository.DeleteService(key)
	if err != nil {
		return err
	}
	return c.cacheConfigRepository.DeleteConfig(key)
}

// UpdateServiceStatus updates the status of a service
func (c ServiceController) UpdateServiceStatus(key string, status core.ServiceStatus) {
	c.serviceRepository.UpdateServiceStatus(key, status)
}

// ReverseProxy forwards traffic to a service
func (c ServiceController) ReverseProxy(w http.ResponseWriter, r *http.Request, service core.Service) error {
	p, err := proxy.NewProxy(service, Cache.handleResponse)
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
