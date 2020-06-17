package model

// ServiceDaoI interface
type ServiceDaoI interface {
	StoreService(key string, url string, healthURL string) error
	GetService(key string) (*Service, error)
	getAllServices() []string
	updateServiceStatus(key string, status ServiceStatus)
}
