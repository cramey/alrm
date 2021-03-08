package server

import (
	"fmt"
	"net/http"
	"git.binarythought.com/cdramey/alrm/api"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api":
		g := r.URL.Query()
		c := g.Get("cmd")
		if c == "" {
			http.Error(w, "no command given", http.StatusBadRequest)
			return
		}
		cmd, err := api.ParseCommand([]byte(c))
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing command: %s", err.Error()),
				http.StatusBadRequest)
			return
		}

		err = cmd.Validate(s.config.APIKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch cmd.Command {
		case "shutdown":
			s.shutdownc <- false

		case "restart":
			s.shutdownc <- true

		default:
			http.Error(w, fmt.Sprintf("unknown command: %s", cmd.Command),
				http.StatusBadRequest)
		}

	default:
		http.Error(w, "File not found", http.StatusNotFound)
	}
}
