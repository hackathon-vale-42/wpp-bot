package api

import (
	"fmt"
	"log/slog"
	"os"
)

type TwilioInfo struct {
	PhoneNumber              string
	SubscribeConfirmationSid string
	ChildCareSid             string
	ChristmasSid             string
	FuneralSid               string
	GlassesSid               string
	WellhubSid               string
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
		SubscribeConfirmationSid: getEnvOrPanic("SUBSCRIBE_CONFIRMATION_SID"),
		ChildCareSid:             getEnvOrPanic("CHILD_CARE_SID"),
		ChristmasSid:             getEnvOrPanic("CHRISTMAS_SID"),
		FuneralSid:               getEnvOrPanic("FUNERAL_SID"),
		GlassesSid:               getEnvOrPanic("GLASSES_SID"),
		WellhubSid:               getEnvOrPanic("WELLHUB_SID"),
	}
}
