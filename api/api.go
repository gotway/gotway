package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gotway/gotway/config"
	"github.com/gotway/gotway/log"
)

// NewAPI Starts a new HTTP server
func NewAPI() {
	router := createRouter()
	addr := fmt.Sprintf(":%s", config.Port)
	log.Logger.Info("Server listening on port ", config.Port)
	log.Logger.Info("Environment: ", config.Env)
	var err error
	if config.TLS {
		err = http.ListenAndServeTLS(addr, config.TLScert, config.TLSkey, router)
	} else {
		err = http.ListenAndServe(addr, router)
	}
	if err != nil {
		log.Logger.Fatal(err)
	}
}

func createRouter() *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", healthHandler).Methods(http.MethodGet)
	addServiceRouter(api)
	addCacheRouter(api)

	proxy := router.PathPrefix("/{service}").Subrouter()
	proxy.Use(cacheMiddleware)
	proxy.Schemes(getSchemes()...).HandlerFunc(proxyHandler)

	return router
}

func addServiceRouter(root *mux.Router) {
	root.HandleFunc("/services", getServicesHandler).Methods(http.MethodGet)

	service := root.PathPrefix("/service").Subrouter()
	service.Methods(http.MethodPost).HandlerFunc(registerServiceHandler)

	serviceID := service.PathPrefix("/{service}").Subrouter()
	serviceID.Methods(http.MethodGet).HandlerFunc(getServiceHandler)
	serviceID.Methods(http.MethodDelete).HandlerFunc(deleteServiceHandler)
}

func addCacheRouter(root *mux.Router) {
	cache := root.PathPrefix("/cache").Subrouter()

	cache.PathPrefix("/{service}").Methods(http.MethodGet).HandlerFunc(getCacheHandler)

	cache.Methods(http.MethodDelete).HandlerFunc(deleteCacheHandler)
}

func getSchemes() []string {
	var schemes = []string{"http"}
	if config.TLS {
		schemes = append(schemes, "https")
	}
	return schemes
}
