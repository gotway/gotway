package http

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gotway/gotway/internal/cache"
	"github.com/gotway/gotway/internal/config"
	"github.com/gotway/gotway/internal/http/client"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/service"
	"github.com/gotway/gotway/pkg/log"
)

type middleware struct {
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
		service, err := getServiceFromReq(r)
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
		service, err := getServiceFromReq(r)
		if err != nil {
			m.handleInternalError(err, w)
			return
		}

		backendReq, err := getRequestForBackend(r, service)
		if err != nil {
			m.handleInternalError(err, w)
			return
		}
		c := client.New(client.ClientOptions{
			Timeout: config.GatewayTimeout,
			Type:    service.Type,
		})
		res, err := c.Request(backendReq)
		if err != nil {
			m.logger.Error("error requesting service ", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		decorated := r.WithContext(context.WithValue(r.Context(), responseKey, res))
		next.ServeHTTP(w, decorated)
	})
}

func (m *middleware) cacheOut(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Debug("cacheOut")
		res, err := getResponseFromReq(r)
		if err != nil {
			m.handleInternalError(err, w)
			return
		}
		service, err := getServiceFromReq(r)
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
		res, err := getResponseFromReq(r)
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

func getServiceFromReq(r *http.Request) (model.Service, error) {
	service, ok := r.Context().Value(serviceKey).(model.Service)
	if !ok {
		return model.Service{}, errors.New("service not found in request context")
	}
	return service, nil
}

func getResponseFromReq(r *http.Request) (*http.Response, error) {
	res, ok := r.Context().Value(responseKey).(*http.Response)
	if !ok {
		return nil, errors.New("response not found in request context")
	}
	return res, nil
}

func getRequestForBackend(incomingReq *http.Request, service model.Service) (*http.Request, error) {
	url, err := url.Parse(service.Backend.URL)
	if err != nil {
		return nil, err
	}

	url.Path = incomingReq.URL.Path
	if service.Type == model.ServiceTypeGRPC {
		url.Scheme = "https"
	}

	return http.NewRequest(incomingReq.Method, url.String(), incomingReq.Body)
}
