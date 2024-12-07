package api

import (
	"net/http"
)

type health struct {
	Status string `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) (int, any) {
	return writeJson(w, http.StatusOK, health{Status: "ok"})
}
