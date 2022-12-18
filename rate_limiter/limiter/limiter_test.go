package limiter

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {

	endpointCalledCount := int64(0)

	r := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	l := NewLimiter()          //Create our limiter instance
	l.SetRateLimit("/test", 5) //set rate limit to 5 for /test endpoint.

	hf := l.Middleware(func(w http.ResponseWriter, r *http.Request) { //This is our actual handler
		atomic.AddInt64(&endpointCalledCount, 1) //We will increase the count in the actual handler.
		time.Sleep(1 * time.Second)              //We are sleeping here so that we get the realistic count increase.
	})

	r1 := httptest.NewRequest("GET", "/test", nil)
	r1.RemoteAddr = "test-addr" //We create a new request here with different IP addr.

	for i := 0; i < 200; i++ { //We start 200 go routines and call the handler with our request. Since we have rate limiting enabled for /test endpoint our actual handler should only be called 5 times.
		go hf(w, r)
		go hf(w, r1) //We use the same recorder as we want to increase the same count variable.
	}

	time.Sleep(5 * time.Second) //We sleep here as we want to wait for all the go routines to finish.

	if endpointCalledCount != 10 { //We are check for 10 here as our handler should get called 5 times for default remoteAddr and 5 times for test-addr
		t.Error("endpointCalledCount does not match", endpointCalledCount)
	}

}
