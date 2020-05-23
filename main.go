package main

import (
	"github.com/gosmo-devs/microgateway/api"
	"github.com/gosmo-devs/microgateway/log"
)

func main() {
	log.InitLogger()
	api.NewAPI()
}
