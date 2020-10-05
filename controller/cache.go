package controller

import (
	"net/http"

	"github.com/gotway/gotway/core"
	"github.com/gotway/gotway/model"
)

// CacheControllerI interface
type CacheControllerI interface {
	IsCacheableRequest(r *http.Request) bool
	GetCache(r *http.Request, serviceKey string) (core.Cache, error)
	GetCacheDetail(r *http.Request, serviceKey string) (core.CacheDetail, error)
	DeleteCacheByPath(paths []core.CachePath) error
	DeleteCacheByTags(tags []string) error
	ListenResponses()
	HandleResponse(serviceKey string, r *http.Response) error
}

// CacheController controller
type CacheController struct {
	cacheConfigRepository model.CacheConfigRepositoryI
	cacheRepository       model.CacheRepositoryI
	resChan               chan response
}

// NewCacheController creates a new cache controller
func NewCacheController(cacheConfigRepository model.CacheConfigRepositoryI,
	cacheRepository model.CacheRepositoryI) CacheController {
	return CacheController{
		cacheConfigRepository: cacheConfigRepository,
		cacheRepository:       cacheRepository,
		resChan:               make(chan response),
	}
}

// IsCacheableRequest determines if a request's response can be retrieved from cache
func (c CacheController) IsCacheableRequest(r *http.Request) bool {
	return r.Method == http.MethodGet
}

// GetCache gets a cached response for a request
func (c CacheController) GetCache(r *http.Request, serviceKey string) (core.Cache, error) {
	path, err := core.GetServiceRelativePath(r, serviceKey)
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
func (c CacheController) GetCacheDetail(r *http.Request, serviceKey string) (core.CacheDetail, error) {
	path, err := core.GetServiceRelativePath(r, serviceKey)
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
