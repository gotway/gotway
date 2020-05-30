package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/controller"
	"github.com/gosmo-devs/microgateway/log"
)

// NewAPI Starts a new HTTP server
func NewAPI() {
	router := createRouting()
	addr := fmt.Sprintf(":%s", config.PORT)
	log.Info("Server listening on port ", config.PORT)
	log.Info("Environment: ", config.ENV)
	log.Fatal(http.ListenAndServe(addr, router))
}

func createRouting() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/salute", saluteHandler).Methods("GET")
	router.HandleFunc("/api/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/register", registerHandler).Methods("POST")
	return router
}

func saluteHandler(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(2)%2 == 0 {
		apiError := apiError{errors.New("not found"), "No salute found for this request", 404, r.URL.String()}
		log.Error(apiError.Formatted())
		http.Error(w, "No salute for you", http.StatusNotFound)
	} else {
		fmt.Fprintf(w, "Hello\n")
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	decoded := json.NewDecoder(r.Body)

	var json Register
	err := decoded.Decode(&json)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if json.Key == nil {
		http.Error(w, "Missing field 'Key' from JSON", http.StatusBadRequest)
	}

	if json.Url == nil {
		http.Error(w, "Missing field 'Url' from JSON", http.StatusBadRequest)
	}

	if json.HealthEndpoint == nil {
		http.Error(w, "Missing field 'HealthEndpoint' from JSON", http.StatusBadRequest)
	}

	// TODO not sure if this fits here or in the controller
	if json.TTL == nil {
		*json.TTL, err = strconv.Atoi(config.DEFAULT_SERVICE_TTL)
		if err != nil {
			http.Error(w, "Invalid format for field 'TTL'", http.StatusBadRequest)
		}
	}

	registered := controller.RegisterEndpoint(*json.Key, *json.Url, *json.HealthEndpoint, *json.TTL)
	if registered == false {
		http.Error(w, "Error registering service", http.StatusInternalServerError)
	}

	// TODO improve the errors given to the client (e.g. "Service already registered", "Internal Server Error"...)
}

// curl -X POST -d "{\"key\":\"stock\",\"url\":\"http://stock.microservice.com\",\"healthEndpoint\":\"http://stock.microsrevice.com/api/health\",\"ttl\":1}" http://localhost:8000/api/register
