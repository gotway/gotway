package client

import (
	"context"
	"net/url"
	"sync"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type clientGRPC struct {
	clients map[string]healthpb.HealthClient
	conns   []*grpc.ClientConn
	mux     sync.Mutex
	options Options
}

// HealthCheck performs health check
func (c *clientGRPC) HealthCheck(url *url.URL) error {
	client, err := c.getClient(url.Host, c.options)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.options.Timeout)
	defer cancel()
	health, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		return err
	}
	if health.Status != healthpb.HealthCheckResponse_SERVING {
		return ErrServiceNotAvailable
	}
	return nil
}

// Release connections
func (c *clientGRPC) Release() {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, c := range c.conns {
		c.Close()
	}
}

func (c *clientGRPC) getClient(server string, options Options) (healthpb.HealthClient, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if client, ok := c.clients[server]; ok {
		return client, nil
	}

	conn, err := grpc.Dial(server,
		grpc.WithBlock(),
		grpc.WithTimeout(options.Timeout),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	healthClient := healthpb.NewHealthClient(conn)

	c.clients[server] = healthClient
	c.conns = append(c.conns, conn)

	return healthClient, nil
}

func newClientGRPC(options Options) *clientGRPC {
	return &clientGRPC{
		clients: make(map[string]healthpb.HealthClient),
		options: options,
	}
}
