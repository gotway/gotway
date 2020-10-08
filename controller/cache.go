package controller

import (
	"net/http"

	"github.com/gotway/gotway/core"
	"github.com/gotway/gotway/model"
)

// CacheController controller
type CacheController struct {
	cacheRepository   model.CacheRepositoryI
	serviceRepository model.ServiceRepositoryI
	resChan           chan response
}

func newCacheController(cacheRepository model.CacheRepositoryI,
	serviceRepository model.ServiceRepositoryI) CacheController {
	return CacheController{
		cacheRepository:   cacheRepository,
		serviceRepository: serviceRepository,
		resChan:           make(chan response),
	}
}

// IsCacheableRequest determines if a request's response can be retrieved from cache
func (c CacheController) IsCacheableRequest(r *http.Request) bool {
	return r.Method == http.MethodGet
}

// GetCache gets a cached response for a request
func (c CacheController) GetCache(r *http.Request, pathPrefix, serviceKey string) (core.Cache, error) {
	path, err := core.GetServiceRelativePathPrefixed(r, pathPrefix, serviceKey)
	if err != nil {
		return core.Cache{}, err
	}
	cache, err := c.cacheRepository.GetCache(path, serviceKey)
	if err != nil {
		return core.Cache{}, err
	}
	return cache, nil
}

// GetCacheDetail gets a cache with extra info
func (c CacheController) GetCacheDetail(r *http.Request, pathPrefix, serviceKey string) (core.CacheDetail, error) {
	path, err := core.GetServiceRelativePathPrefixed(r, pathPrefix, serviceKey)
	if err != nil {
		return core.CacheDetail{}, err
	}
	cacheDetail, err := c.cacheRepository.GetCacheDetail(path, serviceKey)
	if err != nil {
		return core.CacheDetail{}, err
	}
	return cacheDetail, nil
}

// DeleteCacheByPath deletes cache defined by its path
func (c CacheController) DeleteCacheByPath(paths []core.CachePath) error {
	return c.cacheRepository.DeleteCacheByPath(paths)
}

// DeleteCacheByTags deletes cache with tags
func (c CacheController) DeleteCacheByTags(tags []string) error {
	return c.cacheRepository.DeleteCacheByTags(tags)
}
