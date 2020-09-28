package health

import (
	"context"
	"net/url"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/core"
)

type clientGRPC struct {
	service core.Service
}

func getConn(server string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(config.HealthCheckTimeout),
		grpc.WithInsecure(),
	}
	return grpc.Dial(server, opts...)
}

func (c clientGRPC) getHealthURL() (*url.URL, error) {
	url, err := url.Parse(c.service.URL)
	if err != nil {
		return nil, err
	}
	return url, nil
}

// HealthCheck performs health check
func (c clientGRPC) HealthCheck() error {
	healthURL, err := c.getHealthURL()
	if err != nil {
		return err
	}
	conn, err := getConn(healthURL.Host)
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
		return errServiceNotAvailable
	}
	return nil
}
