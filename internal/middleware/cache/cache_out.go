package cache

import (
	"net/http"

	"github.com/gotway/gotway/internal/cache"
	httpError "github.com/gotway/gotway/internal/http/error"
	"github.com/gotway/gotway/internal/middleware"
	"github.com/gotway/gotway/internal/requestcontext"
	"github.com/gotway/gotway/pkg/log"
)

type cacheOut struct {
	cacheCtrl cache.Controller
	logger    log.Logger
}

func (c *cacheOut) MiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.logger.Debug("cache out")
		res, err := requestcontext.GetResponse(r)
		if err != nil {
			httpError.Handle(err, w, c.logger)
			return
		}
		ingress, err := requestcontext.GetIngress(r)
		if err != nil {
			httpError.Handle(err, w, c.logger)
			return
		}

		params := cache.Params{
			Service:  ingress.Spec.Service.Name,
			TTL:      ingress.Spec.Cache.TTL,
			Statuses: ingress.Spec.Cache.Statuses,
			Tags:     ingress.Spec.Cache.Tags,
		}
		if err := c.cacheCtrl.HandleResponse(res, params); err != nil {
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
