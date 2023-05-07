package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Root Path\n"))
	})

	http.HandleFunc("/first", handleFirst)
	http.HandleFunc("/second", handleSecond)

	fmt.Println("Started Server")
	fmt.Println(http.ListenAndServe(":7000", nil))

}

func handleFirst(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("First\n"))
}

func handleSecond(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Second\n"))
}
