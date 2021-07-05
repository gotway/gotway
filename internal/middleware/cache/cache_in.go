package cache

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gotway/gotway/internal/cache"
	httpError "github.com/gotway/gotway/internal/http/error"
	"github.com/gotway/gotway/internal/middleware"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/request"
	"github.com/gotway/gotway/pkg/log"
)

type cacheIn struct {
	cacheController cache.Controller
	logger          log.Logger
}

func (c *cacheIn) MiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.logger.Debug("cache in")
		service, err := request.GetService(r)
		if err != nil {
			httpError.Handle(err, w, c.logger)
			return
		}

		if !c.cacheController.IsCacheableRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		c.logger.Debug("checking cache")
		cache, err := c.cacheController.GetCache(r, service)
		if err != nil {
			if !errors.Is(err, model.ErrCacheNotFound) {
				c.logger.Error(err)
			}
			next.ServeHTTP(w, r)
			return
		}

		c.logger.Debug("cached response")
		for key, header := range cache.Headers {
			w.Header().Set(key, strings.Join(header[:], ","))
		}
		w.WriteHeader(cache.StatusCode)
		w.Write(cache.Body)
	})
}

func NewCacheIn(
	cacheController cache.Controller,
	logger log.Logger,
) middleware.Middleware {

	return &cacheIn{cacheController, logger}
}
