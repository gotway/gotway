package controller

import (
	"github.com/gosmo-devs/microgateway/log"
	"github.com/gosmo-devs/microgateway/model"
)

// RegisterService adds a new entry in Redis for an endpoint
func RegisterService(key string, url string, healthEndpoint string) error {

	if !model.ServiceDao.StoreService(key, url, healthEndpoint) {
		log.Warnf("Service %s was already registered", key)
		return ErrAlreadyRegistered
	}

	return nil
}
