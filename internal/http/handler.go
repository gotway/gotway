package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gotway/gotway/internal/model"
)

func (s *Server) getServicesHandler(w http.ResponseWriter, r *http.Request) {
	services, err := s.serviceController.GetServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (s *Server) createService(w http.ResponseWriter, r *http.Request) {
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

	err = s.serviceController.CreateService(service)
	if err != nil {
		s.logger.Error("error creating service", err)
		if errors.Is(err, model.ErrServiceAlreadyRegistered) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getService(w http.ResponseWriter, r *http.Request) {
	key := getServiceKey(r)
	serviceDetail, err := s.serviceController.GetService(key)
	if err != nil {
		s.handleServiceError(err, w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serviceDetail)
}

func (s *Server) deleteService(w http.ResponseWriter, r *http.Request) {
	key := getServiceKey(r)
	err := s.serviceController.DeleteService(key)
	if err != nil {
		s.handleServiceError(err, w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteCache(w http.ResponseWriter, r *http.Request) {
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
		err := s.cacheController.DeleteCacheByPath(payload.Paths)
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
		err := s.cacheController.DeleteCacheByTags(payload.Tags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) proxy(w http.ResponseWriter, r *http.Request) {
	service, err := getRequestService(r)
	if !service.IsHealthy() {
		http.Error(
			w,
			fmt.Sprintf("'%s' service is not responding", service.ID),
			http.StatusBadGateway,
		)
		return
	}

	err = s.serviceController.ReverseProxy(w, r, service, s.cacheController.HandleResponse)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleServiceError(err error, w http.ResponseWriter, r *http.Request) {
	s.logger.Error(err)
	if errors.Is(err, model.ErrServiceNotFound) {
		http.Error(w, fmt.Sprintf("service not found"), http.StatusNotFound)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func getServiceKey(r *http.Request) string {
	params := mux.Vars(r)
	return params["service"]
}
