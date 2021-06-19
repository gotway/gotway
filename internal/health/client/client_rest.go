package client

import (
	"net/http"
	"net/url"
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

func newClientREST(options Options) clientREST {
	return clientREST{http.Client{Timeout: options.Timeout}}
}
