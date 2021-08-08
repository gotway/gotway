package http

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gotway/gotway/internal/cache"
	"github.com/gotway/gotway/internal/middleware"
	kubeCtrl "github.com/gotway/gotway/pkg/kubernetes/controller"
	"github.com/gotway/gotway/pkg/log"
)

type ServerOptions struct {
	Port string

	TLSenabled bool
	TLScert    string
	TLSkey     string
}

type Server struct {
	options     ServerOptions
	server      *http.Server
	handler     *handler
	middlewares []middleware.Middleware
	logger      log.Logger
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
	api.HandleFunc("/ingresses", s.handler.getIngresses).Methods(http.MethodGet)
	api.HandleFunc("/cache", s.handler.deleteCache).Methods(http.MethodDelete)
}

func (s *Server) addGatewayRouter(root *mux.Router) {
	gateway := root.PathPrefix("/").Subrouter()
	for _, m := range s.middlewares {
		gateway.Use(m.MiddlewareFunc)
	}
	gateway.PathPrefix("/").HandlerFunc(s.handler.writeResponse)
}

func NewServer(
	options ServerOptions,
	middlewares []middleware.Middleware,
	kubeCtrl *kubeCtrl.Controller,
	cacheCtrl cache.Controller,
	logger log.Logger,
) *Server {

	addr := ":" + options.Port

	return &Server{
		options: options,
		server:  &http.Server{Addr: addr},
		handler: newHandler(
			kubeCtrl,
			cacheCtrl,
			logger.WithField("type", "handler"),
		),
		middlewares: middlewares,
		logger:      logger,
	}
}
