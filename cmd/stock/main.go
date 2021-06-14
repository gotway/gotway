package main

import (
	"github.com/gotway/gotway/cmd/stock/api"
	"github.com/gotway/gotway/cmd/stock/redis"
)

func main() {
	redis.Init()
	api.NewAPI()
}
