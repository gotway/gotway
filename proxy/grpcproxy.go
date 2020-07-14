package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gosmo-devs/microgateway/model"
	"golang.org/x/net/http2"
)

type grpcProxy struct {
	service *model.Service
}

func (p grpcProxy) getTargetURL(r *http.Request) (*url.URL, error) {
	url, err := url.Parse(p.service.URL)
	if err != nil {
		return nil, err
	}
	url.Scheme = "https"
	url.Path = r.URL.Path
	return url, nil
}

func (p grpcProxy) ReverseProxy(w http.ResponseWriter, r *http.Request) error {
	target, err := p.getTargetURL(r)
	if err != nil {
		return err
	}

	proxy := &httputil.ReverseProxy{
		Director: getDirector(target),
		Transport: &http2.Transport{
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
		ModifyResponse: func(res *http.Response) error {
			logProxy(r, res, target)
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			handleError(w, err)
		},
	}
	proxy.ServeHTTP(w, r)

	return nil
}
