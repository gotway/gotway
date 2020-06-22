package main

import (
	"github.com/gosmo-devs/microsamples/stock/api"
	"github.com/gosmo-devs/microsamples/stock/redis"
)

func main() {
	redis.Init()
	api.NewAPI()
}
