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
	http.Handle(m.options.Path, promhttp.Handler())
	m.logger.Infof("metrics server listening in %v:%s", m.options.Path, m.options.Port)

	if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		m.logger.Error("error staring metrics server ", err)
		return
	}
}

func (m *Metrics) Stop() {
	if err := m.server.Shutdown(context.Background()); err != nil {
		m.logger.Error("error stopping metrics server ", err)
		return
	}
	m.logger.Info("stopped metrics server")
}

func New(options Options, logger log.Logger) *Metrics {
	addr := ":" + options.Port
	server := &http.Server{Addr: addr}
	return &Metrics{options, server, logger}
}
