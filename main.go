package main

import (
	"github.com/gotway/gotway/api"
	"github.com/gotway/gotway/controller"
	"github.com/gotway/gotway/health"
	"github.com/gotway/gotway/log"
	"github.com/gotway/gotway/model"
)

func main() {
	log.Init()
	model.Init()
	controller.Init()
	health.Init()
	api.NewAPI()
}
