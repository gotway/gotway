package model

import (
	"net/http"
	"strconv"
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
				log.Debug("Checking service health")
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
	for _, serviceKey := range services {
		service, serviceErr := ServiceDao.GetService(serviceKey)
		if serviceErr != nil {
			log.Error(serviceErr)
			continue
		}
		_, err := http.Get(service.HealthURL)
		if err != nil {
			if service.Status == Healthy {
				idleServices = append(idleServices, service.Key)
			}
		} else {
			if service.Status == Idle {
				healthyServices = append(healthyServices, service.Key)
			}
		}
	}

	return healthyServices, idleServices
}
