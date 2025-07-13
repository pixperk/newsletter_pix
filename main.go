package main

import (
	"log"
	"net/http"

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

func main() {
	godotenv.Load()
	database.InitDB()

	http.HandleFunc("/subscribe", enableCORS(handlers.SubscribeHandler))
	http.HandleFunc("/unsubscribe", enableCORS(handlers.UnsubscribeHandler))
	http.HandleFunc("/send", enableCORS(handlers.SendHandler(database.DB)))

	log.Println("Listening on :9000")
	http.ListenAndServe(":9000", nil)
}
