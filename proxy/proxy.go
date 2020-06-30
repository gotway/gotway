package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gosmo-devs/microgateway/log"
	"github.com/gosmo-devs/microgateway/model"
)

// Proxy object
type Proxy struct {
	service *model.Service
}

// ReverseProxy forwards traffic to a service
func (p *Proxy) ReverseProxy(w http.ResponseWriter, r *http.Request) error {
	target, err := getTargetURL(r, p.service)
	if err != nil {
		return err
	}

	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.Header.Add("X-Forwarded-Host", r.Host)
			r.Header.Add("X-Origin-Host", target.Host)
			r.URL.Scheme = target.Scheme
			r.URL.Host = target.Host
			r.URL.Path = target.Path
		},
		ModifyResponse: func(res *http.Response) error {
			log.Infof("%s %s => %s %d", r.Method, r.URL, target, res.StatusCode)
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Error(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		},
	}
	proxy.ServeHTTP(w, r)

	return nil
}

// NewProxy returns a new Proxy
func NewProxy(service *model.Service) *Proxy {
	return &Proxy{service}
}

func getTargetURL(r *http.Request, service *model.Service) (*url.URL, error) {
	path := strings.Split(r.URL.String(), service.Key)[1]
	target, err := url.Parse(service.URL + path)
	if err != nil {
		return nil, err
	}
	return target, nil
}
