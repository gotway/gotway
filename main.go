package main

import (
	"github.com/gosmo-devs/microgateway/api"
	"github.com/gosmo-devs/microgateway/log"
	"github.com/gosmo-devs/microgateway/redis"
)

func main() {
	log.InitLogger()
	redis.Init()
	api.NewAPI()
}
