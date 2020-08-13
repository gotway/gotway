package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gosmo-devs/microgateway/controller"
	"github.com/gosmo-devs/microgateway/model"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getServicesHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	offsetStr := q.Get("offset")
	limitStr := q.Get("limit")
	offset, limit, err := processPaginationParams(offsetStr, limitStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	servicePage, err := controller.GetServices(offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servicePage)
}

func registerServiceHandler(w http.ResponseWriter, r *http.Request) {
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

	err = controller.RegisterService(service)
	if err != nil {
		if errors.Is(err, model.ErrServiceAlreadyRegistered) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func getServiceHandler(w http.ResponseWriter, r *http.Request) {
	key := getServiceKey(r)
	service, err := controller.GetService(key)
	if err != nil {
		handleServiceError(err, w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}

func deleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	key := getServiceKey(r)
	err := controller.DeleteService(key)
	if err != nil {
		handleServiceError(err, w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	key := getServiceKey(r)
	service, err := controller.GetService(key)
	if err != nil {
		handleServiceError(err, w, r)
		return
	}

	if !service.IsHealthy() {
		http.Error(w, fmt.Sprintf("'%s' service is not responding", key), http.StatusBadGateway)
		return
	}

	err = controller.ReverseProxy(w, r, service)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func getServiceKey(r *http.Request) string {
	params := mux.Vars(r)
	return params["service"]
}

func handleServiceError(err error, w http.ResponseWriter, r *http.Request) {
	key := getServiceKey(r)
	if errors.Is(err, model.ErrServiceNotFound) {
		http.Error(w, fmt.Sprintf("'%s' service not found", key), http.StatusNotFound)
		return
	}
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func processPaginationParams(offsetStr string, limitStr string) (int, int, error) {
	offset, err := processIntParam(offsetStr, 0)
	if err != nil {
		return 0, 0, err
	}
	limit, err := processIntParam(limitStr, 10)
	if err != nil {
		return 0, 0, err
	}
	if offset > limit {
		return 0, 0, errors.New("Offset cannot be greater than limit")
	}
	return offset, limit, nil
}

func processIntParam(paramStr string, defaultValue int) (int, error) {
	if len(paramStr) == 0 {
		return defaultValue, nil
	}
	param, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, err
	}
	if param < 0 {
		return 0, errors.New("Param cannot not be negative")
	}
	return param, nil
}
