package pprof

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	"github.com/gotway/gotway/pkg/log"
)

type Options struct {
	Port string
}

type PProf struct {
	options Options
	logger  log.Logger
	server  *http.Server
}

func (p *PProf) Start() {
	p.logger.Infof("profiling server listening in :%s", p.options.Port)
	if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		p.logger.Error("error staring profiling server ", err)
	}
}

func (p *PProf) Stop() {
	p.logger.Info("stopping profiling server")
	if err := p.server.Shutdown(context.Background()); err != nil {
		p.logger.Error("error stopping profiling server ", err)
	}
}

func New(options Options, logger log.Logger) *PProf {
	addr := ":" + options.Port
	server := &http.Server{Addr: addr}
	return &PProf{options, logger, server}
}
