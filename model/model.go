package model

import (
	"github.com/gotway/gotway/config"
	"github.com/gotway/gotway/core"
)

// ServiceRepositoryI interface
type ServiceRepositoryI interface {
	StoreService(service core.ServiceDetail) error
	GetAllServiceKeys() []string
	GetService(key string) (core.Service, error)
	GetServiceDetail(key string) (core.ServiceDetail, error)
	GetServices(keys ...string) ([]core.Service, error)
	DeleteService(key string) error
	UpdateServiceStatus(key string, status core.ServiceStatus) error
	GetServiceCache(key string) (core.CacheConfig, error)
	IsCacheableStatusCode(key string, statusCode int) bool
}

// CacheRepositoryI implementation
type CacheRepositoryI interface {
	StoreCache(cache core.CacheDetail, serviceKey string) error
	GetCache(path string, serviceKey string) (core.Cache, error)
	GetCacheDetail(path string, serviceKey string) (core.CacheDetail, error)
	DeleteCacheByPath(paths []core.CachePath) error
	DeleteCacheByTags(tags []string) error
}

// ServiceRepository instance
var ServiceRepository ServiceRepositoryI

// CacheRepository instance
var CacheRepository CacheRepositoryI

// Init initializes the repositories
func Init() {
	client := newRedisClient(config.RedisServer)
	ServiceRepository = newServiceRepositoryRedis(client)
	CacheRepository = newCacheRepositoryRedis(client)
}
