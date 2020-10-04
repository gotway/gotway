package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gotway/microsamples/stock/config"
)

// NewAPI Starts a new HTTP server
func NewAPI() {
	router := mux.NewRouter()

	router.HandleFunc("/health", health).Methods(http.MethodGet)
	router.HandleFunc("/list", upsertStockList).Methods(http.MethodPost)
	router.HandleFunc("/list", getStockList).Methods(http.MethodGet)
	router.HandleFunc("/{id}", upsertStock).Methods(http.MethodPost)
	router.HandleFunc("/{id}", getStock).Methods(http.MethodGet)

	log.Print("Server listening on port ", config.Port)
	addr := fmt.Sprintf(":%d", config.Port)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal(err)
	}
}
