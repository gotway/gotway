package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/controller"
	"github.com/gosmo-devs/microgateway/log"
	"github.com/gosmo-devs/microgateway/model"
)

// NewAPI Starts a new HTTP server
func NewAPI() {
	router := createRouting()
	addr := fmt.Sprintf(":%s", config.Port)
	log.Info("Server listening on port ", config.Port)
	log.Info("Environment: ", config.Env)
	var err error
	if config.TLS {
		err = http.ListenAndServeTLS(addr, config.TLScert, config.TLSkey, router)
	} else {
		err = http.ListenAndServe(addr, router)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func createRouting() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/register", registerServiceHandler).Methods("POST")
	router.PathPrefix("/{service}").HandlerFunc(proxyHandler)
	return router
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["service"]

	service, err := controller.GetService(key)
	if err != nil {
		if errors.Is(err, model.ErrServiceNotFound) {
			http.Error(w, fmt.Sprintf("'%s' service not found", key), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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
