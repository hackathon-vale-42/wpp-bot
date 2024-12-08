package api

import (
	"fmt"
	"log/slog"
	"os"
)

type TwilioInfo struct {
	PhoneNumber              string
	SubscribeConfirmationSid string
}

func getEnvOrPanic(key string) string {
	value, found := os.LookupEnv(key)
	if !found {
		err := fmt.Sprintf(key + " not set")

		slog.Error(err)
		panic(err)
	}

	return value
}

func NewTwilioInfo() *TwilioInfo {
	return &TwilioInfo{
		PhoneNumber:              getEnvOrPanic("TWILIO_PHONE_NUMBER"),
		SubscribeConfirmationSid: getEnvOrPanic("TWILIO_SUBSCRIBE_CONFIRMATION_SID"),
	}
}