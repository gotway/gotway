package cache

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	httpError "github.com/gotway/gotway/internal/http/error"
	"github.com/gotway/gotway/internal/middleware"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/request"
	"github.com/gotway/gotway/pkg/log"
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
		service, err := request.GetService(r)
		if err != nil {
			httpError.Handle(err, w, g.logger)
			return
		}

		backendReq, err := getBackendRequest(r, service)
		if err != nil {
			httpError.Handle(err, w, g.logger)
			return
		}

		res, err := g.client.Do(backendReq)
		if err != nil {
			g.logger.Error("error requesting service ", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		g.log(r, res, backendReq.URL)

		next.ServeHTTP(w, request.WithResponse(r, res))
	})
}

func (g *gateway) log(req *http.Request, res *http.Response, target *url.URL) {
	g.logger.Infof("%s %s => %s %d", req.Method, req.URL, target, res.StatusCode)
}

func getBackendRequest(r *http.Request, service model.Service) (*http.Request, error) {
	url := service.Backend.URL + r.URL.Path
	if r.URL.RawQuery != "" {
		url = fmt.Sprintf("%s?%s", url, r.URL.RawQuery)
	}
	backendReq, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		return nil, err
	}
	backendReq.Header.Add("X-Forwarded-Host", r.Host)
	backendReq.Header.Add("X-Origin-Host", backendReq.Host)
	return backendReq, nil
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
