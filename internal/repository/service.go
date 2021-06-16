package repository

import (
	"context"
	"fmt"
	"strconv"

	goRedis "github.com/go-redis/redis/v8"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/pkg/redis"
)

type ServiceRepo interface {
	StoreService(service model.ServiceDetail) error
	GetAllServiceKeys() []string
	GetService(key string) (model.Service, error)
	GetServiceDetail(key string) (model.ServiceDetail, error)
	GetServices(keys ...string) ([]model.Service, error)
	DeleteService(key string) error
	UpdateServicesStatus(status model.ServiceStatus, keys ...string) error
	GetServiceCache(key string) (model.CacheConfig, error)
	IsCacheableStatusCode(key string, statusCode int) bool
}

type ServiceRepoRedis struct {
	redis redis.Cmdable
}

const (
	serviceSet   = "service"
	maxTxRetries = 1000
)

var ctx = context.Background()

// StoreService stores a service into redis
func (s ServiceRepoRedis) StoreService(service model.ServiceDetail) error {
	serviceMap := map[string]interface{}{
		"type":       service.Type,
		"url":        service.URL,
		"path":       service.Path,
		"healthPath": service.HealthPath,
		"status":     model.ServiceStatusHealthy,
	}

	serviceKey := getServiceRedisKey(service.Path)
	cacheKey := getServiceCacheRedisKey(service.Path)
	statusesKey := getStatusesRedisKey(service.Path)
	tagsKey := getTagsRedisKey(service.Path)
	keys := []string{serviceKey, cacheKey, statusesKey, tagsKey}

	txFn := func(tx *goRedis.Tx) error {
		pipe := tx.TxPipeline()

		storeService := pipe.HSet(ctx, serviceKey, serviceMap)
		pipe.SAdd(ctx, serviceSet, service.Path)
		if !service.Cache.IsEmpty() {
			pipe.Set(ctx, cacheKey, service.Cache.TTL, 0)
			pipe.SAdd(ctx, statusesKey, service.Cache.Statuses.Serialize())
			pipe.SAdd(ctx, tagsKey, service.Cache.Tags)
		}

		_, err := pipe.Exec(ctx)
		if err != nil {
			return err
		}

		if storeService.Val() == 0 {
			return model.ErrServiceAlreadyRegistered
		}

		return nil
	}

	return s.redis.OptimisticLockTx(ctx, maxTxRetries, txFn, keys...)
}

// GetAllServiceKeys gets all service keys
func (s ServiceRepoRedis) GetAllServiceKeys() []string {
	return s.redis.SMembers(ctx, serviceSet).Val()
}

// GetService gets a service from redis
func (s ServiceRepoRedis) GetService(key string) (model.Service, error) {
	redisKey := getServiceRedisKey(key)
	values := s.redis.HGetAll(ctx, redisKey).Val()
	return newService(values)
}

// GetServiceDetail gets a service with extra info
func (s ServiceRepoRedis) GetServiceDetail(key string) (model.ServiceDetail, error) {
	serviceKey := getServiceRedisKey(key)
	cacheKey := getServiceCacheRedisKey(key)
	statusesKey := getStatusesRedisKey(key)
	tagsKey := getTagsRedisKey(key)
	keys := []string{serviceKey, cacheKey, statusesKey, tagsKey}

	serviceDetail := model.ServiceDetail{}
	txFn := func(tx *goRedis.Tx) error {
		pipe := tx.TxPipeline()

		getService := pipe.HGetAll(ctx, serviceKey)
		getTTL := pipe.Get(ctx, cacheKey)
		getStatuses := pipe.SMembers(ctx, statusesKey)
		getTags := pipe.SMembers(ctx, tagsKey)

		pipe.Exec(ctx)

		serviceValues, err := getService.Result()
		service, err := newService(serviceValues)
		if err != nil {
			return err
		}

		ttl, errTTL := getTTL.Result()
		statuses, errStatuses := getStatuses.Result()
		tags, errTags := getTags.Result()
		errs := []error{errTTL, errStatuses, errTags}

		var cacheConfig model.CacheConfig
		if redis.AnyEmptyErr(errs...) {
			cacheConfig = model.DefaultCacheConfig
		} else {
			cacheConfig, err = newCacheConfig(ttl, statuses, tags)
			if err != nil {
				return err
			}
		}

		serviceDetail = model.ServiceDetail{
			Service: service,
			Cache:   cacheConfig,
		}

		return nil
	}

	return serviceDetail, s.redis.OptimisticLockTx(ctx, maxTxRetries, txFn, keys...)
}

// GetServices gets services from redis
func (s ServiceRepoRedis) GetServices(keys ...string) ([]model.Service, error) {
	var services []model.Service
	txFn := func(tx *goRedis.Tx) error {
		pipe := tx.TxPipeline()

		cmds := make([]*goRedis.StringStringMapCmd, len(keys))
		for i, s := range keys {
			redisKey := getServiceRedisKey(s)
			cmds[i] = pipe.HGetAll(ctx, redisKey)
		}

		_, err := pipe.Exec(ctx)
		if err != nil {
			return err
		}

		for _, cmd := range cmds {
			values := cmd.Val()
			service, err := newService(values)
			if err != nil {
				return err
			}
			services = append(services, service)
		}

		return nil
	}

	return services, s.redis.OptimisticLockTx(ctx, maxTxRetries, txFn, keys...)
}

