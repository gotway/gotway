package http

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gotway/gotway/internal/cache"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/service"
	"github.com/gotway/gotway/pkg/log"
)

type middlewareOptions struct {
	gatewayTimeout time.Duration
}

type middleware struct {
	client            *http.Client
	cacheController   cache.Controller
	serviceController service.Controller
	logger            log.Logger
}

const (
	serviceKey  = "service"
	responseKey = "response"
)

func (m *middleware) matchService(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Debug("matchService")
		services, err := m.serviceController.GetServices()
		if err != nil {
			m.logger.Error("services not found", err)
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}

		for _, s := range services {
			if s.MatchRequest(r) {
				decorated := r.WithContext(context.WithValue(r.Context(), serviceKey, s))
				next.ServeHTTP(w, decorated)
				return
			}
		}

		http.Error(w, "Service not found", http.StatusNotFound)
	})
}

func (m *middleware) cacheIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Debug("cacheIn")
		service, err := getServiceFromRequest(r)
		if err != nil {
			m.handleInternalError(err, w)
			return
		}

		if !m.cacheController.IsCacheableRequest(r, service) {
			next.ServeHTTP(w, r)
			return
		}

		m.logger.Debug("checking cache")
		cache, err := m.cacheController.GetCache(r, service)
		if err != nil {
			if !errors.Is(err, model.ErrCacheNotFound) {
				m.logger.Error(err)
			}
			next.ServeHTTP(w, r)
			return
		}

		m.logger.Debug("cached response")
		for key, header := range cache.Headers {
			w.Header().Set(key, strings.Join(header[:], ","))
		}
		w.WriteHeader(cache.StatusCode)
		w.Write(cache.Body)
	})
}

func (m *middleware) gateway(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Debug("gateway")
		service, err := getServiceFromRequest(r)
		if err != nil {
			m.handleInternalError(err, w)
			return
		}

		backendReq, err := getBackendRequest(r, service)
		if err != nil {
			m.handleInternalError(err, w)
			return
		}

		res, err := m.client.Do(backendReq)
		if err != nil {
			m.logger.Error("error requesting service ", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		m.log(r, res, backendReq.URL)

		decorated := r.WithContext(context.WithValue(r.Context(), responseKey, res))
		next.ServeHTTP(w, decorated)
	})
}

func (m *middleware) cacheOut(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Debug("cacheOut")
		res, err := getResponseFromRequest(r)
		if err != nil {
			m.handleInternalError(err, w)
			return
		}
		service, err := getServiceFromRequest(r)
		if err != nil {
			m.handleInternalError(err, w)
			return
		}

		if err := m.cacheController.HandleResponse(res, service); err != nil {
			m.handleInternalError(err, w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *middleware) writeResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := getResponseFromRequest(r)
		if err != nil {
			m.handleInternalError(err, w)
			return
		}

		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			m.handleInternalError(err, w)
			return
		}

		m.logger.Debug("writeResponse")
		for key, header := range res.Header {
			w.Header().Set(key, strings.Join(header[:], ","))
		}
		w.WriteHeader(res.StatusCode)
		w.Write(bytes)
	})
}

func (m *middleware) handleInternalError(err error, w http.ResponseWriter) {
	m.logger.Error(err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func (p *middleware) log(req *http.Request, res *http.Response, target *url.URL) {
	p.logger.Infof("%s %s => %s %d", req.Method, req.URL, target, res.StatusCode)
}

func getServiceFromRequest(r *http.Request) (model.Service, error) {
	service, ok := r.Context().Value(serviceKey).(model.Service)
	if !ok {
		return model.Service{}, errors.New("service not found in request context")
	}
	return service, nil
}

func getResponseFromRequest(r *http.Request) (*http.Response, error) {
	res, ok := r.Context().Value(responseKey).(*http.Response)
	if !ok {
		return nil, errors.New("response not found in request context")
	}
	return res, nil
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

func newMiddleware(
	options middlewareOptions,
	serviceController service.Controller,
	cacheController cache.Controller,
	logger log.Logger,
) *middleware {

	return &middleware{
		client:            &http.Client{Timeout: options.gatewayTimeout},
		serviceController: serviceController,
		cacheController:   cacheController,
		logger:            logger,
	}
}
