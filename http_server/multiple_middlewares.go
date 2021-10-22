package main

import (
	"net/http"
)

type middleWare func(http.HandlerFunc) http.HandlerFunc

func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World"))
	})

	http.HandleFunc("/public", handlePublic)
	http.HandleFunc("/private", wrap(handlePrivate, authenticateApiKey))
	http.HandleFunc("/private_ip", wrap(handlePrivate, authenticateApiKey, authenticateIP))

	http.ListenAndServe(":7000", nil)

}

func wrap(l http.HandlerFunc, midW ...middleWare) http.HandlerFunc {

	//do it in reverse to preserve order
	for i := len(midW) - 1; i >= 0; i-- {
		l = midW[i](l)
	}

	return l
}

func authenticateApiKey(next http.HandlerFunc) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {

		if r.Header.Get("api-key") == "" {
			rw.Write([]byte("Auth Error\n"))
			return
		}

		next(rw, r)
	}

}

func authenticateIP(next http.HandlerFunc) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {

		if r.RemoteAddr != "0.0.0.0" {
			rw.Write([]byte("IP should be 0.0.0.0\n"))
			return
		}

		next(rw, r)
	}

}

func handlePublic(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Public Endpoint Called\n"))
}

func handlePrivate(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Private Endpoint Called\n"))
}
