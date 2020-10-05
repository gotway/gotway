package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	pb "github.com/gotway/gotway/microservices/route/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/testdata"
)

var (
	server         = flag.String("server", "localhost:11000", "The server address in the format of host:port")
	timeoutSeconds = flag.Int("timeout", 10, "Request timeout in seconds")
	timeout        = time.Duration(*timeoutSeconds) * time.Second
	tls            = flag.Bool("tls", false, "Enable TLS")
	tlsCa          = flag.String("tls_ca", testdata.Path("ca.pem"), "Certificate authority file for TLS")
	tlsServerHost  = flag.String("tls_server_host", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func health(client healthpb.HealthClient) {
	log.Println("HEALTH")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	health, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		log.Printf("Error checking health: %v", err)
		return
	}
	log.Print(health)
}

func getFeature(client pb.RouteClient, point *pb.Point) {
	log.Println("GET FEATURE")
	log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Printf("Error getting feature: %v", err)
		return
	}
	log.Print(feature)
}

func listFeatures(client pb.RouteClient, rect *pb.Rectangle) {
	log.Println("LIST FEATURES")
	log.Printf("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	stream, err := client.ListFeatures(ctx, rect)
	if err != nil {
		log.Printf("Error listing features: %v", err)
		return
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Printf("Error receiving features: %v", err)
			break
		}
		log.Print(feature)
	}
}

func recordRoute(client pb.RouteClient) {
	log.Println("RECORD ROUTE")
	points := randomPoints()
	log.Printf("Transversing %d points", len(points))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	stream, err := client.RecordRoute(ctx)
	if err != nil {
		log.Printf("Error recording route: %v", err)
		return
	}
	for _, point := range points {
		if err := stream.Send(point); err != nil {
			log.Printf("Error sending points: %v", err)
			return
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("Error receiving reply: %v", err)
		return
	}
	log.Printf("Route summary: %v", reply)
}

func routeChat(client pb.RouteClient) {
	log.Println("ROUTE CHAT")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	stream, err := client.RouteChat(ctx)
	if err != nil {
		log.Printf("Error route chatting: %v", err)
		return
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Printf("Failed to receive a note: %v", err)
			}
			log.Printf("Got message \"%s\" at point (%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	for _, note := range notes() {
		if err := stream.Send(note); err != nil {
			log.Printf("Failed to send a note: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func getConn(server string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(timeout),
	}
	if *tls {
		creds, err := credentials.NewClientTLSFromFile(*tlsCa, *tlsServerHost)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	return grpc.Dial(server, opts...)
}

func main() {
	flag.Parse()

	conn, err := getConn(*server)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRouteClient(conn)
	healthClient := healthpb.NewHealthClient(conn)

	health(healthClient)
	getFeature(client, validPoint)
	getFeature(client, invalidPoint)
	listFeatures(client, rect)
	recordRoute(client)
	routeChat(client)
}
