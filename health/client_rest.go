package health

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gotway/gotway/config"
	"github.com/gotway/gotway/core"
)

type clientREST struct {
	service core.Service
}

func (c clientREST) getHealthURL() (*url.URL, error) {
	healthPath, err := c.service.HealthPathForType()
	if err != nil {
		return nil, err
	}
	urlString := fmt.Sprintf("%s/%s", c.service.URL, healthPath)
	url, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	return url, nil
}

// HealthCheck performs health check
func (c clientREST) HealthCheck() error {
	healthURL, err := c.getHealthURL()
	if err != nil {
		return err
	}

	client := http.Client{Timeout: config.HealthCheckTimeout}

	res, err := client.Get(healthURL.String())
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errServiceNotAvailable
	}
	return nil
}
