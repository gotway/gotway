package model

import (
	"github.com/gosmo-devs/microgateway/config"
)

// ServiceDao implementation
var ServiceDao ServiceDaoI

// Init Initialize the databases configured
func Init() {
	initDatabase()
	initHealthcheck()
}

func initDatabase() {
	switch config.Database {
	case "redis":
		ServiceDao = redisServiceDao()
	default:
		panic("Database configuration value not recognized")
	}
}
