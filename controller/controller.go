package controller

import (
	"net/http"

	"github.com/gotway/gotway/core"
	"github.com/gotway/gotway/model"
)

// ServiceControllerI interface
type ServiceControllerI interface {
	GetServices(offset, limit int) (core.ServicePage, error)
	GetAllServiceKeys() []string
	RegisterService(serviceDetail core.ServiceDetail) error
	GetService(key string) (core.Service, error)
	GetServiceDetail(key string) (core.ServiceDetail, error)
	DeleteService(key string) error
	UpdateServiceStatus(key string, status core.ServiceStatus) error
	ReverseProxy(w http.ResponseWriter, r *http.Request, service core.Service) error
}

// CacheControllerI interface
type CacheControllerI interface {
	IsCacheableRequest(r *http.Request) bool
	GetCache(r *http.Request, pathPrefix, serviceKey string) (core.Cache, error)
	GetCacheDetail(r *http.Request, pathPrefix, serviceKey string) (core.CacheDetail, error)
	DeleteCacheByPath(paths []core.CachePath) error
	DeleteCacheByTags(tags []string) error
	ListenResponses()
	HandleResponse(serviceKey string, r *http.Response) error
}

// Service controller instance
var Service ServiceControllerI

// Cache controller instance
var Cache CacheControllerI

// Init initializes controllers with its dependencies
func Init() {
	Service = newServiceController(model.ServiceRepository)
	Cache = newCacheController(model.CacheRepository, model.ServiceRepository)
	Cache.ListenResponses()
}
