package limiter

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type status struct {
	Limit uint64            //The rate limit set for the endpoint
	Count map[string]uint64 //Key will be RemoteAddr and value will be the number of current connections
}

type limiter struct {
	det map[string]*status //Key will be the endpoint here.
	mux sync.RWMutex       //For handling concurrent read write on maps.
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
	defer l.mux.Unlock() //defering here so we dont forget later

	status := l.getStatus(r)
	if status == nil { //if the endpoint was not initialized with set rate limit status will be nil
		return nil
	}

	if status.Count[r.RemoteAddr] >= status.Limit { //Check if the limit is exceeded.
		return fmt.Errorf("rate limit reached for %s", r.RequestURI)
	}

	status.Count[r.RemoteAddr]++ //Increase the count for this IP.
	return nil
}

func (l *limiter) decr(r *http.Request) {
	l.mux.Lock()
	defer l.mux.Unlock() //defering here so we dont forget later

	status := l.getStatus(r)
	if status == nil { //if the endpoint was not initialized with set rate limit status will be nil
		return
	}

	status.Count[r.RemoteAddr]-- //Decrease the count for this IP.
}

func (l *limiter) Middleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := l.incr(r); err != nil { //If limit is exceeded we return too many requests
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(err.Error()))
			return
		}
		defer l.decr(r)         //Decrement will happen after the actual handler is executed.
		handler.ServeHTTP(w, r) //Actual handler
	}
}

func (l *limiter) getStatus(r *http.Request) *status {

	status, ok := l.det[r.RequestURI] //Check if the endoint exists in the map.

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
