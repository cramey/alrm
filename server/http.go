package server

import (
	"fmt"
	"net/http"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/shutdown":
		fmt.Fprintf(w, "shutting down .. ")
		s.shutdownc <- false
	case "/restart":
		fmt.Fprintf(w, "restarting .. ")
		s.shutdownc <- true
	default:
		fmt.Fprintf(w, "Hello, world!")
	}
}
