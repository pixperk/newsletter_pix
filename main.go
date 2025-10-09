package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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
	Database  string    `json:"database"`
	Uptime    string    `json:"uptime"`
}

var startTime = time.Now()

// HealthCheckHandler handles health check requests
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check database connection
	dbStatus := "healthy"
	if err := database.DB.Ping(); err != nil {
		dbStatus = "unhealthy"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// Calculate uptime
	uptime := time.Since(startTime).Round(time.Second)

	response := HealthCheckResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "newsletter-api",
		Version:   "1.0.0",
		Database:  dbStatus,
		Uptime:    uptime.String(),
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	godotenv.Load()
	database.InitDB()

	// Create rate limiters for different endpoints
	// Email endpoints: 5 requests per 15 minutes per IP
	emailRateLimit := handlers.NewRateLimiter(5, 15*time.Minute)

	// General endpoints: 100 requests per 15 minutes per IP
	generalRateLimit := handlers.NewRateLimiter(100, 15*time.Minute)

	// Admin endpoints (send): 10 requests per hour per IP
	adminRateLimit := handlers.NewRateLimiter(10, 1*time.Hour)

	// Apply rate limiting to routes
	http.HandleFunc("/health", enableCORS(handlers.RateLimitMiddleware(generalRateLimit)(HealthCheckHandler)))
	http.HandleFunc("/subscribe", enableCORS(handlers.RateLimitMiddleware(emailRateLimit)(handlers.SubscribeHandler)))
	http.HandleFunc("/subscribe/verify", enableCORS(handlers.RateLimitMiddleware(emailRateLimit)(handlers.VerifySubscribeHandler)))
	http.HandleFunc("/subscribe/confirm", enableCORS(handlers.RateLimitMiddleware(emailRateLimit)(handlers.VerifyConfirmHandler)))
	http.HandleFunc("/unsubscribe", enableCORS(handlers.RateLimitMiddleware(emailRateLimit)(handlers.UnsubscribeHandler)))
	http.HandleFunc("/send", enableCORS(handlers.RateLimitMiddleware(adminRateLimit)(handlers.SendHandler(database.DB))))
	http.HandleFunc("/test-send", enableCORS(handlers.RateLimitMiddleware(adminRateLimit)(handlers.TestSendHandler)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on :%s", port)
	log.Printf("Rate limiting enabled:")
	log.Printf("- Email endpoints: 5 requests per 15 minutes per IP")
	log.Printf("- General endpoints: 100 requests per 15 minutes per IP")
	log.Printf("- Admin endpoints: 10 requests per hour per IP")
	http.ListenAndServe(":"+port, nil)
}
