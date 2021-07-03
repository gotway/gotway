package service

import (
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/repository"
	"github.com/gotway/gotway/pkg/log"
)

type Controller interface {
	CreateService(service model.Service) error
	GetServices() ([]model.Service, error)
	GetService(key string) (model.Service, error)
	DeleteService(key string) error
	UpsertService(service model.Service) error
}

type BasicController struct {
	serviceRepo repository.ServiceRepo
	logger      log.Logger
}

// CreateService creates a new service
func (c BasicController) CreateService(service model.Service) error {
	service.Status = model.ServiceStatusIdle
	return c.serviceRepo.Create(service)
}

// GetServices gets services paginated
func (c BasicController) GetServices() ([]model.Service, error) {
	return c.serviceRepo.GetAll()
}

// GetService gets a service
func (c BasicController) GetService(key string) (model.Service, error) {
	return c.serviceRepo.Get(key)
}

// DeleteService deletes a service
func (c BasicController) DeleteService(key string) error {
	return c.serviceRepo.Delete(key)
}

// UpdateServiceStatus updates the status of a service
func (c BasicController) UpsertService(service model.Service) error {
	return c.serviceRepo.Upsert(service)
}

func NewController(
	serviceRepo repository.ServiceRepo,
	logger log.Logger,
) Controller {
	return BasicController{serviceRepo, logger}
}
