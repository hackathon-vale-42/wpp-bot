package api

import (
	"log/slog"
	"os"
)

type clientInfo struct {
	PhoneNumber              string
	SubscribeConfirmationSid string
	AlreadySubscribedSid     string
}

func NewClientInfo() *clientInfo {
	phoneNumber, found := os.LookupEnv("TWILIO_PHONE_NUMBER")
	if !found {
		slog.Error("Missing TWILIO_PHONE_NUMBER environment variable")
		return nil
	}

	subscribeConfirmationSid, found := os.LookupEnv("SUBSCRIBE_CONFIRMATION_SID")
	if !found {
		slog.Error("Missing SUBSCRIBE_CONFIRMATION_SID environment variable")
		return nil
	}

	alreadySubscribedSid, found := os.LookupEnv("ALREADY_SUBSCRIBED_SID")
	if !found {
		slog.Error("Missing ALREADY_SUBSCRIBED_SID environment variable")
		return nil
	}

	return &clientInfo{
		PhoneNumber:              phoneNumber,
		SubscribeConfirmationSid: subscribeConfirmationSid,
		AlreadySubscribedSid:     alreadySubscribedSid,
	}
}
