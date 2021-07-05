package error

import (
	"errors"
	"net/http"

	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/pkg/log"
)

func Handle(err error, w http.ResponseWriter, logger log.Logger) {
	logger.Error(err)
	badRequestErrors := []error{
		model.ErrServiceAlreadyRegistered,
		model.ErrInvalidCacheConfig,
		model.ErrInvalidDeleteCache,
	}
	notFoundErrors := []error{
		model.ErrServiceNotFound,
		model.ErrCacheNotFound,
		model.ErrCacheConfigNotFound,
	}
	for _, e := range badRequestErrors {
		if errors.Is(err, e) {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
	}
	for _, e := range notFoundErrors {
		if errors.Is(err, e) {
			http.Error(w, e.Error(), http.StatusNotFound)
			return
		}
	}
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
