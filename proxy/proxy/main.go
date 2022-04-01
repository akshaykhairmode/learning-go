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

	//get the endpoint as a flag
	endpoint := flag.String("e", "", "")
	flag.Parse()

	if *endpoint == "" {
		log.Fatal("endpoint is required")
	}

	log.Println("starting server")

	//on every request call this handler
	log.Fatal(http.ListenAndServe(":9000", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		//create a new request
		req, err := http.NewRequest(r.Method, *endpoint, r.Body)
		if err != nil {
			processErr(rw, err)
			return
		}

		//Set the request headers as whatever was passed by caller
		req.Header = r.Header.Clone()

		//create a http client, timeout should be mentioned or it will never timeout.
		client := http.Client{
			Timeout: 5 * time.Second,
		}

		//Get dump of our request
		reqData, err := httputil.DumpRequest(req, true)
		if err != nil {
			processErr(rw, err)
			return
		}

		log.Println("Forward Request Data", string(reqData))

		//Actually forward the request to our endpoint
		resp, err := client.Do(req)
		if err != nil {
			processErr(rw, err)
			return
		}
		defer resp.Body.Close()

		//Get dump of our response
		respData, err := httputil.DumpResponse(resp, true)
		if err != nil {
			processErr(rw, err)
			return
		}

		log.Println("Forward Request Response", string(respData))

		//set the statuscode whatever we got from the response
		rw.WriteHeader(resp.StatusCode)

		//Copy the response headers to the actual response
		for k, v := range resp.Header {
			rw.Header()[k] = v
		}

		//Copy the response body to the actual response
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
