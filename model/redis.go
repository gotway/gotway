package model

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/log"
)

var client *redis.Client

type serviceDaoRedis struct {
}

func redisServiceDao() ServiceDaoI {
	initializeClient()
	return serviceDaoRedis{}
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
func (dao serviceDaoRedis) StoreService(service Service) error {
	redisKey := getRedisKey(service.Path)
	healthPath, err := service.HealthPathForType()
	if err != nil {
		return err
	}
	serviceMap := map[string]interface{}{
		"type":       service.Type,
		"url":        service.URL,
		"path":       service.Path,
		"healthPath": *healthPath,
		"status":     Healthy,
	}
	saved, err := client.HSet(context.Background(), redisKey, serviceMap).Result()
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
func (dao serviceDaoRedis) GetService(key string) (*Service, error) {
	redisKey := getRedisKey(key)
	values := client.HGetAll(context.Background(), redisKey).Val()
	if len(values) == 0 {
		return nil, ErrServiceNotFound
	}
	service, err := newService(values)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// GetAllServices gets all service keys
func (dao serviceDaoRedis) GetAllServices() []string {
	keyPattern := "service:*"
	servicesKeys := client.Keys(context.Background(), keyPattern)
	var keys []string
	for _, redisKey := range servicesKeys.Val() {
		keys = append(keys, getKey(redisKey))
	}
	return keys
}

// UpdateServiceStatus updates the status of a service
func (dao serviceDaoRedis) UpdateServiceStatus(key string, status ServiceStatus) {
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

func newService(values map[string]string) (*Service, error) {
	serviceType := ServiceType(values["type"])
	if err := serviceType.validate(); err != nil {
		return nil, err
	}
	serviceStatus := ServiceStatus(values["status"])
	if err := serviceStatus.validate(); err != nil {
		return nil, err
	}
	return &Service{
		Type:       serviceType,
		URL:        values["url"],
		Path:       values["path"],
		HealthPath: values["healthPath"],
		Status:     serviceStatus,
	}, nil
}
