package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gosmo-devs/microgateway/config"
	"github.com/gosmo-devs/microgateway/log"
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
	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/health", healthHandler).Methods("GET")
	createServiceRouting(apiRouter)
	router.PathPrefix("/{service}").HandlerFunc(proxyHandler)

	return router
}

func createServiceRouting(rootRouter *mux.Router) {
	rootRouter.HandleFunc("/services", getServicesHandler).Methods("GET")

	serviceRouter := rootRouter.PathPrefix("/service").Subrouter()
	serviceRouter.Methods("POST").HandlerFunc(registerServiceHandler)

	idRouter := serviceRouter.PathPrefix("/{service}").Subrouter()
	idRouter.Methods("GET").HandlerFunc(getServiceHandler)
	idRouter.Methods("DELETE").HandlerFunc(deleteServiceHandler)
}
