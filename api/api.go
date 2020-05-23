package api

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/log"
)

// NewAPI Starts a new HTTP server
func NewAPI() {
	router := mux.NewRouter()
	router.HandleFunc("/api/salute", saluteHandler).Methods("GET")
	router.HandleFunc("/api/health", healthHandler).Methods("GET")
	log.Info("Server listening on port ", config.PORT)
	log.Info("Environment: ", config.ENV)
	addr := fmt.Sprintf(":%s", config.PORT)
	log.Fatal(http.ListenAndServe(addr, router))
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
