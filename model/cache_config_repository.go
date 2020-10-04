package model

import "github.com/gotway/gotway/core"

// CacheConfigRepositoryI interface
type CacheConfigRepositoryI interface {
	StoreConfig(config core.CacheConfig, serviceKey string) error
	GetConfig(serviceKey string) (core.CacheConfig, error)
	DeleteConfig(serviceKey string) error
	IsCacheableStatusCode(serviceKey string, statusCode int) bool
}
