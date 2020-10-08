package controller

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gotway/gotway/core"
	mocks "github.com/gotway/gotway/mocks/model"
	"github.com/stretchr/testify/assert"
)

func TestIsCacheable(t *testing.T) {
	cacheConfigRepo := new(mocks.CacheConfigRepositoryI)
	cacheRepo := new(mocks.CacheRepositoryI)
	controller := NewCacheController(cacheConfigRepo, cacheRepo)

	cacheableReq, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/service/foo", nil)
	notCacheableReq, _ := http.NewRequest(http.MethodPost, "http://api.gotway.com/service/foo", nil)

	tests := []struct {
		name     string
		req      *http.Request
		wantBool bool
	}{
		{
			name:     "Is GET cacheable",
			req:      cacheableReq,
			wantBool: true,
		},
		{
			name:     "Is POST cacheable",
			req:      notCacheableReq,
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
	cacheConfigRepo := new(mocks.CacheConfigRepositoryI)
	cacheRepo := new(mocks.CacheRepositoryI)
	controller := NewCacheController(cacheConfigRepo, cacheRepo)

	reqCacheError, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/service/foo", nil)
	cacheError := errors.New("Cache not found")
	cacheRepo.On("GetCache", "/foo", "service").Return(core.Cache{}, cacheError)

	reqSuccess, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/catalog/products", nil)
	cache := core.Cache{
		Path:       "/products",
		StatusCode: 200,
	}
	cacheRepo.On("GetCache", "/products", "catalog").Return(cache, nil)

	tests := []struct {
		name       string
		req        *http.Request
		serviceKey string
		wantCache  core.Cache
		wantErr    error
	}{
		{
			name:       "Service path error",
			req:        reqSuccess,
			serviceKey: "bar",
			wantCache:  core.Cache{},
			wantErr: &core.ErrServiceNotFoundInURL{
				URL:         reqSuccess.URL,
				ServicePath: "bar",
			},
		},
		{
			name:       "Cache not found error",
			req:        reqCacheError,
			serviceKey: "service",
			wantCache:  core.Cache{},
			wantErr:    cacheError,
		},
		{
			name:       "Get cache successfully",
			req:        reqSuccess,
			serviceKey: "catalog",
			wantCache:  cache,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := controller.GetCache(tt.req, tt.serviceKey)

			assert.Equal(t, tt.wantCache, cache)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestGetCacheDetail(t *testing.T) {
	cacheConfigRepo := new(mocks.CacheConfigRepositoryI)
	cacheRepo := new(mocks.CacheRepositoryI)
	controller := NewCacheController(cacheConfigRepo, cacheRepo)

	reqCacheError, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/service/foo", nil)
	cacheError := errors.New("Cache not found")
	cacheRepo.On("GetCacheDetail", "/foo", "service").Return(core.CacheDetail{}, cacheError)

	reqSuccess, _ := http.NewRequest(http.MethodGet, "http://api.gotway.com/catalog/products", nil)
	cacheDetail := core.CacheDetail{
		Cache: core.Cache{
			Path:       "/products",
			StatusCode: 200,
		},
		TTL:  10,
		Tags: []string{"foo"},
	}
	cacheRepo.On("GetCacheDetail", "/products", "catalog").Return(cacheDetail, nil)

	tests := []struct {
		name            string
		req             *http.Request
		serviceKey      string
		wantCacheDetail core.CacheDetail
		wantErr         error
	}{
		{
			name:            "Service path error",
			req:             reqSuccess,
			serviceKey:      "bar",
			wantCacheDetail: core.CacheDetail{},
			wantErr: &core.ErrServiceNotFoundInURL{
				URL:         reqSuccess.URL,
				ServicePath: "bar",
			},
		},
		{
			name:            "Cache detail not found error",
			req:             reqCacheError,
			serviceKey:      "service",
			wantCacheDetail: core.CacheDetail{},
			wantErr:         cacheError,
		},
		{
			name:            "Get cache successfully",
			req:             reqSuccess,
			serviceKey:      "catalog",
			wantCacheDetail: cacheDetail,
			wantErr:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheDetail, err := controller.GetCacheDetail(tt.req, tt.serviceKey)

			assert.Equal(t, tt.wantCacheDetail, cacheDetail)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestDeleteCacheByPath(t *testing.T) {
	cacheConfigRepo := new(mocks.CacheConfigRepositoryI)
	cacheRepo := new(mocks.CacheRepositoryI)
	controller := NewCacheController(cacheConfigRepo, cacheRepo)

	paths := []core.CachePath{
		{
			ServicePath: "service",
			Path:        "foo",
		},
	}
	cacheRepo.On("DeleteCacheByPath", paths).Return(nil)

	err := controller.DeleteCacheByPath(paths)

	assert.Nil(t, err)
	cacheRepo.AssertExpectations(t)
}

func TestDeleteCacheByTags(t *testing.T) {
	cacheConfigRepo := new(mocks.CacheConfigRepositoryI)
	cacheRepo := new(mocks.CacheRepositoryI)
	controller := NewCacheController(cacheConfigRepo, cacheRepo)

	tags := []string{"foo"}
	cacheRepo.On("DeleteCacheByTags", tags).Return(nil)

	err := controller.DeleteCacheByTags(tags)

	assert.Nil(t, err)
	cacheRepo.AssertExpectations(t)
}
