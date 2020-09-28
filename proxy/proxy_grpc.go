package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"golang.org/x/net/http2"
)

type proxyGRPC struct {
	proxy
}

func (p proxyGRPC) getTargetURL(r *http.Request) (*url.URL, error) {
	url, err := url.Parse(p.service.URL)
	if err != nil {
		return nil, err
	}
	url.Scheme = "https"
	url.Path = r.URL.Path
	return url, nil
}

func (p proxyGRPC) ReverseProxy(w http.ResponseWriter, r *http.Request) error {
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
			return p.handleResponse(p.service.Path, res)
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			handleError(w, err)
		},
	}
	proxy.ServeHTTP(w, r)

	return nil
}
