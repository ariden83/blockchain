package middleware

import (
	"net/http"
	"time"
)

func JSONHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Accept-ranges", "items")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		now := time.Now()
		w.Header().Set("Date", now.String())
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
