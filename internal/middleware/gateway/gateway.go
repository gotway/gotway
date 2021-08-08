package cache

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	httpError "github.com/gotway/gotway/internal/http/error"
	"github.com/gotway/gotway/internal/middleware"
	"github.com/gotway/gotway/internal/requestcontext"
	"github.com/gotway/gotway/pkg/log"

	crdv1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1"
)

type GatewayOptions struct {
	Timeout time.Duration
}

type gateway struct {
	client *http.Client
	logger log.Logger
}

func (g *gateway) MiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g.logger.Debug("gateway")
		ingress, err := requestcontext.GetIngress(r)
		if err != nil {
			httpError.Handle(err, w, g.logger)
			return
		}

		if !ingress.Status.IsServiceHealthy {
			http.Error(w, "service not available", http.StatusServiceUnavailable)
			return
		}

		serviceReq, err := getServiceRequest(r, ingress)
		if err != nil {
			httpError.Handle(err, w, g.logger)
			return
		}

		res, err := g.client.Do(serviceReq)
		if err != nil {
			g.logger.Error("error requesting service ", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		g.log(r, res, serviceReq.URL)

		next.ServeHTTP(w, requestcontext.WithResponse(r, res))
	})
}

func (g *gateway) log(req *http.Request, res *http.Response, target *url.URL) {
	g.logger.Infof("%s %s => %s %d", req.Method, req.URL, target, res.StatusCode)
}

func getServiceRequest(r *http.Request, ingress crdv1alpha1.IngressHTTP) (*http.Request, error) {
	url := ingress.Spec.Service.URL + r.URL.Path
	if r.URL.RawQuery != "" {
		url = fmt.Sprintf("%s?%s", url, r.URL.RawQuery)
	}
	serviceReq, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		return nil, err
	}
	serviceReq.Header.Add("X-Forwarded-Host", r.Host)
	serviceReq.Header.Add("X-Origin-Host", serviceReq.Host)
	return serviceReq, nil
}

func New(
	options GatewayOptions,
	logger log.Logger,
) middleware.Middleware {

	return &gateway{
		client: &http.Client{Timeout: options.Timeout},
		logger: logger,
	}
}
