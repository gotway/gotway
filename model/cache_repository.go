package model

import (
	"github.com/gosmo-devs/microgateway/core"
)

type deleteItem struct {
	path       string
	serviceKey string
}

// CacheRepositoryI implementation
type CacheRepositoryI interface {
	StoreCache(cache core.CacheDetail, serviceKey string) error
	GetCache(path string, serviceKey string) (core.Cache, error)
	GetCacheDetail(path string, serviceKey string) (core.CacheDetail, error)
	DeleteCacheByPath(paths []core.CachePath) error
	DeleteCacheByTags(tags []string) error
}
