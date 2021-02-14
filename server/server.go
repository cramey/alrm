package server

import (
	"fmt"
	"git.binarythought.com/cdramey/alrm/config"
	"time"
)

type Server struct {
	workers  []*worker
	cfg      *config.Config
	debuglvl int
}

func (srv *Server) Start() {
	for _, w := range srv.workers {
		go w.start()
	}

	t := time.NewTicker(srv.cfg.Interval)
	defer t.Stop()
	for {
		select {
		case r := <-t.C:
			if srv.debuglvl > 0 {
				fmt.Printf("interval check at %s\n", r)
			}
			for _, w := range srv.workers {
				w.wake()
			}
		}
	}
}

func NewServer(cfg *config.Config, debuglvl int) *Server {
	srv := &Server{cfg: cfg, debuglvl: debuglvl}
	for _, g := range cfg.Groups {
		srv.workers = append(srv.workers, makeworker(g, debuglvl))
	}
	return srv
}
