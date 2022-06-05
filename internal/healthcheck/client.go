package healthcheck

import (
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

func (c client) healthCheck(url *url.URL) (bool, error) {
	res, err := c.client.Get(url.String())
	if err != nil {
		return false, nil
	}
	return res.StatusCode == http.StatusOK, nil
}

func newClient(options clientOptions) client {
	return client{http.Client{Timeout: options.timeout}}
}
