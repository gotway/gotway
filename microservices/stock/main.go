package main

import (
	"github.com/gotway/gotway/microservices/stock/api"
	"github.com/gotway/gotway/microservices/stock/redis"
)

func main() {
	redis.Init()
	api.NewAPI()
}
