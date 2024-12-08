package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type BroadcastBody struct {
	TemplateId string `json:"template_id"`
}

func (s *Server) broadcast(w http.ResponseWriter, r *http.Request) (int, any) {
	var body BroadcastBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("Decode error", "errorKind", err)
		return writeJson(w, http.StatusBadRequest, nil)
	}

	s.messageAll(body.TemplateId)

	return writeJson(w, http.StatusOK, nil)
}
