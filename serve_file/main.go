package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var InternalServerError = []byte("Internal Server Error")
var NotFound = []byte("Page Not Found")

var smallFile []byte

const (
	LargeFileName = "large.txt"
	SmallFileName = "small.txt"
)

//seq 1 10000 | xargs -n 1 -P 10000 -I {} wget -P files --content-disposition http://localhost:9000/download-large-file >> cmd_out.log 2>&1
//seq 1 10000 | xargs -n 1 -P 10000 -I {} curl -o /dev/null http://localhost:9000/download-large-file

func main() {

	f, err := os.Open(SmallFileName)
	if err != nil {
		panic(err)
	}

	smallFile, err = io.ReadAll(f)
	if err != nil {
		panic(err)
	}

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
	case "/download-small-file":
		handleSmallFile(w, r)
	case "/download-large-file":
		handleLargeFile(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write(NotFound)
	}

}

func setHeaders(w http.ResponseWriter, name, len string) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	w.Header().Set("Content-Length", len)
	w.Header().Set("Cache-Control", "private")
	w.Header().Set("Pragma", "private")
	w.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
}

func handleLargeFile(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(LargeFileName)
	if err != nil {
		handleError(err, w)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		handleError(err, w)
		return
	}

	setHeaders(w, LargeFileName, strconv.Itoa(int(info.Size())))
	w.WriteHeader(http.StatusOK)
	n, err := io.Copy(w, f)
	if err != nil {
		handleError(err, w)
		return
	}
	log.Printf("LARGE :: Written : %d", n)
}

func handleSmallFile(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, SmallFileName, strconv.Itoa(len(smallFile)))
	n, err := w.Write(smallFile)
	if err != nil {
		handleError(err, w)
		return
	}
	log.Printf("SMALL :: Written : %d", n)
}

func handleError(err error, w http.ResponseWriter) {
	log.Printf("error while opening file : %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(InternalServerError)
}
