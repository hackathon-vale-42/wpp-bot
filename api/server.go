package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
  "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
  
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
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

type TemplateBody struct {
	TemplateId string `json:"template_id"`
	From       string `json:"from"`
}

type Server struct {
	TwilioClient *twilio.RestClient
	PhoneNumbers map[string]interface{}
}

type handlerFunc func(http.ResponseWriter, *http.Request) (int, any)

func NewServer() *Server {
	twilioClient := twilio.NewRestClient()
	if twilioClient == nil {
		slog.Error("Couldn't connect to twilio client")
		panic("Couldn't connect to twilio client")
	}

	return &Server{
		TwilioClient: twilioClient,
		PhoneNumbers: make(map[string]interface{}),
	}
}

func (s *Server) sendTemplate(w http.ResponseWriter, r *http.Request) {

	var body TemplateBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("Decode error", "error", err)
		http.Error(w, "Invalid body", http.StatusBadRequest)
	}

	// send template to each phone number subscribed
	messageParams := &twilioApi.CreateMessageParams{}
	messageParams.SetContentSid(body.TemplateId)
	messageParams.SetFrom(fmt.Sprintf("whatsapp:%s", body.From))
	for key := range s.PhoneNumbers {
		messageParams.SetTo(fmt.Sprintf("whatsapp:%s", key))
		resp, err := s.TwilioClient.Api.CreateMessage(messageParams)
		slog.Info(fmt.Sprintf("send for client: %s", key), "response", resp)
		if err != nil {
			slog.Error("error on twilio call", "kind", err)
			http.Error(w, "Error on twilio api call", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (s *Server) Run(listenAddr string) error {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("POST /message", s.sendTemplate)
	http.HandleFunc("POST /subscribe-phone-number", s.SubscribePhoneNumber)

	http.Handle("/health", promMiddleware(healthHandler))

	slog.Info("Server started", "ListenAddr", listenAddr)

	return http.ListenAndServe(listenAddr, nil)
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
