package limiter

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type status struct {
	Limit uint64
	Count map[string]uint64
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
		Count: make(map[string]uint64),
	}
}

func (l *limiter) incr(r *http.Request) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	status := l.getStatus(r)
	if status == nil {
		return nil
	}

	if status.Count[r.RemoteAddr] >= status.Limit {
		return fmt.Errorf("rate limit reached for %s", r.RequestURI)
	}

	status.Count[r.RemoteAddr]++
	return nil
}

func (l *limiter) decr(r *http.Request) {
	l.mux.Lock()
	defer l.mux.Unlock()

	status := l.getStatus(r)
	if status == nil {
		return
	}

	status.Count[r.RemoteAddr]--
}

func (l *limiter) Middleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := l.incr(r); err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(err.Error()))
			return
		}
		defer l.decr(r)
		handler.ServeHTTP(w, r)
	}
}

func (l *limiter) getStatus(r *http.Request) *status {

	status, ok := l.det[r.RequestURI]

	if !ok {
		log.Printf("endpoint %s not found", r.RequestURI)
		return nil
	}

	return status
}

func NewLimiter() limiter {
	return limiter{
		det: map[string]*status{},
		mux: sync.RWMutex{},
	}
}
