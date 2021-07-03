package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/gotway/gotway/internal/model"
	"golang.org/x/net/http2"
)

type ClientOptions struct {
	Timeout time.Duration
	Type    model.ServiceType
}

type Client struct {
	httpClient http.Client
}

func (c *Client) Request(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func getTransport(options ClientOptions) http.RoundTripper {
	if options.Type == model.ServiceTypeGRPC {
		return &http2.Transport{
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		}
	}
	return &http.Transport{}
}

func New(options ClientOptions) *Client {
	return &Client{http.Client{
		Timeout:   options.Timeout,
		Transport: getTransport(options),
	}}
}
