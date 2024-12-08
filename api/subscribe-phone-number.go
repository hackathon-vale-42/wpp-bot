package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type WhatsappMessage struct {
	WaId          string `json:"WaId"`
	Body          string `json:"Body"`
	SmsMessageSid string `json:"SmsMessageSid"`
	AccountSid    string `json:"AccountSid"`
	ProfileName   string `json:"ProfileName"`
	To            string `json:"To"`
	NumMedia      string `json:"NumMedia"`
	From          string `json:"From"`
	SmsStatus     string `json:"SmsStatus"`
}

func (s *Server) SubscribePhoneNumber(w http.ResponseWriter, r *http.Request) {
	var body WhatsappMessage

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("Decode error", "error", err)
		http.Error(w, "Invalid body", http.StatusBadRequest)
	}

	phoneNumber, _ := strings.CutPrefix(body.From, "whatsapp:")

	phone, found := s.PhoneNumbers[phoneNumber]
	if !found {
		s.PhoneNumbers[phoneNumber] = struct{}{}
		slog.Info(fmt.Sprintf("New phone number added: %s", phoneNumber))
	} else {
		slog.Info(fmt.Sprintf("The phone number %s was already subscribed", phone))
	}
	return
}
