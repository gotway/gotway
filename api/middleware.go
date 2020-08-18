package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gosmo-devs/microgateway/controller"
	"github.com/gosmo-devs/microgateway/core"
	"github.com/gosmo-devs/microgateway/log"
)

func cacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !controller.Cache.IsCacheableRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		log.Logger.Debug("Checking cache")
		serviceKey := getServiceKey(r)
		cache, err := controller.Cache.GetCache(r, serviceKey)
		if err != nil {
			if !errors.Is(err, core.ErrCacheNotFound) {
				log.Logger.Error(err)
			}
			next.ServeHTTP(w, r)
			return
		}

		log.Logger.Debug("Cached response")
		bodyBytes, err := ioutil.ReadAll(cache.Body)
		if err != nil {
			log.Logger.Error(err)
			next.ServeHTTP(w, r)
			return
		}
		for key, header := range cache.Headers {
			w.Header().Set(key, strings.Join(header[:], ","))
		}
		w.WriteHeader(cache.StatusCode)
		w.Write(bodyBytes)
	})
}
