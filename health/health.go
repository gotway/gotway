package health

import (
	"sync"
	"time"

	"github.com/gotway/gotway/controller"
	"github.com/gotway/gotway/core"

	"github.com/gotway/gotway/config"
	"github.com/gotway/gotway/log"
)

// Init initializes service health check
func Init() {
	ticker := time.NewTicker(config.HealthCheckInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Logger.Debug("Checking services health")
				updateServiceStatus()
			}
		}
	}()
}

func updateServiceStatus() {
	setToHealthy, setToIdle := getServicesToChangeStatus()

	for _, service := range setToHealthy {
		err := controller.Service.UpdateServiceStatus(service, core.ServiceStatusHealthy)
		if err != nil {
			log.Logger.Error(err)
		}
	}
	for _, service := range setToIdle {
		err := controller.Service.UpdateServiceStatus(service, core.ServiceStatusIdle)
		if err != nil {
			log.Logger.Error(err)
		}
	}
}

func getServicesToChangeStatus() (setToHealthy []string, setToIdle []string) {
	var healthyServices []string
	var idleServices []string

	services := controller.Service.GetAllServiceKeys()
	var wg sync.WaitGroup
	for _, serviceKey := range services {
		wg.Add(1)

		go func(serviceKey string) {
			defer wg.Done()

			service, err := controller.Service.GetService(serviceKey)
			if err != nil {
				log.Logger.Error(err)
				return
			}
			client, err := NewClient(service)
			if err != nil {
				log.Logger.Error(err)
				return
			}

			err = client.HealthCheck()
			if err != nil {
				if service.Status == core.ServiceStatusHealthy {
					log.Logger.Infof("Service %s is now idle. Cause: %v", service.Path, err)
					idleServices = append(idleServices, service.Path)
				}
			} else {
				if service.Status == core.ServiceStatusIdle {
					log.Logger.Infof("Service %s is now healthy", service.Path)
					healthyServices = append(healthyServices, service.Path)
				}
			}

		}(serviceKey)
	}
	wg.Wait()
	return healthyServices, idleServices
}
