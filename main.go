package main

import (
	"github.com/gosmo-devs/microgateway/api"
	"github.com/gosmo-devs/microgateway/log"
	"github.com/gosmo-devs/microgateway/model"
)

func main() {
	log.InitLogger()
	model.Init()
	api.NewAPI()
}
