package proxy

import (
	"net/http"
	"net/url"

	"github.com/gotway/gotway/core"
	"github.com/gotway/gotway/log"
)

// Proxy interface
type Proxy interface {
	getTargetURL(r *http.Request) (*url.URL, error)
	ReverseProxy(w http.ResponseWriter, r *http.Request) error
}

// ResponseHandler is a function hook for handling responses
type ResponseHandler = func(serviceKey string, res *http.Response) error

type proxy struct {
	service        core.Service
	handleResponse ResponseHandler
}

// NewProxy instanciates a new Proxy
func NewProxy(service core.Service, handleResponse ResponseHandler) (Proxy, error) {
	switch service.Type {
	case core.ServiceTypeREST:
		return proxyREST{
			proxy: proxy{
				service,
				handleResponse,
			},
		}, nil
	case core.ServiceTypeGRPC:
		return proxyGRPC{
			proxy: proxy{
				service,
				handleResponse,
			},
		}, nil
	default:
		return nil, core.ErrInvalidServiceType
	}
}

func getDirector(target *url.URL) func(r *http.Request) {
	return func(r *http.Request) {
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Add("X-Origin-Host", target.Host)
		r.URL.Scheme = target.Scheme
		r.URL.Host = target.Host
		r.URL.Path = target.Path
	}
}

func logProxy(req *http.Request, res *http.Response, target *url.URL) {
	log.Logger.Infof("%s %s => %s %d", req.Method, req.URL, target, res.StatusCode)
}

func handleError(w http.ResponseWriter, err error) {
	log.Logger.Error(err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
