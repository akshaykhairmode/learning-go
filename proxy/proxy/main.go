package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

func main() {

	endpoint := flag.String("e", "", "")
	flag.Parse()

	if *endpoint == "" {
		log.Fatal("endpoint is required")
	}

	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":9000", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		req, err := http.NewRequest(r.Method, *endpoint, r.Body)
		if err != nil {
			processErr(rw, err)
			return
		}

		reqData, err := httputil.DumpRequest(req, true)
		if err != nil {
			processErr(rw, err)
			return
		}

		log.Println("Forward Request Data", string(reqData))

		req.Header = r.Header.Clone()

		client := http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			processErr(rw, err)
			return
		}
		defer resp.Body.Close()

		respData, err := httputil.DumpResponse(resp, true)
		if err != nil {

			return
		}

		log.Println("Forward Request Response", string(respData))

		rw.WriteHeader(resp.StatusCode)

		for k, v := range resp.Header {
			rw.Header()[k] = v
		}

		_, err = io.Copy(rw, resp.Body)
		if err != nil {
			log.Println(err)
			rw.Write([]byte("error"))
			return
		}

	})))

}

func processErr(rw http.ResponseWriter, err error) {
	log.Println(err)
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write([]byte("error"))
}
