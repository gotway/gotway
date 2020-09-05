package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mmontes11/go-grpc-routes/config"
	"github.com/mmontes11/go-grpc-routes/server"
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
