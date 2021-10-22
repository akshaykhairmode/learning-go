package main

import (
	"net/http"
)

func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World"))
	})

	http.HandleFunc("/public", handlePublic)
	http.HandleFunc("/private", authenticate(handlePrivate))

	http.ListenAndServe(":7000", nil)

}

func authenticate(next http.HandlerFunc) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {

		if r.Header.Get("api-key") == "" {
			rw.Write([]byte("Auth Error"))
			return
		}

		next(rw, r)
	}

}

func handlePublic(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Public Endpoint Called"))
}

func handlePrivate(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Private Endpoint Called"))
}
