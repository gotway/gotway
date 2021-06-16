package client

import (
	"net/http"
	"net/url"

	"github.com/gotway/gotway/internal/config"
)

type clientREST struct {
	client http.Client
}

// HealthCheck performs health check
func (c clientREST) HealthCheck(url *url.URL) error {
	res, err := c.client.Get(url.String())
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return ErrServiceNotAvailable
	}
	return nil
}

func newClientREST() clientREST {
	return clientREST{http.Client{Timeout: config.HealthCheckTimeout}}
}
