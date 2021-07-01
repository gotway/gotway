package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gotway/gotway/internal/cache"
	"github.com/gotway/gotway/internal/middleware"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/service"
	"github.com/gotway/gotway/pkg/log"
)

type ServerOptions struct {
	Port string

	TLSenabled bool
	TLScert    string
	TLSkey     string
}

type Server struct {
	server            *http.Server
	options           ServerOptions
	middleware        *middleware.Middleware
	serviceController service.Controller
	cacheController   cache.Controller
	logger            log.Logger
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
	s.logger.Info("stopping server")
	if err := s.server.Shutdown(context.Background()); err != nil {
		s.logger.Error("error stopping server ", err)
	}
}

func (s *Server) createRouter() *mux.Router {
	root := mux.NewRouter()

	s.addApiRouter(root)

	proxy := root.MatcherFunc(s.matchRequest).Subrouter()
	proxy.Use(s.middleware.RequestDecorator)
	proxy.Use(s.middleware.Cache)
	proxy.PathPrefix("/").HandlerFunc(s.proxy)

	return root
}

func (s *Server) addApiRouter(root *mux.Router) {
	api := root.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
	s.addServiceRouter(api)
	s.addCacheRouter(api)
}

func (s *Server) addServiceRouter(root *mux.Router) {
	root.HandleFunc("/services", s.getServicesHandler).Methods(http.MethodGet)

	service := root.PathPrefix("/service").Subrouter()
	service.Methods(http.MethodPost).HandlerFunc(s.createService)

	serviceID := service.PathPrefix("/{service}").Subrouter()
	serviceID.Methods(http.MethodGet).HandlerFunc(s.getService)
	serviceID.Methods(http.MethodDelete).HandlerFunc(s.deleteService)
}

func (s *Server) addCacheRouter(root *mux.Router) {
	root.PathPrefix("/cache").Methods(http.MethodDelete).HandlerFunc(s.deleteCache)
}

func (s *Server) matchRequest(r *http.Request, _ *mux.RouteMatch) bool {
	service, err := getRequestService(r)
	if err != nil {
		return false
	}
	return service.MatchRequest(r)
}

func getRequestService(r *http.Request) (model.Service, error) {
	service, ok := r.Context().Value("service").(model.Service)
	if !ok {
		return model.Service{}, errors.New("service not found in request context")
	}
	return service, nil
}

func NewServer(options ServerOptions, cacheController cache.Controller,
	serviceController service.Controller, logger log.Logger) *Server {

	addr := ":" + options.Port
	server := &http.Server{Addr: addr}
	middleware := middleware.New(
		cacheController,
		serviceController,
		logger.WithField("type", "middleware"),
	)

	return &Server{
		server:            server,
		options:           options,
		middleware:        middleware,
		serviceController: serviceController,
		cacheController:   cacheController,
		logger:            logger,
	}
}
