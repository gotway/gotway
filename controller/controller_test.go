package controller

import (
	"testing"

	"github.com/gotway/gotway/model"

	mocks "github.com/gotway/gotway/mocks/model"
)

func TestInit(t *testing.T) {
	serviceRepository := new(mocks.ServiceRepositoryI)
	cacheRepository := new(mocks.CacheRepositoryI)

	model.ServiceRepository = serviceRepository
	model.CacheRepository = cacheRepository

	Init()
}
