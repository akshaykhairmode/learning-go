package main

import (
	"net/http"
	limiter "rate-limiter/limiter"
	"time"
)

func main() {

	rateLimiter := limiter.NewLimiter()
	rateLimiter.SetRateLimit("/rate-limit", 5)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(204)
	})

	http.HandleFunc("/rate-limit", rateLimiter.Middleware(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(204)
	}))

	http.ListenAndServe(":9000", nil)

}
