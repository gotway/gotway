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
)

// NewAPI Starts a new HTTP server
func NewAPI() {
	router := createRouting()
	addr := fmt.Sprintf(":%s", config.Port)
	log.Info("Server listening on port ", config.Port)
	log.Info("Environment: ", config.Env)
	log.Fatal(http.ListenAndServe(addr, router))
}

func createRouting() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/register", registerServiceHandler).Methods("POST")
	return router
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func registerServiceHandler(w http.ResponseWriter, r *http.Request) {
	decoded := json.NewDecoder(r.Body)

	var json Register
	err := decoded.Decode(&json)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if json.Key == nil {
		http.Error(w, "Missing field 'Key' from JSON", http.StatusBadRequest)
		return
	}

	if json.URL == nil {
		http.Error(w, "Missing field 'Url' from JSON", http.StatusBadRequest)
		return
	}

	if json.HealthURL == nil {
		http.Error(w, "Missing field 'HealthEndpoint' from JSON", http.StatusBadRequest)
		return
	}

	err = controller.RegisterService(*json.Key, *json.URL, *json.HealthURL)
	if err != nil {
		if errors.Is(err, controller.ErrAlreadyRegistered) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
