package model

import (
	"github.com/gosmo-devs/microgateway/config"
)

// ServiceDao implementation
var ServiceDao ServiceDaoI

// Init Initialize the databases configured
func Init() {
	switch config.Database {
	case "redis":
		ServiceDao = RedisServiceDao()
	default:
		panic("Database configuration value not recognized")
	}
}

// StoreService stores a service
func StoreService(key string, url string, healthURL string) {
	ServiceDao.StoreService(key, url, healthURL)
}
