package main

import "net/http"

func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Root Path"))
	})

	http.HandleFunc("/first", handleFirst)
	http.HandleFunc("/second", handleSecond)

	http.ListenAndServe(":7000", nil)

}

func handleFirst(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("First"))
}

func handleSecond(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Second"))
}
