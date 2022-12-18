package limiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {

	endpointCalledCount := 0

	r := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	l := NewLimiter()
	l.SetRateLimit("/test", 5)

	hf := l.Middleware(func(w http.ResponseWriter, r *http.Request) {
		endpointCalledCount++
		time.Sleep(1 * time.Second)

	})

	for i := 0; i < 200; i++ {
		go hf(w, r)
	}

	r1 := httptest.NewRequest("GET", "/test", nil)
	r1.RemoteAddr = "test-addr"

	for i := 0; i < 200; i++ {
		go hf(w, r1)
	}

	if endpointCalledCount != 10 {
		t.Error("endpointCalledCount does not match", endpointCalledCount)
	}

}
