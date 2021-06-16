package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	goRedis "github.com/go-redis/redis/v8"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/pkg/redis"
)

type CacheRepo interface {
	StoreCache(cache model.CacheDetail, serviceKey string) error
	GetCache(path string, serviceKey string) (model.Cache, error)
	GetCacheDetail(path string, serviceKey string) (model.CacheDetail, error)
	DeleteCacheByPath(paths []model.CachePath) error
	DeleteCacheByTags(tags []string) error
}

type CacheRepoRedis struct {
	redis redis.Cmdable
}

// StoreCache stores a cache
func (r CacheRepoRedis) StoreCache(cache model.CacheDetail, serviceKey string) error {
	cacheMap, err := newCacheMap(cache.StatusCode, cache.Body)
	if err != nil {
		return nil
	}
	headersMap := newCacheHeadersMap(cache.Headers)
	ttl := time.Duration(cache.TTL)

	cacheKey := getCacheRedisKey(cache.Path, serviceKey)
	headersKey := getCacheHeadersRedisKey(cache.Path, serviceKey)
	tagsKey := getCacheTagsRedisKey(cache.Path, serviceKey)
	keys := []string{cacheKey, headersKey, tagsKey}

	pipe := r.redis.TxPipeline()

	pipe.HSet(ctx, cacheKey, cacheMap)
	pipe.HSet(ctx, headersKey, headersMap)
	pipe.SAdd(ctx, tagsKey, cache.Tags)
	for _, key := range keys {
		pipe.Expire(ctx, key, ttl)
	}

	_, err = pipe.Exec(ctx)

	return err
}

// GetCache gets a cache
func (r CacheRepoRedis) GetCache(path string, serviceKey string) (model.Cache, error) {
	cacheKey := getCacheRedisKey(path, serviceKey)
	headersKey := getCacheHeadersRedisKey(path, serviceKey)

	pipe := r.redis.TxPipeline()

	getCache := pipe.HGetAll(ctx, cacheKey)
	getHeaders := pipe.HGetAll(ctx, headersKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return model.Cache{}, err
	}

	cacheMap, headersMap := getCache.Val(), getHeaders.Val()
	if len(cacheMap) == 0 || len(headersMap) == 0 {
		return model.Cache{}, model.ErrCacheNotFound
	}

	return newCache(path, cacheMap, headersMap)
}

// GetCacheDetail gets a extended version of a cache
func (r CacheRepoRedis) GetCacheDetail(
	path string,
	serviceKey string,
) (model.CacheDetail, error) {
	cache, err := r.GetCache(path, serviceKey)
	if err != nil {
		return model.CacheDetail{}, err
	}

	cacheKey := getCacheRedisKey(path, serviceKey)
	ttl, err := r.redis.TTL(context.Background(), cacheKey).Result()
	if err != nil {
		return model.CacheDetail{}, err
	}

	cacheTagsKey := getCacheTagsRedisKey(path, serviceKey)
	tags, err := r.redis.SMembers(context.Background(), cacheTagsKey).Result()
	if err != nil {
		return model.CacheDetail{}, err
	}

	cacheDetail := model.CacheDetail{
		Cache: cache,
		TTL:   model.CacheTTL(ttl),
		Tags:  tags,
	}

	return cacheDetail, nil
}

// DeleteCacheByPath deletes caches by specifying its path
func (r CacheRepoRedis) DeleteCacheByPath(paths []model.CachePath) error {
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

// DeleteCacheByTags deletes caches by specifying its tags
func (r CacheRepoRedis) DeleteCacheByTags(tags []string) error {
	tmpTagsToDeleteKey := fmt.Sprintf("tmp:tags:delete:%d", time.Now().UnixNano())
	err := r.redis.SAdd(context.Background(), tmpTagsToDeleteKey, tags).Err()
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
		keys, cursor, err = r.redis.Scan(ctx, cursor, "cache:*:tags", 20).Result()
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

func (r CacheRepoRedis) deleteCaches(cacheRedisKeys ...string) error {
	var keys []string
	for _, cacheKey := range cacheRedisKeys {
		cacheHeadersKey := fmt.Sprintf("%s:headers", cacheKey)
		cacheTagsKey := fmt.Sprintf("%s:tags", cacheKey)

		newKeys := []string{cacheKey, cacheHeadersKey, cacheTagsKey}
		keys = append(keys, newKeys...)
	}

	return r.redis.Del(context.Background(), keys...).Err()
}

func (r CacheRepoRedis) deleteCacheByCacheTagsKey(cacheTagsKey string) error {
	parts := strings.Split(cacheTagsKey, ":tags")
	return r.deleteCaches(parts[0])
}

func getCacheRedisKey(path string, serviceKey string) string {
	return fmt.Sprintf("cache:%s:%s", serviceKey, path)
}

func getCacheHeadersRedisKey(path string, serviceKey string) string {
	return fmt.Sprintf("%s:headers", getCacheRedisKey(path, serviceKey))
}

func getCacheTagsRedisKey(path string, serviceKey string) string {
	return fmt.Sprintf("%s:tags", getCacheRedisKey(path, serviceKey))
}

func newCacheMap(statusCode int, body io.Reader) (map[string]interface{}, error) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"statusCode": strconv.Itoa(statusCode),
		"body":       string(bodyBytes),
	}, nil
}

var cacheHeadersSeparator = ","

func newCacheHeadersMap(headers http.Header) map[string]interface{} {
	cacheHeaders := make(map[string]interface{})
	for key, values := range headers {
		cacheHeaders[key] = strings.Join(values[:], cacheHeadersSeparator)
	}
	return cacheHeaders
}

func newCache(
	path string,
	cacheMap map[string]string,
	headersMap map[string]string,
) (model.Cache, error) {
	statusCode, err := strconv.Atoi(cacheMap["statusCode"])
	if err != nil {
		return model.Cache{}, err
	}

	headers := make(map[string][]string)
	for key, header := range headersMap {
		headers[key] = strings.Split(header, cacheHeadersSeparator)
	}

	body := model.CacheBody{
		Reader: ioutil.NopCloser(bytes.NewBufferString(cacheMap["body"])),
	}

	return model.Cache{
		Path:       path,
		StatusCode: statusCode,
		Headers:    headers,
		Body:       body,
	}, nil
}

func NewCacheRepoRedis(redis redis.Cmdable) CacheRepoRedis {
	return CacheRepoRedis{redis}
}
