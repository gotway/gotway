// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	core "github.com/gotway/gotway/core"
	mock "github.com/stretchr/testify/mock"
)

// CacheRepositoryI is an autogenerated mock type for the CacheRepositoryI type
type CacheRepositoryI struct {
	mock.Mock
}

// DeleteCacheByPath provides a mock function with given fields: paths
func (_m *CacheRepositoryI) DeleteCacheByPath(paths []core.CachePath) error {
	ret := _m.Called(paths)

	var r0 error
	if rf, ok := ret.Get(0).(func([]core.CachePath) error); ok {
		r0 = rf(paths)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteCacheByTags provides a mock function with given fields: tags
func (_m *CacheRepositoryI) DeleteCacheByTags(tags []string) error {
	ret := _m.Called(tags)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(tags)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCache provides a mock function with given fields: path, serviceKey
func (_m *CacheRepositoryI) GetCache(path string, serviceKey string) (core.Cache, error) {
	ret := _m.Called(path, serviceKey)

	var r0 core.Cache
	if rf, ok := ret.Get(0).(func(string, string) core.Cache); ok {
		r0 = rf(path, serviceKey)
	} else {
		r0 = ret.Get(0).(core.Cache)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(path, serviceKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCacheDetail provides a mock function with given fields: path, serviceKey
func (_m *CacheRepositoryI) GetCacheDetail(path string, serviceKey string) (core.CacheDetail, error) {
	ret := _m.Called(path, serviceKey)

	var r0 core.CacheDetail
	if rf, ok := ret.Get(0).(func(string, string) core.CacheDetail); ok {
		r0 = rf(path, serviceKey)
	} else {
		r0 = ret.Get(0).(core.CacheDetail)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(path, serviceKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreCache provides a mock function with given fields: cache, serviceKey
func (_m *CacheRepositoryI) StoreCache(cache core.CacheDetail, serviceKey string) error {
	ret := _m.Called(cache, serviceKey)

	var r0 error
	if rf, ok := ret.Get(0).(func(core.CacheDetail, string) error); ok {
		r0 = rf(cache, serviceKey)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
