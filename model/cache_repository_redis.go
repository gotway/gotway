package model

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

	"github.com/go-redis/redis/v8"
	"github.com/gotway/gotway/core"
)

// CacheRepositoryRedis is a redis implementation for cache repo
type CacheRepositoryRedis struct {
	client *redis.Client
}

func newCacheRepositoryRedis(client *redis.Client) CacheRepositoryRedis {
	return CacheRepositoryRedis{client}
}

// StoreCache stores a cache
func (r CacheRepositoryRedis) StoreCache(cache core.CacheDetail, serviceKey string) error {
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

	pipe := r.client.TxPipeline()

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
func (r CacheRepositoryRedis) GetCache(path string, serviceKey string) (core.Cache, error) {
	cacheKey := getCacheRedisKey(path, serviceKey)
	headersKey := getCacheHeadersRedisKey(path, serviceKey)

	pipe := r.client.TxPipeline()

	getCache := pipe.HGetAll(ctx, cacheKey)
	getHeaders := pipe.HGetAll(ctx, headersKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return core.Cache{}, err
	}

	cacheMap, headersMap := getCache.Val(), getHeaders.Val()
	if len(cacheMap) == 0 || len(headersMap) == 0 {
		return core.Cache{}, core.ErrCacheNotFound
	}

	return newCache(path, cacheMap, headersMap)
}

// GetCacheDetail gets a extended version of a cache
func (r CacheRepositoryRedis) GetCacheDetail(path string, serviceKey string) (core.CacheDetail, error) {
	cache, err := r.GetCache(path, serviceKey)
	if err != nil {
		return core.CacheDetail{}, err
	}

	cacheKey := getCacheRedisKey(path, serviceKey)
	ttl, err := r.client.TTL(context.Background(), cacheKey).Result()
	if err != nil {
		return core.CacheDetail{}, err
	}

	cacheTagsKey := getCacheTagsRedisKey(path, serviceKey)
	tags, err := r.client.SMembers(context.Background(), cacheTagsKey).Result()
	if err != nil {
		return core.CacheDetail{}, err
	}

	cacheDetail := core.CacheDetail{
		Cache: cache,
		TTL:   core.CacheTTL(ttl),
		Tags:  tags,
	}

	return cacheDetail, nil
}

// DeleteCacheByPath deletes caches by specifying its path
func (r CacheRepositoryRedis) DeleteCacheByPath(paths []core.CachePath) error {
	cacheKeys := make([]string, len(paths))
	for index, item := range paths {
		cacheKeys[index] = getCacheRedisKey(item.Path, item.ServicePath)
	}

	ok, notFoundIndex, err := allExists(r.client, cacheKeys...)
	if err != nil {
		return err
	}
	if !ok && notFoundIndex >= 0 {
		return &core.ErrCachePathNotFound{
			CachePath: paths[notFoundIndex],
		}
	}

	return deleteCaches(r.client, cacheKeys...)
}

// DeleteCacheByTags deletes caches by specifying its tags
func (r CacheRepositoryRedis) DeleteCacheByTags(tags []string) error {
	tmpTagsToDeleteKey := fmt.Sprintf("tmp:tags:delete:%d", time.Now().UnixNano())
	err := r.client.SAdd(context.Background(), tmpTagsToDeleteKey, tags).Err()
	if err != nil {
		return err
	}
	defer func() error {
		return r.client.Del(context.Background(), tmpTagsToDeleteKey).Err()
	}()

	var cursor uint64
	var wg sync.WaitGroup

	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(context.Background(), cursor, "cache:*:tags", 20).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			wg.Add(1)
			go func(keys, tags []string, wg *sync.WaitGroup) {
				defer wg.Done()

				pipe := r.client.Pipeline()
				var cmds []*redis.StringSliceCmd
				for _, key := range keys {
					cmd := pipe.SInter(context.Background(), key, tmpTagsToDeleteKey)
					cmds = append(cmds, cmd)
				}

				_, err := pipe.Exec(context.Background())
				if err != nil {
					return
				}

				for index, cmd := range cmds {
					intersection := cmd.Val()
					if len(intersection) > 0 {
						deleteCacheByCacheTagsKey(r.client, keys[index])
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

func newCache(path string, cacheMap map[string]string, headersMap map[string]string) (core.Cache, error) {
	statusCode, err := strconv.Atoi(cacheMap["statusCode"])
	if err != nil {
		return core.Cache{}, err
	}

	headers := make(map[string][]string)
	for key, header := range headersMap {
		headers[key] = strings.Split(header, cacheHeadersSeparator)
	}

	body := core.CacheBody{
		Reader: ioutil.NopCloser(bytes.NewBufferString(cacheMap["body"])),
	}

	return core.Cache{
		Path:       path,
		StatusCode: statusCode,
		Headers:    headers,
		Body:       body,
	}, nil
}

func deleteCaches(client redis.Cmdable, cacheRedisKeys ...string) error {
	var keys []string
	for _, cacheKey := range cacheRedisKeys {
		cacheHeadersKey := fmt.Sprintf("%s:headers", cacheKey)
		cacheTagsKey := fmt.Sprintf("%s:tags", cacheKey)

		newKeys := []string{cacheKey, cacheHeadersKey, cacheTagsKey}
		keys = append(keys, newKeys...)
	}

	return client.Del(context.Background(), keys...).Err()
}

func deleteCacheByCacheTagsKey(client redis.Cmdable, cacheTagsKey string) error {
	parts := strings.Split(cacheTagsKey, ":tags")
	return deleteCaches(client, parts[0])
}
