package api

import (
	"net/http"
)

type health struct {
	Status string `json:"status"`
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) (int, any) {
	return writeJSON(w, http.StatusOK, health{Status: "ok"})
}
