package limiter

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
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

func (l *limiter) incr(r *http.Request, now string) error {
	l.mux.Lock()
	defer l.mux.Unlock() //defering here so we dont forget later

	status := l.getStatus(r)
	if status == nil { //if the endpoint was not initialized with set rate limit status will be nil
		return nil
	}

	if status.Count[r.RemoteAddr+now] >= status.Limit { //Check if the limit is exceeded.
		return fmt.Errorf("rate limit reached for %s", r.RequestURI)
	}

	status.Count[r.RemoteAddr+now]++ //Increase the count for this IP.
	return nil
}

func (l *limiter) decr(r *http.Request, now string) {
	l.mux.Lock()
	defer l.mux.Unlock() //defering here so we dont forget later

	status := l.getStatus(r)
	if status == nil { //if the endpoint was not initialized with set rate limit status will be nil
		return
	}

	status.Count[r.RemoteAddr+now]-- //Decrease the count for this IP.

	if status.Count[r.RemoteAddr+now] <= 0 { //Delete the key once counter is 0, this is important now because for each second we will have one key.
		delete(status.Count, r.RemoteAddr+now)
	}

}

func (l *limiter) Middleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := strconv.FormatInt(time.Now().Unix(), 10)
		if err := l.incr(r, now); err != nil { //If limit is exceeded we return too many requests
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(err.Error()))
			return
		}
		defer l.decr(r, now)    //Decrement will happen after the actual handler is executed.
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
