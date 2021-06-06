package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gotway/gotway/internal/core"
)

type proxyREST struct {
	proxy
}

func (p proxyREST) getTargetURL(r *http.Request) (*url.URL, error) {
	path, err := core.GetServiceRelativePath(r, p.service.Path)
	if err != nil {
		return nil, err
	}
	target, err := url.Parse(p.service.URL + path)
	if err != nil {
		return nil, err
	}
	return target, nil
}

// ReverseProxy forwards traffic to a service
func (p proxyREST) ReverseProxy(w http.ResponseWriter, r *http.Request) error {
	target, err := p.getTargetURL(r)
	if err != nil {
		return err
	}

	proxy := &httputil.ReverseProxy{
		Director: getDirector(target),
		ModifyResponse: func(res *http.Response) error {
			p.log(r, res, target)
			return p.handleResponse(p.service.Path, res)
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			p.handleError(w, err)
		},
	}
	proxy.ServeHTTP(w, r)

	return nil
}
