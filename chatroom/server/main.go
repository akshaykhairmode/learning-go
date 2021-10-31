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
	chatroom.HandleFunc("/delete/{id}", CR(DeleteChatroom)).Methods(http.MethodDelete)
	// chatroom.HandleFunc("/add", AddUserToChatroom).Methods(http.MethodPut)
	chatroom.HandleFunc("/connect/{id}", WS(handleChatroomWS))

	// user := r.PathPrefix("/user").Subrouter()
	// user.HandleFunc("/connect/{id}")

	log.Println("Registered Handlers")

	// chat := r.PathPrefix("/chat").Subrouter()
	log.Println("Started Server")
	http.ListenAndServe(":"+port, r)
}
