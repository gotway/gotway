package cache

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/repository"
	"github.com/gotway/gotway/pkg/log"
	"github.com/pquerna/cachecontrol/cacheobject"
)

type Controller interface {
	IsCacheableRequest(r *http.Request, service model.Service) bool
	GetCache(r *http.Request, service model.Service) (model.Cache, error)
	DeleteCacheByPath(paths []model.CachePath) error
	DeleteCacheByTags(tags []string) error
	ListenResponses(ctx context.Context)
	HandleResponse(r *http.Response, service model.Service) error
}

type response struct {
	httpResponse *http.Response
	bodyBytes    []byte
	service      model.Service
}

type BasicController struct {
	cacheRepo    repository.CacheRepo
	serviceRepo  repository.ServiceRepo
	pendingCache chan response
	logger       log.Logger
}

// IsCacheableRequest determines if a request's response can be retrieved from cache
func (c BasicController) IsCacheableRequest(r *http.Request, service model.Service) bool {
	return r.Method == http.MethodGet && service.Type == model.ServiceTypeREST
}

// GetCache gets a cached response for a request and a service
func (c BasicController) GetCache(r *http.Request, service model.Service) (model.Cache, error) {
	cache, err := c.cacheRepo.Get(r.URL.Path, service.Name)
	if err != nil {
		return model.Cache{}, err
	}
	return cache, nil
}

// DeleteCacheByPath deletes cache defined by its path
func (c BasicController) DeleteCacheByPath(paths []model.CachePath) error {
	return c.cacheRepo.DeleteByPath(paths)
}

// DeleteCacheByTags deletes cache with tags
func (c BasicController) DeleteCacheByTags(tags []string) error {
	return c.cacheRepo.DeleteByTags(tags)
}

// ListenResponses starts listening for responses
func (c BasicController) ListenResponses(ctx context.Context) {
	c.logger.Info("starting cache handler")
	for i := 0; i < config.CacheNumWorkers; i++ {
		go func() {
			c.checkCache(ctx)
		}()
	}
}

// HandleResponse handles a response ans sends it to the channel
func (c BasicController) HandleResponse(r *http.Response, service model.Service) error {
	if !c.isCacheableResponse(r, service) {
		return nil
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	c.pendingCache <- response{
		httpResponse: r,
		bodyBytes:    bodyBytes,
		service:      service,
	}

	return nil
}

func (c BasicController) isCacheableResponse(r *http.Response, service model.Service) bool {
	if !c.IsCacheableRequest(r.Request, service) || headersDisallowCaching(r) {
		return false
	}
	for _, s := range service.Cache.Statuses {
		if s == r.StatusCode {
			return true
		}
	}
	return false
}

func (c BasicController) checkCache(ctx context.Context) {
	c.logger.Info("stopping cache handler")
	for {
		select {
		case <-ctx.Done():
			return
		case response := <-c.pendingCache:
			c.logger.Debug("caching response")
			if err := c.cacheResponse(response); err != nil {
				c.logger.Error("error caching response", err)
			}
		}
	}
}

func (c BasicController) cacheResponse(res response) error {
	path := getPath(res.httpResponse.Request)
	ttl := getTTL(res.httpResponse, res.service.Cache)
	tags := getTags(res.httpResponse, res.service.Cache)

	cache := model.Cache{
		Path:       path,
		StatusCode: res.httpResponse.StatusCode,
		Body:       res.bodyBytes,
		TTL:        ttl,
		Tags:       tags,
	}

	return c.cacheRepo.Create(cache, res.service.Name)
}

func getPath(r *http.Request) string {
	path := r.URL.Path
	query := r.URL.RawQuery
	if query != "" {
		return fmt.Sprintf("%s?%s", path, query)
	}
	return path
}

func getTTL(r *http.Response, config model.CacheConfig) model.CacheTTL {
	ttl, err := getCacheTTLHeader(r)
	var seconds int64
	if err != nil {
		seconds = config.TTL
	} else {
		seconds = ttl
	}
	return model.NewCacheTTL(seconds)
}

func getTags(r *http.Response, config model.CacheConfig) []string {
	tags, err := getCacheTagsHeader(r)
	if err != nil {
		return config.Tags
	}
	return tags
}

func getCacheTTLHeader(r *http.Response) (int64, error) {
	cacheControl := r.Header.Get("Cache-Control")
	directives, err := cacheobject.ParseResponseCacheControl(cacheControl)
	if err != nil {
		return 0, err
	}
	ttl := directives.SMaxAge
	if ttl < 0 {
		return 0, errors.New("Cache-Control header not found")
	}
	return int64(ttl), nil
}

func getCacheTagsHeader(r *http.Response) ([]string, error) {
	cacheTags := r.Header.Values("X-Cache-Tags")
	if len(cacheTags) == 0 {
		return cacheTags, errors.New("X-Cache-Tags header not found")
	}
	return cacheTags, nil
}

func headersDisallowCaching(r *http.Response) bool {
	ttl, err := getCacheTTLHeader(r)
	if err != nil {
		return false
	}
	return ttl == 0
}

func NewController(
	cacheRepo repository.CacheRepo,
	serviceRepo repository.ServiceRepo,
	logger log.Logger,
) Controller {

	return &BasicController{
		cacheRepo:    cacheRepo,
		serviceRepo:  serviceRepo,
		pendingCache: make(chan response, config.CacheBufferSize),
		logger:       logger,
	}
}
