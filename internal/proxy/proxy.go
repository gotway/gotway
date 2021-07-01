package proxy

import (
	"net/http"
	"net/url"

	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/pkg/log"
)

// Proxy interface
type Proxy interface {
	ReverseProxy(w http.ResponseWriter, r *http.Request) error
}

// ResponseHandler is a function hook for handling responses
type ResponseHandler = func(r *http.Response, service model.Service) error

type proxy struct {
	service        model.Service
	handleResponse ResponseHandler
	logger         log.Logger
}

func (p *proxy) log(req *http.Request, res *http.Response, target *url.URL) {
	p.logger.Infof("%s %s => %s %d", req.Method, req.URL, target, res.StatusCode)
}

func (p *proxy) handleError(w http.ResponseWriter, err error) {
	p.logger.Error("proxy error ", err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
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

// New instanciates a new Proxy
func New(service model.Service, handleResponse ResponseHandler, logger log.Logger) (Proxy, error) {
	proxy := proxy{
		service,
		handleResponse,
		logger,
	}
	switch service.Type {
	case model.ServiceTypeREST:
		return proxyREST{proxy}, nil
	case model.ServiceTypeGRPC:
		return proxyGRPC{proxy}, nil
	default:
		return nil, model.ErrInvalidServiceType
	}
}
