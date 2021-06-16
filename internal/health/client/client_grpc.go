package client

import (
	"context"
	"net/url"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/gotway/gotway/internal/config"
)

type clientGRPC struct{}

// HealthCheck performs health check
func (c clientGRPC) HealthCheck(url *url.URL) error {
	conn, err := getConn(url.Host)
	if err != nil {
		return err
	}

	healthClient := healthpb.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), config.HealthCheckTimeout)
	defer cancel()

	health, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		return err
	}
	if health.Status != healthpb.HealthCheckResponse_SERVING {
		return ErrServiceNotAvailable
	}
	return nil
}

func getConn(server string) (*grpc.ClientConn, error) {
	return grpc.Dial(server,
		grpc.WithBlock(),
		grpc.WithTimeout(config.HealthCheckTimeout),
		grpc.WithInsecure(),
	)
}

func newClientGRPC() clientGRPC {
	return clientGRPC{}
}
