package api

import (
	"log/slog"
	"net/http"
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

func (s *Server) subscribePhoneNumber(w http.ResponseWriter, r *http.Request) (int, any) {
	phoneNumber := r.PostFormValue("From")
	if phoneNumber == "" {
		slog.Error("Parse form error, missing 'From' value")
		return writeJson(w, http.StatusBadRequest, nil)
	}

	_, found := s.PhoneNumbers[phoneNumber]
	if found {
		slog.Warn("Phone number already subscribed", "phoneNumber", phoneNumber)
		return writeJson(w, http.StatusBadRequest, nil)
	}

	s.PhoneNumbers[phoneNumber] = struct{}{}
	slog.Info("Phone number subscribed", "phoneNumber", phoneNumber)

	s.messageOne(s.TwilioInfo.SubscribeConfirmationSid, phoneNumber)

	return writeJson(w, http.StatusOK, nil)
}
