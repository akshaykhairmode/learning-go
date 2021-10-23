package main

import (
	"net/http"
)

type myHandler struct{}

func (mh myHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	switch r.RequestURI {
	case "/first":
		rw.Write([]byte("First"))
	case "/second":
		rw.Write([]byte("Second"))
	case "/third":
		rw.Write([]byte("Third"))
	default:
		rw.Write([]byte("Zero"))
	}

}

func main() {

	handler := new(myHandler)

	mux := http.NewServeMux()
	mux.Handle("/", handler)
	mux.Handle("/first", handler)
	mux.Handle("/second", handler)
	mux.Handle("/third", handler)

	http.ListenAndServe(":7000", mux)

}
