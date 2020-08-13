package controller

import (
	"net/http"

	"github.com/gosmo-devs/microgateway/model"
	"github.com/gosmo-devs/microgateway/proxy"
)

// GetServices get services paginated
func GetServices(offset, limit int) (*model.ServicePage, error) {
	keys := model.ServiceDao.GetAllServiceKeys()
	if len(keys) == 0 || offset > len(keys) {
		return nil, model.ErrServiceNotFound
	}
	lowerIndex := offset
	upperIndex := min(offset+limit, len(keys))
	slicedKeys := keys[lowerIndex:upperIndex]
	if len(slicedKeys) == 0 {
		return nil, model.ErrServiceNotFound
	}
	services, err := model.ServiceDao.GetServices(slicedKeys...)
	if err != nil {
		return nil, err
	}
	servicePage := model.ServicePage{Services: services, TotalCount: len(keys)}
	return &servicePage, nil
}

// RegisterService adds a new service
func RegisterService(service model.Service) error {
	return model.ServiceDao.StoreService(service)
}

// GetService gets a service
func GetService(key string) (*model.Service, error) {
	return model.ServiceDao.GetService(key)
}

// DeleteService deletes a service
func DeleteService(key string) error {
	_, err := model.ServiceDao.GetService(key)
	if err != nil {
		return err
	}
	return model.ServiceDao.DeleteService(key)
}

// ReverseProxy forwards traffic to a service
func ReverseProxy(w http.ResponseWriter, r *http.Request, service *model.Service) error {
	p, err := proxy.NewProxy(service)
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
