package model

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/log"
)

var client *redis.Client

// ServiceDaoRedis struct implementing ServiceDaoI
type ServiceDaoRedis struct {
}

func redisServiceDao() ServiceDaoI {
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
func (dao ServiceDaoRedis) StoreService(key string, url string, healthURL string) error {
	redisKey := getRedisKey(key)
	saved, err := client.HSet(context.Background(), redisKey, "url", url, "healthUrl", healthURL, "status", Healthy).Result()
	if err != nil {
		log.Error(err)
		return err
	}
	if saved == 0 {
		return ErrServiceAlreadyRegistered
	}
	return nil
}

// GetService gets a service from redis
func (dao ServiceDaoRedis) GetService(key string) (*Service, error) {
	redisKey := getRedisKey(key)
	values := client.HGetAll(context.Background(), redisKey).Val()
	if len(values) == 0 {
		return nil, ErrServiceNotFound
	}
	service, err := newService(key, values)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (dao ServiceDaoRedis) getAllServices() []string {
	keyPattern := "service:*"
	servicesKeys := client.Keys(context.Background(), keyPattern)
	var keys []string
	for _, redisKey := range servicesKeys.Val() {
		keys = append(keys, getKey(redisKey))
	}
	return keys
}

func (dao ServiceDaoRedis) updateServiceStatus(key string, status ServiceStatus) {
	redisKey := getRedisKey(key)
	if err := status.validate(); err != nil {
		log.Error(err)
		return
	}
	_, err := client.HSet(context.Background(), redisKey, "status", status).Result()
	if err != nil {
		log.Error(err)
	}
}

func getRedisKey(key string) string {
	return "service:" + key
}

func getKey(redisKey string) string {
	return strings.Split(redisKey, ":")[1]
}

func newService(key string, values map[string]string) (*Service, error) {
	serviceStatus := ServiceStatus(values["status"])
	if err := serviceStatus.validate(); err != nil {
		return nil, err
	}
	return &Service{Key: key, URL: values["url"], HealthURL: values["healthUrl"], Status: serviceStatus}, nil
}
