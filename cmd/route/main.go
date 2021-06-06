package main

import (
	"fmt"
	"log"
	"net"

	"github.com/gotway/gotway/cmd/route/config"
	"github.com/gotway/gotway/cmd/route/server"
)

func main() {
	server := server.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	} else {
		log.Printf("Server listening on port %s", config.Port)
	}
	server.Serve(lis)
}
