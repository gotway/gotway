package model

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gotway/gotway/core"
	"github.com/gotway/gotway/log"
)

type cacheConfigRepositoryRedis struct{}

func (c cacheConfigRepositoryRedis) StoreConfig(config core.CacheConfig, serviceKey string) error {
	cacheRedisKey := getCacheConfigRedisKey(serviceKey)
	err := client.Set(context.Background(), cacheRedisKey, config.TTL, 0).Err()
	if err != nil {
		log.Logger.Error(err)
		return err
	}

	statusesRedisKey := getStatusesRedisKey(serviceKey)
	err = client.SAdd(context.Background(), statusesRedisKey, config.Statuses.Serialize()).Err()
	if err != nil {
		log.Logger.Error(err)
		return err
	}

	tagsRedisKey := getTagsRedisKey(serviceKey)
	err = client.SAdd(context.Background(), tagsRedisKey, config.Tags).Err()
	if err != nil {
		log.Logger.Error(err)
		return err
	}

	return nil
}

func (c cacheConfigRepositoryRedis) GetConfig(serviceKey string) (core.CacheConfig, error) {
	cacheRedisKey := getCacheConfigRedisKey(serviceKey)
	defaultTTL := client.Get(context.Background(), cacheRedisKey).Val()
	if defaultTTL == "" {
		return core.CacheConfig{}, core.ErrCacheConfigNotFound
	}

	statusesRedisKey := getStatusesRedisKey(serviceKey)
	statuses := client.SMembers(context.Background(), statusesRedisKey).Val()
	if statuses == nil {
		return core.CacheConfig{}, core.ErrCacheConfigNotFound
	}

	tagsRedisKey := getTagsRedisKey(serviceKey)
	tags := client.SMembers(context.Background(), tagsRedisKey).Val()
	if tags == nil {
		return core.CacheConfig{}, core.ErrCacheConfigNotFound
	}

	cache, err := newCacheConfig(defaultTTL, statuses, tags)
	if err != nil {
		return core.CacheConfig{}, err
	}
	return cache, nil
}

func (c cacheConfigRepositoryRedis) DeleteConfig(serviceKey string) error {
	cacheRedisKey := getCacheConfigRedisKey(serviceKey)
	statusesRedisKey := getStatusesRedisKey(serviceKey)
	tagsRedisKey := getTagsRedisKey(serviceKey)
	return client.Del(context.Background(), cacheRedisKey, statusesRedisKey, tagsRedisKey).Err()
}

func (c cacheConfigRepositoryRedis) IsCacheableStatusCode(serviceKey string, statusCode int) bool {
	statusesRedisKey := getStatusesRedisKey(serviceKey)
	return client.SIsMember(context.Background(), statusesRedisKey, statusCode).Val()
}

func getCacheConfigRedisKey(key string) string {
	return fmt.Sprintf("service:%s:cache", key)
}

func getStatusesRedisKey(key string) string {
	return getCacheConfigRedisKey(key) + ":statuses"
}

func getTagsRedisKey(key string) string {
	return getCacheConfigRedisKey(key) + ":tags"
}

func newCacheConfig(ttlString string, statusesStr []string, tags []string) (core.CacheConfig, error) {
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
