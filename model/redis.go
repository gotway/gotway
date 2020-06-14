package model

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/log"
)

var client *redis.Client

// ServiceDaoRedis struct implementing ServiceDaoI
type ServiceDaoRedis struct {
}

// RedisServiceDao returns a redis DAO
func RedisServiceDao() ServiceDaoI {
	initializeClient()
	return ServiceDaoRedis{}
}

func initializeClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     config.RedisServer,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Error("Error connecting to redis server")
		panic(err)
	}

	log.Info("Connected to redis server ", config.RedisServer)
}

// StoreService stores a service into redis
func (dao ServiceDaoRedis) StoreService(key string, url string, healthURL string) bool {

	redisKey := "service:" + key + ":hooks"

	saved, err := client.HSet(context.Background(), redisKey, "url", url, "healthUrl", healthURL).Result()
	if err != nil {
		return false
	}

	return saved != 0
}
