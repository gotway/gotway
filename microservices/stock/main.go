package main

import (
	"github.com/gotway/microsamples/stock/api"
	"github.com/gotway/microsamples/stock/redis"
)

func main() {
	redis.Init()
	api.NewAPI()
}
