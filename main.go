package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	database "github.com/pixperk/newsletter/db"
	"github.com/pixperk/newsletter/handlers"
)

// CORS middleware to handle cross-origin requests
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Secret")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next(w, r)
	}
}

// HealthCheckResponse represents the health check response structure
type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}

var startTime = time.Now()

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uptime := time.Since(startTime).Round(time.Second)

	response := HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "newsletter-api",
		Version:   "1.0.0",
		Uptime:    uptime.String(),
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	godotenv.Load()
	database.InitDB()

	http.HandleFunc("/health", enableCORS(HealthCheckHandler))
	http.HandleFunc("/subscribe", enableCORS(handlers.SubscribeHandler))
	http.HandleFunc("/unsubscribe", enableCORS(handlers.UnsubscribeHandler))
	http.HandleFunc("/send", enableCORS(handlers.SendHandler(database.DB)))
	http.HandleFunc("/test-send", enableCORS(handlers.TestSendHandler))

	log.Println("Listening on :9000")
	http.ListenAndServe(":9000", nil)
}
