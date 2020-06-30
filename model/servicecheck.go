package model

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/log"
)

func initHealthcheck() {
	interval, _ := strconv.Atoi(config.ServiceCheckInterval)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

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
		ServiceDao.updateServiceStatus(service, Healthy)
	}
	for _, service := range setToIdle {
		ServiceDao.updateServiceStatus(service, Idle)
	}
}

func getServicesToChangeStatus() (setToHealthy []string, setToIdle []string) {
	var healthyServices []string
	var idleServices []string

	services := ServiceDao.getAllServices()
	var wg sync.WaitGroup
	for _, serviceKey := range services {
		wg.Add(1)
		go func(serviceKey string) {
			defer wg.Done()
			service, serviceErr := ServiceDao.GetService(serviceKey)
			if serviceErr != nil {
				log.Error(serviceErr)
			} else {
				_, err := http.Get(service.HealthURL)
				if err != nil {
					if service.Status == Healthy {
						log.Infof("Service %s is now idle. Cause: %v", service.Key, err)
						idleServices = append(idleServices, service.Key)
					}
				} else {
					if service.Status == Idle {
						log.Infof("Service %s is now healty", service.Key)
						healthyServices = append(healthyServices, service.Key)
					}
				}
			}
		}(serviceKey)
	}
	wg.Wait()
	return healthyServices, idleServices
}
