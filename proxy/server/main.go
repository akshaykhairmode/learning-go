package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {

	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":9001", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		reqData, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("error"))
			return
		}

		log.Println("Request at endpoint :", string(reqData))

		rw.WriteHeader(http.StatusOK)
		rw.Header().Add("test-header", "test-header-value")
		rw.Write([]byte("Endpoint called"))
	})))

}
