package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

const (
	subscribeTemplateId = "HXe47b0c14279142ed9840ba67155309e5"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)
)

type handlerFunc func(http.ResponseWriter, *http.Request) (int, any)

type Server struct {
	TwilioClient      *twilio.RestClient
	TwilioPhoneNumber string
	PhoneNumbers      map[string]interface{}
}

func NewServer() *Server {
	twilioClient := twilio.NewRestClient()
	if twilioClient == nil {
		slog.Error("Couldn't connect to twilio client")
		panic("Couldn't connect to twilio client")
	}

	twilioPhoneNumber, ok := os.LookupEnv("TWILIO_PHONE_NUMBER")
	if !ok {
		slog.Error("TWILIO_PHONE_NUMBER not set")
		panic("TWILIO_PHONE_NUMBER not set")
	}

	return &Server{
		TwilioClient:      twilioClient,
		TwilioPhoneNumber: twilioPhoneNumber,
		PhoneNumbers:      make(map[string]interface{}),
	}
}

func (s *Server) Run(listenAddr string) error {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)

	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/health", promMiddleware(s.healthHandler))

	http.Handle("POST /subscribe-phone-number", promMiddleware(s.subscribePhoneNumber))
	http.Handle("POST /broadcast", promMiddleware(s.broadcast))

	slog.Info("Server started", "ListenAddr", listenAddr)

	return http.ListenAndServe(listenAddr, nil)
}

func (s *Server) messageOne(contentSid string, to string) {
	messageParams := &twilioApi.CreateMessageParams{}

	messageParams.SetContentSid(contentSid)
	messageParams.SetFrom(s.TwilioPhoneNumber)
	messageParams.SetTo(to)

	resp, err := s.TwilioClient.Api.CreateMessage(messageParams)
	if err != nil {
		slog.Error("Failed to send message", "to", to, "errorKind", err)
		return
	}

	slog.Info("Message sent", "to", to, "response", resp)
}

func (s *Server) messageAll(contentSid string) {
	messageParams := &twilioApi.CreateMessageParams{}

	messageParams.SetContentSid(contentSid)
	messageParams.SetFrom(s.TwilioPhoneNumber)

	for key := range s.PhoneNumbers {
		messageParams.SetTo(key)

		resp, err := s.TwilioClient.Api.CreateMessage(messageParams)
		if err != nil {
			slog.Error("Failed to send message", "to", key, "errorKind", err)
			continue
		}

		slog.Info("Message sent", "to", key, "response", resp)
	}
}

func writeJson(w http.ResponseWriter, status int, v any) (int, any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("writeJson: failed to encode")
		panic(err)
	}

	return status, v
}
