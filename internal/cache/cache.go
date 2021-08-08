package cache

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/repository"
	"github.com/gotway/gotway/pkg/log"
	"github.com/pquerna/cachecontrol/cacheobject"
)

type Options struct {
	NumWorkers int
	BufferSize int
}

type Params struct {
	Service  string
	TTL      int64
	Statuses []int
	Tags     []string
}

type Controller interface {
	Start(ctx context.Context)
	HandleResponse(r *http.Response, params Params) error
	IsCacheableRequest(r *http.Request) bool
	IsCacheableResponse(r *http.Response, params Params) bool
	GetCache(r *http.Request, service string) (model.Cache, error)
	DeleteCacheByPath(paths []model.CachePath) error
	DeleteCacheByTags(tags []string) error
}

type response struct {
	httpResponse *http.Response
	bodyBytes    []byte
	params       Params
}

type BasicController struct {
	options      Options
	cacheRepo    repository.CacheRepo
	pendingCache chan response
	logger       log.Logger
}

// ListenResponses starts listening for responses
func (c BasicController) Start(ctx context.Context) {
	c.logger.Info("starting cache controller")
	var logOnce sync.Once

	for i := 0; i < c.options.NumWorkers; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					logOnce.Do(func() {
						c.logger.Info("stopping cache controller")
					})
					return
				case response := <-c.pendingCache:
					c.logger.Debug("caching response")
					if err := c.cacheResponse(response); err != nil {
						c.logger.Error("error caching response", err)
					}
				}
			}
		}()
	}
}

// HandleResponse handles a response ans sends it to the channel
func (c BasicController) HandleResponse(r *http.Response, params Params) error {
	if !c.IsCacheableResponse(r, params) {
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
		params:       params,
	}

	return nil
}

// IsCacheableRequest determines if a request's response can be retrieved from cache
func (c BasicController) IsCacheableRequest(r *http.Request) bool {
	return r.Method == http.MethodGet
}

// IsCacheableResponse determines if a response can be stored in cache
func (c BasicController) IsCacheableResponse(r *http.Response, params Params) bool {
	if !c.IsCacheableRequest(r.Request) || headersDisallowCaching(r) {
		return false
	}
	for _, s := range params.Statuses {
		if s == r.StatusCode {
			return true
		}
	}
	return false
}

// GetCache gets a cached response for a request and a service
func (c BasicController) GetCache(r *http.Request, service string) (model.Cache, error) {
	cache, err := c.cacheRepo.Get(r.URL.Path, service)
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

func (c BasicController) cacheResponse(res response) error {
	path := getPath(res.httpResponse.Request)
	ttl := getTTL(res.httpResponse, res.params)
	tags := getTags(res.httpResponse, res.params)

	cache := model.Cache{
		Path:       path,
		StatusCode: res.httpResponse.StatusCode,
		Headers:    res.httpResponse.Header,
		Body:       res.bodyBytes,
		TTL:        ttl,
		Tags:       tags,
	}

	return c.cacheRepo.Create(cache, res.params.Service)
}

func getPath(r *http.Request) string {
	path := r.URL.Path
	if r.URL.RawQuery != "" {
		return fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
	}
	return path
}

func getTTL(r *http.Response, params Params) model.CacheTTL {
	ttl, err := getCacheTTLHeader(r)
	var seconds int64
	if err != nil {
		seconds = params.TTL
	} else {
		seconds = ttl
	}
	return model.NewCacheTTL(seconds)
}

func getTags(r *http.Response, params Params) []string {
	tags, err := getCacheTagsHeader(r)
	if err != nil {
		return params.Tags
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
	options Options,
	cacheRepo repository.CacheRepo,
	logger log.Logger,
) Controller {
	return &BasicController{
		options:      options,
		cacheRepo:    cacheRepo,
		pendingCache: make(chan response, options.BufferSize),
		logger:       logger,
	}
}
