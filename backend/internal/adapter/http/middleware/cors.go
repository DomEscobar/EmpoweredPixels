package middleware

import (
	"net/http"
	"strings"
)

// AllowedOrigins defines the list of trusted origins for CORS
// Update this list based on your deployment environment
var AllowedOrigins = []string{
	"http://localhost:3000",      // Local development frontend
	"http://localhost:5173",      // Vite dev server alternative
	"http://127.0.0.1:3000",      // Local testing
	"http://127.0.0.1:5173",      // Local testing
	// Add production domain(s) here, e.g.:
	// "https://app.example.com",
	// "https://www.example.com",
}

func isOriginAllowed(origin string) bool {
	for _, allowed := range AllowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Validate origin against whitelist
		if origin != "" && isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
