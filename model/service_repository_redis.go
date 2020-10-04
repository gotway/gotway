package model

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gotway/gotway/core"
	"github.com/gotway/gotway/log"
)

var serviceSet = "service"

type serviceRepositoryRedis struct{}

// StoreService stores a service into redis
func (s serviceRepositoryRedis) StoreService(service core.Service) error {
	redisKey := getRedisKey(service.Path)
	healthPath, err := service.HealthPathForType()
	if err != nil {
		return err
	}
	serviceMap := map[string]interface{}{
		"type":       service.Type,
		"url":        service.URL,
		"path":       service.Path,
		"healthPath": healthPath,
		"status":     core.ServiceStatusHealthy,
	}
	saved, err := client.HSet(context.Background(), redisKey, serviceMap).Result()
	if err != nil {
		log.Logger.Error(err)
		return err
	}
	if saved == 0 {
		return core.ErrServiceAlreadyRegistered
	}
	_, err = client.SAdd(context.Background(), serviceSet, service.Path).Result()
	if err != nil {
		log.Logger.Error(err)
	}
	return err
}

// ExistService determines if a service is registered
func (s serviceRepositoryRedis) ExistService(key string) error {
	isMember := client.SIsMember(context.Background(), serviceSet, key).Val()
	if !isMember {
		return core.ErrServiceNotFound
	}
	return nil
}

// GetAllServiceKeys gets all service keys
func (s serviceRepositoryRedis) GetAllServiceKeys() []string {
	return client.SMembers(context.Background(), serviceSet).Val()
}

// GetService gets a service from redis
func (s serviceRepositoryRedis) GetService(key string) (core.Service, error) {
	redisKey := getRedisKey(key)
	values := client.HGetAll(context.Background(), redisKey).Val()
	return processServiceMap(values)
}

// GetServices gets services from redis
func (s serviceRepositoryRedis) GetServices(keys ...string) ([]core.Service, error) {
	pipe := client.Pipeline()
	var cmds []*redis.StringStringMapCmd
	for _, s := range keys {
		redisKey := getRedisKey(s)
		cmd := pipe.HGetAll(context.Background(), redisKey)
		cmds = append(cmds, cmd)
	}
	_, err := pipe.Exec(context.Background())
	if err != nil {
		return nil, err
	}
	var services []core.Service
	for _, cmd := range cmds {
		values := cmd.Val()
		service, err := processServiceMap(values)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}

// DeleteService deletes a service from redis
func (s serviceRepositoryRedis) DeleteService(key string) error {
	redisKey := getRedisKey(key)
	err := client.Del(context.Background(), redisKey).Err()
	if err != nil {
		return err
	}
	_, err = client.SRem(context.Background(), serviceSet, key).Result()
	return err
}

// UpdateServiceStatus updates the status of a service
func (s serviceRepositoryRedis) UpdateServiceStatus(key string, status core.ServiceStatus) {
	redisKey := getRedisKey(key)
	if err := status.Validate(); err != nil {
		log.Logger.Error(err)
		return
	}
	_, err := client.HSet(context.Background(), redisKey, "status", status).Result()
	if err != nil {
		log.Logger.Error(err)
	}
}

func getRedisKey(key string) string {
	return "service:" + key
}

func getKey(redisKey string) string {
	return strings.Split(redisKey, ":")[1]
}

func processServiceMap(values map[string]string) (core.Service, error) {
	if len(values) == 0 {
		return core.Service{}, core.ErrServiceNotFound
	}
	service, err := newService(values)
	if err != nil {
		return core.Service{}, err
	}
	return service, nil
}

func newService(values map[string]string) (core.Service, error) {
	serviceType := core.ServiceType(values["type"])
	if err := serviceType.Validate(); err != nil {
		return core.Service{}, err
	}
	serviceStatus := core.ServiceStatus(values["status"])
	if err := serviceStatus.Validate(); err != nil {
		return core.Service{}, err
	}
	return core.Service{
		Type:       serviceType,
		URL:        values["url"],
		Path:       values["path"],
		HealthPath: values["healthPath"],
		Status:     serviceStatus,
	}, nil
}
