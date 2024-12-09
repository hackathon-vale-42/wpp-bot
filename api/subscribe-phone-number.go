package api

import (
	"log/slog"
	"net/http"
	"time"
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
		return http.StatusBadRequest, nil
	}

	_, found := s.PhoneNumbers[phoneNumber]
	if found {
		slog.Warn("Phone number already subscribed", "phoneNumber", phoneNumber)
		s.messageOne(s.TwilioInfo.AlreadySubscribedSid, phoneNumber)
		return writeJSON(w, http.StatusConflict, nil)
	}

	if _, err := s.RedisClient.Set(s.Ctx, phoneNumber, nil, 24*time.Hour).Result(); err != nil {
		slog.Error("Failed to set phone number on redisClient", "phoneNumber", phoneNumber, "error", err)
		return writeJSON(w, http.StatusInternalServerError, nil)
	}

	s.PhoneNumbers[phoneNumber] = struct{}{}
	slog.Info("Phone number subscribed", "phoneNumber", phoneNumber)

	s.messageOne(s.TwilioInfo.SubscribeConfirmationSid, phoneNumber)

	w.WriteHeader(http.StatusNoContent)
	return http.StatusNoContent, nil
}
