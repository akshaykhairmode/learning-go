package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	readBuffSize = 2 << 10
	writeBuffSize
)

const (
	nameHeader = "WS-NAME"
	idHeader   = "WS-ID"
)

var port string

func main() {

	log.SetOutput(os.Stdout)

	flag.StringVar(&port, "p", "9000", "port")
	flag.Parse()

	r := mux.NewRouter()

	chatroom := r.PathPrefix("/chatroom").Subrouter()
	chatroom.HandleFunc("/create/{name}", CR(CreateChatroom)).Methods(http.MethodPost)
	chatroom.HandleFunc("/connect", clientMW(chatroomWSHandler))

	client := r.PathPrefix("/client").Subrouter()
	client.HandleFunc("/list", ListAllClients).Methods(http.MethodGet)
	// client.HandleFunc("/connect")

	log.Println("Registered Handlers")

	// chat := r.PathPrefix("/chat").Subrouter()
	log.Printf("Started Server on port : %v", port)
	http.ListenAndServe(":"+port, r)
}
