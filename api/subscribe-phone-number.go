package api

import (
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

	phoneFrom := r.PostFormValue("From")
	if phoneFrom == "" {
		slog.Error("Parse form error, missing 'From' value")
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return

	}

	phoneNumber, _ := strings.CutPrefix(phoneFrom, "whatsapp:")

	phone, found := s.PhoneNumbers[phoneNumber]
	if !found {
		s.PhoneNumbers[phoneNumber] = struct{}{}
		slog.Info(fmt.Sprintf("New phone number added: %s", phoneNumber))
	} else {
		slog.Info(fmt.Sprintf("The phone number %s was already subscribed", phone))
	}
	return
}
