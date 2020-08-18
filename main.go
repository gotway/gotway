package main

import (
	"github.com/gosmo-devs/microgateway/api"
	"github.com/gosmo-devs/microgateway/controller"
	"github.com/gosmo-devs/microgateway/health"
	"github.com/gosmo-devs/microgateway/log"
	"github.com/gosmo-devs/microgateway/model"
)

func main() {
	log.Init()
	model.Init()
	controller.Init()
	health.Init()
	api.NewAPI()
}
