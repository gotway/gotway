package cache_test

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gotway/gotway/internal/cache"
	"github.com/gotway/gotway/internal/mocks"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type cacheResponse struct {
	httpResponse *http.Response
	bodyBytes    []byte
	params       cache.Params
}

func TestIsCacheable(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	controller := cache.NewController(cache.Options{10, 10}, cacheRepo, log.Log)

	getReq, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/service/foo", nil)
	postReq, _ := http.NewRequest(http.MethodPost, "http://api.gotway.com/service/foo", nil)

	tests := []struct {
		name     string
		req      *http.Request
		wantBool bool
	}{
		{
			name:     "GET",
			req:      getReq,
			wantBool: true,
		},
		{
			name:     "POST",
			req:      postReq,
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isCacheable := controller.IsCacheableRequest(tt.req)

			assert.Equal(t, isCacheable, tt.wantBool)
		})
	}
}

func TestGetCache(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	controller := cache.NewController(cache.Options{10, 10}, cacheRepo, log.Log)

	reqCacheError, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/foo", nil)
	cacheError := errors.New("Cache not found")
	cacheRepo.On("Get", "/foo", "service").Return(model.Cache{}, cacheError)

	reqSuccess, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/products", nil)
	reqPrefix, _ := http.NewRequest(
		http.MethodGet,
		"http://api.gotway.com/products",
		nil,
	)
	cache := model.Cache{
		Path:       "/products",
		StatusCode: 200,
		TTL:        10,
		Tags:       []string{"foo"},
	}
	cacheRepo.On("Get", "/products", "catalog").Return(cache, nil)

	tests := []struct {
		name      string
		req       *http.Request
		service   string
		wantCache model.Cache
		wantErr   error
	}{
		{
			name:      "Cache detail not found error",
			req:       reqCacheError,
			service:   "service",
			wantCache: model.Cache{},
			wantErr:   cacheError,
		},
		{
			name:      "Get cache successfully",
			req:       reqSuccess,
			service:   "catalog",
			wantCache: cache,
			wantErr:   nil,
		},
		{
			name:      "Get cache successfully",
			req:       reqPrefix,
			service:   "catalog",
			wantCache: cache,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheDetail, err := controller.GetCache(tt.req, tt.service)

			assert.Equal(t, tt.wantCache, cacheDetail)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestDeleteCacheByPath(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	controller := cache.NewController(cache.Options{10, 10}, cacheRepo, log.Log)

	paths := []model.CachePath{
		{
			Service: "service",
			Path:    "foo",
		},
	}
	cacheRepo.On("DeleteByPath", paths).Return(nil)

	err := controller.DeleteCacheByPath(paths)

	assert.Nil(t, err)
	cacheRepo.AssertExpectations(t)
}

func TestDeleteCacheByTags(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	controller := cache.NewController(cache.Options{10, 10}, cacheRepo, log.Log)

	tags := []string{"foo"}
	cacheRepo.On("DeleteByTags", tags).Return(nil)

	err := controller.DeleteCacheByTags(tags)

	assert.Nil(t, err)
	cacheRepo.AssertExpectations(t)
}

func TestListenResponses(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	controller := cache.NewController(cache.Options{10, 10}, cacheRepo, log.Log)

	bodyBytes := []byte("{}")
	body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	url, _ := url.Parse("http://api.gotway.com/products?offset=0&limit=10")
	httpRes := &http.Response{
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    url,
		},
		StatusCode: http.StatusOK,
		Body:       body,
	}
	catalogParams := cache.Params{
		Service:  "catalog",
		Statuses: []int{http.StatusOK},
	}
	stockParams := cache.Params{
		Service: "stock",
	}
	cacheableRes := cacheResponse{
		httpResponse: httpRes,
		bodyBytes:    bodyBytes,
		params:       catalogParams,
	}
	nonCacheableRes := cacheResponse{
		httpResponse: httpRes,
		bodyBytes:    bodyBytes,
		params:       stockParams,
	}

	cacheRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go controller.Start(ctx)

	for _, r := range []cacheResponse{cacheableRes, nonCacheableRes} {
		if err := controller.HandleResponse(r.httpResponse, r.params); err != nil {
			t.Errorf("got unexpected error: %v", err)
		}
	}

	time.Sleep(1 * time.Second)
	cacheRepo.AssertNumberOfCalls(t, "Create", 1)
}

func TestListenCacheControlResponses(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	controller := cache.NewController(cache.Options{10, 10}, cacheRepo, log.Log)

	params := cache.Params{
		Service:  "foo",
		Statuses: []int{http.StatusOK},
	}
	bodyBytes := []byte("{}")
	body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	url, _ := url.Parse("http://api.gotway.com/catalog/products")
	header := http.Header{}
	header.Set("Cache-Control", "s-maxage=10")
	TTLRes := cacheResponse{
		httpResponse: &http.Response{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    url,
			},
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       body,
		},
		bodyBytes: bodyBytes,
		params:    params,
	}
	noTTLRes := cacheResponse{
		httpResponse: &http.Response{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    url,
			},
			StatusCode: http.StatusOK,
			Body:       body,
		},
		bodyBytes: bodyBytes,
		params:    params,
	}
	zeroTTLHeader := http.Header{}
	zeroTTLHeader.Set("Cache-Control", "s-maxage=0")
	zeroTTLRes := cacheResponse{
		httpResponse: &http.Response{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    url,
			},
			StatusCode: http.StatusOK,
			Header:     zeroTTLHeader,
			Body:       body,
		},
		bodyBytes: bodyBytes,
		params:    params,
	}

	cacheRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go controller.Start(ctx)

	for _, r := range []cacheResponse{TTLRes, noTTLRes, zeroTTLRes} {
		if err := controller.HandleResponse(r.httpResponse, r.params); err != nil {
			t.Errorf("Got unexpected error: %v", err)
		}
	}

	time.Sleep(1 * time.Second)
	cacheRepo.AssertNumberOfCalls(t, "Create", 2)
}

