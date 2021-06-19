package metrics

import (
	"context"
	"net/http"

	"github.com/gotway/gotway/pkg/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Options struct {
	Path string
	Port string
}

type Metrics struct {
	options Options
	server  *http.Server
	logger  log.Logger
}

func (m *Metrics) Start() {
	m.logger.Infof("metrics server listening in %v:%s", m.options.Path, m.options.Port)
	http.Handle(m.options.Path, promhttp.Handler())
	if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		m.logger.Error("error staring metrics server ", err)
	}
}

func (m *Metrics) Stop() {
	m.logger.Info("stopping metrics server")
	if err := m.server.Shutdown(context.Background()); err != nil {
		m.logger.Error("error stopping metrics server ", err)
	}
}

func New(options Options, logger log.Logger) *Metrics {
	addr := ":" + options.Port
	server := &http.Server{Addr: addr}
	return &Metrics{options, server, logger}
}
