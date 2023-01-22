package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var NotFound = []byte("Page Not Found")

func main() {

	server := http.Server{
		Handler: http.HandlerFunc(handle),
		Addr:    ":9000",
	}

	log.Println("starting server")
	log.Println(server.ListenAndServe())

}

func handle(w http.ResponseWriter, r *http.Request) {

	log.Println("path", r.URL.Path)

	switch r.URL.Path {
	case "/receive-file":
		upload(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write(NotFound)
	}

}

func upload(w http.ResponseWriter, r *http.Request) {

	//If its not multipart, We will expect file data in body.
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		handleFileInBody(w, r)
		return
	}

	handleFileInForm(w, r)

}

func handleFileInBody(w http.ResponseWriter, r *http.Request) {

	//Check if body if empty or not
	if r.ContentLength <= 0 {
		w.WriteHeader(400)
		w.Write([]byte("Got content length <= 0"))
		return
	}

	f, err := getFile("")
	if err != nil {
		somethingWentWrong(w)
		return
	}
	defer f.Close()

	written, err := io.Copy(f, r.Body)
	if err != nil {
		log.Println("copy error", err)
		somethingWentWrong(w)
		return
	}

	success(w)

	log.Println("Written", written)
}

func handleFileInForm(w http.ResponseWriter, r *http.Request) {
	f, fh, err := r.FormFile("file")
	if err != nil {
		log.Println("formfile error", err)
		somethingWentWrong(w)
		return
	}

	if fh.Size <= 0 {
		w.WriteHeader(400)
		w.Write([]byte("Got File length <= 0"))
		return
	}

	outFile, err := getFile(fh.Filename)
	if err != nil {
		log.Println("getFile error", err)
		somethingWentWrong(w)
		return
	}

	written, err := io.Copy(outFile, f)
	if err != nil {
		log.Println("copy error", err)
		somethingWentWrong(w)
		return
	}

	success(w)
	log.Println("Written", written)
}

func getFile(fname string) (*os.File, error) {
	var fileName string

	now := time.Now()
	if fname != "" {
		fileName = strconv.Itoa(int(now.Unix())) + "_" + fname
	} else {
		fileName = "temp_" + strconv.Itoa(int(now.Unix())) + ".txt"
	}

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("create file error", err)
		return nil, err
	}

	return f, nil

}

func success(w http.ResponseWriter) {
	w.WriteHeader(200)
}

func somethingWentWrong(w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte("something went wrong"))
}
