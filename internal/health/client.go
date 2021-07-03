package health

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

type clientOptions struct {
	timeout time.Duration
}

type client struct {
	client http.Client
}

func (c client) healthCheck(url *url.URL) error {
	res, err := c.client.Get(url.String())
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("service not available")
	}
	return nil
}

func newClient(options clientOptions) client {
	return client{http.Client{Timeout: options.timeout}}
}
