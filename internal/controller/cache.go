package controller

import (
	"context"
	"net/http"

	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/repository"
	"github.com/gotway/gotway/pkg/log"
)

type CacheController interface {
	IsCacheableRequest(r *http.Request) bool
	GetCache(r *http.Request, pathPrefix, serviceKey string) (model.Cache, error)
	GetCacheDetail(r *http.Request, pathPrefix, serviceKey string) (model.CacheDetail, error)
	DeleteCacheByPath(paths []model.CachePath) error
	DeleteCacheByTags(tags []string) error
	ListenResponses(ctx context.Context)
	HandleResponse(serviceKey string, r *http.Response) error
	isCacheableResponse(r *http.Response, serviceKey string) bool
}

type BasicCacheController struct {
	cacheRepo   repository.CacheRepo
	serviceRepo repository.ServiceRepo
	resChan     chan response
	logger      log.Logger
}

// IsCacheableRequest determines if a request's response can be retrieved from cache
func (c BasicCacheController) IsCacheableRequest(r *http.Request) bool {
	return r.Method == http.MethodGet
}

// GetCache gets a cached response for a request
func (c BasicCacheController) GetCache(
	r *http.Request,
	pathPrefix, serviceKey string,
) (model.Cache, error) {
	path, err := model.GetServiceRelativePathPrefixed(r, pathPrefix, serviceKey)
	if err != nil {
		return model.Cache{}, err
	}
	cache, err := c.cacheRepo.GetCache(path, serviceKey)
	if err != nil {
		return model.Cache{}, err
	}
	return cache, nil
}

// GetCacheDetail gets a cache with extra info
func (c BasicCacheController) GetCacheDetail(
	r *http.Request,
	pathPrefix, serviceKey string,
) (model.CacheDetail, error) {
	path, err := model.GetServiceRelativePathPrefixed(r, pathPrefix, serviceKey)
	if err != nil {
		return model.CacheDetail{}, err
	}
	cacheDetail, err := c.cacheRepo.GetCacheDetail(path, serviceKey)
	if err != nil {
		return model.CacheDetail{}, err
	}
	return cacheDetail, nil
}

// DeleteCacheByPath deletes cache defined by its path
func (c BasicCacheController) DeleteCacheByPath(paths []model.CachePath) error {
	return c.cacheRepo.DeleteCacheByPath(paths)
}

// DeleteCacheByTags deletes cache with tags
func (c BasicCacheController) DeleteCacheByTags(tags []string) error {
	return c.cacheRepo.DeleteCacheByTags(tags)
}

func NewCacheController(
	cacheRepo repository.CacheRepo,
	serviceRepo repository.ServiceRepo,
	logger log.Logger,
) CacheController {

	return &BasicCacheController{
		cacheRepo:   cacheRepo,
		serviceRepo: serviceRepo,
		resChan:     make(chan response),
		logger:      logger,
	}
}
