package controller

import (
	"strconv"

	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/log"
	"github.com/gosmo-devs/microgateway/redis"
)

// RegisterEndpoint adds a new entry in Redis for an endpoint
func RegisterEndpoint(key string, url string, healthEndpoint string, ttl int) bool {

	log.Debug("hola")

	if maxTTL, _ := strconv.Atoi(config.MAX_SERVICE_TTL); ttl > maxTTL {
		log.Warnf("Service tried to register with a TTL of %d. Maximum allowed is %s", ttl, config.MAX_SERVICE_TTL)
		return false
	}

	if !redis.Store(key, url, healthEndpoint, ttl) {
		log.Warnf("Service %s was already registered", key)
		return false
	}

	return true
}
