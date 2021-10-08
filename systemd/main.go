package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	port := ":9000"

	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Reached Handler\n"))
		log.Printf("Logging Request URL : %v", r.RequestURI)
	})

	log.Fatal(http.ListenAndServe(port, nil))

}
