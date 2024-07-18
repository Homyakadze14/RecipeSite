package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Logging requests
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		info := fmt.Sprintf("uri: %s method: %s duration: %s", r.RequestURI, r.Method, duration)
		log.Print(info)
	})
}