func TestListenCacheTagsResponses(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	controller := cache.NewController(cache.Options{10, 10}, cacheRepo, log.Log)

	params := cache.Params{
		Service:  "foo",
		Statuses: []int{http.StatusOK},
	}
	bodyBytes := []byte("{}")
	body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	url, _ := url.Parse("http://api.gotway.com/catalog/products")
	header := http.Header{}
	header.Set("X-Cache-Tags", "products")
	tagsRes := cacheResponse{
		httpResponse: &http.Response{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    url,
			},
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       body,
		},
		bodyBytes: bodyBytes,
		params:    params,
	}
	noTagsRes := cacheResponse{
		httpResponse: &http.Response{
			Request: &http.Request{
				Method: http.MethodGet,
				URL:    url,
			},
			StatusCode: http.StatusOK,
			Body:       body,
		},
		bodyBytes: bodyBytes,
		params:    params,
	}

	cacheRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go controller.Start(ctx)

	for _, r := range []cacheResponse{tagsRes, noTagsRes} {
		if err := controller.HandleResponse(r.httpResponse, r.params); err != nil {
			t.Errorf("got unexpected error: %v", err)
		}
	}

	time.Sleep(1 * time.Second)
	cacheRepo.AssertNumberOfCalls(t, "Create", 2)
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Error reading")
}

func TestErrReadingBody(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	controller := cache.NewController(cache.Options{10, 10}, cacheRepo, log.Log)

	url, _ := url.Parse("http://api.gotway.com/catalog/products")
	testRequest := httptest.NewRequest(http.MethodPost, "/foo", errReader(0))
	defer testRequest.Body.Close()

	params := cache.Params{
		Service:  "service",
		Statuses: []int{http.StatusOK},
	}
	res := &http.Response{
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    url,
		},
		StatusCode: http.StatusOK,
		Body:       testRequest.Body,
	}

	cacheRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go controller.Start(ctx)

	err := controller.HandleResponse(res, params)
	assert.NotNil(t, err)
}

func TestCachePolicy(t *testing.T) {
	cacheRepo := new(mocks.CacheRepo)
	controller := cache.NewController(
		cache.Options{10, 10},
		cacheRepo,
		log.Log,
	)

	url, _ := url.Parse("http://api.gotway.com/catalog/products")
	notCacheableHeader := http.Header{}
	notCacheableHeader.Set("Cache-Control", "s-maxage=0")
	cacheableHeader := http.Header{}
	cacheableHeader.Set("Cache-Control", "s-maxage=10")
	invalidCacheableHeader := http.Header{}
	invalidCacheableHeader.Set("Cache-Control", "s-maxage")

	cacheableParams := cache.Params{
		Statuses: []int{http.StatusOK},
	}
	notCacheableService := cache.Params{}

	tests := []struct {
		name            string
		response        *http.Response
		params          cache.Params
		wantIsCacheable bool
	}{
		{
			name:   "Not cacheable by method",
			params: cacheableParams,
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodPost,
					URL:    url,
				},
				StatusCode: http.StatusOK,
				Header:     cacheableHeader,
			},
			wantIsCacheable: false,
		},
		{
			name:   "Not cacheable by headers",
			params: cacheableParams,
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    url,
				},
				StatusCode: http.StatusOK,
				Header:     notCacheableHeader,
			},
			wantIsCacheable: false,
		},
		{
			name:   "Not cacheable by status",
			params: notCacheableService,
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    url,
				},
				StatusCode: http.StatusBadRequest,
				Header:     cacheableHeader,
			},
			wantIsCacheable: false,
		},
		{
			name:   "Cacheable",
			params: cacheableParams,
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    url,
				},
				StatusCode: http.StatusOK,
				Header:     cacheableHeader,
			},
			wantIsCacheable: true,
		},
		{
			name:   "Error parsing cache header",
			params: cacheableParams,
			response: &http.Response{
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    url,
				},
				StatusCode: http.StatusOK,
				Header:     invalidCacheableHeader,
			},
			wantIsCacheable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isCacheable := controller.IsCacheableResponse(tt.response, tt.params)

			assert.Equal(t, tt.wantIsCacheable, isCacheable)
		})
	}

}
