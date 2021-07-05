package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/pkg/redis"

	goRedis "github.com/go-redis/redis/v8"
)

type ServiceRepo interface {
	Create(service model.Service) error
	GetAll() ([]model.Service, error)
	Get(key string) (model.Service, error)
	Delete(key string) error
	Upsert(service model.Service) error
}

type ServiceRepoRedis struct {
	redis redis.Cmdable
}

var ctx = context.Background()

// Create stores a service into redis
func (s ServiceRepoRedis) Create(service model.Service) error {
	if s.redis.Exists(ctx, getServiceRedisKey(service.ID)).Val() == 1 {
		return model.ErrServiceAlreadyRegistered
	}
	return s.Upsert(service)
}

// Get gets all services from redis
func (s ServiceRepoRedis) GetAll() ([]model.Service, error) {
	keys, err := s.redis.Keys(ctx, "service::*").Result()
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, model.ErrServiceNotFound
	}

	services := make([]model.Service, len(keys))
	for i, key := range keys {
		rediskey := strings.TrimPrefix(key, "service::")
		s, err := s.Get(rediskey)
		if err != nil {
			return nil, redisServiceError(err)
		}
		services[i] = s
	}

	return services, nil
}

// Get gets a service from redis
func (s ServiceRepoRedis) Get(key string) (model.Service, error) {
	val, err := s.redis.Get(ctx, getServiceRedisKey(key)).Result()
	if err != nil {
		return model.Service{}, redisServiceError(err)
	}
	var Service model.Service
	if err := json.Unmarshal([]byte(val), &Service); err != nil {
		return model.Service{}, err
	}
	return Service, nil
}

// Delete deletes a service from redis
func (s ServiceRepoRedis) Delete(key string) error {
	val, err := s.redis.Del(ctx, getServiceRedisKey(key)).Result()
	if err != nil {
		return err
	}
	if val == 0 {
		return model.ErrServiceNotFound
	}
	return nil
}

// Upsert creates or updates a service in redis
func (s ServiceRepoRedis) Upsert(service model.Service) error {
	bytes, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("error serializing service %v", err)
	}
	return s.redis.Set(ctx, getServiceRedisKey(service.ID), string(bytes), 0).Err()
}

func getServiceRedisKey(key string) string {
	return "service::" + strings.ToLower(key)
}

func redisServiceError(err error) error {
	if err == nil {
		return nil
	}
	if err == goRedis.Nil {
		return model.ErrServiceNotFound
	}
	return err
}

func NewServiceRepoRedis(redis redis.Cmdable) ServiceRepo {
	return ServiceRepoRedis{redis}
}
