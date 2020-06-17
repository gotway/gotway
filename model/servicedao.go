package model

// ServiceDaoI interface
type ServiceDaoI interface {
	StoreService(key string, url string, healthURL string) bool
	getAllServices() []string
	getStatusAndHealthURL(redisKey string) (string, string)
	updateServiceStatus(redisKey string, status serviceStatus)
}
