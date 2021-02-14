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
	httpsrv  http.Server
}

func (srv *Server) Start() {
	for _, w := range srv.workers {
		go w.start(srv.debuglvl)
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
				select {
				case w.wake <- true:
				default:
				}
			}
		}
	}
}

func NewServer(cfg *config.Config, debuglvl int) *Server {
	srv := &Server{cfg: cfg, debuglvl: debuglvl}
	for _, g := range cfg.Groups {
		w := &worker{group: g, wake: make(chan bool)}
		srv.workers = append(srv.workers, w)
	}
	return srv
}
