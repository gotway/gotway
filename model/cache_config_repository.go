package model

import "github.com/gosmo-devs/microgateway/core"

// CacheConfigRepositoryI interface
type CacheConfigRepositoryI interface {
	StoreConfig(config core.CacheConfig, serviceKey string) error
	GetConfig(serviceKey string) (core.CacheConfig, error)
	DeleteConfig(serviceKey string) error
	IsCacheableStatusCode(serviceKey string, statusCode int) bool
}
