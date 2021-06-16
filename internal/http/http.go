package http

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gotway/gotway/internal/cache"
	"github.com/gotway/gotway/internal/controller"
	"github.com/gotway/gotway/pkg/log"
)

type ServerOptions struct {
	Port string

	TLSenabled bool
	TLScert    string
	TLSkey     string
}

type Server struct {
	server  *http.Server
	options ServerOptions

	cacheController   cache.Controller
	serviceController controller.ServiceController

	logger log.Logger
}

func (s *Server) Start() {
	addr := ":" + s.options.Port
	router := s.createRouter()

	s.logger.Infof("server listening on port %s", s.options.Port)
	var err error
	if s.options.TLSenabled {
		err = http.ListenAndServeTLS(addr, s.options.TLScert, s.options.TLSkey, router)
	} else {
		err = http.ListenAndServe(addr, router)
	}
	if err != nil && err != http.ErrServerClosed {
		s.logger.Error("error starting server ", err)
	}
}

func (s *Server) Stop() {
	s.logger.Info("stopping server...")
	if err := s.server.Shutdown(context.Background()); err != nil {
		s.logger.Error("error stopping server ", err)
	}
}

func (s *Server) getSchemes() []string {
	var schemes = []string{"http"}
	if s.options.TLSenabled {
		schemes = append(schemes, "https")
	}
	return schemes
}

func (s *Server) createRouter() *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	s.addServiceRouter(api)
	s.addCacheRouter(api)

	proxy := router.PathPrefix("/{service}").Subrouter()
	proxy.Use(s.cacheMiddleware)
	proxy.Schemes(s.getSchemes()...).HandlerFunc(s.proxyHandler)

	return router
}

func (s *Server) addServiceRouter(root *mux.Router) {
	root.HandleFunc("/services", s.getServicesHandler).Methods(http.MethodGet)

	service := root.PathPrefix("/service").Subrouter()
	service.Methods(http.MethodPost).HandlerFunc(s.registerServiceHandler)

	serviceID := service.PathPrefix("/{service}").Subrouter()
	serviceID.Methods(http.MethodGet).HandlerFunc(s.getServiceHandler)
	serviceID.Methods(http.MethodDelete).HandlerFunc(s.deleteServiceHandler)
}

func (s *Server) addCacheRouter(root *mux.Router) {
	cache := root.PathPrefix("/cache").Subrouter()

	cache.PathPrefix("/{service}").Methods(http.MethodGet).HandlerFunc(s.getCacheHandler)

	cache.Methods(http.MethodDelete).HandlerFunc(s.deleteCacheHandler)
}

func NewServer(options ServerOptions, cacheController cache.Controller,
	serviceController controller.ServiceController, logger log.Logger) *Server {

	addr := ":" + options.Port
	server := &http.Server{Addr: addr}
	return &Server{
		server:            server,
		options:           options,
		cacheController:   cacheController,
		serviceController: serviceController,
		logger:            logger,
	}
}
