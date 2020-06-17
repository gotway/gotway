package controller

import (
	"net/http"

	"github.com/gosmo-devs/microgateway/model"
	"github.com/gosmo-devs/microgateway/proxy"
)

// RegisterService adds a new service
func RegisterService(key string, url string, healthURL string) error {
	return model.ServiceDao.StoreService(key, url, healthURL)
}

// GetService gets a service
func GetService(key string) (*model.Service, error) {
	return model.ServiceDao.GetService(key)
}

// ReverseProxy forwards traffic to a service
func ReverseProxy(w http.ResponseWriter, r *http.Request, service *model.Service) error {
	p := proxy.NewProxy(service)
	return p.ReverseProxy(w, r)
}
