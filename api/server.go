package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
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

type handlerFunc func(http.ResponseWriter, *http.Request) (int, any)

type Server struct {
	TwilioClient *twilio.RestClient
	TwilioInfo   *clientInfo
	Ctx          context.Context
	RedisClient  *redis.Client
	PhoneNumbers map[string]struct{}
}

func NewServer() *Server {
	twilioClient := twilio.NewRestClient()
	if twilioClient == nil {
		slog.Error("Couldn't connect to twilio client")
		return nil
	}

	twilioInfo := NewClientInfo()
	if twilioInfo == nil {
		slog.Error("Couldn't create twilio info")
		return nil
	}

	url, found := os.LookupEnv("REDIS_URL")
	if !found {
		slog.Error("Missing REDIS_URL environment variable")
		return nil
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		slog.Error("Failed to parse REDIS_URL", "error", err)
		return nil
	}

	redisClient := redis.NewClient(opts)

	return &Server{
		TwilioClient: twilioClient,
		TwilioInfo:   twilioInfo,
		Ctx:          context.Background(),
		RedisClient:  redisClient,
		PhoneNumbers: make(map[string]struct{}),
	}
}

func (s *Server) LoadPhoneNumbers() error {
	keys, _, err := s.RedisClient.Scan(s.Ctx, 0, "whatsapp:*", 50).Result()
	if err != nil {
		slog.Error("Failed to load phone numbers", "error", err)
		return err
	}

	for _, key := range keys {
		s.PhoneNumbers[key] = struct{}{}
	}

	return nil
}

func (s *Server) Run(listenAddr string) error {
	if err := s.LoadPhoneNumbers(); err != nil {
		return err
	}

	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)

	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/health", promMiddleware(s.healthHandler))

	http.Handle("POST /subscribe-phone-number", promMiddleware(s.subscribePhoneNumber))
	http.Handle("POST /broadcast", promMiddleware(s.broadcast))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/app/static"))))

	slog.Info("Server started", "listenAddr", listenAddr)
	return http.ListenAndServe(listenAddr, nil)
}

func (s *Server) messageOne(contentSid string, to string) {
	messageParams := &twilioApi.CreateMessageParams{}

	messageParams.SetContentSid(contentSid)
	messageParams.SetFrom(s.TwilioInfo.PhoneNumber)
	messageParams.SetTo(to)

	resp, err := s.TwilioClient.Api.CreateMessage(messageParams)
	if err != nil {
		slog.Error("Failed to send message", "to", to, "error", err)
		return
	}

	slog.Info("Message sent", "to", to, "response", resp)
}

func (s *Server) messageAll(contentSid string) {
	messageParams := &twilioApi.CreateMessageParams{}

	messageParams.SetContentSid(contentSid)
	messageParams.SetFrom(s.TwilioInfo.PhoneNumber)

	for key := range s.PhoneNumbers {
		messageParams.SetTo(key)

		resp, err := s.TwilioClient.Api.CreateMessage(messageParams)
		if err != nil {
			slog.Error("Failed to send message", "to", key, "error", err)
			continue
		}

		slog.Info("Message sent", "to", key, "response", resp)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) (int, any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("writeJson: failed to encode")
		panic(err)
	}

	return status, v
}
