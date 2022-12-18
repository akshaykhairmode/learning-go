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

//date;xargs -I % -P 10 curl -I "localhost:9000/health" < <(printf '%s\n' {1..50}) 2>&1 | grep "204 No Content" | wc -l;date

//date;xargs -I % -P 10 curl -I "localhost:9000/rate-limit" < <(printf '%s\n' {1..50}) 2>&1 | grep "204 No Content" | wc -l;date
