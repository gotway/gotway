package model

import "github.com/gosmo-devs/microgateway/core"

// ServiceRepositoryI interface
type ServiceRepositoryI interface {
	StoreService(service core.Service) error
	ExistService(key string) error
	GetAllServiceKeys() []string
	GetService(key string) (core.Service, error)
	GetServices(keys ...string) ([]core.Service, error)
	DeleteService(key string) error
	UpdateServiceStatus(key string, status core.ServiceStatus)
}
