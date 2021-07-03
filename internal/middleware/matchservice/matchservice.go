package matchservice

import (
	"net/http"

	httpError "github.com/gotway/gotway/internal/http/error"
	"github.com/gotway/gotway/internal/middleware"
	"github.com/gotway/gotway/internal/request"
	"github.com/gotway/gotway/internal/service"
	"github.com/gotway/gotway/pkg/log"
)

type matchService struct {
	serviceController service.Controller
	logger            log.Logger
}

func (m *matchService) MiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Debug("match service")
		services, err := m.serviceController.GetServices()
		if err != nil {
			httpError.Handle(err, w, m.logger)
			return
		}

		for _, s := range services {
			if s.MatchRequest(r) {
				next.ServeHTTP(w, request.WithService(r, s))
				return
			}
		}
		http.Error(w, "Service not found", http.StatusNotFound)
	})
}

func New(
	serviceController service.Controller,
	logger log.Logger,
) middleware.Middleware {

	return &matchService{
		serviceController: serviceController,
		logger:            logger,
	}
}
