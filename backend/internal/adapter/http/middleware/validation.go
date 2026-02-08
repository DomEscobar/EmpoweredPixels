package middleware

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

// ValidateJSON middleware validates that request bodies are valid JSON
// and that required Content-Type header is present for POST/PUT requests
func ValidateJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only validate for methods that expect a body
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			contentType := r.Header.Get("Content-Type")
			
			// Check Content-Type header
			if contentType == "" {
				log.Printf("validation error: missing Content-Type for %s %s", r.Method, r.RequestURI)
				http.Error(w, `{"error": "Content-Type header required"}`, http.StatusBadRequest)
				return
			}

			// Ensure it's JSON
			if !strings.Contains(contentType, "application/json") {
				log.Printf("validation error: invalid Content-Type %s for %s %s", contentType, r.Method, r.RequestURI)
				http.Error(w, `{"error": "Content-Type must be application/json"}`, http.StatusBadRequest)
				return
			}

			// Validate request body is valid JSON by attempting to parse
			if r.ContentLength > 0 {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					log.Printf("validation error: failed to read body for %s %s: %v", r.Method, r.RequestURI, err)
					http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
					return
				}

				// Validate JSON structure
				var jsonData interface{}
				if err := json.Unmarshal(body, &jsonData); err != nil {
					log.Printf("validation error: invalid JSON for %s %s: %v", r.Method, r.RequestURI, err)
					http.Error(w, `{"error": "invalid JSON in request body"}`, http.StatusBadRequest)
					return
				}

				// Reset body for downstream handlers
				r.ContentLength = int64(len(body))
				r.Body = io.NopCloser(strings.NewReader(string(body)))
			}
		}

		next.ServeHTTP(w, r)
	})
}

// ValidateAuthHeader ensures Authorization header is properly formatted
func ValidateAuthHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			// Authorization header must be "Bearer <token>"
			if !strings.HasPrefix(auth, "Bearer ") {
				log.Printf("validation error: malformed Authorization header for %s %s", r.Method, r.RequestURI)
				http.Error(w, `{"error": "invalid Authorization header format"}`, http.StatusBadRequest)
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			if token == "" {
				log.Printf("validation error: empty token in Authorization header for %s %s", r.Method, r.RequestURI)
				http.Error(w, `{"error": "Authorization token cannot be empty"}`, http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// ValidateQueryParams checks for suspicious or malformed query parameters
func ValidateQueryParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters to ensure they're valid
		if err := r.ParseForm(); err != nil {
			log.Printf("validation error: invalid query parameters for %s %s: %v", r.Method, r.RequestURI, err)
			http.Error(w, `{"error": "invalid query parameters"}`, http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
