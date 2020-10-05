package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gotway/gotway/microservices/catalog/config"
)

// NewAPI Starts a new HTTP server
func NewAPI() {
	router := mux.NewRouter()

	router.HandleFunc("/health", health).Methods(http.MethodGet)
	router.HandleFunc("/products", getProducts).Methods(http.MethodGet)
	router.HandleFunc("/product", createProduct).Methods(http.MethodPost)
	router.HandleFunc("/product/{id}", getProduct).Methods(http.MethodGet)
	router.HandleFunc("/product/{id}", deleteProduct).Methods(http.MethodDelete)
	router.HandleFunc("/product/{id}", updateProduct).Methods(http.MethodPut)

	log.Print("Server listening on port ", config.Port)
	addr := fmt.Sprintf(":%s", config.Port)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal(err)
	}
}
