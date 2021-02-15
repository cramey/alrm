package server

import (
	"context"
	"fmt"
	"git.binarythought.com/cdramey/alrm/config"
	"net"
	"net/http"
	"time"
)

type Server struct {
	workers   []*worker
	cfg       *config.Config
	shutdownc chan bool
	debuglvl  int
	http      http.Server
}

func (srv *Server) Start() (bool, error) {
	listen, err := net.Listen("tcp", srv.cfg.Listen)
	if err != nil {
		return false, err
	}

	for _, w := range srv.workers {
		go w.start()
	}

	srv.http = http.Server{Handler: srv}
	go srv.http.Serve(listen)

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
		case b := <-srv.shutdownc:
			srv.http.Shutdown(context.Background())
			for _, w := range srv.workers {
				w.shutdown()
			}
			return b, nil
		}
	}
}

func NewServer(cfg *config.Config, debuglvl int) *Server {
	srv := &Server{
		cfg:       cfg,
		debuglvl:  debuglvl,
		shutdownc: make(chan bool, 1),
	}
	for _, g := range cfg.Groups {
		srv.workers = append(
			srv.workers, makeworker(g, debuglvl),
		)
	}
	return srv
}
