package repository

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	goRedis "github.com/go-redis/redis/v8"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/pkg/redis"
)

type CacheRepo interface {
	Create(cache model.Cache, serviceKey string) error
	Get(path string, serviceKey string) (model.Cache, error)
	DeleteByPath(paths []model.CachePath) error
	DeleteByTags(tags []string) error
}

var (
	maxTxRetries = 1000
)

type CacheRepoRedis struct {
	redis redis.Cmdable
}

func (r CacheRepoRedis) Create(cache model.Cache, serviceKey string) error {
	bytes, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	cacheKey := getCacheRedisKey(cache.Path, serviceKey)
	tagsKey := getCacheTagsRedisKey(cache.Path, serviceKey)
	keys := []string{cacheKey, tagsKey}

	txFn := func(tx *goRedis.Tx) error {
		pipe := tx.TxPipeline()

		pipe.SAdd(ctx, tagsKey, cache.Tags)
		pipe.Set(ctx, cacheKey, string(bytes), time.Duration(cache.TTL))

		_, err := pipe.Exec(ctx)
		return err
	}

	return r.redis.OptimisticLockTx(ctx, maxTxRetries, txFn, keys...)
}

// Get gets a cache
func (r CacheRepoRedis) Get(path string, serviceKey string) (model.Cache, error) {
	cacheKey := getCacheRedisKey(path, serviceKey)

	result, err := r.redis.Get(ctx, cacheKey).Result()
	if err != nil {
		return model.Cache{}, err
	}

	var cache model.Cache
	if err := json.Unmarshal([]byte(result), &cache); err != nil {
		return model.Cache{}, err
	}
	return cache, nil
}

// DeleteByPath deletes caches by specifying its path
func (r CacheRepoRedis) DeleteByPath(paths []model.CachePath) error {
	cacheKeys := make([]string, len(paths))
	for index, item := range paths {
		cacheKeys[index] = getCacheRedisKey(item.Path, item.ServicePath)
	}

	ok, notFoundIndex, err := r.redis.AllExist(ctx, cacheKeys...)
	if err != nil {
		return err
	}
	if !ok && notFoundIndex >= 0 {
		return &model.ErrCachePathNotFound{
			CachePath: paths[notFoundIndex],
		}
	}

	return r.deleteCaches(cacheKeys...)
}

// DeleteByTags deletes caches by specifying its tags
func (r CacheRepoRedis) DeleteByTags(tags []string) error {
	tmpTagsToDeleteKey := fmt.Sprintf("tmp::tags::delete::%d", time.Now().UnixNano())
	err := r.redis.SAdd(ctx, tmpTagsToDeleteKey, tags).Err()
	if err != nil {
		return err
	}
	defer func() error {
		return r.redis.Del(ctx, tmpTagsToDeleteKey).Err()
	}()

	var cursor uint64
	var wg sync.WaitGroup

	for {
		var keys []string
		var err error
		keys, cursor, err = r.redis.Scan(ctx, cursor, "cache::*::tags", 20).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			wg.Add(1)
			go func(keys, tags []string, wg *sync.WaitGroup) {
				defer wg.Done()

				pipe := r.redis.Pipeline()
				var cmds []*goRedis.StringSliceCmd
				for _, key := range keys {
					cmd := pipe.SInter(ctx, key, tmpTagsToDeleteKey)
					cmds = append(cmds, cmd)
				}

				_, err := pipe.Exec(ctx)
				if err != nil {
					return
				}

				for index, cmd := range cmds {
					intersection := cmd.Val()
					if len(intersection) > 0 {
						r.deleteCacheByCacheTagsKey(keys[index])
					}
				}
			}(keys, tags, &wg)
		}

		if cursor == 0 {
			break
		}
	}
	wg.Wait()

	return nil
}

func (r CacheRepoRedis) deleteCaches(cacheKeys ...string) error {
	var keys []string
	for _, cacheKey := range cacheKeys {
		cacheTagsKey := fmt.Sprintf("%s::tags", cacheKey)
		keys = append(keys, cacheKey, cacheTagsKey)
	}
	return r.redis.Del(ctx, keys...).Err()
}

func (r CacheRepoRedis) deleteCacheByCacheTagsKey(cacheTagsKey string) error {
	key := strings.TrimSuffix(cacheTagsKey, "::tags")
	return r.deleteCaches(key)
}

func getCacheRedisKey(path, serviceKey string) string {
	return fmt.Sprintf("cache::%s::%s", serviceKey, path)
}

func getCacheTagsRedisKey(path, serviceKey string) string {
	return fmt.Sprintf("%s::tags", getCacheRedisKey(path, serviceKey))
}

func NewCacheRepoRedis(redis redis.Cmdable) CacheRepo {
	return CacheRepoRedis{redis}
}
