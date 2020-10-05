package controller

import "github.com/gotway/gotway/model"

// Service controller instance
var Service ServiceControllerI

// Cache controller instance
var Cache CacheControllerI

// Init initializes controllers with its dependencies
func Init() {
	Service = NewServiceController(model.ServiceRepository, model.CacheConfigRepository)
	Cache = NewCacheController(model.CacheConfigRepository, model.CacheRepository)
	Cache.ListenResponses()
}
