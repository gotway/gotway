package health

import (
	"sync"
	"time"

	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/log"
	"github.com/gosmo-devs/microgateway/model"
)

// Init initializes service health check
func Init() {
	ticker := time.NewTicker(config.HealthCheckInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Debug("Checking services health")
				updateServiceStatus()
			}
		}
	}()
}

func updateServiceStatus() {
	setToHealthy, setToIdle := getServicesToChangeStatus()

	for _, service := range setToHealthy {
		model.ServiceDao.UpdateServiceStatus(service, model.Healthy)
	}
	for _, service := range setToIdle {
		model.ServiceDao.UpdateServiceStatus(service, model.Idle)
	}
}

func getServicesToChangeStatus() (setToHealthy []string, setToIdle []string) {
	var healthyServices []string
	var idleServices []string

	services := model.ServiceDao.GetAllServiceKeys()
	var wg sync.WaitGroup
	for _, serviceKey := range services {
		wg.Add(1)
		go func(serviceKey string) {
			defer wg.Done()
			service, err := model.ServiceDao.GetService(serviceKey)
			if err != nil {
				log.Error(err)
				return
			}
			client, err := NewClient(service)
			if err != nil {
				log.Error(err)
				return
			}
			err = client.HealthCheck()
			if err != nil {
				if service.Status == model.Healthy {
					log.Infof("Service %s is now idle. Cause: %v", service.Path, err)
					idleServices = append(idleServices, service.Path)
				}
			} else {
				if service.Status == model.Idle {
					log.Infof("Service %s is now healthy", service.Path)
					healthyServices = append(healthyServices, service.Path)
				}
			}
		}(serviceKey)
	}
	wg.Wait()
	return healthyServices, idleServices
}
