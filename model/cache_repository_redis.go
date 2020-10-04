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

type cacheRepositoryRedis struct{}

func (r cacheRepositoryRedis) StoreCache(cache core.CacheDetail, serviceKey string) error {
	cacheMap, err := newCacheMap(cache.StatusCode, cache.Body)
	if err != nil {
		return nil
	}
	cacheKey := getCacheRedisKey(cache.Path, serviceKey)
	ttl := time.Duration(cache.TTL)
	err = hsetTTL(cacheKey, cacheMap, ttl)
	if err != nil {
		return err
	}

	cacheHeadersKey := getCacheHeadersRedisKey(cache.Path, serviceKey)
	cacheHeadersMap := newCacheHeadersMap(cache.Headers)
	err = hsetTTL(cacheHeadersKey, cacheHeadersMap, ttl)
	if err != nil {
		return err
	}

	cacheTagsKey := getCacheTagsRedisKey(cache.Path, serviceKey)
	err = saddTTL(cacheTagsKey, ttl, cache.Tags)
	if err != nil {
		return err
	}

	return nil
}

func (r cacheRepositoryRedis) GetCache(path string, serviceKey string) (core.Cache, error) {
	cacheKey := getCacheRedisKey(path, serviceKey)
	cacheHeadersKey := getCacheHeadersRedisKey(path, serviceKey)
	keys := []string{cacheKey, cacheHeadersKey}

	pipe := client.Pipeline()
	var cmds []*redis.StringStringMapCmd
	for _, key := range keys {
		cmd := pipe.HGetAll(context.Background(), key)
		cmds = append(cmds, cmd)
	}

	_, err := pipe.Exec(context.Background())
	if err != nil {
		return core.Cache{}, nil
	}

	cacheMap, headersMap := cmds[0].Val(), cmds[1].Val()
	if len(cacheMap) == 0 || len(headersMap) == 0 {
		return core.Cache{}, core.ErrCacheNotFound
	}

	cache, err := newCache(path, cacheMap, headersMap)
	if err != nil {
		return core.Cache{}, err
	}

	return cache, nil
}

func (r cacheRepositoryRedis) GetCacheDetail(path string, serviceKey string) (core.CacheDetail, error) {
	cache, err := r.GetCache(path, serviceKey)
	if err != nil {
		return core.CacheDetail{}, err
	}

	cacheKey := getCacheRedisKey(path, serviceKey)
	ttl, err := client.TTL(context.Background(), cacheKey).Result()
	if err != nil {
		return core.CacheDetail{}, err
	}

	cacheTagsKey := getCacheTagsRedisKey(path, serviceKey)
	tags, err := client.SMembers(context.Background(), cacheTagsKey).Result()
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

func (r cacheRepositoryRedis) DeleteCacheByPath(paths []core.CachePath) error {
	cacheKeys := make([]string, len(paths))
	for index, item := range paths {
		cacheKeys[index] = getCacheRedisKey(item.Path, item.ServicePath)
	}

	ok, notFoundIndex, err := indexedExists(cacheKeys...)
	if err != nil {
		return err
	}
	if !ok && notFoundIndex >= 0 {
		return &core.ErrCachePathNotFound{
			CachePath: paths[notFoundIndex],
		}
	}

	return deleteCaches(cacheKeys...)
}

func (r cacheRepositoryRedis) DeleteCacheByTags(tags []string) error {
	tmpTagsToDeleteKey := fmt.Sprintf("tmp:tags:delete:%d", time.Now().UnixNano())
	err := client.SAdd(context.Background(), tmpTagsToDeleteKey, tags).Err()
	if err != nil {
		return err
	}
	defer func() error {
		return client.Del(context.Background(), tmpTagsToDeleteKey).Err()
	}()

	var cursor uint64
	var wg sync.WaitGroup

	for {
		var keys []string
		var err error
		keys, cursor, err = client.Scan(context.Background(), cursor, "cache:*:tags", 20).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			wg.Add(1)
			go func(keys, tags []string, wg *sync.WaitGroup) {
				defer wg.Done()

				pipe := client.Pipeline()
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
						deleteCacheByCacheTagsKey(keys[index])
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

func deleteCaches(cacheRedisKeys ...string) error {
	var keys []string
	for _, cacheKey := range cacheRedisKeys {
		cacheHeadersKey := fmt.Sprintf("%s:headers", cacheKey)
		cacheTagsKey := fmt.Sprintf("%s:tags", cacheKey)

		newKeys := []string{cacheKey, cacheHeadersKey, cacheTagsKey}
		keys = append(keys, newKeys...)
	}

	return client.Del(context.Background(), keys...).Err()
}

func deleteCacheByCacheTagsKey(cacheTagsKey string) error {
	parts := strings.Split(cacheTagsKey, ":tags")
	return deleteCaches(parts[0])
}
