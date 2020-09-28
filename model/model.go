package model

import (
	"github.com/gosmo-devs/microgateway/config"
)

// ServiceRepository implementation
var ServiceRepository ServiceRepositoryI

// CacheConfigRepository implementation
var CacheConfigRepository CacheConfigRepositoryI

// CacheRepository implementation
var CacheRepository CacheRepositoryI

// Init initializes the databases configured
func Init() {
	switch config.Database {
	case "redis":
		initRedisClient()
		ServiceRepository = serviceRepositoryRedis{}
		CacheConfigRepository = cacheConfigRepositoryRedis{}
		CacheRepository = cacheRepositoryRedis{}
	default:
		panic("Database configuration value not recognized")
	}
}
