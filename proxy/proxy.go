package proxy

import (
	"net/http"
	"net/url"

	"github.com/gosmo-devs/microgateway/log"
	"github.com/gosmo-devs/microgateway/model"
)

// Proxy interface
type Proxy interface {
	getTargetURL(r *http.Request) (*url.URL, error)
	ReverseProxy(w http.ResponseWriter, r *http.Request) error
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
	log.Infof("%s %s => %s %d", req.Method, req.URL, target, res.StatusCode)
}

func handleError(w http.ResponseWriter, err error) {
	log.Error(err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

// NewProxy instanciates a new Proxy
func NewProxy(service *model.Service) (Proxy, error) {
	switch service.Type {
	case model.REST:
		return restProxy{service}, nil
	case model.GRPC:
		return grpcProxy{service}, nil
	default:
		return nil, model.ErrInvalidServiceType
	}
}
