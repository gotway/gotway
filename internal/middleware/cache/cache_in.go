package cache

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gotway/gotway/internal/cache"
	httpError "github.com/gotway/gotway/internal/http/error"
	"github.com/gotway/gotway/internal/middleware"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/requestcontext"
	"github.com/gotway/gotway/pkg/log"
)

type cacheIn struct {
	cacheCtrl cache.Controller
	logger    log.Logger
}

func (c *cacheIn) MiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.logger.Debug("cache in")
		ingress, err := requestcontext.GetIngress(r)
		if err != nil {
			httpError.Handle(err, w, c.logger)
			return
		}

		if !c.cacheCtrl.IsCacheableRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		c.logger.Debug("checking cache")
		cache, err := c.cacheCtrl.GetCache(r, ingress.Spec.Service.Name)
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
		_, _ = w.Write(cache.Body)
	})
}

func NewCacheIn(
	cacheController cache.Controller,
	logger log.Logger,
) middleware.Middleware {
	return &cacheIn{cacheController, logger}
}
