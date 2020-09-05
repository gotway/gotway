package model

// ServiceDaoI interface
type ServiceDaoI interface {
	StoreService(service Service) error
	GetService(key string) (*Service, error)
	GetAllServices() []string
	UpdateServiceStatus(key string, status ServiceStatus)
}
