package controller

import (
	"net/http"

	"github.com/gosmo-devs/microgateway/model"
	"github.com/gosmo-devs/microgateway/proxy"
)

// RegisterService adds a new service
func RegisterService(service model.Service) error {
	return model.ServiceDao.StoreService(service)
}

// GetService gets a service
func GetService(key string) (*model.Service, error) {
	return model.ServiceDao.GetService(key)
}

// ReverseProxy forwards traffic to a service
func ReverseProxy(w http.ResponseWriter, r *http.Request, service *model.Service) error {
	p, err := proxy.NewProxy(service)
	if err != nil {
		return err
	}
	return p.ReverseProxy(w, r)
}
