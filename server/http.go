package server

import (
	"fmt"
	"git.binarythought.com/cdramey/alrm/api"
	"net/http"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api":
		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("form parse error: %s", err.Error()),
				http.StatusBadRequest)
			return
		}

		c := r.FormValue("cmd")
		if c == "" {
			http.Error(w, "no command given", http.StatusBadRequest)
			return
		}
		cmd, err := api.ParseCommand([]byte(c))
		if err != nil {
			http.Error(w, fmt.Sprintf("command parse error: %s", err.Error()),
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
