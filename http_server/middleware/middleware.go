package main

import (
	"net/http"
)

const (
	APIKEY = "12345"
	TOKEN  = "token"
)

type middleWare func(http.HandlerFunc) http.HandlerFunc

func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World\n"))
	})

	http.HandleFunc("/landing", authorize(landing))
	http.HandleFunc("/profile", wrap(profile, authorize, authenticate))
	http.HandleFunc("/client_api", authenticate(client))

	http.ListenAndServe(":7000", nil)

}

func authenticate(next http.HandlerFunc) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {

		if r.Header.Get("TOKEN") != TOKEN {
			rw.Write([]byte("Invalid User\n"))
			return
		}

		next(rw, r)
	}

}

func authorize(next http.HandlerFunc) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {

		if r.Header.Get("api-key") != APIKEY {
			rw.Write([]byte("Auth Error\n"))
			return
		}

		next(rw, r)
	}

}

func wrap(l http.HandlerFunc, midW ...middleWare) http.HandlerFunc {

	//do it in reverse to preserve order
	for i := len(midW) - 1; i >= 0; i-- {
		l = midW[i](l)
	}

	return l
}

func landing(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Some Data for landing screen\n"))
}

func profile(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Some Data for profile\n"))
}

func client(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Some Data for client\n"))
}
