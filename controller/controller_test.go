package controller

import (
	"testing"

	"github.com/gotway/gotway/model"

	mocks "github.com/gotway/gotway/mocks/model"
)

func TestInit(t *testing.T) {
	serviceRepository := new(mocks.ServiceRepositoryI)
	cacheConfigRepository := new(mocks.CacheConfigRepositoryI)
	cacheRepository := new(mocks.CacheRepositoryI)

	model.ServiceRepository = serviceRepository
	model.CacheConfigRepository = cacheConfigRepository
	model.CacheRepository = cacheRepository

	Init()
}
