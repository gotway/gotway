package cache

import (
	"net/http"

	"github.com/gotway/gotway/internal/cache"
	httpError "github.com/gotway/gotway/internal/http/error"
	"github.com/gotway/gotway/internal/middleware"
	"github.com/gotway/gotway/internal/request"
	"github.com/gotway/gotway/pkg/log"
)

type cacheOut struct {
	cacheController cache.Controller
	logger          log.Logger
}

func (c *cacheOut) MiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.logger.Debug("cache out")
		res, err := request.GetResponse(r)
		if err != nil {
			httpError.Handle(err, w, c.logger)
			return
		}
		service, err := request.GetService(r)
		if err != nil {
			httpError.Handle(err, w, c.logger)
			return
		}

		if err := c.cacheController.HandleResponse(res, service); err != nil {
			httpError.Handle(err, w, c.logger)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewCacheOut(
	cacheController cache.Controller,
	logger log.Logger,
) middleware.Middleware {

	return &cacheOut{cacheController, logger}
}
