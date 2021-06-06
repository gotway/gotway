package repository

import (
	"context"
	"fmt"
	"strconv"

	goRedis "github.com/go-redis/redis/v8"
	"github.com/gotway/gotway/internal/core"
	"github.com/gotway/gotway/pkg/redis"
)

type ServiceRepo interface {
	StoreService(service core.ServiceDetail) error
	GetAllServiceKeys() []string
	GetService(key string) (core.Service, error)
	GetServiceDetail(key string) (core.ServiceDetail, error)
	GetServices(keys ...string) ([]core.Service, error)
	DeleteService(key string) error
	UpdateServiceStatus(key string, status core.ServiceStatus) error
	GetServiceCache(key string) (core.CacheConfig, error)
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
func (s ServiceRepoRedis) StoreService(service core.ServiceDetail) error {
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
			return core.ErrServiceAlreadyRegistered
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
func (s ServiceRepoRedis) GetService(key string) (core.Service, error) {
	redisKey := getServiceRedisKey(key)
	values := s.redis.HGetAll(ctx, redisKey).Val()
	return newService(values)
}

// GetServiceDetail gets a service with extra info
func (s ServiceRepoRedis) GetServiceDetail(key string) (core.ServiceDetail, error) {
	serviceKey := getServiceRedisKey(key)
	cacheKey := getServiceCacheRedisKey(key)
	statusesKey := getStatusesRedisKey(key)
	tagsKey := getTagsRedisKey(key)
	keys := []string{serviceKey, cacheKey, statusesKey, tagsKey}

	serviceDetail := core.ServiceDetail{}
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

		var cacheConfig core.CacheConfig
		if redis.AnyEmptyErr(errs...) {
			cacheConfig = core.DefaultCacheConfig
		} else {
			cacheConfig, err = newCacheConfig(ttl, statuses, tags)
			if err != nil {
				return err
			}
		}

		serviceDetail = core.ServiceDetail{
			Service: service,
			Cache:   cacheConfig,
		}

		return nil
	}

	return serviceDetail, s.redis.OptimisticLockTx(ctx, maxTxRetries, txFn, keys...)
}

// GetServices gets services from redis
func (s ServiceRepoRedis) GetServices(keys ...string) ([]core.Service, error) {
	var services []core.Service
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
			return core.ErrServiceNotFound
		}

		return err
	}

	return s.redis.OptimisticLockTx(ctx, maxTxRetries, txFn, keys...)
}

// UpdateServiceStatus updates the status of a service
func (s ServiceRepoRedis) UpdateServiceStatus(key string, status core.ServiceStatus) error {
	serviceKey := getServiceRedisKey(key)
	if err := status.Validate(); err != nil {
		return err
	}
	return s.redis.HSet(ctx, serviceKey, "status", status).Err()
}

// GetServiceCache gets cache configuration of a service
func (s ServiceRepoRedis) GetServiceCache(key string) (core.CacheConfig, error) {
	cacheKey := getServiceCacheRedisKey(key)
	statusesKey := getStatusesRedisKey(key)
	tagsKey := getTagsRedisKey(key)
	keys := []string{cacheKey, statusesKey, tagsKey}

	cache := core.CacheConfig{}
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
			return core.ErrCacheConfigNotFound
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

func newService(values map[string]string) (core.Service, error) {
	if len(values) == 0 {
		return core.Service{}, core.ErrServiceNotFound
	}
	service, err := processServiceMap(values)
	if err != nil {
		return core.Service{}, err
	}
	return service, nil
}

func processServiceMap(values map[string]string) (core.Service, error) {
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

func newCacheConfig(
	ttlString string,
	statusesStr []string,
	tags []string,
) (core.CacheConfig, error) {
	ttl, err := strconv.ParseInt(ttlString, 10, 64)
	if err != nil {
		return core.CacheConfig{}, err
	}

	statuses := make([]int, len(statusesStr))
	for i, s := range statusesStr {
		status, err := strconv.Atoi(s)
		if err != nil {
			return core.CacheConfig{}, err
		}
		statuses[i] = status
	}

	return core.CacheConfig{
		TTL:      int64(ttl),
		Statuses: statuses,
		Tags:     tags,
	}, nil
}

func NewServiceRepoRedis(redis redis.Cmdable) ServiceRepo {
	return ServiceRepoRedis{redis}
}
