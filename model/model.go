package model

import (
	"github.com/gosmo-devs/microgateway/config"
)

// ServiceDao implementation
var ServiceDao ServiceDaoI

// Init initializes the databases configured
func Init() {
	switch config.Database {
	case "redis":
		ServiceDao = redisServiceDao()
	default:
		panic("Database configuration value not recognized")
	}
}
