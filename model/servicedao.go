package model

// ServiceDaoI interface
type ServiceDaoI interface {
	StoreService(service Service) error
	GetAllServiceKeys() []string
	GetService(key string) (*Service, error)
	GetServices(keys ...string) ([]Service, error)
	DeleteService(key string) error
	UpdateServiceStatus(key string, status ServiceStatus)
}
