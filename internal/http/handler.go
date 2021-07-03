package http

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gotway/gotway/internal/cache"
	httpError "github.com/gotway/gotway/internal/http/error"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/request"
	"github.com/gotway/gotway/internal/service"
	"github.com/gotway/gotway/pkg/log"
)

type handler struct {
	serviceController service.Controller
	cacheController   cache.Controller
	logger            log.Logger
}

func (h *handler) getServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.serviceController.GetServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (h *handler) createService(w http.ResponseWriter, r *http.Request) {
	decoded := json.NewDecoder(r.Body)
	var service model.Service
	err := decoded.Decode(&service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = service.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.serviceController.CreateService(service)
	if err != nil {
		h.logger.Error("error creating service", err)
		if errors.Is(err, model.ErrServiceAlreadyRegistered) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) getService(w http.ResponseWriter, r *http.Request) {
	key := getServiceKey(r)
	serviceDetail, err := h.serviceController.GetService(key)
	if err != nil {
		httpError.Handle(err, w, h.logger)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serviceDetail)
}

func (h *handler) deleteService(w http.ResponseWriter, r *http.Request) {
	key := getServiceKey(r)
	err := h.serviceController.DeleteService(key)
	if err != nil {
		httpError.Handle(err, w, h.logger)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) deleteCache(w http.ResponseWriter, r *http.Request) {
	decoded := json.NewDecoder(r.Body)

	var payload model.DeleteCache
	err := decoded.Decode(&payload)
	if err != nil {
		http.Error(w, model.ErrInvalidDeleteCache.Error(), http.StatusBadRequest)
		return
	}

	err = payload.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(payload.Paths) > 0 {
		err := h.cacheController.DeleteCacheByPath(payload.Paths)
		if err != nil {
			if _, ok := err.(*model.ErrCachePathNotFound); ok {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if len(payload.Tags) > 0 {
		err := h.cacheController.DeleteCacheByTags(payload.Tags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) writeResponse(w http.ResponseWriter, r *http.Request) {
	res, err := request.GetResponse(r)
	if err != nil {
		httpError.Handle(err, w, h.logger)
		return
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		httpError.Handle(err, w, h.logger)
		return
	}

	h.logger.Debug("writeResponse")
	for key, header := range res.Header {
		w.Header().Set(key, strings.Join(header[:], ","))
	}
	w.WriteHeader(res.StatusCode)
	w.Write(bytes)
}

func getServiceKey(r *http.Request) string {
	params := mux.Vars(r)
	return params["service"]
}

func newHandler(
	serviceController service.Controller,
	cacheController cache.Controller,
	logger log.Logger,
) *handler {

	return &handler{
		serviceController: serviceController,
		cacheController:   cacheController,
		logger:            logger,
	}
}
