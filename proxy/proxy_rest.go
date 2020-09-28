package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gosmo-devs/microgateway/util"
)

type proxyREST struct {
	proxy
}

func (p proxyREST) getTargetURL(r *http.Request) (*url.URL, error) {
	path := util.GetServiceRelativePath(r, p.service.Path)
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
