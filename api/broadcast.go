package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type broadcastBody struct {
	ContentSid string `json:"content_sid"`
}

func (s *Server) broadcast(w http.ResponseWriter, r *http.Request) (int, any) {
	var body broadcastBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("Decode error", "error", err)
		return writeJSON(w, http.StatusBadRequest, err)
	}

	s.messageAll(body.ContentSid)
	return writeJSON(w, http.StatusOK, nil)
}
