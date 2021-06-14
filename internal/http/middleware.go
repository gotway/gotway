package http

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gotway/gotway/internal/core"
)

func (s *Server) cacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.cacheController.IsCacheableRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		s.logger.Debug("Checking cache")
		serviceKey := getServiceKey(r)
		cache, err := s.cacheController.GetCache(r, "", serviceKey)
		if err != nil {
			if !errors.Is(err, core.ErrCacheNotFound) {
				s.logger.Error(err)
			}
			next.ServeHTTP(w, r)
			return
		}

		s.logger.Debug("Cached response")
		bodyBytes, err := ioutil.ReadAll(cache.Body)
		if err != nil {
			s.logger.Error(err)
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