// DeleteService deletes a service from redis
func (s ServiceRepoRedis) DeleteService(key string) error {
	serviceKey := getServiceRedisKey(key)
	cacheRedisKey := getServiceCacheRedisKey(key)
	statusesRedisKey := getStatusesRedisKey(key)
	tagsRedisKey := getTagsRedisKey(key)
	keys := []string{
		serviceKey,
		cacheRedisKey,
		statusesRedisKey,
		tagsRedisKey,
		serviceSet,
	}

	txFn := func(tx *goRedis.Tx) error {
		pipe := tx.TxPipeline()

		delService := pipe.Del(ctx, serviceKey, cacheRedisKey, statusesRedisKey, tagsRedisKey)
		pipe.SRem(ctx, serviceSet, key)

		_, err := pipe.Exec(ctx)

		deleted := delService.Val()
		if deleted == 0 {
			return model.ErrServiceNotFound
		}

		return err
	}

	return s.redis.OptimisticLockTx(ctx, maxTxRetries, txFn, keys...)
}

// UpdateServiceStatus updates the status of a service
func (s ServiceRepoRedis) UpdateServicesStatus(status model.ServiceStatus, keys ...string) error {
	pipe := s.redis.TxPipeline()
	for _, key := range keys {
		serviceKey := getServiceRedisKey(key)
		if err := status.Validate(); err != nil {
			return err
		}
		pipe.HSet(ctx, serviceKey, "status", status)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetServiceCache gets cache configuration of a service
func (s ServiceRepoRedis) GetServiceCache(key string) (model.CacheConfig, error) {
	cacheKey := getServiceCacheRedisKey(key)
	statusesKey := getStatusesRedisKey(key)
	tagsKey := getTagsRedisKey(key)
	keys := []string{cacheKey, statusesKey, tagsKey}

	cache := model.CacheConfig{}
	txFn := func(tx *goRedis.Tx) error {
		pipe := tx.TxPipeline()

		getTTL := pipe.Get(ctx, cacheKey)
		getStatuses := pipe.SMembers(ctx, statusesKey)
		getTags := pipe.SMembers(ctx, tagsKey)

		pipe.Exec(ctx)

		ttl, errTTL := getTTL.Result()
		statuses, errStatuses := getStatuses.Result()
		tags, errTags := getTags.Result()
		errs := []error{errTTL, errStatuses, errTags}

		if redis.AnyEmptyErr(errs...) {
			return model.ErrCacheConfigNotFound
		}

		var err error
		cache, err = newCacheConfig(ttl, statuses, tags)

		return err
	}

	return cache, s.redis.OptimisticLockTx(ctx, maxTxRetries, txFn, keys...)
}

// IsCacheableStatusCode determines if a status code is cacheable for a service
func (s ServiceRepoRedis) IsCacheableStatusCode(key string, statusCode int) bool {
	statusesRedisKey := getStatusesRedisKey(key)
	return s.redis.SIsMember(ctx, statusesRedisKey, statusCode).Val()
}

func getServiceRedisKey(key string) string {
	return "service:" + key
}

func getServiceCacheRedisKey(key string) string {
	return fmt.Sprintf("service:%s:cache", key)
}

func getStatusesRedisKey(key string) string {
	return getServiceCacheRedisKey(key) + ":statuses"
}

func getTagsRedisKey(key string) string {
	return getServiceCacheRedisKey(key) + ":tags"
}

func newService(values map[string]string) (model.Service, error) {
	if len(values) == 0 {
		return model.Service{}, model.ErrServiceNotFound
	}
	service, err := processServiceMap(values)
	if err != nil {
		return model.Service{}, err
	}
	return service, nil
}

func processServiceMap(values map[string]string) (model.Service, error) {
	serviceType := model.ServiceType(values["type"])
	if err := serviceType.Validate(); err != nil {
		return model.Service{}, err
	}
	serviceStatus := model.ServiceStatus(values["status"])
	if err := serviceStatus.Validate(); err != nil {
		return model.Service{}, err
	}
	return model.Service{
		Type:       serviceType,
		URL:        values["url"],
		Path:       values["path"],
		HealthPath: values["healthPath"],
		Status:     serviceStatus,
	}, nil
}

func newCacheConfig(
	ttlString string,
	statusesStr []string,
	tags []string,
) (model.CacheConfig, error) {
	ttl, err := strconv.ParseInt(ttlString, 10, 64)
	if err != nil {
		return model.CacheConfig{}, err
	}

	statuses := make([]int, len(statusesStr))
	for i, s := range statusesStr {
		status, err := strconv.Atoi(s)
		if err != nil {
			return model.CacheConfig{}, err
		}
		statuses[i] = status
	}

	return model.CacheConfig{
		TTL:      int64(ttl),
		Statuses: statuses,
		Tags:     tags,
	}, nil
}

func NewServiceRepoRedis(redis redis.Cmdable) ServiceRepo {
	return ServiceRepoRedis{redis}
}
