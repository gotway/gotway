package client

import (
	"context"
	"net/url"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type clientGRPC struct {
	options Options
}

// HealthCheck performs health check
func (c clientGRPC) HealthCheck(url *url.URL) error {
	conn, err := getConn(url.Host, c.options)
	if err != nil {
		return err
	}
	defer conn.Close()

	healthClient := healthpb.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), c.options.Timeout)
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

func getConn(server string, options Options) (*grpc.ClientConn, error) {
	return grpc.Dial(server,
		grpc.WithBlock(),
		grpc.WithTimeout(options.Timeout),
		grpc.WithInsecure(),
	)
}

func newClientGRPC(options Options) clientGRPC {
	return clientGRPC{options}
}
