package cache

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gotway/gotway/internal/model"
	"github.com/pquerna/cachecontrol/cacheobject"
)

type response struct {
	serviceKey   string
	httpResponse *http.Response
	body         *io.ReadCloser
}

// ListenResponses starts listening for responses
func (c BasicController) ListenResponses(ctx context.Context) {
	c.logger.Info("starting cache handler")
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("stopping cache handler")
			return
		case response := <-c.resChan:
			c.logger.Debug("caching response")
			if err := c.cacheResponse(response); err != nil {
				c.logger.Error("error caching response", err)
			}
		}
	}
}

// HandleResponse handles a response ans sends it to the channel
func (c BasicController) HandleResponse(serviceKey string, r *http.Response) error {
	if !c.isCacheableResponse(r, serviceKey) {
		return nil
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	response := response{
		serviceKey:   serviceKey,
		httpResponse: r,
		body:         &body,
	}

	c.resChan <- response

	return nil
}

// IsCacheableResponse checks if a response is cacheable
func (c BasicController) isCacheableResponse(r *http.Response, serviceKey string) bool {
	if !c.IsCacheableRequest(r.Request) || headersDisallowCaching(r) {
		return false
	}
	return c.serviceRepo.IsCacheableStatusCode(serviceKey, r.StatusCode)
}

func (c BasicController) cacheResponse(res response) error {
	config, err := c.serviceRepo.GetServiceCache(res.serviceKey)
	if err != nil {
		return err
	}
	path := getPath(res.httpResponse.Request)
	ttl := getTTL(res.httpResponse, config)
	tags := getTags(res.httpResponse, config)

	cache := model.CacheDetail{
		Cache: model.Cache{
			Path:       path,
			StatusCode: res.httpResponse.StatusCode,
			Headers:    res.httpResponse.Header,
			Body: model.CacheBody{
				Reader: *res.body,
			},
		},
		TTL:  ttl,
		Tags: tags,
	}

	return c.cacheRepo.StoreCache(cache, res.serviceKey)
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
