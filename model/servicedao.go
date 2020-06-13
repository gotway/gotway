package model

// ServiceDaoI interface
type ServiceDaoI interface {
	StoreService(key string, url string, healthURL string) bool
}
