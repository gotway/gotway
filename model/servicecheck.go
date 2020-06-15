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
	for _, service := range services {
		status, healthURL := ServiceDao.getStatusAndHealthURL(service)
		_, err := http.Get(healthURL)
		if err != nil {
			if status == Healthy.String() {
				idleServices = append(idleServices, service)
			}
		} else {
			if status == Idle.String() {
				healthyServices = append(healthyServices, service)
			}
		}
	}

	return healthyServices, idleServices
}
