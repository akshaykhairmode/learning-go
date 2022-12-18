package limiter

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type status struct {
	Limit, Curr uint64
}

type limiter struct {
	det map[string]*status
	mux sync.RWMutex
}

func (l *limiter) SetRateLimit(endpoint string, limit uint64) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.det[endpoint] = &status{
		Limit: limit,
		Curr:  0,
	}
}

func (l *limiter) incr(endpoint string) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	status, ok := l.det[endpoint]
	if !ok {
		log.Println("incr: endpoint not found")
		return nil
	}

	if status.Curr >= status.Limit {
		return fmt.Errorf("rate limit reached for %s", endpoint)
	}

	status.Curr++
	return nil
}

func (l *limiter) decr(endpoint string) {
	l.mux.Lock()
	defer l.mux.Unlock()

	status, ok := l.det[endpoint]
	if !ok {
		log.Println("decr: endpoint not found")
		return
	}

	status.Curr--
}

func (l *limiter) Middleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := l.incr(r.RequestURI); err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(err.Error()))
			return
		}
		defer l.decr(r.RequestURI)
		handler.ServeHTTP(w, r)
	}
}
func NewLimiter() limiter {
	return limiter{
		det: map[string]*status{},
		mux: sync.RWMutex{},
	}
}
