package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	database "github.com/pixperk/newsletter/db"
	"github.com/pixperk/newsletter/handlers"
)

var allowedOrigins = map[string]bool{
	"https://pixperk.tech":     true,
	"https://www.pixperk.tech": true,
}

const maxBodySize = 1 << 20 // 1 MB

func securityHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		next(w, r)
	}
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Secret")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func limitBody(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		next(w, r)
	}
}

func wrap(h http.HandlerFunc, rl *handlers.RateLimiter) http.HandlerFunc {
	return securityHeaders(enableCORS(limitBody(handlers.RateLimitMiddleware(rl)(h))))
}

// HealthCheckResponse represents the health check response structure
type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
	Database  string    `json:"database"`
	Uptime    string    `json:"uptime"`
}

var startTime = time.Now()

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dbStatus := "healthy"
	status := "ok"
	if err := database.DB.Ping(); err != nil {
		dbStatus = "unhealthy"
		status = "degraded"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	uptime := time.Since(startTime).Round(time.Second)

	response := HealthCheckResponse{
		Status:    status,
		Timestamp: time.Now(),
		Service:   "newsletter-api",
		Version:   "1.0.0",
		Database:  dbStatus,
		Uptime:    uptime.String(),
	}

	json.NewEncoder(w).Encode(response)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "pong"})
}

func main() {
	godotenv.Load()
	database.InitDB()
	defer database.DB.Close()

	emailRateLimit := handlers.NewRateLimiter(5, 15*time.Minute)
	generalRateLimit := handlers.NewRateLimiter(100, 15*time.Minute)
	adminRateLimit := handlers.NewRateLimiter(10, 1*time.Hour)

	http.HandleFunc("/ping", wrap(PingHandler, generalRateLimit))
	http.HandleFunc("/health", wrap(HealthCheckHandler, generalRateLimit))
	http.HandleFunc("/subscribe", wrap(handlers.SubscribeHandler, emailRateLimit))
	http.HandleFunc("/subscribe/verify", wrap(handlers.VerifySubscribeHandler, emailRateLimit))
	http.HandleFunc("/subscribe/confirm", wrap(handlers.VerifyConfirmHandler, emailRateLimit))
	http.HandleFunc("/unsubscribe", wrap(handlers.UnsubscribeHandler, emailRateLimit))
	http.HandleFunc("/send", wrap(handlers.SendHandler(database.DB), adminRateLimit))
	http.HandleFunc("/test-send", wrap(handlers.TestSendHandler, adminRateLimit))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped gracefully")
}
