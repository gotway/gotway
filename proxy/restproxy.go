package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gosmo-devs/microgateway/model"
)

type restProxy struct {
	service *model.Service
}

func (p restProxy) getTargetURL(r *http.Request) (*url.URL, error) {
	path := strings.Split(r.URL.String(), p.service.Path)[1]
	target, err := url.Parse(p.service.URL + path)
	if err != nil {
		return nil, err
	}
	return target, nil
}

// ReverseProxy forwards traffic to a service
func (p restProxy) ReverseProxy(w http.ResponseWriter, r *http.Request) error {
	target, err := p.getTargetURL(r)
	if err != nil {
		return err
	}

	proxy := &httputil.ReverseProxy{
		Director: getDirector(target),
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
