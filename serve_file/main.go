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
	//Represents binary file
	w.Header().Set("Content-Type", "application/octet-stream")
	//Tells client what filename should be used.
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	//The length of the data.
	w.Header().Set("Content-Length", len)
	//No cache headers.
	w.Header().Set("Cache-Control", "private")
	//No cache headers.
	w.Header().Set("Pragma", "private")
	//No cache headers.
	w.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
}

func handleLargeFile(w http.ResponseWriter, r *http.Request) {
	//Open file
	f, err := os.Open(LargeFileName)
	if err != nil {
		handleError(err, w)
		return
	}
	defer f.Close()

	//read the file info
	info, err := f.Stat()
	if err != nil {
		handleError(err, w)
		return
	}

	//Set the headers
	setHeaders(w, LargeFileName, strconv.Itoa(int(info.Size())))
	w.WriteHeader(http.StatusOK)

	//Copy without loading everything in memory
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
