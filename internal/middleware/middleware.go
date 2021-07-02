package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gotway/gotway/internal/cache"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/service"
	"github.com/gotway/gotway/pkg/log"
)

type Middleware struct {
	cacheController   cache.Controller
	serviceController service.Controller
	logger            log.Logger
}

func (m *Middleware) MatchService(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services, err := m.serviceController.GetServices()
		if err != nil {
			m.logger.Error("services not found", err)
			return
		}

		for _, s := range services {
			if s.MatchRequest(r) {
				decorated := r.WithContext(context.WithValue(r.Context(), "service", s))
				next.ServeHTTP(w, decorated)
				return
			}
		}

		http.Error(w, "Service not found", http.StatusNotFound)
	})
}

func (m *Middleware) Cache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service, ok := r.Context().Value("service").(model.Service)
		if !ok {
			m.logger.Error("service not found in request context")
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

func New(
	cacheController cache.Controller,
	serviceController service.Controller,
	logger log.Logger,
) *Middleware {
	return &Middleware{cacheController, serviceController, logger}
}
