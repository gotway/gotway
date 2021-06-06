package controller

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gotway/gotway/core"
	"github.com/gotway/gotway/log"
	"github.com/pquerna/cachecontrol/cacheobject"
)

type response struct {
	serviceKey   string
	httpResponse *http.Response
	body         *io.ReadCloser
}

// ListenResponses starts listening for responses
func (c CacheController) ListenResponses() {
	go func() {
		for {
			select {
			case response := <-c.resChan:
				log.Logger.Debug("Caching response")
				err := c.cacheResponse(response)
				if err != nil {
					log.Logger.Error(err)
				}
			}
		}
	}()
}

// HandleResponse handles a response ans sends it to the channel
func (c CacheController) HandleResponse(serviceKey string, r *http.Response) error {
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
func (c CacheController) isCacheableResponse(r *http.Response, serviceKey string) bool {
	if !c.IsCacheableRequest(r.Request) || headersDisallowCaching(r) {
		return false
	}
	return c.serviceRepository.IsCacheableStatusCode(serviceKey, r.StatusCode)
}

func (c CacheController) cacheResponse(res response) error {
	config, err := c.serviceRepository.GetServiceCache(res.serviceKey)
	if err != nil {
		return err
	}
	path := getPath(res.httpResponse.Request)
	ttl := getTTL(res.httpResponse, config)
	tags := getTags(res.httpResponse, config)

	cache := core.CacheDetail{
		Cache: core.Cache{
			Path:       path,
			StatusCode: res.httpResponse.StatusCode,
			Headers:    res.httpResponse.Header,
			Body: core.CacheBody{
				Reader: *res.body,
			},
		},
		TTL:  ttl,
		Tags: tags,
	}

	return c.cacheRepository.StoreCache(cache, res.serviceKey)
}

func getPath(r *http.Request) string {
	path := r.URL.Path
	query := r.URL.RawQuery
	if query != "" {
		return fmt.Sprintf("%s?%s", path, query)
	}
	return path
}

func getTTL(r *http.Response, config core.CacheConfig) core.CacheTTL {
	ttl, err := getCacheTTLHeader(r)
	var seconds int64
	if err != nil {
		seconds = config.TTL
	} else {
		seconds = ttl
	}
	return core.NewCacheTTL(seconds)
}

func getTags(r *http.Response, config core.CacheConfig) []string {
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
