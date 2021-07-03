package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gotway/gotway/internal/cache"
	"github.com/gotway/gotway/internal/service"
	"github.com/gotway/gotway/pkg/log"
)

type ServerOptions struct {
	Port           string
	GatewayTimeout time.Duration

	TLSenabled bool
	TLScert    string
	TLSkey     string
}

type Server struct {
	options    ServerOptions
	server     *http.Server
	handler    *handler
	middleware *middleware
	logger     log.Logger
}

func (s *Server) Start() {
	http.Handle("/", s.createRouter())
	s.logger.Infof("server listening on port %s", s.options.Port)

	var err error
	if s.options.TLSenabled {
		err = s.server.ListenAndServeTLS(s.options.TLScert, s.options.TLSkey)
	} else {
		err = s.server.ListenAndServe()
	}
	if err != nil && err != http.ErrServerClosed {
		s.logger.Error("error starting server ", err)
		return
	}
}

func (s *Server) Stop() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		s.logger.Error("error stopping server ", err)
		return
	}
	s.logger.Info("stopped server")
}

func (s *Server) createRouter() *mux.Router {
	root := mux.NewRouter()
	s.addApiRouter(root)
	s.addGatewayRouter(root)
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

func (s *Server) addGatewayRouter(root *mux.Router) {
	middlewares := []mux.MiddlewareFunc{
		s.middleware.matchService,
		s.middleware.cacheIn,
		s.middleware.gateway,
		s.middleware.cacheOut,
		s.middleware.writeResponse,
	}

	proxy := root.PathPrefix("/").Subrouter()
	proxy.PathPrefix("/")
	for _, m := range middlewares {
		proxy.Use(m)
	}
}

func (s *Server) addServiceRouter(root *mux.Router) {
	root.HandleFunc("/services", s.handler.getServices).Methods(http.MethodGet)

	service := root.PathPrefix("/service").Subrouter()
	service.Methods(http.MethodPost).HandlerFunc(s.handler.createService)

	serviceID := service.PathPrefix("/{service}").Subrouter()
	serviceID.Methods(http.MethodGet).HandlerFunc(s.handler.getService)
	serviceID.Methods(http.MethodDelete).HandlerFunc(s.handler.deleteService)
}

func (s *Server) addCacheRouter(root *mux.Router) {
	root.PathPrefix("/cache").Methods(http.MethodDelete).HandlerFunc(s.handler.deleteCache)
}

func NewServer(
	options ServerOptions,
	cacheController cache.Controller,
	serviceController service.Controller,
	logger log.Logger) *Server {

	addr := ":" + options.Port
	return &Server{
		options: options,
		server:  &http.Server{Addr: addr},
		handler: newHandler(
			serviceController,
			cacheController,
			logger.WithField("type", "handler"),
		),
		middleware: newMiddleware(
			middlewareOptions{gatewayTimeout: options.GatewayTimeout},
			serviceController,
			cacheController,
			logger.WithField("type", "middleware"),
		),
		logger: logger,
	}
}
